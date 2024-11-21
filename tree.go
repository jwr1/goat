package goat

import (
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func rebuildTree(newWidget Widget, thisElement *Element, constraints Constraints) error {
	if thisElement.queueBuild || !thisElement.isInitialized {
		goto build
	}

	// If new widget is different type, then recreate whole widget tree.
	{
		oldWidgetType := reflect.TypeOf(thisElement.widget)
		newWidgetType := reflect.TypeOf(newWidget)

		if oldWidgetType != newWidgetType {
			destroyTree(thisElement)
			err := rebuildTree(newWidget, thisElement, constraints)
			if err != nil {
				return err
			}
			return nil
		}
	}

	// If widget props have changed, then perform a build.
	if !reflect.DeepEqual(thisElement.widget, newWidget) {
		thisElement.queuePaint = true
		goto build
	}

	// If widget constraints have changed, then perform a build.
	if constraints != thisElement.prevConstraints {
		goto build
	}

	// Even if this element doesn't need to be rebuilt, all descendants still need to be checked
	for _, childElement := range thisElement.children {
		err := rebuildTree(childElement.widget, childElement, childElement.prevConstraints)
		if err != nil {
			return err
		}
	}
	return nil

build:
	thisElement.queueBuild = false
	thisElement.widget = newWidget
	thisElement.prevConstraints = constraints

	switch newWidget := newWidget.(type) {
	case StateWidget:
		setupHooks(thisElement)
		childWidget, err := newWidget.Build()
		newEffects := currentHookContext.effects
		resetHooks()
		if err != nil {
			return err
		}

		if !thisElement.isInitialized {
			thisElement.effects = newEffects
			for i := 0; i < len(thisElement.effects); i++ {
				thisElement.effects[i].cleanup = thisElement.effects[i].setup()
			}
		} else {
			if len(thisElement.effects) != len(newEffects) {
				return fmt.Errorf("different number of effects found in widget %s, last build had %d, this build has %d", reflect.TypeOf(newWidget).String(), len(thisElement.effects), len(newEffects))
			}

			for i := 0; i < len(thisElement.effects); i++ {
				if slices.Equal(thisElement.effects[i].dependencies, newEffects[i].dependencies) {
					continue
				}

				if thisElement.effects[i].cleanup != nil {
					thisElement.effects[i].cleanup()
				}

				thisElement.effects[i] = newEffects[i]
				thisElement.effects[i].cleanup = thisElement.effects[i].setup()
			}
		}

		childElement, ok := thisElement.children[0]
		if !ok {
			childElement = &Element{}
			thisElement.children = map[int]*Element{
				0: childElement,
			}
		}
		err = rebuildTree(childWidget, childElement, constraints)
		if err != nil {
			return err
		}
		thisElement.size = childElement.size

	case RenderWidget:
		oldChildren := thisElement.children
		newChildren := make(map[int]*Element)

		layoutContext := LayoutContext{
			Constraints: constraints,
			LayoutChild: func(key int, c Widget, constraints Constraints) (Size, error) {
				childElement, ok := oldChildren[key]
				if !ok {
					childElement = &Element{parent: thisElement}
				} else {
					delete(oldChildren, key)
				}
				err := rebuildTree(c, childElement, constraints)
				if err != nil {
					return Size{}, err
				}

				newChildren[key] = childElement
				return childElement.size, nil
			},
			PositionChild: func(key int, pos Pos) error {
				childElement, ok := newChildren[key]
				if !ok {
					return fmt.Errorf("LayoutChild() must be called before PositionChild()")
				}

				childElement.pos = pos

				return nil
			},
		}

		size, err := newWidget.Layout(layoutContext)
		if err != nil {
			return err
		}
		if !constraints.Check(size) {
			return &ConstraintViolationErr{
				constraints: constraints,
				realSize:    size,
				widgetName:  reflect.TypeOf(newWidget).String(),
			}
		}

		// Destroy any children that have not been layed out in this build
		for _, e := range oldChildren {
			destroyTree(e)
		}

		thisElement.children = newChildren

		// if size != thisElement.size {
		thisElement.queuePaint = true
		// }

		thisElement.size = size
	default:
		panic("widget not implemented")
	}

	thisElement.isInitialized = true

	return nil
}

func renderTree(thisElement *Element) (Canvas, error) {
	var resultCanvas Canvas

	switch widget := thisElement.widget.(type) {
	case StateWidget:
		childElement := thisElement.children[0]
		canvas, err := renderTree(childElement)
		if err != nil {
			return Canvas{}, err
		}
		resultCanvas = canvas

	case RenderWidget:
		if thisElement.queuePaint {
			thisElement.queuePaint = false
			renderContext := PaintContext{
				Canvas: NewCanvas(thisElement.size),
				Size:   thisElement.size,
			}

			err := widget.Paint(renderContext)
			if err != nil {
				return Canvas{}, err
			}

			thisElement.renderCanvas = renderContext.Canvas
		}

		resultCanvas = NewCanvas(thisElement.size)
		resultCanvas.OverlayCanvas(0, 0, thisElement.renderCanvas)

	default:
		panic("widget not implemented")
	}

	for _, childElement := range thisElement.children {
		childElement.renderAbsPos = thisElement.renderAbsPos.Add(childElement.pos)
		childCanvas, err := renderTree(childElement)
		if err != nil {
			return Canvas{}, err
		}
		resultCanvas.OverlayCanvas(childElement.pos.X, childElement.pos.Y, childCanvas)
	}

	return resultCanvas, nil
}

func destroyTree(thisElement *Element) {
	for _, childElement := range thisElement.children {
		destroyTree(childElement)
	}

	for i := 0; i < len(thisElement.effects); i++ {
		if thisElement.effects[i].cleanup != nil {
			thisElement.effects[i].cleanup()
		}
	}

	*thisElement = Element{}
}

func stringifyTree(thisElement *Element, builder *strings.Builder, indent int) {
	strIndent := strings.Repeat(" ", indent)

	builder.WriteString(strIndent)
	builder.WriteString("widget: ")
	builder.WriteString(reflect.TypeOf(thisElement.widget).String())
	builder.WriteString("\n")
	builder.WriteString(strIndent)
	builder.WriteString("size: ")
	builder.WriteString(thisElement.size.String())
	builder.WriteString("\n")
	builder.WriteString(strIndent)
	builder.WriteString(fmt.Sprintf(
		"constraints: Min: %s, Max: %s",
		thisElement.prevConstraints.Min.String(),
		thisElement.prevConstraints.Max.String(),
	))
	builder.WriteString("\n")

	builder.WriteString(strIndent)
	builder.WriteString("props:\n")

	widgetType := reflect.TypeOf(thisElement.widget)
	propsIndent := strings.Repeat(" ", indent+2)

	for i := 0; i < widgetType.NumField(); i++ {
		field := widgetType.Field(i)

		fieldName := field.Name
		fieldType := field.Type

		if fieldType.String() == "goat.Widget" || fieldType.String() == "[]goat.Widget" {
			continue
		}

		fieldValue := reflect.ValueOf(thisElement.widget).Field(i).Interface()

		builder.WriteString(propsIndent)
		builder.WriteString(fmt.Sprintf("%s %s: %v\n", fieldName, fieldType, fieldValue))
	}

	childKeys := make([]int, 0, len(thisElement.children))
	for k := range thisElement.children {
		childKeys = append(childKeys, k)
	}

	sort.Ints(childKeys)

	for _, key := range childKeys {
		builder.WriteString(strIndent)
		builder.WriteString("child ")
		builder.WriteString(strconv.Itoa(key))
		builder.WriteString(":\n")
		stringifyTree(thisElement.children[key], builder, indent+2)
	}
}
