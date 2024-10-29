package main

import (
	"strconv"

	goatw "github.com/jwr1/goat/widget"

	"github.com/jwr1/goat"

	"github.com/gdamore/tcell/v2"
)

type app struct {
	goat.Widget
}

var _ goat.StateWidget = app{}

func (w app) Build() (goat.Widget, error) {
	value, setValue := goat.UseState(0)

	buttonPad := goat.EdgeInsertsSymmetric(0, 2)

	return goatw.Center{
		Child: goatw.Column{
			MainAxisShrinkWrap: true,
			CrossAxisAlignment: goatw.CrossAxisAlignmentCenter,
			Children: []goat.Widget{
				goatw.Text{
					Text: "The counter is currently at:",
				},
				goatw.Padding{
					Padding: goat.EdgeInsertsAll(1),
					Child: goatw.Text{
						Text: strconv.Itoa(value),
					},
				},
				goatw.Row{
					MainAxisShrinkWrap: true,
					Children: []goat.Widget{
						goatw.Button{
							Label:   "[R] Reset",
							Padding: buttonPad,
							OnActivate: func() {
								setValue(0)
							},
							KeyActivators: []string{"r"},
						},
						goatw.SizedBox{Width: 1},
						goatw.Button{
							Label:   "[↑] Increment",
							Padding: buttonPad,
							OnActivate: func() {
								setValue(value + 1)
							},
							KeyActivators: []string{tcell.KeyNames[tcell.KeyUp]},
						},
						goatw.SizedBox{Width: 1},
						goatw.Button{
							Label:   "[↓] Decrement",
							Padding: buttonPad,
							OnActivate: func() {
								setValue(value - 1)
							},
							KeyActivators: []string{tcell.KeyNames[tcell.KeyDown]},
						},
					},
				},
			},
		},
	}, nil
}

func main() {
	err := goat.RunApp(app{})
	if err != nil {
		panic(err.Error())
	}
}
