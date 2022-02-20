package mainLogic

import (
	"ControllerGo/src/osSpecific"
	"fmt"
	"math"
	"sync"
	"time"
)

var mouseMovement = Coords{}
var scrollMovement = Coords{}

type Coords struct {
	_x, _y float64
	mu     sync.Mutex
}

func (coords *Coords) reset() {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	coords._x = 0
	coords._y = 0
}

func (coords *Coords) setX(x float64) {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	coords._x = x
}

func (coords *Coords) setY(y float64) {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	coords._y = y
}

func (coords *Coords) setValues(x, y float64) {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	coords._x = x
	coords._y = y
}

func (coords *Coords) getValues() (float64, float64) {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	return coords._x, coords._y
}

func applyDeadzone(value *float64) {
	if math.Abs(*value) < Deadzone {
		*value = 0.0
	}
}

func convertRange(input, outputMax float64) float64 {
	sign := getSignMakeAbs(&input)

	if input == 0.0 {
		return 0.0
	}

	if input > 1.0 {
		panic(fmt.Sprintf("Axis input value is greater than 1.0. Current value: %v\n", input))
	}

	outputMin := 1.0

	output := outputMin + ((outputMax-outputMin)/(inputRange))*(input-Deadzone)
	applySign(sign, &output)
	return output
}

func applyPower(force *float64) {
	sign := getSignMakeAbs(force)
	*force = math.Pow(*force, forcePower)
	applySign(sign, force)
}

func mouseForce(val float64) int32 {
	force := convertRange(val, mouseMaxMove)
	//printForce(force, "before")
	applyPower(&force)
	//printForce(force, "after")
	return int32(force)
}

func printForce(force float64, prefix string) {
	if force != 0.0 {
		fmt.Printf("%s: %0.3f\n", prefix, force)
	}
}

func printPair[T Number](_x, _y T, prefix string) {
	x, y := float64(_x), float64(_y)
	fmt.Printf("%s: %0.2f %0.2f\n", prefix, x, y)
}

func calcForces() (int32, int32) {
	x, y := mouseMovement.getValues()
	xForce := mouseForce(x)
	yForce := -mouseForce(y)

	//if x != 0.0 || y != 0.0{
	//	printPair(x, y, "x, y")
	//	printPair(xForce, yForce, "force")
	//	fmt.Println()
	//}
	return xForce, yForce
}

func RunMouseMoveThread() {
	for {
		xForce, yForce := calcForces()
		if (xForce != 0) || (yForce != 0) {
			//fmt.Printf("%v %v\n", xForce, yForce)
			osSpecific.MoveMouse(xForce, yForce)
		}
		time.Sleep(mouseInterval)
	}
}

func calcScrollInterval(value float64) time.Duration {
	input := math.Abs(value)
	scroll := convertRange(input, scrollIntervalRange)
	scrollInterval := scrollSlowestInterval - math.Round(scroll)
	return time.Duration(scrollInterval) * time.Millisecond
}

func getDirection(val float64, horizontal bool) int32 {
	if horizontal && math.Abs(val) < horizontalScrollThreshold {
		return 0
	}
	switch {
	case val == 0.0:
		return 0
	case val > 0:
		return 1
	case val < 0:
		return -1
	}
	panic("direction error")
}

func getDirections(x, y float64) (int32, int32) {
	hDir, vDir := getDirection(x, true), getDirection(y, false)
	hDir *= -1

	if hDir != 0 {
		vDir = 0
	}
	return hDir, vDir
}

func RunScrollThread() {
	for {
		x, y := scrollMovement.getValues()
		hDir, vDir := getDirections(x, y)

		if hDir != 0 {
			osSpecific.ScrollHorizontal(hDir)
		}
		if vDir != 0 {
			osSpecific.ScrollVertical(vDir)
		}

		scrollInterval := DefaultWaitInterval
		if hDir != 0 || vDir != 0 {
			scrollVal := y
			if hDir != 0 {
				scrollVal = x
			}
			scrollInterval = calcScrollInterval(scrollVal)
		}

		time.Sleep(scrollInterval)
	}
}
