package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

const (
	debugLogFile = "/tmp/hetzner-driver.log"
)

type jsonParameter struct {
	FSGroup		   string `json:"kubernetes.io/fsGroup"`
	FSType			string `json:"kubernetes.io/fsType"`
	PVOrVolumeName	string `json:"kubernetes.io/pvOrVolumeName"`
	PodName		   string `json:"kubernetes.io/pod.name"`
	PodNamespace	  string `json:"kubernetes.io/pod.namespace"`
	PodUID			string `json:"kubernetes.io/pod.uid"`
	ReadWrite		 string `json:"kubernetes.io/readwrite"`
	ServiceAccount	string `json:"kubernetes.io/serviceAccount.name"`
}

func main() {
	var command string
	var mountDir string
	var jsonOptions string

	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	if len(os.Args) > 2 {
		mountDir = os.Args[2]
	}
	if len(os.Args) > 3 {
		jsonOptions = os.Args[3]
	}

	debug(fmt.Sprintf("%s %s %s", command, mountDir, jsonOptions))

	switch command {
	case "init":
		fmt.Print("{\"status\": \"Success\", \"capabilities\": {\"attach\": false}}")
		os.Exit(0)
	case "mount":
		mount(mountDir, jsonOptions)
	case "unmount":
		unmount(mountDir)
	default:
		fmt.Print("{\"status\": \"Not supported\"}")
		os.Exit(1)
	}
}

func debug(message string) {
	if _, err := os.Stat(debugLogFile); err == nil {
		f, _ := os.OpenFile(debugLogFile, os.O_APPEND|os.O_WRONLY, 0600)
		defer f.Close()
		f.WriteString(fmt.Sprintln(message))
	}
}

func success() {
	debug("SUCCESS")

	fmt.Print("{\"status\": \"Success\"}")

	os.Exit(0)
}

func failure(err error) {
	debug(fmt.Sprintf("FAILURE - %s", err.Error()))

	failureMap := map[string]string{"status": "Failure", "message": err.Error()}
	jsonMessage, _ := json.Marshal(failureMap)
	fmt.Print(string(jsonMessage))

	os.Exit(1)
}

func mount(mountDir, jsonOptions string) {
	// TODO: To be implemented
	success()
}

func unmount(mountDir string) {
	// TODO: To be implemented
	success()
}

func run(cmd string, args ...string) (string, error) {
	debug(fmt.Sprintf("Running %s %s", cmd, args))
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Error running %s %v: %v, %s", cmd, args, err, out)
	}
	return string(out), nil
}
