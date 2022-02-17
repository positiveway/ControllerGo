package osSpecific

import (
	"fmt"
	"os"
	"path/filepath"
)

const LeftMouse = -3
const RightMouse = -4
const MiddleMouse = -5

const EventServerNotRunningMsg = "Event server is not running"

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetCurFileDir() string {
	ex, err := os.Executable()
	CheckErr(err)
	exPath := filepath.Dir(ex)
	fmt.Printf("Exec path: %s\n", exPath)
	return exPath
}
