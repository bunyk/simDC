package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bunyk/simDC/geometry"
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

func (gp GridPoint) ToRight() GridPoint {
	return GridPoint{gp.X + 1, gp.Y}
}

type CircuitBoard struct {
	Filename string

	gridSignals map[GridPoint]bool
	WireGroups  map[GridPoint]*WireGroup

	Switches map[GridPoint]bool
	Lamps    map[GridPoint]bool
	Chips    []ChipInstance

	mutex sync.RWMutex
}

func NewCircuitBoard() *CircuitBoard {
	cb := &CircuitBoard{}

	cb.gridSignals = make(map[GridPoint]bool)
	cb.WireGroups = make(map[GridPoint]*WireGroup)
	cb.Switches = make(map[GridPoint]bool)
	cb.Lamps = make(map[GridPoint]bool)

	cb.Filename = "circuit_save.json"
	if len(os.Args) > 1 {
		cb.Filename = os.Args[1]
	}
	cb.Load()
	return cb
}

func (cb *CircuitBoard) SplitWire(pos GridPoint) {
	seenWG := make(map[*WireGroup]bool)
	for _, wg := range cb.WireGroups {
		if seenWG[wg] {
			continue
		}
		for _, wire := range wg.Wires {
			if pos == wire.A || pos == wire.B { // one of the ends
				continue // not interested
			}
			if wire.AsLine().Contains(pos.Pos()) {
				nw := Wire{pos, wire.B}
				wire.B = pos
				wg.Wires = append(wg.Wires, nw)
				cb.SetSignal(pos, wg.Signal)
			}
		}
		seenWG[wg] = true
	}
}

func (cb *CircuitBoard) AddWire(x1, y1, x2, y2 int) {
	a := GridPoint{x1, y1}
	b := GridPoint{x2, y2}
	cb.SplitWire(a)
	cb.SplitWire(b)
	signal := cb.GetSignal(a) || cb.GetSignal(b)

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
		cb.SetSignal(w.A, signal)
		cb.SetSignal(w.B, signal)
	}
}

func (cb *CircuitBoard) AddChip(chip string, x, y int) {
	p := GridPoint{x, y}
	ci := NewChipInstance(chip, p)
	cb.Chips = append(cb.Chips, ci)

	go func(ch ChipInstance) {
		ch.Process(cb)
	}(ci)
}

func (cb *CircuitBoard) Save() {
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

func (cb *CircuitBoard) MarshalJSON() ([]byte, error) {
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
		if sw {
			m.Switches = append(m.Switches, []int{
				pos.X, pos.Y,
			})
		}
	}
	for pos, l := range cb.Lamps {
		if l {
			m.Lamps = append(m.Lamps, []int{
				pos.X, pos.Y,
			})
		}
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
	for _, wire := range mcb.Wires {
		cb.AddWire(wire[0], wire[1], wire[2], wire[3])
	}
	for _, sw := range mcb.Switches {
		cb.AddSwitch(sw[0], sw[1])
	}
	for _, l := range mcb.Lamps {
		cb.AddLamp(l[0], l[1])
	}
	for _, ch := range mcb.Chips {
		cb.AddChip(ch.Class, ch.Location.X, ch.Location.Y)
	}
}

func (cb *CircuitBoard) AddSwitch(x, y int) {
	// TODO: check if place is free
	p := GridPoint{x, y}
	cb.Switches[p] = true
}

func (cb *CircuitBoard) AddLamp(x, y int) {
	// TODO: check if place is free
	p := GridPoint{x, y}
	cb.Lamps[p] = true
}

func (cb *CircuitBoard) PressSwitch(x, y int) {
	p := GridPoint{x, y}
	if !cb.Switches[p] {
		return // No switch here
	}
	p.X += 1 // output is to the right
	// Switch switch
	cb.SetSignal(p, !cb.GetSignal(p))
}

func (cb *CircuitBoard) GetSignal(p GridPoint) bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.gridSignals[p]
}

func (cb *CircuitBoard) SetSignal(p GridPoint, signal bool) {
	// fmt.Printf("cb.SetSignal(%#v, %#v)\n", p, signal)
	cb.mutex.Lock()

	prev := cb.gridSignals[p]
	cb.gridSignals[p] = signal // still explicitly set signal so we see node
	if prev == signal {        // signal not changed
		cb.mutex.Unlock()
		return // no need to update anything
	}

	cb.mutex.Unlock()

	if cb.WireGroups[p] != nil {
		cb.WireGroups[p].SetSignal(signal, cb)
	}
	for _, chip := range cb.Chips {
		if chip.HasInputAt(p) {
			go func(ch ChipInstance) {
				ch.Process(cb)
			}(chip)
		}
	}
}

func (cb *CircuitBoard) Draw(win *pixelgl.Window) {
	//start := time.Now()
	imd := imdraw.New(nil)
	imd.Precision = 10
	for p, sw := range cb.Switches {
		if !sw {
			continue
		}
		drawSwitch(imd, win, p.Pos(), cb.GetSignal(p.ToRight()))
	}

	for _, chip := range cb.Chips {
		chip.Draw(imd, win)
	}
	seenWG := make(map[*WireGroup]bool)
	for _, group := range cb.WireGroups {
		if seenWG[group] {
			continue
		}
		group.Draw(imd)
		seenWG[group] = true
	}

	cb.mutex.RLock()
	for pos, signal := range cb.gridSignals {
		drawNode(imd, pos.Pos(), signal)
	}
	cb.mutex.RUnlock()

	for p, lamp := range cb.Lamps {
		if !lamp {
			continue
		}
		drawLamp(imd, p.Pos(), cb.GetSignal(p))
	}
	imd.Draw(win)
	// fmt.Printf("Rendered %d switches, %d chips, %d wire groups and %d lamps in %s\n",
	// 	len(cb.Switches), len(cb.Chips), len(cb.WireGroups), len(cb.Lamps), time.Since(start),
	// )
}

func (cb *CircuitBoard) CutThrough(a, b pixel.Vec) {
	cut := pixel.L(a, b)
	cb.cutWires(cut)
	cb.cutLamps(cut)
	cb.cutSwitches(cut)
	cb.cutChips(cut)
}

func (cb *CircuitBoard) cutWires(cut pixel.Line) {
	// Gather list of all wires
	var wires []Wire
	seenWG := make(map[*WireGroup]bool)
	for _, wg := range cb.WireGroups {
		if seenWG[wg] {
			continue
		}
		wires = append(wires, wg.Wires...)
		seenWG[wg] = true
	}
	// Cleanup wireGroups
	cb.WireGroups = make(map[GridPoint]*WireGroup)

	// Remove wires that are cut, and add remaining wires back
	for _, wire := range wires {
		if !geometry.LineSegmentsIntersect(wire.AsLine(), cut) {
			cb.AddWire(wire.A.X, wire.A.Y, wire.B.X, wire.B.Y)
		}
	}

}

func (cb *CircuitBoard) cutLamps(cut pixel.Line) {
	for p := range cb.Lamps {
		if geometry.LineCircleIntersect(cut, p.Pos(), LAMP_RADIUS) {
			delete(cb.Lamps, p)
		}
	}
}

func (cb *CircuitBoard) cutSwitches(cut pixel.Line) {
	for p := range cb.Switches {
		pos := p.Pos()
		if geometry.LineRectangleIntersect(cut, pixel.R(pos.X, pos.Y-GRID_SIZE/2, pos.X+GRID_SIZE, pos.Y+GRID_SIZE/2)) {
			delete(cb.Switches, p)
		}
	}
}

func (cb *CircuitBoard) cutChips(cut pixel.Line) {
	chips := make([]ChipInstance, 0, len(cb.Chips))
	for _, c := range cb.Chips {
		pos := c.Location.Pos()
		height := ChipClasses[c.Class].Height()
		if !geometry.LineRectangleIntersect(cut, pixel.R(
			pos.X, pos.Y-GRID_SIZE*(float64(height)-0.5),
			pos.X+GRID_SIZE, pos.Y+GRID_SIZE/2,
		)) {
			chips = append(chips, c)
		}
	}
	cb.Chips = chips
}
