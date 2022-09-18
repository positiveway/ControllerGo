package mainLogic

import (
	"ControllerGo/osSpec"
	"github.com/positiveway/gofuncs"
	"math"
	"time"
)

type PositionT struct {
	CfgStruct
	x, y  float64
	Reset func()
}

func MakeEmptyPosition(cfg *ConfigsT) *PositionT {
	position := &PositionT{}
	position.Init(cfg)

	position.Reset = position.GetResetFunc()
	position.Reset()

	return position
}

func MakePosition(x, y float64) *PositionT {
	return &PositionT{x: x, y: y}
}

func (pos *PositionT) GetResetFunc() func() {
	switch pos.cfg.ControllerInUse {
	case SteamController:
		return func() {
			pos.x = gofuncs.NaN()
			pos.y = gofuncs.NaN()
		}
	case DualShock:
		return func() {
			pos.x = 0
			pos.y = 0
		}
	default:
		PanicUnsupportedController()
	}
	panic("")
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
	CfgButtonsLockStruct
	curPos, prevMousePos, transformedPos *PositionT
	magnitude                            float64
	shiftedAngle                         uint
	radius                               float64
	newValueHandled                      bool
	zone                                 ZoneT
	zoneCanBeUsed, zoneChanged           bool
	zoneRotation                         float64
	awaitingCentralPosition              bool

	firstTouchTime   time.Time
	leftClickBtn     BtnOrAxisT
	leftClickCmdInfo *CommandInfoT

	convertRange ConvertRangeFuncT
	setValue     SetValueFuncT
	detectZone   DetectZoneFuncT
	calcRadius,
	Reset,
	moveMouseSC func()

	//fromMaxPossiblePos *PositionT
	//normalizedMagnitude
}

func (pad *PadStickPositionT) Init(zoneRotation float64, isOnLeftSide bool, cfg *ConfigsT, buttons *ButtonsT) {
	pad.CfgButtonsLockStruct.Init(cfg, buttons)

	pad.leftClickBtn, pad.leftClickCmdInfo = buttons.CreateVirtualButton(CommandT{osSpec.LeftMouse})

	pad.curPos = MakeEmptyPosition(cfg)
	pad.prevMousePos = MakeEmptyPosition(cfg)
	pad.transformedPos = MakeEmptyPosition(cfg)
	//pad.fromMaxPossiblePos = MakeEmptyPosition()

	if isOnLeftSide {
		zoneRotation *= -1
	}
	pad.zoneRotation = zoneRotation

	zeroTime, err := time.Parse("2006-01-02", "2006-01-02")
	gofuncs.CheckErr(err)
	pad.firstTouchTime = zeroTime

	pad.Reset = pad.GetResetFunc()
	pad.setValue = pad.GetSetValueFunc()
	pad.convertRange = pad.GetConvertRangeFunc()
	pad.calcRadius = pad.GetCalcRadiusFunc()
	pad.detectZone = pad.GetDetectZoneFunc()

	pad.Reset()
	pad.Validate()
}

func checkRotation(rotation float64) {
	if gofuncs.Abs(rotation) > 360 {
		gofuncs.Panic("Incorrect rotation: %v", rotation)
	}
}

func (pad *PadStickPositionT) InitMoveSCFunc(highPrecisionMode *HighPrecisionModeT) {
	pad.moveMouseSC = pad.GetMoveMouseSCFunc(highPrecisionMode)
}

func (pad *PadStickPositionT) Validate() {
	checkRotation(pad.zoneRotation)
}

func (pad *PadStickPositionT) GetResetFunc() func() {
	curPos := pad.curPos
	prevMousePos := pad.prevMousePos
	transformedPos := pad.transformedPos

	return func() {
		curPos.Reset()
		prevMousePos.Reset()
		transformedPos.Reset()
		//pad.fromMaxPossiblePos.Reset()

		//fmt.Println("reset")

		if pad.moveMouseSC != nil {
			//fmt.Println("release")
			pad.buttons.releaseButton(pad.leftClickBtn)
		}

		// should be inside reset func
		//to set values correctly during initial object creation
		pad.ReCalculateValues()
	}
}

func (pad *PadStickPositionT) GetCalcRadiusFunc() func() {
	minStandardRadius := pad.cfg.PadsSticks.MinStandardRadius
	if minStandardRadius < 1.0 {
		gofuncs.Panic("Radius can't be less than 1.0, current value: %v", minStandardRadius)
	}

	return func() {
		pad.radius = gofuncs.Max(pad.magnitude, minStandardRadius)
	}
}

func (pad *PadStickPositionT) ReCalculateValues() {
	//never assign position (pointer field) directly
	var _transformedPos *PositionT

	pad.newValueHandled = false

	_transformedPos, pad.shiftedAngle, pad.magnitude = pad.curPos.CalcTransformedPos(pad.zoneRotation)
	pad.transformedPos.Update(_transformedPos)

	pad.calcRadius()

	//pad.fromMaxPossiblePos.Update(pad.transformedPos.CalcFromMaxPossible(pad.radius))
}

type SetValueFuncT = func(fieldPointer *float64, value float64)

func (pad *PadStickPositionT) GetSetValueFunc() SetValueFuncT {
	notInit := math.IsNaN
	curPos := pad.curPos

	return func(fieldPointer *float64, value float64) {
		pad.Lock()
		defer pad.Unlock()

		*fieldPointer = value

		if value == 0 {
			x, y := curPos.x, curPos.y
			if (x == 0 && y == 0) ||
				(notInit(x) && y == 0) || (x == 0 && notInit(y)) {
				pad.Reset()
			}
		} else { // x != 0 && y != 0
			pad.ReCalculateValues()

			if pad.moveMouseSC != nil {
				pad.moveMouseSC()
			}
		}
	}
}

func (pad *PadStickPositionT) SetX(value float64) {
	pad.setValue(&(pad.curPos.x), value)
}

func (pad *PadStickPositionT) SetY(value float64) {
	pad.setValue(&(pad.curPos.y), value)
}

type ConvertRangeFuncT = func(input, outputMax float64) float64

func (pad *PadStickPositionT) GetConvertRangeFunc() ConvertRangeFuncT {
	FloatEqualityMargin := pad.cfg.Math.FloatEqualityMargin
	inputMin := pad.cfg.PadsSticks.Stick.DeadzoneDS
	outputMin := pad.cfg.Math.OutputMin

	return func(input, outputMax float64) float64 {
		gofuncs.PanicAnyNotInit(input)

		if input == 0 {
			return 0
		}

		isNegative, input := gofuncs.GetIsNegativeAndAbs(input)

		if input > pad.radius+FloatEqualityMargin {
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

func (pos *PositionT) CalcFromMaxPossible(radius float64) *PositionT {
	calcFromMaxPossible := func(x, y float64) float64 {
		maxPossibleX := math.Sqrt(gofuncs.Sqr(radius) - gofuncs.Sqr(y))
		if maxPossibleX == 0 {
			return 0
		}

		ratioFromMaxPossible := x / maxPossibleX

		if ratioFromMaxPossible > radius {
			if ratioFromMaxPossible > radius+pos.cfg.Math.FloatEqualityMargin {
				gofuncs.Panic("Incorrect calculations")
			}
			ratioFromMaxPossible = radius
		}
		return ratioFromMaxPossible
	}

	//important to use temp values then assign
	posFromMaxPossible := MakeEmptyPosition(pos.cfg)
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
