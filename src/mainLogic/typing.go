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
		gofuncs.AssignWithDuplicateCheck(layout, position, code)
	}
	return layout
}

func genTypingBoundariesMap() ZoneBoundariesMapT {
	return genEqualThresholdBoundariesMap(true,
		MakeAngleMargin(
			Cfg.Typing.AngleMargin.Diagonal,
			Cfg.Typing.AngleMargin.Straight,
			Cfg.Typing.AngleMargin.Straight),
		Cfg.Typing.ThresholdPct,
		1.0)
}

func TypeLetter() {
	if Cfg.PadsSticks.Mode.GetMode() != TypingMode {
		return
	}
	Cfg.Typing.LeftPS.ReCalculateZone(TypingBoundariesMap)
	Cfg.Typing.RightPS.ReCalculateZone(TypingBoundariesMap)

	if Cfg.Typing.LeftPS.zoneCanBeUsed && Cfg.Typing.RightPS.zoneCanBeUsed {
		if Cfg.Typing.LeftPS.zoneChanged || Cfg.Typing.RightPS.zoneChanged {
			if !Cfg.Typing.LeftPS.awaitingCentralPosition || !Cfg.Typing.RightPS.awaitingCentralPosition {
				Cfg.Typing.LeftPS.awaitingCentralPosition = true
				Cfg.Typing.RightPS.awaitingCentralPosition = true

				position := SticksPositionT{Cfg.Typing.LeftPS.zone, Cfg.Typing.RightPS.zone}
				if code, found := typingLayout[position]; found {
					osSpec.TypeKey(code)
				}
			}
		}
	}
}
