package mainLogic

import (
	"github.com/positiveway/gofuncs"
	"strings"
)

type CodeTypeT string

const (
	CTAbs CodeTypeT = "ABS"
	CTKey CodeTypeT = "KEY"
)

type CodeT int

const (
	//CodeStickX  CodeT = 0
	//CodeStickY  CodeT = 1
	//CodeLeftPadX  CodeT = 16
	//CodeLeftPadY  CodeT = 17
	//CodeRightPadX CodeT = 3
	//CodeRightPadY CodeT = 4
	CodeLeftWing  CodeT = 336
	CodeRightWing CodeT = 337
)

var UnknownCodesResolvingMap = map[CodeT]BtnOrAxisT{
	//CodeStickX: AxisStickX,
	//CodeStickY: AxisStickY,
	//CodeLeftPadX:  AxisLeftPadX,
	//CodeLeftPadY:  AxisLeftPadY,
	//CodeRightPadX: AxisRightPadX,
	//CodeRightPadY: AxisRightPadY,
	CodeLeftWing:  BtnLeftWing,
	CodeRightWing: BtnRightWing,
}

type BtnOrAxisT string

const (
	AxisStickX    BtnOrAxisT = "StickX"
	AxisStickY    BtnOrAxisT = "StickY"
	AxisStickZ    BtnOrAxisT = "StickZ"
	AxisRightPadX BtnOrAxisT = "RightPadX"
	AxisRightPadY BtnOrAxisT = "RightPadY"
	AxisRightPadZ BtnOrAxisT = "RightPadZ"
	AxisLeftPadX  BtnOrAxisT = "LeftPadX"
	AxisLeftPadY  BtnOrAxisT = "LeftPadY"
	AxisUnknown   BtnOrAxisT = "Unknown"
)

var _AxisMap = map[uint8]BtnOrAxisT{
	'u': AxisStickX,
	'v': AxisStickY,
	'w': AxisStickZ,
	'x': AxisRightPadX,
	'y': AxisRightPadY,
	'z': AxisRightPadZ,
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
	BtnDPadUp               = BtnLeftPad
	BtnDPadDown             = BtnLeftPad
	BtnDPadLeft             = BtnLeftPad
	BtnDPadRight            = BtnLeftPad

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
	EvAxisChanged    EventTypeT = "AxisChanged"
	EvButtonChanged  EventTypeT = "ButtonChanged"
	EvButtonReleased EventTypeT = "ButtonReleased"
	EvButtonPressed  EventTypeT = "ButtonPressed"
	EvButtonRepeated EventTypeT = "ButtonRepeated"
	EvConnected      EventTypeT = "Connected"
	EvDisconnected   EventTypeT = "Disconnected"
	EvDropped        EventTypeT = "Dropped"
	EvPadReleased    EventTypeT = "PadReleased"
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
		gofuncs.AssignWithDuplicateCheck(mapping, k, v)
	}
	for k, v := range _BtnMap {
		gofuncs.AssignWithDuplicateCheck(mapping, k, v)
	}
	return mapping
}

var BtnAxisMap = genBtnAxisMap()
