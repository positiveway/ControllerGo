//go:build linux

package osSpecific

import (
	"github.com/bendahl/uinput"
)

var mouse uinput.Mouse
var keyboard uinput.Keyboard

func CloseInputResources() {
	mouse.Close()
	keyboard.Close()
}

func InitInput() {
	var err error
	// initialize mouse and check for possible errors
	mouse, err = uinput.CreateMouse("/dev/uinput", []byte("testmouse"))
	CheckErr(err)
	// always do this after the initialization in order to guarantee that the device will be properly closed

	// initialize keyboard and check for possible errors
	keyboard, err = uinput.CreateKeyboard("/dev/uinput", []byte("testkeyboard"))
	if err != nil {
		return
	}
	// always do this after the initialization in order to guarantee that the device will be properly closed
}

func PressKeyOrMouse(key int) {
	switch key {
	case LeftMouse:
		mouse.LeftPress()
	case RightMouse:
		mouse.RightPress()
	case MiddleMouse:
		mouse.MiddlePress()
	default:
		keyboard.KeyDown(key)
	}
}

func ReleaseKeyOrMouse(key int) {
	switch key {
	case LeftMouse:
		mouse.LeftRelease()
	case RightMouse:
		mouse.RightRelease()
	case MiddleMouse:
		mouse.MiddleRelease()
	default:
		keyboard.KeyUp(key)
	}
}

func TypeKey(key int) {
	keyboard.KeyPress(key)
}

func MoveMouse(x, y int32) {
	mouse.Move(x, y)
}

func ScrollHorizontal(direction int32) {
	mouse.Wheel(true, direction)
}

func ScrollVertical(direction int32) {
	mouse.Wheel(false, direction)
}
