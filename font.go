package main

import (
	"image"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// DrawText draws a string at position (x, y)
// Returns the total width of the text drawn
func (sc *Screen) DrawText(x, y int, text string) int {
	// Measure text to create appropriately sized image
	width := sc.measureText(text)
	height := sc.fontFace.Metrics().Height.Ceil()
	ascent := sc.fontFace.Metrics().Ascent.Ceil()

	if width == 0 {
		return 0
	}

	// Create image and draw text onto it
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: sc.fontFace,
		Dot:  fixed.P(0, ascent),
	}
	drawer.DrawString(text)

	// Copy pixels to our display
	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			_, _, _, a := img.At(px, py).RGBA()
			color := 0
			if a > 0 {
				color = 1
			}
			sc.Buffer.Set(x+px, y+py, color)
		}
	}

	return width
}

// DrawTextCentered draws text centered horizontally at position (cx, y)
func (sc *Screen) DrawTextCentered(cx, y int, text string) {
	width := sc.measureText(text)
	sc.DrawText(cx-width/2, y, text)
}

// DrawTextRight draws text right-aligned ending at position (x, y)
func (sc *Screen) DrawTextRight(x, y int, text string) {
	width := sc.measureText(text)
	sc.DrawText(x-width, y, text)
}

// measureText returns the width of the text without drawing it
func (sc *Screen) measureText(text string) int {
	var width fixed.Int26_6
	for _, c := range text {
		adv, ok := sc.fontFace.GlyphAdvance(c)
		if ok {
			width += adv
		}
	}
	return width.Ceil()
}
