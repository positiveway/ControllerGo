package mainLogic

import (
	"ControllerGo/osSpec"
	"github.com/positiveway/gofuncs"
	"math"
	"time"
)

type MoveByPixelFuncT = func(moveByPixelX, moveByPixelY int)
type FilterMoveFuncT = func(input float64, isX bool) float64

func (pad *PadStickPositionT) GetMoveOrScrollSCFunc(initialActionsFunc func(), moveByPixelFunc MoveByPixelFuncT) func(speed float64) {
	transformedPos := pad.transformedPos
	prevMousePos := pad.prevMousePos

	notInit := math.IsNaN

	calcPixelsToMoveMouse := func(value, prevValue, speed float64) int {
		if notInit(value) || notInit(prevValue) {
			return 0
		}
		diff := value - prevValue
		pixels := gofuncs.FloatToIntRound[int](diff * speed)

		return pixels
	}

	return func(speed float64) {
		if initialActionsFunc != nil {
			initialActionsFunc()
		}

		moveX := calcPixelsToMoveMouse(transformedPos.x, prevMousePos.x, speed)
		moveY := calcPixelsToMoveMouse(transformedPos.y, prevMousePos.y, speed)
		prevMousePos.Update(transformedPos)

		if moveX != 0 || moveY != 0 {
			moveByPixelFunc(moveX, moveY)
		}
	}
}

func (pad *PadStickPositionT) GetMoveMouseSCFunc(highPrecisionMode *HighPrecisionModeT) func() {
	transformedPos := pad.transformedPos
	prevMousePos := pad.prevMousePos

	cfg := pad.cfg
	doubleTouchMaxInterval := gofuncs.NumberToMillis(cfg.Mouse.DoubleTouchMaxInterval)

	buttons := pad.buttons
	leftClickBtn, leftClickCmdInfo := pad.leftClickBtn, pad.leftClickCmdInfo

	notInit := math.IsNaN
	isInit := func(value float64) bool { return !math.IsNaN(value) }

	doubleClickFunc := func() {
		if notInit(prevMousePos.x) || notInit(prevMousePos.y) {
			if isInit(transformedPos.x) && isInit(transformedPos.y) {
				diff := time.Now().Sub(pad.firstTouchTime)
				if diff < doubleTouchMaxInterval {
					//fmt.Println("click")
					buttons.pressIfNotAlready(leftClickBtn, leftClickCmdInfo)
				} else {
					pad.firstTouchTime = time.Now()
				}
			}
		}
	}

	moveMouseFunc := pad.GetMoveOrScrollSCFunc(doubleClickFunc, GetMouseMoveFunc())

	return func() {
		moveMouseFunc(highPrecisionMode.curMouseSpeed)
	}
}

func (pad *PadStickPositionT) GetScrollSCFunc(highPrecisionMode *HighPrecisionModeT) func() {
	moveByPixelFunc := func(moveByPixelX, moveByPixelY int) {
		if moveByPixelX != 0 {
			osSpec.ScrollHorizontal(moveByPixelX)
		}
		if moveByPixelY != 0 {
			osSpec.ScrollVertical(moveByPixelY)
		}
	}

	transformedPos := pad.transformedPos
	prevMousePos := pad.prevMousePos

	notInit := math.IsNaN
	isInit := func(value float64) bool { return !math.IsNaN(value) }

	pressCtrlIfZoomFunc := func() {
		if notInit(prevMousePos.x) || notInit(prevMousePos.y) {
			if isInit(transformedPos.x) && isInit(transformedPos.y) {
				highPrecisionMode.PressCtrl()
			}
		}
	}

	scrollFunc := pad.GetMoveOrScrollSCFunc(pressCtrlIfZoomFunc, moveByPixelFunc)

	return func() {
		scrollFunc(highPrecisionMode.curScrollSpeed)
	}
}

//func RepetitionsToInterval(repetitions float64) float64 {
//	return 1000 / repetitions
//}

func GetMoveInInterval(cfg *ConfigsT,
	padStick *PadStickPositionT, position *PositionT,
	moveFunc MoveByPixelFuncT, filterFunc FilterMoveFuncT) func(repetitionIntervals *IntervalRangeT) {

	calcMovement := func(input float64, isX bool, moveInterval *RepeatedTimerT, repetitionIntervals *IntervalRangeT) int {
		var moveByPixel int

		if moveInterval.DecreaseInterval() {
			moveInterval.SetInterval(repetitionIntervals.Fastest)

			if filterFunc != nil {
				input = filterFunc(input, isX)
			}

			if !gofuncs.IsEmptyOrNotInit(input) {
				moveByPixel = gofuncs.SignAsInt(input)
				moveInterval.SetInterval(padStick.calcRefreshInterval(input, repetitionIntervals.Slowest, repetitionIntervals.Fastest))
			}
		}
		return moveByPixel
	}

	intervalTimers := MakeIntervalTimers2(cfg)

	return func(repetitionIntervals *IntervalRangeT) {
		//slowestInterval := RepetitionsToInterval(minRepetitionPerSec)
		//fastestInterval := RepetitionsToInterval(maxRepetitionPerSec)

		padStick.Lock()

		moveByPixelX := calcMovement(position.x, true, intervalTimers.X, repetitionIntervals)
		moveByPixelY := calcMovement(position.y, false, intervalTimers.Y, repetitionIntervals)

		padStick.Unlock()

		if moveByPixelX != 0 || moveByPixelY != 0 {
			moveFunc(moveByPixelX, moveByPixelY)
		}
	}
}

func (dependentVars *DependentVariablesT) GetScrollMoveDSFunc() MoveByPixelFuncT {
	highPrecisionMode := dependentVars.HighPrecisionMode

	return func(moveByPixelX, moveByPixelY int) {
		highPrecisionMode.PressCtrl()

		if moveByPixelX != 0 {
			osSpec.ScrollHorizontal(moveByPixelX)
		}
		if moveByPixelY != 0 {
			osSpec.ScrollVertical(moveByPixelY)
		}
	}
}

func (dependentVars *DependentVariablesT) GetScrollFilterDSFunc() FilterMoveFuncT {
	scrollFilterValue := dependentVars.cfg.Scroll.HorizontalThresholdPct
	scrollPadStick := dependentVars.ScrollPS

	return func(input float64, isX bool) float64 {
		if isX && math.Abs(input) <= scrollFilterValue*scrollPadStick.radius {
			input = 0
		}
		return input
	}
}

func GetMouseMoveFunc() MoveByPixelFuncT {
	return osSpec.MoveMouse
}
