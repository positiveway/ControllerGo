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

const mouseMinMove float64 = 0
const mouseMaxMove float64 = 20

func convertRange(input, outputStart, outputEnd float64) (output float64) {
	sign := math.Signbit(input)
	input = math.Abs(input)

	inputStart := 0.0
	inputEnd := 1.0

	output = outputStart + ((outputEnd-outputStart)/(inputEnd-inputStart))*(input-inputStart)
	if sign {
		output *= -1
	}
	return
}

func mouseForce(val float64) float64 {
	return convertRange(val, mouseMinMove, mouseMaxMove)
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
		time.Sleep(INTERVAL)
	}
}
