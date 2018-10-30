package driver

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

const defaultFSType = "ext4"

// Mount initiate the host mount
func (d *Driver) Mount(mountDir string) {
	// TODO: check if volume was created
	// TODO: Detach if volume is attached (!! Maybe it's not necessary to detach before attaching?!)
	volume := GetVolume(d.client, d.options.PVOrVolumeName)
	server := GetServer(d.client)
	if !server.Locked {
		Debug("Detach volume for " + mountDir)
		_, _, errDetach := d.client.Volume.Detach(context.Background(), volume)

		if errDetach != nil {
			Debug("Volume was not attached to a server")
		}
	}

	Debug("Attach volume for " + mountDir)
	_, _, errAttach := d.client.Volume.Attach(context.Background(), volume, server)
	if errAttach != nil {
		Failure(errAttach)
	}

	time.Sleep(10)
	err := d.mountAttachedVolume(volume, mountDir)

	if err != nil {
		Debug(err.Error())
		Failure(err)
	}

	Success()
}

func (d Driver) mountAttachedVolume(volume *hcloud.Volume, mountDir string) error {
	blkid, err := RunCommand("blkid", volume.LinuxDevice)
	if err != nil && !strings.Contains(err.Error(), "exit status 2") {
		Failure(err)
	}

	// if device is not formatted, format it
	if blkid == "" {
		if _, err := RunCommand("mkfs", "-t", defaultFSType, volume.LinuxDevice); err != nil {
			Debug(err.Error())
		}
	}

	Debug("ioutil.WriteFile")
	jsonData := []byte(d.rawOptions)
	if err := ioutil.WriteFile(fmt.Sprintf("%s.json", mountDir), jsonData, 0600); err != nil {
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
