package mainLogic

import (
	"ControllerGo/src/osSpecific"
	"fmt"
	"math"
	"time"
)

var mouseMovement = Coords{}
var scrollMovement = Coords{}

func applyPower(force *float64) {
	sign := getSignMakeAbs(force)
	*force = math.Pow(*force, forcePower)
	applySign(sign, force)
}

func mouseForce(val float64, magnitude float64) int32 {
	force := convertRange(val, mouseMaxMove)
	//printForce(force, "before")
	applyPower(&force)
	//if magnitude >= MaxAccelRadiusThreshold {
	//	force *= MaxAccelMultiplier
	//}
	//printForce(force, "after")
	return int32(force)
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

func calcForces() (int32, int32) {
	//coordsMetrics := mouseMovement.getMetrics()
	//coordsMetrics.correctValuesNearRadius()
	x, y := mouseMovement.getValues()
	magnitude := calcMagnitude(x, y)

	xForce := mouseForce(x, magnitude)
	yForce := -mouseForce(y, magnitude)

	//if x != 0.0 || y != 0.0{
	//	printPair(x, y, "x, y")
	//	printPair(xForce, yForce, "force")
	//	fmt.Println()
	//}
	return xForce, yForce
}

func RunMouseMoveThread() {
	for {
		xForce, yForce := calcForces()
		if (xForce != 0) || (yForce != 0) {
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
	case val == 0.0:
		return 0
	case val > 0:
		return -1
	case val < 0:
		return 1
	}
	panic("direction error")
}

func getDirections(x, y float64) (int32, int32) {
	hDir, vDir := getDirection(x, true), getDirection(y, false)
	//hDir *= -1

	if hDir != 0 {
		vDir = 0
	}
	return hDir, vDir
}

func RunScrollThread() {
	for {
		x, y := scrollMovement.getValues()
		magnitude := calcMagnitude(x, y)
		hDir, vDir := getDirections(x, y)

		scrollInterval := time.Duration(scrollFastestInterval) * time.Millisecond
		if magnitude != 0 {
			scrollInterval = calcScrollInterval(magnitude)
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
