package main

import (
	"math/rand"

	"github.com/gdamore/tcell"
)

type mino struct {
	point
	shape int
	state int
}

func newMino() *mino {
	m := &mino{startP, rand.Intn(7), 0}
	return m
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
	ps := func() []point {
		ret := make([]point, 4)
		diffs := minoShapes[m.shape]
		for i, d := range diffs {
			switch m.state {
			case 0:
				// なにもしないよ！
			case 1:
				d.y, d.x = d.x, -d.y
			case 2:
				d.y, d.x = -d.y, -d.x
			case 3:
				d.y, d.x = -d.x, d.y
			}
			ret[i] = m.sum(d)
		}
		return ret
	}()

	ok := func() bool {
		if c == tcell.ColorDefault {
			return true
		}

		for _, p := range ps {
			if p.field() {
				return false
			}
		}
		return true
	}()

	if ok {
		for _, p := range ps {
			drawSquare(c, p.y, p.x)
		}
	} else {
		m.state--
		if m.state < 0 {
			m.state = 3
		}
		m.draw(c)
	}
}

func (m *mino) move(k tcell.Key) {
	white := tcell.ColorWhite
	def := tcell.ColorDefault

	m.draw(def)

	switch k {
	case tcell.KeyLeft, tcell.KeyRight, tcell.KeyDown:
		m.add(direction[k])
	case tcell.KeyUp:
		// 真下に着地する処理
	case tcell.KeyEnter:
		m.state = (m.state + 1) % 4
	}

	if m.isCollided() {
		m.sub(direction[k])
	}

	if m.isLanding() {
		m.draw(white)
		new := newMino()
		*m = *new
	}
	m.draw(white)
}

func (m *mino) isCollided() bool {
	for _, diff := range minoShapes[m.shape] {
		p := m.sum(diff)
		if p.field() {
			return true
		}
	}
	return false
}

func (m *mino) isLanding() bool {
	for _, diff := range minoShapes[m.shape] {
		p := m.sum(diff)
		if field[p.y+1][p.x] {
			return true
		}
	}
	return false
}
