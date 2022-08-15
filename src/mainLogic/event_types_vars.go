package mainLogic

var UnknownCodesResolvingMapSC map[CodeT]BtnOrAxisT

var BtnAxisMap BtnAxisMapT

var (
	AxisLeftStickX,
	AxisLeftStickY,
	AxisLeftStickZ,
	AxisRightPadStickX,
	AxisRightPadStickY,
	AxisRightPadStickZ BtnOrAxisT
)

var _AxisMap BtnAxisMapT

var PadAndStickAxes []BtnOrAxisT

var (
	BtnLeftStick,
	BtnRightPadStick BtnOrAxisT

	BtnLeftWingSC,
	BtnRightWingSC BtnOrAxisT

	BtnStickUpSC,
	BtnStickDownSC,
	BtnStickLeftSC,
	BtnStickRightSC BtnOrAxisT

	BtnDPadUp,
	BtnDPadDown,
	BtnDPadLeft,
	BtnDPadRight BtnOrAxisT
)

var StickZoneToBtnMapSC map[Zone]BtnOrAxisT

var AllAvailableButtons []BtnOrAxisT

var _BtnMap BtnAxisMapT

var BtnSynonyms Synonyms

var EventTypeMap map[uint8]EventTypeT
