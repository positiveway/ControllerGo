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

type MoveByPixelFunc = func(moveByPixelX, moveByPixelY int)
type FilterMoveFunc = func(input float64, isX bool, padStick *PadStickPosition) float64

func calcMovement(input float64, isX bool, moveInterval, tickerInterval float64,
	padStick *PadStickPosition, repetitionIntervals *RepetitionIntervals,
	filterFunc FilterMoveFunc) (float64, int) {

	var moveByPixel int

	moveInterval -= tickerInterval
	if moveInterval <= 0 {
		moveInterval = repetitionIntervals.fastest

		if filterFunc != nil {
			input = filterFunc(input, isX, padStick)
		}

		if !gofuncs.IsNotInitOrEmpty(input) {
			moveByPixel = gofuncs.SignAsInt(input)
			moveInterval = padStick.calcRefreshInterval(input, repetitionIntervals.slowest, repetitionIntervals.fastest)
		}
	}
	return moveInterval, moveByPixel
}
func RunMoveInIntervalThread(
	tickerInterval float64,
	padStick *PadStickPosition, position *Position,
	repetitionIntervals *RepetitionIntervals,
	allowedModes []ModeT, filterFunc FilterMoveFunc, moveFunc MoveByPixelFunc) {

	go func() {
		var moveIntervalX, moveIntervalY float64

		ticker := time.NewTicker(gofuncs.NumberToMillis(tickerInterval))

		for range ticker.C {
			if gofuncs.Contains(allowedModes, Cfg.PadsSticksMode.GetMode()) {
				//slowestInterval := RepetitionsToInterval(minRepetitionPerSec)
				//fastestInterval :=RepetitionsToInterval(maxRepetitionPerSec)

				var moveByPixelX, moveByPixelY int

				padStick.Lock()

				moveIntervalX, moveByPixelX = calcMovement(position.x, true, moveIntervalX, tickerInterval, padStick, repetitionIntervals, filterFunc)
				moveIntervalY, moveByPixelY = calcMovement(position.y, false, moveIntervalY, tickerInterval, padStick, repetitionIntervals, filterFunc)

				padStick.Unlock()

				moveFunc(moveByPixelX, moveByPixelY)
			}
		}
	}()
}

func moveMouseByPixelDS(moveByPixelX, moveByPixelY int) {
	if moveByPixelX != 0 || moveByPixelY != 0 {
		osSpec.MoveMouse(moveByPixelX, moveByPixelY)
	}
}

func RunMouseThreadsDS() {
	mousePadStick := Cfg.mousePadStick
	position := mousePadStick.transformedPos

	RunMoveInIntervalThread(1, mousePadStick, position, Cfg.mouseIntervalsDS, Cfg.MouseAllowedMods, nil, moveMouseByPixelDS)
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

func RunScrollThreads() {
	scrollPadStick := Cfg.scrollPadStick
	position := scrollPadStick.transformedPos

	RunMoveInIntervalThread(10, scrollPadStick, position, Cfg.scrollIntervals, Cfg.ScrollAllowedMods, filterScrollHorizontal, moveScrollByPixel)

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
