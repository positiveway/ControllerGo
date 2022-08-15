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

func moveMouseDS(input float64, isX bool) {
	mouseMoveInterval := gofuncs.NumberToMillis(Cfg.mouseFastestIntervalDS)

	if Cfg.PadsSticksMode.GetMode() == TypingMode {
		time.Sleep(mouseMoveInterval)
	}

	mousePadStick := Cfg.mousePadStick

	mousePadStick.Lock()

	if input != 0 {
		moveMouse := gofuncs.SignAsInt(input)
		mouseMoveInterval = mousePadStick.calcRefreshInterval(input, Cfg.mouseSlowestIntervalDS, Cfg.mouseFastestIntervalDS)
		if isX {
			osSpec.MoveMouse(moveMouse, 0)
		} else {
			osSpec.MoveMouse(0, moveMouse)
		}
	}

	mousePadStick.Unlock()

	time.Sleep(mouseMoveInterval)
}

func RunMouseThreadXDS() {
	for {
		transformedPos := Cfg.mousePadStick.transformedPos

		moveMouseDS(transformedPos.x, true)
	}
}

func RunMouseThreadYDS() {
	for {
		transformedPos := Cfg.mousePadStick.transformedPos

		moveMouseDS(transformedPos.y, false)
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

		scrollPadStick := Cfg.scrollPadStick
		transformedPos := scrollPadStick.transformedPos

		scrollPadStick.Lock()

		hDir, vDir := getDirections(transformedPos.x, transformedPos.y)

		if scrollPadStick.magnitude != 0 {
			scrollInterval = scrollPadStick.calcRefreshInterval(scrollPadStick.magnitude, Cfg.scrollSlowestInterval, Cfg.scrollFastestInterval)
		}

		scrollPadStick.Unlock()

		if hDir != 0 {
			osSpec.ScrollHorizontal(hDir)
		}
		if vDir != 0 {
			osSpec.ScrollVertical(vDir)
		}

		time.Sleep(scrollInterval)
	}
}
