package goat

import (
	"fmt"
)

type Equality[T any] interface {
	Equal(T) bool
}

type Pos struct {
	X int
	Y int
}

func (p Pos) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

func (pos Pos) Add(other Pos) Pos {
	return Pos{
		X: pos.X + other.X,
		Y: pos.Y + other.Y,
	}
}

func (pos Pos) Sub(other Pos) Pos {
	return Pos{
		X: pos.X - other.X,
		Y: pos.Y - other.Y,
	}
}

type RenderViewport struct {
	AbsoluteStart Pos
	AbsoluteEnd   Pos
	LocalStart    Pos
	LocalEnd      Pos
}

type Size struct {
	Width  int
	Height int
}

var SizeZero = Size{}

func SizeSquare(value int) Size {
	return Size{
		Width:  value,
		Height: value,
	}
}

func (s Size) Add(other Size) Size {
	return Size{
		Width:  s.Width + other.Width,
		Height: s.Height + other.Height,
	}
}

func (s Size) Sub(other Size) Size {
	return Size{
		Width:  s.Width - other.Width,
		Height: s.Height - other.Height,
	}
}

func (s Size) String() string {
	return fmt.Sprintf("(%d,%d)", s.Width, s.Height)
}

func (s Size) IsNegative() bool {
	return s.Width < 0 || s.Height < 0
}

func (s Size) Clamp(constraints Constraints) Size {
	return Size{
		Width:  min(max(constraints.Min.Width, s.Width), constraints.Max.Width),
		Height: min(max(constraints.Min.Height, s.Height), constraints.Max.Height),
	}
}

// Creates constraints that forbid sizes larger than this size
func (s Size) LooseConstraints() Constraints {
	return Constraints{Max: s}
}

// Creates constraints that is respected only by this size
func (s Size) TightConstraints() Constraints {
	return Constraints{Min: s, Max: s}
}

// Creates a Size with the original size plus edge inserts
func (s Size) AddEdgeInserts(edgeInserts EdgeInserts) Size {
	return Size{
		Width:  s.Width + (edgeInserts.Left + edgeInserts.Right),
		Height: s.Height + (edgeInserts.Top + edgeInserts.Bottom),
	}
}

// Creates a Size with the original size plus edge inserts
func (s Size) SubEdgeInserts(edgeInserts EdgeInserts) Size {
	return Size{
		Width:  s.Width - (edgeInserts.Left + edgeInserts.Right),
		Height: s.Height - (edgeInserts.Top + edgeInserts.Bottom),
	}
}

type Constraints struct {
	Min Size
	Max Size
}

func (c Constraints) Check(size Size) bool {
	return c.Min.Width <= size.Width && size.Width <= c.Max.Width &&
		c.Min.Height <= size.Height && size.Height <= c.Max.Height
}

type ConstraintViolationErr struct {
	constraints Constraints
	realSize    Size
	widgetName  string
}

func (err *ConstraintViolationErr) Error() string {
	return fmt.Sprintf(
		"constraints violated for widget %s: Min: %s, Max: %s, Given size: %s",
		err.widgetName,
		err.constraints.Min.String(),
		err.constraints.Max.String(),
		err.realSize.String(),
	)
}

type EdgeInserts struct {
	Top    int
	Left   int
	Right  int
	Bottom int
}

func (ei EdgeInserts) String() string {
	return fmt.Sprintf("(%d,%d,%d,%d)", ei.Top, ei.Left, ei.Right, ei.Bottom)
}

func EdgeInsertsAll(value int) EdgeInserts {
	return EdgeInserts{
		Top:    value,
		Left:   value,
		Right:  value,
		Bottom: value,
	}
}

func EdgeInsertsSymmetric(vertical int, horizontal int) EdgeInserts {
	return EdgeInserts{
		Top:    vertical,
		Left:   horizontal,
		Right:  horizontal,
		Bottom: vertical,
	}
}
