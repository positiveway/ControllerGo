package mainLogic

import (
	"ControllerGo/src/platformSpecific"
	"math"
	"time"
)

var scrollMovement = Coords{}
var mousePad = makeSoloPadPosition()

const CoordNotInitialized = -10000

const changeThreshold float64 = 0.001

var mouseSpeed float64 = 300

func moveMouse(value float64, prevValue *float64, isX bool) {
	if *prevValue == CoordNotInitialized {
		*prevValue = value
		return
	}

	diff := value - *prevValue
	pixels := floatToInt32(diff * mouseSpeed)
	*prevValue = value
	if pixels != 0 {
		if isX {
			platformSpecific.MoveMouse(pixels, 0)
		} else {
			platformSpecific.MoveMouse(0, pixels)
		}
	}
}

type SoloPadPosition struct {
	prevX, prevY float64
}

func makeSoloPadPosition() SoloPadPosition {
	pad := SoloPadPosition{}
	pad.reset()
	return pad
}

func (pad *SoloPadPosition) setX() {
	moveMouse(event.value, &pad.prevX, true)
}

func (pad *SoloPadPosition) setY() {
	moveMouse(event.value, &pad.prevY, false)
}

func (pad *SoloPadPosition) reset() {
	pad.prevX = CoordNotInitialized
	pad.prevY = CoordNotInitialized
}

func calcScrollInterval(input float64) time.Duration {
	return calcRefreshInterval(input, scrollSlowestInterval, scrollFastestInterval)
}

func getDirection(val float64, horizontal bool) int32 {
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

func getDirections(x, y float64) (int32, int32) {
	hDir := getDirection(x, true)
	vDir := getDirection(y, false)
	//hDir *= -1

	if hDir != 0 {
		vDir = 0
	}
	return hDir, vDir
}

func RunScrollThread() {
	var hDir, vDir int32
	for {
		scrollMovement.updateValues()
		hDir, vDir = getDirections(scrollMovement.x, scrollMovement.y)

		scrollInterval := time.Duration(scrollFastestInterval) * time.Millisecond
		if scrollMovement.magnitude != 0 {
			scrollInterval = calcScrollInterval(scrollMovement.magnitude)
		}

		if hDir != 0 {
			platformSpecific.ScrollHorizontal(hDir)
		}
		if vDir != 0 {
			platformSpecific.ScrollVertical(vDir)
		}

		time.Sleep(scrollInterval)
	}
}
