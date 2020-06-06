package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
)

// ごちゃごちゃしているので整理しましょう

var sc tcell.Screen

var keyCh chan tcell.Key

func main() {
	rand.Seed(time.Now().UnixNano())

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	sc, _ = tcell.NewScreen()
	sc.Init()
	defer sc.Fini()
	sc.SetStyle(tcell.StyleDefault)
	sc.Clear()

	ctx, cancel := context.WithCancel(context.Background())
	gameInit(ctx, cancel)
	game(ctx)
}

func game(ctx context.Context) {
	favColor := tcell.ColorPaleGreen
	m := newMino()
	m.draw(favColor)

	for {
		sc.Show()
		select {
		case <-ctx.Done():
			return
		case k := <-keyCh:
			m.move(k, favColor)
		}
	}
}

func gameInit(ctx context.Context, cancel context.CancelFunc) {
	maxP = point{21, 12}
	startP = point{0, 5}

	field = func() [][]bool {
		ret := make([][]bool, maxP.y)
		for i := range ret {
			ret[i] = make([]bool, maxP.x)
		}
		return ret
	}()

	keyCh = make(chan tcell.Key)

	go poll(ctx, cancel)
	go fall(ctx)

	drawWall()
}

func poll(ctx context.Context, cancel context.CancelFunc) {
	for {
		ev := sc.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				cancel()
				return
			default:
				keyCh <- ev.Key()
			}
		}
	}
}

func fall(ctx context.Context) {
	for {
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			return
		default:
			keyCh <- tcell.KeyDown
		}
	}
}
