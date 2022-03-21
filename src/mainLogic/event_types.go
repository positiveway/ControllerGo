package mainLogic

import (
	"strings"
)

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
	CodeLeftWing  CodeT = 336
	CodeRightWing CodeT = 337
)

var UnknownCodesResolvingMap = map[CodeT]BtnOrAxisT{
	CodeLeftPadX:  AxisLeftPadX,
	CodeLeftPadY:  AxisLeftPadY,
	CodeRightPadX: AxisRightPadX,
	CodeRightPadY: AxisRightPadY,
	CodeLeftWing:  BtnLeftWing,
	CodeRightWing: BtnRightWing,
}

type BtnOrAxisT string

const (
	AxisLeftStickX BtnOrAxisT = "LeftStickX"
	AxisLeftStickY BtnOrAxisT = "LeftStickY"
	AxisLeftZ      BtnOrAxisT = "LeftZ"
	AxisRightPadX  BtnOrAxisT = "RightPadX"
	AxisRightPadY  BtnOrAxisT = "RightPadY"
	AxisRightZ     BtnOrAxisT = "RightZ"
	AxisLeftPadX   BtnOrAxisT = "LeftPadX"
	AxisLeftPadY   BtnOrAxisT = "LeftPadY"
	AxisUnknown    BtnOrAxisT = "Unknown"
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
	BtnLeftWing      BtnOrAxisT = "LeftWing"
	BtnRightWing     BtnOrAxisT = "RightWing"
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
	BtnLeftStick     BtnOrAxisT = "LeftStick"
	BtnRightStick    BtnOrAxisT = "RightStick"
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
		"LeftStick":     BtnLeftStick,
		"RightStick":    BtnRightStick,
	}
	for key, val := range synonyms {
		synonyms[addHoldSuffix(key)] = addHoldSuffix(val)
	}
	return synonyms
}

var BtnSynonyms = genBtnSynonyms()

var AllAvailableButtons = []BtnOrAxisT{
	BtnLeftWing,
	BtnRightWing,
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
	BtnLeftStick,
	BtnRightStick,
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
	'n': BtnLeftStick,
	'o': BtnRightStick,
	'p': BtnDPadUp,
	'q': BtnDPadDown,
	'r': BtnDPadLeft,
	's': BtnDPadRight,
	't': BtnUnknown,
}

type EventTypeT string

const (
	EvAxisChanged     EventTypeT = "AxisChanged"
	EvButtonChanged   EventTypeT = "ButtonChanged"
	EvButtonReleased  EventTypeT = "ButtonReleased"
	EvButtonPressed   EventTypeT = "ButtonPressed"
	EvButtonRepeated  EventTypeT = "ButtonRepeated"
	EvConnected       EventTypeT = "Connected"
	EvDisconnected    EventTypeT = "Disconnected"
	EvDropped         EventTypeT = "Dropped"
	EvPadFirstTouched EventTypeT = "PadFirstTouched"
	EvPadReleased     EventTypeT = "PadReleased"
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
