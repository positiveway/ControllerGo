package mainLogic

var CurPressedStickButtonSC *BtnOrAxisT

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

type ZoneToBtnMapT map[ZoneT]BtnOrAxisT

type AvailableButtonsT []BtnOrAxisT
