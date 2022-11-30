package ui

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type SwitchTool struct {
}

func cellAlign(mp pixel.Vec) pixel.Vec {
	return pixel.V(
		math.Floor(mp.X/GRID_SIZE)*GRID_SIZE,
		math.Floor(mp.Y/GRID_SIZE+0.5)*GRID_SIZE,
	)
}

func (st *SwitchTool) Update(win *pixelgl.Window, cb *CircuitBoard, mp pixel.Vec) {
	mp = cellAlign(mp)
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		cb.AddSwitch(
			int(mp.X/GRID_SIZE),
			int(mp.Y/GRID_SIZE),
		)
	}
}

func (st *SwitchTool) Draw(win *pixelgl.Window, mp pixel.Vec) {
	mp = cellAlign(mp)
	imd := imdraw.New(nil)
	drawSwitch(imd, win, mp, false)
	imd.Draw(win)
}
