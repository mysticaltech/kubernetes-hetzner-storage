package driver

import (
	"context"
	"syscall"
)

// Unmount removes mounted volume
func (d *Driver) Unmount(mountDir string) {
	Debug("findmnt: " + mountDir)
	_, err := RunCommand("findmnt", "-n", "-o", "SOURCE", "--target", mountDir)
	if err != nil {
		Debug(err.Error())
	}

	Debug("syscall.Unmount: " + mountDir)
	if err := syscall.Unmount(mountDir, 0); err != nil {
		Failure(err)
	}

	Debug("Detach hetzner volume from server")
	volume := GetVolume(d.client, d.options.PVOrVolumeName)
	_, _, errDetach := d.client.Volume.Detach(context.Background(), volume)

	if errDetach != nil {
		Failure(errDetach)
	}

	// Delete json file with token in it
	//Debug("os.Remove")
	//if err := os.Remove(jsonOptionsFile); err != nil {
	//	failure(err)
	//}

	Success()
}
