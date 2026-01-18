package robotfarm

import "fmt"

type orangePos struct {
	col   int
	row   int
	state string
}

func orangesRotting(grid [][]int) int {

	parallelQueue := make([]orangePos, 0)

	freshOranges := 0
	robots := 0

	MAX_COL := len(grid)
	MAX_ROW := len(grid[0])

	for col := 0; col < len(grid); col++ {
		for row := 0; row < len(grid[0]); row++ {

			switch grid[col][row] {
			case 2:
				parallelQueue = append(parallelQueue, orangePos{col: col, row: row, state: "rotten"})
			case 3:
				parallelQueue = append(parallelQueue, orangePos{col: col, row: row, state: "robot"})
				robots++
			case 1:
				freshOranges++
			}
		}
	}

	freshOrangesStart := freshOranges

	parallelQueue = append(parallelQueue, orangePos{col: -1, row: -1, state: "end"})

	// Why not 0?
	elapsedTime := -1
	directions := []orangePos{
		{col: -1, row: 0, state: "back"}, {col: 0, row: -1, state: "up"},
		{col: 1, row: 0, state: "left"}, {col: 0, row: 1, state: "down"}}
	for {
		if len(parallelQueue) == robots {
			break
		}

		// POP the top orangePos
		currentOrange := parallelQueue[0]
		parallelQueue = parallelQueue[1:]

		if currentOrange.state == "end" {
			if len(parallelQueue) > robots {
				parallelQueue = append(parallelQueue, orangePos{col: -1, row: -1, state: "end"})
			}

			elapsedTime++
			continue
		} else if currentOrange.state == "robot" {
			// Remove rotten fruit around it -or-
			// Spray neem oil to prevent spread
			for _, dir := range directions {

				colDir, rowDir := currentOrange.col+dir.col, currentOrange.row+dir.row

				// spread the rotten in valid directions only
				if MAX_COL > colDir && colDir >= 0 && MAX_ROW > rowDir && rowDir >= 0 {

					if grid[colDir][rowDir] == 2 {
						grid[colDir][rowDir] = 0

						// search the queue and change the state to "cleared"
						for i, q := range parallelQueue {
							if q.col == colDir && q.row == rowDir {
								parallelQueue[i].state = "cleared"
								break
							}
						}
					}
				}

			}
			parallelQueue = append(parallelQueue, currentOrange)

		} else if currentOrange.state == "rotten" {
			for _, dir := range directions {

				colDir, rowDir := currentOrange.col+dir.col, currentOrange.row+dir.row

				// spread the rotten in valid directions only
				if MAX_COL > colDir && colDir >= 0 && MAX_ROW > rowDir && rowDir >= 0 {

					if grid[colDir][rowDir] == 1 {
						grid[colDir][rowDir] = 2
						freshOranges--
						parallelQueue = append(parallelQueue, orangePos{col: colDir, row: rowDir, state: "rotten"})
					}
				}

			}

		}
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

func printGrid(grid [][]int) {

	for _, row := range grid {
		fmt.Println(row)
	}
}
