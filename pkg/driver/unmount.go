package driver

import (
	"context"
)

// Unmount removes mounted volume
func (d *Driver) Unmount(mountDir string) {
	volume := GetVolume(d.client, d.options.PVOrVolumeName)
	_, _, err := d.client.Volume.Detach(context.Background(), volume)

	if err != nil {
		Failure(err)
	}

	Success()
}
