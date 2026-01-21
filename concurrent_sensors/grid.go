package concurrentsensors

import "fmt"

type gridPosition struct {
	row int
	col int
}

type gridSpace struct {
	gridPosition
	gridObjects       []gridObjectTable // Optional
	oshaRobots        map[string]*oshaRobot
	oshaSensors       []oshaSensor // Optional
	H2S_Level         int
	radioactive_level int
	H2SPocketVolume   int
	// geoOrigin   *geoPosition      // Optional lat/long origin for the top-left corner of the grid space
	// altitude *int // Optional

}

func (g *gridSpace) updateH2S() {
	fmt.Println("Update H2S - Pocket Volume: ", g.H2SPocketVolume)
	if g.H2SPocketVolume > 0 {
		fmt.Println("Update H2S: ", g.H2S_Level)
		g.H2S_Level += 2
		g.H2SPocketVolume -= 2
	}
}

func (g *gridSpace) updateRadioactive(r int) {

	g.radioactive_level += r
}

// type geoPosition struct {
// 	lat  float64
// 	long float64
// }
