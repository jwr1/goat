package goat

func UseState[T comparable](initialValue T) (T, func(T)) {
	value, setValue := UseStateFunc(func() T { return initialValue })

	return value, func(newValue T) {
		setValue(func(T) T { return newValue })
	}
}

func UseStateFunc[T comparable](initialValue func() T) (T, func(func(T) T)) {
	triggerRender := UseTriggerRender()

	ref := UseRefFunc[T](func() *T {
		value := initialValue()
		return &value
	})

	return *ref, func(setter func(T) T) {
		newValue := setter(*ref)

		// Don't rerender if the old and new state values are the same
		if newValue == *ref {
			return
		}

		*ref = newValue

		triggerRender()
	}
}

// Like UseState, but will not check if a new value is equal to the old and will always trigger a rerender when set.
// Only use this if your state contains incomparable values.
func UseRawState[T any](initialValue T) (T, func(T)) {
	value, setValue := UseRawStateFunc(func() T { return initialValue })

	return value, func(newValue T) {
		setValue(func(T) T { return newValue })
	}
}

// Like UseRawState, but uses functions to retrieve the initialization and setter values.
func UseRawStateFunc[T any](initialValue func() T) (T, func(func(T) T)) {
	triggerRender := UseTriggerRender()

	ref := UseRefFunc[T](func() *T {
		value := initialValue()
		return &value
	})

	return *ref, func(setter func(T) T) {
		*ref = setter(*ref)
		triggerRender()
	}
}

// A hook that lets you reference a value that’s not needed for rendering.
// If the value you're using is needed for rendering (which is the case the majority of the time), then use one of the UseState hooks.
func UseRef[T any](initialValue *T) *T {
	return UseRefFunc(func() *T { return initialValue })
}

// A hook that runs a setup function when the widget is first mounted
func UseSetup(setup func()) {
	UseEffect(func() func() {
		setup()
		return nil
	}, []any{})
}

// A hook that runs a cleanup function when the widget is unmounted/destroyed.
func UseCleanup(cleanup func()) {
	UseEffect(func() func() {
		return cleanup
	}, []any{})
}
