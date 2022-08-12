package mainLogic

import (
	"github.com/positiveway/gofuncs"
	"time"
)

type ControllerInUseT struct {
	SteamController, DualShock, SteamDeck bool
}

func MakeControllerInUse(isSteamControllerInUse bool) *ControllerInUseT {
	return &ControllerInUseT{SteamController: isSteamControllerInUse, DualShock: !isSteamControllerInUse, SteamDeck: !isSteamControllerInUse}
}

func (c *ConfigsT) setConfigConstants() {
	//Debug
	gofuncs.PrintDebugInfo = false

	c.RunFromTerminal = false
}

func (c *ConfigsT) setConfigVars() {
	c.ControllerInUse = MakeControllerInUse(true)

	//Math
	c.OutputMin = 0.0
	c.PadRadius = 1.2

	// web socket
	c.SocketPort = 1234
	c.SocketIP = "0.0.0.0"

	//Mode
	c.padsMode = MakePadsMode(c.toIntConfig("PadsMode"))

	//commands
	c.holdRefreshInterval = 15 * time.Millisecond
	c.TriggerThreshold = c.toFloatConfig("TriggerThreshold")
	c.holdingThreshold = c.toMillisConfig("holdingThreshold")

	//mouse
	c.mouseInterval = c.toMillisConfig("mouseInterval")
	c.mouseSpeed = c.toFloatConfig("mouseSpeed")
	c.mouseEdgeThreshold = c.toFloatConfig("mouseEdgeThreshold")

	//Pads/Stick
	c.LeftPadRotation = c.toFloatConfig("LeftPadRotation")
	c.RightPadRotation = c.toFloatConfig("RightPadRotation")
	c.StickRotation = c.toFloatConfig("StickRotation")

	c.StickAngleMargin = c.toIntConfig("StickAngleMargin")
	c.StickThreshold = c.toPctConfig("StickThresholdPct")
	c.StickEdgeThreshold = c.toPctConfig("StickEdgeThresholdPct")

	c.StickBoundariesMap = genEqualThresholdBoundariesMap(false,
		makeAngleMargin(0, c.StickAngleMargin, c.StickAngleMargin),
		c.StickThreshold,
		c.StickEdgeThreshold)

	c.StickDeadzone = c.toFloatConfig("StickDeadzone")

	//scroll
	c.scrollFastestInterval = c.toIntToFloatConfig("scrollFastestInterval")
	c.scrollSlowestInterval = c.toIntToFloatConfig("scrollSlowestInterval")

	c.horizontalScrollThreshold = c.toFloatConfig("horizontalScrollThreshold")

	//typing
	c.TypingStraightAngleMargin = c.toIntConfig("TypingStraightAngleMargin")
	c.TypingDiagonalAngleMargin = c.toIntConfig("TypingDiagonalAngleMargin")
	c.TypingThreshold = c.toPctConfig("TypingThresholdPct")
}

var Cfg *ConfigsT

type ConfigsT struct {
	// Math
	OutputMin float64
	PadRadius float64

	// Mode
	RunFromTerminal bool
	ControllerInUse *ControllerInUseT
	padsMode        *PadsMode

	// commands
	holdRefreshInterval time.Duration
	TriggerThreshold    float64
	holdingThreshold    time.Duration

	// mouse
	mouseInterval      time.Duration
	mouseSpeed         float64
	mouseEdgeThreshold float64

	// Pads/Stick
	LeftPadRotation, RightPadRotation, StickRotation float64

	StickAngleMargin                   int
	StickThreshold, StickEdgeThreshold float64

	StickBoundariesMap ZoneBoundariesMap

	StickDeadzone float64

	// scroll
	scrollFastestInterval, scrollSlowestInterval float64

	horizontalScrollThreshold float64

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
