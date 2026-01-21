package concurrentsensors

import (
	"fmt"
	"sync/atomic"
)

type oshaRobot struct {
	oshaSensor
	H2S_CAPACITY    int
	H2S_Storage     int
	radiation_level int
	MAX_RADIATION   int
	state           string
	gridMap         [][]gridSpace
}

func (robot oshaRobot) processEnvironment(space *gridSpace, grid [][]gridSpace) {
	go func(rb oshaRobot, sp *gridSpace) {
		fmt.Println("THREAD START: Processing environment for robot ID: ", rb.getID())

		rb.update()

		for _, obj := range space.gridObjects {

			if goObj, ok := obj.(gridObject); ok {
				rb.checkCollisions(goObj)
			}

		}

		// capture robot and obj for the goroutine to avoid loop variable capture
		rb.renderWork(sp)
		rb.recordProgress()
		fmt.Println("THREAD STOP")
		fmt.Println()
	}(robot, space)

}

func (o oshaRobot) update() {

	fmt.Println("Update ID: ", o.getID(), " ->Expected Type: ", OSHA_ROBOT)
	// Read onboard sensor data

	// Read distributed sensor data from channel

	// COMMANDS - H2S DETECTED

	o.updateState()

}

func (o oshaRobot) updateState() {
	fmt.Println("Updating state for ID: ", o.getID())
}

func (o oshaRobot) checkCollisions(gridObj gridObject) {
	fmt.Println("Checking collisions for robot ID: ", o.getID(), " and object: ", gridObj)
}

func (o oshaRobot) renderWork(space *gridSpace) {
	fmt.Println("Rendering work for robot ID: ", o.getID())

	switch o.state {
	case STATE_MOVE:
		fmt.Println("Rendering work for robot ID: ", o.getID(), " in STATE_MOVE")
	case STATE_TURN_OFF_EQUIP:
		fmt.Println("Rendering work for robot ID: ", o.getID(), " in STATE_TURN_OFF_EQUIP")
	case STATE_ABSORB_H2S:
		fmt.Println("Rendering work for robot ID: ", o.getID(), " in STATE_ABSORB_H2S")
	case STATE_HOLD_POSITION:
		fmt.Println("Rendering work for robot ID: ", o.getID(), " in STATE_HOLD_POSITION")
	}
}

func (o oshaRobot) recordProgress() {
	fmt.Println("Recording progress for robot ID: ", o.getID())
}

// ---STATE STEPS---
const (
	STATE_MOVE           = "MOVE"
	STATE_TURN_OFF_EQUIP = "TURN_OFF_EQUIP"
	STATE_ABSORB_H2S     = "ABSORB_H2S"
	STATE_HOLD_POSITION  = "HOLD_POSITION"
)

// ---Initializers---
const (
	OSHA_ROBOT = "OSHA_ROBOT"
	H2S_SMALL  = 3
	H2S_MEDIUM = 5
	H2S_LARGE  = 8
)

var oshaRobotCounter int64

// Constructor Function
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
		state:         STATE_MOVE,
		gridMap:       gridMap,
	}
}
