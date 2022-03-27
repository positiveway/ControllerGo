package mainLogic

import "ControllerGo/src/osSpec"

const NoneStr = "None"

type SticksPosition [2]Zone
type TypingLayout map[SticksPosition]int

var TypingBoundariesMap BoundariesMap

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

var typingInitBoundaries = InitBoundaries{
	AngleRight:     ZoneRight,
	AngleUpRight:   ZoneUpRight,
	AngleUp:        ZoneUp,
	AngleUpLeft:    ZoneUpLeft,
	AngleLeft:      ZoneLeft,
	AngleDownLeft:  ZoneDownLeft,
	AngleDown:      ZoneDown,
	AngleDownRight: ZoneDownRight,
}

func genTypingBoundariesMap() BoundariesMap {
	return genBoundariesMap(typingInitBoundaries,
		makeAngleMargin(TypingDiagonalAngleMargin, TypingStraightAngleMargin, TypingStraightAngleMargin),
		makeThreshold(TypingThreshold, TypingThreshold, TypingThreshold),
		makeThreshold(1.0, 1.0, 1.0))
}

type PadTyping struct {
	layout                    TypingLayout
	leftPadZone, rightPadZone Zone
	awaitingNeutralPos        bool
	leftCoords, rightCoords   *Coords
	leftCanUse, leftChanged   bool
	rightCanUse, rightChanged bool
}

func makePadTyping() PadTyping {
	return PadTyping{
		layout:             loadTypingLayout(),
		leftPadZone:        NeutralZone,
		rightPadZone:       NeutralZone,
		awaitingNeutralPos: false,
		leftCoords:         makeCoords(),
		rightCoords:        makeCoords(),
	}
}

var joystickTyping PadTyping

func zoneCanBeUsed(zone Zone) bool {
	return zone != UnmappedZone && zone != NeutralZone
}

func (padTyping *PadTyping) zoneChanged(zone Zone, prevZone *Zone) bool {
	if zone != UnmappedZone {
		if *prevZone != zone {
			*prevZone = zone
			if zone == NeutralZone {
				padTyping.awaitingNeutralPos = false
			}
			return true
		}
	}
	return false
}

func (padTyping *PadTyping) calcNewZone(prevZone *Zone, coords *Coords) (bool, bool) {
	coords.updateValues()
	coords.updateAngle()

	zone := detectZone(coords.magnitude, coords.angle, TypingBoundariesMap)
	//print("x: %0.2f; y: %0.2f; magn: %0.2f; angle: %v; zone: %s", coords.x, coords.y, coords.magnitude, coords.angle, zone)
	canUse := zoneCanBeUsed(zone)
	changed := padTyping.zoneChanged(zone, prevZone)
	return canUse, changed
}

func (padTyping *PadTyping) updateLeftZone() {
	//print("Left")
	padTyping.leftCanUse, padTyping.leftChanged = padTyping.calcNewZone(&padTyping.leftPadZone, padTyping.leftCoords)
	padTyping.typeLetter()
}
func (padTyping *PadTyping) updateRightZone() {
	//print("Right")
	padTyping.rightCanUse, padTyping.rightChanged = padTyping.calcNewZone(&padTyping.rightPadZone, padTyping.rightCoords)
	padTyping.typeLetter()
}

func (padTyping *PadTyping) typeLetter() {
	if padTyping.leftCanUse && padTyping.rightCanUse {
		//print("%s %s", padTyping.leftPadZone, padTyping.rightPadZone)
		//print("%v %v", leftCanUse, rightCanUse)
		//print("%v %v", leftChanged, rightChanged)

		if padTyping.leftChanged || padTyping.rightChanged {
			if !padTyping.awaitingNeutralPos {
				padTyping.awaitingNeutralPos = true
				position := SticksPosition{padTyping.leftPadZone, padTyping.rightPadZone}
				if code, found := padTyping.layout[position]; found {
					osSpec.TypeKey(code)
				}
			}
		}
	}
}
