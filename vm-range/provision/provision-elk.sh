#!/bin/bash
apt-get update && apt-get install -y wget apt-transport-https openjdk-11-jdk
apt install -y bpftrace


# Elasticsearch
wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | apt-key add -
echo "deb https://artifacts.elastic.co/packages/7.x/apt stable main" > /etc/apt/sources.list.d/elastic-7.x.list
apt-get update && apt-get install -y elasticsearch kibana
systemctl enable elasticsearch && systemctl start elasticsearch
systemctl enable kibana && systemctl start kibana
