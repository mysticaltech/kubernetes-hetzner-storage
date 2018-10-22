package common

import (
	"context"
	"net"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// GetClient wraps the hetzner client with just the token
func GetClient(token string) *hcloud.Client {
	return hcloud.NewClient(hcloud.WithToken(token))
}

// GetVolume wraps the hetzner volume finder by name
func GetVolume(client *hcloud.Client, name string) *hcloud.Volume {
	volume, _, _ := client.Volume.GetByName(context.Background(), name)

	return volume
}

// GetServer retreive the server information equals the host machine with application run on
func GetServer(client *hcloud.Client) *hcloud.Server {
	// get all hetzner nodes
	servers, err := client.Server.All(context.Background())

	if err != nil {
		Failure(err)
	}

	// get all interface ips
	ifaces, err := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil {
				continue
			}

			for _, server := range servers {
				if server.PublicNet.IPv4.IP.String() == ip.String() {
					return server
				}
			}
		}
	}

	return nil
}
