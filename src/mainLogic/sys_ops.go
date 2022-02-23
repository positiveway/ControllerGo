package mainLogic

import (
	"ControllerGo/src/osSpecific"
	"os"
)

func SetSelfPriority() {
	pid := os.Getppid()
	osSpecific.SetHighPriority(pid)
}
