package mainLogic

import (
	"ControllerGo/src/platformSpecific"
	"math"
	"time"
)

func DsEventChanged() {
	switch GamesModeOn {
	case false:
		switch typingMode.mode {
		case false:
			switch event.btnOrAxis {
			case AxisLeftStickX:
				mouseMovement.setX(event.value)
			case AxisLeftStickY:
				mouseMovement.setY(event.value)
			case AxisRightStickX:
				scrollMovement.setX(event.value)
			case AxisRightStickY:
				scrollMovement.setY(event.value)
			}
		case true:
			switch event.btnOrAxis {
			case AxisLeftStickX:
				joystickTyping.leftCoords.setDirectlyX(event.value)
				joystickTyping.updateLeftZone()
			case AxisLeftStickY:
				joystickTyping.leftCoords.setDirectlyY(event.value)
				joystickTyping.updateLeftZone()
			case AxisRightStickX:
				joystickTyping.rightCoords.setDirectlyX(event.value)
				joystickTyping.updateRightZone()
			case AxisRightStickY:
				joystickTyping.rightCoords.setDirectlyY(event.value)
				joystickTyping.updateRightZone()
			}
		}
	case true:
		switch event.btnOrAxis {
		case AxisLeftStickX:
			movementCoords.setX(event.value)
		case AxisLeftStickY:
			movementCoords.setY(event.value)
		case AxisRightStickX:
			mouseMovement.setX(event.value)
		case AxisRightStickY:
			mouseMovement.setY(event.value)
		}
	}
}

var mouseMovement = Coords{}

func applyPower(force float64) float64 {
	sign, force := getSignAndAbs(force)
	force = math.Pow(force, forcePower)
	return applySign(sign, force)
}

func mouseForce(input float64) int32 {
	force := convertRange(input, mouseMaxMove)
	//printForce(force, "before")
	force = applyPower(force)
	//if magnitude >= MaxAccelRadiusThreshold {
	//	force *= MaxAccelMultiplier
	//}
	//printForce(force, "after")
	return floatToInt32(force)
}

func RunMouseMoveThread() {
	var xForce, yForce int32
	for {
		//coordsMetrics := mouseMovement.getMetrics()
		//coordsMetrics.correctValuesNearRadius()
		mouseMovement.updateValues()

		xForce = mouseForce(mouseMovement.x)
		yForce = -mouseForce(mouseMovement.y)

		//if x != 0.0 || y != 0.0{
		//	printPair(x, y, "x, y")
		//	printPair(xForce, yForce, "force")
		//	fmt.Println()
		//}

		if xForce != 0 || yForce != 0 {
			//print("%v %v", xForce, yForce)
			platformSpecific.MoveMouse(xForce, yForce)
		}

		time.Sleep(mouseInterval)
	}
}
