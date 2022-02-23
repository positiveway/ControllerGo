//go:build linux

package osSpecific

import (
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

func GetLocale() string {
	cmd := exec.Command("sh", GetLocaleExecPath)
	stdout, err := cmd.Output()

	if err != nil {
		panic(err.Error())
	}

	return string(stdout)
}
