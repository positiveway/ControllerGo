package mainLogic

import (
	"fmt"
	"github.com/positiveway/gofuncs"
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

const CentralNeutralZone Zone = "⬤"
const UnmappedZone Zone = "❌"
const EdgeZoneSuffix Zone = "_Edge"

type Direction struct {
	zoneThreshold, edgeThreshold float64
	zone                         Zone
}

func makeDirection(zone Zone, zoneThreshold, edgeThreshold float64) Direction {
	return Direction{zone: zone, zoneThreshold: zoneThreshold, edgeThreshold: edgeThreshold}
}

type ZoneBoundariesMap map[int]Direction

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

func isEdgeZone(zone Zone) bool {
	return gofuncs.EndsWith(string(zone), string(EdgeZoneSuffix))
}

func resolveThreeDirectionalValue[T gofuncs.Number](angle int, diagonal, horizontal, vertical T) T {
	if isDiagonal(angle) {
		return diagonal
	}
	if isHorizontal(angle) {
		return horizontal
	}
	if isVertical(angle) {
		return vertical
	}
	gofuncs.Panic("Incorrect base angle")
	panic("")
}

func getThresholdValue(angle int, threshold Threshold) float64 {
	return resolveThreeDirectionalValue(angle, threshold.diagonal, threshold.horizontal, threshold.vertical)
}

func getAngleMargin(angle int, angleMargin AngleMargin) int {
	return resolveThreeDirectionalValue(angle, angleMargin.diagonal, angleMargin.horizontal, angleMargin.vertical)
}

func checkZoneThreshold(zoneThreshold Threshold) {
	if gofuncs.AnyEqual([][]float64{
		{zoneThreshold.diagonal, 0},
		{zoneThreshold.vertical, 0},
		{zoneThreshold.horizontal, 0},
	}) {
		gofuncs.Panic("Threshold can't be zero")
	}
}

func checkEdgeThreshold(zoneThreshold, edgeThreshold Threshold) {
	if gofuncs.AnyGreaterOrEqual([][]float64{
		{zoneThreshold.diagonal, edgeThreshold.diagonal},
		{zoneThreshold.vertical, edgeThreshold.vertical},
		{zoneThreshold.horizontal, edgeThreshold.horizontal},
	}) {
		gofuncs.Panic("Edge threshold can't be less or equal to Zone threshold")
	}
}

func checkAngleMargin(angleMargin AngleMargin) {
	if gofuncs.AnyGreaterOrEqual([][]int{
		{angleMargin.diagonal + angleMargin.horizontal, 45},
		{angleMargin.diagonal + angleMargin.vertical, 45},
	}) {
		gofuncs.Panic("With this margin of angle areas will overlap")
	}
}

func genRange(lowerBound, upperBound int, _boundariesMap ZoneBoundariesMap, zone Zone, zoneThreshold, edgeThreshold float64) {
	lowerBound += 360
	upperBound += 360

	for angle := lowerBound; angle <= upperBound; angle++ {
		resolvedAngle := resolveAngle(angle)
		gofuncs.AssignWithDuplicateCheck(_boundariesMap, resolvedAngle, makeDirection(zone, zoneThreshold, edgeThreshold))
	}
}

func genInitBoundaries(includeDiagonalZones bool) InitBoundaries {
	initBoundaries := InitBoundaries{
		AngleRight:     ZoneRight,
		AngleUpRight:   ZoneUpRight,
		AngleUp:        ZoneUp,
		AngleUpLeft:    ZoneUpLeft,
		AngleLeft:      ZoneLeft,
		AngleDownLeft:  ZoneDownLeft,
		AngleDown:      ZoneDown,
		AngleDownRight: ZoneDownRight,
	}
	resBoundaries := InitBoundaries{}
	for angle, zone := range initBoundaries {
		if isDiagonal(angle) && !includeDiagonalZones {
			continue
		}
		resBoundaries[angle] = zone
	}
	return resBoundaries
}

func genEqualThresholdBoundariesMap(includeDiagonalZones bool, angleMargin AngleMargin, zoneThreshold, edgeThreshold float64) ZoneBoundariesMap {
	return genBoundariesMap(includeDiagonalZones, angleMargin,
		makeThreshold(zoneThreshold, zoneThreshold, zoneThreshold),
		makeThreshold(edgeThreshold, edgeThreshold, edgeThreshold))
}

func genBoundariesMap(includeDiagonalZones bool, angleMargin AngleMargin, zoneThreshold, edgeThreshold Threshold) ZoneBoundariesMap {
	checkZoneThreshold(zoneThreshold)
	checkEdgeThreshold(zoneThreshold, edgeThreshold)
	checkAngleMargin(angleMargin)

	_boundariesMap := ZoneBoundariesMap{}
	for baseAngle, direction := range genInitBoundaries(includeDiagonalZones) {
		margin := getAngleMargin(baseAngle, angleMargin)
		zoneThresholdValue := getThresholdValue(baseAngle, zoneThreshold)
		edgeThresholdValue := getThresholdValue(baseAngle, edgeThreshold)
		genRange(baseAngle-margin, baseAngle+margin, _boundariesMap, direction, zoneThresholdValue, edgeThresholdValue)
	}
	//printAnglesForZones(_boundariesMap)
	return _boundariesMap
}

func printAnglesForZones(_boundariesMap ZoneBoundariesMap) {
	for _, zone := range AllZones {
		gofuncs.Print("%v: ", zone)
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

func detectZone(magnitude float64, angle int, zoneRotation int, boundariesMap ZoneBoundariesMap) Zone {
	angle += zoneRotation
	if direction, found := boundariesMap[angle]; found {
		if magnitude > direction.zoneThreshold {
			zone := direction.zone

			if magnitude > direction.edgeThreshold {
				zone += EdgeZoneSuffix
			}
			return zone

		} else {
			return CentralNeutralZone
		}
	} else {
		return UnmappedZone
	}
}

func (pad *PadPosition) ReCalculateZone(zoneBoundariesMap ZoneBoundariesMap) {
	if pad.newValueHandled {
		return
	}
	pad.newValueHandled = true

	zone := detectZone(pad.magnitude, pad.angle, pad.zoneRotation, zoneBoundariesMap)
	//printDebug("x: %0.2f; y: %0.2f; magn: %0.2f; angle: %v; zone: %s", pad.x, pad.y, pad.magnitude, pad.angle, zone)

	if zone == UnmappedZone {
		pad.zoneCanBeUsed = false
		pad.zoneChanged = false
		return
	}

	pad.zoneCanBeUsed = zone != CentralNeutralZone
	pad.zoneChanged = pad.zone != zone
	pad.zone = zone

	if pad.zoneChanged && pad.zone == CentralNeutralZone {
		pad.awaitingCentralPosition = false
	}
}
