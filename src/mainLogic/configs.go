package mainLogic

import (
	"github.com/positiveway/gofuncs"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type ControllerInUseT struct {
	Steam, DS bool
}

func MakeControllerInUse(isSteamInUse bool) *ControllerInUseT {
	return &ControllerInUseT{Steam: isSteamInUse, DS: !isSteamInUse}
}

// Math
const OutputMin float64 = 0.0
const PadRadius = 1.2

// Mode
var ControllerInUse *ControllerInUseT = MakeControllerInUse(true)
var padsMode *PadsMode

// commands
var TriggerThreshold float64
var holdingThreshold time.Duration

// mouse
var mouseInterval time.Duration
var mouseSpeed float64
var mouseEdgeThreshold float64

// Pads/Stick
var PadsRotation int
var StickRotation int

var StickAngleMargin int
var StickThreshold float64
var StickEdgeThreshold float64

var StickBoundariesMap ZoneBoundariesMap

var StickDeadzone float64

// scroll
var scrollFastestInterval float64
var scrollSlowestInterval float64

var horizontalScrollThreshold float64

// typing
var TypingStraightAngleMargin int
var TypingDiagonalAngleMargin int
var TypingThreshold float64

// common
var DefaultRefreshInterval time.Duration

// web socket
const SocketPort int = 1234
const SocketIP string = "0.0.0.0"

// path
var BaseDir string
var LayoutsDir string
var LayoutInUse string
var Configs = map[string]string{}

func ReadLayoutFile(pathFromLayoutsDir string, skipLines int) [][]string {
	file := filepath.Join(LayoutsDir, pathFromLayoutsDir)
	lines := gofuncs.ReadLines(file)
	lines = lines[skipLines:]

	var linesParts [][]string
	for _, line := range lines {
		line = gofuncs.Strip(line)
		if gofuncs.IsEmptyStripStr(line) || gofuncs.StartsWithAnyOf(line, ";", "//") {
			continue
		}
		parts := gofuncs.SplitByAnyOf(line, "&|>:,=")
		for ind, part := range parts {
			parts[ind] = gofuncs.Strip(part)
		}
		linesParts = append(linesParts, parts)
	}
	return linesParts
}

func loadConfigs() {
	linesParts := ReadLayoutFile(path.Join(LayoutInUse, "configs.csv"), 0)
	for _, parts := range linesParts {
		constName := parts[0]
		constValue := parts[1]

		constName = strings.ToLower(constName)
		gofuncs.AssignWithDuplicateCheck(Configs, constName, constValue)
	}
}

func getConfig(constName string) string {
	constName = strings.ToLower(constName)
	return gofuncs.GetOrPanic(Configs, constName, "No such name in config")
}

func toBoolConfig(name string) bool {
	return gofuncs.StrToBool(getConfig(name))
}

func toIntConfig(name string) int {
	return gofuncs.StrToInt(getConfig(name))
}

func toMillisConfig(name string) time.Duration {
	return gofuncs.StrToMillis(getConfig(name))
}

func toFloatConfig(name string) float64 {
	return gofuncs.StrToFloat(getConfig(name))
}

func toIntToFloatConfig(name string) float64 {
	return gofuncs.StrToIntToFloat(getConfig(name))
}

func toPctConfig(name string) float64 {
	return gofuncs.StrToPct(getConfig(name))
}

func setConfigVars() {
	loadConfigs()

	//Mode
	padsMode = MakePadsMode(toIntConfig("PadsMode"))

	//commands
	TriggerThreshold = toFloatConfig("TriggerThreshold")
	holdingThreshold = toMillisConfig("holdingThreshold")

	//mouse
	mouseInterval = toMillisConfig("mouseInterval")
	mouseSpeed = toFloatConfig("mouseSpeed")
	mouseEdgeThreshold = toFloatConfig("mouseEdgeThreshold")

	//Pads/Stick
	PadsRotation = toIntConfig("PadsRotation")
	StickRotation = toIntConfig("StickRotation")

	StickAngleMargin = toIntConfig("StickAngleMargin")
	StickThreshold = toPctConfig("StickThresholdPct")
	StickEdgeThreshold = toPctConfig("StickEdgeThresholdPct")

	StickBoundariesMap = genEqualThresholdBoundariesMap(false,
		makeAngleMargin(0, StickAngleMargin, StickAngleMargin),
		StickThreshold,
		StickEdgeThreshold)

	StickDeadzone = toFloatConfig("StickDeadzone")

	//scroll
	scrollFastestInterval = toIntToFloatConfig("scrollFastestInterval")
	scrollSlowestInterval = toIntToFloatConfig("scrollSlowestInterval")

	horizontalScrollThreshold = toFloatConfig("horizontalScrollThreshold")

	//typing
	TypingStraightAngleMargin = toIntConfig("TypingStraightAngleMargin")
	TypingDiagonalAngleMargin = toIntConfig("TypingDiagonalAngleMargin")
	TypingThreshold = toPctConfig("TypingThresholdPct")

	//common
	DefaultRefreshInterval = toMillisConfig("DefaultRefreshInterval")
}
