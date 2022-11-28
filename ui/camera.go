package ui

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const MAX_ZOOM = 2.0
const MIN_ZOOM = 0.5

type Camera struct {
	Pos       pixel.Vec
	Speed     float64
	Zoom      float64
	ZoomSpeed float64
}

func (c *Camera) Update(win *pixelgl.Window, dt float64) {
	if win.Pressed(pixelgl.KeyLeft) {
		c.Pos.X -= c.Speed * dt
	}
	if win.Pressed(pixelgl.KeyRight) {
		c.Pos.X += c.Speed * dt
	}
	if win.Pressed(pixelgl.KeyDown) {
		c.Pos.Y -= c.Speed * dt
	}
	if win.Pressed(pixelgl.KeyUp) {
		c.Pos.Y += c.Speed * dt
	}

	c.Zoom *= math.Pow(c.ZoomSpeed, win.MouseScroll().Y)
	if c.Zoom > MAX_ZOOM {
		c.Zoom = MAX_ZOOM
	}
	if c.Zoom < MIN_ZOOM {
		c.Zoom = MIN_ZOOM
	}
}

func (c Camera) Matrix(win *pixelgl.Window) pixel.Matrix {
	return pixel.IM.Scaled(c.Pos, c.Zoom).Moved(
		win.Bounds().Center().Sub(c.Pos),
	)
}
