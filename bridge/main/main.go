package main

import (
	"fmt"
	"github.com/rajatchopra/cni-examples/bridge/ipam"
	"flag"
	"net"
	"net/http"
)

var ipamAllocator ipam.Ipam

func main() {
	var subnet string
	flag.StringVar(&subnet, "subnet", "10.1.0.0/16", "Subnet for the IPAM. Defaults to 10.1.0.0/24")
	flag.Parse()
	_, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		fmt.Printf("ERROR[1]: Not a valid CIDR: %v", subnet)
		return
	}
	i, err := ipam.NewIPAM(ipnet.String())
	if err != nil {
		fmt.Printf("ERROR[2]: %v", err)
		return
	}
	ipamAllocator = i
	gw, err := ipamAllocator.DefaultGateway()
	if err != nil {
		fmt.Printf("ERROR[3]: %v", err)
		return
	}
	ipamAllocator.SetGateway(gw)

	serve()
	return
}

func serve() {
	http.HandleFunc("/", ipamHandler)       // set router
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		fmt.Printf("ERROR[4]: %v", err)
	}
}

func ipamHandler(w http.ResponseWriter, r *http.Request) {
	switch(r.Method) {
	case "GET":
		r.ParseForm()
		if (r.URL.Path == "/gateway") {
			gw, err := ipamAllocator.Gateway()
			if err != nil {
				fmt.Printf("ERROR[3]: %v", err)
				return
			}
			fmt.Fprintf(w,"%s", gw.String())
			return
		}
		ip, err := ipamAllocator.GetNextIP()
		if err != nil {
			fmt.Printf("ERROR[2]: %v", err)
			return
		}
		fmt.Fprintf(w, "%s", ip.String())
	case "DELETE":
		//ipamAllocator.ReleaseIP(r.String())
		r.ParseForm()  
	        ipStr := r.URL.Path[1:len(r.URL.Path)] 
		ip := net.ParseIP(ipStr)
		if ip == nil {
			http.Error(w, fmt.Errorf("Invalid IP: %s", ipStr).Error(), 500)
			return
		}
		err := ipamAllocator.ReleaseIP(&ip)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
        default:
		fmt.Printf("Unknown type of request %v\n", r)
	}
	return
}
