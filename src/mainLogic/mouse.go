package mainLogic

import (
	"ControllerGo/src/platformSpecific"
	"math"
	"time"
)

var scrollMovement = Coords{}
var mousePad = makeTouchPosition()

const CoordNotInitialized = -10000

var TimeNotInitialized = time.Now().Add(time.Hour)

var AccelIntervalNum float64 = 16
var AccelInterval = numberToMillis(AccelIntervalNum)

type TouchPadPosition struct {
	prevX, prevY                     float64
	prevAccelPointX, prevAccelPointY float64
	accelTimeStart                   time.Time
	accel                            float64
}

func makeTouchPosition() TouchPadPosition {
	pad := TouchPadPosition{}
	pad.reset()
	return pad
}

func (pad *TouchPadPosition) setX() {
	pad.updateAccelX()
	if pixels := pad.calcPixels(&pad.prevX); pixels != 0 {
		//print("x: %v", pixels)
		platformSpecific.MoveMouse(pixels, 0)
	}
}

func (pad *TouchPadPosition) setY() {
	pad.updateAccelY()
	if pixels := pad.calcPixels(&pad.prevY); pixels != 0 {
		//print("y: %v", pixels)
		platformSpecific.MoveMouse(0, pixels)
	}
}

func (pad *TouchPadPosition) resetAccel() {
	pad.updateAccelValues(CoordNotInitialized, CoordNotInitialized, 0, TimeNotInitialized)
}

func (pad *TouchPadPosition) reset() {
	pad.prevX = CoordNotInitialized
	pad.prevY = CoordNotInitialized

	pad.resetAccel()
}

const changeThreshold float64 = 0.005
const pixelsThreshold = 2

func diffIgnoreNotInit(curValue, prevValue float64) float64 {
	if curValue == CoordNotInitialized || prevValue == CoordNotInitialized {
		return 0
	}
	return curValue - prevValue
}
func diffCheckInit(curValue, prevValue float64) float64 {
	if curValue == CoordNotInitialized || prevValue == CoordNotInitialized {
		panicMsg("Value for diff is not initialized")
	}
	return curValue - prevValue
}

func (pad *TouchPadPosition) updateAccelValues(x, y, accel float64, startTime time.Time) {
	pad.accel = accel
	pad.accelTimeStart = startTime
	pad.prevAccelPointX = x
	pad.prevAccelPointY = y
}

func (pad *TouchPadPosition) initAccelTime(startTime time.Time) {
	pad.accelTimeStart = startTime
}

func (pad *TouchPadPosition) updateAccel(x, y float64) {
	timeNow := time.Now()
	if pad.accelTimeStart == TimeNotInitialized {
		pad.initAccelTime(timeNow)
		return
	}
	timeDiff := timeNow.Sub(pad.accelTimeStart)
	if timeDiff >= AccelInterval {

		dist := calcDistance(diffIgnoreNotInit(x, pad.prevAccelPointX), diffIgnoreNotInit(y, pad.prevAccelPointY))
		//pad.accel = 2 * dist / math.Pow(float64(timeDiff)/timeDiv, 2.0)
		accel := dist / AccelIntervalNum
		accel *= 100
		print("accel: %0.2f", accel)

		pad.updateAccelValues(x, y, accel, timeNow)
	}
}

func (pad *TouchPadPosition) updateAccelX() {
	pad.updateAccel(event.value, pad.prevY)
}

func (pad *TouchPadPosition) updateAccelY() {
	pad.updateAccel(pad.prevX, event.value)
}

func (pad *TouchPadPosition) calcPixels(prevValue *float64) int32 {
	curValue := event.value

	if *prevValue == CoordNotInitialized {
		*prevValue = curValue
		return 0
	}

	diff := diffCheckInit(curValue, *prevValue)
	if math.Abs(diff) <= changeThreshold {
		return 0
	}
	*prevValue = curValue
	pixels := floatToInt32(diff * mouseMaxMove)
	//if math.Abs(float64(pixels)) < pixelsThreshold {
	//	return 0
	//}
	return pixels
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
