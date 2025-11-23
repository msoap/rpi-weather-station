package main

import (
	"github.com/quasilyte/bitsweetfont"
	"golang.org/x/image/font"
)

type DispDrawler interface {
	Clear()
	Update()
	Size() (width, height int)
	SetPixel(x, y int, on bool)
	GetPixel(x, y int) bool
}

type Screen struct {
	width, height int
	disp          DispDrawler
	fontFace      font.Face
}

func NewScreen(disp DispDrawler) Screen {
	width, height := disp.Size()
	return Screen{
		width:    width,
		height:   height,
		disp:     disp,
		fontFace: bitsweetfont.New1(),
	}
}

func (sc *Screen) SetPixel(x, y int, on bool) {
	sc.disp.SetPixel(x, y, on)
}

func (sc *Screen) GetPixel(x, y int) bool {
	return sc.disp.GetPixel(x, y)
}

func (sc *Screen) Clear() {
	sc.disp.Clear()
}

func (sc *Screen) Update() {
	sc.disp.Update()
}

func (sc *Screen) DrawHorizontalLine(x, y int, width int, on bool) {
	for i := 0; i < width; i++ {
		if x+i < 0 || x+i >= sc.width || y < 0 || y >= sc.height {
			continue
		}
		sc.SetPixel(x+i, y, on)
	}
}

func (sc *Screen) DrawVerticalLine(x, y int, height int, on bool) {
	for i := 0; i < height; i++ {
		if x < 0 || x >= sc.width || y+i < 0 || y+i >= sc.height {
			continue
		}
		sc.SetPixel(x, y+i, on)
	}
}

func (sc *Screen) Rectangle(x, y, w, h int, on bool) {
	sc.DrawHorizontalLine(x, y, w, on)
	sc.DrawHorizontalLine(x, y+h-1, w, on)
	sc.DrawVerticalLine(x, y, h, on)
	sc.DrawVerticalLine(x+w-1, y, h, on)
}

func (sc *Screen) Box(x, y, w, h int, on bool) {
	for i := range w {
		for j := range h {
			sc.SetPixel(x+i, y+j, on)
		}
	}
}

// Line - draw line using the Bresenham's algorithm
func (sc *Screen) Line(x1, y1, x2, y2 int, on bool) {
	dx := abs(x2 - x1)
	dy := -abs(y2 - y1)

	sx := sgn(x2 - x1)
	sy := sgn(y2 - y1)

	e := dx + dy
	x0, y0 := x1, y1

	for {
		sc.SetPixel(x0, y0, on)

		if x0 == x2 && y0 == y2 {
			break
		}
		e2 := 2 * e

		if e2 >= dy {
			if x0 == x2 {
				break
			}
			e = e + dy
			x0 = x0 + sx
		}

		if e2 <= dx {
			if y0 == y2 {
				break
			}
			e = e + dx
			y0 = y0 + sy
		}
	}
}

// Circle - draw a circle using the Midpoint Circle Algorithm
func (sc *Screen) Circle(x, y, r int, on bool) {
	if r < 0 {
		return
	}

	x1, y1, err := -r, 0, 2-2*r
	for {
		sc.SetPixel(x-x1, y+y1, on)
		sc.SetPixel(x-y1, y-x1, on)
		sc.SetPixel(x+x1, y-y1, on)
		sc.SetPixel(x+y1, y+x1, on)
		rr := err
		if rr > x1 {
			x1++
			err += x1*2 + 1
		}
		if rr <= y1 {
			y1++
			err += y1*2 + 1
		}
		if x1 >= 0 {
			break
		}
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func sgn(a int) int {
	if a == 0 {
		return 0
	}
	if a > 0 {
		return 1
	}
	return -1
}
