#!/bin/bash
apt-get update
apt-get install -y apache2
apt install -y bpftrace
systemctl enable apache2
systemctl start apache2
