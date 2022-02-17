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
	setLayoutDir(LayoutInUse)
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
