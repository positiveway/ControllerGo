package mainLogic

import (
	"ControllerGo/src/osSpec"
	"math"
	"sync"
	"time"
)

var scrollMovement = makeCoords()
var mousePad = makeSoloPadPosition()

func calcMove(value, prevValue float64) int {
	if isNotInit(prevValue) {
		return 0
	}

	diff := value - prevValue
	pixels := floatToInt(diff * mouseSpeed)

	return pixels
}

func RunMouseThread() {
	for {
		mousePad.mu.Lock()

		moveX := calcMove(mousePad.x, mousePad.prevX)
		moveY := calcMove(mousePad.y, mousePad.prevY)
		mousePad.update()

		mousePad.mu.Unlock()

		if moveX != 0 || moveY != 0 {
			osSpec.MoveMouse(moveX, moveY)
		}

		time.Sleep(mouseInterval)
	}
}

type PadPosition struct {
	x, y         float64
	prevX, prevY float64
	mu           sync.Mutex
}

func makeSoloPadPosition() *PadPosition {
	pad := PadPosition{}
	pad.reset()
	return &pad
}

func (pad *PadPosition) update() {
	pad.prevX = pad.x
	pad.prevY = pad.y
}

func (pad *PadPosition) setX() {
	pad.mu.Lock()
	defer pad.mu.Unlock()
	pad.x = event.value
}

func (pad *PadPosition) setY() {
	pad.mu.Lock()
	defer pad.mu.Unlock()
	pad.y = event.value
}

func (pad *PadPosition) reset() {
	pad.mu.Lock()
	defer pad.mu.Unlock()

	pad.x = NaN()
	pad.y = NaN()
	pad.prevX = pad.x
	pad.prevY = pad.y
}

func calcScrollInterval(input float64) time.Duration {
	return calcRefreshInterval(input, scrollSlowestInterval, scrollFastestInterval)
}

func getDirection(val float64, horizontal bool) int {
	if isNotInit(val) {
		return 0
	}
	if horizontal && math.Abs(val) < horizontalScrollThreshold {
		return 0
	}
	switch {
	case val == 0:
		return 0
	case val > 0:
		return -1
	case val < 0:
		return 1
	}
	panic("direction error")
}

func getDirections(x, y float64) (int, int) {
	hDir := getDirection(x, true)
	vDir := getDirection(y, false)
	//hDir *= -1

	if hDir != 0 {
		vDir = 0
	}
	return hDir, vDir
}

func RunScrollThread() {
	var hDir, vDir int
	for {
		scrollMovement.mu.Lock()

		scrollMovement.updateValues()
		hDir, vDir = getDirections(scrollMovement.x, scrollMovement.y)

		scrollInterval := numberToMillis(scrollFastestInterval)
		if scrollMovement.magnitude != 0 {
			scrollInterval = calcScrollInterval(scrollMovement.magnitude)
		}

		scrollMovement.mu.Unlock()

		if hDir != 0 {
			osSpec.ScrollHorizontal(hDir)
		}
		if vDir != 0 {
			osSpec.ScrollVertical(vDir)
		}

		time.Sleep(scrollInterval)
	}
}
