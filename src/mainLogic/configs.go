package mainLogic

import (
	"fmt"
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
		panic(fmt.Sprintf("No such name in config %s\n", constName))
	}
}

func toBoolConfig(name string) bool {
	return ToBool(getConfig(name))
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

	//games
	GamesModeOn = toBoolConfig("GamesModeOn")

	//commands
	TriggerThreshold = toFloatConfig("TriggerThreshold")
	holdThreshold = toMillisecondsConfig("holdThreshold")

	//mouse
	mouseMaxMove = toIntToFloatConfig("mouseMaxMove")
	forcePower = toFloatConfig("forcePower")

	MaxAccelMultiplier = toFloatConfig("MaxAccelMultiplier")
	MaxAccelRadiusThreshold = toFloatConfig("MaxAccelRadiusThreshold")
	MaxAccelAngleMargin = toIntConfig("MaxAccelAngleMargin")
	initMaxAccelValues()

	Deadzone = toFloatConfig("Deadzone")
	inputRange = 1.0 - Deadzone

	mouseInterval = toMillisecondsConfig("mouseInterval")

	//scroll
	scrollFastestInterval = toIntToFloatConfig("scrollFastestInterval")
	scrollSlowestInterval = toIntToFloatConfig("scrollSlowestInterval")
	scrollIntervalRange = scrollSlowestInterval - scrollFastestInterval

	horizontalScrollThreshold = toFloatConfig("horizontalScrollThreshold")

	//typing
	RightAngleMargin = toIntConfig("RightAngleMargin")
	DiagonalAngleMargin = toIntConfig("DiagonalAngleMargin")
	magnitudeThresholdPct = toIntToFloatConfig("magnitudeThresholdPct")
	MagnitudeThreshold = magnitudeThresholdPct / 100

	//common
	DefaultRefreshInterval = toMillisecondsConfig("DefaultRefreshInterval")
}

//games
var GamesModeOn bool

//commands
var TriggerThreshold float64
var holdThreshold time.Duration

//mouse
var mouseMaxMove float64
var forcePower float64
var Deadzone float64
var inputRange float64
var MaxAccelRadiusThreshold float64
var MaxAccelAngleMargin int
var MaxAccelMinAngle, MaxAccelMaxAngle int
var MaxAccelMultiplier float64

//var mouseScaleFactor float64 = 3
//var mouseIntervalInt int = int(math.Round(mouseMaxMove*mouseScaleFactor))
var mouseInterval time.Duration

//scroll
var scrollFastestInterval float64
var scrollSlowestInterval float64

var scrollIntervalRange float64
var horizontalScrollThreshold float64

//typing
var RightAngleMargin int
var DiagonalAngleMargin int
var magnitudeThresholdPct float64
var MagnitudeThreshold float64

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
