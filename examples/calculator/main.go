package main

import (
	"strconv"

	goatw "github.com/jwr1/goat/widget"

	"github.com/jwr1/goat"

	"github.com/gdamore/tcell/v2"
)

var squareButtonPad = goat.EdgeInsertsSymmetric(1, 3)

type operation int

const (
	operationNone operation = iota + 0
	operationAdd
	operationSub
	operationMul
	operationDiv
)

var operationStr = map[operation]string{
	operationAdd: "+",
	operationSub: "-",
	operationMul: "*",
	operationDiv: "/",
}

type app struct {
	goat.Widget
}

var _ goat.StateWidget = app{}

func (w app) Build() (goat.Widget, error) {
	total, setTotal := goat.UseState(0)
	operationType, setOperationType := goat.UseState(operationNone)
	operationValue, setOperationValue := goat.UseState(0)

	applyOperation := func() {
		if operationValue == 0 {
			return
		}

		switch operationType {
		case operationNone:
			setTotal(operationValue)
		case operationAdd:
			setTotal(total + operationValue)
		case operationSub:
			setTotal(total - operationValue)
		case operationMul:
			setTotal(total * operationValue)
		case operationDiv:
			setTotal(total / operationValue)
		}
	}

	onOperationClick := func(op operation) {
		applyOperation()

		setOperationType(op)
		setOperationValue(0)
	}

	onDigitClick := func(digit int) {
		setTotal(total)
		setOperationType(operationType)

		setOperationValue(operationValue*10 + digit)
	}

	return goatw.Center{
		Child: goatw.Column{
			MainAxisShrinkWrap: true,
			CrossAxisAlignment: goatw.CrossAxisAlignmentEnd,
			Children: []goat.Widget{
				goatw.Text{Text: strconv.Itoa(total)},
				goatw.Text{Text: operationStr[operationType] + strconv.Itoa(operationValue)},
				goatw.Row{
					MainAxisShrinkWrap: true,
					Children: []goat.Widget{
						goatw.Button{
							Label:   "←",
							Padding: squareButtonPad,
							OnActivate: func() {
								setOperationValue(operationValue / 10)
							},
							KeyActivators: []string{tcell.KeyNames[tcell.KeyBackspace]},
						},
						goatw.Button{
							Label:   "C",
							Padding: squareButtonPad,
							OnActivate: func() {
								setOperationType(operationNone)
								setOperationValue(0)
							},
							KeyActivators: []string{tcell.KeyNames[tcell.KeyDelete]},
						},
						goatw.Button{
							Label: "AC",
							Padding: goat.EdgeInserts{
								Top:    1,
								Left:   2,
								Right:  3,
								Bottom: 1,
							},
							OnActivate: func() {
								setTotal(0)
								setOperationType(operationNone)
								setOperationValue(0)
							},
							KeyActivators: []string{"c"},
						},
						goatw.Button{
							Label:   "÷",
							Padding: squareButtonPad,
							OnActivate: func() {
								onOperationClick(operationDiv)
							},
							KeyActivators: []string{"/"},
						},
					},
				},
				goatw.Row{
					MainAxisShrinkWrap: true,
					Children: []goat.Widget{
						digitButton{Digit: 7, OnClick: onDigitClick},
						digitButton{Digit: 8, OnClick: onDigitClick},
						digitButton{Digit: 9, OnClick: onDigitClick},
						goatw.Button{
							Label:   "×",
							Padding: squareButtonPad,
							OnActivate: func() {
								onOperationClick(operationMul)
							},
							KeyActivators: []string{"*"},
						},
					},
				},
				goatw.Row{
					MainAxisShrinkWrap: true,
					Children: []goat.Widget{
						digitButton{Digit: 4, OnClick: onDigitClick},
						digitButton{Digit: 5, OnClick: onDigitClick},
						digitButton{Digit: 6, OnClick: onDigitClick},
						goatw.Button{
							Label:   "-",
							Padding: squareButtonPad,
							OnActivate: func() {
								onOperationClick(operationSub)
							},
							KeyActivators: []string{"-"},
						},
					},
				},
				goatw.Row{
					MainAxisShrinkWrap: true,
					Children: []goat.Widget{
						digitButton{Digit: 1, OnClick: onDigitClick},
						digitButton{Digit: 2, OnClick: onDigitClick},
						digitButton{Digit: 3, OnClick: onDigitClick},
						goatw.Button{
							Label:   "+",
							Padding: squareButtonPad,
							OnActivate: func() {
								onOperationClick(operationAdd)
							},
							KeyActivators: []string{"+"},
						},
					},
				},
				goatw.Row{
					MainAxisShrinkWrap: true,
					Children: []goat.Widget{
						goatw.Button{
							Label: "0",
							Padding: goat.EdgeInserts{
								Top:    1,
								Left:   7,
								Right:  6,
								Bottom: 1,
							},
							OnActivate: func() {
								onDigitClick(0)
							},
							KeyActivators: []string{"0"},
						},
						goatw.Button{
							Label: "=",
							Padding: goat.EdgeInserts{
								Top:    1,
								Left:   6,
								Right:  7,
								Bottom: 1,
							},
							OnActivate: func() {
								applyOperation()

								setOperationType(operationNone)
								setOperationValue(0)
							},
							KeyActivators: []string{"=", tcell.KeyNames[tcell.KeyEnter]},
						},
					},
				},
			},
		},
	}, nil
}

type digitButton struct {
	goat.Widget

	Digit   int
	OnClick func(digit int)
}

var _ goat.StateWidget = digitButton{}

func (w digitButton) Build() (goat.Widget, error) {
	str := strconv.Itoa(w.Digit)

	return goatw.Button{
		Label:   str,
		Padding: squareButtonPad,
		OnActivate: func() {
			w.OnClick(w.Digit)
		},
		KeyActivators: []string{str},
	}, nil
}

func main() {
	err := goat.RunApp(app{})
	if err != nil {
		panic(err.Error())
	}
}
