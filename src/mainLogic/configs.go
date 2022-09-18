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

func PanicUnsupportedController() {
	gofuncs.Panic("Unsupported controller type")
}

func checkEnumCfg[T comparable](allEnumVariants []T, cfgValue T) T {
	if !gofuncs.Contains(allEnumVariants, cfgValue) {
		gofuncs.Panic("Incorrect enum type: %v", cfgValue)
	}

	return cfgValue
}

func toControllerCfg(rawCfg *RawConfigsT) ControllerInUseT {
	allControllers := []ControllerInUseT{SteamController, DualShock, SteamDeck}
	return checkEnumCfg(allControllers, ControllerInUseT(rawCfg.ControllerInUse))
}

func toPadsSticksModeCfg(rawCfg *RawConfigsT) *PadsSticksModeT {
	allModes := []ModeT{MouseMode, GamingMode}
	return MakePadsSticksMode(checkEnumCfg(allModes, ModeT(rawCfg.PadsSticks.Mode)))
}

func (c *ConfigsT) setConfigConstants() {
	//Debug
	gofuncs.PrintDebugInfo = false

	c.System.RunFromTerminal = false
}

type DependentVariablesT struct {
	CfgStruct
	HighPrecisionMode *HighPrecisionModeT

	LeftPad, RightPadStick, LeftStick,
	MousePS, ScrollPS *PadStickPositionT

	allBtnAxis *AllBtnAxis
	Buttons    *ButtonsT
	Typing     *TypingT

	CurPressedStickButtonSC *BtnOrAxisT
}

func MakeDependentVariables() *DependentVariablesT {
	cfg, rawCfg := MakeConfigs()

	dependentVars := &DependentVariablesT{}
	dependentVars.Init(cfg)

	dependentVars.allBtnAxis = MakeAllBtnAxis(cfg)

	dependentVars.CurPressedStickButtonSC = InitCurStickButton()

	dependentVars.Buttons = &ButtonsT{}
	dependentVars.HighPrecisionMode = &HighPrecisionModeT{}
	dependentVars.Typing = &TypingT{}

	dependentVars.RightPadStick = &PadStickPositionT{}
	dependentVars.LeftStick = &PadStickPositionT{}
	dependentVars.LeftPad = &PadStickPositionT{}

	switch cfg.ControllerInUse {
	case SteamController:
		dependentVars.Typing.RightPS = dependentVars.RightPadStick
		dependentVars.Typing.LeftPS = dependentVars.LeftPad

	case DualShock:
		dependentVars.Typing.RightPS = dependentVars.RightPadStick
		dependentVars.Typing.LeftPS = dependentVars.LeftStick
	}

	dependentVars.MousePS = dependentVars.Typing.RightPS
	dependentVars.ScrollPS = dependentVars.Typing.LeftPS

	if !cfg.Mouse.OnRightStickPad {
		gofuncs.Swap(dependentVars.MousePS, dependentVars.ScrollPS)
	}

	dependentVars.Buttons.Init(cfg, dependentVars.HighPrecisionMode, dependentVars.allBtnAxis)

	dependentVars.RightPadStick.Init(
		rawCfg.PadsSticks.Rotation.RightPad, false, cfg, dependentVars.Buttons)
	dependentVars.LeftStick.Init(
		rawCfg.PadsSticks.Rotation.LeftStick, true, cfg, dependentVars.Buttons)
	dependentVars.LeftPad.Init(
		rawCfg.PadsSticks.Rotation.LeftPad, true, cfg, dependentVars.Buttons)

	dependentVars.HighPrecisionMode.Init(cfg, dependentVars.Buttons)
	dependentVars.Typing.Init(cfg, dependentVars.Buttons)

	switch cfg.ControllerInUse {
	case SteamController:
		dependentVars.MousePS.InitMoveSCFunc(dependentVars.HighPrecisionMode)
	}

	return dependentVars
}

func (c *ConfigsT) setConfigVars(rawCfg *RawConfigsT) {
	c.System.TickerInterval = 1
	c.System.GCPercent = 10000

	c.ControllerInUse = toControllerCfg(rawCfg)
	c.PadsSticks.Mode = toPadsSticksModeCfg(rawCfg)

	c.SharedCfgT = rawCfg.SharedCfgT

	c.Math.OutputMin = 0
	c.Math.FloatEqualityMargin = 0.000000000000001

	gofuncs.SetDefaultIfValueIsEmpty(
		&c.WebSocket.Port, 1234)

	gofuncs.SetDefaultIfValueIsEmpty(
		&c.WebSocket.IP, "0.0.0.0")

	gofuncs.SetDefaultIfValueIsEmpty(
		&(c.Buttons.HoldRepeatInterval), 40)

	c.PadsSticks.MinStandardRadius = gofuncs.GetValueOrDefaultIfEmpty(
		rawCfg.PadsSticks.MinStandardRadius, 1.0)

	c.SharedCfgT.ValidateConvert(c.PadsSticks.Mode)

	switch c.ControllerInUse {
	case SteamController:
		c.Mouse.Speed.Validate()

		//init Stick map
		stickAngleMarginSC := rawCfg.PadsSticks.Stick.AngleMargin
		stickThresholdSC := gofuncs.NumberToPct(rawCfg.PadsSticks.Stick.ThresholdPct)
		stickEdgeThresholdSC := gofuncs.NumberToPct(rawCfg.PadsSticks.Stick.EdgeThresholdPct)

		c.PadsSticks.Stick.BoundariesMapSC = genEqualThresholdBoundariesMap(false,
			MakeAngleMargin(0, stickAngleMarginSC, stickAngleMarginSC),
			stickThresholdSC,
			stickEdgeThresholdSC)

	case DualShock:
		c.Mouse.Intervals.Validate()
		c.PadsSticks.Stick.DeadzoneDS = rawCfg.PadsSticks.Stick.Deadzone
		gofuncs.PanicAnyNotPositive(c.PadsSticks.Stick.DeadzoneDS)
	}

	switch c.PadsSticks.Mode.GetMode() {
	case GamingMode:
		c.Gaming.MoveIntervals = rawCfg.Gaming.MoveIntervals
	}
}

type MouseCfgT struct {
	OnRightStickPad        bool                 `json:"OnRightStickPad"`
	Intervals              PrecisionsIntervalsT `json:"Intervals"`
	Speed                  PrecisionsSpeedT     `json:"Speed"`
	EdgeThresholdPct       float64              `json:"EdgeThresholdPct"`
	DoubleTouchMaxInterval float64              `json:"DoubleTouchMaxInterval"`
	ClickReleaseInterval   float64              `json:"ClickReleaseInterval"`
}

func (rawMouseCfg *MouseCfgT) ValidateConvert() {
	gofuncs.PanicAnyNotPositive(
		rawMouseCfg.DoubleTouchMaxInterval,
		rawMouseCfg.ClickReleaseInterval)

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
	ThresholdPct float64 `json:"ThresholdPct"`
	AngleMargin  struct {
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
		MinStandardRadius float64

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
