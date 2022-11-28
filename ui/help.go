package ui

import (
	"fmt"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type Help struct {
	Hidden bool
	text   *text.Text
}

const HELP_TEXT = `
F1 - Hide/show help text

Arrow keys - move visible area
Mouse scroll - zoom visible area

    Tools:
0 - Turn on power & click switches
1 - Add wires
2 - Cut wires
3 - Switch
4 - Indicator (aka LED, Lamp)
5 - NOT
6 - AND
7 - OR

	Actions:
S - Save
`

var ADD_WIRE_KEY = pixelgl.Key1

func NewHelp() (*Help, error) {
	txt := text.New(pixel.ZV, FontAtlas)
	txt.Color = colornames.Green
	fmt.Fprint(txt, HELP_TEXT)

	return &Help{text: txt}, nil
}

func (h *Help) Update(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeyF1) {
		h.Hidden = !h.Hidden
	}
}

func (h Help) Draw(win *pixelgl.Window) {
	if h.Hidden {
		return
	}
	win.SetMatrix(pixel.IM)
	h.text.Draw(win, pixel.IM.Scaled(pixel.ZV, 2).Moved(pixel.V(10, win.Bounds().Max.Y)))
}
