package mainLogic

import (
	"fmt"
	"github.com/positiveway/gofuncs"
	"math"
	"sort"
)

const (
	AngleRight     uint = 0
	AngleUpRight   uint = 45
	AngleUp        uint = 90
	AngleUpLeft    uint = 135
	AngleLeft      uint = 180
	AngleDownLeft  uint = 225
	AngleDown      uint = 270
	AngleDownRight uint = 315
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

type ZonesT []ZoneT

func InitAllZones() ZonesT {
	return ZonesT{
		ZoneRight,
		ZoneUpRight,
		ZoneUp,
		ZoneUpLeft,
		ZoneLeft,
		ZoneDownLeft,
		ZoneDown,
		ZoneDownRight,
	}
}

const CentralNeutralZone ZoneT = "⬤"
const UnmappedZone ZoneT = "❌"
const EdgeZoneSuffix ZoneT = "_Edge"

type DirectionT struct {
	zoneThresholdPct, edgeThresholdPct float64
	zone                               ZoneT
}

func MakeDirection(zone ZoneT, zoneThresholdPct, edgeThresholdPct float64) *DirectionT {
	return &DirectionT{zone: zone, zoneThresholdPct: zoneThresholdPct, edgeThresholdPct: edgeThresholdPct}
}

type ZoneBoundariesMapT map[uint]*DirectionT

type InitBoundariesT map[uint]ZoneT

type ThresholdPctT struct {
	diagonal, horizontal, vertical float64
}

func MakeThresholdPct(diagonal, horizontal, vertical float64) *ThresholdPctT {
	return &ThresholdPctT{diagonal: diagonal, horizontal: horizontal, vertical: vertical}
}

type AngleMarginT struct {
	Diagonal   uint `json:"Diagonal"`
	Horizontal uint `json:"Horizontal"`
	Vertical   uint `json:"Vertical"`
}

func MakeAngleMargin(diagonal, horizontal, vertical uint) *AngleMarginT {
	return &AngleMarginT{Diagonal: diagonal, Horizontal: horizontal, Vertical: vertical}
}

func isDiagonal(angle uint) bool {
	return math.Mod(float64(angle), 90) == 45
}

func isHorizontal(angle uint) bool {
	return angle == 0 || angle == 180
}

func isVertical(angle uint) bool {
	return angle == 90 || angle == 270
}

func isEdgeZone(zone ZoneT) bool {
	return gofuncs.EndsWith(string(zone), string(EdgeZoneSuffix))
}

func resolveThreeDirectionalValue[T gofuncs.Number](angle uint, diagonal, horizontal, vertical T) T {
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

func getThresholdValue(angle uint, threshold *ThresholdPctT) float64 {
	return resolveThreeDirectionalValue(angle, threshold.diagonal, threshold.horizontal, threshold.vertical)
}

func getAngleMargin(angle uint, angleMargin *AngleMarginT) uint {
	return resolveThreeDirectionalValue(angle, angleMargin.Diagonal, angleMargin.Horizontal, angleMargin.Vertical)
}

func (zoneThreshold *ThresholdPctT) Validate() {
	gofuncs.PanicAnyNotPositive(
		zoneThreshold.diagonal,
		zoneThreshold.vertical,
		zoneThreshold.horizontal)
}

func (zoneThreshold *ThresholdPctT) ValidateEdgeThreshold(edgeThreshold *ThresholdPctT) {
	if gofuncs.AnyGreaterOrEqual([][]float64{
		{zoneThreshold.diagonal, edgeThreshold.diagonal},
		{zoneThreshold.vertical, edgeThreshold.vertical},
		{zoneThreshold.horizontal, edgeThreshold.horizontal},
	}) {
		gofuncs.Panic("Edge threshold can't be less or equal to Zone threshold")
	}
}

func (angleMargin *AngleMarginT) Validate() {
	gofuncs.PanicAnyNotPositive(angleMargin.Horizontal, angleMargin.Vertical)

	if gofuncs.AnyGreaterOrEqual([][]uint{
		{angleMargin.Diagonal + angleMargin.Horizontal, 45},
		{angleMargin.Diagonal + angleMargin.Vertical, 45},
	}) {
		gofuncs.Panic("With this margin of angle areas will overlap")
	}
}

func genRange(lowerBound, upperBound uint, _boundariesMap ZoneBoundariesMapT, zone ZoneT, zoneThreshold, edgeThreshold float64) {
	lowerBound += 360
	upperBound += 360

	for angle := lowerBound; angle <= upperBound; angle++ {
		resolvedAngle := resolveCircleAngle(angle)
		gofuncs.AssignWithDuplicateKeyValueCheck(_boundariesMap, resolvedAngle, MakeDirection(zone, zoneThreshold, edgeThreshold), false)
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

func genEqualThresholdBoundariesMap(includeDiagonalZones bool, angleMargin *AngleMarginT, zoneThreshold, edgeThreshold float64) ZoneBoundariesMapT {
	return genBoundariesMap(includeDiagonalZones, angleMargin,
		MakeThresholdPct(zoneThreshold, zoneThreshold, zoneThreshold),
		MakeThresholdPct(edgeThreshold, edgeThreshold, edgeThreshold))
}

func genBoundariesMap(includeDiagonalZones bool, angleMargin *AngleMarginT, zoneThreshold, edgeThreshold *ThresholdPctT) ZoneBoundariesMapT {
	angleMargin.Validate()
	zoneThreshold.Validate()
	zoneThreshold.ValidateEdgeThreshold(edgeThreshold)

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
	for _, zone := range InitAllZones() {
		gofuncs.Print("%v: ", zone)
		var needAngles []int
		for angle, dir := range _boundariesMap {
			if dir.zone == zone {
				needAngles = append(needAngles, int(angle))
			}
		}
		sort.Ints(needAngles)
		fmt.Println(needAngles)
	}
}

type DetectZoneFuncT func(boundariesMap ZoneBoundariesMapT) ZoneT

func (pad *PadStickPositionT) GetDetectZoneFunc() DetectZoneFuncT {
	FloatEqualityMargin := pad.cfg.Math.FloatEqualityMargin

	isGreaterThanThreshold := func(thresholdPct float64) bool {
		return pad.magnitude > thresholdPct*pad.radius+FloatEqualityMargin
	}

	return func(boundariesMap ZoneBoundariesMapT) ZoneT {
		if direction, found := boundariesMap[pad.shiftedAngle]; found {
			if isGreaterThanThreshold(direction.zoneThresholdPct) {
				zone := direction.zone

				if isGreaterThanThreshold(direction.edgeThresholdPct) {
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
}

func (pad *PadStickPositionT) ReCalculateZone(zoneBoundariesMap ZoneBoundariesMapT) {
	if pad.newValueHandled {
		return
	}
	pad.newValueHandled = true

	zone := pad.detectZone(zoneBoundariesMap)
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
