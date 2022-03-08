package mainLogic

import (
	"ControllerGo/src/platformSpecific"
	"math"
	"time"
)

func UnifiedDiffIgnoreNotInit(curValue, prevValue float64) float64 {
	if isNan(curValue, prevValue) {
		return 0
	}
	return curValue - prevValue
}

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

type UnifiedPadPosition struct {
	x, y           float64
	prevX, prevY   float64
	startTimePoint time.Time
	prevSpeed      float64
}

func updateCoord(value float64, prevValue *float64, pixels int32) {
	if pixels != 0 || isNan(value, *prevValue) {
		*prevValue = value
	}
}

func (pad *UnifiedPadPosition) updateX(pixels int32) {
	updateCoord(pad.x, &pad.prevX, pixels)
}

func (pad *UnifiedPadPosition) updateY(pixels int32) {
	updateCoord(pad.y, &pad.prevY, pixels)
}

func (pad *UnifiedPadPosition) updateValues(speed float64, startTime time.Time) {
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

func (pad *UnifiedPadPosition) diffX() float64 {
	return UnifiedDiffIgnoreNotInit(pad.x, pad.prevX)
}

func (pad *UnifiedPadPosition) diffY() float64 {
	return UnifiedDiffIgnoreNotInit(pad.y, pad.prevY)
}

func (pad *UnifiedPadPosition) updateMetrics() {
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

func makeTouchPosition() UnifiedPadPosition {
	pad := UnifiedPadPosition{}
	pad.reset()
	return pad
}

func (pad *UnifiedPadPosition) setX() {
	pad.x = event.value
	pad.updateMetrics()
}

func (pad *UnifiedPadPosition) setY() {
	pad.y = event.value
	pad.updateMetrics()
}

func (pad *UnifiedPadPosition) reset() {
	pad.x = math.NaN()
	pad.y = math.NaN()
	pad.prevX = pad.x
	pad.prevY = pad.y

	pad.updateValues(0, TimeNotInitialized)
}

func (pad *UnifiedPadPosition) initAccelTime(startTime time.Time) {
	pad.startTimePoint = startTime
}
