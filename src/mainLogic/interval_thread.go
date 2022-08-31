package mainLogic

import (
	"github.com/positiveway/gofuncs"
	"time"
)

type Interval struct {
	repeatInterval, intervalLeft float64
}

func MakeInterval(repeatInterval float64) *Interval {
	interval := &Interval{}
	interval.SetInterval(repeatInterval)
	return interval
}

func (i *Interval) SetInterval(repeatInterval float64) {
	gofuncs.PanicAnyNotInitOrEmpty(repeatInterval)

	i.repeatInterval = repeatInterval
	i.intervalLeft = repeatInterval
}

func (i *Interval) ResetInterval() bool {
	if i.intervalLeft <= 0 {
		i.intervalLeft = i.repeatInterval
		return true
	}
	return false
}

func (i *Interval) DecreaseInterval() bool {
	i.intervalLeft -= Cfg.tickerInterval
	return i.ResetInterval()
}

type Intervals2 struct {
	X, Y *Interval
}

func MakeIntervals2() *Intervals2 {
	return &Intervals2{
		X: MakeInterval(0),
		Y: MakeInterval(0),
	}
}

func RunGlobalEventsThread() {
	ticker := time.NewTicker(gofuncs.NumberToMillis(Cfg.tickerInterval))

	mouseIntervals := MakeIntervals2()
	scrollIntervals := MakeIntervals2()

	for range ticker.C {
		switch Cfg.ControllerInUse {
		case DualShock:
			MoveMouse(mouseIntervals)
		}
		MoveScroll(scrollIntervals)
		RepeatCommand()
	}
}
