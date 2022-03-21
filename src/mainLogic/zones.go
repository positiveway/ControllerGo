package mainLogic

import (
	"fmt"
	"math"
	"sort"
)

const (
	AngleRight     int = 0
	AngleUpRight   int = 45
	AngleUp        int = 90
	AngleUpLeft    int = 135
	AngleLeft      int = 180
	AngleDownLeft  int = 225
	AngleDown      int = 270
	AngleDownRight int = 315
)

type Zone string

const (
	ZoneRight     Zone = "Right"
	ZoneUpRight   Zone = "UpRight"
	ZoneUp        Zone = "Up"
	ZoneUpLeft    Zone = "UpLeft"
	ZoneLeft      Zone = "Left"
	ZoneDownLeft  Zone = "DownLeft"
	ZoneDown      Zone = "Down"
	ZoneDownRight Zone = "DownRight"
)

var AllZones = []Zone{
	ZoneRight,
	ZoneUpRight,
	ZoneUp,
	ZoneUpLeft,
	ZoneLeft,
	ZoneDownLeft,
	ZoneDown,
	ZoneDownRight,
}

type Direction struct {
	threshold float64
	zone      Zone
}

type BoundariesMap map[int]Direction

type InitBoundaries map[int]Zone

type Threshold struct {
	diagonal, horizontal, vertical float64
}

func makeThreshold(diagonal, horizontal, vertical float64) Threshold {
	return Threshold{diagonal: diagonal, horizontal: horizontal, vertical: vertical}
}

type AngleMargin struct {
	diagonal, horizontal, vertical int
}

func makeAngleMargin(diagonal, horizontal, vertical int) AngleMargin {
	return AngleMargin{diagonal: diagonal, horizontal: horizontal, vertical: vertical}
}

func genRange(lowerBound, upperBound int, _boundariesMap BoundariesMap, zone Zone, threshold float64) {
	lowerBound += 360
	upperBound += 360

	for angle := lowerBound; angle <= upperBound; angle++ {
		resolvedAngle := resolveAngle(angle)
		AssignWithDuplicateCheck(_boundariesMap, resolvedAngle, Direction{zone: zone, threshold: threshold})
	}
}

func isDiagonal(angle int) bool {
	return math.Mod(float64(angle), 90) == 45
}

func isHorizontal(angle int) bool {
	return angle == 0 || angle == 180
}

func isVertical(angle int) bool {
	return angle == 90 || angle == 270
}

func resolveThreeDirectionalValue[T Number](angle int, diagonal, horizontal, vertical T) T {
	if isDiagonal(angle) {
		return diagonal
	}
	if isHorizontal(angle) {
		return horizontal
	}
	if isVertical(angle) {
		return vertical
	}
	panicMsg("Incorrect base angle")
	panic("")
}

func getThresholdValue(angle int, threshold Threshold) float64 {
	return resolveThreeDirectionalValue(angle, threshold.diagonal, threshold.horizontal, threshold.vertical)
}

func getAngleMargin(angle int, angleMargin AngleMargin) int {
	return resolveThreeDirectionalValue(angle, angleMargin.diagonal, angleMargin.horizontal, angleMargin.vertical)
}

func checkThreshold(threshold Threshold) {
	if threshold.diagonal == 0 || threshold.vertical == 0 || threshold.horizontal == 0 {
		panicMsg("Threshold can't be zero")
	}
}

func checkAngleMargin(angleMargin AngleMargin) {
	if angleMargin.diagonal+angleMargin.horizontal > 45 ||
		angleMargin.diagonal+angleMargin.vertical > 45 {
		panicMsg("With this margin of angle areas will overlap")
	}
}

func genBoundariesMap(initBoundaries InitBoundaries, angleMargin AngleMargin, threshold Threshold) BoundariesMap {
	//newMapping := map[string]AngleRange{
	//	ZoneRight:   {350, 22},
	//	ZoneUpRight: {24, 71},
	//}
	//print(newMapping)

	checkThreshold(threshold)
	checkAngleMargin(angleMargin)

	_boundariesMap := BoundariesMap{}
	for baseAngle, direction := range initBoundaries {
		margin := getAngleMargin(baseAngle, angleMargin)
		thresholdValue := getThresholdValue(baseAngle, threshold)
		genRange(baseAngle-margin, baseAngle+margin, _boundariesMap, direction, thresholdValue)
	}
	//printAnglesForZones(_boundariesMap)
	return _boundariesMap
}

func printAnglesForZones(_boundariesMap BoundariesMap) {
	for _, zone := range AllZones {
		print("%v: ", zone)
		var needAngles []int
		for angle, dir := range _boundariesMap {
			if dir.zone == zone {
				needAngles = append(needAngles, angle)
			}
		}
		sort.Ints(needAngles)
		fmt.Println(needAngles)
	}
}

func detectZone(magnitude float64, angle int, boundariesMap BoundariesMap) Zone {
	if direction, found := boundariesMap[angle]; found {
		if magnitude > direction.threshold {
			return direction.zone
		} else {
			return NeutralZone
		}
	} else {
		return EdgeZone
	}
}
