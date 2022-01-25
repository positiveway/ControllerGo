package main

import (
	"fmt"
	"strconv"
)

func check_err(err error) {
	if err != nil {
		panic(err)
	}
}

type Event struct {
	deviceID  int
	eventType string
	btnOrAxis string
	value     float64
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

const (
	BtnSouth         string = "BtnSouth"
	BtnEast                 = "BtnEast"
	BtnNorth                = "BtnNorth"
	BtnWest                 = "BtnWest"
	BtnC                    = "BtnC"
	BtnZ                    = "BtnZ"
	BtnLeftTrigger          = "BtnLeftTrigger"
	BtnLeftTrigger2         = "BtnLeftTrigger2"
	BtnRightTrigger         = "BtnRightTrigger"
	BtnRightTrigger2        = "BtnRightTrigger2"
	BtnSelect               = "BtnSelect"
	BtnStart                = "BtnStart"
	BtnMode                 = "BtnMode"
	BtnLeftThumb            = "BtnLeftThumb"
	BtnRightThumb           = "BtnRightThumb"
	BtnDPadUp               = "BtnDPadUp"
	BtnDPadDown             = "BtnDPadDown"
	BtnDPadLeft             = "BtnDPadLeft"
	BtnDPadRight            = "BtnDPadRight"
	BtnUnknown              = "BtnUnknown"
)

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

func contains_str(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func makeEvent(id, eventType, btnOrAxis, value string) Event {
	deviceId, err := strconv.Atoi(id)
	check_err(err)
	valueFloat, err := strconv.ParseFloat(value, 32)
	check_err(err)
	eventType, ok := EventTypeMap[eventType]
	if !ok {
		panic(fmt.Sprintf("no element %v\n", eventType))
	}
	if contains_str(ButtonEvents, eventType) {
		btnOrAxis = BtnMap[btnOrAxis]
	} else if eventType == EvAxisChanged {
		btnOrAxis = AxisMap[btnOrAxis]
	}
	event := Event{
		deviceID:  deviceId,
		value:     valueFloat,
		eventType: eventType,
		btnOrAxis: btnOrAxis,
	}
	return event
}
