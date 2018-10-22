package provisioner

import (
	"os"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"k8s.io/utils/exec"
)

type hetznerProvisioner struct {
	runner exec.Interface
	token  string
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
