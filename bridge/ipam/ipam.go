package ipam

import (
	"github.com/openshift/origin/pkg/util/netutils"
	"net"
)

type Ipam interface {
	GetNextIP() (*net.IP, error)
	ReleaseIP(ip *net.IP) error
}

type IpamAllocator struct {
	*netutils.SubnetAllocator
}

func NewIPAM(network string) (Ipam, error) {
	if network == "" {
		network = "10.1.0.0/16"
	}
	var inUse []string
	subnetAllocator, err := netutils.NewSubnetAllocator(network, 0, inUse)
	if err != nil {
		return nil, err
	}
	return &IpamAllocator{subnetAllocator}, nil
}

func (ipam *IpamAllocator) GetNextIP() (*net.IP, error) {
	ipnet, err := ipam.GetNetwork()
	if err != nil {
		return nil, err
	}
	return &ipnet.IP, nil
}

func (ipam *IpamAllocator) ReleaseIP(ip *net.IP) error {
	_, ipnet, err := net.ParseCIDR(ip.String() + "/32")
	if err != nil {
		return err
	}
	return ipam.ReleaseNetwork(ipnet)
}
