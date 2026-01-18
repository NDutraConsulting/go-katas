package robotfarm

type stateBundle struct {
	parallelQueue *[]orangePos
	currentOrange *orangePos
	MAX_COL       int
	MAX_ROW       int
	freshOranges  *int
}

func (sb stateBundle) destructure() (*[]orangePos, *orangePos, int, int, *int) {
	return sb.parallelQueue, sb.currentOrange, sb.MAX_COL, sb.MAX_ROW, sb.freshOranges
}
