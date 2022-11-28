package main

import (
	"github.com/bunyk/simcomp/ui"
	"github.com/faiface/pixel/pixelgl"
)

func run() {
	u, err := ui.NewUI()
	if err != nil {
		panic(err)
	}
	u.MainLoop()
}

func main() {
	pixelgl.Run(run)
}
