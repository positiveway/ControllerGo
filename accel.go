package main

import "math"

const calcEveryN int = 3
const accelTimeInterval = float64(mouseIntervalInt * calcEveryN)

const accelMultiplier float64 = 0.1
const accelThreshold float64 = 5.0

var mouseAccel float64 = 0.0
var counter int = 0

func calcAccel(x, y float64) {
	counter++
	if counter > calcEveryN {
		counter = 0
		prevX, prevY := prevMouse.getValues()

		distance := math.Sqrt(math.Pow(x-prevX, 2) + math.Pow(y-prevY, 2))
		distance *= math.Pow(10, 5)

		mouseAccel = 2 * distance / math.Pow(accelTimeInterval, 2)
		mouseAccel *= accelMultiplier
		if mouseAccel < accelThreshold {
			mouseAccel = 0.0
		}
		//mouseAccel = round(mouseAccel, 1)

		prevMouse.setValues(x, y)
	}
}

func applyAccel(force float64) float64 {
	return force * (max(mouseAccel, 1.0))
}
