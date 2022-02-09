package main

import "math"

const calcEveryN int = 30
const accelTimeInterval = float64(mouseIntervalInt * calcEveryN)

var counter int = 0

func calcAccel(x, y float64) float64 {
	counter++
	if counter < calcEveryN {
		return 0
	} else {
		counter = 0
	}
	prevX, prevY := prevMouse.getValues()
	distance := math.Sqrt(math.Pow(x-prevX, 2) + math.Pow(y-prevY, 2))
	distance *= math.Pow(10, 5)
	accel := 2 * distance / math.Pow(accelTimeInterval, 2)
	accel = round(accel, 1)
	prevMouse.setValues(x, y)
	//if accel != 0 {
	//	fmt.Printf("%.8f\n", accel)
	//}
	return accel
}
