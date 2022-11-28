package ui

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var FontAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)

func signalColor(signal bool) color.Color {
	if signal {
		return colornames.Red
	}
	return colornames.Black
}

func drawWire(imd *imdraw.IMDraw, a, b pixel.Vec, col color.Color) {
	imd.Color = col
	imd.Push(a)
	imd.Push(b)
	imd.Line(WIRE_WIDTH)

	imd.Push(a)
	imd.Circle(HOLE_RADIUS, 0)
	imd.Push(b)
	imd.Circle(HOLE_RADIUS, 0)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func drawChip(win *pixelgl.Window, mp pixel.Vec, title string, inputs []bool, outputs []bool) {
	imd := imdraw.New(nil)

	for i, s := range inputs {
		imd.Color = signalColor(s)
		imd.Push(mp.Add(pixel.V(
			0,
			-GRID_SIZE*float64(i),
		)))
		imd.Circle(HOLE_RADIUS, 0)
	}
	for i, s := range outputs {
		imd.Color = signalColor(s)
		imd.Push(mp.Add(pixel.V(
			GRID_SIZE,
			-GRID_SIZE*(float64(i)),
		)))
		imd.Circle(HOLE_RADIUS, 0)
	}

	imd.Color = colornames.Darkgray
	height := max(len(inputs), len(outputs))
	imd.Push(mp.Add(pixel.V(0, GRID_SIZE/2)))
	imd.Push(mp.Add(pixel.V(
		GRID_SIZE,
		-GRID_SIZE*(float64(height)-0.5),
	)))
	imd.Rectangle(0)
	imd.Draw(win)

	txt := text.New(pixel.ZV, FontAtlas)
	txt.Color = colornames.Black
	fmt.Fprint(txt, title)
	textPos := mp.Sub(txt.Bounds().Center()).
		Add(pixel.V(GRID_SIZE/2, 0))
	txt.Draw(win, pixel.IM.Moved(textPos))
}

const SWITCH_WIDTH = GRID_SIZE
const SWITCH_HEIGHT = GRID_SIZE
const LAMP_RADIUS = GRID_SIZE / 2

func drawSwitch(win *pixelgl.Window, mp pixel.Vec, state bool) {
	title := "OFF"
	if state {
		title = "ON"
	}
	drawChip(win, mp, title, []bool{}, []bool{state})
}

func drawLamp(imd *imdraw.IMDraw, mp pixel.Vec, state bool) {
	imd.Color = signalColor(state)
	imd.Push(mp)
	imd.Circle(LAMP_RADIUS, 0)
}
