//go:build !windows

package osSpecific

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

const DefaultProjectDir string = "/home/user/GolandProjects/ControllerGo"

var GetLocaleExecPath string

func InitResources() {
	GetLocaleExecPath = GetCurFileDir() + "/getLocale.sh"
}

func SetHighPriority(pid int) {
	err := syscall.Setpriority(syscall.PRIO_PROCESS, pid, -20)
	CheckErr(err)
}

func CheckProcess(pid int) {
	process, err := os.FindProcess(pid)
	if err != nil {
		panic(fmt.Sprintf("Failed to find process: %s\n", err))
	} else {
		err := process.Signal(syscall.Signal(0))
		if err != nil {
			panic(fmt.Sprintf("%s\nprocess.Signal on pid %d returned: %v\n",
				EventServerNotRunningMsg, pid, err))
		}
	}
}

func GetLocale() string {
	cmd := exec.Command("sh", GetLocaleExecPath)
	stdout, err := cmd.Output()

	if err != nil {
		panic(err.Error())
	}

	return string(stdout)
}
