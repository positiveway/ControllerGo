package main

import (
	"math"
	"sync"
	"time"
)

const mouseInterval = 25 * time.Millisecond
const mouseMaxMove float64 = 12

const scrollMinValue float64 = 35
const scrollMaxValue float64 = 200

const scrollMaxRange float64 = scrollMaxValue - scrollMinValue
const horizontalScrollThreshold float64 = 0.45

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

func (coords *Coords) getValues() (float64, float64) {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	return coords._x, coords._y
}

func convertRange(input, outputEnd float64) (output float64) {
	sign := math.Signbit(input)
	input = math.Abs(input)

	outputStart := 0.0
	inputStart := 0.0
	inputEnd := 1.0

	output = outputStart + ((outputEnd-outputStart)/(inputEnd-inputStart))*(input-inputStart)
	if sign {
		output *= -1
	}
	return
}

func mouseForce(val float64) int32 {
	return int32(convertRange(val, mouseMaxMove))
}

func (coords *Coords) CalcForces() (xForce, yForce int32) {
	x, y := coords.getValues()
	xForce = mouseForce(x)
	yForce = -mouseForce(y)
	return
}

func moveMouse() {
	for {
		xForce, yForce := mouseMovement.CalcForces()
		if (xForce != 0) || (yForce != 0) {
			//fmt.Printf("%v %v\n", xForce, yForce)
			mouse.Move(xForce, yForce)
		}
		time.Sleep(mouseInterval)
	}
}

func calcScrollInterval(value float64) time.Duration {
	input := math.Abs(value)
	scroll := convertRange(input, scrollMaxRange)
	scrollInterval := scrollMaxValue - math.Round(scroll)
	return time.Duration(scrollInterval) * time.Millisecond
}

func getDirection(val float64, horizontal bool) int32 {
	if horizontal && math.Abs(val) < horizontalScrollThreshold {
		return 0
	}
	switch {
	case val == 0:
		return 0
	case val > 0:
		return 1
	case val < 0:
		return -1
	}
	panic("direction error")
}

func (coords *Coords) getDirections() (hDir, vDir int32) {
	x, y := coords.getValues()
	hDir, vDir = getDirection(x, true), getDirection(y, false)
	hDir *= -1

	if hDir != 0 {
		vDir = 0
	}
	return
}

func scroll() {
	for {
		hDir, vDir := scrollMovement.getDirections()

		x, y := scrollMovement.getValues()
		scrollVal := y
		if hDir != 0 {
			scrollVal = x
		}
		scrollInterval := calcScrollInterval(scrollVal)

		if hDir != 0 {
			mouse.Wheel(true, hDir)
		}
		if vDir != 0 {
			mouse.Wheel(false, vDir)
		}
		time.Sleep(scrollInterval)
	}
}
