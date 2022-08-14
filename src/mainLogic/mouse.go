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

func moveMouseSC() {
	if Cfg.PadsSticksMode.GetMode() == TypingMode {
		return
	}

	if gofuncs.AnyNotInit(Cfg.mousePadStick.curPos.x, Cfg.mousePadStick.curPos.y) {
		return
	}

	//Cfg.mousePadStick.Lock()

	moveX := calcMove(Cfg.mousePadStick.transformedPos.x, Cfg.mousePadStick.prevMousePos.x)
	moveY := calcMove(Cfg.mousePadStick.transformedPos.y, Cfg.mousePadStick.prevMousePos.y)
	RightPadStick.UpdatePrevMousePos()

	//Cfg.mousePadStick.Unlock()

	if moveX != 0 || moveY != 0 {
		osSpec.MoveMouse(moveX, moveY)
	}
}

func RunMouseThreadDS() {
	ticker := time.NewTicker(Cfg.mouseInterval)
	for range ticker.C {
		//moveMouseSC()
	}
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

		if Cfg.PadsSticksMode.GetMode() != MouseMode {
			time.Sleep(scrollInterval)
			continue
		}

		Cfg.scrollPadStick.Lock()

		hDir, vDir := getDirections(Cfg.scrollPadStick.transformedPos.x, Cfg.scrollPadStick.transformedPos.y)

		scrollInterval = Cfg.scrollPadStick.calcRefreshInterval(Cfg.scrollPadStick.magnitude, Cfg.scrollSlowestInterval, Cfg.scrollFastestInterval)

		Cfg.scrollPadStick.Unlock()

		if hDir != 0 {
			osSpec.ScrollHorizontal(hDir)
		}
		if vDir != 0 {
			osSpec.ScrollVertical(vDir)
		}

		time.Sleep(scrollInterval)
	}
}
