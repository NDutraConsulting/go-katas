package concurrentsensors

import (
	"fmt"
)

type gridObject struct {
	grid  gridPosition
	ID    int64
	label string
	// geoPos *geoPosition //Optional
}

const (
	TYPE_OIL_WELL        = "OIL_WELL"
	TYPE_LITHIUM_MINE    = "LITHIUM_MINE"
	TYPE_RADIOACTIVE_ORE = "TYPE_RADIOACTIVE_ORE"
)

func (o gridObject) processEnvironment(space *gridSpace) {
	fmt.Println("Processing environment for object ID: ", o.getID())
	// Redux style update
	switch o.label {
	case TYPE_OIL_WELL:
		// Handle oil well update
		fmt.Println("Update ID: ", o.getID(), " ->Expected Type: ", TYPE_OIL_WELL)

		// Pump HS2
		space.updateH2S()

	case TYPE_LITHIUM_MINE:
		// Handle lithium mine update
		fmt.Println("Update ID: ", o.getID(), "Expected Type: ", TYPE_LITHIUM_MINE)
	case TYPE_RADIOACTIVE_ORE:
		// Handle radioactive ore update
		fmt.Println("Update ID: ", o.getID(), "Expected Type: ", TYPE_RADIOACTIVE_ORE)
		space.updateRadioactive(1)
	}
}

func (g gridObject) col() int {
	return g.grid.col
}
func (g gridObject) row() int {
	return g.grid.row
}
func (g gridObject) getID() string {
	return fmt.Sprintf("%s-%d", g.label, g.ID)
}

type gridObjectTable interface {
	processEnvironment(space *gridSpace)
	col() int
	row() int
	getID() string
}
