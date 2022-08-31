package mainLogic

import (
	"ControllerGo/osSpec"
	"github.com/positiveway/gofuncs"
	"math"
)

func calcMove(value, prevValue float64) int {
	if gofuncs.AnyNotInit(value, prevValue) {
		return 0
	}

	diff := value - prevValue
	pixels := gofuncs.FloatToIntRound[int](diff * Cfg.mouseSpeedSC)

	return pixels
}

func moveMouseSC() {
	if Cfg.PadsSticksMode.GetMode() == TypingMode {
		return
	}

	transformedPos := Cfg.mousePadStick.transformedPos
	prevMousePos := Cfg.mousePadStick.prevMousePos

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
type FilterMoveFunc = func(input float64, isX bool, padStick *PadStickPosition) float64

func calcMovement(input float64, isX bool, moveInterval *Interval,
	padStick *PadStickPosition, repetitionIntervals *RepetitionIntervals,
	filterFunc FilterMoveFunc) int {

	var moveByPixel int

	if moveInterval.DecreaseInterval() {
		moveInterval.SetInterval(repetitionIntervals.fastest)

		if filterFunc != nil {
			input = filterFunc(input, isX, padStick)
		}

		if !gofuncs.IsNotInitOrEmpty(input) {
			moveByPixel = gofuncs.SignAsInt(input)
			moveInterval.SetInterval(padStick.calcRefreshInterval(input, repetitionIntervals.slowest, repetitionIntervals.fastest))
		}
	}
	return moveByPixel
}

func MoveInInterval(
	moveIntervals *Intervals2,
	padStick *PadStickPosition, position *Position,
	repetitionIntervals *RepetitionIntervals,
	moveFunc MoveByPixelFunc, filterFunc FilterMoveFunc) {

	//slowestInterval := RepetitionsToInterval(minRepetitionPerSec)
	//fastestInterval :=RepetitionsToInterval(maxRepetitionPerSec)

	padStick.Lock()

	moveByPixelX := calcMovement(position.x, true, moveIntervals.X, padStick, repetitionIntervals, filterFunc)
	moveByPixelY := calcMovement(position.y, false, moveIntervals.Y, padStick, repetitionIntervals, filterFunc)

	padStick.Unlock()

	moveFunc(moveByPixelX, moveByPixelY)

}

func MoveMouse(moveIntervals *Intervals2) {
	mousePadStick := Cfg.mousePadStick
	position := mousePadStick.transformedPos

	MoveInInterval(moveIntervals, mousePadStick, position,
		Cfg.mouseIntervalsDS, moveMouseByPixelDS, nil)
}

func moveMouseByPixelDS(moveByPixelX, moveByPixelY int) {
	if moveByPixelX != 0 || moveByPixelY != 0 {
		osSpec.MoveMouse(moveByPixelX, moveByPixelY)
	}
}

func moveScrollByPixel(moveByPixelX, moveByPixelY int) {
	if moveByPixelX != 0 {
		osSpec.ScrollHorizontal(moveByPixelX)
	}
	if moveByPixelY != 0 {
		osSpec.ScrollVertical(moveByPixelY)
	}
}

func filterScrollHorizontal(input float64, isX bool, padStick *PadStickPosition) float64 {
	if isX && gofuncs.Abs(input) <= Cfg.scrollHorizontalThreshold*padStick.radius {
		input = 0
	}
	return input
}

func MoveScroll(moveIntervals *Intervals2) {
	scrollPadStick := Cfg.scrollPadStick
	position := scrollPadStick.transformedPos

	MoveInInterval(moveIntervals, scrollPadStick, position,
		Cfg.scrollIntervals, moveScrollByPixel, filterScrollHorizontal)
}

func getDirection(val float64, horizontal bool) int {
	if gofuncs.IsNotInit(val) {
		return 0
	}
	if horizontal && math.Abs(val) < Cfg.scrollHorizontalThreshold {
		return 0
	}
	return gofuncs.SignAsInt(val)
}

func getDirections(x, y float64) (int, int) {
	hDir := getDirection(x, true)
	vDir := getDirection(y, false)

	if hDir != 0 {
		vDir = 0
	}
	return hDir, vDir
}
