package main

import (
	"os"

	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"k8s.io/utils/exec"
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

