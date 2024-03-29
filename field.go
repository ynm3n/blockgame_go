package main

import (
	"github.com/gdamore/tcell"
)

type point struct {
	y, x int
}

func (p *point) add(a point) {
	p.y += a.y
	p.x += a.x
}

func (p *point) sum(a point) point {
	return point{p.y + a.y, p.x + a.x}
}

func (p *point) sub(a point) {
	p.y -= a.y
	p.x -= a.x
}

func (p *point) field() bool {
	return field[p.y][p.x]
}

var (
	field  [][]bool
	maxP   point
	startP point
)

var direction = map[tcell.Key]point{
	tcell.KeyRight: {0, 1},
	tcell.KeyDown:  {1, 0},
	tcell.KeyLeft:  {0, -1},
}

func drawSquare(c tcell.Color, y, x int) {
	style := tcell.StyleDefault.Background(c)
	space := ' '
	sc.SetContent(2*x, y, space, nil, style)
	sc.SetContent(2*x+1, y, space, nil, style)
	if c == tcell.ColorDefault {
		field[y][x] = false
	} else {
		field[y][x] = true
	}
}

func drawWall() {
	white := tcell.ColorWhite
	for y := 0; y < maxP.y; y++ {
		drawSquare(white, y, 0)
		drawSquare(white, y, maxP.x-1)
	}

	for x := 0; x < maxP.x; x++ {
		drawSquare(white, maxP.y-1, x)
	}
}

func clearLine() {
	diffs := make([]int, maxP.y)
	cnt := 0
	minY := 0
	for y := maxP.y - 2; y >= 0; y-- {
		full := true
		empty := true
		for x := 1; x <= maxP.x-2; x++ {
			if field[y][x] {
				empty = false
			} else {
				full = false
			}
		}

		if full {
			diffs[y] = 0
			cnt++
		} else {
			diffs[y] = cnt
		}
		if empty {
			minY = y + 1
			break
		}
	}

	for y := maxP.y - 1; y >= minY; y-- {
		d := diffs[y]
		if d == 0 {
			continue
		}

		for x := 1; x <= maxP.x-2; x++ {
			if field[y][x] {
				drawSquare(tcell.ColorWhite, y+d, x)
			} else {
				drawSquare(tcell.ColorDefault, y+d, x)
			}
		}
	}
	for y := minY; y < minY+cnt; y++ {
		for x := 1; x <= maxP.x-2; x++ {
			drawSquare(tcell.ColorDefault, y, x)
		}
	}
}
