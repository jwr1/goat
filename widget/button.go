package goatw

import (
	"unicode/utf8"

	. "github.com/jwr1/goat"

	"github.com/gdamore/tcell/v2"
)

type ButtonState int

const (
	ButtonStateIdle ButtonState = iota
	ButtonStateHover
	ButtonStateActive
)

type Button struct {
	Widget

	Label         string
	Padding       EdgeInserts
	OnActivate    func()
	KeyActivators []string
}

var _ StateWidget = Button{}

func (w Button) Build() (Widget, error) {
	buttonState, setButtonState := UseState(ButtonStateIdle)

	UseEvent(func(context EventContext) {
		switch event := context.Event.(type) {
		case *tcell.EventKey:
			if w.OnActivate == nil || len(w.KeyActivators) == 0 {
				return
			}

			if event.Key() == tcell.KeyRune {
				for _, keyActivator := range w.KeyActivators {
					r, _ := utf8.DecodeRuneInString(keyActivator)
					if r == event.Rune() {
						w.OnActivate()
					}
				}
			} else {
				keyName := tcell.KeyNames[event.Key()]
				for _, keyActivator := range w.KeyActivators {
					if keyActivator == keyName {
						w.OnActivate()
					}
				}
			}
		case *tcell.EventMouse:
			x, y := event.Position()
			intersects := context.RenderPos.X <= x && x < context.RenderPos.X+context.RenderSize.Width.Int() &&
				context.RenderPos.Y <= y && y < context.RenderPos.Y+context.RenderSize.Height.Int()
			if intersects {
				setButtonState(ButtonStateHover)

				if event.Buttons()&tcell.ButtonPrimary != 0 {
					setButtonState(ButtonStateActive)
				} else {
					setButtonState(ButtonStateHover)
					if buttonState == ButtonStateActive {
						if w.OnActivate != nil {
							w.OnActivate()
						}
					}
				}
			} else {
				setButtonState(ButtonStateIdle)
			}
		}
	})

	var bgColor Color
	switch buttonState {
	case ButtonStateIdle:
		bgColor = ColorRGB(100, 0, 0)
	case ButtonStateHover:
		bgColor = ColorRGB(200, 0, 0)
	case ButtonStateActive:
		bgColor = ColorRGB(255, 0, 0)
	}

	return Background{
		Background: bgColor,
		Child: Padding{
			Padding: w.Padding,
			Child:   Text{Text: w.Label},
		},
	}, nil
}
