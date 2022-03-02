package mainLogic

import (
	"ControllerGo/src/platformSpecific"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
)

const RunFromTerminal = false

func InitPath() {
	if RunFromTerminal {
		BaseDir = filepath.Dir(filepath.Dir(platformSpecific.GetCurFileDir()))
	} else {
		BaseDir = platformSpecific.DefaultProjectDir
	}
	LayoutsDir = filepath.Join(BaseDir, "layouts")
}

func setLayoutDir() {
	LayoutInUse = ReadFile(filepath.Join(LayoutsDir, "layout_to_use.txt"))
	LayoutInUse = strings.TrimSpace(LayoutInUse)

	curLayoutDir := path.Join(LayoutsDir, LayoutInUse)
	if _, err := os.Stat(curLayoutDir); os.IsNotExist(err) {
		panicMsg("Layout folder with such name doesn't exist: %s", LayoutInUse)
	}
}

func InitSettings() {
	InitPath()
	setLayoutDir()
	setConfigVars()
	initCodeMapping()
	joystickTyping = makeJoystickTyping()
	commandsLayout = loadCommandsLayout()
	boundariesMap = genBoundariesMap()
}

var found bool
var err error

func RunMain() {
	InitSettings()
	platformSpecific.InitResources()

	platformSpecific.InitInput()
	defer platformSpecific.CloseInputResources()
	defer releaseAll()

	go RunMouseMoveThread()
	go RunScrollThread()
	go RunReleaseHoldThread()

	if GamesModeOn {
		RunMovementThread()
	}

	runtime.GC()
	debug.SetGCPercent(-1)

	RunWebSocket()
}
