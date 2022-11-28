package ui

import "github.com/faiface/pixel/imdraw"

type Wire struct {
	A, B GridPoint
}

type WireGroup struct {
	Signal bool
	Wires  []Wire
}

func (wg WireGroup) Draw(imd *imdraw.IMDraw) {
	for _, wire := range wg.Wires {
		drawWire(imd, wire.A.Pos(), wire.B.Pos(), signalColor(wg.Signal))
	}
}

func (wg WireGroup) ConnnectedPoints() (res []GridPoint) {
	seen := make(map[GridPoint]bool)

	addP := func(p GridPoint) {
		if seen[p] {
			return
		}
		res = append(res, p)
		seen[p] = true
	}
	for _, wire := range wg.Wires {
		addP(wire.A)
		addP(wire.B)
	}
	return
}

func (wg *WireGroup) SetSignal(s bool, cb *CircuitBoard) {
	wg.Signal = s

	for _, p := range wg.ConnnectedPoints() {
		cb.SetElementSignal(p, s)
	}
}
