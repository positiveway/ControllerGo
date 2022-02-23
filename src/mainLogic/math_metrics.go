package mainLogic

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type Coords struct {
	_x, _y float64
	mu     sync.Mutex
}

func (coords *Coords) reset() {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	coords._x = 0
	coords._y = 0
}

func (coords *Coords) setX(x float64) {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	coords._x = x
}

func (coords *Coords) setY(y float64) {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	coords._y = y
}

func (coords *Coords) setValues(x, y float64) {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	coords._x = x
	coords._y = y
}

func (coords *Coords) getRawValues() (float64, float64) {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	return coords._x, coords._y
}

func (coords *Coords) getNormalizedValues() (float64, float64) {
	x, y := coords.getRawValues()

	applyDeadzones(&x, &y)
	normalizeIncorrectEdgeValues(&x, &y)

	return x, y
}

func (coords *Coords) getMetrics() CoordsMetrics {
	x, y := coords.getNormalizedValues()
	magnitude := calcMagnitude(x, y)
	resolvedAngle := calcResolvedAngle(x, y)
	oneQuarterAngle := calcOneQuarterAngle(resolvedAngle)

	//mappedX, mappedY := mapCircleToSquare(x, y)

	if magnitude != 0 {
		//fmt.Printf("x: %0.2f, y: %0.2f, mappedX: %0.2f, mappedY: %0.2f, magn: %0.2f\n", coordsMetrics.x, coordsMetrics.y, coordsMetrics.mappedX, coordsMetrics.mappedY, coordsMetrics.magnitude)
	}

	return CoordsMetrics{
		x:               x,
		y:               y,
		magnitude:       magnitude,
		angle:           resolvedAngle,
		oneQuarterAngle: oneQuarterAngle,
	}
}

type CoordsMetrics struct {
	x, y                   float64
	mappedX, mappedY       float64
	magnitude              float64
	angle, oneQuarterAngle int
}

func setToRadiusValue(val *float64) {
	*val = math.Copysign(1.0, *val)
}

func (metrics *CoordsMetrics) correctValuesNearRadius() {
	if metrics.magnitude > MaxAccelRadiusThreshold {
		if MaxAccelMinAngle < metrics.oneQuarterAngle && metrics.oneQuarterAngle < MaxAccelMaxAngle {
			setToRadiusValue(&metrics.x)
			setToRadiusValue(&metrics.y)
		}
	}
}

func initMaxAccelValues() {
	if MaxAccelAngleMargin > 45 {
		panic(fmt.Sprintf("Incorrect value of \"MaxAccelAngleMargin\": %v\n", MaxAccelAngleMargin))
	}
	MaxAccelMinAngle = 45 - MaxAccelAngleMargin
	MaxAccelMaxAngle = 45 + MaxAccelAngleMargin
}

func applyDeadzone(value *float64) {
	if math.Abs(*value) < Deadzone {
		*value = 0.0
	}
}

func applyDeadzones(x, y *float64) {
	applyDeadzone(x)
	applyDeadzone(y)
}

const twoSqrt2 float64 = 2.0 * math.Sqrt2

func mapCircleToSquare(u, v float64) (float64, float64) {
	u2 := u * u
	v2 := v * v
	subtermx := 2.0 + u2 - v2
	subtermy := 2.0 + v2 - u2
	u *= twoSqrt2
	v *= twoSqrt2
	x := 0.5 * (math.Sqrt(subtermx+u) - math.Sqrt(subtermx-u))
	y := 0.5 * (math.Sqrt(subtermy+v) - math.Sqrt(subtermy-v))
	normalizeIncorrectEdgeValues(&x, &y)
	return x, y
}

func calcOneQuarterAngle(resolvedAngle int) int {
	return int(math.Mod(float64(resolvedAngle), 90))
}

func resolveAngle(angle float64) int {
	angle = math.Mod(angle+360, 360)
	return int(angle)
}

const radiansMultiplier float64 = 180 / math.Pi

func calcResolvedAngle(x, y float64) int {
	degrees := math.Atan2(y, x) * radiansMultiplier
	return resolveAngle(degrees)
}

func calcMagnitude(x, y float64) float64 {
	return math.Hypot(x, y)
}

func normalizeIncorrectEdgeValues(x, y *float64) {
	magnitude := calcMagnitude(*x, *y)
	if magnitude > 1.0 {
		*x /= magnitude
		*y /= magnitude
	}
}

func normalizeCoords(x, y *float64, magnitude float64) {
	*x *= magnitude
	*y *= magnitude
}

const outputRangeMin float64 = 1.0

func convertRange(input, outputMax float64) float64 {
	sign := getSignMakeAbs(&input)

	if input == 0.0 {
		return 0.0
	}

	if input > 1.0 {
		panic(fmt.Sprintf("Axis input value is greater than 1.0. Current value: %v\n", input))
	}

	output := outputRangeMin + ((outputMax-outputRangeMin)/inputRange)*(input-Deadzone)
	applySign(sign, &output)
	return output
}

func calcRefreshInterval(input, intervalRange, slowestInterval float64) time.Duration {
	input = math.Abs(input)
	refreshInterval := convertRange(input, intervalRange)
	refreshInterval = slowestInterval - math.Round(refreshInterval)
	return time.Duration(refreshInterval) * time.Millisecond
}
