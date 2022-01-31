package main

import "math"

func isGreater(oldValue, newValue float64) bool {
	return math.Abs(newValue) > math.Abs(oldValue)
}

type Coords struct {
	x, y float64
}

func (coords *Coords) Reset() {
	coords.x = 0
	coords.y = 0
}

const mouseAccelStep = 0.5
const mouseUpdateEveryN = 5
const mouseDistThreshold = 5.0

type Movement struct {
	curCoords     Coords
	prevCoords    Coords
	accel         float64
	accelStep     float64
	updateEveryN  int
	updateCounter int
	distThreshold float64
}

func makeMovement(accelStep float64, updateEveryN int) Movement {
	return Movement{
		accelStep:     accelStep,
		updateEveryN:  updateEveryN,
		distThreshold: mouseDistThreshold,
	}
}

func (movement *Movement) SetX(x float64) {
	if isGreater(movement.curCoords.x, x) {
		movement.updateCounter++
		movement.curCoords.x = x
	}
}

func (movement *Movement) SetY(y float64) {
	if isGreater(movement.curCoords.y, y) {
		movement.updateCounter++
		movement.curCoords.y = y
	}
}

func (movement *Movement) Reset() {
	movement.curCoords.Reset()
	movement.prevCoords.Reset()
	movement.accel = 0
}

func (movement *Movement) Distance() float64 {
	xDif := movement.curCoords.x - movement.prevCoords.x
	yDif := movement.curCoords.y - movement.prevCoords.y
	return math.Sqrt(math.Pow(xDif, 2) + math.Pow(yDif, 2))
}

func (movement *Movement) UpdateAccel() {
	if movement.updateCounter >= movement.updateEveryN {
		dist := movement.Distance()
		if dist > movement.distThreshold {
			movement.accel += movement.accelStep
		} else {
			movement.accel = 0
		}

		movement.updateCounter = 0
		movement.prevCoords = movement.curCoords
		movement.curCoords.Reset()
	}
	//movement.updateCounter = 0
}

func (movement *Movement) ApplyAccel(value float64) float64 {
	accel := 1 + movement.accel
	//accel = math.Pow(accel, 2)
	return value * accel
}

func (movement *Movement) calcForce(value float64) int32 {
	force := movement.ApplyAccel(value)
	return int32(force)
}

func (movement *Movement) CalcForces() (xForce, yForce int32) {
	xForce = movement.calcForce(movement.curCoords.x)
	yForce = movement.calcForce(movement.curCoords.y)
	return
}
