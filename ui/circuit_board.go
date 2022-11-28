package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type GridPoint struct {
	X int
	Y int
}

func (gp GridPoint) Pos() pixel.Vec {
	return pixel.V(
		float64(gp.X)*GRID_SIZE,
		float64(gp.Y)*GRID_SIZE,
	)
}

type CircuitBoard struct {
	Filename   string
	WireGroups map[GridPoint]*WireGroup
	Switches   map[GridPoint]*bool
	Lamps      map[GridPoint]*bool
	Chips      []ChipInstance

	mutex sync.Mutex
}

func NewCircuitBoard() *CircuitBoard {
	cb := &CircuitBoard{}
	cb.WireGroups = make(map[GridPoint]*WireGroup)
	cb.Switches = make(map[GridPoint]*bool)
	cb.Lamps = make(map[GridPoint]*bool)
	cb.Filename = "circuit_save.json"
	if len(os.Args) > 1 {
		cb.Filename = os.Args[1]
	}
	cb.Load()
	return cb
}

func (cb *CircuitBoard) AddWire(x1, y1, x2, y2 int, signal bool) {
	a := GridPoint{x1, y1}
	b := GridPoint{x2, y2}

	ng := WireGroup{ // Create new group with this wire
		Wires:  []Wire{{A: a, B: b}},
		Signal: signal,
	}
	if cb.WireGroups[a] != nil { // add elements from wires on first end of current wire
		ng.Wires = append(ng.Wires, cb.WireGroups[a].Wires...)
	}
	if cb.WireGroups[b] != nil && cb.WireGroups[b] != cb.WireGroups[a] { // and on second end
		ng.Wires = append(ng.Wires, cb.WireGroups[b].Wires...)
	}
	// Now make all grid points in group point to same group
	for _, w := range ng.Wires {
		cb.WireGroups[w.A] = &ng
		cb.WireGroups[w.B] = &ng
	}
}

func (cb *CircuitBoard) AddChip(chip string, x, y int) {
	p := GridPoint{x, y}
	cb.Chips = append(cb.Chips, NewChipInstance(chip, p))

}

func (cb CircuitBoard) Save() {
	data, err := json.Marshal(cb)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.WriteFile(cb.Filename, data, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Saved to", cb.Filename)
}

type marshallableCB struct {
	Wires    [][]int
	Switches [][]int
	Lamps    [][]int
	Chips    []ChipInstance
}

func (cb CircuitBoard) MarshalJSON() ([]byte, error) {
	var m marshallableCB

	seenWG := make(map[*WireGroup]bool)
	for _, wg := range cb.WireGroups {
		if seenWG[wg] {
			continue
		}
		for _, wire := range wg.Wires {
			m.Wires = append(m.Wires, []int{
				wire.A.X, wire.A.Y,
				wire.B.X, wire.B.Y,
			})
		}
		seenWG[wg] = true
	}
	for pos, sw := range cb.Switches {
		if sw == nil {
			continue
		}
		m.Switches = append(m.Switches, []int{
			pos.X, pos.Y,
		})
	}
	for pos, l := range cb.Lamps {
		if l == nil {
			continue
		}
		m.Lamps = append(m.Lamps, []int{
			pos.X, pos.Y,
		})
	}
	m.Chips = cb.Chips
	return json.Marshal(m)
}

func (cb *CircuitBoard) Load() {
	fmt.Println("Loading from", cb.Filename)
	data, err := os.ReadFile(cb.Filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	var mcb marshallableCB
	err = json.Unmarshal(data, &mcb)
	if err != nil {
		fmt.Println(err)
		return
	}
	cb.Chips = mcb.Chips
	for _, wire := range mcb.Wires {
		cb.AddWire(wire[0], wire[1], wire[2], wire[3], false)
	}
	for _, sw := range mcb.Switches {
		cb.AddSwitch(sw[0], sw[1])
	}
	for _, l := range mcb.Lamps {
		cb.AddLamp(l[0], l[1])
	}
}

func (cb *CircuitBoard) AddSwitch(x, y int) {
	p := GridPoint{x, y}
	if cb.Switches[p] != nil {
		return
	}
	state := false
	cb.Switches[p] = &state
}

func (cb *CircuitBoard) AddLamp(x, y int) {
	p := GridPoint{x, y}
	if cb.Lamps[p] != nil {
		return
	}
	state := false
	cb.Lamps[p] = &state
}

func (cb *CircuitBoard) PressSwitch(x, y int) {
	// Switch switch
	p := GridPoint{x, y}
	if cb.Switches[p] == nil {
		return // No switch here
	}
	signal := !*cb.Switches[p] // new signal
	*cb.Switches[p] = signal

	// Switch wire
	p.X += 1 // they are connected to the right of switch
	cb.SetWireSignal(p, signal)
}

func (cb *CircuitBoard) SetWireSignal(p GridPoint, signal bool) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	if cb.WireGroups[p] != nil {
		cb.WireGroups[p].SetSignal(signal, cb)
	}
}

func (cb *CircuitBoard) SetElementSignal(p GridPoint, signal bool) {
	if cb.Lamps[p] != nil {
		cb.Lamps[p] = &signal
	}
	for _, chip := range cb.Chips {
		// check if we have input in this coordinates
		inputPinIndex := chip.Location.Y - p.Y
		if (chip.Location.X == p.X) && (inputPinIndex >= 0) && (inputPinIndex < ChipClasses[chip.Class].InputsCount) {
			go func(chip ChipInstance) {
				chip.SetInputSignal(inputPinIndex, signal)
				for i, out := range chip.Outputs {
					cb.SetWireSignal(GridPoint{
						chip.Location.X + 1,
						chip.Location.Y - i,
					}, out)
				}
			}(chip)
		}
	}
}

func (cb CircuitBoard) Draw(win *pixelgl.Window) {
	//start := time.Now()
	imd := imdraw.New(nil)
	for p, sw := range cb.Switches {
		if sw == nil {
			continue
		}
		drawSwitch(win, p.Pos(), *sw)
	}

	for _, chip := range cb.Chips {
		chip.Draw(win)
	}
	seenWG := make(map[*WireGroup]bool)
	for _, group := range cb.WireGroups {
		if seenWG[group] {
			continue
		}
		group.Draw(imd)
		seenWG[group] = true
	}

	for p, lamp := range cb.Lamps {
		if lamp == nil {
			continue
		}
		drawLamp(imd, p.Pos(), *lamp)
	}
	imd.Draw(win)
	// fmt.Printf("Rendered %d switches, %d chips, %d wire groups and %d lamps in %s\n",
	// 	len(cb.Switches), len(cb.Chips), len(cb.WireGroups), len(cb.Lamps), time.Since(start),
	// )
}

func (cb *CircuitBoard) CutWires(a, b pixel.Vec) {
	// Gather list of all wires, with their signals
	var wires []Wire
	var signals []bool
	seenWG := make(map[*WireGroup]bool)
	for _, wg := range cb.WireGroups {
		if seenWG[wg] {
			continue
		}
		wires = append(wires, wg.Wires...)
		for range wg.Wires {
			signals = append(signals, wg.Signal)
		}
		seenWG[wg] = true
	}
	// Cleanup wireGroups
	cb.WireGroups = make(map[GridPoint]*WireGroup)

	// Remove wires that are cut, and add remaining wires back
	for i, wire := range wires {
		c := wire.A.Pos()
		d := wire.B.Pos()
		if !lineSegmentsIntersect(c.X, c.Y, d.X, d.Y, a.X, a.Y, b.X, b.Y) {
			cb.AddWire(wire.A.X, wire.A.Y, wire.B.X, wire.B.Y, signals[i])
		}
	}

}

func (cb *CircuitBoard) CutLamps(a, b pixel.Vec) {
	for p := range cb.Lamps {
		pos := p.Pos()
		if lineCircleIntersect(a.X, a.Y, b.X, b.Y, pos.X, pos.Y, LAMP_RADIUS) {
			delete(cb.Lamps, p)
		}
	}
}

func (cb *CircuitBoard) CutSwitches(a, b pixel.Vec) {
	for p := range cb.Switches {
		pos := p.Pos()
		if lineRectangleIntersect(a.X, a.Y, b.X, b.Y, pos.X, pos.Y-GRID_SIZE/2, pos.X+GRID_SIZE, pos.Y+GRID_SIZE/2) {
			delete(cb.Switches, p)
		}
	}
}

func (cb *CircuitBoard) CutChips(a, b pixel.Vec) {
	chips := make([]ChipInstance, 0, len(cb.Chips))
	for _, c := range cb.Chips {
		pos := c.Location.Pos()
		height := max(len(c.Inputs), len(c.Outputs))
		if !lineRectangleIntersect(
			a.X, a.Y, b.X, b.Y,
			pos.X, pos.Y-GRID_SIZE*(float64(height)-0.5),
			pos.X+GRID_SIZE, pos.Y+GRID_SIZE/2,
		) {
			chips = append(chips, c)
		}
	}
	cb.Chips = chips
}
