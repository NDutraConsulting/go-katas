package concurrentsensors

import "fmt"

type gridPosition struct {
	row int
	col int
}

type gridSpace struct {
	gridPosition
	gridObjects      []gridObjectTable // Optional
	oshaRobots       map[string]*oshaRobot
	oshaSensors      []oshaSensor // Optional
	h2sLevel         int
	radioactiveLevel int
	h2sPocketVolume  int
}

func (g *gridSpace) updateH2S() {
	fmt.Println("Update H2S - Pocket Volume: ", g.h2sPocketVolume)
	if g.h2sPocketVolume > 0 {
		fmt.Println("Update H2S: ", g.h2sLevel)
		g.h2sLevel += 2
		g.h2sPocketVolume -= 2
	}
}

func (g *gridSpace) updateRadioactive(r int) {

	g.radioactiveLevel += r
}
