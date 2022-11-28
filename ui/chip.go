package ui

import (
	"fmt"
	"time"

	"github.com/faiface/pixel/pixelgl"
)

type ChipClass struct {
	Title        string
	InputsCount  int
	OutputsCount int
	Delay        time.Duration
	Logic        func([]bool) []bool
}

var ChipClasses = map[string]ChipClass{
	"NOT": ChipClass{
		InputsCount:  1,
		OutputsCount: 1,
		Delay:        time.Second / 4,
		Logic: func(inputs []bool) []bool {
			return []bool{!inputs[0]}
		},
	},
	"AND": ChipClass{
		InputsCount:  2,
		OutputsCount: 1,
		Delay:        time.Second / 2,
		Logic: func(inputs []bool) []bool {
			return []bool{inputs[0] && inputs[1]}
		},
	},
	"OR": ChipClass{
		InputsCount:  2,
		OutputsCount: 1,
		Delay:        time.Second / 2,
		Logic: func(inputs []bool) []bool {
			return []bool{inputs[0] || inputs[1]}
		},
	},
}

type ChipInstance struct {
	Class    string
	Location GridPoint
	Inputs   []bool
	Outputs  []bool
}

func NewChipInstance(class string, pos GridPoint) ChipInstance {
	return ChipInstance{
		Class:    class,
		Location: pos,
		Inputs:   make([]bool, ChipClasses[class].InputsCount),
		Outputs:  make([]bool, ChipClasses[class].OutputsCount),
	}
}

func (ci *ChipInstance) SetInputSignal(index int, signal bool) {
	if index < 0 || index >= len(ci.Inputs) {
		fmt.Printf("Error setting signal with index %d for %s, it has %d inputs\n", index, ci.Class, len(ci.Inputs))
		return
	}
	ci.Inputs[index] = signal
	time.Sleep(ChipClasses[ci.Class].Delay)
	ci.Outputs = ChipClasses[ci.Class].Logic(ci.Inputs)
}

func (ci ChipInstance) Draw(win *pixelgl.Window) {
	drawChip(
		win,
		ci.Location.Pos(),
		ci.Class,
		ci.Inputs,
		ci.Outputs,
	)
}
