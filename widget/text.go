package goatw

import (
	"strings"
	"unicode"

	. "github.com/jwr1/goat"

	"github.com/gdamore/tcell/v2"
)

type Text struct {
	Widget

	Text  string
	Style tcell.Style
}

var _ RenderWidget = Text{}

func (w Text) Layout(context LayoutContext) (Size, error) {
	x, y, maxLineWidth := 0, 0, 0

	for _, r := range wordWrap(w.Text, context.Constraints.Max.Width) {
		switch r {
		case '\n':
			x = 0
			y++
		default:
			x++
			if x > maxLineWidth {
				maxLineWidth = x
			}
		}
	}

	return Size{
		Width:  max(maxLineWidth, context.Constraints.Min.Width),
		Height: max(y+1, context.Constraints.Min.Height),
	}, nil
}

func (w Text) Paint(context PaintContext) error {
	x, y := 0, 0
	for _, r := range wordWrap(w.Text, context.Size.Width) {
		switch r {
		case '\n':
			x = 0
			y++
		default:
			context.Canvas.SetCell(x, y, r, w.Style)
			x++
		}
	}

	return nil
}

func wordWrap(text string, maxWidth int) string {
	var output strings.Builder
	var line strings.Builder
	var block strings.Builder

	curLineWidth := 0

	nextLine := func() {
		output.WriteString(strings.TrimRightFunc(line.String(), unicode.IsSpace))
		output.WriteRune('\n')
		line.Reset()
		curLineWidth = 0
	}

	consumeBlock := func() {
		if block.Len() == 0 {
			return
		}

		// If the block doesn't fit on this line, but would on the next.
		if maxWidth-curLineWidth < block.Len() && block.Len() < maxWidth {
			nextLine()
		}

		for _, r := range block.String() {
			if curLineWidth >= maxWidth {
				nextLine()
			}

			line.WriteRune(r)
			curLineWidth++
		}

		block.Reset()
	}

	for _, r := range text {
		switch r {
		case '\n':
			consumeBlock()
			nextLine()
		case ' ':
			consumeBlock()
			if curLineWidth >= maxWidth {
				nextLine()
			} else {
				line.WriteRune(' ')
				curLineWidth++
			}
		default:
			block.WriteRune(r)
		}
	}

	consumeBlock()
	if line.Len() > 0 {
		nextLine()
	}

	return strings.TrimRightFunc(output.String(), unicode.IsSpace)
}
