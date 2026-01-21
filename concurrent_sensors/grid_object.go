package concurrentsensors

import (
	"fmt"
)

type gridObject struct {
	grid  gridPosition
	ID    int64
	label string
}

const (
	typeOilWell        = "OIL_WELL"
	typeLithiumMine    = "LITHIUM_MINE"
	typeRadioactiveOre = "typeRadioactiveOre"
)

func (o gridObject) processEnvironment(space *gridSpace) {
	fmt.Println("Processing environment for object ID: ", o.getID())
	// Redux style update
	switch o.label {
	case typeOilWell:
		// Handle oil well update
		fmt.Println("Update ID: ", o.getID(), " ->Expected Type: ", typeOilWell)

		// Pump HS2
		space.updateH2S()

	case typeLithiumMine:
		// Handle lithium mine update
		fmt.Println("Update ID: ", o.getID(), "Expected Type: ", typeLithiumMine)
	case typeRadioactiveOre:
		// Handle radioactive ore update
		fmt.Println("Update ID: ", o.getID(), "Expected Type: ", typeRadioactiveOre)
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
