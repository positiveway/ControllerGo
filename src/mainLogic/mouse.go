package mainLogic

import (
	"ControllerGo/osSpec"
	"github.com/positiveway/gofuncs"
)

var MoveMouseSC func()

func GetMoveMouseSCFunc() func() {
	transformedPos := Cfg.PadsSticks.MousePS.transformedPos
	prevMousePos := Cfg.PadsSticks.MousePS.prevMousePos
	highPrecisionMode := Cfg.PadsSticks.HighPrecisionMode

	calcPixelsToMoveMouse := func(value, prevValue float64) int {
		if gofuncs.AnyNotInit(value, prevValue) {
			return 0
		}

		diff := value - prevValue
		pixels := gofuncs.FloatToIntRound[int](diff * highPrecisionMode.curMouseSpeed)

		return pixels
	}

	return func() {
		if Cfg.PadsSticks.Mode.GetMode() == TypingMode {
			return
		}

		moveX := calcPixelsToMoveMouse(transformedPos.x, prevMousePos.x)
		moveY := calcPixelsToMoveMouse(transformedPos.y, prevMousePos.y)
		prevMousePos.Update(transformedPos)

		if moveX != 0 || moveY != 0 {
			osSpec.MoveMouse(moveX, moveY)
		}
	}
}

//func RepetitionsToInterval(repetitions float64) float64 {
//	return 1000 / repetitions
//}

type MoveByPixelFuncT = func(moveByPixelX, moveByPixelY int)
type FilterMoveFuncT = func(input float64, isX bool) float64

func GetMoveInInterval(
	padStick *PadStickPositionT, position *PositionT,
	moveFunc MoveByPixelFuncT, filterFunc FilterMoveFuncT) func(repetitionIntervals *IntervalRangeT) {

	calcMovement := func(input float64, isX bool, moveInterval *IntervalTimerT, repetitionIntervals *IntervalRangeT) int {
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

	intervalTimers := MakeIntervalTimers2()

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

func GetScrollMoveFunc() MoveByPixelFuncT {
	highPrecisionMode := Cfg.PadsSticks.HighPrecisionMode

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

func GetScrollFilterFunc() FilterMoveFuncT {
	scrollFilterValue := Cfg.Scroll.HorizontalThresholdPct
	scrollPadStick := Cfg.PadsSticks.ScrollPS

	return func(input float64, isX bool) float64 {
		if isX && gofuncs.Abs(input) <= scrollFilterValue*scrollPadStick.radius {
			input = 0
		}
		return input
	}
}

func GetMouseMoveFunc() MoveByPixelFuncT {
	return osSpec.MoveMouse
}
