#!/bin/bash
apt-get update
apt-get install -y rabbitmq-server
apt install -y bpftrace

systemctl enable rabbitmq-server
systemctl start rabbitmq-server
