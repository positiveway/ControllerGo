package mainLogic

import (
	"github.com/positiveway/gofuncs"
)

type ControllerInUseT string

const (
	SteamController ControllerInUseT = "SteamController"
	DualShock       ControllerInUseT = "DualShock"
	SteamDeck       ControllerInUseT = "SteamDeck"
)

func checkEnumCfg[T comparable](allEnumVariants []T, cfgValue T) T {
	if !gofuncs.Contains(allEnumVariants, cfgValue) {
		gofuncs.Panic("Incorrect enum type: %v", cfgValue)
	}

	return cfgValue
}

func toControllerCfg() ControllerInUseT {
	allControllers := []ControllerInUseT{SteamController, DualShock, SteamDeck}
	return checkEnumCfg(allControllers, ControllerInUseT(_RawCfg.ControllerInUse))
}

func toPadsSticksModeCfg() *PadsSticksModeT {
	allModes := []ModeT{MouseMode, GamingMode}
	return MakePadsSticksMode(checkEnumCfg(allModes, ModeT(_RawCfg.PadsSticks.Mode)))
}

func (c *ConfigsT) setConfigConstants() {
	//Debug
	gofuncs.PrintDebugInfo = false

	c.System.RunFromTerminal = false
}

func (c *ConfigsT) initDependent() {
	c.PadsSticks.HighPrecisionMode = MakeHighPrecisionMode()

	switch c.ControllerInUse {
	case SteamController:
		RightPadStick = MakePadPosition(
			_RawCfg.PadsSticks.Rotation.RightPad, false)
		LeftStick = MakePadPosition(
			_RawCfg.PadsSticks.Rotation.LeftStick, true)
		LeftPad = MakePadPosition(
			_RawCfg.PadsSticks.Rotation.LeftPad, true)

		c.PadsSticks.MousePS = RightPadStick
		c.PadsSticks.ScrollPS = LeftPad

		c.Typing.RightPS = RightPadStick
		c.Typing.LeftPS = LeftPad

	case DualShock:
		RightPadStick = MakePadPosition(
			_RawCfg.PadsSticks.Rotation.RightStick, false)
		LeftStick = MakePadPosition(
			_RawCfg.PadsSticks.Rotation.LeftStick, true)

		c.PadsSticks.MousePS = RightPadStick
		c.PadsSticks.ScrollPS = LeftStick

		c.Typing.RightPS = RightPadStick
		c.Typing.LeftPS = LeftStick
	}

	if !c.Mouse.OnRightStickPad {
		gofuncs.Swap(c.PadsSticks.MousePS, c.PadsSticks.ScrollPS)
	}
}

func (c *ConfigsT) setConfigVars() {
	c.System.TickerInterval = 1
	c.System.GCPercent = 10000

	c.ControllerInUse = toControllerCfg()
	c.PadsSticks.Mode = toPadsSticksModeCfg()

	c.SharedCfgT = _RawCfg.SharedCfgT

	c.Math.OutputMin = 0
	c.Math.FloatEqualityMargin = 0.000000000000001

	gofuncs.SetDefaultIfValueIsEmpty(
		&c.WebSocket.Port, 1234)

	gofuncs.SetDefaultIfValueIsEmpty(
		&c.WebSocket.IP, "0.0.0.0")

	gofuncs.SetDefaultIfValueIsEmpty(
		&(c.Buttons.HoldRepeatInterval), 40)

	c.PadsSticks.MinStandardRadius = gofuncs.GetValueOrDefaultIfEmpty(
		_RawCfg.PadsSticks.MinStandardRadius, 1.0)

	c.SharedCfgT.ValidateConvert(c.PadsSticks.Mode)

	switch c.ControllerInUse {
	case SteamController:
		c.Mouse.Speed.Validate()

		//init Stick map
		stickAngleMarginSC := _RawCfg.PadsSticks.Stick.AngleMargin
		stickThresholdSC := gofuncs.NumberToPct(_RawCfg.PadsSticks.Stick.ThresholdPct)
		stickEdgeThresholdSC := gofuncs.NumberToPct(_RawCfg.PadsSticks.Stick.EdgeThresholdPct)

		c.PadsSticks.Stick.BoundariesMapSC = genEqualThresholdBoundariesMap(false,
			MakeAngleMargin(0, stickAngleMarginSC, stickAngleMarginSC),
			stickThresholdSC,
			stickEdgeThresholdSC)

	case DualShock:
		c.Mouse.Intervals.Validate()
		c.PadsSticks.Stick.DeadzoneDS = _RawCfg.PadsSticks.Stick.Deadzone
		gofuncs.PanicAnyNotPositive(c.PadsSticks.Stick.DeadzoneDS)
	}

	switch c.PadsSticks.Mode.GetMode() {
	case GamingMode:
		c.Gaming.MoveIntervals = _RawCfg.Gaming.MoveIntervals
	}
}

var Cfg *ConfigsT
var _RawCfg *RawConfigsT

type MouseCfgT struct {
	OnRightStickPad  bool                 `json:"OnRightStickPad"`
	Intervals        PrecisionsIntervalsT `json:"Intervals"`
	Speed            PrecisionsSpeedT     `json:"Speed"`
	EdgeThresholdPct float64              `json:"EdgeThresholdPct"`
}

func (rawMouseCfg *MouseCfgT) ValidateConvert() {
	gofuncs.NumberToPctInPlace(&rawMouseCfg.EdgeThresholdPct)
}

type ScrollCfgT struct {
	HorizontalThresholdPct float64              `json:"HorizontalThresholdPct"`
	Intervals              PrecisionsIntervalsT `json:"Intervals"`
}

func (scrollCfg *ScrollCfgT) ValidateConvert() {
	gofuncs.NumberToPctInPlace(&scrollCfg.HorizontalThresholdPct)
	scrollCfg.Intervals.Validate()
}

type GamingCfgT struct {
	MoveIntervals IntervalRangeT `json:"MoveIntervals"`
}

func (gamingCfg *GamingCfgT) ValidateConvert() {
	gamingCfg.MoveIntervals.Validate()
}

type ButtonsCfgT struct {
	TriggerThreshold      float64 `json:"TriggerThresholdPct"`
	HoldingStateThreshold float64 `json:"HoldingStateThresholdMs"`
	HoldRepeatInterval    float64 `json:"HoldRepeatIntervalMs"`
}

func (rawButtonsCfg *ButtonsCfgT) ValidateConvert() {
	gofuncs.NumberToPctInPlace(&rawButtonsCfg.TriggerThreshold)
	gofuncs.PanicAnyNotInteger(rawButtonsCfg.HoldingStateThreshold, rawButtonsCfg.HoldingStateThreshold)
}

type TypingCfgT struct {
	LeftPS, RightPS *PadStickPositionT
	ThresholdPct    float64 `json:"ThresholdPct"`
	AngleMargin     struct {
		Straight uint `json:"Straight"`
		Diagonal uint `json:"Diagonal"`
	} `json:"AngleMargin"`
}

func (typingCfg *TypingCfgT) ValidateConvert() {
	gofuncs.NumberToPctInPlace(&typingCfg.ThresholdPct)
}

type WebSocketCfgT struct {
	Port int    `json:"Port"`
	IP   string `json:"IP"`
}

func (webSocketCfg *WebSocketCfgT) ValidateConvert() {
	if !gofuncs.IsEmpty(webSocketCfg.Port) {
		gofuncs.PanicAnyNotPositive(webSocketCfg.Port)
	}
}

type SharedCfgT struct {
	Mouse     MouseCfgT     `json:"Mouse"`
	Scroll    ScrollCfgT    `json:"Scroll"`
	Gaming    GamingCfgT    `json:"Gaming"`
	Buttons   ButtonsCfgT   `json:"Buttons"`
	Typing    TypingCfgT    `json:"Typing"`
	WebSocket WebSocketCfgT `json:"WebSocket"`
}

func (sharedCfg *SharedCfgT) ValidateConvert(padSticksMode *PadsSticksModeT) {
	sharedCfg.Mouse.ValidateConvert()
	sharedCfg.Scroll.ValidateConvert()
	if padSticksMode.GetMode() == GamingMode {
		sharedCfg.Gaming.ValidateConvert()
	}
	sharedCfg.Buttons.ValidateConvert()
	sharedCfg.Typing.ValidateConvert()
	sharedCfg.WebSocket.ValidateConvert()
}

type ConfigsT struct {
	ControllerInUse ControllerInUseT

	System struct {
		TickerInterval  float64
		GCPercent       int
		RunFromTerminal bool
	}

	PadsSticks struct {
		Mode              *PadsSticksModeT
		HighPrecisionMode *HighPrecisionModeT
		MinStandardRadius float64

		MousePS, ScrollPS *PadStickPositionT

		Stick struct {
			BoundariesMapSC ZoneBoundariesMapT
			DeadzoneDS      float64
		}
	}

	Math struct {
		OutputMin           float64
		FloatEqualityMargin float64
	}

	SharedCfgT

	Path struct {
		AllLayoutsDir, CurLayoutDir string
	}
}

type RawConfigsT struct {
	ControllerInUse string `json:"ControllerInUse"`
	PadsSticks      struct {
		Mode              string  `json:"Mode"`
		MinStandardRadius float64 `json:"MinStandardRadius"`
		Rotation          struct {
			LeftStick  float64 `json:"LeftStick"`
			RightStick float64 `json:"RightStick"`
			LeftPad    float64 `json:"LeftPad"`
			RightPad   float64 `json:"RightPad"`
			Stick      float64 `json:"Stick"`
		} `json:"Rotation"`
		Stick struct {
			Deadzone         float64 `json:"Deadzone"`
			AngleMargin      uint    `json:"AngleMargin"`
			ThresholdPct     uint    `json:"ThresholdPct"`
			EdgeThresholdPct uint    `json:"EdgeThresholdPct"`
		} `json:"Stick"`
	} `json:"Pads/Sticks"`
	SharedCfgT
}
