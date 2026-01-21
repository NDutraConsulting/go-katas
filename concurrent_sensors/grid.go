package concurrentsensors

type gridPosition struct {
	row int
	col int
}

type gridSpace struct {
	gridPosition
	gridObjects       []gridObjectTable // Optional
	oshaRobots        []oshaRobot       // Optional
	oshaSensors       []oshaSensor      // Optional
	H2S_Level         int
	radioactive_level int
	H2SPocketVolume   int
	// geoOrigin   *geoPosition      // Optional lat/long origin for the top-left corner of the grid space
	// altitude *int // Optional

}

func (g *gridSpace) updateH2S() {

	if g.H2SPocketVolume > 0 {
		g.H2S_Level += 1
		g.H2SPocketVolume -= 1
	}
}

func (g *gridSpace) updateRadioactive(r int) {

	g.radioactive_level += r
}

// type geoPosition struct {
// 	lat  float64
// 	long float64
// }
