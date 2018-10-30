package driver

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

const (
	debugLogFile = "/tmp/hetzner-storage.log"
)

// Debug writes log to file
func Debug(message string) {
	if _, err := os.Stat(debugLogFile); err == nil {
		f, _ := os.OpenFile(debugLogFile, os.O_APPEND|os.O_WRONLY, 0600)
		defer f.Close()
		f.WriteString(fmt.Sprintln(message))
	}
}

// Success write debug log and exits the program with signal 0
func Success() {
	Debug("SUCCESS")
	fmt.Print("{\"status\": \"Success\"}")
	os.Exit(0)
}

// Failure writes debug log and exits the program with signal 1
func Failure(err error) {
	Debug(fmt.Sprintf("FAILURE - %s", err.Error()))

	failureMap := map[string]string{"status": "Failure", "message": err.Error()}
	jsonMessage, _ := json.Marshal(failureMap)
	fmt.Print(string(jsonMessage))

	os.Exit(1)
}

// RunCommand executes command(s) on host machine
func RunCommand(cmd string, args ...string) (string, error) {
	Debug(fmt.Sprintf("Running %s %s", cmd, args))
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Error running %s %v: %v, %s", cmd, args, err, out)
	}
	return string(out), nil
}

func NotSupported() {
	fmt.Print("{\"status\": \"Not supported\"}")
	os.Exit(1)
}

func Initialize() {
	fmt.Print("{\"status\": \"Success\", \"capabilities\": {\"attach\": false}}")
	os.Exit(0)
}

