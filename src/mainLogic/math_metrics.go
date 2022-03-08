package mainLogic

import (
	"math"
	"sync"
	"time"
)

type Coords struct {
	x, y      float64
	magnitude float64
	angle     int
	mu        sync.Mutex
}

func makeCoords() *Coords {
	coords := Coords{}
	coords.reset()
	return &coords
}

func (coords *Coords) setDirectlyX() {
	coords.x = event.value
}

func (coords *Coords) setDirectlyY() {
	coords.y = event.value
}

func (coords *Coords) setX() {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	coords.setDirectlyX()
}

func (coords *Coords) setY() {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	coords.setDirectlyY()
}

func (coords *Coords) printCurState() {
	printPair(coords.x, coords.y, "(x, y): ")
}

func (coords *Coords) reset() {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	coords.x = math.NaN()
	coords.y = math.NaN()
}

func (coords *Coords) updateValues() {
	coords.x, coords.y, coords.magnitude = normalizeIncorrectEdgeValues(coords.x, coords.y)
}

func (coords *Coords) updateAngle() {
	coords.angle = calcResolvedAngle(coords.x, coords.y)
}

func resolveAngle(angle float64) int {
	angle = math.Mod(angle+360, 360)
	return floatToInt(angle)
}

const radiansMultiplier float64 = 180 / math.Pi

func calcResolvedAngle(x, y float64) int {
	if isNan(x, y) {
		return 0
	}
	angle := math.Atan2(y, x) * radiansMultiplier
	return resolveAngle(angle)
}

func calcDistance(x, y float64) float64 {
	if isNan(x, y) {
		return 0
	}
	return math.Hypot(x, y)
}

func normalizeIncorrectEdgeValues(x, y float64) (float64, float64, float64) {
	magnitude := calcDistance(x, y)
	if magnitude > 1.0 {
		x /= magnitude
		y /= magnitude
		magnitude = 1.0
	}
	return x, y, magnitude
}

const outputMin float64 = 0.0

func convertRange(input, outputMax float64) float64 {
	panicIsNan(input)

	if input == 0.0 {
		return 0.0
	}

	sign, input := getSignAndAbs(input)

	if input > 1.0 {
		panicMsg("Axis input value is greater than 1.0. Current value: %v", input)
	}

	output := outputMin + ((outputMax-outputMin)/inputRange)*(input-Deadzone)
	return applySign(sign, output)
}

func calcRefreshInterval(input, slowestInterval, fastestInterval float64) time.Duration {
	input = math.Abs(input)
	refreshInterval := convertRange(input, slowestInterval-fastestInterval)
	refreshInterval = slowestInterval - refreshInterval
	return time.Duration(floatToInt64(refreshInterval)) * time.Millisecond
}

func applyDeadzone(value float64) float64 {
	if isNan(value) {
		return value
	}
	if math.Abs(value) < Deadzone {
		value = 0.0
	}
	return value
}

func printPair[T Number](_x, _y T, prefix string) {
	x, y := float64(_x), float64(_y)
	print("%s: %0.2f %0.2f", prefix, x, y)
}

//func calcOneQuarterAngle(resolvedAngle int) int {
//	return floatToInt(math.Mod(float64(resolvedAngle), 90))
//}

//const twoSqrt2 float64 = 2.0 * math.Sqrt2

//func mapCircleToSquare(u, v float64) (float64, float64) {
//	u2 := u * u
//	v2 := v * v
//	subtermx := 2.0 + u2 - v2
//	subtermy := 2.0 + v2 - u2
//	u *= twoSqrt2
//	v *= twoSqrt2
//	x := 0.5 * (math.Sqrt(subtermx+u) - math.Sqrt(subtermx-u))
//	y := 0.5 * (math.Sqrt(subtermy+v) - math.Sqrt(subtermy-v))
//	normalizeIncorrectEdgeValues(&x, &y)
//	return x, y
//}
