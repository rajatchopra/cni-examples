#!/bin/bash



function destroy() {
	ip link set dev ocibr0 down || echo
	brctl delbr ocibr0 || echo
	killall ocicni-localbridge-ipam || echo
}

function install() {
	mkdir -p /etc/cni/net.d
	cp scripts/10-mynet.conf /etc/cni/net.d/

	mkdir -p /opt/cni/bin
	cp scripts/localbridge /opt/cni/bin/
	cp scripts/loopback /opt/cni/bin/

	make clean; make
	cp _output/bin/ipam /usr/bin/ocicni-localbridge-ipam
}

function setup() {

	# build ipam and install scripts
	install

	# run ipam
	/usr/bin/ocicni-localbridge-ipam &

	# test it works.. just get an ip (discard it)
	curl "http://localhost:9090/" &> /dev/null
	if [ $? -ne 0 ]; then
		echo "IPAM failed"
		exit 1
	fi

	brctl addbr ocibr0
	ip link set dev ocibr0 up

	gateway_ip=`curl -s http://localhost:9090/gateway`
	ip addr add ${gateway_ip}/24 dev ocibr0
}

destroy
setup
