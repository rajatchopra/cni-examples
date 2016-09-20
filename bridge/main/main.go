package main

import (
	"fmt"
	"github.com/rajatchopra/cni-examples/bridge/ipam"
	"net"
	"net/http"
)

var ipamAllocator ipam.Ipam

func main() {
	i, err := ipam.NewIPAM("")
	ipamAllocator = i
	if err != nil {
		fmt.Printf("ERROR[1]: %v", err)
	}

	serve()
	return
}

func serve() {
	http.HandleFunc("/", ipamHandler)       // set router
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		fmt.Printf("ERROR[3]: %v", err)
	}
}

func ipamHandler(w http.ResponseWriter, r *http.Request) {
	switch(r.Method) {
	case "GET":
		r.ParseForm()
		if (r.URL.Path == "/gateway") {
			gw, err := ipamAllocator.Gateway()
			if err != nil {
				http.Error(w, err.Error(), 500)
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
