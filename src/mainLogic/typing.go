package mainLogic

import (
	"ControllerGo/osSpec"
	"github.com/positiveway/gofuncs"
)

const NoneStr = "None"

type SticksPositionT [2]ZoneT
type TypingLayoutT map[SticksPositionT]int

var TypingBoundariesMap ZoneBoundariesMapT
var typingLayout TypingLayoutT

func initTyping() {
	TypingBoundariesMap = genTypingBoundariesMap()
	typingLayout = loadTypingLayout()
}

func loadTypingLayout() TypingLayoutT {
	linesParts := Cfg.ReadLayoutFile("typing.csv", 2)

	layout := TypingLayoutT{}
	for _, parts := range linesParts {
		leftPadStickZone, rightPadStickZone, letter := ZoneT(parts[0]), ZoneT(parts[1]), parts[2]
		if !gofuncs.Contains(AllZones, leftPadStickZone) {
			gofuncs.PanicMisspelled(leftPadStickZone)
		}
		if !gofuncs.Contains(AllZones, rightPadStickZone) {
			gofuncs.PanicMisspelled(rightPadStickZone)
		}
		if letter == NoneStr {
			continue
		}
		code := getCodeFromLetter(letter)
		position := SticksPositionT{leftPadStickZone, rightPadStickZone}
		gofuncs.AssignWithDuplicateCheck(layout, position, code)
	}
	return layout
}

func genTypingBoundariesMap() ZoneBoundariesMapT {
	return genEqualThresholdBoundariesMap(true,
		makeAngleMargin(Cfg.TypingDiagonalAngleMargin, Cfg.TypingStraightAngleMargin, Cfg.TypingStraightAngleMargin),
		Cfg.TypingThreshold,
		Cfg.MinStandardPadRadius)
}

func TypeLetter() {
	if Cfg.PadsSticksMode.GetMode() != TypingMode {
		return
	}
	Cfg.LeftTypingPS.ReCalculateZone(TypingBoundariesMap)
	Cfg.RightTypingPS.ReCalculateZone(TypingBoundariesMap)

	if Cfg.LeftTypingPS.zoneCanBeUsed && Cfg.RightTypingPS.zoneCanBeUsed {
		if Cfg.LeftTypingPS.zoneChanged || Cfg.RightTypingPS.zoneChanged {
			if !Cfg.LeftTypingPS.awaitingCentralPosition || !Cfg.RightTypingPS.awaitingCentralPosition {
				Cfg.LeftTypingPS.awaitingCentralPosition = true
				Cfg.RightTypingPS.awaitingCentralPosition = true

				position := SticksPositionT{Cfg.LeftTypingPS.zone, Cfg.RightTypingPS.zone}
				if code, found := typingLayout[position]; found {
					osSpec.TypeKey(code)
				}
			}
		}
	}
}
