package main

import (
	"github.com/bendahl/uinput"
	"path/filepath"
)

var mouse uinput.Mouse
var keyboard uinput.Keyboard

func setLayoutDir(layoutName string) {
	layoutDir = filepath.Join(BaseDir, "Layouts", layoutName)
}

func initPath() {
	BaseDir = filepath.Dir(getCurFileDir())
	EventServerExecPath = filepath.Join(BaseDir, "Build", "ControllerRust")
	getLocaleExecPath = getCurFileDir() + "/getLocale.sh"
}

func loadConfigs() {
	initPath()
	setLayoutDir("Linux")
	convertLetterToCodeMapping()
	joystickTyping = makeJoystickTyping()
	commandsLayout = loadCommandsLayout()
	boundariesMap = genBoundariesMap()
}

func main() {
	loadConfigs()
	setSelfPriority()

	startEventServer()
	defer killEventServer()

	var err error
	// initialize mouse and check for possible errors
	mouse, err = uinput.CreateMouse("/dev/uinput", []byte("testmouse"))
	check_err(err)
	// always do this after the initialization in order to guarantee that the device will be properly closed
	defer mouse.Close()

	// initialize keyboard and check for possible errors
	keyboard, err = uinput.CreateKeyboard("/dev/uinput", []byte("testkeyboard"))
	if err != nil {
		return
	}
	// always do this after the initialization in order to guarantee that the device will be properly closed
	defer keyboard.Close()

	go moveMouse()
	go scroll()

	mainWS()
}
