package robotfarm

import "fmt"

// ============== Constants & Types ==============

// Use constants instead of strings for state
const (
	StateRotten  = "rotten"
	StateRobot   = "robot"
	StateEnd     = "end"
	StateCleared = "cleared"
)

// Grid cell values
const (
	CellEmpty  = 0
	CellFresh  = 1
	CellRotten = 2
	CellRobot  = 3
)

type position struct {
	col int
	row int
}

type queueItem struct {
	position
	state string
}

var rottenDirections = []position{
	{col: -1, row: 0}, // left
	{col: 1, row: 0},  // right
	{col: 0, row: -1}, // up
	{col: 0, row: 1},  // down
}

var robotDirections = []position{
	{col: -1, row: 0},  // left
	{col: 1, row: 0},   // right
	{col: 0, row: -1},  // up
	{col: 0, row: 1},   // down
	{col: 1, row: 1},   // down-right
	{col: -1, row: 1},  // down-left
	{col: 1, row: -1},  // up-right
	{col: -1, row: -1}, // up-left
}

// ============== Orchard (main state container) ==============

type orchard struct {
	grid         [][]int
	queue        []queueItem
	maxCol       int
	maxRow       int
	freshOranges int
	robots       int
}

func newOrchard(grid [][]int) *orchard {
	o := &orchard{
		grid:   grid,
		queue:  make([]queueItem, 0),
		maxCol: len(grid),
		maxRow: len(grid[0]),
	}
	o.scan()
	return o
}

// Scan grid and initialize queue
func (o *orchard) scan() {
	for col := 0; col < o.maxCol; col++ {
		for row := 0; row < o.maxRow; row++ {
			switch o.grid[col][row] {
			case CellRotten:
				o.enqueue(col, row, StateRotten)
			case CellRobot:
				o.enqueue(col, row, StateRobot)
				o.robots++
			case CellFresh:
				o.freshOranges++
			}
		}
	}
}

// ============== Queue Operations ==============

func (o *orchard) enqueue(col, row int, state string) {
	o.queue = append(o.queue, queueItem{position{col, row}, state})
}

func (o *orchard) dequeue() queueItem {
	item := o.queue[0]
	o.queue = o.queue[1:]
	return item
}

func (o *orchard) addEndMarker() {
	o.enqueue(-1, -1, StateEnd)
}

func (o *orchard) hasWork() bool {
	return len(o.queue) > o.robots
}

// ============== Position Helpers ==============

func (o *orchard) isValid(p position) bool {
	return p.col >= 0 && p.col < o.maxCol && p.row >= 0 && p.row < o.maxRow
}

func (o *orchard) cellAt(p position) int {
	return o.grid[p.col][p.row]
}

func (o *orchard) setCell(p position, value int) {
	o.grid[p.col][p.row] = value
}

func (o *orchard) neighbors(p position, directions []position) []position {
	result := make([]position, 0, 4)
	for _, dir := range directions {
		neighbor := position{col: p.col + dir.col, row: p.row + dir.row}
		if o.isValid(neighbor) {
			result = append(result, neighbor)
		}
	}
	return result
}

// ============== State Handlers ==============

func (o *orchard) processRotten(item queueItem) {
	for _, neighbor := range o.neighbors(item.position, rottenDirections) {
		if o.cellAt(neighbor) == CellFresh {
			o.setCell(neighbor, CellRotten)
			o.freshOranges--
			o.enqueue(neighbor.col, neighbor.row, StateRotten)
		}
	}
}

func (o *orchard) processRobot(item queueItem) {
	for _, neighbor := range o.neighbors(item.position, robotDirections) {
		if o.cellAt(neighbor) == CellRotten {
			o.setCell(neighbor, CellEmpty)
			o.markCleared(neighbor)
		}
	}
	// Robots persist in queue
	o.enqueue(item.col, item.row, StateRobot)
}

func (o *orchard) markCleared(p position) {
	for i := range o.queue {
		if o.queue[i].col == p.col && o.queue[i].row == p.row {
			o.queue[i].state = StateCleared
			break
		}
	}
}

// ============== Debug Helpers ==============

func (o *orchard) printStats(startFresh, elapsed int) {
	fmt.Println("Elapsed Time:", elapsed)
	fmt.Printf("Oranges Start: %d, Remaining: %d\n", startFresh, o.freshOranges)
	if startFresh > 0 {
		fmt.Printf("Percent Remaining: %.2f%%\n", float64(o.freshOranges)/float64(startFresh)*100)
	}
	fmt.Println("Oranges Lost:", startFresh-o.freshOranges)
	o.printGrid()
}

func (o *orchard) printGrid() {
	for _, row := range o.grid {
		fmt.Println(row)
	}
}
