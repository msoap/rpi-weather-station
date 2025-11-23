package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/msoap/tcg"
)

type termScreen struct {
	tcg *tcg.Tcg
}

func NewTermScreen() (*termScreen, chan struct{}, error) {
	tcgObj, err := tcg.New(tcg.Mode2x3,
		tcg.WithClipCenter(dispW/2, dispH/3+1),
		tcg.WithColor("blue"),
		tcg.WithBackgroundColor("black"),
	)
	if err != nil {
		return nil, nil, err
	}

	exitCh := make(chan struct{})
	go func() {
		for {
			switch ev := tcgObj.TCellScreen.PollEvent().(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
					close(exitCh)
					return
				}
			}
		}
	}()

	return &termScreen{
		tcg: tcgObj,
	}, exitCh, nil
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
