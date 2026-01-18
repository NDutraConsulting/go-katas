package robotfarm

import "fmt"

const debug = false

// Constants for state
const (
	StateRotten  = "rotten"
	StateRobot   = "robot"
	StateEnd     = "end"
	StateCleared = "cleared"
	StateFallen  = "fallen"
)

// Grid cell values
const (
	CellEmpty  = 0
	CellFresh  = 1
	CellRotten = 2
	CellRobot  = 3
)

type orangePos struct {
	col   int
	row   int
	state string
}

func orangesRotting(grid [][]int) int {

	fmt.Println("\n-START-")
	printGrid(grid)

	parallelQueue := make([]orangePos, 0)

	freshOranges := 0
	robots := 0

	MAX_COL := len(grid)
	MAX_ROW := len(grid[0])

	for col := 0; col < MAX_COL; col++ {
		for row := 0; row < MAX_ROW; row++ {

			switch grid[col][row] {
			case CellRotten:
				parallelQueue = append(parallelQueue, orangePos{col: col, row: row, state: StateRotten})
			case CellRobot:
				parallelQueue = append(parallelQueue, orangePos{col: col, row: row, state: StateRobot})
				robots++
			case CellFresh:
				freshOranges++
			}
		}
	}

	freshOrangesStart := freshOranges

	parallelQueue = append(parallelQueue, orangePos{col: -1, row: -1, state: StateEnd})

	// Why not 0?
	elapsedTime := -1
	rottenDirections := rottenDirections()

	// Introduces additional steps?
	// robotDirections := robotDirections()

	for {
		// STOP BFS if queue is empty or only contains only robots
		// Robots are never removed from the queue because the persist
		// This means that the min length of the queue will be the total number of robots
		if len(parallelQueue) == robots {
			break
		}

		// POP the top orangePos
		currentOrange := parallelQueue[0]
		parallelQueue = parallelQueue[1:]

		printDebug(currentOrange)

		// StateBundle might create too much complexity
		orchardState := stateBundle{
			parallelQueue: &parallelQueue,
			currentOrange: &currentOrange,
			MAX_COL:       MAX_COL,
			MAX_ROW:       MAX_ROW,
			freshOranges:  &freshOranges}

		switch currentOrange.state {
		case StateEnd:
			if len(parallelQueue) > robots {
				parallelQueue = append(parallelQueue, orangePos{col: -1, row: -1, state: StateEnd})
			}

			elapsedTime++
			continue
		case StateRobot:
			robot(orchardState, rottenDirections, grid)
		case StateRotten:
			rotten(orchardState, rottenDirections, grid)
		}
		printDebug(currentOrange)

	}

	fmt.Println("Elapsed Time: ", elapsedTime)
	fmt.Printf("Oranges Start: %d, Oranges Remaining: %d\n", freshOrangesStart, freshOranges)
	fmt.Println("Percent Remaining:", (float64(freshOranges) / float64(freshOrangesStart)), "%")
	fmt.Println("Oranges Lost:", float64(freshOrangesStart)-float64(freshOranges))
	printGrid(grid)

	result := elapsedTime
	if freshOranges > 0 {
		result = -1
	}

	return result
}

// State handlers
func robot(sb stateBundle, directions []gridObject, grid [][]int) {
	parallelQueue, currentOrange, MAX_COL, MAX_ROW, _ := sb.destructure()

	// Remove rotten fruit around it -or-
	// Spray neem oil to prevent spread
	for _, dir := range directions {

		colDir, rowDir := currentOrange.col+dir.col, currentOrange.row+dir.row

		// spread the rotten in valid directions only
		if MAX_COL > colDir && colDir >= 0 && MAX_ROW > rowDir && rowDir >= 0 {

			if grid[colDir][rowDir] == 2 {
				grid[colDir][rowDir] = 0

				// search the queue and change the state to "cleared"
				for i, q := range *parallelQueue {
					if q.col == colDir && q.row == rowDir {
						(*parallelQueue)[i].state = StateCleared
						break
					}
				}
			}
		}

	}
	*parallelQueue = append(*parallelQueue, *currentOrange)
}

func rotten(sb stateBundle, directions []gridObject, grid [][]int) {

	parallelQueue, currentOrange, MAX_COL, MAX_ROW, freshOranges := sb.destructure()

	for _, dir := range directions {

		colDir, rowDir := currentOrange.col+dir.col, currentOrange.row+dir.row

		// spread the rotten in valid directions only
		if MAX_COL > colDir && colDir >= 0 && MAX_ROW > rowDir && rowDir >= 0 {

			if grid[colDir][rowDir] == 1 {
				grid[colDir][rowDir] = 2
				*freshOranges--
				*parallelQueue = append(*parallelQueue, orangePos{col: colDir, row: rowDir, state: StateRotten})
			}
		}

	}

	// Pointer state setting check
	if debug {
		currentOrange.state = StateFallen
		printDebug(*currentOrange)
	}

}

func rottenDirections() []gridObject {
	return []gridObject{
		{col: -1, row: 0, direction: "back"}, {col: 0, row: -1, direction: "up"},
		{col: 1, row: 0, direction: "left"}, {col: 0, row: 1, direction: "down"}}
}

func robotDirections() []gridObject {
	return []gridObject{
		{col: -1, row: 0, direction: "back"}, {col: 0, row: -1, direction: "up"},
		{col: 1, row: 0, direction: "left"}, {col: 0, row: 1, direction: "down"},
		{col: 1, row: 1, direction: "down-right"}, {col: -1, row: 1, direction: "down-left"},
		{col: -1, row: -1, direction: "up-left"}, {col: 1, row: -1, direction: "up-right"}}
}

func printGrid(grid [][]int) {

	for _, row := range grid {
		fmt.Println(row)
	}
}

func printDebug(cOrange orangePos) {
	if debug {
		fmt.Println(cOrange)
	}
}
