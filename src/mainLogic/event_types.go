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

const HoldSuffix = "_Hold"

func addHoldSuffix(btn BtnOrAxisT) BtnOrAxisT {
	return BtnOrAxisT(string(btn) + HoldSuffix)
}

func removeHoldSuffix(btn BtnOrAxisT) BtnOrAxisT {
	return BtnOrAxisT(strings.TrimSuffix(string(btn), HoldSuffix))
}

const (
	BtnLeftWing     BtnOrAxisT = "LeftWing"
	BtnRightWing    BtnOrAxisT = "RightWing"
	BtnB            BtnOrAxisT = "B"
	BtnY            BtnOrAxisT = "Y"
	BtnX            BtnOrAxisT = "X"
	BtnA            BtnOrAxisT = "A"
	BtnC            BtnOrAxisT = "BtnC"
	BtnZ            BtnOrAxisT = "BtnZ"
	BtnLeftButton   BtnOrAxisT = "LB"
	BtnLeftTrigger  BtnOrAxisT = "LT"
	BtnRightButton  BtnOrAxisT = "RB"
	BtnRightTrigger BtnOrAxisT = "RT"
	BtnSelect       BtnOrAxisT = "Select"
	BtnStart        BtnOrAxisT = "Start"
	BtnMode         BtnOrAxisT = "Mode"
	BtnStick        BtnOrAxisT = "Stick"

	BtnRightPad  BtnOrAxisT = "RightPad"
	BtnLeftPad   BtnOrAxisT = "LeftPad"
	BtnDPadUp    BtnOrAxisT = BtnLeftPad
	BtnDPadDown  BtnOrAxisT = BtnLeftPad
	BtnDPadLeft  BtnOrAxisT = BtnLeftPad
	BtnDPadRight BtnOrAxisT = BtnLeftPad

	BtnStickUp    BtnOrAxisT = "StickUp"
	BtnStickDown  BtnOrAxisT = "StickDown"
	BtnStickLeft  BtnOrAxisT = "StickLeft"
	BtnStickRight BtnOrAxisT = "StickRight"

	BtnUnknown BtnOrAxisT = "BtnUnknown"
)

type Synonyms map[BtnOrAxisT]BtnOrAxisT

func genBtnSynonyms() Synonyms {
	synonyms := Synonyms{
		"LeftButton":   BtnLeftButton,
		"LeftTrigger":  BtnLeftTrigger,
		"RightButton":  BtnRightButton,
		"RightTrigger": BtnRightTrigger,
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
	BtnA,
	BtnB,
	BtnY,
	BtnX,
	BtnC,
	BtnZ,
	BtnLeftButton,
	BtnLeftTrigger,
	BtnRightButton,
	BtnRightTrigger,
	BtnSelect,
	BtnStart,
	BtnMode,
	BtnStick,
	BtnLeftPad,
	BtnRightPad,

	BtnStickUp,
	BtnStickDown,
	BtnStickLeft,
	BtnStickRight,
}

var _BtnMap = map[uint8]BtnOrAxisT{
	'a': BtnA,
	'b': BtnB,
	'c': BtnY,
	'd': BtnX,
	'e': BtnC,
	'f': BtnZ,
	'g': BtnLeftButton,
	'h': BtnLeftTrigger,
	'i': BtnRightButton,
	'j': BtnRightTrigger,
	'k': BtnSelect,
	'l': BtnStart,
	'm': BtnMode,
	'n': BtnStick,
	'o': BtnRightPad,
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
