package goat

import (
	"github.com/gdamore/tcell/v2"
)

type canvasCell struct {
	rune  rune
	style tcell.Style
}

type Canvas struct {
	size  Size
	cells []canvasCell
}

func NewCanvas(size Size) Canvas {
	return Canvas{
		size:  size,
		cells: make([]canvasCell, size.Width*size.Height),
	}
}

func (c *Canvas) Size() Size {
	return c.size
}

func (c *Canvas) GetCell(x, y int) (rune, tcell.Style) {
	cell := &c.cells[y*c.size.Width+x]
	return cell.rune, cell.style
}

func (c *Canvas) SetCell(x, y int, rune rune, style tcell.Style) {
	cell := &c.cells[y*c.size.Width+x]
	cell.rune = rune
	cell.style = style
}

func (c *Canvas) FillStyle(x, y, width, height int, style tcell.Style) {
	for i := y; i < y+height; i++ {
		for j := x; j < x+width; j++ {
			cell := &c.cells[i*c.size.Width+j]
			cell.style = style
			c.cells[i*c.size.Width+j] = *cell
		}
	}
}

func (c *Canvas) OverlayCanvas(x, y int, topCanvas Canvas) {
	topCanvasIndex := 0
	for i := y; i < y+topCanvas.size.Height; i++ {
		for j := x; j < x+topCanvas.size.Width; j++ {
			bottomCell := &c.cells[i*c.size.Width+j]
			topCell := topCanvas.cells[topCanvasIndex]

			bottomCell.rune = topCell.rune
			fg, bg, attr := topCell.style.Decompose()
			bottomCell.style = bottomCell.style.Attributes(attr)
			if fg != tcell.ColorDefault {
				bottomCell.style = bottomCell.style.Foreground(fg)
			}
			if bg != tcell.ColorDefault {
				bottomCell.style = bottomCell.style.Background(bg)
			}

			topCanvasIndex += 1
		}
	}
}
