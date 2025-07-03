#!/bin/bash
apt-get update
apt-get install -y redis-server
apt install -y bpftrace


# Modify Redis config
REDIS_CONF="/etc/redis/redis.conf"

# Set 'protected-mode no'
sed -i 's/^protected-mode .*/protected-mode no/' $REDIS_CONF

# Set 'bind 0.0.0.0'
sed -i 's/^bind .*/bind 0.0.0.0/' $REDIS_CONF

# Optional: make sure Redis listens on the right port (6379)
# sed -i 's/^port .*/port 6379/' $REDIS_CONF

# Restart and enable service
systemctl enable redis-server
systemctl restart redis-server
