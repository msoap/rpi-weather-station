package main

import (
	"image"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Assuming your object has this interface
type PixelDrawer interface {
	SetPixel(x, y int, on bool)
}

// DrawText draws a string at position (x, y)
// Returns the total width of the text drawn
func DrawText(d PixelDrawer, x, y int, text string) int {
	face := basicfont.Face7x13

	// Measure text to create appropriately sized image
	width := MeasureText(text)
	height := face.Metrics().Height.Ceil()
	ascent := face.Metrics().Ascent.Ceil()

	if width == 0 {
		return 0
	}

	// Create image and draw text onto it
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(0, ascent),
	}
	drawer.DrawString(text)

	// Copy pixels to our display
	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			_, _, _, a := img.At(px, py).RGBA()
			if a > 0 {
				d.SetPixel(x+px, y+py, true)
			}
		}
	}

	return width
}

// DrawChar draws a single character at position (x, y)
// Returns the width of the character drawn
func DrawChar(d PixelDrawer, x, y int, c rune) int {
	return DrawText(d, x, y, string(c))
}

// DrawTextCentered draws text centered horizontally at position (cx, y)
func DrawTextCentered(d PixelDrawer, cx, y int, text string) {
	width := MeasureText(text)
	DrawText(d, cx-width/2, y, text)
}

// DrawTextRight draws text right-aligned ending at position (x, y)
func DrawTextRight(d PixelDrawer, x, y int, text string) {
	width := MeasureText(text)
	DrawText(d, x-width, y, text)
}

// MeasureText returns the width of the text without drawing it
func MeasureText(text string) int {
	face := basicfont.Face7x13
	var width fixed.Int26_6
	for _, c := range text {
		adv, ok := face.GlyphAdvance(c)
		if ok {
			width += adv
		}
	}
	return width.Ceil()
}

// FontHeight returns the height of the basicfont
func FontHeight() int {
	return basicfont.Face7x13.Metrics().Height.Ceil()
}

// FontAscent returns the ascent (pixels above baseline)
func FontAscent() int {
	return basicfont.Face7x13.Metrics().Ascent.Ceil()
}

// Ensure draw package is used (for side effects if needed)
var _ = draw.Draw
