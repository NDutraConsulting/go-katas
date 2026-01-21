package concurrentsensors

import (
	"fmt"
	"sync/atomic"
)

const (
	commandHazardDetected = "HAZARD_DETECTED"
	commandRemoveHazard   = "REMOVE_HAZARD"
	commandStop           = "STOP"
	commandFindHazard     = "FIND_HAZARD"
)

type h2sRecord struct {
	h2sLevel  int
	timestamp int64
	state     string
}

type oshaSensor struct {
	gridObject
	h2sThreshold int
	h2sRecords   []h2sRecord
	currentSpace *gridSpace
}

func (o oshaSensor) renderWork(space *gridSpace) {
	fmt.Println("Rendering work for sensor ID: ", o.getID())
}

func (o oshaSensor) update() {
	fmt.Println("Update work for sensor ID: ", o.getID())

	// --- Send Warnings ---
	// COMMANDS - H2S DETECTED

}

var oshaSensorCounter int64

const (
	oshaSENSOR              = "OSHA_SENSOR"
	oshaH2SLow              = 3
	oshaH2SHazard           = 5
	oshaSensorStateNormal   = "NORMAL"
	oshaSensorStateWarning  = "WARNING"
	oshaSensorStateCritical = "CRITICAL"
)

// Constructor Function
func NewOshaSensor(row, col, h2sThreshold int) oshaSensor {
	return oshaSensor{
		// Increment and get the new value atomically
		gridObject: gridObject{
			ID:    atomic.AddInt64(&oshaSensorCounter, 1),
			grid:  gridPosition{row: row, col: col},
			label: "OSHA Sensor",
		},
		h2sThreshold: h2sThreshold,
		h2sRecords:   []h2sRecord{},
	}
}
