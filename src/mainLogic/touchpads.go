package mainLogic

import (
	"github.com/positiveway/gofuncs"
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
	_x, _y                     float64
	actualX, actualY           float64
	prevX, prevY               float64
	magnitude                  float64
	angle                      int
	newValueHandled            bool
	lock                       sync.Mutex
	zone                       Zone
	zoneCanBeUsed, zoneChanged bool
	zoneRotation               int
	awaitingCentralPosition    bool
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
	pad.prevX = pad.actualX
	pad.prevY = pad.actualY
}

var maxMagnitude = 1.0

func checkMagnitude(x, y float64) {
	magnitude := calcDistance(x, y)
	if magnitude > maxMagnitude {
		maxMagnitude = magnitude
		//print("New max magn: %.3f", maxMagnitude)
	}
	if magnitude > PadRadius {
		gofuncs.Panic("Magnitude is greater than Pad radius")
	}
}

func (pad *PadPosition) ReCalculateValues() {
	pad.newValueHandled = false

	checkMagnitude(pad._x, pad._y)
	pad.actualX, pad.actualY = pad._x, pad._y

	pad.angle = calcResolvedAngle(pad.actualX, pad.actualY)
}

func calcFromMaxPossible(x, y float64) float64 {
	maxPossibleX := math.Sqrt(gofuncs.Sqr(PadRadius) - gofuncs.Sqr(y))
	ratioFromMaxPossible := x / maxPossibleX
	return ratioFromMaxPossible * PadRadius
}

func (pad *PadPosition) CalcCoordsFromMaxPossible() {
	//important to use temp values then assign
	xFromMaxPossible := calcFromMaxPossible(pad.actualX, pad.actualY)
	yFromMaxPossible := calcFromMaxPossible(pad.actualY, pad.actualX)
	if !gofuncs.AnyNotInit(xFromMaxPossible, yFromMaxPossible) {
		print("before: %.2f, %.2f after: %.2f, %.2f", pad.actualX, pad.actualY, xFromMaxPossible, yFromMaxPossible)
	}
	pad.actualX, pad.actualY = xFromMaxPossible, yFromMaxPossible
}

func (pad *PadPosition) setValue(fieldPointer *float64) {
	pad.Lock()
	defer pad.Unlock()

	*fieldPointer = event.value

	pad.ReCalculateValues()
}

func (pad *PadPosition) SetX() {
	pad.setValue(&pad._x)
}

func (pad *PadPosition) SetY() {
	pad.setValue(&pad._y)
}

func (pad *PadPosition) Reset() {
	pad.Lock()
	defer pad.Unlock()

	pad._x = gofuncs.NaN()
	pad._y = gofuncs.NaN()
	pad.prevX = gofuncs.NaN()
	pad.prevY = gofuncs.NaN()

	pad.ReCalculateValues()
}

func resolveAngle[T gofuncs.Number](angle T) int {
	resolvedAngle := math.Mod(float64(angle)+360, 360)
	return gofuncs.FloatToIntRound[int](resolvedAngle)
}

const RadiansMultiplier float64 = 180 / math.Pi

func calcResolvedAngle(x, y float64) int {
	if gofuncs.AnyNotInit(x, y) {
		return 0
	}
	angle := math.Atan2(y, x) * RadiansMultiplier
	return resolveAngle(angle)
}

func calcDistance(x, y float64) float64 {
	if gofuncs.AnyNotInit(x, y) {
		return 0
	}
	return math.Hypot(x, y)
}

//func normalizeIncorrectEdgeValues(x, y float64) (float64, float64, float64) {
//	magnitude := calcDistance(x, y)
//	if magnitude > PadRadius {
//		x /= magnitude
//		y /= magnitude
//		magnitude = PadRadius
//	}
//	return x, y, magnitude
//}

const OutputMin float64 = 0.0
const PadRadius = math.Sqrt2

func convertRange(input, outputMax float64) float64 {
	gofuncs.PanicAnyNotInit(input)

	if input == 0.0 {
		return 0.0
	}

	isNegative, input := gofuncs.GetIsNegativeAndAbs(input)

	if input > PadRadius {
		gofuncs.Panic("Axis input value is greater than %v. Current value: %v", PadRadius, input)
	}

	output := OutputMin + ((outputMax-OutputMin)/(PadRadius-StickDeadzone))*(input-StickDeadzone)
	return gofuncs.ApplySign(isNegative, output)
}

func calcRefreshInterval(input, slowestInterval, fastestInterval float64) time.Duration {
	input = math.Abs(input)
	refreshInterval := convertRange(input, slowestInterval-fastestInterval)
	refreshInterval = slowestInterval - refreshInterval
	return time.Duration(gofuncs.FloatToIntRound[int64](refreshInterval)) * time.Millisecond
}

func applyDeadzone(value float64) float64 {
	if gofuncs.IsNotInit(value) {
		return value
	}
	if math.Abs(value) < StickDeadzone {
		value = 0.0
	}
	return value
}

//func calcOneQuarterAngle(resolvedAngle int) int {
//	return floatToInt(math.Mod(float64(resolvedAngle), 90))
//}
