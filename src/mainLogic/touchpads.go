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

type Position struct {
	x, y float64
}

func MakePosition() *Position {
	return &Position{}
}

func (pos *Position) Reset() {
	pos.x = gofuncs.NaN()
	pos.y = gofuncs.NaN()
}

func (pos *Position) Update(newPos *Position) {
	pos.x, pos.y = newPos.x, newPos.y
}

func (pos *Position) GetCopy() *Position {
	return &Position{x: pos.x, y: pos.y}
}

func calcFromMaxPossible(x, y float64) float64 {
	maxPossibleX := math.Sqrt(gofuncs.Sqr(PadRadius) - gofuncs.Sqr(y))
	ratioFromMaxPossible := x / maxPossibleX
	return ratioFromMaxPossible * PadRadius
}

func (pos *Position) CalcFromMaxPossible() *Position {
	//important to use temp values then assign
	posFromMaxPossible := MakePosition()
	posFromMaxPossible.x = calcFromMaxPossible(pos.x, pos.y)
	posFromMaxPossible.y = calcFromMaxPossible(pos.y, pos.x)
	if !gofuncs.AnyNotInit(pos.x, pos.y) {
		gofuncs.Print("before: %.2f, %.2f after: %.2f, %.2f", pos.x, pos.y, posFromMaxPossible.x, posFromMaxPossible.y)
		if gofuncs.AnyNotInit(posFromMaxPossible.x, posFromMaxPossible.y) {
			gofuncs.Panic("Incorrect calculations")
		}
	}
	return posFromMaxPossible
}

func calcDistance(x, y float64) float64 {
	if gofuncs.AnyNotInit(x, y) {
		return 0
	}
	return math.Hypot(x, y)
}

var maxMagnitude = 1.0

func (pos *Position) CalcAndCheckMagnitude() float64 {
	magnitude := calcDistance(pos.x, pos.y)
	if magnitude > maxMagnitude {
		maxMagnitude = magnitude
		gofuncs.Print("New max magn: %.3f", maxMagnitude)
	}
	if magnitude > PadRadius {
		gofuncs.Panic("Magnitude is greater than Pad radius: %v", magnitude)
	}
	return magnitude
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

func (pos *Position) CalcResolvedAngle() int {
	return calcResolvedAngle(pos.x, pos.y)
}

type PadPosition struct {
	curPos, prevPos, fromMaxPossiblePos *Position
	magnitude                           float64
	angle                               int
	newValueHandled                     bool
	lock                                sync.Mutex
	zone                                Zone
	zoneCanBeUsed, zoneChanged          bool
	zoneRotation                        int
	awaitingCentralPosition             bool
}

func MakePadPosition(zoneRotation int) *PadPosition {
	pad := PadPosition{}
	pad.curPos = MakePosition()
	pad.prevPos = MakePosition()
	pad.fromMaxPossiblePos = MakePosition()

	pad.zoneRotation = zoneRotation
	pad.Reset()

	return &pad
}

func (pad *PadPosition) Lock() {
	pad.lock.Lock()
}

func (pad *PadPosition) Unlock() {
	pad.lock.Unlock()
}

func (pad *PadPosition) UpdatePrevValues() {
	pad.prevPos.Update(pad.curPos)
}

func (pad *PadPosition) ReCalculateValues() {
	pad.newValueHandled = false

	pad.magnitude = pad.curPos.CalcAndCheckMagnitude()
	pad.angle = pad.curPos.CalcResolvedAngle()
	pad.fromMaxPossiblePos.Update(pad.curPos.CalcFromMaxPossible())
}

func (pad *PadPosition) setValue(fieldPointer *float64) {
	pad.Lock()
	defer pad.Unlock()

	*fieldPointer = event.value

	pad.ReCalculateValues()
}

func (pad *PadPosition) SetX() {
	pad.setValue(&(pad.curPos.x))
}

func (pad *PadPosition) SetY() {
	pad.setValue(&(pad.curPos.y))
}

func (pad *PadPosition) Reset() {
	pad.Lock()
	defer pad.Unlock()

	pad.curPos.Reset()
	//don't reset prev value to calc proper delta from prev to zero
	pad.prevPos.Reset()
	pad.fromMaxPossiblePos.Reset()

	pad.ReCalculateValues()
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
