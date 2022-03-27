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

const NeutralZone Zone = "⬤"
const UnmappedZone Zone = "❌"
const EdgeZoneSuffix Zone = "_Edge"

type Direction struct {
	zoneThreshold, edgeThreshold float64
	zone                         Zone
}

func makeDirection(zone Zone, zoneThreshold, edgeThreshold float64) Direction {
	return Direction{zone: zone, zoneThreshold: zoneThreshold, edgeThreshold: edgeThreshold}
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

func checkZoneThreshold(zoneThreshold Threshold) {
	if anyEqual([][]float64{
		{zoneThreshold.diagonal, 0},
		{zoneThreshold.vertical, 0},
		{zoneThreshold.horizontal, 0},
	}) {
		panicMsg("Threshold can't be zero")
	}
}

func checkEdgeThreshold(zoneThreshold, edgeThreshold Threshold) {
	if anyGreaterOrEqual([][]float64{
		{zoneThreshold.diagonal, edgeThreshold.diagonal},
		{zoneThreshold.vertical, edgeThreshold.vertical},
		{zoneThreshold.horizontal, edgeThreshold.horizontal},
	}) {
		panicMsg("Incorrect edge threshold")
	}
}

func checkAngleMargin(angleMargin AngleMargin) {
	if anyGreaterOrEqual([][]int{
		{angleMargin.diagonal + angleMargin.horizontal, 45},
		{angleMargin.diagonal + angleMargin.vertical, 45},
	}) {
		panicMsg("With this margin of angle areas will overlap")
	}
}

func genRange(lowerBound, upperBound int, _boundariesMap BoundariesMap, zone Zone, zoneThreshold, edgeThreshold float64) {
	lowerBound += 360
	upperBound += 360

	for angle := lowerBound; angle <= upperBound; angle++ {
		resolvedAngle := resolveAngle(angle)
		AssignWithDuplicateCheck(_boundariesMap, resolvedAngle, makeDirection(zone, zoneThreshold, edgeThreshold))
	}
}

func genBoundariesMap(initBoundaries InitBoundaries, angleMargin AngleMargin, zoneThreshold, edgeThreshold Threshold) BoundariesMap {
	//newMapping := map[string]AngleRange{
	//	ZoneRight:   {350, 22},
	//	ZoneUpRight: {24, 71},
	//}
	//print(newMapping)

	checkZoneThreshold(zoneThreshold)
	checkEdgeThreshold(zoneThreshold, edgeThreshold)
	checkAngleMargin(angleMargin)

	_boundariesMap := BoundariesMap{}
	for baseAngle, direction := range initBoundaries {
		margin := getAngleMargin(baseAngle, angleMargin)
		zoneThresholdValue := getThresholdValue(baseAngle, zoneThreshold)
		edgeThresholdValue := getThresholdValue(baseAngle, edgeThreshold)
		genRange(baseAngle-margin, baseAngle+margin, _boundariesMap, direction, zoneThresholdValue, edgeThresholdValue)
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
		if magnitude > direction.zoneThreshold {
			zone := direction.zone

			if magnitude > direction.edgeThreshold {
				zone += EdgeZoneSuffix
			}
			return zone

		} else {
			return NeutralZone
		}
	} else {
		return UnmappedZone
	}
}
