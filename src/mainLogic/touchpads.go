package mainLogic

import (
	"github.com/positiveway/gofuncs"
	"math"
	"sync"
)

var (
	LeftPad, RightPadStick, LeftStick *PadStickPositionT
)

type PositionT struct {
	x, y float64
}

func MakeEmptyPosition() *PositionT {
	position := &PositionT{}
	position.Reset()
	return position
}

func MakePosition(x, y float64) *PositionT {
	return &PositionT{x: x, y: y}
}

func (pos *PositionT) Reset() {
	switch Cfg.ControllerInUse {
	case SteamController:
		pos.x = gofuncs.NaN()
		pos.y = gofuncs.NaN()
	case DualShock:
		pos.x = 0
		pos.y = 0
	}
}

func (pos *PositionT) Update(newPos *PositionT) {
	pos.UpdateRaw(newPos.x, newPos.y)
}

func (pos *PositionT) UpdateRaw(x, y float64) {
	pos.x, pos.y = x, y
}

func (pos *PositionT) GetCopy() *PositionT {
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

func resolveCircleAngle[T gofuncs.Number](angle T) uint {
	resolvedAngle := resolveRawCircleAngle(angle)
	return gofuncs.FloatToIntRound[uint](resolvedAngle)
}

func calcRawAngle(x, y float64) float64 {
	if isEmptyPos(x, y) {
		return 0
	}

	angleInRads := math.Atan2(y, x)
	angleInDegrees := angleInRads * RadiansMultiplier
	return angleInDegrees
}

func calcResolvedAngle(x, y float64) uint {
	return resolveCircleAngle(calcRawAngle(x, y))
}

func calcRawResolvedAngle(x, y float64) float64 {
	return resolveRawCircleAngle(calcRawAngle(x, y))
}

func (pos *PositionT) CalcTransformedPos(rotationShift float64) (*PositionT, uint, float64) {
	x, y := pos.x, pos.y

	magnitude := calcDistance(x, y)
	shiftedAngle := resolveCircleAngle(calcRawAngle(x, y) + rotationShift)
	transformedPos := MakePosition(x, y)

	return transformedPos, shiftedAngle, magnitude
}

type PadStickPositionT struct {
	curPos, prevMousePos, transformedPos *PositionT
	magnitude                            float64
	shiftedAngle                         uint
	radius                               float64
	newValueHandled                      bool
	lock                                 sync.Mutex
	zone                                 ZoneT
	zoneCanBeUsed, zoneChanged           bool
	zoneRotation                         float64
	awaitingCentralPosition              bool

	convertRange ConvertRangeFuncT

	//fromMaxPossiblePos *PositionT
	//normalizedMagnitude
}

func MakePadPosition(zoneRotation float64, isOnLeftSide bool) *PadStickPositionT {
	pad := PadStickPositionT{}

	pad.curPos = MakeEmptyPosition()
	pad.prevMousePos = MakeEmptyPosition()
	pad.transformedPos = MakeEmptyPosition()
	//pad.fromMaxPossiblePos = MakeEmptyPosition()

	if isOnLeftSide {
		zoneRotation *= -1
	}
	pad.zoneRotation = zoneRotation

	pad.Reset()
	pad.Validate()

	pad.convertRange = pad.GetConvertRangeFunc()

	return &pad
}

func checkRotation(rotation float64) {
	if gofuncs.Abs(rotation) > 360 {
		gofuncs.Panic("Incorrect rotation: %v", rotation)
	}
}

func (pad *PadStickPositionT) Validate() {
	checkRotation(pad.zoneRotation)
}

func (pad *PadStickPositionT) Reset() {
	pad.Lock()
	defer pad.Unlock()

	pad.curPos.Reset()
	//don't reset prev value to calc proper delta from prev to zero
	pad.prevMousePos.Reset()
	pad.transformedPos.Reset()
	//pad.fromMaxPossiblePos.Reset()

	pad.ReCalculateValues()
}

func calcRadius(magnitude float64) float64 {
	return gofuncs.Max(magnitude, Cfg.PadsSticks.MinStandardRadius)
}

func calcFromMaxPossible(x, y, radius float64) float64 {
	maxPossibleX := math.Sqrt(gofuncs.Sqr(radius) - gofuncs.Sqr(y))
	if maxPossibleX == 0 {
		return 0
	}

	ratioFromMaxPossible := x / maxPossibleX

	if ratioFromMaxPossible > radius {
		if ratioFromMaxPossible > radius+Cfg.Math.FloatEqualityMargin {
			gofuncs.Panic("Incorrect calculations")
		}
		ratioFromMaxPossible = radius
	}
	return ratioFromMaxPossible
}

func (pos *PositionT) CalcFromMaxPossible(radius float64) *PositionT {
	//important to use temp values then assign
	posFromMaxPossible := MakeEmptyPosition()
	posFromMaxPossible.x = calcFromMaxPossible(pos.x, pos.y, radius)
	posFromMaxPossible.y = calcFromMaxPossible(pos.y, pos.x, radius)
	if !gofuncs.AnyNotInit(pos.x, pos.y) {
		//gofuncs.Print("before: %.2f, %.2f after: %.2f, %.2f", pos.x, pos.y, posFromMaxPossible.x, posFromMaxPossible.y)
		if gofuncs.AnyNotInit(posFromMaxPossible.x, posFromMaxPossible.y) {
			gofuncs.Panic("Incorrect calculations")
		}
	}
	return posFromMaxPossible
}

func (pad *PadStickPositionT) ReCalculateValues() {
	//never assign position (pointer field) directly
	var _transformedPos *PositionT

	pad.newValueHandled = false

	_transformedPos, pad.shiftedAngle, pad.magnitude = pad.curPos.CalcTransformedPos(pad.zoneRotation)
	pad.transformedPos.Update(_transformedPos)

	pad.radius = calcRadius(pad.magnitude)

	//pad.fromMaxPossiblePos.Update(pad.transformedPos.CalcFromMaxPossible(pad.radius))
}

func (pad *PadStickPositionT) setValue(fieldPointer *float64, value float64) {
	pad.Lock()
	defer pad.Unlock()

	*fieldPointer = value

	pad.ReCalculateValues()

	switch Cfg.ControllerInUse {
	case SteamController:
		MoveMouseSC()
	}
}

func (pad *PadStickPositionT) SetX(value float64) {
	pad.setValue(&(pad.curPos.x), value)
}

func (pad *PadStickPositionT) SetY(value float64) {
	pad.setValue(&(pad.curPos.y), value)
}

func (pad *PadStickPositionT) Lock() {
	pad.lock.Lock()
}

func (pad *PadStickPositionT) Unlock() {
	pad.lock.Unlock()
}

type ConvertRangeFuncT = func(input, outputMax float64) float64

func (pad *PadStickPositionT) GetConvertRangeFunc() ConvertRangeFuncT {
	inputMin := Cfg.PadsSticks.Stick.DeadzoneDS
	outputMin := Cfg.Math.OutputMin

	return func(input, outputMax float64) float64 {
		gofuncs.PanicAnyNotInit(input)

		if input == 0 {
			return 0
		}

		isNegative, input := gofuncs.GetIsNegativeAndAbs(input)

		if input > pad.radius {
			gofuncs.Panic("Axis input value is greater than %v. Current value: %v", pad.radius, input)
		}

		output := outputMin + (outputMax-outputMin)/(pad.radius-inputMin)*(input-inputMin)
		return gofuncs.ApplySign(isNegative, output)
	}
}

func (pad *PadStickPositionT) calcRefreshInterval(input, slowestInterval, fastestInterval float64) float64 {
	input = math.Abs(input)

	refreshInterval := pad.convertRange(input, slowestInterval-fastestInterval)
	refreshInterval = slowestInterval - refreshInterval

	return float64(gofuncs.FloatToIntRound[int64](refreshInterval))
}

func calcOneQuarterAngle[T gofuncs.Number](resolvedAngle T) T {
	return T(math.Mod(float64(resolvedAngle), 90))
}

func applyDeadzone(value float64) float64 {
	if gofuncs.IsNotInit(value) {
		return value
	}
	if math.Abs(value) <= Cfg.PadsSticks.Stick.DeadzoneDS {
		value = 0
	}
	return value
}

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
