package mainLogic

import (
	"ControllerGo/src/platformSpecific"
	"math"
	"time"
)

var scrollMovement = Coords{}
var mousePad = makeTouchPosition()

const NotInitialized = -10000

type TouchPadPosition struct {
	prevX, prevY float64
	touchHappen  bool
}

func makeTouchPosition() TouchPadPosition {
	pad := TouchPadPosition{}
	pad.reset()
	return pad
}

func (pad *TouchPadPosition) setX() {
	if pixels := pad.calcPixels(event.value, &pad.prevX); pixels != 0 {
		//print("x: %v", pixels)
		platformSpecific.MoveMouse(pixels, 0)
	}
}

func (pad *TouchPadPosition) setY() {
	if pixels := pad.calcPixels(event.value, &pad.prevY); pixels != 0 {
		//print("y: %v", pixels)
		platformSpecific.MoveMouse(0, pixels)
	}
}

func (pad *TouchPadPosition) reset() {
	pad.prevX = NotInitialized
	pad.prevY = NotInitialized
	pad.touchHappen = false
}

const changeThreshold float64 = 0.01
const pixelsThreshold = 2

func (pad *TouchPadPosition) calcPixels(curValue float64, prevValue *float64) int32 {
	switch curValue {
	case AxisPressed:
		pad.touchHappen = true
		return 0
	case AxisReleased:
		pad.reset()
		return 0
	}

	if !pad.touchHappen {
		return 0
	} else {
		if *prevValue == NotInitialized {
			*prevValue = curValue
			return 0
		}
	}

	diff := curValue - *prevValue
	if math.Abs(diff) <= changeThreshold {
		return 0
	}

	pixels := floatToInt32(diff * mouseMaxMove)
	//if math.Abs(float64(pixels)) < pixelsThreshold {
	//	return 0
	//}
	*prevValue = curValue
	return pixels
}

func distance(x1, y1, x2, y2 float64) float64 {
	return math.Hypot(x2-x1, y2-y1)
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
