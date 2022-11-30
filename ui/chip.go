package ui

import (
	"fmt"
	"time"

	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type ChipClass struct {
	InputsCount  int
	OutputsCount int
	Delay        time.Duration
	Logic        func([]bool) []bool
}

func (cc ChipClass) Height() int {
	return max(cc.InputsCount, cc.OutputsCount)
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
}

func NewChipInstance(class string, pos GridPoint) ChipInstance {
	return ChipInstance{
		Class:    class,
		Location: pos,
	}
}

func (ci ChipInstance) HasInputAt(pos GridPoint) bool {
	if pos.X != ci.Location.X {
		return false
	}
	index := ci.Location.Y - pos.Y
	fmt.Println("index=", index)
	if index < 0 {
		return false
	}
	if index >= ChipClasses[ci.Class].InputsCount {
		return false
	}
	fmt.Printf("%s has input at %#v\n", ci.Class, pos)
	return true
}

func (ci *ChipInstance) Process(cb *CircuitBoard) {
	class := ChipClasses[ci.Class]
	// Get input
	inputs := make([]bool, class.InputsCount)
	for i := range inputs {
		inputs[i] = cb.GetSignal(GridPoint{
			ci.Location.X, ci.Location.Y - i,
		})
	}
	fmt.Printf("%s at %#v will be updated using %#v inputs\n", ci.Class, ci.Location, inputs)
	outputs := class.Logic(inputs) // compute result

	time.Sleep(class.Delay) // but slowly

	fmt.Printf("%s at %#v produces output %#v\n", ci.Class, ci.Location, outputs)

	// Produce output
	for i, out := range outputs {
		cb.SetSignal(GridPoint{
			ci.Location.X + 1,
			ci.Location.Y - i,
		}, out)
	}
}

func (ci ChipInstance) Draw(imd *imdraw.IMDraw, win *pixelgl.Window) {
	drawChip(imd, ci.Location.Pos(), ci.Class, ChipClasses[ci.Class].Height())
}
