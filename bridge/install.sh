#!/bin/bash

mkdir -p /etc/cni/net.d
cp scripts/10-mynet.conf /etc/cni/net.d/

mkdir -p /opt/cni/bin
cp scripts/localbridge /opt/cni/bin/
cp scripts/loopback /opt/cni/bin/
