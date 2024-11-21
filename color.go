package goat

import (
	"fmt"
	imageColor "image/color"
)

type Color struct {
	R, G, B, A uint8
}

func (c Color) String() string {
	return fmt.Sprintf("(%d,%d,%d,%d)", c.R, c.G, c.B, c.A)
}

func ColorRGB(r, g, b uint8) Color {
	return Color{r, g, b, 0xff}
}

func ColorFromImageColor(c imageColor.Color) Color {
	nrgba, _ := imageColor.NRGBAModel.Convert(c).(imageColor.NRGBA)
	return Color(nrgba)
}

func (bottom Color) Blend(top Color) Color {
	// Shortcut if top color is fully opaque or fully transparent
	switch top.A {
	case 0xFF:
		return top
	case 0x00:
		return bottom
	}

	var (
		topR = uint16(top.R)
		topG = uint16(top.G)
		topB = uint16(top.B)
		topA = uint16(top.A)

		bottomR = uint16(bottom.R)
		bottomG = uint16(bottom.G)
		bottomB = uint16(bottom.B)
		bottomA = uint16(bottom.A)
	)

	resultA := topA + bottomA/0xff*(0xff-topA)

	blendChannel := func(bottomC, topC uint16) uint8 {
		return uint8((topC*topA + bottomC*bottomA/0xff*(0xff-topA)) / resultA)
	}

	return Color{
		blendChannel(bottomR, topR),
		blendChannel(bottomG, topG),
		blendChannel(bottomB, topB),
		uint8(resultA),
	}
}
