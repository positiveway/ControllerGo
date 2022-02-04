package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

func isGreater(oldValue, newValue float64) bool {
	return math.Abs(newValue) > math.Abs(oldValue)
}

type Coords struct {
	_x, _y float64
	mu     sync.Mutex
}

func (coords *Coords) reset() {
	coords.setX(0)
	coords.setY(0)
}

func round(val float64) string {
	return fmt.Sprintf("%.3f", val)
}

func (coords *Coords) print() {
	fmt.Printf("X: %s, Y: %s\n", round(coords._x), round(coords._y))
}

func (coords *Coords) setX(x float64) {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	coords._x = x
	//coords.print()
}

func (coords *Coords) setY(y float64) {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	coords._y = y
	//coords.print()
}

func (coords *Coords) getX() float64 {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	return coords._x
}

func (coords *Coords) getY() float64 {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	return coords._y
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

func mouseForce(val float64) float64 {
	return convertRange(val, mouseMaxMove)
}

func (coords *Coords) CalcForces() (x, y int32) {
	x = int32(mouseForce(coords.getX()))
	y = int32(-mouseForce(coords.getY()))
	return
}

func moveMouse() {
	for {
		xForce, yForce := mouseMovement.CalcForces()
		if (xForce != 0) || (yForce != 0) {
			//fmt.Printf("%v %v\n", xForce, yForce)
			err := mouse.Move(xForce, yForce)
			check_err(err)
		}
		time.Sleep(mouseInterval)
	}
}

func (coords *Coords) calcScrollInterval() time.Duration {
	input := math.Abs(coords.getY())
	scroll := convertRange(input, scrollMaxRange)
	scrollInterval := scrollMaxValue - int64(math.Round(scroll))
	return time.Duration(scrollInterval) * time.Millisecond
}

func getDirection(val float64) int32 {
	if val == 0 {
		return 0
	} else if val > 0 {
		return 1
	} else if val < 1 {
		return -1
	}
	panic("direction error")
}

func (coords *Coords) getDirections() (int32, int32) {
	return getDirection(coords.getX()), getDirection(coords.getY())
}

func scroll() {
	for {
		scrollInterval := scrollMovement.calcScrollInterval()
		dirX, dirY := scrollMovement.getDirections()
		if dirX != 0 {
			err := mouse.Wheel(true, dirX)
			check_err(err)
		}
		if dirY != 0 {
			err := mouse.Wheel(false, dirY)
			check_err(err)
		}
		time.Sleep(scrollInterval)
	}
}
