package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type WireTool struct {
	WireEnd *pixel.Vec
}

func (wt *WireTool) Update(win *pixelgl.Window, cb *CircuitBoard, mp pixel.Vec) {
	mp = gridAlign(mp)
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		wt.WireEnd = &mp
		return
	}
	if win.JustReleased(pixelgl.MouseButtonLeft) && wt.WireEnd != nil {
		addWire(cb, *wt.WireEnd, mp)
		wt.WireEnd = nil
	}
}

func (wt *WireTool) Draw(win *pixelgl.Window, mp pixel.Vec) {
	mp = gridAlign(mp)
	imd := imdraw.New(nil)
	if wt.WireEnd == nil { // Selecting first end of wire
		imd.Color = colornames.Black
		imd.Push(mp)
		imd.Circle(HOLE_RADIUS, 0)
	} else { // Selecting second end
		drawWire(imd, *wt.WireEnd, mp, colornames.Black)
	}
	imd.Draw(win)
}

func addWire(cb *CircuitBoard, a, b pixel.Vec) {
	if a.Eq(b) {
		return // Wires of length 0 make no sense
	}
	cb.AddWire(
		int(a.X/GRID_SIZE),
		int(a.Y/GRID_SIZE),
		int(b.X/GRID_SIZE),
		int(b.Y/GRID_SIZE),
		false,
	)
}
