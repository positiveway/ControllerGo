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
	pixels := gofuncs.FloatToIntRound[int](diff * Cfg.mouseSpeed)

	return pixels
}

func RunMouseThread() {
	ticker := time.NewTicker(Cfg.mouseInterval)
	for range ticker.C {
		if Cfg.padsMode.GetMode() == TypingMode {
			continue
		}

		RightPad.Lock()

		moveX := calcMove(RightPad.shiftedPos.x, RightPad.prevPos.x)
		moveY := calcMove(RightPad.shiftedPos.y, RightPad.prevPos.y)
		RightPad.UpdatePrevValues()

		RightPad.Unlock()

		if moveX != 0 || moveY != 0 {
			osSpec.MoveMouse(moveX, moveY)
		}
	}
}

func calcScrollInterval(input float64) time.Duration {
	return calcRefreshInterval(input, Cfg.scrollSlowestInterval, Cfg.scrollFastestInterval)
}

func getDirection(val float64, horizontal bool) int {
	if gofuncs.IsNotInit(val) {
		return 0
	}
	if horizontal && math.Abs(val) < Cfg.horizontalScrollThreshold {
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

func RunScrollThread() {
	for {
		scrollInterval := gofuncs.NumberToMillis(Cfg.scrollFastestInterval)

		if Cfg.padsMode.GetMode() != ScrollingMode {
			time.Sleep(scrollInterval)
			continue
		}

		LeftPad.Lock()

		hDir, vDir := getDirections(LeftPad.fromMaxPossiblePos.x, LeftPad.fromMaxPossiblePos.y)

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
