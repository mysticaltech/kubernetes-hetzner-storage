package driver

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/hetznercloud/hcloud-go/hcloud"
	h "github.com/stevenklar/kubernetes-hetzner-storage/pkg/hetzner"
)

// JSONParameter contains import kubernetes details about pod and volume
type JSONParameter struct {
	FSGroup        string `json:"kubernetes.io/fsGroup"`
	FSType         string `json:"kubernetes.io/fsType"`
	PVOrVolumeName string `json:"kubernetes.io/pvOrVolumeName"`
	PodName        string `json:"kubernetes.io/pod.name"`
	PodNamespace   string `json:"kubernetes.io/pod.namespace"`
	PodUID         string `json:"kubernetes.io/pod.uid"`
	ReadWrite      string `json:"kubernetes.io/readwrite"`
	ServiceAccount string `json:"kubernetes.io/serviceAccount.name"`
	Token          string `json:"stevenklar/hetzner-cloud-driver/token"`
}

// Driver contains options and client information
type Driver struct {
	options    JSONParameter
	rawOptions string
	client     *hcloud.Client
}

// Run executes the driver routine
func Run() {
	var command, mountDir, jsonOptions string

	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	if len(os.Args) > 2 {
		mountDir = os.Args[2]
	}
	if len(os.Args) > 3 {
		jsonOptions = os.Args[3]
	} else {
		jsonOptions = getJsonOptionsByFile(mountDir, jsonOptions)
	}

	Debug(fmt.Sprintf("%s %s %s", command, mountDir, jsonOptions))

	switch command {
	case "init":
		Initialize()
	case "mount":
		driver := newDriver(jsonOptions)
		driver.Mount(mountDir)
	case "unmount":
		driver := newDriver(jsonOptions)
		driver.Unmount(mountDir)
	default:
		NotSupported()
	}
}

func getJsonOptionsByFile(mountDir string, jsonOptions string) string {
	jsonOptionsFile := fmt.Sprintf("%s.json", mountDir)
	Debug("Read json options file: " + jsonOptionsFile)
	byt, err := ioutil.ReadFile(jsonOptionsFile)
	if err != nil {
		Debug("Error reading options json file")
		Failure(err)
	}
	jsonOptions = string(byt)
	return jsonOptions
}

func newDriver(jsonOptions string) *Driver {
	byt := []byte(jsonOptions)
	options := JSONParameter{}
	if err := json.Unmarshal(byt, &options); err != nil {
		Failure(err)
	}

	client := h.GetClient(options.Token)
	driver := Driver{options: options, rawOptions: jsonOptions, client: client}

	return &driver
}

// GetVolume wraps the hetzner volume finder by name
func GetVolume(client *hcloud.Client, name string) *hcloud.Volume {
	volume, _, _ := client.Volume.GetByName(context.Background(), name)

	return volume
}

// GetServer retreive the server information equals the host machine with application run on
func GetServer(client *hcloud.Client) *hcloud.Server {
	// get all hetzner nodes
	servers, err := client.Server.All(context.Background())

	if err != nil {
		Failure(err)
	}

	// get all interface ips
	ifaces, err := net.Interfaces()

	if err != nil {
		Debug("No interfaces found")
		return nil
	}

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
