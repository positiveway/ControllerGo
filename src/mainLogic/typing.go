package mainLogic

import (
	"ControllerGo/src/platformSpecific"
	"fmt"
	"sort"
)

const NeutralZone ZoneT = "⬤"
const EdgeZone ZoneT = "❌"

const NoneStr = "None"

type SticksPosition [2]ZoneT
type TypingLayout map[SticksPosition]int
type AngleRange [2]int

func loadTypingLayout() TypingLayout {
	linesParts := ReadLayoutFile("typing.csv", 2)

	layout := TypingLayout{}
	for _, parts := range linesParts {
		leftStick, rightStick, letter := ZoneT(parts[0]), ZoneT(parts[1]), parts[2]
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

type BoundariesMap map[int]ZoneT

var boundariesMap BoundariesMap

func genRange(lowerBound, upperBound int, _boundariesMap BoundariesMap, direction ZoneT) {
	lowerBound += 360
	upperBound += 360

	for angle := lowerBound; angle <= upperBound; angle++ {
		resolvedAngle := resolveAngle(float64(angle))
		AssignWithDuplicateCheck(_boundariesMap, resolvedAngle, direction)
	}
}

func printValuesForDir(_boundariesMap BoundariesMap) {
	direction := ZoneRight
	var needAngles []int
	for angle, dir := range _boundariesMap {
		if dir == direction {
			needAngles = append(needAngles, angle)
		}
	}
	sort.Ints(needAngles)
	fmt.Println(needAngles)
}

func genBoundariesMap() BoundariesMap {
	//newMapping := map[string]AngleRange{
	//	ZoneRight:   {350, 22},
	//	ZoneUpRight: {24, 71},
	//}
	//print(newMapping)

	if RightAngleMargin+DiagonalAngleMargin > 45 {
		panic("With this margin of angle areas will overlap")
	}

	mapping := map[ZoneT]AngleRange{
		ZoneRight:     {0, RightAngleMargin},
		ZoneUpRight:   {45, DiagonalAngleMargin},
		ZoneUp:        {90, RightAngleMargin},
		ZoneUpLeft:    {135, DiagonalAngleMargin},
		ZoneLeft:      {180, RightAngleMargin},
		ZoneDownLeft:  {225, DiagonalAngleMargin},
		ZoneDown:      {270, RightAngleMargin},
		ZoneDownRight: {315, DiagonalAngleMargin},
	}

	_boundariesMap := BoundariesMap{}
	for direction, angleRange := range mapping {
		angle, angleMargin := angleRange[0], angleRange[1]
		genRange(angle-angleMargin, angle+angleMargin, _boundariesMap, direction)
	}
	//printValuesForDir(_boundariesMap)
	return _boundariesMap
}

type PadTyping struct {
	layout                    TypingLayout
	leftPadZone, rightPadZone ZoneT
	awaitingNeutralPos        bool
	leftCoords, rightCoords   Coords
	leftCanUse, leftChanged   bool
	rightCanUse, rightChanged bool
}

func makePadTyping() PadTyping {
	return PadTyping{
		layout:             loadTypingLayout(),
		leftPadZone:        NeutralZone,
		rightPadZone:       NeutralZone,
		awaitingNeutralPos: false,
	}
}

var joystickTyping PadTyping

func detectZone(magnitude float64, angle int) ZoneT {
	if magnitude > MagnitudeThreshold {
		//print("%v", angle)
		return getOrDefault(boundariesMap, angle, EdgeZone)
	} else {
		return NeutralZone
	}
}

func zoneCanBeUsed(zone ZoneT) bool {
	return zone != EdgeZone && zone != NeutralZone
}

func (padTyping *PadTyping) zoneChanged(zone ZoneT, prevZone *ZoneT) bool {
	if zone != EdgeZone {
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

func (padTyping *PadTyping) calcNewZone(prevZone *ZoneT, coords *Coords) (bool, bool) {
	coords.updateValues()
	coords.updateAngle()

	zone := detectZone(coords.magnitude, coords.angle)
	//print("x: %0.2f; y: %0.2f; magn: %0.2f; angle: %v; zone: %s", coords.x, coords.y, coords.magnitude, coords.angle, zone)
	canUse := zoneCanBeUsed(zone)
	changed := padTyping.zoneChanged(zone, prevZone)
	return canUse, changed
}

func (padTyping *PadTyping) updateLeftZone() {
	//print("Left")
	padTyping.leftCanUse, padTyping.leftChanged = padTyping.calcNewZone(&padTyping.leftPadZone, &padTyping.leftCoords)
	padTyping.typeLetter()
}
func (padTyping *PadTyping) updateRightZone() {
	//print("Right")
	padTyping.rightCanUse, padTyping.rightChanged = padTyping.calcNewZone(&padTyping.rightPadZone, &padTyping.rightCoords)
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
					platformSpecific.TypeKey(code)
				}
			}
		}
	}
}
