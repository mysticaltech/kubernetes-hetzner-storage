package main

import (
	"context"
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
)

func (p *hetznerProvisioner) Delete(volume *v1.PersistentVolume) error {
	glog.Infof("Delete called for volume: %s", volume.Name)

	client := p.getClient(p.token)
	hetznerVolume, _, err := client.Volume.GetByName(context.Background(), volume.Name)

	if err != nil {
		glog.Infof("Delete failed for volume: %s", volume.Name)
		return err
	}

	response, err := client.Volume.Delete(context.Background(), hetznerVolume)
	glog.Infoln(response)

	return err
}

