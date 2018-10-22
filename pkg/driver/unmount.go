package driver

import (
	"context"

	c "github.com/stevenklar/kubernetes-hetzner-storage/pkg/common"
)

// Unmount removes mounted volume
func (d *Driver) Unmount(mountDir string) {
	volume := c.GetVolume(d.client, d.options.PVOrVolumeName)
	_, _, err := d.client.Volume.Detach(context.Background(), volume)

	if err != nil {
		c.Failure(err)
	}

	c.Success()
}
