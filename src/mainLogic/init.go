package mainLogic

import (
	"ControllerGo/src/osSpecific"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const RunFromTerminal = false

func InitPath() {
	if RunFromTerminal {
		BaseDir = filepath.Dir(filepath.Dir(osSpecific.GetCurFileDir()))
	} else {
		BaseDir = osSpecific.DefaultProjectDir
	}
	LayoutsDir = filepath.Join(BaseDir, "layouts")
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

	osSpecific.InitInput()
	defer osSpecific.CloseInputResources()

	go RunMouseMoveThread()
	go RunScrollThread()
	go RunReleaseHoldThread()

	RunWebSocket()
}
