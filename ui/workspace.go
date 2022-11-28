package ui

import (
	"math"
	"os"

	"github.com/bunyk/simcomp/sprites"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const GRID_SIZE = 50.0
const HOLE_RADIUS = 7.0
const WIRE_WIDTH = 3.0

type Tool interface {
	Update(win *pixelgl.Window, board *CircuitBoard, mousePosition pixel.Vec)
	Draw(win *pixelgl.Window, mousePosition pixel.Vec)
}

type Workspace struct {
	Camera Camera

	// For drawing background
	HoleSprite *pixel.Sprite
	HoleBatch  *pixel.Batch

	// For wires & elements they connect
	CircuitBoard *CircuitBoard

	// Tools
	UsedTool Tool
}

func NewWorkspace() (*Workspace, error) {
	camera := Camera{
		Speed:     500.0,
		ZoomSpeed: 1.2,
		Zoom:      1.0,
	}

	holePic, err := sprites.LoadPicture("icons/hole.png")
	if err != nil {
		return nil, err
	}

	return &Workspace{
		HoleSprite:   pixel.NewSprite(holePic, holePic.Bounds()),
		HoleBatch:    pixel.NewBatch(&pixel.TrianglesData{}, holePic),
		Camera:       camera,
		CircuitBoard: NewCircuitBoard(),
		UsedTool:     &WireTool{},
	}, nil

}

func (w *Workspace) Update(win *pixelgl.Window, dt float64) {
	w.Camera.Update(win, dt)

	switch {
	case win.JustPressed(pixelgl.Key0):
		w.UsedTool = &Finger{}
	case win.JustPressed(pixelgl.Key1):
		w.UsedTool = &WireTool{}
	case win.JustPressed(pixelgl.Key2):
		w.UsedTool = &Scissors{}
	case win.JustPressed(pixelgl.Key3):
		w.UsedTool = &SwitchTool{}
	case win.JustPressed(pixelgl.Key4):
		w.UsedTool = &LampTool{}
	case win.JustPressed(pixelgl.Key5):
		w.UsedTool = ChipTool{"NOT"}
	case win.JustPressed(pixelgl.Key6):
		w.UsedTool = ChipTool{"AND"}
	case win.JustPressed(pixelgl.Key7):
		w.UsedTool = ChipTool{"OR"}

	case win.JustPressed(pixelgl.KeyS):
		w.CircuitBoard.Save()

	case win.JustPressed(pixelgl.KeyEscape):
		os.Exit(0)
	}

	cm := w.Camera.Matrix(win)
	mp := cm.Unproject(win.MousePosition())

	w.UsedTool.Update(win, w.CircuitBoard, mp)
}

func gridAlign(v pixel.Vec) pixel.Vec {
	return pixel.V(
		math.Floor(v.X/GRID_SIZE+0.5)*GRID_SIZE,
		math.Floor(v.Y/GRID_SIZE+0.5)*GRID_SIZE,
	)
}

func (w *Workspace) Draw(win *pixelgl.Window) {
	// Zoom & move
	cm := w.Camera.Matrix(win)
	win.SetMatrix(cm)

	// Draw grid holes
	w.HoleBatch.Clear()
	wmin := cm.Unproject(win.Bounds().Min)
	wmax := cm.Unproject(win.Bounds().Max)
	for x := math.Ceil(wmin.X/GRID_SIZE) * GRID_SIZE; x <= wmax.X; x += GRID_SIZE {
		for y := math.Ceil(wmin.Y/GRID_SIZE) * GRID_SIZE; y <= wmax.Y; y += GRID_SIZE {
			w.HoleSprite.Draw(w.HoleBatch, pixel.IM.Moved(pixel.V(x, y)))
		}
	}
	w.HoleBatch.Draw(win)

	w.CircuitBoard.Draw(win)

	mp := cm.Unproject(win.MousePosition())
	w.UsedTool.Draw(win, mp)
}
