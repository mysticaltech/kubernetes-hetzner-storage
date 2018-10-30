package hetzner

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
)

// GetClient wraps the hetzner client with just the token
func GetClient(token string) *hcloud.Client {
	return hcloud.NewClient(hcloud.WithToken(token))
}
