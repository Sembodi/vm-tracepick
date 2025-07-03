#!/bin/bash
apt-get update
apt-get install -y nginx
apt install -y bpftrace

systemctl enable nginx
systemctl start nginx
