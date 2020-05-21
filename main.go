package main

import (
	"context"
	"time"

	"github.com/gdamore/tcell"
)

type direction int

const (
	up direction = iota
	right
	down
	left
)

var sc tcell.Screen

func main() {
	maxX, maxY = 12, 21
	startP.x, startP.y = 5, 0

	field = func() [][]bool {
		ret := make([][]bool, maxY)
		for i := range ret {
			ret[i] = make([]bool, maxX)
		}
		return ret
	}()

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	sc, _ = tcell.NewScreen()
	sc.Init()
	defer sc.Fini()
	sc.SetStyle(tcell.StyleDefault)
	sc.Clear()

	game()
}

func game() {
	ctx, cancel := context.WithCancel(context.Background())
	dch := make(chan direction)

	screenInit()

	go poll(ctx, cancel, dch)
	go fall(ctx, dch)

LOOP:
	for p := startP; ; {
		select {
		case <-ctx.Done():
			break LOOP
		case d := <-dch:
			p = move(p, d)
		}

		sc.Show()
	}
}

func poll(ctx context.Context, cancel context.CancelFunc, dch chan direction) {
	for {
		ev := sc.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				cancel()
				return
			case tcell.KeyUp:
				dch <- up
			case tcell.KeyRight:
				dch <- right
			case tcell.KeyDown:
				dch <- down
			case tcell.KeyLeft:
				dch <- left
			}
		}
	}
}

func fall(ctx context.Context, dch chan direction) {
	for {
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			return
		default:
			dch <- down
		}
	}
}
