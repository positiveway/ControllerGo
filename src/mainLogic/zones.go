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

type ZoneT string

const (
	ZoneRight     ZoneT = "Right"
	ZoneUpRight   ZoneT = "UpRight"
	ZoneUp        ZoneT = "Up"
	ZoneUpLeft    ZoneT = "UpLeft"
	ZoneLeft      ZoneT = "Left"
	ZoneDownLeft  ZoneT = "DownLeft"
	ZoneDown      ZoneT = "Down"
	ZoneDownRight ZoneT = "DownRight"
)

var AllZones = []ZoneT{
	ZoneRight,
	ZoneUpRight,
	ZoneUp,
	ZoneUpLeft,
	ZoneLeft,
	ZoneDownLeft,
	ZoneDown,
	ZoneDownRight,
}

const CentralNeutralZone ZoneT = "⬤"
const UnmappedZone ZoneT = "❌"
const EdgeZoneSuffix ZoneT = "_Edge"

type DirectionT struct {
	zoneThresholdPct, edgeThresholdPct float64
	zone                               ZoneT
}

func makeDirection(zone ZoneT, zoneThresholdPct, edgeThresholdPct float64) DirectionT {
	return DirectionT{zone: zone, zoneThresholdPct: zoneThresholdPct, edgeThresholdPct: edgeThresholdPct}
}

type ZoneBoundariesMapT map[int]DirectionT

type InitBoundariesT map[int]ZoneT

type ThresholdT struct {
	diagonal, horizontal, vertical float64
}

func makeThreshold(diagonal, horizontal, vertical float64) ThresholdT {
	return ThresholdT{diagonal: diagonal, horizontal: horizontal, vertical: vertical}
}

type AngleMarginT struct {
	diagonal, horizontal, vertical int
}

func makeAngleMargin(diagonal, horizontal, vertical int) AngleMarginT {
	return AngleMarginT{diagonal: diagonal, horizontal: horizontal, vertical: vertical}
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

func isEdgeZone(zone ZoneT) bool {
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

func getThresholdValue(angle int, threshold ThresholdT) float64 {
	return resolveThreeDirectionalValue(angle, threshold.diagonal, threshold.horizontal, threshold.vertical)
}

func getAngleMargin(angle int, angleMargin AngleMarginT) int {
	return resolveThreeDirectionalValue(angle, angleMargin.diagonal, angleMargin.horizontal, angleMargin.vertical)
}

func checkZoneThreshold(zoneThreshold ThresholdT) {
	if gofuncs.AnyEqual([][]float64{
		{zoneThreshold.diagonal, 0},
		{zoneThreshold.vertical, 0},
		{zoneThreshold.horizontal, 0},
	}) {
		gofuncs.Panic("Threshold can't be zero")
	}
}

func checkEdgeThreshold(zoneThreshold, edgeThreshold ThresholdT) {
	if gofuncs.AnyGreaterOrEqual([][]float64{
		{zoneThreshold.diagonal, edgeThreshold.diagonal},
		{zoneThreshold.vertical, edgeThreshold.vertical},
		{zoneThreshold.horizontal, edgeThreshold.horizontal},
	}) {
		gofuncs.Panic("Edge threshold can't be less or equal to Zone threshold")
	}
}

func checkAngleMargin(angleMargin AngleMarginT) {
	if gofuncs.AnyGreaterOrEqual([][]int{
		{angleMargin.diagonal + angleMargin.horizontal, 45},
		{angleMargin.diagonal + angleMargin.vertical, 45},
	}) {
		gofuncs.Panic("With this margin of angle areas will overlap")
	}
}

func genRange(lowerBound, upperBound int, _boundariesMap ZoneBoundariesMapT, zone ZoneT, zoneThreshold, edgeThreshold float64) {
	lowerBound += 360
	upperBound += 360

	for angle := lowerBound; angle <= upperBound; angle++ {
		resolvedAngle := resolveCircleAngle(angle)
		gofuncs.AssignWithDuplicateCheck(_boundariesMap, resolvedAngle, makeDirection(zone, zoneThreshold, edgeThreshold))
	}
}

func genInitBoundaries(includeDiagonalZones bool) InitBoundariesT {
	initBoundaries := InitBoundariesT{
		AngleRight:     ZoneRight,
		AngleUpRight:   ZoneUpRight,
		AngleUp:        ZoneUp,
		AngleUpLeft:    ZoneUpLeft,
		AngleLeft:      ZoneLeft,
		AngleDownLeft:  ZoneDownLeft,
		AngleDown:      ZoneDown,
		AngleDownRight: ZoneDownRight,
	}
	resBoundaries := InitBoundariesT{}
	for angle, zone := range initBoundaries {
		if isDiagonal(angle) && !includeDiagonalZones {
			continue
		}
		resBoundaries[angle] = zone
	}
	return resBoundaries
}

func genEqualThresholdBoundariesMap(includeDiagonalZones bool, angleMargin AngleMarginT, zoneThreshold, edgeThreshold float64) ZoneBoundariesMapT {
	return genBoundariesMap(includeDiagonalZones, angleMargin,
		makeThreshold(zoneThreshold, zoneThreshold, zoneThreshold),
		makeThreshold(edgeThreshold, edgeThreshold, edgeThreshold))
}

func genBoundariesMap(includeDiagonalZones bool, angleMargin AngleMarginT, zoneThreshold, edgeThreshold ThresholdT) ZoneBoundariesMapT {
	checkZoneThreshold(zoneThreshold)
	checkEdgeThreshold(zoneThreshold, edgeThreshold)
	checkAngleMargin(angleMargin)

	_boundariesMap := ZoneBoundariesMapT{}
	for baseAngle, direction := range genInitBoundaries(includeDiagonalZones) {
		margin := getAngleMargin(baseAngle, angleMargin)
		zoneThresholdValue := getThresholdValue(baseAngle, zoneThreshold)
		edgeThresholdValue := getThresholdValue(baseAngle, edgeThreshold)
		genRange(baseAngle-margin, baseAngle+margin, _boundariesMap, direction, zoneThresholdValue, edgeThresholdValue)
	}
	//printAnglesForZones(_boundariesMap)
	return _boundariesMap
}

func printAnglesForZones(_boundariesMap ZoneBoundariesMapT) {
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

func detectZone(magnitude, radius float64, angle int, boundariesMap ZoneBoundariesMapT) ZoneT {
	if direction, found := boundariesMap[angle]; found {
		if magnitude > direction.zoneThresholdPct*radius {
			zone := direction.zone

			if magnitude > direction.edgeThresholdPct*radius {
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

func (pad *PadStickPositionT) ReCalculateZone(zoneBoundariesMap ZoneBoundariesMapT) {
	if pad.newValueHandled {
		return
	}
	pad.newValueHandled = true

	zone := detectZone(pad.magnitude, pad.radius, pad.shiftedAngle, zoneBoundariesMap)
	//printDebug("(x: %0.2f, y: %0.2f); magn: %0.2f; angle: %v; zone: %s", pad.x, pad.y, pad.magnitude, pad.shiftedAngle, zone)

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
