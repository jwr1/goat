package goat

import (
	"github.com/gdamore/tcell/v2"
)

type hookContext struct {
	element    *element
	stateIndex int
	eventFuncs []func(context EventContext)
	effects    []effect
}

var curHookContext hookContext

var globalHookEventListeners = make(map[*element][]func(context EventContext))

func setupHooks(e *element) {
	curHookContext = hookContext{element: e}
}

func resetHooks() {
	globalHookEventListeners[curHookContext.element] = curHookContext.eventFuncs
	curHookContext = hookContext{}
}

func verifyHooks() {
	if curHookContext.element == nil {
		panic("Hook was used outside of widget")
	}
}

func UseState[T comparable](defaultValue T) (T, func(T)) {
	state, setState := UseRawState(defaultValue)

	return state, func(newValue T) {
		if newValue != state {
			// TreeLock.Lock()
			// defer TreeLock.Unlock()

			setState(newValue)
		}
	}
}

// Like UseState, but will not check if a new value is equal to the old and will always trigger a rerender when set.
// It's recommended to only use this if your state contains incomparable values.
func UseRawState[T any](defaultValue T) (T, func(T)) {
	verifyHooks()

	curIndex := curHookContext.stateIndex
	curElement := curHookContext.element
	curHookContext.stateIndex += 1

	// State does not already exist
	if curIndex >= len(curElement.states) {
		curElement.states = append(curElement.states, defaultValue)
	}

	return curElement.states[curIndex].(T), func(newValue T) {
		// TreeLock.Lock()
		// defer TreeLock.Unlock()

		curElement.states[curIndex] = newValue
		curElement.queueBuild = true
		curElement.queueRender = true
	}
}

type EventContext struct {
	Event      tcell.Event
	RenderPos  Pos
	RenderSize Size
}

func UseEvent(fn func(context EventContext)) {
	verifyHooks()

	curHookContext.eventFuncs = append(curHookContext.eventFuncs, fn)
}

// A hook that lets you synchronize a widget with an external system.
// The setup function will be run when the widget is first mounted, and also whenever the effect's dependencies change.
// Setup can optionally return a cleanup function, which will be run when the widget is unmounted and also before a dependency change setup is triggered.
func UseEffect(setup func() func(), dependencies []any) {
	verifyHooks()

	curHookContext.effects = append(curHookContext.effects, effect{
		setup:        setup,
		dependencies: dependencies,
	})
}
