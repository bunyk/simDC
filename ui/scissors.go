package ui

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Scissors struct {
	CutStart *pixel.Vec
}

const CUT_WIDTH = 1

func (s *Scissors) Update(win *pixelgl.Window, cb *CircuitBoard, mp pixel.Vec) {
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		s.CutStart = &mp
		return
	}
	if win.JustReleased(pixelgl.MouseButtonLeft) && s.CutStart != nil {
		cutThrough(cb, *s.CutStart, mp)
		s.CutStart = nil
	}
}
func (s Scissors) Draw(win *pixelgl.Window, mp pixel.Vec) {
	if s.CutStart != nil {
		imd := imdraw.New(nil)
		imd.Color = colornames.Red
		imd.Push(*s.CutStart)
		imd.Push(mp)
		imd.Line(CUT_WIDTH)
		imd.Draw(win)
	}
}

func cutThrough(cb *CircuitBoard, a, b pixel.Vec) {
	fmt.Println("Cutting", a, b)
	cb.CutWires(a, b)
	cb.CutLamps(a, b)
	cb.CutSwitches(a, b)
	cb.CutChips(a, b)
}
