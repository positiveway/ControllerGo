package mainLogic

import (
	"math"
	"strconv"
	"strings"
)

type Adjustment [2]float64

const adjMultX = 0.08
const adjMultY = 0.04

var AxesAdjustments = map[BtnOrAxisT]Adjustment{
	AxisRightPadX: {adjMultX, adjMultX},
	AxisRightPadY: {0.0, adjMultY},
	AxisLeftPadX:  {adjMultX, adjMultX},
	AxisLeftPadY:  {0.0, adjMultY},
}

var adjustmentThreshold float64 = 0.8

func checkAdj(value *float64) {
	if *value < 0 {
		panicMsg("Adjustment value can't be negative")
	}
	if math.Abs(*value) >= 0.2 {
		panicMsg("Adjustment value is too high")
	}
	*value += 1
}

func checkAdjustments() {
	for axis, adjustment := range AxesAdjustments {
		negAdj, posAdj := adjustment[0], adjustment[1]
		checkAdj(&negAdj)
		checkAdj(&posAdj)
		AxesAdjustments[axis] = Adjustment{negAdj, posAdj}
	}
}

type Event struct {
	eventType EventTypeT
	btnOrAxis BtnOrAxisT
	value     float64
	codeType  CodeTypeT
	code      CodeT
}

func (event *Event) transformToPadEvent() {
	if event.btnOrAxis == BtnUnknown && event.codeType == CTAbs {
		if axis, found := CodeToAxisMap[event.code]; found {
			switch event.eventType {
			case EvButtonPressed:
				event.eventType = EvPadFirstTouched
			case EvButtonReleased:
				event.eventType = EvPadReleased
			default:
				return
			}
			event.btnOrAxis = axis
			event.value = 0
			event.code = 0
			event.codeType = ""
		}
	}
}

func applyAdjustments(value float64, axis BtnOrAxisT) float64 {
	if adjustment, found := AxesAdjustments[axis]; found {
		if math.Abs(value) > adjustmentThreshold {
			negAdj, posAdj := adjustment[0], adjustment[1]

			switch {
			case value > 0:
				value = math.Min(value*posAdj, 1.0)
			case value < 0:
				value = math.Max(value*negAdj, -1.0)
			}
		}
	}
	return value
}

func (event *Event) transformStickToButtons() {

}

func (event *Event) transformAndFilter() {
	if event.eventType == EvAxisChanged {
		if _, found := AxesAdjustments[event.btnOrAxis]; found {
			if event.value == 0.0 {
				return
			}
		}
	}

	//fmt.Printf("Before: ")
	//event.print()

	event.transformToPadEvent()

	//fmt.Printf("After: ")
	event.print()

	matchEvent()
}

func (event *Event) update(msg string) {
	event.eventType, found = EventTypeMap[msg[0]]
	if !found {
		PanicMisspelled(string(msg[0]))
	}
	if event.eventType != EvConnected && event.eventType != EvDisconnected && event.eventType != EvDropped {
		event.btnOrAxis, found = BtnAxisMap[msg[1]]
		if !found {
			PanicMisspelled(string(msg[1]))
		}
		if event.eventType == EvAxisChanged || event.eventType == EvButtonChanged {
			msg = msg[2:]
			valueAndCode := strings.Split(msg, ";")

			event.value, err = strconv.ParseFloat(valueAndCode[0], 32)
			CheckErr(err)

			if strings.HasSuffix(msg, ";") {
				return
			}
			typeAndCode := strings.Split(valueAndCode[1], "(")
			event.codeType = CodeTypeT(typeAndCode[0])

			code := typeAndCode[1]
			codeNum, err := strconv.Atoi(code[:len(code)-1])
			CheckErr(err)
			event.code = CodeT(codeNum)
		}
	}
	event.transformAndFilter()
}

func (event *Event) print() {
	print("%s %s %s %v %0.2f",
		TrimAnyPrefix(string(event.eventType), "Ev"),
		TrimAnyPrefix(string(event.btnOrAxis), "Btn", "Axis"),
		event.codeType,
		event.code,
		event.value)
}

type CodeTypeT string

const (
	CTAbs CodeTypeT = "ABS"
	CTKey CodeTypeT = "KEY"
)

type CodeT int

const (
	CodeLeftPadX  CodeT = 16
	CodeLeftPadY  CodeT = 17
	CodeRightPadX CodeT = 3
	CodeRightPadY CodeT = 4
)

var CodeToAxisMap = map[CodeT]BtnOrAxisT{
	CodeLeftPadX:  AxisLeftPadX,
	CodeLeftPadY:  AxisLeftPadY,
	CodeRightPadX: AxisRightPadX,
	CodeRightPadY: AxisRightPadY,
}

type BtnOrAxisT string

const (
	AxisLeftStickX BtnOrAxisT = "AxisLeftStickX"
	AxisLeftStickY BtnOrAxisT = "AxisLeftStickY"
	AxisLeftZ      BtnOrAxisT = "AxisLeftZ"
	AxisRightPadX  BtnOrAxisT = "AxisRightPadX"
	AxisRightPadY  BtnOrAxisT = "AxisRightPadY"
	AxisRightZ     BtnOrAxisT = "AxisRightZ"
	AxisLeftPadX   BtnOrAxisT = "AxisLeftPadX"
	AxisLeftPadY   BtnOrAxisT = "AxisLeftPadY"
	AxisUnknown    BtnOrAxisT = "AxisUnknown"
)

var _AxisMap = map[uint8]BtnOrAxisT{
	'u': AxisLeftStickX,
	'v': AxisLeftStickY,
	'w': AxisLeftZ,
	'x': AxisRightPadX,
	'y': AxisRightPadY,
	'z': AxisRightZ,
	'0': AxisLeftPadX,
	'1': AxisLeftPadY,
	'2': AxisUnknown,
}

const HoldSuffix = "Hold"

func addHoldSuffix(btn BtnOrAxisT) BtnOrAxisT {
	return BtnOrAxisT(string(btn) + HoldSuffix)
}

func removeHoldSuffix(btn BtnOrAxisT) BtnOrAxisT {
	return BtnOrAxisT(strings.TrimSuffix(string(btn), HoldSuffix))
}

const (
	BtnSouth         BtnOrAxisT = "South"
	BtnEast          BtnOrAxisT = "East"
	BtnNorth         BtnOrAxisT = "North"
	BtnWest          BtnOrAxisT = "West"
	BtnC             BtnOrAxisT = "BtnC"
	BtnZ             BtnOrAxisT = "BtnZ"
	BtnLeftTrigger   BtnOrAxisT = "LB"
	BtnLeftTrigger2  BtnOrAxisT = "LT"
	BtnRightTrigger  BtnOrAxisT = "RB"
	BtnRightTrigger2 BtnOrAxisT = "RT"
	BtnSelect        BtnOrAxisT = "Select"
	BtnStart         BtnOrAxisT = "Start"
	BtnMode          BtnOrAxisT = "Mode"
	BtnLeftThumb     BtnOrAxisT = "LeftThumb"
	BtnRightThumb    BtnOrAxisT = "RightThumb"
	BtnDPadUp        BtnOrAxisT = "DPadUp"
	BtnDPadDown      BtnOrAxisT = "DPadDown"
	BtnDPadLeft      BtnOrAxisT = "DPadLeft"
	BtnDPadRight     BtnOrAxisT = "DPadRight"
	BtnUnknown       BtnOrAxisT = "BtnUnknown"
)

type Synonyms map[BtnOrAxisT]BtnOrAxisT

func genBtnSynonyms() Synonyms {
	synonyms := Synonyms{
		"LeftTrigger":   BtnLeftTrigger,
		"LeftTrigger2":  BtnLeftTrigger2,
		"RightTrigger":  BtnRightTrigger,
		"RightTrigger2": BtnRightTrigger2,
		"LeftStick":     BtnLeftThumb,
		"RightStick":    BtnRightThumb,
	}
	for key, val := range synonyms {
		synonyms[addHoldSuffix(key)] = addHoldSuffix(val)
	}
	return synonyms
}

var BtnSynonyms = genBtnSynonyms()

var AllOriginalButtons = []BtnOrAxisT{
	BtnSouth,
	BtnEast,
	BtnNorth,
	BtnWest,
	BtnC,
	BtnZ,
	BtnLeftTrigger,
	BtnLeftTrigger2,
	BtnRightTrigger,
	BtnRightTrigger2,
	BtnSelect,
	BtnStart,
	BtnMode,
	BtnLeftThumb,
	BtnRightThumb,
	BtnDPadUp,
	BtnDPadDown,
	BtnDPadLeft,
	BtnDPadRight,
}

var _BtnMap = map[uint8]BtnOrAxisT{
	'a': BtnSouth,
	'b': BtnEast,
	'c': BtnNorth,
	'd': BtnWest,
	'e': BtnC,
	'f': BtnZ,
	'g': BtnLeftTrigger,
	'h': BtnLeftTrigger2,
	'i': BtnRightTrigger,
	'j': BtnRightTrigger2,
	'k': BtnSelect,
	'l': BtnStart,
	'm': BtnMode,
	'n': BtnLeftThumb,
	'o': BtnRightThumb,
	'p': BtnDPadUp,
	'q': BtnDPadDown,
	'r': BtnDPadLeft,
	's': BtnDPadRight,
	't': BtnUnknown,
}

type EventTypeT string

const (
	EvAxisChanged     EventTypeT = "EvAxisChanged"
	EvButtonChanged   EventTypeT = "EvButtonChanged"
	EvButtonReleased  EventTypeT = "EvButtonReleased"
	EvButtonPressed   EventTypeT = "EvButtonPressed"
	EvButtonRepeated  EventTypeT = "EvButtonRepeated"
	EvConnected       EventTypeT = "EvConnected"
	EvDisconnected    EventTypeT = "EvDisconnected"
	EvDropped         EventTypeT = "EvDropped"
	EvPadFirstTouched EventTypeT = "EvPadFirstTouched"
	EvPadReleased     EventTypeT = "EvPadReleased"
)

var EventTypeMap = map[uint8]EventTypeT{
	'a': EvAxisChanged,
	'b': EvButtonChanged,
	'c': EvButtonReleased,
	'd': EvButtonPressed,
	'e': EvButtonRepeated,
	'f': EvConnected,
	'g': EvDisconnected,
	'h': EvDropped,
}

func genBtnAxisMap() map[uint8]BtnOrAxisT {
	mapping := map[uint8]BtnOrAxisT{}
	for k, v := range _AxisMap {
		AssignWithDuplicateCheck(mapping, k, v)
	}
	for k, v := range _BtnMap {
		AssignWithDuplicateCheck(mapping, k, v)
	}
	return mapping
}

var BtnAxisMap = genBtnAxisMap()
