//go:build !windows

package main

import "github.com/bendahl/uinput"

const Linux = "Linux"
const Windows = "Windows"

var mouse uinput.Mouse
var keyboard uinput.Keyboard

func pressKeyOrMouse(key int) {
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

func releaseKeyOrMouse(key int) {
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
