package main

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/msoap/tcg"
)

type termScreen struct {
	tcg *tcg.Tcg
}

func NewTermScreen() (*termScreen, error) {
	tcgObj, err := tcg.New(tcg.Mode2x3)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			switch ev := tcgObj.TCellScreen.PollEvent().(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
					tcgObj.Finish()
					os.Exit(0)
				}
			}
		}
	}()

	return &termScreen{
		tcg: tcgObj,
	}, nil
}

func (ts *termScreen) Clear() {
	ts.tcg.Buf.Clear()
}

func (ts *termScreen) Update() {
	ts.tcg.Show()
}

func (ts *termScreen) Size() (width, height int) {
	return dispW, dispH
}

func (ts *termScreen) SetPixel(x, y int, on bool) {
	color := 0
	if on {
		color = 1
	}
	ts.tcg.Buf.Set(x, y, color)
}

func (ts *termScreen) Finish() {
	ts.tcg.Finish()
}
