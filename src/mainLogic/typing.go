package mainLogic

import (
	"ControllerGo/src/osSpecific"
	"fmt"
	"math"
)

const NeutralZone = "⬤"
const EdgeZone = "❌"

const NoneStr = "None"

type SticksPosition = [2]string
type TypingLayout = map[SticksPosition]int
type AngleRange = [2]int

func loadTypingLayout() TypingLayout {
	linesParts := ReadLayoutFile("typing.csv", 2)

	layout := TypingLayout{}
	for _, parts := range linesParts {
		leftStick, rightStick, letter := parts[0], parts[1], parts[2]
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

type BoundariesMap = map[int]string

var boundariesMap BoundariesMap

func genRange(lowerBound, upperBound int, _boundariesMap BoundariesMap, direction string) {
	for angle := lowerBound; angle <= upperBound; angle++ {
		_boundariesMap[angle] = direction
	}
}

func genBoundariesMap() BoundariesMap {
	newMapping := map[string]AngleRange{
		ZoneRight:   {350, 22},
		ZoneUpRight: {24, 71},
	}
	fmt.Println(newMapping)
	// %360

	if RightAngleMargin+DiagonalAngleMargin > 45 {
		panic("With this margin of angle areas will overlap")
	}

	mapping := map[string]AngleRange{
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

		genRange(angle, angle+angleMargin, _boundariesMap, direction)
		if angle == 0 {
			genRange(360-angleMargin, 360, _boundariesMap, direction)
		} else {
			genRange(angle-angleMargin, angle, _boundariesMap, direction)
		}
	}
	return _boundariesMap
}

type JoystickTyping struct {
	layout                        TypingLayout
	leftStickZone, rightStickZone string
	awaitingNeutralPos            bool
	leftCoords, rightCoords       Coords
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

func calcAngle(x, y float64) int {
	val := math.Atan2(y, x) * (180 / math.Pi)
	angle := int(val)
	if angle < 0 {
		angle = 360 + angle
	}
	return angle
}

func calcMagnitude(x, y float64) float64 {
	magnitude := math.Sqrt(x*x + y*y)
	if magnitude > 1.0 {
		magnitude = 1.0
	}
	return magnitude
}

func detectZone(magnitude float64, angle int) string {
	if magnitude > MagnitudeThreshold {
		//fmt.Printf("%v\n", angle)
		return getOrDefault(boundariesMap, angle, EdgeZone)
	} else {
		return NeutralZone
	}
}

func zoneCanBeUsed(zone string) bool {
	return zone != EdgeZone && zone != NeutralZone
}

func (jTyping *JoystickTyping) zoneChanged(zone string, prevZone *string) bool {
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

func (jTyping *JoystickTyping) calcNewZone(prevZone *string, coords *Coords) (bool, bool) {
	x, y := coords.getValues()
	angle := calcAngle(x, y)
	magnitude := calcMagnitude(x, y)

	zone := detectZone(magnitude, angle)
	canUse := zoneCanBeUsed(zone)
	changed := jTyping.zoneChanged(zone, prevZone)
	return canUse, changed
}

func (jTyping *JoystickTyping) updateZones() {
	leftCanUse, leftChanged := jTyping.calcNewZone(&jTyping.leftStickZone, &jTyping.leftCoords)
	rightCanUse, rightChanged := jTyping.calcNewZone(&jTyping.rightStickZone, &jTyping.rightCoords)

	if leftCanUse && rightCanUse {
		//fmt.Printf("%s %s\n", jTyping.leftStickZone, jTyping.rightStickZone)
		//fmt.Printf("%v %v\n", leftCanUse, rightCanUse)
		//fmt.Printf("%v %v\n", leftChanged, rightChanged)

		if leftChanged || rightChanged {
			if !jTyping.awaitingNeutralPos {
				jTyping.awaitingNeutralPos = true
				position := SticksPosition{jTyping.leftStickZone, jTyping.rightStickZone}
				if code, found := jTyping.layout[position]; found {
					osSpecific.TypeKey(code)
				}
			}
		}
	}
}
