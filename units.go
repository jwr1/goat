package goat

import (
	"fmt"
	"math"
	"strconv"
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

type Dimension struct {
	value int
}

var (
	DimensionZero     = Dimension{0}
	DimensionInfinite = Dimension{math.MaxInt}
)

// Creates a new Dimension from an integer
func DimensionInt(value int) Dimension {
	return Dimension{value}
}

// Returns the underlying integer value
func (d Dimension) Int() int {
	return d.value
}

// Returns the string representation, either it's integer value, or "Inf"
func (d Dimension) String() string {
	if d.IsInf() {
		return "Inf"
	}

	return strconv.Itoa(d.Int())
}

// Reports whether the Dimension is less than zero
func (d Dimension) IsNeg() bool {
	return d.value < DimensionZero.value
}

// Reports whether the Dimension is equal to zero
func (d Dimension) IsZero() bool {
	return d == DimensionZero
}

// Reports whether the Dimension is infinite
func (d Dimension) IsInf() bool {
	return d == DimensionInfinite
}

func (d Dimension) Add(other Dimension) Dimension {
	if d.IsInf() || other.IsInf() {
		return DimensionInfinite
	}

	return Dimension{d.value + other.value}
}
func (d Dimension) AddInt(other int) Dimension {
	if d.IsInf() {
		return DimensionInfinite
	}

	return Dimension{d.value + other}
}
func (d Dimension) Sub(other Dimension) Dimension {
	if d.IsInf() || other.IsInf() {
		return DimensionInfinite
	}

	return Dimension{d.value - other.value}
}
func (d Dimension) SubInt(other int) Dimension {
	if d.IsInf() {
		return DimensionInfinite
	}

	return Dimension{d.value - other}
}

type Size struct {
	Width  Dimension
	Height Dimension
}

var (
	SizeZero = Size{
		Width:  DimensionZero,
		Height: DimensionZero,
	}
	SizeInfinite = Size{
		Width:  DimensionInfinite,
		Height: DimensionInfinite,
	}
)

func SizeInt(width, height int) Size {
	return Size{
		Width:  DimensionInt(width),
		Height: DimensionInt(height),
	}
}

func SizeSquare(value int) Size {
	return Size{
		Width:  DimensionInt(value),
		Height: DimensionInt(value),
	}
}

func SizeSquareVisual(value int) Size {
	return Size{
		Width:  DimensionInt(value * 2),
		Height: DimensionInt(value),
	}
}

func (s Size) Add(other Size) Size {
	return Size{
		Width:  s.Width.Add(other.Width),
		Height: s.Height.Add(other.Height),
	}
}

func (s Size) Sub(other Size) Size {
	return Size{
		Width:  s.Width.Sub(other.Width),
		Height: s.Height.Sub(other.Height),
	}
}

func (s Size) String() string {
	return fmt.Sprintf("(%s,%s)", s.Width.String(), s.Height.String())
}

func (s Size) HasNeg() bool {
	return s.Width.IsNeg() || s.Height.IsNeg()
}
func (s Size) HasZero() bool {
	return s.Width.IsZero() || s.Height.IsZero()
}
func (s Size) HasInf() bool {
	return s.Width.IsInf() || s.Height.IsInf()
}

func (s Size) Clamp(constraints Constraints) Size {
	return Size{
		Width:  DimensionInt(min(max(constraints.Min.Width.Int(), s.Width.Int()), constraints.Max.Width.Int())),
		Height: DimensionInt(min(max(constraints.Min.Height.Int(), s.Height.Int()), constraints.Max.Height.Int())),
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
		Width:  s.Width.AddInt(edgeInserts.Left + edgeInserts.Right),
		Height: s.Height.AddInt(edgeInserts.Top + edgeInserts.Bottom),
	}
}

// Creates a Size with the original size plus edge inserts
func (s Size) SubEdgeInserts(edgeInserts EdgeInserts) Size {
	return Size{
		Width:  s.Width.SubInt(edgeInserts.Left + edgeInserts.Right),
		Height: s.Height.SubInt(edgeInserts.Top + edgeInserts.Bottom),
	}
}

type Constraints struct {
	Min Size
	Max Size
}

func (c Constraints) Check(size Size) bool {
	return c.Min.Width.Int() <= size.Width.Int() && size.Width.Int() <= c.Max.Width.Int() &&
		c.Min.Height.Int() <= size.Height.Int() && size.Height.Int() <= c.Max.Height.Int()
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
