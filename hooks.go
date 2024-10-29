package goat

import (
	"github.com/gdamore/tcell/v2"
)

type hookContext struct {
	element        *element
	stateIndex     int
	eventListeners []func(context EventContext)
}

var curHookContext hookContext

var globalHookEventListeners = make(map[*element][]func(context EventContext))

func setupHooks(e *element) {
	curHookContext = hookContext{element: e}
}

func resetHooks() {
	globalHookEventListeners[curHookContext.element] = curHookContext.eventListeners
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
	if curIndex >= len(curElement.state) {
		curElement.state = append(curElement.state, defaultValue)
	}

	return curElement.state[curIndex].(T), func(newValue T) {
		// TreeLock.Lock()
		// defer TreeLock.Unlock()

		curElement.state[curIndex] = newValue
		curElement.queueBuild = true
		curElement.queueRender = true
	}
}

type EventContext struct {
	Event      tcell.Event
	RenderPos  Pos
	RenderSize Size
}

func UseEvent(cb func(context EventContext)) {
	verifyHooks()

	curHookContext.eventListeners = append(curHookContext.eventListeners, cb)
}
