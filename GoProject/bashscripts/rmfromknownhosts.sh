#!/bin/bash

PATTERN=$1

# Remove lines containing "localhost" from known_hosts
grep -v "$PATTERN" ~/.ssh/known_hosts > /tmp/known_hosts.tmp && mv /tmp/known_hosts.tmp ~/.ssh/known_hosts

echo "Lines containing '$PATTERN' have been removed from ~/.ssh/known_hosts"
