package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const (
	debugLogFile = "/tmp/hetzner-driver.log"
)

type jsonParameter struct {
	FSGroup			string `json:"kubernetes.io/fsGroup"`
	FSType			string `json:"kubernetes.io/fsType"`
	PVOrVolumeName	string `json:"kubernetes.io/pvOrVolumeName"`
	PodName		    string `json:"kubernetes.io/pod.name"`
	PodNamespace	string `json:"kubernetes.io/pod.namespace"`
	PodUID			string `json:"kubernetes.io/pod.uid"`
	ReadWrite		string `json:"kubernetes.io/readwrite"`
	ServiceAccount	string `json:"kubernetes.io/serviceAccount.name"`
	Token			string ``
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

func getClient(token string) *hcloud.Client {
	return hcloud.NewClient(hcloud.WithToken(token))
}

func getServer(client *hcloud.Client) *hcloud.Server {
	// get all hetzner nodes
	servers, err := client.Server.All(context.Background())

	if err != nil {
		failure(err)
	}

	// get all interface ips
	ifaces, err := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil {
				continue
			}

			for _, server := range servers {
				if server.PublicNet.IPv4.IP.String() == ip.String() {
					return server
				}
			}
		}
	}

	return nil
}

func mountAttachedVolume(volume *hcloud.Volume, mountDir string) error {
	// mkdir /mnt/%volume.name%
	// mount -o discord,defaults %volume.linux_device% /mnt/%volume.name%

	blkid, err := run("blkid", volume.LinuxDevice)
	if err != nil && !strings.Contains(err.Error(), "exit status 2") {
		failure(err)
	}

	if blkid == "" {
		if _, err := run("mkfs", "-t", "ext4", volume.LinuxDevice); err != nil {
			failure(err)
		}
	}

	debug("ioutil.WriteFile")
	if err := ioutil.WriteFile(fmt.Sprintf("%s.json", mountDir), nil, 0600); err != nil {
		failure(err)
	}

	debug("os.MkdirAll")
	if err := os.MkdirAll(mountDir, 0755); err != nil {
		failure(err)
	}

	debug("syscall.Mount")
	if err := syscall.Mount(volume.LinuxDevice, mountDir, "ext4", 0, ""); err != nil {
		failure(err)
	}

	return nil
}

func mount(mountDir, jsonOptions string) {
	byt := []byte(jsonOptions)
	options := jsonParameter{}
	if err := json.Unmarshal(byt, &options); err != nil {
		failure(err)
	}

	client := getClient(options.Token)

	// TODO: check if volume was created
	// TODO: Detach if volume is attached (!! Maybe it's not necessary to detach before attaching?!)

	volumeName := options.PVOrVolumeName
	volume, _, _:= client.Volume.GetByName(context.Background(), volumeName)
	server := getServer(client)
	_, _, err := client.Volume.Attach(context.Background(), volume, server)

	if err != nil {
		failure(err)
	}

	// TODO: Retrieve attached volume information
	mountAttachedVolume(volume, mountDir)

	if err != nil {
		failure(err)
	}

	success()
}

func unmount(mountDir string) {
	client := getClient("")

	volumeName := "" // TODO: get volume name from options?
	volume, _, _:= client.Volume.GetByName(context.Background(), volumeName)
	_, _, err := client.Volume.Detach(context.Background(), volume)

	if err != nil {
		failure(err)
	}

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
