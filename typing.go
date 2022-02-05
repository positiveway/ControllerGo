package main

import (
	"math"
	"os"
	"strings"
)

const neutralZone = "⬤"
const edgeZone = "❌"
const angleMargin = 15
const magnitudeThresholdPct = 75
const magnitudeThreshold = magnitudeThresholdPct / 100

type tuple2 = [2]string
type Layout = map[tuple2]string

func loadLayout() Layout {
	dat, err := os.ReadFile("layout.csv")
	check_err(err)
	lines := strings.Split(string(dat), "\n")
	lines = lines[1:]

	layout := Layout{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}
		parts := strings.Split(line, ", ")
		letter, leftStick, rightStick := parts[0], parts[1], parts[2]
		position := tuple2{leftStick, rightStick}
		if _, found := layout[position]; found {
			panic("duplicate position")
		}
		layout[position] = letter
	}
	return layout
}

type BoundariesMap = map[int]string

var boundariesMap = genBoundariesMap()

func genRange(lowerBound, upperBound int, _boundariesMap BoundariesMap, direction string) {
	for angle := lowerBound; angle < upperBound; angle++ {
		_boundariesMap[angle] = direction
	}
}

func genBoundariesMap() BoundariesMap {
	mapping := map[int]string{
		0:   "Right",
		45:  "UpRight",
		90:  "Up",
		135: "UpLeft",
		180: "Left",
		225: "DownLeft",
		270: "Down",
		315: "DownRight",
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
	layout                        Layout
	leftStickZone, rightStickZone string
	awaitingNeutralPos            bool
}

func makeJoystickTyping() JoystickTyping {
	return JoystickTyping{
		layout:             loadLayout(),
		leftStickZone:      neutralZone,
		rightStickZone:     neutralZone,
		awaitingNeutralPos: false,
	}
}

var joystickTyping = makeJoystickTyping()

func calcAngle(x, y float64) int {
	val := math.Atan2(y, x) * (180 / math.Pi)
	angle := int(val)
	if angle < 0 {
		angle = 360 + angle
	}
	return angle
}

func calcMagnitude(x, y float64) float64 {
	val := math.Pow(x, 2) + math.Pow(y, 2)
	norm := math.Sqrt(val / 2)
	return norm
}
