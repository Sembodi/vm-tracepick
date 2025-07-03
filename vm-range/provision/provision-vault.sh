#!/bin/bash
apt-get update
apt-get install -y unzip curl
apt install -y bpftrace

VAULT_VERSION="1.8.0"
cd /tmp
curl -O https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_linux_amd64.zip
unzip vault_${VAULT_VERSION}_linux_amd64.zip
mv vault /usr/local/bin/
useradd --system --home /etc/vault.d --shell /bin/false vault
mkdir -p /etc/vault.d
vault server -dev -dev-root-token-id=root &
