package ui

import (
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

func drawChip(imd *imdraw.IMDraw, win *pixelgl.Window, mp pixel.Vec, title string, inputs []bool, outputs []bool) {
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

	labels[title] = append(labels[title], mp.Add(pixel.V(GRID_SIZE/2, 0)))
}

var labels = make(map[string][]pixel.Vec) // map strings to list of locations where they should appear

const SWITCH_WIDTH = GRID_SIZE
const SWITCH_HEIGHT = GRID_SIZE
const LAMP_RADIUS = GRID_SIZE / 2

func drawSwitch(imd *imdraw.IMDraw, win *pixelgl.Window, mp pixel.Vec, state bool) {
	title := "OFF"
	if state {
		title = "ON"
	}
	drawChip(imd, win, mp, title, []bool{}, []bool{state})
}

func drawLamp(imd *imdraw.IMDraw, mp pixel.Vec, state bool) {
	imd.Color = signalColor(state)
	imd.Push(mp)
	imd.Circle(LAMP_RADIUS, 0)
}
