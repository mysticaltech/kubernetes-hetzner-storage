package driver

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hetznercloud/hcloud-go/hcloud"
	c "github.com/stevenklar/kubernetes-hetzner-storage/pkg/common"
)

// JsonParameter contains import kubernetes details about pod and volume
type JsonParameter struct {
	FSGroup        string `json:"kubernetes.io/fsGroup"`
	FSType         string `json:"kubernetes.io/fsType"`
	PVOrVolumeName string `json:"kubernetes.io/pvOrVolumeName"`
	PodName        string `json:"kubernetes.io/pod.name"`
	PodNamespace   string `json:"kubernetes.io/pod.namespace"`
	PodUID         string `json:"kubernetes.io/pod.uid"`
	ReadWrite      string `json:"kubernetes.io/readwrite"`
	ServiceAccount string `json:"kubernetes.io/serviceAccount.name"`
	Token          string ``
}

// Driver contains options and client information
type Driver struct {
	options JsonParameter
	client  *hcloud.Client
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
	}

	c.Debug(fmt.Sprintf("%s %s %s", command, mountDir, jsonOptions))

	switch command {
	case "init":
		initialize()
	case "mount":
		driver := newDriver(jsonOptions)
		driver.Mount(mountDir)
	case "unmount":
		driver := newDriver(jsonOptions)
		driver.Unmount(mountDir)
	default:
		fmt.Print("{\"status\": \"Not supported\"}")
		os.Exit(1)
	}
}

func initialize() {
	fmt.Print("{\"status\": \"Success\", \"capabilities\": {\"attach\": false}}")
	os.Exit(0)
}

func newDriver(jsonOptions string) *Driver {
	byt := []byte(jsonOptions)
	options := JsonParameter{}
	if err := json.Unmarshal(byt, &options); err != nil {
		c.Failure(err)
	}

	client := c.GetClient(options.Token)
	driver := Driver{options: options, client: client}

	return &driver
}
