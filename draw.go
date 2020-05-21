package main

import (
	"github.com/gdamore/tcell"
)

type point struct {
	x, y int
}

var (
	maxX   int
	maxY   int
	field  [][]bool
	startP point
)

func drawSquare(c tcell.Color, p point) {
	x, y := p.x, p.y
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
	for y := 0; y < maxY; y++ {
		drawSquare(white, point{0, y})
		drawSquare(white, point{maxX - 1, y})
	}

	for x := 0; x < maxX; x++ {
		drawSquare(white, point{x, maxY - 1})
	}
}

func screenInit() {
	white := tcell.ColorWhite
	drawWall()
	drawSquare(white, startP)
	sc.Show()
}

func move(old point, d direction) point {
	white := tcell.ColorWhite
	def := tcell.ColorDefault

	next := old
	switch d {
	case right:
		if !field[next.y][next.x+1] {
			next.x++
		}
	case left:
		if !field[next.y][next.x-1] {
			next.x--
		}
	case down:
		if !field[next.y+1][next.x] {
			next.y++
		}
	case up:
		for !field[next.y+1][next.x] {
			next.y++
		}
	}

	if field[next.y+1][next.x] {
		drawSquare(white, next)
		next = startP
	}

	drawSquare(white, next)
	if old != next {
		drawSquare(def, old)
	}

	return next
}
