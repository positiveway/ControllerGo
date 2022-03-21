//go:build linux

package osSpec

import (
	"os/exec"
)

const DefaultProjectDir string = "/home/user/GolandProjects/ControllerGo"

var GetLocaleExecPath string

func InitResources() {
	GetLocaleExecPath = GetCurFileDir() + "/getLocale.sh"
}

func GetLocale() string {
	cmd := exec.Command("sh", GetLocaleExecPath)
	stdout, err := cmd.Output()

	if err != nil {
		panic(err.Error())
	}

	return string(stdout)
}
