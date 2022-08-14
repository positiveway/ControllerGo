package mainLogic

import (
	"github.com/positiveway/gofuncs"
	"math"
	"sync"
	"time"
)

type ModeType string

const (
	TypingMode ModeType = "Typing"
	MouseMode  ModeType = "Mouse"
	GamingMode ModeType = "Gaming"
)

type PadsSticksMode struct {
	currentMode, defaultMode ModeType
	lock                     sync.Mutex
}

func MakePadsSticksMode(defaultMode ModeType) *PadsSticksMode {
	return &PadsSticksMode{currentMode: defaultMode, defaultMode: defaultMode}
}

func (mode *PadsSticksMode) SwitchMode() {
	mode.lock.Lock()
	defer mode.lock.Unlock()

	if mode.currentMode == mode.defaultMode {
		mode.currentMode = TypingMode
	} else {
		mode.currentMode = mode.defaultMode
	}
}

func (mode *PadsSticksMode) GetMode() ModeType {
	mode.lock.Lock()
	defer mode.lock.Unlock()

	return mode.currentMode
}

var (
	LeftPad, RightPadStick, LeftStick *PadStickPosition
)

type Position struct {
	x, y float64
}

func MakeEmptyPosition() *Position {
	return MakePosition(0, 0)
}

func MakePosition(x, y float64) *Position {
	return &Position{x: x, y: y}
}

func (pos *Position) Reset() {
	pos.x = gofuncs.NaN()
	pos.y = gofuncs.NaN()
}

func (pos *Position) Update(newPos *Position) {
	pos.x, pos.y = newPos.x, newPos.y
}

func (pos *Position) GetCopy() *Position {
	return MakePosition(pos.x, pos.y)
}

func isEmptyPos(x, y float64) bool {
	return gofuncs.AnyNotInit(x, y) || (x == 0 && y == 0)
}

func calcDistance(x, y float64) float64 {
	if isEmptyPos(x, y) {
		return 0
	}
	return math.Hypot(x, y)
}

const RadiansMultiplier float64 = 180 / math.Pi

func resolveRawCircleAngle[T gofuncs.Number](angle T) float64 {
	return math.Mod(float64(angle)+360, 360)
}

func resolveCircleAngle[T gofuncs.Number](angle T) int {
	resolvedAngle := resolveRawCircleAngle(angle)
	return gofuncs.FloatToIntRound[int](resolvedAngle)
}

func calcRawAngle(x, y float64) float64 {
	if isEmptyPos(x, y) {
		return 0
	}

	angleInRads := math.Atan2(y, x)
	angleInDegrees := angleInRads * RadiansMultiplier
	return angleInDegrees
}

func calcResolvedAngle(x, y float64) int {
	return resolveCircleAngle(calcRawAngle(x, y))
}

func calcRawResolvedAngle(x, y float64) float64 {
	return resolveRawCircleAngle(calcRawAngle(x, y))
}

func (pos *Position) CalcTransformedPos(rotationShift float64) (*Position, int, float64) {
	magnitude := calcDistance(pos.x, pos.y)
	shiftedAngle := resolveCircleAngle(calcRawAngle(pos.x, pos.y) + rotationShift)
	return MakePosition(pos.x, pos.y), shiftedAngle, magnitude
}

type PadStickPosition struct {
	curPos, prevMousePos, transformedPos *Position
	magnitude                            float64
	shiftedAngle                         int
	radius                               float64
	newValueHandled                      bool
	lock                                 sync.Mutex
	zone                                 Zone
	zoneCanBeUsed, zoneChanged           bool
	zoneRotation                         float64
	awaitingCentralPosition              bool

	//fromMaxPossiblePos *Position
	//normalizedMagnitude
}

func MakePadPosition() *PadStickPosition {
	pad := PadStickPosition{}

	pad.curPos = MakeEmptyPosition()
	pad.prevMousePos = MakeEmptyPosition()
	pad.transformedPos = MakeEmptyPosition()

	//pad.fromMaxPossiblePos = MakeEmptyPosition()

	pad.zoneRotation = gofuncs.NaN()
	pad.Reset()

	return &pad
}

func (pad *PadStickPosition) Reset() {
	pad.Lock()
	defer pad.Unlock()

	pad.curPos.Reset()
	//don't reset prev value to calc proper delta from prev to zero
	pad.prevMousePos.Reset()
	pad.transformedPos.Reset()

	//pad.fromMaxPossiblePos.Reset()

	pad.ReCalculateValues()
}

func (pad *PadStickPosition) UpdatePrevMousePos() {
	pad.prevMousePos.Update(pad.transformedPos)
}

func calcRadius(magnitude float64) float64 {
	return gofuncs.Max(magnitude, Cfg.MinStandardPadRadius)
}

func (pad *PadStickPosition) ReCalculateValues() {
	pad.newValueHandled = false

	pad.transformedPos, pad.shiftedAngle, pad.magnitude = pad.curPos.CalcTransformedPos(pad.zoneRotation)
	pad.radius = calcRadius(pad.magnitude)
	//pad.fromMaxPossiblePos.Update(pad.shiftedPos.CalcFromMaxPossible(pad.radius))
}

func (pad *PadStickPosition) setValue(fieldPointer *float64) {
	pad.Lock()
	defer pad.Unlock()

	*fieldPointer = Event.value

	pad.ReCalculateValues()

	switch Cfg.ControllerInUse {
	case SteamController:
		moveMouseSC()
	}
}

func (pad *PadStickPosition) SetX() {
	pad.setValue(&(pad.curPos.x))
}

func (pad *PadStickPosition) SetY() {
	pad.setValue(&(pad.curPos.y))
}

func (pad *PadStickPosition) Lock() {
	pad.lock.Lock()
}

func (pad *PadStickPosition) Unlock() {
	pad.lock.Unlock()
}

func (pad *PadStickPosition) convertRange(input, outputMax float64) float64 {
	gofuncs.PanicAnyNotInit(input)

	if input == 0 {
		return 0
	}

	isNegative, input := gofuncs.GetIsNegativeAndAbs(input)

	if input > pad.radius {
		gofuncs.Panic("Axis input value is greater than %v. Current value: %v", pad.radius, input)
	}

	inputMin := Cfg.StickDeadzoneDS

	output := Cfg.OutputMin + ((outputMax-Cfg.OutputMin)/(pad.radius-inputMin))*(input-inputMin)
	return gofuncs.ApplySign(isNegative, output)
}

func (pad *PadStickPosition) calcRefreshInterval(input, slowestInterval, fastestInterval float64) time.Duration {
	input = math.Abs(input)

	//TODO: Check
	if input == 0 {
		return gofuncs.NumberToMillis(fastestInterval)
	}

	refreshInterval := pad.convertRange(input, slowestInterval-fastestInterval)
	refreshInterval = slowestInterval - refreshInterval

	return time.Duration(gofuncs.FloatToIntRound[int64](refreshInterval)) * time.Millisecond
}

func calcOneQuarterAngle[T gofuncs.Number](resolvedAngle T) T {
	return T(math.Mod(float64(resolvedAngle), 90))
}

func applyDeadzone(value float64) float64 {
	if gofuncs.IsNotInit(value) {
		return value
	}
	if math.Abs(value) <= Cfg.StickDeadzoneDS {
		value = 0
	}
	return value
}

//func calcFromMaxPossible(x, y, radius float64) float64 {
//	maxPossibleX := math.Sqrt(gofuncs.Sqr(radius) - gofuncs.Sqr(y))
//	ratioFromMaxPossible := x / maxPossibleX
//	return ratioFromMaxPossible * radius
//}
//
//func (pos *Position) CalcFromMaxPossible(radius float64) *Position {
//	//important to use temp values then assign
//	posFromMaxPossible := MakeEmptyPosition()
//	posFromMaxPossible.x = calcFromMaxPossible(pos.x, pos.y, radius)
//	posFromMaxPossible.y = calcFromMaxPossible(pos.y, pos.x, radius)
//	if !gofuncs.AnyNotInit(pos.x, pos.y) {
//		//gofuncs.Print("before: %.2f, %.2f after: %.2f, %.2f", pos.x, pos.y, posFromMaxPossible.x, posFromMaxPossible.y)
//		if gofuncs.AnyNotInit(posFromMaxPossible.x, posFromMaxPossible.y) {
//			gofuncs.Panic("Incorrect calculations")
//		}
//	}
//	return posFromMaxPossible
//}

//var maxMagnitude = 1.0
//
//func calcAndCheckMagnitude(x, y float64) float64 {
//	magnitude := calcDistance(x, y)
//	//gofuncs.Print("", magnitude)
//	if magnitude > maxMagnitude {
//		maxMagnitude = magnitude
//		gofuncs.Print("New max magn: %.3f", maxMagnitude)
//	}
//	return magnitude
//}

//func normalizeIncorrectEdgeValues(x, y float64) (float64, float64, float64) {
//	magnitude := calcDistance(x, y)
//	if magnitude > PadRadius {
//		x /= magnitude
//		y /= magnitude
//		magnitude = PadRadius
//	}
//	return x, y, magnitude
//}
