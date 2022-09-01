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

func (c *ConfigsT) toPadsSticksModeCfg() *PadsSticksModeT {
	allModes := []ModeT{MouseMode, GamingMode}
	modeType := ModeT(c.getConfig("PadsSticksMode"))
	checkEnumCfg(allModes, modeType)
	return MakePadsSticksMode(modeType)
}

func (c *ConfigsT) setConfigConstants() {
	//Debug
	gofuncs.PrintDebugInfo = false

	c.RunFromTerminal = false
}

func (c *ConfigsT) initTouchpads() {
	switch c.ControllerInUse {
	case SteamController:
		RightPadStick = MakePadPosition(c.toFloatConfig("RightPadRotation"))
		LeftStick = MakePadPosition(c.toFloatConfig("StickRotation"))
		LeftPad = MakePadPosition(c.toFloatConfig("LeftPadRotation"))

		c.mousePadStick = RightPadStick
		c.scrollPadStick = LeftPad

		c.LeftTypingPS = LeftPad
		c.RightTypingPS = RightPadStick
	case DualShock:
		RightPadStick = MakePadPosition(c.toFloatConfig("RightStickRotation"))
		LeftStick = MakePadPosition(c.toFloatConfig("LeftStickRotation"))

		c.mousePadStick = RightPadStick
		c.scrollPadStick = LeftStick

		c.LeftTypingPS = LeftStick
		c.RightTypingPS = RightPadStick
	}
}

const FloatEqualityMargin = 0.000000000000001

func (c *ConfigsT) setConfigVars() {
	c.tickerInterval = 1
	c.GCPercent = 10000

	c.ControllerInUse = c.toControllerCfg()

	//Math
	c.OutputMin = 0
	c.MinStandardPadRadius = 1.0

	// web socket
	c.SocketPort = 1234
	c.SocketIP = "0.0.0.0"

	//Mode
	c.PadsSticksMode = c.toPadsSticksModeCfg()
	c.HighPrecisionMode = MakeHighPrecisionMode()

	//mouse/scroll
	//c.mouseOnRightStickPad = c.toBoolConfig("mouseOnRightStickPad")

	//Pads/Stick
	switch c.ControllerInUse {
	case SteamController:
		//mouse
		c.mouseIntervalSC = c.toMillisConfig("mouseIntervalMs")
		c.mouseSpeedSC = c.toFloatConfig("mouseSpeed")

		//stick
		stickAngleMarginSC := c.toIntConfig("StickAngleMargin")
		stickThresholdSC := c.toPctConfig("StickThresholdPct")
		stickEdgeThresholdSC := c.toPctConfig("StickEdgeThresholdPct")

		//init Stick map
		c.StickBoundariesMapSC = genEqualThresholdBoundariesMap(false,
			makeAngleMargin(0, stickAngleMarginSC, stickAngleMarginSC),
			stickThresholdSC,
			stickEdgeThresholdSC)

	case DualShock:
		//mouse
		c.mouseIntervalsDS = MakeIntervalRange(
			c.toIntToFloatConfig("mouseSlowestIntervalMs"),
			c.toIntToFloatConfig("mouseFastestIntervalMs"))

		//stick
		c.StickDeadzoneDS = c.toFloatConfig("StickDeadzone")
	}

	switch c.PadsSticksMode.GetMode() {
	case GamingMode:
		c.gamingMoveIntervals = MakeIntervalRange(
			c.toIntToFloatConfig("gamingMoveSlowestMs"),
			c.toIntToFloatConfig("gamingMoveFastestMs"),
		)
	}

	//commands
	c.TriggerThreshold = c.toPctConfig("TriggerThresholdPct")
	c.holdRepeatInterval = 40
	c.holdingStateThreshold = c.toFloatConfig("holdingStateThresholdMs")

	//gaming

	//mouse
	c.mouseEdgeThreshold = c.toPctConfig("mouseEdgeThresholdPct")

	//scroll
	c.scrollIntervals = MakeIntervalRange(
		c.toIntToFloatConfig("scrollSlowestIntervalMs"),
		c.toIntToFloatConfig("scrollFastestIntervalMs"))

	c.scrollHorizontalThreshold = c.toPctConfig("scrollHorizontalThresholdPct")

	//typing
	c.TypingStraightAngleMargin = c.toIntConfig("TypingStraightAngleMargin")
	c.TypingDiagonalAngleMargin = c.toIntConfig("TypingDiagonalAngleMargin")
	c.TypingThreshold = c.toPctConfig("TypingThresholdPct")
}

var Cfg *ConfigsT

type ConfigsT struct {
	tickerInterval float64
	GCPercent      int

	mouseOnRightStickPad bool

	mousePadStick, scrollPadStick *PadStickPositionT
	LeftTypingPS, RightTypingPS   *PadStickPositionT

	// Math
	OutputMin            float64
	MinStandardPadRadius float64

	// Mode
	RunFromTerminal bool
	ControllerInUse ControllerInUseT

	PadsSticksMode    *PadsSticksModeT
	HighPrecisionMode *HighPrecisionModeT

	// commands
	holdRepeatInterval, holdingStateThreshold float64
	TriggerThreshold                          float64

	//games
	gamingMoveIntervals *IntervalRangeT

	// mouse
	mouseIntervalsDS *IntervalRangeT

	mouseIntervalSC time.Duration
	mouseSpeedSC    float64

	mouseEdgeThreshold float64

	// scroll
	scrollIntervals           *IntervalRangeT
	scrollHorizontalThreshold float64

	//stick
	StickBoundariesMapSC ZoneBoundariesMapT

	StickDeadzoneDS float64

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
