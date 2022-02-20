package mainLogic

import (
	"fmt"
	"path"
	"strings"
	"time"
)

//path
var BaseDir string
var LayoutsDir string
var LayoutInUse string
var EventServerExecPath string
var Configs = map[string]string{}

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
		panic(fmt.Sprintf("No such name in config %s\n", constName))
	}
}

func toIntConfig(name string) int {
	return ToInt(getConfig(name))
}

func toMillisecondsConfig(name string) time.Duration {
	return ToMilliseconds(getConfig(name))
}

func toFloatConfig(name string) float64 {
	return ToFloat(getConfig(name))
}

func toIntToFloatConfig(name string) float64 {
	return ToIntToFloat(getConfig(name))
}

func setConfigVars() {
	loadConfigs()

	//commands
	TriggerThreshold = toFloatConfig("TriggerThreshold")
	holdThreshold = toMillisecondsConfig("holdThreshold")

	//mouse
	mouseMaxMove = toIntToFloatConfig("mouseMaxMove")
	forcePower = toFloatConfig("forcePower")
	Deadzone = toFloatConfig("Deadzone")

	mouseInterval = toMillisecondsConfig("mouseInterval")

	//scroll
	scrollFastestInterval = toIntToFloatConfig("scrollFastestInterval")
	scrollSlowestInterval = toIntToFloatConfig("scrollSlowestInterval")
	scrollIntervalRange = scrollSlowestInterval - scrollFastestInterval

	horizontalScrollThreshold = toFloatConfig("horizontalScrollThreshold")

	//typing
	angleMargin = toIntConfig("angleMargin")
	magnitudeThresholdPct = toIntToFloatConfig("magnitudeThresholdPct")
	MagnitudeThreshold = magnitudeThresholdPct / 100

	//common
	DefaultWaitInterval = toMillisecondsConfig("DefaultWaitInterval")
}

//commands
var TriggerThreshold float64
var holdThreshold time.Duration

//mouse
var mouseMaxMove float64
var forcePower float64
var Deadzone float64

//var mouseScaleFactor float64 = 3
//var mouseIntervalInt int = int(math.Round(mouseMaxMove*mouseScaleFactor))
var mouseInterval time.Duration

//scroll
var scrollFastestInterval float64
var scrollSlowestInterval float64

var scrollIntervalRange float64
var horizontalScrollThreshold float64

//typing
var angleMargin int
var magnitudeThresholdPct float64
var MagnitudeThreshold float64

//common
var DefaultWaitInterval time.Duration

//web socket
const SocketPort int = 1234
const SocketIP string = "0.0.0.0"
