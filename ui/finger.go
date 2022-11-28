package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Finger struct {
}

func (f *Finger) Update(win *pixelgl.Window, cb *CircuitBoard, mp pixel.Vec) {
	mp = cellAlign(mp)
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		cb.PressSwitch(
			int(mp.X/GRID_SIZE),
			int(mp.Y/GRID_SIZE),
		)
	}
}

func (f Finger) Draw(win *pixelgl.Window, mp pixel.Vec) {
}
