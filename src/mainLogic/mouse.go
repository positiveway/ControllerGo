package mainLogic

import (
	"ControllerGo/osSpec"
	"github.com/positiveway/gofuncs"
	"math"
	"time"
)

func calcMove(value, prevValue float64) int {
	if gofuncs.IsNotInit(prevValue) {
		return 0
	}

	diff := value - prevValue
	pixels := gofuncs.FloatToIntRound[int](diff * mouseSpeed)

	return pixels
}

func RunMouseThread() {
	ticker := time.NewTicker(mouseInterval)
	for range ticker.C {
		if padsMode.GetMode() == TypingMode {
			continue
		}

		RightPad.Lock()

		moveX := calcMove(RightPad.actualX, RightPad.prevX)
		moveY := calcMove(RightPad.actualY, RightPad.prevY)
		RightPad.UpdatePrevValues()

		RightPad.Unlock()

		if moveX != 0 || moveY != 0 {
			osSpec.MoveMouse(moveX, moveY)
		}
	}
}

func calcScrollInterval(input float64) time.Duration {
	return calcRefreshInterval(input, scrollSlowestInterval, scrollFastestInterval)
}

func getDirection(val float64, horizontal bool) int {
	if gofuncs.IsNotInit(val) {
		return 0
	}
	if horizontal && math.Abs(val) < horizontalScrollThreshold {
		return 0
	}
	return gofuncs.SignAsNumber(val)
}

func getDirections(x, y float64) (int, int) {
	hDir := getDirection(x, true)
	vDir := getDirection(y, false)

	if hDir != 0 {
		vDir = 0
	}
	return hDir, vDir
}

func RunScrollThread() {
	for {
		scrollInterval := gofuncs.NumberToMillis(scrollFastestInterval)

		if padsMode.GetMode() != ScrollingMode {
			time.Sleep(scrollInterval)
			continue
		}

		LeftPad.Lock()

		LeftPad.CalcCoordsFromMaxPossible()
		hDir, vDir := getDirections(LeftPad.actualX, LeftPad.actualY)

		if LeftPad.magnitude != 0 {
			scrollInterval = calcScrollInterval(LeftPad.magnitude)
		}

		LeftPad.Unlock()

		if hDir != 0 {
			osSpec.ScrollHorizontal(hDir)
		}
		if vDir != 0 {
			osSpec.ScrollVertical(vDir)
		}

		time.Sleep(scrollInterval)
	}
}
