package concurrentsensors

import (
	"fmt"
	"time"
)

func makeBoard(ROW, COL int) [][]gridSpace {

	oilFeild := make([][]gridSpace, ROW)
	for i := range oilFeild {
		oilFeild[i] = make([]gridSpace, COL)
		for j := range COL {
			oilFeild[i][j] = gridSpace{
				gridPosition:     gridPosition{row: j, col: i},
				gridObjects:      []gridObjectTable{},
				oshaRobots:       map[string]*oshaRobot{},
				oshaSensors:      []oshaSensor{},
				h2sLevel:         0,
				h2sPocketVolume:  0,
				radioactiveLevel: 0,
			}
		}
	}

	return oilFeild

}

func sim1() bool {

	ROW, COL := 3, 3
	oilFeild := makeBoard(ROW, COL)

	printGrid(oilFeild)

	oilWells := []gridObjectTable{
		gridObject{
			grid:  gridPosition{row: 0, col: 0},
			ID:    1,
			label: typeOilWell},
		gridObject{
			grid:  gridPosition{row: 1, col: 2},
			ID:    2,
			label: typeOilWell},
	}

	for _, well := range oilWells {
		gridObjects := oilFeild[well.col()][well.row()].gridObjects
		oilFeild[well.col()][well.row()].gridObjects = append(gridObjects, well)
		oilFeild[well.col()][well.row()].h2sPocketVolume = 100
	}

	oshaRobotsInit := []oshaRobot{
		NewOshaRobot(0, 1, robotH2SCapLARGE, oilFeild),
		NewOshaRobot(2, 2, robotH2SCapSMALL, oilFeild),
	}

	for _, robot := range oshaRobotsInit {
		oilFeild[robot.col()][robot.row()].oshaRobots[robot.getID()] = &robot

		// Start with all robots searching for H2S at 0,0
		robot.setTargetSpace(&oilFeild[0][0])
	}

	minutes := 1
	simulateOperations(oilFeild, minutes)

	return false
}

func simulateOperations(grid [][]gridSpace, minutes int) {

	seconds := minutes * 10

	// Simulate sensor data ingestion from the sensors and robots at the mining operation
	for i := range seconds {
		fmt.Println("\n\nTime Domain - Tick: ", i+1)
		for r := range grid {
			fmt.Println("\n|----------- Process Spaces in ROW: ", r, "----------|")
			for c := range grid[r] {

				space := &grid[r][c]

				if len(space.gridObjects) == 0 {
					continue
				}

				fmt.Println("START Blocking thread for setting up the Environtment")
				// Process environment objects
				for _, obj := range space.gridObjects {

					// read incoming event streams
					obj.processEnvironment(space)

				}
				fmt.Println("STOP Blocking thread for setting up the Environtment")

			}

			for c := range grid[r] {

				space := &grid[r][c]
				// Process sensors
				for _, sensor := range space.oshaSensors {
					sensor.processEnvironment(space)
				}

				// Process robots
				for _, robot := range space.oshaRobots {
					robot.processEnvironment()
				}

			}
		}

		time.Sleep(30 * time.Millisecond)
	}
}

func printGrid(grid [][]gridSpace) {

	for _, row := range grid {
		fmt.Println("|-------------------------------------------------------------------------------|")
		fmt.Print("| ")
		for _, space := range row {
			fmt.Print(space)
			fmt.Print(" | ")
		}

		fmt.Println()
	}
	fmt.Println("|-------------------------------------------------------------------------------|")

}
