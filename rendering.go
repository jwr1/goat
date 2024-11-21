package goat

type LayoutContext struct {
	Constraints   Constraints
	LayoutChild   func(key int, c Widget, constraints Constraints) (Size, error)
	PositionChild func(key int, pos Pos) error
	// PositionChildViewport func(key int, pos Pos, childStart Pos, childEnd Pos) error
}

type PaintContext struct {
	Canvas Canvas
	Size   Size
}

type Widget interface {
	isWidget()
}

type RenderWidget interface {
	Widget
	Layout(context LayoutContext) (Size, error)
	Paint(context PaintContext) error
}

type StateWidget interface {
	Widget
	Build() (Widget, error)
}

type effect struct {
	setup        func() func()
	cleanup      func()
	dependencies []any
}

type Element struct {
	isInitialized   bool
	widget          Widget
	size            Size
	pos             Pos
	renderCanvas    Canvas
	renderAbsPos    Pos
	prevConstraints Constraints

	queueBuild bool
	queuePaint bool

	parent   *Element
	children map[int]*Element

	refs    []any
	effects []effect

	renderData       any
	renderParentData any
}

func (e *Element) RenderData() any {
	return e.renderData
}

func (e *Element) SetRenderData(renderData any) {
	e.renderData = renderData
}

func (e *Element) RenderParentData() any {
	return e.renderData
}

func (e *Element) SetRenderParentData(renderParentData any) {
	e.renderParentData = renderParentData
}

func (e *Element) MarkNeedsBuild() {
	e.queueBuild = true
}

func (e *Element) MarkNeedsPaint() {
	e.queuePaint = true
}

func (e *Element) Parent() *Element {
	return e.parent
}

func (e *Element) Children() map[int]*Element {
	return e.children
}
