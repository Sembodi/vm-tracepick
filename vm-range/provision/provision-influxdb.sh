#!/bin/bash
apt-get update
wget -qO- https://repos.influxdata.com/influxdb.key | apt-key add -
echo "deb https://repos.influxdata.com/ubuntu jammy stable" | tee /etc/apt/sources.list.d/influxdb.list
apt-get update
apt-get install -y influxdb
apt install -y bpftrace

systemctl enable influxdb
systemctl start influxdb
