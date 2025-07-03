#!/bin/bash
apt-get update
apt-get install -y software-properties-common
apt install -y bpftrace

add-apt-repository "deb https://packages.grafana.com/oss/deb stable main"
wget -q -O - https://packages.grafana.com/gpg.key | apt-key add -
apt-get update
apt-get install -y grafana
systemctl enable grafana-server
systemctl start grafana-server
