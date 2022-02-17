package mainLogic

import (
	"ControllerGo/src/osSpecific"
	"path/filepath"
)

func setLayoutDir(layoutName string) {
	LayoutDir = filepath.Join(BaseDir, "Layouts", layoutName)
}

func loadConfigs() {
	InitPath()
	setLayoutDir("Linux")
	convertLetterToCodeMapping()
	joystickTyping = makeJoystickTyping()
	commandsLayout = loadCommandsLayout()
	boundariesMap = genBoundariesMap()
}

func InitSettings() {
	loadConfigs()
}

const RunFromTerminal = true

func RunMain() {
	InitSettings()
	osSpecific.RunOsLogic()

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
