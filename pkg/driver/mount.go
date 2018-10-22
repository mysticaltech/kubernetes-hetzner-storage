package driver

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"github.com/hetznercloud/hcloud-go/hcloud"
	c "github.com/stevenklar/kubernetes-hetzner-storage/pkg/common"
)

// Mount initiate the host mount
func (d *Driver) Mount(mountDir string) {
	// TODO: check if volume was created
	// TODO: Detach if volume is attached (!! Maybe it's not necessary to detach before attaching?!)
	volume := c.GetVolume(d.client, d.options.PVOrVolumeName)
	server := c.GetServer(d.client)
	_, _, err := d.client.Volume.Attach(context.Background(), volume, server)

	if err != nil {
		c.Failure(err)
	}

	// TODO: Retrieve attached volume information
	mountAttachedVolume(volume, mountDir)

	if err != nil {
		c.Failure(err)
	}

	c.Success()
}

func mountAttachedVolume(volume *hcloud.Volume, mountDir string) error {
	blkid, err := c.RunCommand("blkid", volume.LinuxDevice)
	if err != nil && !strings.Contains(err.Error(), "exit status 2") {
		c.Failure(err)
	}

	// if device is not formatted, format it
	if blkid == "" {
		if _, err := c.RunCommand("mkfs", "-t", "ext4", volume.LinuxDevice); err != nil {
			c.Failure(err)
		}
	}

	c.Debug("ioutil.WriteFile")
	if err := ioutil.WriteFile(fmt.Sprintf("%s.json", mountDir), nil, 0600); err != nil {
		c.Failure(err)
	}

	c.Debug("os.MkdirAll")
	if err := os.MkdirAll(mountDir, 0755); err != nil {
		c.Failure(err)
	}

	c.Debug("syscall.Mount")
	if err := syscall.Mount(volume.LinuxDevice, mountDir, "ext4", 0, ""); err != nil {
		c.Failure(err)
	}

	return nil
}
