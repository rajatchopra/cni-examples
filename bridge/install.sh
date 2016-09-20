#!/bin/bash

mkdir -p /etc/cni/net.d
cp scripts/10-mynet.conf /etc/cni/net.d/

mkdir -p /opt/cni/bin
cp scripts/localbridge /opt/cni/bin/
cp scripts/loopback /opt/cni/bin/

cd main/
go build -o /usr/bin/ocicni-localbridge-ipam main.go


function destroy() {
	ip link set dev ocibr0 down || echo
	brctl delbr ocibr0 || echo
	killall ocicni-localbridge-ipam || echo
}

function setup() {
	/usr/bin/ocicni-localbridge-ipam &

	# get couple of IPs out of the way, the current logic will distribute the gateway as one of the valid IPs too
	curl "http://localhost:9090/" &> /dev/null
	curl "http://localhost:9090/" &> /dev/null
	curl "http://localhost:9090/" &> /dev/null
	curl "http://localhost:9090/" &> /dev/null

	brctl addbr ocibr0
	ip link set dev ocibr0 up

	gateway_ip=`curl -s http://localhost:9090/gateway`
	ip addr add ${gateway_ip}/24 dev ocibr0
}

destroy
setup
