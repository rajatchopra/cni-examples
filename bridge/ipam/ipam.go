package ipam

import (
	"github.com/openshift/origin/pkg/util/netutils"
	"net"
)

type Ipam interface {
	GetNextIP() (*net.IP, error)
	ReleaseIP(ip *net.IP) error
	Gateway() (*net.IP, error)
}

type IpamAllocator struct {
	allocator *netutils.SubnetAllocator
	network *net.IPNet
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
	_, ipnet, _ := net.ParseCIDR(network)
	return &IpamAllocator{allocator: subnetAllocator, network: ipnet}, nil
}

func (ipam *IpamAllocator) Gateway() (*net.IP, error) {
	ip := netutils.GenerateDefaultGateway(ipam.network)
	return &ip, nil
}

func (ipam *IpamAllocator) GetNextIP() (*net.IP, error) {
	ipnet, err := ipam.allocator.GetNetwork()
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
	return ipam.allocator.ReleaseNetwork(ipnet)
}
