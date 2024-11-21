package goat

import (
	"github.com/gdamore/tcell/v2"
)

type hookContext struct {
	element    *Element
	refIndex   int
	effects    []effect
	eventFuncs []func(context EventContext)
}

var currentHookContext hookContext

var globalHookEventListeners = make(map[*Element][]func(context EventContext))

func setupHooks(e *Element) {
	currentHookContext = hookContext{element: e}
}

func resetHooks() {
	globalHookEventListeners[currentHookContext.element] = currentHookContext.eventFuncs
	currentHookContext = hookContext{}
}

func getHookContext() *hookContext {
	if currentHookContext.element == nil {
		panic("Hook was used outside of widget")
	}

	return &currentHookContext
}

// A hook used to queue a widget render. This should almost never be needed for standard use cases; try using one of the UseState hooks instead.
func UseTriggerRender() func() {
	context := getHookContext()
	curElement := context.element

	return func() {
		curElement.queueBuild = true
		curElement.queuePaint = true
	}
}

// The most primitive hook to associate persistent data with a widget. A getter function is passed in that will be run once to retrieve the initial value and a reference to the variable is returned by the hook.
//
// If you need to pass in the initial value directly, use UseRef instead. If the value you're using is needed for rendering (which is the case the majority of the time), then use one of the UseState hooks.
func UseRefFunc[T any](initialValue func() *T) *T {
	context := getHookContext()

	curIndex := context.refIndex
	curElement := context.element
	context.refIndex += 1

	// Ref does not already exist
	if curIndex >= len(curElement.refs) {
		curElement.refs = append(curElement.refs, initialValue())
	}

	return curElement.refs[curIndex].(*T)
}

type EventContext struct {
	Event      tcell.Event
	RenderPos  Pos
	RenderSize Size
}

func UseEvent(fn func(context EventContext)) {
	context := getHookContext()

	context.eventFuncs = append(context.eventFuncs, fn)
}

// A hook that lets you synchronize a widget with an external system.
// The setup function will be run when the widget is first mounted, and also whenever the effect's dependencies change.
// Setup can optionally return a cleanup function, which will be run when the widget is unmounted and also before a dependency change setup is triggered.
func UseEffect(setup func() func(), dependencies []any) {
	context := getHookContext()

	context.effects = append(context.effects, effect{
		setup:        setup,
		dependencies: dependencies,
	})
}
