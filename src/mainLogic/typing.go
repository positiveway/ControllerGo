package mainLogic

import (
	"ControllerGo/osSpec"
	"github.com/positiveway/gofuncs"
)

const NoneStr = "None"

type SticksPositionT [2]ZoneT
type TypingLayoutT map[SticksPositionT]int

type TypingT struct {
	CfgStruct
	LeftPS, RightPS *PadStickPositionT
	typeLetter      func()
}

func (typing *TypingT) Init(cfg *ConfigsT) {
	typing.CfgStruct.Init(cfg)

	typing.typeLetter = typing.GetTypeLetterFunc()
}

func (typing *TypingT) loadLayout() TypingLayoutT {
	linesParts := gofuncs.ReadLayoutFile(2,
		[]string{typing.cfg.Path.AllLayoutsDir, "typing.csv"})

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

func (typing *TypingT) genBoundariesMap() ZoneBoundariesMapT {
	cfg := typing.cfg
	AngleMargin := cfg.Typing.AngleMargin

	return genEqualThresholdBoundariesMap(true,
		MakeAngleMargin(
			AngleMargin.Diagonal,
			AngleMargin.Straight,
			AngleMargin.Straight),
		cfg.Typing.ThresholdPct,
		1.0)
}

func (typing *TypingT) GetTypeLetterFunc() func() {
	padsSticksMode := typing.cfg.PadsSticks.Mode
	LeftPS := typing.LeftPS
	RightPS := typing.RightPS

	boundariesMap := typing.genBoundariesMap()
	layout := typing.loadLayout()

	return func() {
		if padsSticksMode.CurrentMode != TypingMode {
			return
		}
		LeftPS.ReCalculateZone(boundariesMap)
		RightPS.ReCalculateZone(boundariesMap)

		if LeftPS.zoneCanBeUsed && RightPS.zoneCanBeUsed {
			if LeftPS.zoneChanged || RightPS.zoneChanged {
				if !LeftPS.awaitingCentralPosition || !RightPS.awaitingCentralPosition {
					LeftPS.awaitingCentralPosition = true
					RightPS.awaitingCentralPosition = true

					position := SticksPositionT{LeftPS.zone, RightPS.zone}
					if code, found := layout[position]; found {
						osSpec.TypeKey(code)
					}
				}
			}
		}
	}
}
