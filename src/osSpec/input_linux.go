//go:build linux

package osSpec

import (
	"github.com/bendahl/uinput"
	"github.com/positiveway/gofuncs"
)

const DefaultProjectDir string = "/home/user/GolandProjects/ControllerGo"

var mouse uinput.Mouse
var keyboard uinput.Keyboard

func CloseInputResources() {
	if keyboard != nil {
		keyboard.Close()
	}
	if mouse != nil {
		mouse.Close()
	}
}

func InitInput() {
	var err error

	// initialize keyboard and check for possible errors
	keyboard, err = uinput.CreateKeyboard("/dev/uinput", []byte("testkeyboard"))
	gofuncs.CheckErr(err)

	// initialize mouse and check for possible errors
	mouse, err = uinput.CreateMouse("/dev/uinput", []byte("testmouse"))
	gofuncs.CheckErr(err)
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

func MoveMouse(x, y int) {
	mouse.Move(int32(x), int32(-y))
}

func ScrollHorizontal(direction int) {
	mouse.Wheel(true, int32(direction))
}

func ScrollVertical(direction int) {
	mouse.Wheel(false, int32(direction))
}
