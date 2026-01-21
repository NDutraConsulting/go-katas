package concurrentsensors

import (
	"fmt"
	"sync/atomic"
)

// ---STATE STEPS---
const (
	oshaStateFindH2S      = "FIND_HS2"
	oshaStateTurnOffEquip = "TURN_OFF_EQUIP"
	oshaStateAbsorbH2S    = "ABSORB_H2S"
	oshaStateHoldPosition = "HOLD_POSITION"
)

// ---Initializers---
const (
	oshaRobotTag     = "OSHA_ROBOT"
	robotH2SCapSMALL = 3
	robotH2SCapLARGE = 8
)

type oshaRobot struct {
	oshaSensor
	robotH2SCapacity int
	robotH2SStorage  int
	radiatioLevel    int
	maxRadiation     int
	state            string
	gridMap          [][]gridSpace
	targetSpace      *gridSpace
	currentSpace     *gridSpace
}

// --- GAME LOOP WITH A THREAD BLOCK (GO WORKER)
func (robot *oshaRobot) processEnvironment() {
	go func(rb *oshaRobot) {
		fmt.Println("\nTHREAD START: Processing environment for robot ID: ", rb.getID())

		rb.update()

		for _, obj := range rb.currentSpace.gridObjects {

			if goObj, ok := obj.(gridObject); ok {
				rb.checkCollisions(goObj)
			}

		}

		// Capture robot and obj for the goroutine to avoid loop variable capture
		rb.renderWork()
		rb.recordProgress()
		fmt.Println("THREAD STOP for robot ID: ", rb.getID())
	}(robot)

}

// --- PROCESSING CYCLE FUNCTIONS ---
/**
* (o *oshaRobot) update()
* Updates the robot's environment awareness using sensor data
* and then shifts to the next appropriate state
**/
func (o *oshaRobot) update() {

	fmt.Println("Update ID: ", o.getID(), " ->Expected Type: ", oshaRobotTag)

	surroundingh2sLevel := o.gridMap[o.col()][o.row()].h2sLevel
	if surroundingh2sLevel >= oshaH2SHazard {
		fmt.Println("-------------------> !!!!!! HAZARD H2S detected at surrounding level: ", surroundingh2sLevel)
		o.state = oshaStateAbsorbH2S

	} else if surroundingh2sLevel >= oshaH2SLow {
		fmt.Println("-------------------> H2S detected at surrounding level: ", surroundingh2sLevel)
	}

	// Read onboard sensor data

	// Read distributed sensor data from channel
}

/**
* (o *oshaRobot) checkCollisions(gridObj gridObject)
* Checks for potential collisions with other objects
* and the changes the state of this robot accordingly
**/
func (o *oshaRobot) checkCollisions(gridObj gridObject) {
	fmt.Println("Checking collisions for robot ID: ", o.getID(), " and object: ", gridObj)
}

/**
* (o *oshaRobot) renderWork()
* Renders the work for the robot based on its current state
**/
func (o *oshaRobot) renderWork() {
	fmt.Println("Rendering work for robot ID: ", o.getID())

	robotPosition := o.oshaSensor.grid

	switch o.state {
	case oshaStateFindH2S:
		fmt.Println("Rendering work for robot ID: ", o.getID(), " in oshaStateFindH2S")
		if o.targetSpace != nil {
			if o.targetSpace.gridPosition.col == robotPosition.col && o.targetSpace.row == robotPosition.row {
				o.state = oshaStateAbsorbH2S
				return
			}

			pos := o.oshaSensor.grid
			if o.targetSpace.gridPosition.col > robotPosition.col {
				fmt.Println("Moving horizontally...")
				pos.col++
			} else if o.targetSpace.gridPosition.col < robotPosition.col {
				pos.col--
			}

			if o.targetSpace.row > robotPosition.row {
				fmt.Println("Moving vertically...")
				pos.row++
			} else if o.targetSpace.row < robotPosition.row {
				pos.row--
			}
			o.move(pos.row, pos.col)
		}

	case oshaStateTurnOffEquip:
		fmt.Println("Rendering work for robot ID: ", o.getID(), " in oshaStateTurnOffEquip")
	case oshaStateAbsorbH2S:
		fmt.Println("Rendering work for robot ID: ", o.getID(), " in oshaStateAbsorbH2S")

		h2sLevel := o.currentSpace.h2sLevel

		if h2sLevel > 0 && o.robotH2SStorage < o.robotH2SCapacity {
			fmt.Println("Absorbing H2S at level: ", h2sLevel)
			o.robotH2SStorage += 1
			o.currentSpace.h2sLevel--
		}

	case oshaStateHoldPosition:
		fmt.Println("Rendering work for robot ID: ", o.getID(), " in oshaStateHoldPosition")
	}
}

/**
* (o *oshaRobot) recordProgress()
* Records the progress of the robot in its memory
* and sends to the server storage
**/
func (o *oshaRobot) recordProgress() {
	fmt.Println("Recording progress for robot ID: ", o.getID())
}

// --- UNIQUE ID COUNTER FOR OSHA ROBOT ---
var oshaRobotCounter int64

// --- Constructor Function ---
func NewOshaRobot(row, col, h2sCapacity int, gridMap [][]gridSpace) oshaRobot {
	return oshaRobot{
		// Increment and get the new value atomically
		oshaSensor: oshaSensor{
			gridObject: gridObject{
				ID:    atomic.AddInt64(&oshaRobotCounter, 1),
				grid:  gridPosition{row: row, col: col},
				label: oshaRobotTag,
			},
		},

		robotH2SCapacity: h2sCapacity,
		robotH2SStorage:  0,
		maxRadiation:     100,
		state:            oshaStateFindH2S,
		gridMap:          gridMap,
		currentSpace:     &gridMap[row][col],
	}
}

// --- HELPER FUNCTIONS ---

func (o *oshaRobot) setTargetSpace(targetSpace *gridSpace) {
	o.targetSpace = targetSpace
}

func (o *oshaRobot) move(row, col int) {

	delete(o.currentSpace.oshaRobots, o.getID())

	o.oshaSensor.grid = gridPosition{row: row, col: col}

	o.currentSpace = &o.gridMap[row][col]

	o.currentSpace.oshaRobots[o.getID()] = o
}
