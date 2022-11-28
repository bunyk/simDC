package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type LampTool struct {
}

func (st *LampTool) Update(win *pixelgl.Window, cb *CircuitBoard, mp pixel.Vec) {
	mp = gridAlign(mp)
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		cb.AddLamp(
			int(mp.X/GRID_SIZE),
			int(mp.Y/GRID_SIZE),
		)
	}
}

func (st *LampTool) Draw(win *pixelgl.Window, mp pixel.Vec) {
	mp = gridAlign(mp)
	imd := imdraw.New(nil)
	drawLamp(imd, mp, false)
	imd.Draw(win)
}
