package goatw

import (
	. "github.com/jwr1/goat"
)

type Axis int

const (
	AxisVertical Axis = iota
	AxisHorizontal
)

type MainAxisAlignment int

const (
	MainAxisAlignmentStart MainAxisAlignment = iota
	MainAxisAlignmentEnd
	MainAxisAlignmentCenter
	MainAxisAlignmentSpaceBetween
	MainAxisAlignmentSpaceAround
	MainAxisAlignmentSpaceEvenly
)

type CrossAxisAlignment int

const (
	CrossAxisAlignmentStart CrossAxisAlignment = iota
	CrossAxisAlignmentEnd
	CrossAxisAlignmentCenter
	CrossAxisAlignmentStretch
)

type Flex struct {
	Widget

	Children           []Widget
	Direction          Axis
	MainAxisAlignment  MainAxisAlignment
	MainAxisShrinkWrap bool
	CrossAxisAlignment CrossAxisAlignment
}

var _ RenderWidget = Flex{}

func (w Flex) Layout(context LayoutContext) (Size, error) {
	isHorizontal := w.Direction == AxisHorizontal

	mainAxisSize := func(s Size) int {
		if isHorizontal {
			return s.Width.Int()
		} else {
			return s.Height.Int()
		}
	}
	crossAxisSize := func(s Size) int {
		if isHorizontal {
			return s.Height.Int()
		} else {
			return s.Width.Int()
		}
	}
	sizeFromAxes := func(mainAxisSize, crossAxisSize int) Size {
		if isHorizontal {
			return SizeInt(mainAxisSize, crossAxisSize)
		} else {
			return SizeInt(crossAxisSize, mainAxisSize)
		}
	}
	positionChild := func(key, mainAxisPos, crossAxisPos int) error {
		if isHorizontal {
			return context.PositionChild(key, Pos{X: mainAxisPos, Y: crossAxisPos})
		} else {
			return context.PositionChild(key, Pos{X: crossAxisPos, Y: mainAxisPos})
		}
	}

	remainingSpace := mainAxisSize(context.Constraints.Max)
	childrenSizes := make([]Size, len(w.Children))
	childrenCrossAxisPos := make([]int, len(w.Children))
	finalCrossAxisSize := crossAxisSize(context.Constraints.Min)
	if w.CrossAxisAlignment == CrossAxisAlignmentStretch {
		finalCrossAxisSize = crossAxisSize(context.Constraints.Max)
	}

	for i, child := range w.Children {
		childMinSize := Size{}
		if w.CrossAxisAlignment == CrossAxisAlignmentStretch {
			childMinSize = sizeFromAxes(0, crossAxisSize(context.Constraints.Max))
		}

		childConstrains := Constraints{
			Min: childMinSize,
			Max: sizeFromAxes(remainingSpace, crossAxisSize(context.Constraints.Max)),
		}

		childSize, err := context.LayoutChild(i, child, childConstrains)
		if err != nil {
			return Size{}, err
		}

		childrenSizes[i] = childSize
		remainingSpace -= mainAxisSize(childSize)
		if crossAxisSize(childSize) > finalCrossAxisSize {
			finalCrossAxisSize = crossAxisSize(childSize)
		}
	}

	for i, child := range childrenSizes {
		switch w.CrossAxisAlignment {
		case CrossAxisAlignmentStart | CrossAxisAlignmentStretch:
			childrenCrossAxisPos[i] = 0
		case CrossAxisAlignmentCenter:
			childrenCrossAxisPos[i] = (finalCrossAxisSize - crossAxisSize(child)) / 2
		case CrossAxisAlignmentEnd:
			childrenCrossAxisPos[i] = (finalCrossAxisSize - crossAxisSize(child))
		}
	}

	finalMainAxisSize := mainAxisSize(context.Constraints.Max)
	if w.MainAxisShrinkWrap {
		minMainAxisSize := mainAxisSize(context.Constraints.Max) - remainingSpace
		finalMainAxisSize = max(minMainAxisSize, mainAxisSize(context.Constraints.Min))
		remainingSpace = finalMainAxisSize - minMainAxisSize
	}

	switch w.MainAxisAlignment {
	case MainAxisAlignmentStart:
		curPos := 0
		for i, size := range childrenSizes {
			positionChild(i,
				curPos,
				childrenCrossAxisPos[i],
			)
			curPos += mainAxisSize(size)
		}
	case MainAxisAlignmentCenter:
		curPos := remainingSpace / 2
		for i, size := range childrenSizes {
			positionChild(i,
				curPos,
				childrenCrossAxisPos[i],
			)
			curPos += mainAxisSize(size)
		}
	case MainAxisAlignmentEnd:
		curPos := remainingSpace
		for i, size := range childrenSizes {
			positionChild(i,
				curPos,
				childrenCrossAxisPos[i],
			)
			curPos += mainAxisSize(size)
		}
	case MainAxisAlignmentSpaceBetween:
		gap := remainingSpace / (len(w.Children) - 1)
		curPos := 0
		for i, size := range childrenSizes {
			positionChild(i,
				curPos,
				childrenCrossAxisPos[i],
			)
			curPos += mainAxisSize(size) + gap
		}
	case MainAxisAlignmentSpaceAround:
		gap := remainingSpace / len(w.Children)
		curPos := gap / 2
		for i, size := range childrenSizes {
			positionChild(i,
				curPos,
				childrenCrossAxisPos[i],
			)
			curPos += mainAxisSize(size) + gap
		}
	case MainAxisAlignmentSpaceEvenly:
		gap := remainingSpace / (len(w.Children) + 1)
		curPos := gap
		for i, size := range childrenSizes {
			positionChild(i,
				curPos,
				childrenCrossAxisPos[i],
			)
			curPos += mainAxisSize(size) + gap
		}
	}

	return sizeFromAxes(finalMainAxisSize, finalCrossAxisSize), nil
}

func (w Flex) Paint(context PaintContext) error {
	return nil
}

type Row struct {
	Widget

	Children           []Widget
	MainAxisAlignment  MainAxisAlignment
	MainAxisShrinkWrap bool
	CrossAxisAlignment CrossAxisAlignment
}

var _ StateWidget = Row{}

func (w Row) Build() (Widget, error) {
	return Flex{
		Direction:          AxisHorizontal,
		Children:           w.Children,
		MainAxisAlignment:  w.MainAxisAlignment,
		MainAxisShrinkWrap: w.MainAxisShrinkWrap,
		CrossAxisAlignment: w.CrossAxisAlignment,
	}, nil
}

type Column struct {
	Widget

	Children           []Widget
	MainAxisAlignment  MainAxisAlignment
	MainAxisShrinkWrap bool
	CrossAxisAlignment CrossAxisAlignment
}

var _ StateWidget = Column{}

func (w Column) Build() (Widget, error) {
	return Flex{
		Direction:          AxisVertical,
		Children:           w.Children,
		MainAxisAlignment:  w.MainAxisAlignment,
		MainAxisShrinkWrap: w.MainAxisShrinkWrap,
		CrossAxisAlignment: w.CrossAxisAlignment,
	}, nil
}

type Center struct {
	Widget

	Child        Widget
	WidthFactor  float64
	HeightFactor float64
}

var _ RenderWidget = Center{}

func (w Center) Layout(context LayoutContext) (Size, error) {
	childSize, err := context.LayoutChild(0, w.Child, Constraints{
		Max: context.Constraints.Max,
	})
	if err != nil {
		return Size{}, err
	}

	size := context.Constraints.Max
	if w.WidthFactor >= 1 {
		size.Width = DimensionInt(int(float64(childSize.Width.Int()) * w.WidthFactor))
	}
	if w.HeightFactor >= 1 {
		size.Height = DimensionInt(int(float64(childSize.Height.Int()) * w.HeightFactor))
	}
	size = size.Clamp(context.Constraints)

	remainingSize := size.Sub(childSize)

	context.PositionChild(0, Pos{
		X: remainingSize.Width.Int() / 2,
		Y: remainingSize.Height.Int() / 2,
	})

	return size, nil
}

func (w Center) Paint(context PaintContext) error {
	return nil
}
