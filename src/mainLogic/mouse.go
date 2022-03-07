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

type MouseConfigs struct {
	interval,
	speedMult,
	speedPower,
	accelMult,
	accelPower float64
	intervalTime time.Duration
}

func makeMouseConfigs(
	interval,
	speedMult,
	speedPower,
	accelMult,
	accelPower float64) MouseConfigs {
	return MouseConfigs{
		intervalTime: numberToMillis(interval),
		interval:     interval,
		speedMult:    speedMult,
		speedPower:   speedPower,
		accelMult:    accelMult,
		accelPower:   accelPower,
	}
}

var mouseConfigs = makeMouseConfigs(12, 5000, 1, 0, 1)

const changeThreshold float64 = 0.001

type TouchPadPosition struct {
	x, y           float64
	prevX, prevY   float64
	startTimePoint time.Time
	prevSpeed      float64
}

func updateCoord(value float64, prevValue *float64, pixels int32) {
	if pixels != 0 || *prevValue == CoordNotInitialized || value == CoordNotInitialized {
		*prevValue = value
	}
}

func (pad *TouchPadPosition) updateX(pixels int32) {
	updateCoord(pad.x, &pad.prevX, pixels)
}

func (pad *TouchPadPosition) updateY(pixels int32) {
	updateCoord(pad.y, &pad.prevY, pixels)
}

func (pad *TouchPadPosition) updateValues(speed float64, startTime time.Time) {
	pad.prevSpeed = speed
	pad.startTimePoint = startTime
}

type MouseMetrics struct {
	speed, accel, coef float64
}

func pow(value, power float64) float64 {
	sign, value := getSignAndAbs(value)
	value = math.Pow(value, power)
	return applySign(sign, value)
}

func (m *MouseMetrics) calcCoef() {
	speed := mouseConfigs.speedMult * pow(m.speed, mouseConfigs.speedPower) * mouseConfigs.interval
	accel := mouseConfigs.accelMult * pow(m.accel, mouseConfigs.accelPower) * pow(mouseConfigs.interval, 2)
	m.coef = speed + accel
	print("speed: %0.4f; accel: %0.4f", speed, accel)
}

func (m *MouseMetrics) calcMove(diff float64) int32 {
	pixels := diff * m.coef
	print("diff: %0.4f, pixels: %0.4f", diff, pixels)
	return floatToInt32(pixels)
}

func (pad *TouchPadPosition) diffX() float64 {
	return diffIgnoreNotInit(pad.x, pad.prevX)
}

func (pad *TouchPadPosition) diffY() float64 {
	return diffIgnoreNotInit(pad.y, pad.prevY)
}

func (pad *TouchPadPosition) updateMetrics() {
	timeNow := time.Now()
	timeInt := mouseConfigs.interval

	if pad.startTimePoint == TimeNotInitialized {
		pad.initAccelTime(timeNow)
		return
	}

	timeDiff := timeNow.Sub(pad.startTimePoint)
	if timeDiff >= mouseConfigs.intervalTime {
		m := MouseMetrics{}

		diffX := pad.diffX()
		diffY := pad.diffY()

		dist := calcDistance(diffX, diffY)

		m.speed = dist / timeInt
		m.accel = (m.speed - pad.prevSpeed) / timeInt
		m.calcCoef()

		moveX := m.calcMove(diffX)
		moveY := m.calcMove(diffY)

		pad.updateX(moveX)
		pad.updateY(moveY)

		platformSpecific.MoveMouse(moveX, moveY)

		pad.updateValues(m.speed, timeNow)
	}
	return
}

func makeTouchPosition() TouchPadPosition {
	pad := TouchPadPosition{}
	pad.reset()
	return pad
}

func (pad *TouchPadPosition) setX() {
	pad.x = event.value
	pad.updateMetrics()
}

func (pad *TouchPadPosition) setY() {
	pad.y = event.value
	pad.updateMetrics()
}

func (pad *TouchPadPosition) reset() {
	pad.x = CoordNotInitialized
	pad.y = CoordNotInitialized
	pad.prevX = pad.x
	pad.prevY = pad.y

	pad.updateValues(0, TimeNotInitialized)
}

func diffIgnoreNotInit(curValue, prevValue float64) float64 {
	if curValue == CoordNotInitialized || prevValue == CoordNotInitialized {
		return 0
	}
	return curValue - prevValue
}

func (pad *TouchPadPosition) initAccelTime(startTime time.Time) {
	pad.startTimePoint = startTime
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
