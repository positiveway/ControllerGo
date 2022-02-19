package mainLogic

import (
	"ControllerGo/src/osSpecific"
	"math"
)

const NeutralZone = "⬤"
const EdgeZone = "❌"

const NoneStr = "None"

type SticksPosition = [2]string
type TypingLayout = map[SticksPosition]int

func loadTypingLayout() TypingLayout {
	linesParts := ReadLayoutFile("typing_layout.csv")

	layout := TypingLayout{}
	for _, parts := range linesParts {
		leftStick, rightStick, letter := parts[0], parts[1], parts[2]
		if !contains(AllZones, leftStick) {
			panicMisspelled(leftStick)
		}
		if !contains(AllZones, rightStick) {
			panicMisspelled(rightStick)
		}
		if letter == NoneStr {
			continue
		}
		code := getCodeFromLetter(letter)
		position := SticksPosition{leftStick, rightStick}
		assignWithDuplicateCheck(layout, position, code)
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
	mapping := map[int]string{
		0:   ZoneRight,
		45:  ZoneUpRight,
		90:  ZoneUp,
		135: ZoneUpLeft,
		180: ZoneLeft,
		225: ZoneDownLeft,
		270: ZoneDown,
		315: ZoneDownRight,
	}
	if angleMargin > 22 {
		panic("With this margin of angle areas will overlap")
	}
	_boundariesMap := BoundariesMap{}
	for angle, dir := range mapping {
		genRange(angle, angle+angleMargin, _boundariesMap, dir)
		if angle == 0 {
			genRange(360-angleMargin, 360, _boundariesMap, dir)
		} else {
			genRange(angle-angleMargin, angle, _boundariesMap, dir)
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
	magnitude := calcMagnitude(x, y)
	angle := calcAngle(x, y)
	//fmt.Printf("(%.2f, %.2f): %v %.2f", x, y, angle, magnitude)

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
