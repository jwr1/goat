package goatw

import (
	"fmt"

	. "github.com/jwr1/goat"
)

type SizedBox struct {
	Widget

	Width  int
	Height int
}

var _ RenderWidget = SizedBox{}

func (w SizedBox) Layout(context LayoutContext) (Size, error) {
	return SizeInt(w.Width, w.Height), nil
}

func (w SizedBox) Paint(context PaintContext) error {

	return nil
}

type Padding struct {
	Widget

	Child   Widget
	Padding EdgeInserts
}

var _ RenderWidget = Padding{}

func (w Padding) Layout(context LayoutContext) (Size, error) {
	childConstrains := Constraints{
		Min: context.Constraints.Min.SubEdgeInserts(w.Padding),
		Max: context.Constraints.Max.SubEdgeInserts(w.Padding),
	}
	if childConstrains.Min.HasNeg() {
		childConstrains.Min = SizeZero
	}
	if childConstrains.Max.HasNeg() {
		return Size{}, fmt.Errorf("not enough space for padding given constraints")
	}

	childSize, err := context.LayoutChild(0, w.Child, childConstrains)
	if err != nil {
		return Size{}, err
	}

	err = context.PositionChild(0, Pos{
		X: w.Padding.Left,
		Y: w.Padding.Top,
	})
	if err != nil {
		return Size{}, err
	}

	return childSize.AddEdgeInserts(w.Padding), nil
}

func (w Padding) Paint(context PaintContext) error {
	return nil
}

type Background struct {
	Widget

	Child      Widget
	Background Color
}

var _ RenderWidget = Background{}

func (w Background) Layout(context LayoutContext) (Size, error) {
	size, err := context.LayoutChild(0, w.Child, context.Constraints)
	if err != nil {
		return Size{}, err
	}
	err = context.PositionChild(0, Pos{})
	if err != nil {
		return Size{}, err
	}
	return size, nil
}

func (w Background) Paint(context PaintContext) error {
	context.Canvas.FillBackground(0, 0, context.Size.Width.Int(), context.Size.Height.Int(), w.Background)
	return nil
}
