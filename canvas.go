package goat

import (
	"fmt"
	"image"
)

type Cell struct {
	// Default rune represents no change in text
	Rune rune

	Background, Foreground Color

	TextStyle *CellTextStyle
}

type CellTextStyle struct {
	Bold, Blink, Dim, Italic, Underline, StrikeThrough bool

	Url   string
	UrlId string
}

func (bottom Cell) Blend(top Cell) Cell {
	result := Cell{
		Rune:       top.Rune,
		Background: bottom.Background.Blend(top.Background),
		// TODO: take background into account with transparent foreground
		Foreground: bottom.Foreground.Blend(top.Foreground),
		TextStyle:  top.TextStyle,
	}

	// If top cell does not override the rune or text style, then fallback to the bottom cell's
	if top.Rune == rune(0) {
		result.Rune = bottom.Rune
	}
	if top.TextStyle == nil {
		result.TextStyle = bottom.TextStyle
	}

	return result
}

type Canvas struct {
	size  Size
	cells []Cell
}

func NewCanvas(size Size) Canvas {
	if size.HasInf() {
		panic(fmt.Sprint("canvas cannot be created with infinite size: ", size.String()))
	}

	return Canvas{
		size:  size,
		cells: make([]Cell, size.Width.Int()*size.Height.Int()),
	}
}

func (c *Canvas) Size() Size {
	return c.size
}

func (c *Canvas) GetCell(x, y int) Cell {
	return c.cells[y*c.size.Width.Int()+x]
}

func (c *Canvas) SetCell(x, y int, cell Cell) {
	c.cells[y*c.size.Width.Int()+x] = cell
}

func (c *Canvas) FillBackground(x, y, width, height int, background Color) {
	for i := y; i < y+height; i++ {
		for j := x; j < x+width; j++ {
			cell := &c.cells[i*c.size.Width.Int()+j]
			cell.Background = background
		}
	}
}

func (c *Canvas) OverlayCanvas(x, y int, topCanvas Canvas) {
	topCanvasIndex := 0
	for i := y; i < y+topCanvas.size.Height.Int(); i++ {
		for j := x; j < x+topCanvas.size.Width.Int(); j++ {
			bottomCell := &c.cells[i*c.size.Width.Int()+j]
			topCell := topCanvas.cells[topCanvasIndex]

			*bottomCell = bottomCell.Blend(topCell)

			topCanvasIndex += 1
		}
	}
}

func (c *Canvas) OverlayImage(x, y int, image image.Image) {
	imageBound := image.Bounds()
	imageWidth := imageBound.Dx()
	imageHeight := imageBound.Dy()

	imageX := image.Bounds().Min.X
	imageY := image.Bounds().Min.Y

	for i := y; i < y+imageHeight; i++ {
		for j := x; j < x+imageWidth; j++ {
			c.cells[i*c.size.Width.Int()+j] = Cell{
				Rune:       ' ',
				Background: ColorFromImageColor(image.At(imageX, imageY)),
			}

			imageX++
		}

		imageX = image.Bounds().Min.X
		imageY++
	}
}
