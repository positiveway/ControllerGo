package mainLogic

import (
	"math"
	"sync"
	"time"
)

type TimePassed struct {
	value time.Duration
}

func (t *TimePassed) passedInterval(interval time.Duration) bool {
	t.value += DefaultRefreshInterval
	if t.value >= interval {
		t.value = 0
		return true
	}
	return false
}

type Coords struct {
	x, y           float64
	forceX, forceY float64
	magnitude      float64
	mu             sync.Mutex
}

func (coords *Coords) setDirectlyX(x float64) {
	coords.x = x
}

func (coords *Coords) setDirectlyY(y float64) {
	coords.y = y
}

func (coords *Coords) setX(x float64) {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	coords.x = x
}

func (coords *Coords) setY(y float64) {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	coords.y = y
}

func (coords *Coords) reset() {
	coords.mu.Lock()
	defer coords.mu.Unlock()
	coords.x = 0
	coords.y = 0
}

func (coords *Coords) normalizeValues() {
	applyDeadzone(&coords.x)
	applyDeadzone(&coords.y)
	normalizeIncorrectEdgeValues(&coords.x, &coords.y, &coords.magnitude)
}

func (coords *Coords) getValues() (float64, float64, float64) {
	coords.mu.Lock()
	defer coords.mu.Unlock()

	coords.normalizeValues()
	return coords.x, coords.y, coords.magnitude
}

func (coords *Coords) getMetrics() Metrics {
	coords.mu.Lock()
	defer coords.mu.Unlock()

	coords.normalizeValues()
	return Metrics{
		x:         coords.x,
		y:         coords.y,
		magnitude: coords.magnitude,
		angle:     calcResolvedAngle(coords.x, coords.y),
	}
}

func (coords *Coords) oldGetMetrics() Metrics {
	//x, y := coords.getValues()
	//
	//magnitude := calcMagnitude(x, y)
	//resolvedAngle := calcResolvedAngle(x, y)
	//oneQuarterAngle := calcOneQuarterAngle(resolvedAngle)

	//mappedX, mappedY := mapCircleToSquare(x, y)

	//if magnitude != 0 {
	//fmt.Printf("x: %0.2f, y: %0.2f, mappedX: %0.2f, mappedY: %0.2f, magn: %0.2f\n", coordsMetrics.x, coordsMetrics.y, coordsMetrics.mappedX, coordsMetrics.mappedY, coordsMetrics.magnitude)
	//}

	return Metrics{
		//x:         x,
		//y:         y,
		//magnitude: magnitude,
		//angle:     resolvedAngle,
		//oneQuarterAngle: oneQuarterAngle,
	}
}

type Metrics struct {
	x, y float64
	//mappedX, mappedY float64
	magnitude float64
	angle     int
	//oneQuarterAngle int
}

//func setToRadiusValue(val *float64) {
//	*val = math.Copysign(1.0, *val)
//}

//func (metrics *Metrics) correctValuesNearRadius() {
//	if metrics.magnitude > MaxAccelRadiusThreshold {
//		if MaxAccelMinAngle < metrics.oneQuarterAngle && metrics.oneQuarterAngle < MaxAccelMaxAngle {
//			setToRadiusValue(&metrics.x)
//			setToRadiusValue(&metrics.y)
//		}
//	}
//}

func initMaxAccelValues() {
	if MaxAccelAngleMargin > 45 {
		panicMsg("Incorrect value of \"MaxAccelAngleMargin\": %v\n", MaxAccelAngleMargin)
	}
	MaxAccelMinAngle = 45 - MaxAccelAngleMargin
	MaxAccelMaxAngle = 45 + MaxAccelAngleMargin
}

func applyDeadzone(value *float64) {
	if math.Abs(*value) < Deadzone {
		*value = 0.0
	}
}

const twoSqrt2 float64 = 2.0 * math.Sqrt2

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

func calcOneQuarterAngle(resolvedAngle int) int {
	return floatToInt(math.Mod(float64(resolvedAngle), 90))
}

func resolveAngle(angle float64) int {
	angle = math.Mod(angle+360, 360)
	return floatToInt(angle)
}

const radiansMultiplier float64 = 180 / math.Pi

func calcResolvedAngle(x, y float64) int {
	degrees := math.Atan2(y, x) * radiansMultiplier
	return resolveAngle(degrees)
}

func calcMagnitude(x, y float64) float64 {
	return math.Hypot(x, y)
}

const FloatPrecision int = 8

func normalizeIncorrectEdgeValues(x, y, magnitude *float64) {
	*magnitude = calcMagnitude(*x, *y)
	if *magnitude > 1.0 {
		*x /= *magnitude
		*y /= *magnitude
	}
	trunc(x, FloatPrecision)
	trunc(y, FloatPrecision)
}

const outputRangeMin float64 = 1.0

func convertRange(input, outputMax float64) float64 {
	sign := getSignMakeAbs(&input)

	if input == 0.0 {
		return 0.0
	}

	if input > 1.0 {
		panicMsg("Axis input value is greater than 1.0. Current value: %v\n", input)
	}

	output := outputRangeMin + ((outputMax-outputRangeMin)/inputRange)*(input-Deadzone)
	applySign(sign, &output)
	return output
}

func calcRefreshInterval(input, slowestInterval, fastestInterval float64) time.Duration {
	input = math.Abs(input)
	refreshInterval := convertRange(input, slowestInterval-fastestInterval)
	refreshInterval = slowestInterval - refreshInterval
	return time.Duration(floatToInt64(refreshInterval)) * time.Millisecond
}
