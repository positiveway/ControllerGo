package mainLogic

import (
	"ControllerGo/src/osSpecific"
	"fmt"
	"math"
	"time"
)

var mouseMovement = Coords{}
var scrollMovement = Coords{}

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

func printForce(force float64, prefix string) {
	if force != 0.0 {
		fmt.Printf("%s: %0.3f\n", prefix, force)
	}
}

func printPair[T Number](_x, _y T, prefix string) {
	x, y := float64(_x), float64(_y)
	fmt.Printf("%s: %0.2f %0.2f\n", prefix, x, y)
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
			//fmt.Printf("%v %v\n", xForce, yForce)
			osSpecific.MoveMouse(xForce, yForce)
		}

		time.Sleep(mouseInterval)
	}
}

func calcScrollInterval(input float64) time.Duration {
	return calcRefreshInterval(input, scrollSlowestInterval, scrollFastestInterval)
}

func getDirection(val float64, horizontal bool) int32 {
	if horizontal && math.Abs(val) < horizontalScrollThreshold {
		return 0
	}
	switch {
	case val == 0:
		return 0
	case val > 0:
		return -1
	case val < 0:
		return 1
	}
	panic("direction error")
}

func getDirections(x, y float64) (int32, int32) {
	hDir := getDirection(x, true)
	vDir := getDirection(y, false)
	//hDir *= -1

	if hDir != 0 {
		vDir = 0
	}
	return hDir, vDir
}

func RunScrollThread() {
	var hDir, vDir int32
	for {
		scrollMovement.updateValues()
		hDir, vDir = getDirections(scrollMovement.x, scrollMovement.y)

		scrollInterval := time.Duration(scrollFastestInterval) * time.Millisecond
		if scrollMovement.magnitude != 0 {
			scrollInterval = calcScrollInterval(scrollMovement.magnitude)
		}

		if hDir != 0 {
			osSpecific.ScrollHorizontal(hDir)
		}
		if vDir != 0 {
			osSpecific.ScrollVertical(vDir)
		}

		time.Sleep(scrollInterval)
	}
}
