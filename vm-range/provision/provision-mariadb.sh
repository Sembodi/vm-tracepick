#!/bin/bash
apt-get update
debconf-set-selections <<< 'mariadb-server mysql-server/root_password password root'
debconf-set-selections <<< 'mariadb-server mysql-server/root_password_again password root'
apt-get install -y mariadb-server
apt install -y bpftrace

systemctl enable mariadb
systemctl start mariadb
