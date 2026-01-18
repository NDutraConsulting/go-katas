package robotfarm

import "fmt"

// ============== Public API ==============
func OrangesRotting(grid [][]int) int {
	fmt.Println("\n-START-")
	o := newOrchard(grid)
	startFresh := o.freshOranges

	o.addEndMarker()
	elapsedTime := -1

	o.printGrid()

	// ============== Main BFS ==============
	for o.hasWork() {
		current := o.dequeue()

		if debug {

			fmt.Println("Processing:", current)
			o.printGrid()
		}

		switch current.state {
		case StateEnd:
			elapsedTime++
			if o.hasWork() {
				o.addEndMarker()
				if debug {
					fmt.Println("-next minute starts-")
				}
			}

		case StateRotten:
			o.processRotten(current)

		case StateRobot:
			o.processRobot(current)

		case StateCleared:
			// Skip cleared items
		}
	}

	o.printStats(startFresh, elapsedTime)

	if o.freshOranges > 0 {
		return -1
	}
	return elapsedTime
}
