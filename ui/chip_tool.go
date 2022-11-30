package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
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

func (ct ChipTool) Draw(win *pixelgl.Window, mp pixel.Vec) {
	mp = cellAlign(mp)
	imd := imdraw.New(nil)
	class := ChipClasses[ct.Class]
	drawChip(imd, mp, ct.Class, class.Height())
	drawChipPins(imd, mp, class.InputsCount, class.OutputsCount)
	imd.Draw(win)
}
