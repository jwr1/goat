package main

import (
	"fmt"
	"time"

	goatw "github.com/jwr1/goat/widget"

	"github.com/jwr1/goat"

	"github.com/gdamore/tcell/v2"
)

type app struct {
	goat.Widget
}

var _ goat.StateWidget = app{}

func (w app) Build() (goat.Widget, error) {
	value, setValue := goat.UseStateFunc(func() uint8 { return 0xff / 2 })

	buttonPad := goat.EdgeInsertsSymmetric(0, 1)

	goat.UseEffect(func() func() {
		ticker := time.NewTicker(time.Millisecond * 10)
		done := make(chan bool)

		directionUp := true

		go func() {
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					setValue(func(value uint8) uint8 {
						if directionUp {
							if value == 0xff {
								directionUp = false
							}
						} else {
							if value == 0x00 {
								directionUp = true
							}
						}

						if directionUp {
							return value + 1
						} else {
							return value - 1
						}
					})
				}
			}
		}()

		return func() {
			ticker.Stop()
			done <- true
		}
	}, []any{})

	return goatw.Center{
		Child: goatw.Column{
			MainAxisShrinkWrap: true,
			CrossAxisAlignment: goatw.CrossAxisAlignmentCenter,
			Children: []goat.Widget{
				goatw.Text{
					Text: fmt.Sprint("Opacity: ", value),
				},
				goatw.Row{
					MainAxisShrinkWrap: true,
					Children: []goat.Widget{
						goatw.Button{
							Label:   "[↑] Increment",
							Padding: buttonPad,
							OnActivate: func() {
								setValue(func(value uint8) uint8 { return value + 1 })
							},
							KeyActivators: []string{tcell.KeyNames[tcell.KeyUp]},
						},
						goatw.Button{
							Label:   "[↓] Decrement",
							Padding: buttonPad,
							OnActivate: func() {
								setValue(func(value uint8) uint8 { return value - 1 })
							},
							KeyActivators: []string{tcell.KeyNames[tcell.KeyDown]},
						},
					},
				},
				goatw.Row{
					MainAxisShrinkWrap: true,
					Children: []goat.Widget{
						goatw.Button{
							Label:   "[Z] Set to Zero",
							Padding: buttonPad,
							OnActivate: func() {
								setValue(func(value uint8) uint8 { return 0 })
							},
							KeyActivators: []string{"z"},
						},
						goatw.Button{
							Label:   "[H] Set to Half",
							Padding: buttonPad,
							OnActivate: func() {
								setValue(func(value uint8) uint8 { return 0xff / 2 })
							},
							KeyActivators: []string{"h"},
						},
						goatw.Button{
							Label:   "[F] Set to Full",
							Padding: buttonPad,
							OnActivate: func() {
								setValue(func(value uint8) uint8 { return 0xff })
							},
							KeyActivators: []string{"f"},
						},
					},
				},
				goatw.Background{
					Background: goat.ColorRGB(0xff, 0, 0),
					Child: goatw.Padding{
						Padding: goat.EdgeInsertsSymmetric(1, 2),
						Child: goatw.Background{
							Background: goat.Color{R: 0, G: 0, B: 0xff, A: uint8(value)},
							Child: goatw.Padding{
								Padding: goat.EdgeInsertsSymmetric(1, 2),
								Child: goatw.Text{
									Text: "ABCD",
								},
							},
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
