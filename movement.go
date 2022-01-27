package main

type Coords struct {
	x, y float64
}

func (coords *Coords) Reset() {
	coords.x = 0
	coords.y = 0
}

type Movement struct {
	curCoords  Coords
	prevCoords Coords
	accel      float64
}

func (movement *Movement) SetX(x float64) {
	if x > movement.curCoords.x {
		movement.curCoords.x = x
	}
}

func (movement *Movement) SetY(y float64) {
	if y > movement.curCoords.y {
		movement.curCoords.y = y
	}
}

func (movement *Movement) Reset() {
	movement.curCoords.Reset()
	movement.prevCoords.Reset()
	movement.accel = 0
}

func (movement *Movement) Update() {
	movement.UpdateAccel()
	movement.prevCoords = movement.curCoords
	movement.curCoords.Reset()
}

const AccelStep = 0.01

func (movement *Movement) UpdateAccel() {
	movement.accel += AccelStep
}

func (movement *Movement) applyAccel(value float64) float64 {
	return value * (1 + movement.accel)
}
