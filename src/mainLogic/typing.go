package mainLogic

import (
	"ControllerGo/osSpec"
	"github.com/positiveway/gofuncs"
)

const NoneStr = "None"

type SticksPosition [2]Zone
type TypingLayout map[SticksPosition]int

var TypingBoundariesMap ZoneBoundariesMap
var typingLayout TypingLayout

func initTyping() {
	TypingBoundariesMap = genTypingBoundariesMap()
	typingLayout = loadTypingLayout()
}

func loadTypingLayout() TypingLayout {
	linesParts := Cfg.ReadLayoutFile("typing.csv", 2)

	layout := TypingLayout{}
	for _, parts := range linesParts {
		leftStick, rightStick, letter := Zone(parts[0]), Zone(parts[1]), parts[2]
		if !gofuncs.Contains(AllZones, leftStick) {
			gofuncs.PanicMisspelled(leftStick)
		}
		if !gofuncs.Contains(AllZones, rightStick) {
			gofuncs.PanicMisspelled(rightStick)
		}
		if letter == NoneStr {
			continue
		}
		code := getCodeFromLetter(letter)
		position := SticksPosition{leftStick, rightStick}
		gofuncs.AssignWithDuplicateCheck(layout, position, code)
	}
	return layout
}

func genTypingBoundariesMap() ZoneBoundariesMap {
	return genEqualThresholdBoundariesMap(true,
		makeAngleMargin(Cfg.TypingDiagonalAngleMargin, Cfg.TypingStraightAngleMargin, Cfg.TypingStraightAngleMargin),
		Cfg.TypingThreshold,
		Cfg.MinStandardPadRadius)
}

func TypeLetter() {
	if Cfg.PadsSticksMode.GetMode() != TypingMode {
		return
	}
	LeftPad.ReCalculateZone(TypingBoundariesMap)
	RightPadStick.ReCalculateZone(TypingBoundariesMap)

	if LeftPad.zoneCanBeUsed && RightPadStick.zoneCanBeUsed {
		if LeftPad.zoneChanged || RightPadStick.zoneChanged {
			if !LeftPad.awaitingCentralPosition || !RightPadStick.awaitingCentralPosition {
				LeftPad.awaitingCentralPosition = true
				RightPadStick.awaitingCentralPosition = true

				position := SticksPosition{LeftPad.zone, RightPadStick.zone}
				if code, found := typingLayout[position]; found {
					osSpec.TypeKey(code)
				}
			}
		}
	}
}
