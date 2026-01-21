package concurrentsensors

import (
	"fmt"
	"sync/atomic"
)

type h2sRecord struct {
	H2S_Level int
	timestamp int64
	state     string
}

type oshaSensor struct {
	gridObject
	H2S_Threshold int
	H2S_Records   []h2sRecord
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
	OSHA_SENSOR    = "OSHA_SENSOR"
	H2S_LOW        = 3
	H2S_HAZARD     = 5
	STATE_NORMAL   = "NORMAL"
	STATE_WARNING  = "WARNING"
	STATE_CRITICAL = "CRITICAL"
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
		H2S_Threshold: h2sThreshold,
		H2S_Records:   []h2sRecord{},
	}
}
