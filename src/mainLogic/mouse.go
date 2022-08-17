package mainLogic

import (
	"ControllerGo/osSpec"
	"github.com/positiveway/gofuncs"
	"math"
	"time"
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

type MoveByPixelFunc = func(moveByPixel int, isX bool)
type FilterMoveFunc = func(input float64, isX bool, padStick *PadStickPosition) float64

func moveInInterval(isX bool,
	padStick *PadStickPosition, position *Position,
	repetitionIntervals *RepetitionIntervals,
	allowedModes []ModeT, filterFunc FilterMoveFunc, moveFunc MoveByPixelFunc) {

	for {
		//slowestInterval := RepetitionsToInterval(minRepetitionPerSec)
		//fastestInterval :=RepetitionsToInterval(maxRepetitionPerSec)

		moveInterval := gofuncs.NumberToMillis(repetitionIntervals.fastest)

		if gofuncs.Contains(allowedModes, Cfg.PadsSticksMode.GetMode()) {
			padStick.Lock()

			var input float64
			if isX {
				input = position.x
			} else {
				input = position.y
			}

			if filterFunc != nil {
				input = filterFunc(input, isX, padStick)
			}

			if !gofuncs.IsNotInitOrEmpty(input) {
				moveByPixel := gofuncs.SignAsInt(input)
				moveInterval = padStick.calcRefreshInterval(input, repetitionIntervals.slowest, repetitionIntervals.fastest)
				moveFunc(moveByPixel, isX)
			}

			padStick.Unlock()
		}
		time.Sleep(moveInterval)
	}
}

func runMoveThreads(padStick *PadStickPosition, position *Position,
	repetitionIntervals *RepetitionIntervals,
	allowedModes []ModeT, filterFunc FilterMoveFunc, moveFunc MoveByPixelFunc) {

	go moveInInterval(true, padStick, position, repetitionIntervals, allowedModes, filterFunc, moveFunc)
	go moveInInterval(false, padStick, position, repetitionIntervals, allowedModes, filterFunc, moveFunc)
}

func moveMouseByPixelDS(moveByPixel int, isX bool) {
	if isX {
		osSpec.MoveMouse(moveByPixel, 0)
	} else {
		osSpec.MoveMouse(0, moveByPixel)
	}
}

func RunMouseThreadsDS() {
	mousePadStick := Cfg.mousePadStick
	position := mousePadStick.transformedPos

	runMoveThreads(mousePadStick, position, Cfg.mouseIntervalsDS, Cfg.MouseAllowedMods, nil, moveMouseByPixelDS)
}

func moveScrollByPixel(moveByPixel int, isX bool) {
	if isX {
		osSpec.ScrollHorizontal(moveByPixel)
	} else {
		osSpec.ScrollVertical(moveByPixel)
	}
}

func filterScrollHorizontal(input float64, isX bool, padStick *PadStickPosition) float64 {
	if isX && gofuncs.Abs(input) <= Cfg.scrollHorizontalThreshold*padStick.radius {
		input = 0
	}
	return input
}

func RunScrollThreads() {
	scrollPadStick := Cfg.scrollPadStick
	position := scrollPadStick.transformedPos

	runMoveThreads(scrollPadStick, position, Cfg.scrollIntervals, Cfg.ScrollAllowedMods, filterScrollHorizontal, moveScrollByPixel)

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
