package mainLogic

import "time"

//layout
const LayoutInUse string = "Linux"

//path
const DefaultProjectDir string = "/home/user/GolandProjects/ControllerGo"

//commands
const TriggerThreshold float64 = 0.15
const holdThreshold = 150 * time.Millisecond

//mouse
const mouseMaxMove float64 = 4
const forcePower float64 = 1.5
const deadzone float64 = 0.06

//const mouseScaleFactor float64 = 3
//var mouseIntervalInt int = int(math.Round(mouseMaxMove*mouseScaleFactor))
const mouseIntervalInt int = 12
const mouseInterval = time.Duration(mouseIntervalInt) * time.Millisecond

const scrollFastestInterval float64 = 20
const scrollSlowestInterval float64 = 250

const scrollIntervalRange = scrollSlowestInterval - scrollFastestInterval
const horizontalScrollThreshold float64 = 0.45

//typing
const angleMargin int = 16
const magnitudeThresholdPct float64 = 35
const MagnitudeThreshold = magnitudeThresholdPct / 100

//common
const DefaultWaitInterval time.Duration = 25 * time.Millisecond

//web socket
const SocketPort int = 1234
const SocketIP string = "0.0.0.0"
