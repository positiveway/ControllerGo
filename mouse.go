package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

const mouseMaxMove float64 = 4
const forcePower float64 = 1.5
const deadzone float64 = 0.06

//const mouseScaleFactor float64 = 3
//var mouseIntervalInt int = int(math.Round(mouseMaxMove*mouseScaleFactor))
const mouseIntervalInt int = 12
const mouseInterval time.Duration = time.Duration(mouseIntervalInt) * time.Millisecond

const scrollMinValue float64 = 35
const scrollMaxValue float64 = 200

const scrollMaxRange float64 = scrollMaxValue - scrollMinValue
const horizontalScrollThreshold float64 = 0.45

var mouseMovement = Coords{}
var prevMouse = Coords{}
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

func convertRange(input, outputEnd float64) float64 {
	sign := getSignMakeAbs(&input)

	if input <= deadzone {
		return 0.0
	}

	outputStart := 1.0
	inputStart := 0.0
	inputEnd := 1.0

	output := outputStart + ((outputEnd-outputStart)/(inputEnd-inputStart))*(input-inputStart)
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

func moveMouse() {
	for {
		xForce, yForce := calcForces()
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

func getDirections() (int32, int32) {
	x, y := scrollMovement.getValues()
	hDir, vDir := getDirection(x, true), getDirection(y, false)
	hDir *= -1

	if hDir != 0 {
		vDir = 0
	}
	return hDir, vDir
}

func scroll() {
	for {
		hDir, vDir := getDirections()

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
