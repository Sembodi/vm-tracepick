#!/bin/bash
apt-get update
apt install -y bpftrace

wget https://dl.min.io/server/minio/release/linux-amd64/minio
chmod +x minio
mv minio /usr/local/bin/
useradd -r minio-user -s /sbin/nologin
mkdir -p /usr/local/share/minio /etc/minio
chown minio-user:minio-user /usr/local/share/minio /etc/minio

cat <<EOF >/etc/systemd/system/minio.service
[Unit]
Description=MinIO
After=network.target

[Service]
User=minio-user
ExecStart=/usr/local/bin/minio server /usr/local/share/minio --console-address ":9001"
Environment="MINIO_ROOT_USER=minioadmin"
Environment="MINIO_ROOT_PASSWORD=minioadmin"

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable minio
systemctl start minio
