package ipam

import (
	"github.com/rajatchopra/cni-examples/pkg/netutils"
	"net"
)

type Ipam interface {
	GetNextIP() (*net.IP, error)
	ReleaseIP(ip *net.IP) error
	Gateway() (*net.IP, error)
	SetGateway(*net.IP)
	DefaultGateway() (*net.IP, error)
}

type IpamAllocator struct {
	allocator *netutils.SubnetAllocator
	network *net.IPNet
	gw *net.IP
}

func NewIPAM(network string) (Ipam, error) {
	if network == "" {
		network = "10.1.0.0/16"
	}
	_, ipnet, err := net.ParseCIDR(network)
	if err != nil {
		return nil, err
	}
	var inUse []string
	subnetAllocator, err := netutils.NewSubnetAllocator(network, 0, inUse)
	if err != nil {
		return nil, err
	}
	return &IpamAllocator{allocator: subnetAllocator, network: ipnet}, nil
}

func (ipam *IpamAllocator) DefaultGateway() (*net.IP, error) {
	ip := netutils.GenerateDefaultGateway(ipam.network)
	return &ip, nil
}

func (ipam *IpamAllocator) SetGateway(ip *net.IP) {
	ipam.gw = ip
}

func (ipam *IpamAllocator) Gateway() (*net.IP, error) {
	if ipam.gw == nil {
		return ipam.DefaultGateway()
	}
	return ipam.gw, nil
}

func (ipam *IpamAllocator) GetNextIP() (*net.IP, error) {
	ipnet, err := ipam.allocator.GetNetwork()
	if err != nil {
		return nil, err
	}
	if ipnet.IP.String() == ipam.gw.String() {
		ipnet, err = ipam.allocator.GetNetwork()
		if err != nil {
			return nil, err
		}
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
