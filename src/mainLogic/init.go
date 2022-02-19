package mainLogic

import (
	"ControllerGo/src/osSpecific"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

const RunFromTerminal = true

func InitPath() {
	if RunFromTerminal {
		BaseDir = filepath.Dir(filepath.Dir(osSpecific.GetCurFileDir()))
	} else {
		BaseDir = osSpecific.DefaultProjectDir
	}
	EventServerExecPath = filepath.Join(BaseDir, "Build", runtime.GOOS, "ControllerRust")
	LayoutsDir = filepath.Join(BaseDir, "Layouts")
}

func setLayoutDir() {
	LayoutInUse = ReadFile(filepath.Join(LayoutsDir, "layout_to_use.txt"))
	LayoutInUse = strings.TrimSpace(LayoutInUse)

	curLayoutDir := path.Join(LayoutsDir, LayoutInUse)
	if _, err := os.Stat(curLayoutDir); os.IsNotExist(err) {
		panic(fmt.Sprintf("Layout folder with such name doesn't exist: %s\n", LayoutInUse))
	}
}

func InitSettings() {
	InitPath()
	setLayoutDir()
	setConfigVars()
	convertLetterToCodeMapping()
	joystickTyping = makeJoystickTyping()
	commandsLayout = loadCommandsLayout()
	boundariesMap = genBoundariesMap()
}

func RunMain() {
	InitSettings()
	osSpecific.InitResources()

	SetSelfPriority()

	if RunFromTerminal {
		StartEventServer()
		defer KillEventServer()
	}

	osSpecific.InitInput()
	defer osSpecific.CloseInputResources()

	go RunMouseMoveThread()
	go RunScrollThread()
	go RunReleaseHoldThread()

	RunWebSocket()
}
