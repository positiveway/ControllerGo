package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const EventServerNotRunningMsg = "Event server is not running"

var proc *os.Process
var EventServerExecPath string = getCurFileDir() + "/ControllerRust"

func getCurFileDir() string {
	ex, err := os.Executable()
	check_err(err)
	exPath := filepath.Dir(ex)
	fmt.Printf("Exec path: %s\n", exPath)
	return exPath
}

func checkProcess(pid int) {
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

func checkEventServer() {
	if proc == nil {
		panic(EventServerNotRunningMsg)
	}
	checkProcess(proc.Pid)
}

func startProc(execPath string) {
	cmd := exec.Command(execPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	check_err(err)
	proc = cmd.Process
	pid := proc.Pid
	checkEventServer()
	setHighPriority(pid)
	fmt.Printf("Event service started. PID: %d\n", pid)
}

func startEventServer() {
	startProc(EventServerExecPath)
}

func killEventServer() {
	if proc != nil {
		proc.Kill()
	}
}

func setHighPriority(pid int) {
	err := syscall.Setpriority(syscall.PRIO_PROCESS, pid, -20)
	check_err(err)
}

func setSelfPriority() {
	pid := os.Getppid()
	setHighPriority(pid)
}
