package mainLogic

import (
	"math"
	"sync"
	"time"
)

type ModeType = int

const (
	TypingMode    ModeType = 0
	ScrollingMode ModeType = 1
	GamingMode    ModeType = 2
)

type PadsMode struct {
	currentMode, defaultMode ModeType
	lock                     sync.Mutex
}

func MakePadsMode(defaultMode ModeType) *PadsMode {
	return &PadsMode{currentMode: defaultMode, defaultMode: defaultMode}
}

func (mode *PadsMode) SwitchMode() {
	mode.lock.Lock()
	defer mode.lock.Unlock()

	if mode.currentMode == mode.defaultMode {
		mode.currentMode = TypingMode
	} else {
		mode.currentMode = mode.defaultMode
	}
}

func (mode *PadsMode) GetMode() ModeType {
	mode.lock.Lock()
	defer mode.lock.Unlock()

	return mode.currentMode
}

var (
	LeftPad  = MakePadPosition(PadsRotation)
	RightPad = MakePadPosition(PadsRotation)
	Stick    = MakePadPosition(StickRotation)
)

type PadPosition struct {
	x, y                       float64
	prevX, prevY               float64
	magnitude                  float64
	angle                      int
	newValueHandled            bool
	lock                       sync.Mutex
	zone                       Zone
	zoneCanBeUsed, zoneChanged bool
	zoneRotation               int
	awaitingCentralPostion     bool
}

func MakePadPosition(zoneRotation int) *PadPosition {
	pad := PadPosition{}

	pad.Reset()
	pad.zoneRotation = zoneRotation

	return &pad
}

func (pad *PadPosition) Lock() {
	pad.lock.Lock()
}

func (pad *PadPosition) Unlock() {
	pad.lock.Unlock()
}

func (pad *PadPosition) UpdatePrevValues() {
	pad.prevX = pad.x
	pad.prevY = pad.y
}

func (pad *PadPosition) ReCalculateValues() {
	pad.newValueHandled = false
	pad.x, pad.y, pad.magnitude = normalizeIncorrectEdgeValues(pad.x, pad.y)
	pad.angle = calcResolvedAngle(pad.x, pad.y)
}

func (pad *PadPosition) CalcActualCoords() {
	//important to use temp values then assign
	actualX := calcFromActualMax(pad.x, pad.y)
	actualY := calcFromActualMax(pad.y, pad.x)
	pad.x, pad.y = actualX, actualY
}

func (pad *PadPosition) setValue(fieldPointer *float64) {
	pad.Lock()
	defer pad.Unlock()

	*fieldPointer = event.value

	pad.ReCalculateValues()
}

func (pad *PadPosition) SetX() {
	pad.setValue(&pad.x)
}

func (pad *PadPosition) SetY() {
	pad.setValue(&pad.y)
}

func (pad *PadPosition) printCurState() {
	printPair(pad.x, pad.y, "(x, y): ")
}

func (pad *PadPosition) Reset() {
	pad.Lock()
	defer pad.Unlock()

	pad.x = NaN()
	pad.y = NaN()
	pad.prevX = pad.x
	pad.prevY = pad.y

	pad.ReCalculateValues()
}

func calcFromActualMax(x, y float64) float64 {
	maxPossibleX := math.Sqrt(1 - sqr(y))
	ratioFromMax := x / maxPossibleX
	return ratioFromMax
}

func resolveAngle[T Number](angle T) int {
	resolvedAngle := math.Mod(float64(angle)+360, 360)
	return floatToInt(resolvedAngle)
}

const RadiansMultiplier float64 = 180 / math.Pi

func calcResolvedAngle(x, y float64) int {
	if isNotInit(x, y) {
		return 0
	}
	angle := math.Atan2(y, x) * RadiansMultiplier
	return resolveAngle(angle)
}

func calcDistance(x, y float64) float64 {
	if isNotInit(x, y) {
		return 0
	}
	return math.Hypot(x, y)
}

func normalizeIncorrectEdgeValues(x, y float64) (float64, float64, float64) {
	magnitude := calcDistance(x, y)
	if magnitude > PadRadius {
		x /= magnitude
		y /= magnitude
		magnitude = PadRadius
	}
	return x, y, magnitude
}

const OutputMin float64 = 0.0
const PadRadius = 1.0

func convertRange(input, outputMax float64) float64 {
	panicIfNotInit(input)

	if input == 0.0 {
		return 0.0
	}

	sign, input := getSignAndAbs(input)

	if input > PadRadius {
		panicMsg("Axis input value is greater than %v. Current value: %v", PadRadius, input)
	}

	output := OutputMin + ((outputMax-OutputMin)/(PadRadius-StickDeadzone))*(input-StickDeadzone)
	return applySign(sign, output)
}

func calcRefreshInterval(input, slowestInterval, fastestInterval float64) time.Duration {
	input = math.Abs(input)
	refreshInterval := convertRange(input, slowestInterval-fastestInterval)
	refreshInterval = slowestInterval - refreshInterval
	return time.Duration(floatToInt64(refreshInterval)) * time.Millisecond
}

func applyDeadzone(value float64) float64 {
	if isNotInit(value) {
		return value
	}
	if math.Abs(value) < StickDeadzone {
		value = 0.0
	}
	return value
}

func printPair[T Number](_x, _y T, prefix string) {
	x, y := float64(_x), float64(_y)
	print("%s: %0.2f %0.2f", prefix, x, y)
}

//func calcOneQuarterAngle(resolvedAngle int) int {
//	return floatToInt(math.Mod(float64(resolvedAngle), 90))
//}
