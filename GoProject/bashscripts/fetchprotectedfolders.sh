#!/bin/bash

apt install lynx

lynx -dump -nolist https://refspecs.linuxfoundation.org/FHS_3.0/fhs/index.html | grep -oE '/[^[:space:]<>]+' | sort -u

# This script is unused, but all protected paths in defaultdata/helperFiles/protected_paths come from this command.
