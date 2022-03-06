package mainLogic

import (
	"math"
	"strconv"
	"strings"
)

type Adjustment = [2]float64

const adjMult = 0.08

var AxesAdjustments = map[string]Adjustment{
	AxisRightPadX: {adjMult, adjMult},
	AxisRightPadY: {0.0, adjMult},
	AxisLeftPadX:  {adjMult, adjMult},
	AxisLeftPadY:  {0.0, adjMult},
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
	eventType string
	btnOrAxis string
	value     float64
	codeType  string
	code      int
}

func (event *Event) convertToAxisChanged() {
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

func (event *Event) filterEvents() {
	if event.eventType == EvAxisChanged {
		if adjustment, found := AxesAdjustments[event.btnOrAxis]; found {
			if event.value == 0.0 {
				return
			}
			if math.Abs(event.value) > adjustmentThreshold {
				negAdj, posAdj := adjustment[0], adjustment[1]

				switch true {
				case event.value > 0:
					event.value = math.Min(event.value*posAdj, 1.0)
				case event.value < 0:
					event.value = math.Max(event.value*negAdj, -1.0)
				}
			}
		}
	}

	//fmt.Printf("Before: ")
	//event.print()

	event.convertToAxisChanged()

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
			event.codeType = typeAndCode[0]

			code := typeAndCode[1]
			event.code, err = strconv.Atoi(code[:len(code)-1])
			CheckErr(err)
		}
	}
	event.filterEvents()
}

func (event *Event) print() {
	print("%s %s %s %v %0.2f",
		strings.TrimPrefix(event.eventType, "Ev"),
		strings.TrimPrefix(strings.TrimPrefix(event.btnOrAxis, "Btn"), "Axis"),
		event.codeType,
		event.code,
		event.value)
}

const (
	CTAbs string = "ABS"
	CTKey        = "KEY"
)

const (
	CodeLeftPadX  int = 16
	CodeLeftPadY      = 17
	CodeRightPadX     = 3
	CodeRightPadY     = 4
)

var CodeToAxisMap = map[int]string{
	CodeLeftPadX:  AxisLeftPadX,
	CodeLeftPadY:  AxisLeftPadY,
	CodeRightPadX: AxisRightPadX,
	CodeRightPadY: AxisRightPadY,
}

const (
	AxisLeftStickX string = "AxisLeftStickX"
	AxisLeftStickY        = "AxisLeftStickY"
	AxisLeftZ             = "AxisLeftZ"
	AxisRightPadX         = "AxisRightPadX"
	AxisRightPadY         = "AxisRightPadY"
	AxisRightZ            = "AxisRightZ"
	AxisLeftPadX          = "AxisLeftPadX"
	AxisLeftPadY          = "AxisLeftPadY"
	AxisUnknown           = "AxisUnknown"
)

var _AxisMap = map[uint8]string{
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

func addHoldSuffix(btn string) string {
	return btn + HoldSuffix
}

func removeHoldSuffix(btn string) string {
	return strings.TrimSuffix(btn, HoldSuffix)
}

const (
	BtnSouth         string = "South"
	BtnEast                 = "East"
	BtnNorth                = "North"
	BtnWest                 = "West"
	BtnC                    = "BtnC"
	BtnZ                    = "BtnZ"
	BtnLeftTrigger          = "LB"
	BtnLeftTrigger2         = "LT"
	BtnRightTrigger         = "RB"
	BtnRightTrigger2        = "RT"
	BtnSelect               = "Select"
	BtnStart                = "Start"
	BtnMode                 = "Mode"
	BtnLeftThumb            = "LeftThumb"
	BtnRightThumb           = "RightThumb"
	BtnDPadUp               = "DPadUp"
	BtnDPadDown             = "DPadDown"
	BtnDPadLeft             = "DPadLeft"
	BtnDPadRight            = "DPadRight"
	BtnUnknown              = "BtnUnknown"
)

type Synonyms = map[string]string

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

var AllOriginalButtons = []string{
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

var _BtnMap = map[uint8]string{
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

const (
	EvAxisChanged     string = "EvAxisChanged"
	EvButtonChanged          = "EvButtonChanged"
	EvButtonReleased         = "EvButtonReleased"
	EvButtonPressed          = "EvButtonPressed"
	EvButtonRepeated         = "EvButtonRepeated"
	EvConnected              = "EvConnected"
	EvDisconnected           = "EvDisconnected"
	EvDropped                = "EvDropped"
	EvPadFirstTouched        = "EvPadFirstTouched"
	EvPadReleased            = "EvPadReleased"
)

var ButtonEvents = []string{EvButtonChanged, EvButtonReleased, EvButtonPressed, EvButtonRepeated}

var EventTypeMap = map[uint8]string{
	'a': EvAxisChanged,
	'b': EvButtonChanged,
	'c': EvButtonReleased,
	'd': EvButtonPressed,
	'e': EvButtonRepeated,
	'f': EvConnected,
	'g': EvDisconnected,
	'h': EvDropped,
}

func genBtnAxisMap() map[uint8]string {
	mapping := map[uint8]string{}
	for k, v := range _AxisMap {
		AssignWithDuplicateCheck(mapping, k, v)
	}
	for k, v := range _BtnMap {
		AssignWithDuplicateCheck(mapping, k, v)
	}
	return mapping
}

var BtnAxisMap = genBtnAxisMap()
