package concurrentsensors

import (
	"fmt"
	"sync/atomic"
)

// ---STATE STEPS---
const (
	STATE_MOVE           = "MOVE"
	STATE_FIND_HS2       = "FIND_HS2"
	STATE_TURN_OFF_EQUIP = "TURN_OFF_EQUIP"
	STATE_ABSORB_H2S     = "ABSORB_H2S"
	STATE_HOLD_POSITION  = "HOLD_POSITION"
)

// ---Initializers---
const (
	OSHA_ROBOT     = "OSHA_ROBOT"
	H2S_CAP_SMALL  = 3
	H2S_CAP_MEDIUM = 5
	H2S_CAP_LARGE  = 8
)

type oshaRobot struct {
	oshaSensor
	H2S_CAPACITY    int
	H2S_Storage     int
	radiation_level int
	MAX_RADIATION   int
	state           string
	gridMap         [][]gridSpace
	targetSpace     *gridSpace
	currentSpace    *gridSpace
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
* update()
* Updates the robot's environment awareness using sensor data
* and then shifts to the next appropriate state
**/
func (o *oshaRobot) update() {

	fmt.Println("Update ID: ", o.getID(), " ->Expected Type: ", OSHA_ROBOT)

	surroundingH2SLevel := o.gridMap[o.col()][o.row()].H2S_Level
	if surroundingH2SLevel >= H2S_HAZARD {
		fmt.Println("-------------------> !!!!!! HAZARD H2S detected at surrounding level: ", surroundingH2SLevel)
		o.state = STATE_ABSORB_H2S

	} else if surroundingH2SLevel >= H2S_LOW {
		fmt.Println("-------------------> H2S detected at surrounding level: ", surroundingH2SLevel)
	}

	// Read onboard sensor data

	// Read distributed sensor data from channel
}

/**
* checkCollisions(gridObj gridObject)
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
	case STATE_FIND_HS2:
		fmt.Println("Rendering work for robot ID: ", o.getID(), " in STATE_FIND_HS2")
		if o.targetSpace != nil {
			if o.targetSpace.gridPosition.col == robotPosition.col && o.targetSpace.row == robotPosition.row {
				o.state = STATE_ABSORB_H2S
				return
			}

			pos := o.oshaSensor.grid
			if o.targetSpace.gridPosition.col > robotPosition.col {
				fmt.Println("Moving horizontally...")
				pos.col++
			} else {
				pos.col--
			}

			if o.targetSpace.row > robotPosition.row {
				fmt.Println("Moving vertically...")
				pos.row++
			} else {
				pos.row--
			}
			o.move(pos.row, pos.col)
		}

	case STATE_TURN_OFF_EQUIP:
		fmt.Println("Rendering work for robot ID: ", o.getID(), " in STATE_TURN_OFF_EQUIP")
	case STATE_ABSORB_H2S:
		fmt.Println("Rendering work for robot ID: ", o.getID(), " in STATE_ABSORB_H2S")
		space := o.gridMap[o.oshaSensor.grid.col][o.oshaSensor.grid.row]
		H2S_Level := space.H2S_Level

		if H2S_Level > 0 && o.H2S_Storage < o.H2S_CAPACITY {
			fmt.Println("Absorbing H2S at level: ", H2S_Level)
			o.H2S_Storage += 1
			space.H2S_Level--
		}

	case STATE_HOLD_POSITION:
		fmt.Println("Rendering work for robot ID: ", o.getID(), " in STATE_HOLD_POSITION")
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
				label: OSHA_ROBOT,
			},
		},

		H2S_CAPACITY:  h2sCapacity,
		H2S_Storage:   0,
		MAX_RADIATION: 100,
		state:         STATE_HOLD_POSITION,
		gridMap:       gridMap,
		currentSpace:  &gridMap[row][col],
	}
}

// --- HELPER FUNCTIONS ---

func (o *oshaRobot) move(row, col int) {

	delete(o.currentSpace.oshaRobots, o.getID())

	o.oshaSensor.grid = gridPosition{row: row, col: col}

	o.currentSpace = &o.gridMap[row][col]

	o.currentSpace.oshaRobots[o.getID()] = o
}
