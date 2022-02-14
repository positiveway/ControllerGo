package main

import (
	"math"
)

const NeutralZone = "⬤"
const EdgeZone = "❌"
const UndefinedMapping = "Undefined"
const angleMargin int = 7
const magnitudeThresholdPct float64 = 40
const MagnitudeThreshold float64 = magnitudeThresholdPct / 100

const EmptyStr = ""
const NoneStr = "None"

type SticksPosition = [2]string
type TypingLayout = map[SticksPosition]string

func loadTypingLayout() TypingLayout {
	linesParts := readLayoutFile("typing_layout.csv")

	layout := TypingLayout{}
	for _, parts := range linesParts {
		leftStick, rightStick, letter := parts[0], parts[1], parts[2]
		if !contains(AllZones, leftStick) {
			panicMisspelled(leftStick)
		}
		if !contains(AllZones, rightStick) {
			panicMisspelled(rightStick)
		}
		position := SticksPosition{leftStick, rightStick}
		assignWithDuplicateCheck(layout, position, letter)
	}
	return layout
}

type BoundariesMap = map[int]string

var boundariesMap BoundariesMap

func genRange(lowerBound, upperBound int, _boundariesMap BoundariesMap, direction string) {
	for angle := lowerBound; angle < upperBound; angle++ {
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

func (jTyping *JoystickTyping) detectLetter() string {
	curZones := SticksPosition{jTyping.leftStickZone, jTyping.rightStickZone}
	for _, zone := range curZones {
		if zone == NeutralZone {
			return EmptyStr
		} else if zone == EdgeZone {
			panic("zone to letter error")
		}
	}
	jTyping.awaitingNeutralPos = true
	letter := getOrDefault(jTyping.layout, curZones, UndefinedMapping)
	return letter
}

func (jTyping *JoystickTyping) _updateZone(prevZone *string, coords *Coords) string {
	x, y := coords.getValues()
	magnitude := calcMagnitude(x, y)
	angle := calcAngle(x, y)
	//fmt.Printf("(%.2f, %.2f): %v %.2f", x, y, angle, magnitude)

	newZone := detectZone(magnitude, angle)
	if newZone == EdgeZone {
		return EmptyStr
	}
	if newZone != *prevZone {
		*prevZone = newZone
		//return jTyping.detectLetter()
		if jTyping.awaitingNeutralPos {
			if newZone == NeutralZone {
				jTyping.awaitingNeutralPos = false
			}
		} else {
			return jTyping.detectLetter()
		}
	}
	return EmptyStr
}

func typeLetter(letter string) {
	if letter != UndefinedMapping && letter != NoneStr {
		key := getCodeFromLetter(letter)
		keyboard.KeyPress(key)
	} else {
		//fmt.Println(UndefinedMapping)
	}
}

func (jTyping *JoystickTyping) updateZone(prevZone *string, coords *Coords) {
	letter := jTyping._updateZone(prevZone, coords)
	//fmt.Printf(" %s %s %v\n", jTyping.leftStickZone, jTyping.rightStickZone, jTyping.awaitingNeutralPos)
	if letter != EmptyStr {
		typeLetter(letter)
	}
}

func (jTyping *JoystickTyping) updateZoneLeft() {
	coords := &jTyping.leftCoords
	prevZone := &jTyping.leftStickZone

	jTyping.updateZone(prevZone, coords)
}

func (jTyping *JoystickTyping) updateZoneRight() {
	coords := &jTyping.rightCoords
	prevZone := &jTyping.rightStickZone

	jTyping.updateZone(prevZone, coords)
}