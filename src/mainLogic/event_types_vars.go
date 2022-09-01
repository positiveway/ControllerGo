package mainLogic

var CurPressedStickButtonSC *BtnOrAxisT

var UnknownCodesResolvingMapSC map[CodeT]BtnOrAxisT

var BtnAxisMap BtnAxisMapT

var (
	AxisLeftPadX,
	AxisLeftPadY,

	AxisLeftStickX,
	AxisLeftStickY,
	AxisLeftStickZ,

	AxisRightPadStickX,
	AxisRightPadStickY,
	AxisRightPadStickZ BtnOrAxisT
)

var PadAndStickAxes []BtnOrAxisT

var (
	BtnLeftPad,
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

type ZoneToBtnMapT = map[ZoneT]BtnOrAxisT

var StickZoneToBtnMapSC ZoneToBtnMapT

type AvailableButtonsT []BtnOrAxisT

var EventTypeMap map[uint8]EventTypeT
