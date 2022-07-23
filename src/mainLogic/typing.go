package mainLogic

import "ControllerGo/src/osSpec"

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
	linesParts := ReadLayoutFile("typing.csv", 2)

	layout := TypingLayout{}
	for _, parts := range linesParts {
		leftStick, rightStick, letter := Zone(parts[0]), Zone(parts[1]), parts[2]
		if !contains(AllZones, leftStick) {
			PanicMisspelled(leftStick)
		}
		if !contains(AllZones, rightStick) {
			PanicMisspelled(rightStick)
		}
		if letter == NoneStr {
			continue
		}
		code := getCodeFromLetter(letter)
		position := SticksPosition{leftStick, rightStick}
		AssignWithDuplicateCheck(layout, position, code)
	}
	return layout
}

func genTypingBoundariesMap() ZoneBoundariesMap {
	return genEqualThresholdBoundariesMap(true,
		makeAngleMargin(TypingDiagonalAngleMargin, TypingStraightAngleMargin, TypingStraightAngleMargin),
		TypingThreshold,
		PadRadius)
}

var PrintTypingDebugInfo = false

func TypeLetter() {
	if padsMode.GetMode() != TypingMode {
		return
	}
	LeftPad.ReCalculateZone(TypingBoundariesMap, PrintTypingDebugInfo)
	RightPad.ReCalculateZone(TypingBoundariesMap, PrintTypingDebugInfo)

	if LeftPad.zoneCanBeUsed && RightPad.zoneCanBeUsed {
		if LeftPad.zoneChanged || RightPad.zoneChanged {
			if !LeftPad.awaitingCentralPostion || !RightPad.awaitingCentralPostion {
				LeftPad.awaitingCentralPostion = true
				RightPad.awaitingCentralPostion = true

				position := SticksPosition{LeftPad.zone, RightPad.zone}
				if code, found := typingLayout[position]; found {
					osSpec.TypeKey(code)
				}
			}
		}
	}
}
