package mainLogic

import (
	"fmt"
	"strconv"
)

type Event struct {
	deviceID  int
	eventType string
	btnOrAxis string
	value     float64
}

func makeEvent(id, eventType, btnOrAxis, value string) Event {
	deviceId, err := strconv.Atoi(id)
	CheckErr(err)
	valueFloat, err := strconv.ParseFloat(value, 32)
	CheckErr(err)
	eventType, ok := EventTypeMap[eventType]
	if !ok {
		panic(fmt.Sprintf("no element %v\n", eventType))
	}
	if contains(ButtonEvents, eventType) {
		btnOrAxis = BtnMap[btnOrAxis]
	} else if eventType == EvAxisChanged {
		btnOrAxis = AxisMap[btnOrAxis]
	} else if btnOrAxis == "No" {
		btnOrAxis = "None"
	}
	event := Event{
		deviceID:  deviceId,
		value:     valueFloat,
		eventType: eventType,
		btnOrAxis: btnOrAxis,
	}
	return event
}

const (
	AxisLeftStickX  string = "AxisLeftStickX"
	AxisLeftStickY         = "AxisLeftStickY"
	AxisLeftZ              = "AxisLeftZ"
	AxisRightStickX        = "AxisRightStickX"
	AxisRightStickY        = "AxisRightStickY"
	AxisRightZ             = "AxisRightZ"
	AxisDPadX              = "AxisDPadX"
	AxisDPadY              = "AxisDPadY"
	AxisUnknown            = "AxisUnknown"
)

var AxisMap = map[string]string{
	"LX": AxisLeftStickX,
	"LY": AxisLeftStickY,
	"LZ": AxisLeftZ,
	"RX": AxisRightStickX,
	"RY": AxisRightStickY,
	"RZ": AxisRightZ,
	"DX": AxisDPadX,
	"DY": AxisDPadY,
	"U":  AxisUnknown,
}

const HoldSuffix = "Hold"

const (
	BtnSouth         string = "South"
	BtnEast                 = "East"
	BtnNorth                = "North"
	BtnWest                 = "West"
	BtnSouthHold            = BtnSouth + HoldSuffix
	BtnEastHold             = BtnEast + HoldSuffix
	BtnNorthHold            = BtnNorth + HoldSuffix
	BtnWestHold             = BtnWest + HoldSuffix
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
	BtnDPadUpHold           = BtnDPadUp + HoldSuffix
	BtnDPadDownHold         = BtnDPadDown + HoldSuffix
	BtnDPadLeftHold         = BtnDPadLeft + HoldSuffix
	BtnDPadRightHold        = BtnDPadRight + HoldSuffix
	BtnUnknown              = "BtnUnknown"
)

var BtnSynonyms = map[string]string{
	"LeftTrigger":   BtnLeftTrigger,
	"LeftTrigger2":  BtnLeftTrigger2,
	"RightTrigger":  BtnRightTrigger,
	"RightTrigger2": BtnRightTrigger2,
	"LeftStick":     BtnLeftThumb,
	"RightStick":    BtnRightThumb,
}

var AllButtons = []string{
	BtnSouth,
	BtnEast,
	BtnNorth,
	BtnWest,
	BtnSouthHold,
	BtnEastHold,
	BtnNorthHold,
	BtnWestHold,
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
	BtnDPadUpHold,
	BtnDPadDownHold,
	BtnDPadLeftHold,
	BtnDPadRightHold,
	BtnUnknown,
}

var BtnMap = map[string]string{
	"S":  BtnSouth,
	"E":  BtnEast,
	"N":  BtnNorth,
	"W":  BtnWest,
	"C":  BtnC,
	"Z":  BtnZ,
	"L":  BtnLeftTrigger,
	"L2": BtnLeftTrigger2,
	"R":  BtnRightTrigger,
	"R2": BtnRightTrigger2,
	"Se": BtnSelect,
	"St": BtnStart,
	"M":  BtnMode,
	"LT": BtnLeftThumb,
	"RT": BtnRightThumb,
	"DU": BtnDPadUp,
	"DD": BtnDPadDown,
	"DL": BtnDPadLeft,
	"DR": BtnDPadRight,
	"U":  BtnUnknown,
}

const (
	EvAxisChanged    string = "EvAxisChanged"
	EvButtonChanged         = "EvButtonChanged"
	EvButtonReleased        = "EvButtonReleased"
	EvButtonPressed         = "EvButtonPressed"
	EvButtonRepeated        = "EvButtonRepeated"
	EvConnected             = "EvConnected"
	EvDisconnected          = "EvDisconnected"
	EvDropped               = "EvDropped"
)

var ButtonEvents = []string{EvButtonChanged, EvButtonReleased, EvButtonPressed, EvButtonRepeated}

var EventTypeMap = map[string]string{
	"A":  EvAxisChanged,
	"B":  EvButtonChanged,
	"Rl": EvButtonReleased,
	"P":  EvButtonPressed,
	"Rp": EvButtonRepeated,
	"C":  EvConnected,
	"D":  EvDisconnected,
	"Dr": EvDropped,
}
