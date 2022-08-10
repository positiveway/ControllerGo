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
const (
	OutputMin float64 = 0.0
	PadRadius         = 1.2
)

// Mode
var (
	ControllerInUse *ControllerInUseT = MakeControllerInUse(true)
	padsMode        *PadsMode

	// commands
	TriggerThreshold float64
	holdingThreshold time.Duration

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
	Configs             = map[string]string{}
)

// web socket
const (
	SocketPort int    = 1234
	SocketIP   string = "0.0.0.0"
)

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
	LeftPadRotation = toFloatConfig("LeftPadRotation")
	RightPadRotation = toFloatConfig("RightPadRotation")
	StickRotation = toFloatConfig("StickRotation")

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
}

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
