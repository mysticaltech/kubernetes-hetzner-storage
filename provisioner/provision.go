package main

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const defaultFSType = "ext4"

func (p *hetznerProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	glog.Infof("Provision called for volume: %s", options.PVName)

	if err := p.provisionOnHetznerCloud(options); err != nil {
		glog.Errorf("Failed to provision volume %s, error: %s", options, err.Error())
		return nil, err
	}

	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: options.PVName,
		},
		Spec: v1.PersistentVolumeSpec{
			AccessModes: options.PVC.Spec.AccessModes,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
			},
			PersistentVolumeReclaimPolicy: options.PersistentVolumeReclaimPolicy,
			PersistentVolumeSource: v1.PersistentVolumeSource{
				FlexVolume: &v1.FlexPersistentVolumeSource{
					Driver: driver,
					FSType: defaultFSType,
					Options: map[string]string{
						driverOptionToken: p.token,
					},
				},
			},
		},
	}

	return pv, nil
}

func (p *hetznerProvisioner) provisionOnHetznerCloud(options controller.VolumeOptions) error {
    client := p.getClient(p.token)

    capacity, exists := options.PVC.Spec.Resources.Requests[v1.ResourceStorage]
	if !exists {
		return fmt.Errorf("Capacity was not specified for name label %s", options.PVName)
	}

    hetznerCapacity := (((int(capacity.Value()) / 1024) / 1024) / 1024) // kuberntes uses bytes, hetzner uses gbytes

	glog.Infof("Would create volume with capacity: %s", hetznerCapacity)

    // create volume with given volume options
	opts := hcloud.VolumeCreateOpts{
		Name: options.PVName,
        //Size: hetznerCapacity,
        Size: 10,
		Location: &hcloud.Location{Name: "nbg1"}, // not available in falkenstein yet? - add default option and random picker?
	}

	_, _, err := client.Volume.Create(context.Background(), opts)

    return err
}
