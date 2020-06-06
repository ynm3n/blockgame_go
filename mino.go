package main

import (
	"math/rand"
	"sort"

	"github.com/gdamore/tcell"
)

type mino struct {
	point
	shapeStatus int
	spinStatus  int
}

func newMino() *mino {
	r := rand.Intn(len(minoShapes))
	m := &mino{startP, r, 0}
	return m
}

func (m *mino) shape() []point {
	ps := make([]point, 4)
	diffs := minoShapes[m.shapeStatus]
	for i, d := range diffs {
		switch m.spinStatus {
		case 0:
			// なにもしないよ！
		case 1:
			d.y, d.x = d.x, -d.y
		case 2:
			d.y, d.x = -d.y, -d.x
		case 3:
			d.y, d.x = -d.x, d.y
		}
		ps[i] = m.sum(d)
	}
	return ps
}

// 名前ださくね？
var minoShapes = [][]point{
	{{0, 0}, {0, -1}, {0, 1}, {0, 2}},
	{{0, 0}, {0, 1}, {1, 1}, {1, 0}},
	{{0, 0}, {0, -1}, {0, 1}, {1, -1}},
	{{0, 0}, {0, -1}, {0, 1}, {1, 1}},
	{{0, 0}, {0, -1}, {1, 0}, {1, 1}},
	{{0, 0}, {1, -1}, {1, 0}, {0, 1}},
	{{0, 0}, {0, -1}, {1, 0}, {0, 1}},
}

func (m *mino) draw(c tcell.Color) {
	ps := m.shape()
	for _, p := range ps {
		drawSquare(c, p.y, p.x)
	}
}

func (m *mino) clear() {
	m.draw(tcell.ColorDefault)
}

func (m *mino) move(k tcell.Key, c tcell.Color) {
	white := tcell.ColorWhite

	m.clear()

	switch k {
	case tcell.KeyLeft, tcell.KeyRight, tcell.KeyDown:
		m.add(direction[k])
		if m.isCollided() {
			m.sub(direction[k])
		}
	case tcell.KeyUp:
		down := tcell.KeyDown
		for {
			m.add(direction[down])
			if m.isCollided() {
				m.sub(direction[down])
				break
			}
		}
	case tcell.KeyEnter:
		m.spin()
	}

	if m.isLanding() {
		m.draw(white)
		m.clearLine()
		new := newMino()
		*m = *new
	}
	m.draw(c)
}

func (m *mino) isCollided() bool {
	ps := m.shape()
	for _, p := range ps {
		if p.x < 0 || p.x >= maxP.x {
			return true
		} else if p.y < 0 || p.y >= maxP.y {
			return true
		} else if p.field() {
			return true
		}
	}
	return false
}

func (m *mino) isLanding() bool {
	down := tcell.KeyDown
	m.add(direction[down])
	defer m.sub(direction[down])

	return m.isCollided()
}

func (m *mino) spin() bool {
	if m.shapeStatus == 1 { // 正方形は回す必要なし
		return true
	}

	old := m.spinStatus
	m.spinStatus = (m.spinStatus + 1) % 4

	if m.isCollided() {
		m.spinStatus = old
		return false
	}
	return true
}

func (m *mino) clearLine() {
	ps := m.shape()
	mp := make(map[int]bool)
	for _, p := range ps {
		mp[p.y] = true
	}

	ys := make([]int, 0, len(mp))
	for y := range mp {
		ys = append(ys, y)
	}
	sort.Ints(ys)

	minY := 1
	for _, y := range ys {
		isFull := true
		for x := 1; x <= maxP.x-2; x++ {
			if !field[y][x] {
				isFull = false
				break
			}
		}

		if isFull {
			for y1 := y; y1 >= minY; y1-- {
				end := true
				for x := 1; x <= maxP.x-2; x++ {
					if field[y1-1][x] {
						drawSquare(tcell.ColorWhite, y1, x)
						end = false
					} else {
						drawSquare(tcell.ColorDefault, y1, x)
					}
				}
				if end {
					minY = y1
					break
				}
			}
			for x := 1; x <= maxP.x-2; x++ {
				drawSquare(tcell.ColorDefault, 0, x)
			}
		}
	}
}
