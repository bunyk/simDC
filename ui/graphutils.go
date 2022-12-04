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
}

func drawNode(imd *imdraw.IMDraw, pos pixel.Vec, signal bool) {
	imd.Color = signalColor(signal)
	imd.Push(pos)
	imd.Circle(HOLE_RADIUS, 0)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func drawChip(imd *imdraw.IMDraw, pos pixel.Vec, title string, height int) {

	imd.Color = colornames.Darkgray
	imd.Push(pos.Add(pixel.V(0, GRID_SIZE/2)))
	imd.Push(pos.Add(pixel.V(
		GRID_SIZE,
		-GRID_SIZE*(float64(height)-0.5),
	)))
	imd.Rectangle(0)

	labels[title] = append(labels[title], pos.Add(pixel.V(GRID_SIZE/2, 0)))
}

var labels = make(map[string][]pixel.Vec) // map strings to list of locations where they should appear

func DrawLabels(win *pixelgl.Window) {
	for label, locations := range labels {
		txt := text.New(pixel.ZV, FontAtlas)
		txt.Color = colornames.Black
		fmt.Fprint(txt, label)
		tc := txt.Bounds().Center()
		for _, loc := range locations {
			txt.Draw(win, pixel.IM.Moved(loc.Sub(tc)))
		}
	}
	labels = make(map[string][]pixel.Vec)
}

func drawChipPins(imd *imdraw.IMDraw, pos pixel.Vec, numLeft, numRight int) {
	imd.Color = signalColor(false)
	for i := 0; i < numLeft; i++ {
		imd.Push(pos.Add(pixel.V(
			0,
			-GRID_SIZE*float64(i),
		)))
		imd.Circle(HOLE_RADIUS, 0)
	}
	for i := 0; i < numRight; i++ {
		imd.Push(pos.Add(pixel.V(
			GRID_SIZE,
			-GRID_SIZE*float64(i),
		)))
		imd.Circle(HOLE_RADIUS, 0)
	}
}

const SWITCH_WIDTH = GRID_SIZE
const SWITCH_HEIGHT = GRID_SIZE
const LAMP_RADIUS = GRID_SIZE / 2

func drawSwitch(imd *imdraw.IMDraw, win *pixelgl.Window, pos pixel.Vec, state bool) {
	title := "OFF"
	if state {
		title = "ON"
	}
	drawChip(imd, pos, title, 1)
}

func drawLamp(imd *imdraw.IMDraw, pos pixel.Vec, state bool) {
	imd.Color = signalColor(state)
	imd.Push(pos)
	imd.Circle(LAMP_RADIUS, 0)
}
