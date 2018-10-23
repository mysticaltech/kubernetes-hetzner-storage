package driver

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

const defaultFSType = "ext4"

// Mount initiate the host mount
func (d *Driver) Mount(mountDir string) {
	// TODO: check if volume was created
	// TODO: Detach if volume is attached (!! Maybe it's not necessary to detach before attaching?!)
	volume := GetVolume(d.client, d.options.PVOrVolumeName)
	server := GetServer(d.client)
	_, _, err := d.client.Volume.Attach(context.Background(), volume, server)

	if err != nil {
		Failure(err)
	}

	// TODO: Retrieve attached volume information
	mountAttachedVolume(volume, mountDir)

	if err != nil {
		Failure(err)
	}

	Success()
}

func mountAttachedVolume(volume *hcloud.Volume, mountDir string) error {
	blkid, err := RunCommand("blkid", volume.LinuxDevice)
	if err != nil && !strings.Contains(err.Error(), "exit status 2") {
		Failure(err)
	}

	// if device is not formatted, format it
	if blkid == "" {
		if _, err := RunCommand("mkfs", "-t", defaultFSType, volume.LinuxDevice); err != nil {
			Failure(err)
		}
	}

	Debug("ioutil.WriteFile")
	if err := ioutil.WriteFile(fmt.Sprintf("%s.json", mountDir), nil, 0600); err != nil {
		Failure(err)
	}

	Debug("os.MkdirAll")
	if err := os.MkdirAll(mountDir, 0755); err != nil {
		Failure(err)
	}

	Debug("syscall.Mount")
	if err := syscall.Mount(volume.LinuxDevice, mountDir, defaultFSType, 0, ""); err != nil {
		Failure(err)
	}

	return nil
}
