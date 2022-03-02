package mainLogic

import "math"

type TouchPadPosition struct {
	x, y         float64
	prevX, prevY float64
}

func distance(x1, y1, x2, y2 float64) float64 {
	return math.Hypot(x2-x1, y2-y1)
}
