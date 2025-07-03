#!/bin/bash
apt-get update
apt-get install -y unzip curl
apt install -y bpftrace

CONSUL_VERSION="1.10.0"
cd /tmp
curl -O https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_linux_amd64.zip
unzip consul_${CONSUL_VERSION}_linux_amd64.zip
mv consul /usr/local/bin/
useradd --system --home /etc/consul.d --shell /bin/false consul
mkdir -p /etc/consul.d
consul agent -dev -client=0.0.0.0 &
