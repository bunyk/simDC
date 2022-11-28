package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type ChipTool struct {
	Class string
}

func (ct ChipTool) Update(win *pixelgl.Window, cb *CircuitBoard, mp pixel.Vec) {
	mp = cellAlign(mp)
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		cb.AddChip(ct.Class, int(mp.X/GRID_SIZE), int(mp.Y/GRID_SIZE))
	}
}

var unpoweredPins = []bool{
	false, false, false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false, false, false,
}

func (ct ChipTool) Draw(win *pixelgl.Window, mp pixel.Vec) {
	mp = cellAlign(mp)
	drawChip(
		win,
		mp,
		ct.Class,
		unpoweredPins[0:ChipClasses[ct.Class].InputsCount],
		unpoweredPins[0:ChipClasses[ct.Class].OutputsCount],
	)
}
