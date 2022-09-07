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
	TypeLetter = GetTypeLetterFunc()
	TypingBoundariesMap = genTypingBoundariesMap()
	typingLayout = loadTypingLayout()
}

func loadTypingLayout() TypingLayoutT {
	linesParts := gofuncs.ReadLayoutFile(2,
		[]string{Cfg.Path.AllLayoutsDir, "typing.csv"})

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
		gofuncs.AssignWithDuplicateKeyValueCheck(layout, position, code, false)
	}
	return layout
}

func genTypingBoundariesMap() ZoneBoundariesMapT {
	AngleMargin := Cfg.Typing.AngleMargin
	return genEqualThresholdBoundariesMap(true,
		MakeAngleMargin(
			AngleMargin.Diagonal,
			AngleMargin.Straight,
			AngleMargin.Straight),
		Cfg.Typing.ThresholdPct,
		1.0)
}

func GetTypeLetterFunc() func() {
	padsSticksMode := Cfg.PadsSticks.Mode
	LeftPS := Cfg.Typing.LeftPS
	RightPS := Cfg.Typing.RightPS

	return func() {
		if padsSticksMode.GetMode() != TypingMode {
			return
		}
		LeftPS.ReCalculateZone(TypingBoundariesMap)
		RightPS.ReCalculateZone(TypingBoundariesMap)

		if LeftPS.zoneCanBeUsed && RightPS.zoneCanBeUsed {
			if LeftPS.zoneChanged || RightPS.zoneChanged {
				if !LeftPS.awaitingCentralPosition || !RightPS.awaitingCentralPosition {
					LeftPS.awaitingCentralPosition = true
					RightPS.awaitingCentralPosition = true

					position := SticksPositionT{LeftPS.zone, RightPS.zone}
					if code, found := typingLayout[position]; found {
						osSpec.TypeKey(code)
					}
				}
			}
		}
	}
}

var TypeLetter func()
