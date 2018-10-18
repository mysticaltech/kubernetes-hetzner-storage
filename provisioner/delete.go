package main

import (
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
)

func (p *hetznerProvisioner) Delete(volume *v1.PersistentVolume) error {
	glog.Infof("Delete called for volume: %s", volume.Name)

    // TODO: Delete from hetzner

	return nil
}

