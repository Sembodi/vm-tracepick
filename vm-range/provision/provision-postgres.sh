#!/bin/bash
apt-get update
apt-get install -y postgresql postgresql-contrib
apt install -y bpftrace

systemctl enable postgresql
systemctl start postgresql
