#!/bin/bash
for vm in $(vagrant status | grep running | awk '{print $1}'); do
  echo "[$vm]"
  vagrant port $vm
  echo
done
