package main

import (
	"math/rand"

	"github.com/gdamore/tcell"
)

type mino struct {
	point
	shapeStatus int
	spinStatus  int
}

func newMino() *mino {
	r := rand.Intn(7)
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

func (m *mino) move(k tcell.Key) {
	white := tcell.ColorWhite
	// now := tcell.ColorPaleGreen

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
		new := newMino()
		*m = *new
	}
	m.draw(white)
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
		return false
	}

	old := m.spinStatus
	m.spinStatus = (m.spinStatus + 1) % 4

	if m.isCollided() {
		m.spinStatus = old
		return false
	}
	return true
}
