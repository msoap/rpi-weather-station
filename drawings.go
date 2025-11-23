package main

import (
	"github.com/msoap/tcg"
	"github.com/quasilyte/bitsweetfont"
	"golang.org/x/image/font"
)

type dispDrawler interface {
	Clear()
	Update()
	Size() (width, height int)
	SetPixel(x, y int, on bool)
	Finish()
}

type Screen struct {
	tcg.Buffer
	width, height int
	disp          dispDrawler
	fontFace      font.Face
}

func NewScreen(disp dispDrawler) Screen {
	width, height := disp.Size()
	return Screen{
		Buffer: tcg.NewBuffer(width, height),
		width:  width,
		height: height,
		disp:   disp,
		// fontFace: basicfont.Face7x13,
		fontFace: bitsweetfont.New1(),
	}
}

func (sc *Screen) Clear() {
	sc.Buffer.Clear()
	sc.disp.Clear()
}

func (sc *Screen) Update() {
	for y := 0; y < sc.height; y++ {
		for x := 0; x < sc.width; x++ {
			on := sc.Buffer.At(x, y) == tcg.Black
			sc.disp.SetPixel(x, y, on)
		}
	}

	sc.disp.Update()
}

func (sc *Screen) Finish() {
	sc.disp.Finish()
}
