package main

import (
	"os"

	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"k8s.io/utils/exec"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

type hetznerProvisioner struct {
	runner            exec.Interface
	token             string
}

func NewProvisioner() controller.Provisioner {
	return &hetznerProvisioner{
		runner: exec.New(),
		token:  os.Getenv("HETZNER_TOKEN"),
	}
}

func (p *hetznerProvisioner) getClient(token string) *hcloud.Client {
	client := hcloud.NewClient(hcloud.WithToken(token))
    return client
}

