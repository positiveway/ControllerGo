package mainLogic

import (
	"path"
	"strings"
	"time"
)

func loadConfigs() {
	linesParts := ReadLayoutFile(path.Join(LayoutInUse, "configs.csv"), 0)
	for _, parts := range linesParts {
		constName := parts[0]
		constValue := parts[1]

		constName = strings.ToLower(constName)
		AssignWithDuplicateCheck(Configs, constName, constValue)
	}
}

func getConfig(constName string) string {
	constName = strings.ToLower(constName)
	if constValue, found := Configs[constName]; found {
		return constValue
	} else {
		panicMsg("No such name in config %s", constName)
	}
	return ""
}

func toBoolConfig(name string) bool {
	return strToBool(getConfig(name))
}

func toIntConfig(name string) int {
	return strToInt(getConfig(name))
}

func toMillisConfig(name string) time.Duration {
	return strToMillis(getConfig(name))
}

func toFloatConfig(name string) float64 {
	return strToFloat(getConfig(name))
}

func toIntToFloatConfig(name string) float64 {
	return strToIntToFloat(getConfig(name))
}

func setConfigVars() {
	loadConfigs()

	//games
	GamesModeOn = toBoolConfig("GamesModeOn")

	//commands
	TriggerThreshold = toFloatConfig("TriggerThreshold")
	holdThreshold = toMillisConfig("holdThreshold")

	//mouse

	Deadzone = toFloatConfig("Deadzone")
	inputRange = 1.0 - Deadzone

	//scroll
	scrollFastestInterval = toIntToFloatConfig("scrollFastestInterval")
	scrollSlowestInterval = toIntToFloatConfig("scrollSlowestInterval")

	horizontalScrollThreshold = toFloatConfig("horizontalScrollThreshold")

	//typing
	TypingStraightAngleMargin = toIntConfig("TypingStraightAngleMargin")
	TypingDiagonalAngleMargin = toIntConfig("TypingDiagonalAngleMargin")
	TypingThresholdPct = toIntToFloatConfig("TypingThresholdPct")
	TypingThreshold = TypingThresholdPct / 100

	//common
	DefaultRefreshInterval = toMillisConfig("DefaultRefreshInterval")
}

//games
var GamesModeOn bool

//commands
var TriggerThreshold float64
var holdThreshold time.Duration

//mouse
var Deadzone float64
var inputRange float64

//scroll
var scrollFastestInterval float64
var scrollSlowestInterval float64

var horizontalScrollThreshold float64

//typing
var TypingStraightAngleMargin int
var TypingDiagonalAngleMargin int
var TypingThresholdPct float64
var TypingThreshold float64

//common
var DefaultRefreshInterval time.Duration

//web socket
const SocketPort int = 1234
const SocketIP string = "0.0.0.0"

//path
var BaseDir string
var LayoutsDir string
var LayoutInUse string
var Configs = map[string]string{}
