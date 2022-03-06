package mainLogic

import "ControllerGo/src/platformSpecific"

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
		_boundariesMap[resolvedAngle] = direction
	}
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
	return _boundariesMap
}

type JoystickTyping struct {
	layout                        TypingLayout
	leftStickZone, rightStickZone ZoneT
	awaitingNeutralPos            bool
	leftCoords, rightCoords       Coords
	leftCanUse, leftChanged       bool
	rightCanUse, rightChanged     bool
}

func makeJoystickTyping() JoystickTyping {
	return JoystickTyping{
		layout:             loadTypingLayout(),
		leftStickZone:      NeutralZone,
		rightStickZone:     NeutralZone,
		awaitingNeutralPos: false,
	}
}

var joystickTyping JoystickTyping

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

func (jTyping *JoystickTyping) zoneChanged(zone ZoneT, prevZone *ZoneT) bool {
	if zone != EdgeZone {
		if *prevZone != zone {
			*prevZone = zone
			if zone == NeutralZone {
				jTyping.awaitingNeutralPos = false
			}
			return true
		}
	}
	return false
}

func (jTyping *JoystickTyping) calcNewZone(prevZone *ZoneT, coords *Coords) (bool, bool) {
	coords.updateValues()
	coords.updateAngle()

	zone := detectZone(coords.magnitude, coords.angle)
	canUse := zoneCanBeUsed(zone)
	changed := jTyping.zoneChanged(zone, prevZone)
	return canUse, changed
}

func (jTyping *JoystickTyping) updateLeftZone() {
	jTyping.leftCanUse, jTyping.leftChanged = jTyping.calcNewZone(&jTyping.leftStickZone, &jTyping.leftCoords)
	jTyping.typeLetter()
}
func (jTyping *JoystickTyping) updateRightZone() {
	jTyping.rightCanUse, jTyping.rightChanged = jTyping.calcNewZone(&jTyping.rightStickZone, &jTyping.rightCoords)
	jTyping.typeLetter()
}

func (jTyping *JoystickTyping) typeLetter() {
	if jTyping.leftCanUse && jTyping.rightCanUse {
		//print("%s %s", jTyping.leftStickZone, jTyping.rightStickZone)
		//print("%v %v", leftCanUse, rightCanUse)
		//print("%v %v", leftChanged, rightChanged)

		if jTyping.leftChanged || jTyping.rightChanged {
			if !jTyping.awaitingNeutralPos {
				jTyping.awaitingNeutralPos = true
				position := SticksPosition{jTyping.leftStickZone, jTyping.rightStickZone}
				if code, found := jTyping.layout[position]; found {
					platformSpecific.TypeKey(code)
				}
			}
		}
	}
}
