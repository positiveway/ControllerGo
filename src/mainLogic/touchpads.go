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
	LeftPad, RightPad, Stick *PadPosition
)

func initTouchpads() {
	LeftPad = MakePadPosition(Cfg.LeftPadRotation)
	RightPad = MakePadPosition(Cfg.RightPadRotation)
	Stick = MakePadPosition(Cfg.StickRotation)

}

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

func calcFromMaxPossible(x, y float64) float64 {
	maxPossibleX := math.Sqrt(gofuncs.Sqr(Cfg.PadRadius) - gofuncs.Sqr(y))
	ratioFromMaxPossible := x / maxPossibleX
	return ratioFromMaxPossible * Cfg.PadRadius
}

func (pos *Position) CalcFromMaxPossible() *Position {
	//important to use temp values then assign
	posFromMaxPossible := MakeEmptyPosition()
	posFromMaxPossible.x = calcFromMaxPossible(pos.x, pos.y)
	posFromMaxPossible.y = calcFromMaxPossible(pos.y, pos.x)
	if !gofuncs.AnyNotInit(pos.x, pos.y) {
		//gofuncs.Print("before: %.2f, %.2f after: %.2f, %.2f", pos.x, pos.y, posFromMaxPossible.x, posFromMaxPossible.y)
		if gofuncs.AnyNotInit(posFromMaxPossible.x, posFromMaxPossible.y) {
			gofuncs.Panic("Incorrect calculations")
		}
	}
	return posFromMaxPossible
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

var maxMagnitude = 1.0

func calcAndCheckMagnitude(x, y float64) float64 {
	magnitude := calcDistance(x, y)
	//gofuncs.Print("", magnitude)
	if magnitude > maxMagnitude {
		maxMagnitude = magnitude
		gofuncs.Print("New max magn: %.3f", maxMagnitude)
	}
	if magnitude > Cfg.PadRadius {
		gofuncs.Panic("Magnitude is greater than Pad radius: %v", magnitude)
	}
	return magnitude
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

func _calcShiftedRotationPos(x, y, rotationShift, magnitude float64) (float64, float64, int) {
	if gofuncs.AnyNotInit(x, y) {
		return x, y, 0
	}

	angle := calcRawAngle(x, y)
	shiftedAngle := angle + rotationShift

	shiftedX := gofuncs.Sqrt(gofuncs.Sqr(magnitude) / (gofuncs.Sqr(math.Tan(shiftedAngle*math.Pi/180)) + 1))
	shiftedY := gofuncs.Sqrt(gofuncs.Sqr(magnitude) - gofuncs.Sqr(shiftedX))

	angle = resolveRawCircleAngle(angle)
	shiftedAngle = resolveRawCircleAngle(shiftedAngle)

	if shiftedAngle > 180 {
		shiftedY *= -1
	}
	if shiftedAngle > 90 && shiftedAngle < 270 {
		shiftedX *= -1
	}

	shiftedAngleInt := gofuncs.FloatToIntRound[int](shiftedAngle)

	gofuncs.PrintDebug("Angle: %.2f->%.2f (%.2f), X: %.2f->%.2f, Y: %.2f->%.2f",
		angle, shiftedAngle, calcOneQuarterAngle(shiftedAngle), x, shiftedX, y, shiftedY)

	_resAngle := gofuncs.FloatToIntRound[int](calcRawResolvedAngle(shiftedX, shiftedY))
	if _resAngle != shiftedAngleInt {
		gofuncs.Panic("Incorrect calculations with angle: %v", _resAngle)
	}

	return shiftedX, shiftedY, shiftedAngleInt
}

func rp(x float64) float64 {
	return gofuncs.Round(x, 3)
}

func checkShiftCalculations(x, y, magnitude float64) {
	if isEmptyPos(x, y) {
		return
	}
	shiftedX, shiftedY, _ := _calcShiftedRotationPos(x, y, 0, magnitude)
	if rp(x) != rp(shiftedX) || rp(y) != rp(shiftedY) {
		gofuncs.Panic("Calculations error")
	} else {
		gofuncs.Print("passed")
	}
}

func calcShiftedRotationPos(x, y, rotationShift, magnitude float64) (*Position, int) {
	//checkShiftCalculations(x, y, magnitude)
	shiftedX, shiftedY, shiftedAngle := _calcShiftedRotationPos(x, y, rotationShift, magnitude)
	return MakePosition(shiftedX, shiftedY), shiftedAngle
}

func (pos *Position) CalcShiftedRotationPos(rotationShift float64) (*Position, int, float64) {
	magnitude := calcAndCheckMagnitude(pos.x, pos.y)
	shiftedPos, shiftedAngle := calcShiftedRotationPos(pos.x, pos.y, rotationShift, magnitude)
	return shiftedPos, shiftedAngle, magnitude
}

type PadPosition struct {
	curPos, prevPos, shiftedPos, fromMaxPossiblePos *Position
	//normalizedMagnitude
	magnitude                  float64
	shiftedAngle               int
	newValueHandled            bool
	lock                       sync.Mutex
	zone                       Zone
	zoneCanBeUsed, zoneChanged bool
	zoneRotation               float64
	awaitingCentralPosition    bool
}

func MakePadPosition(zoneRotation float64) *PadPosition {
	pad := PadPosition{}
	pad.curPos = MakeEmptyPosition()
	pad.prevPos = MakeEmptyPosition()
	pad.shiftedPos = MakeEmptyPosition()
	pad.fromMaxPossiblePos = MakeEmptyPosition()

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
	pad.prevPos.Update(pad.shiftedPos)
}

func (pad *PadPosition) ReCalculateValues() {
	pad.newValueHandled = false

	pad.shiftedPos, pad.shiftedAngle, pad.magnitude = pad.curPos.CalcShiftedRotationPos(pad.zoneRotation)
	pad.fromMaxPossiblePos.Update(pad.shiftedPos.CalcFromMaxPossible())
}

func (pad *PadPosition) setValue(fieldPointer *float64) {
	pad.Lock()
	defer pad.Unlock()

	*fieldPointer = event.value

	pad.ReCalculateValues()

	if !gofuncs.AnyNotInit(pad.curPos.x, pad.curPos.y) {
		moveMouse()
	}
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
	pad.shiftedPos.Reset()
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

	if input > Cfg.PadRadius {
		gofuncs.Panic("Axis input value is greater than %v. Current value: %v", Cfg.PadRadius, input)
	}

	output := Cfg.OutputMin + ((outputMax-Cfg.OutputMin)/(Cfg.PadRadius-Cfg.StickDeadzone))*(input-Cfg.StickDeadzone)
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
	if math.Abs(value) < Cfg.StickDeadzone {
		value = 0.0
	}
	return value
}

func calcOneQuarterAngle[T gofuncs.Number](resolvedAngle T) T {
	return T(math.Mod(float64(resolvedAngle), 90))
}
