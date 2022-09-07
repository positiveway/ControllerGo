package mainLogic

import (
	"ControllerGo/osSpec"
	"github.com/positiveway/gofuncs"
)

func calcMove(value, prevValue float64) int {
	if gofuncs.AnyNotInit(value, prevValue) {
		return 0
	}

	diff := value - prevValue
	pixels := gofuncs.FloatToIntRound[int](diff * Cfg.PadsSticks.HighPrecisionMode.curMouseSpeed)

	return pixels
}

func moveMouseSC() {
	if Cfg.PadsSticks.Mode.GetMode() == TypingMode {
		return
	}

	transformedPos := Cfg.PadsSticks.MousePS.transformedPos
	prevMousePos := Cfg.PadsSticks.MousePS.prevMousePos

	moveX := calcMove(transformedPos.x, prevMousePos.x)
	moveY := calcMove(transformedPos.y, prevMousePos.y)
	prevMousePos.Update(transformedPos)

	if moveX != 0 || moveY != 0 {
		osSpec.MoveMouse(moveX, moveY)
	}
}

//func RepetitionsToInterval(repetitions float64) float64 {
//	return 1000 / repetitions
//}

type MoveByPixelFunc = func(moveByPixelX, moveByPixelY int)
type FilterMoveFunc = func(input float64, isX bool, padStick *PadStickPositionT) float64

func calcMovement(input float64, isX bool, moveInterval *IntervalTimerT,
	padStick *PadStickPositionT, repetitionIntervals *IntervalRangeT,
	filterFunc FilterMoveFunc) int {

	var moveByPixel int

	if moveInterval.DecreaseInterval() {
		moveInterval.SetInterval(repetitionIntervals.Fastest)

		if filterFunc != nil {
			input = filterFunc(input, isX, padStick)
		}

		if !gofuncs.IsEmptyOrNotInit(input) {
			moveByPixel = gofuncs.SignAsInt(input)
			moveInterval.SetInterval(padStick.calcRefreshInterval(input, repetitionIntervals.Slowest, repetitionIntervals.Fastest))
		}
	}
	return moveByPixel
}

func MoveInInterval(
	intervalTimers *IntervalTimers2T,
	padStick *PadStickPositionT, position *PositionT,
	repetitionIntervals *IntervalRangeT,
	moveFunc MoveByPixelFunc, filterFunc FilterMoveFunc) {

	//slowestInterval := RepetitionsToInterval(minRepetitionPerSec)
	//fastestInterval := RepetitionsToInterval(maxRepetitionPerSec)

	padStick.Lock()

	moveByPixelX := calcMovement(position.x, true, intervalTimers.X, padStick, repetitionIntervals, filterFunc)
	moveByPixelY := calcMovement(position.y, false, intervalTimers.Y, padStick, repetitionIntervals, filterFunc)

	padStick.Unlock()

	if moveByPixelX != 0 || moveByPixelY != 0 {
		moveFunc(moveByPixelX, moveByPixelY)
	}

	//dirty hack to determine scroll
	if filterFunc != nil {
		if position.x == 0 && position.y == 0 {
			Cfg.PadsSticks.HighPrecisionMode.ReleaseCtrl()
		}
	}
}

//
//func MoveMouse(intervalTimers *IntervalTimers2T, moveIntervalRange *IntervalRangeT) {
//	mousePadStick := Cfg.PadsSticks.MousePS
//	mousePosition := mousePadStick.transformedPos
//
//	MoveInInterval(intervalTimers, mousePadStick, mousePosition,
//		moveIntervalRange, moveMouseByPixelDS, nil)
//}

func moveMouseByPixelDS(moveByPixelX, moveByPixelY int) {
	osSpec.MoveMouse(moveByPixelX, moveByPixelY)
}

func moveScrollByPixel(moveByPixelX, moveByPixelY int) {
	Cfg.PadsSticks.HighPrecisionMode.PressCtrl()

	if moveByPixelX != 0 {
		osSpec.ScrollHorizontal(moveByPixelX)
	}
	if moveByPixelY != 0 {
		osSpec.ScrollVertical(moveByPixelY)
	}
}

func filterScrollHorizontal(input float64, isX bool, padStick *PadStickPositionT) float64 {
	if isX && gofuncs.Abs(input) <= Cfg.Scroll.HorizontalThresholdPct*padStick.radius {
		input = 0
	}
	return input
}

//
//func MoveScroll(intervalTimers *IntervalTimers2T, moveIntervalRange *IntervalRangeT) {
//	scrollPadStick := Cfg.PadsSticks.ScrollPS
//	scrollPosition := scrollPadStick.transformedPos
//
//	MoveInInterval(intervalTimers, scrollPadStick, scrollPosition,
//		moveIntervalRange, moveScrollByPixel, filterScrollHorizontal)
//}
