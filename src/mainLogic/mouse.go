package mainLogic

import (
	"ControllerGo/src/platformSpecific"
	"math"
	"time"
)

var scrollMovement = Coords{}

const maxRightY float64 = 0.97
const maxRightX float64 = 0.97

type TouchPadPosition struct {
	x, y         float64
	prevX, prevY float64
}

func (pad *TouchPadPosition) distance() float64 {
	return distance(pad.prevX, pad.prevY, pad.x, pad.y)
}

func (pad *TouchPadPosition) difX() float64 {
	return pad.x - pad.prevX
}

func (pad *TouchPadPosition) difY() float64 {
	return pad.y - pad.prevY
}

func (pad *TouchPadPosition) update() {
	pad.prevX = pad.x
	pad.prevY = pad.y
}

var mousePad = TouchPadPosition{}

func distance(x1, y1, x2, y2 float64) float64 {
	return math.Hypot(x2-x1, y2-y1)
}

const changeThreshold float64 = 0.01

func calcPixels(curValue, prevValue float64) int32 {

	if prevValue == 0 || curValue == 0 {
		mousePad.update()
		return 0
	}

	diff := curValue - prevValue
	if math.Abs(diff) <= changeThreshold {
		return 0
	}

	pixels := floatToInt32(diff * mouseMaxMove)
	if pixels != 0 {
		mousePad.update()
	}
	return pixels
}

func moveMouse() {
	if pixels := calcPixels(mousePad.x, mousePad.prevX); pixels != 0 {
		platformSpecific.MoveMouse(pixels, 0)
	}
	if pixels := calcPixels(mousePad.y, mousePad.prevY); pixels != 0 {
		platformSpecific.MoveMouse(0, pixels)
	}
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
