package mainLogic

import (
	"github.com/positiveway/gofuncs"
	"time"
)

type ControllerInUseT string

const (
	SteamController ControllerInUseT = "SteamController"
	DualShock       ControllerInUseT = "DualShock"
	SteamDeck       ControllerInUseT = "SteamDeck"
)

func checkEnumCfg[T comparable](allEnumVariants []T, cfgValue T) {
	if !gofuncs.Contains(allEnumVariants, cfgValue) {
		gofuncs.Panic("Incorrect enum type: %v", cfgValue)
	}
}

func (c *ConfigsT) toControllerCfg() ControllerInUseT {
	allControllers := []ControllerInUseT{SteamController, DualShock, SteamDeck}
	controller := ControllerInUseT(c.getConfig("ControllerInUse"))
	checkEnumCfg(allControllers, controller)
	return controller
}

func (c *ConfigsT) toPadsSticksModeCfg() *PadsSticksMode {
	allModes := []ModeType{MouseMode, GamingMode}
	modeType := ModeType(c.getConfig("PadsSticksMode"))
	checkEnumCfg(allModes, modeType)
	return MakePadsSticksMode(modeType)
}

func (c *ConfigsT) setConfigConstants() {
	//Debug
	gofuncs.PrintDebugInfo = false

	c.RunFromTerminal = false
}

func (c *ConfigsT) initTouchpads() {
	LeftPad = MakePadPosition()
	RightPadStick = MakePadPosition()
	LeftStick = MakePadPosition()

	switch c.ControllerInUse {
	case SteamController:
		RightPadStick.zoneRotation = c.toFloatConfig("RightPadRotation")
		LeftStick.zoneRotation = c.toFloatConfig("StickRotation")
		LeftPad.zoneRotation = c.toFloatConfig("LeftPadRotation")

		c.mousePadStick = RightPadStick
		c.scrollPadStick = LeftPad
	case DualShock:
		RightPadStick.zoneRotation = c.toFloatConfig("RightStickRotation")
		LeftStick.zoneRotation = c.toFloatConfig("LeftStickRotation")

		c.mousePadStick = RightPadStick
		c.scrollPadStick = LeftStick
	}
}

func (c *ConfigsT) setConfigVars() {
	c.ControllerInUse = c.toControllerCfg()

	//Math
	c.OutputMin = 0
	c.MinStandardPadRadius = 1.0

	// web socket
	c.SocketPort = 1234
	c.SocketIP = "0.0.0.0"

	//Mode
	c.PadsSticksMode = c.toPadsSticksModeCfg()

	//c.mouseOnRightStickPad = c.toBoolConfig("mouseOnRightStickPad")

	//Pads/Stick
	switch c.ControllerInUse {
	case SteamController:
		c.StickAngleMarginSC = c.toIntConfig("StickAngleMargin")
		c.StickThresholdSC = c.toPctConfig("StickThresholdPct")
		c.StickEdgeThresholdSC = c.toPctConfig("StickEdgeThresholdPct")

		c.StickBoundariesMapSC = genEqualThresholdBoundariesMap(false,
			makeAngleMargin(0, c.StickAngleMarginSC, c.StickAngleMarginSC),
			c.StickThresholdSC,
			c.StickEdgeThresholdSC)

	case DualShock:
		c.StickDeadzoneDS = c.toFloatConfig("StickDeadzone")
	}

	//commands
	c.holdRefreshInterval = 15 * time.Millisecond
	c.TriggerThreshold = c.toPctConfig("TriggerThresholdPct")
	c.holdingThreshold = c.toMillisConfig("holdingThresholdMs")

	//mouse
	c.mouseInterval = c.toMillisConfig("mouseIntervalMs")
	c.mouseSpeed = c.toFloatConfig("mouseSpeed")
	c.mouseEdgeThreshold = c.toPctConfig("mouseEdgeThresholdPct")

	//scroll
	c.scrollFastestInterval = c.toIntToFloatConfig("scrollFastestIntervalMs")
	c.scrollSlowestInterval = c.toIntToFloatConfig("scrollSlowestIntervalMs")

	c.horizontalScrollThreshold = c.toPctConfig("horizontalScrollThresholdPct")

	//typing
	c.TypingStraightAngleMargin = c.toIntConfig("TypingStraightAngleMargin")
	c.TypingDiagonalAngleMargin = c.toIntConfig("TypingDiagonalAngleMargin")
	c.TypingThreshold = c.toPctConfig("TypingThresholdPct")
}

var Cfg *ConfigsT

type ConfigsT struct {
	mouseOnRightStickPad bool

	mousePadStick, scrollPadStick *PadStickPosition

	// Math
	OutputMin            float64
	MinStandardPadRadius float64

	// Mode
	RunFromTerminal bool
	ControllerInUse ControllerInUseT
	PadsSticksMode  *PadsSticksMode

	// commands
	holdRefreshInterval time.Duration
	TriggerThreshold    float64
	holdingThreshold    time.Duration

	// mouse
	mouseInterval      time.Duration
	mouseSpeed         float64
	mouseEdgeThreshold float64

	StickAngleMarginSC                     int
	StickThresholdSC, StickEdgeThresholdSC float64

	StickBoundariesMapSC ZoneBoundariesMap

	StickDeadzoneDS float64

	// scroll
	scrollFastestInterval, scrollSlowestInterval float64
	horizontalScrollThreshold                    float64

	// typing
	TypingStraightAngleMargin, TypingDiagonalAngleMargin int
	TypingThreshold                                      float64

	// path
	BaseDir, LayoutsDir string
	LayoutInUse         string
	RawStrConfigs       map[string]string

	// web socket
	SocketPort int
	SocketIP   string
}
