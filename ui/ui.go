package ui

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type UI struct {
	WindowConfig pixelgl.WindowConfig
	Window       *pixelgl.Window
	Help         *Help
	Workspace    *Workspace
}

func NewUI() (*UI, error) {
	cfg := pixelgl.WindowConfig{
		Title:     "Circuits simulator",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return nil, err
	}
	win.SetSmooth(true)
	help, err := NewHelp()
	if err != nil {
		return nil, err
	}
	workspace, err := NewWorkspace()
	if err != nil {
		return nil, err
	}

	return &UI{
		WindowConfig: cfg,
		Window:       win,
		Help:         help,
		Workspace:    workspace,
	}, nil
}

func (u UI) MainLoop() {
	frames := 0
	second := time.Tick(time.Second)

	last := time.Now()
	for !u.Window.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		u.Workspace.Update(u.Window, dt)
		u.Help.Update(u.Window)
		u.Window.Clear(colornames.Floralwhite)

		u.Workspace.Draw(u.Window)
		u.Help.Draw(u.Window)
		u.Window.Update()

		frames++
		select {
		case <-second:
			u.Window.SetTitle(fmt.Sprintf("%s | FPS: %d", u.WindowConfig.Title, frames))
			frames = 0
		default:
		}
	}
}
