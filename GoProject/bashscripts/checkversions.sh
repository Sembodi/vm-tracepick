#!/bin/bash

REPO=$1
TARGET_VERSION=$2

if [ -z "$REPO" ] || [ -z "$TARGET_VERSION" ]; then
  echo "Usage: $0 <repo> <target_version>"
  exit 1
fi

TAGS=$(curl -s "https://registry.hub.docker.com/v2/repositories/library/$REPO/tags?page_size=100" | jq -r '.results[].name')

# 1. Check for exact match
if echo "$TAGS" | grep -qx "$TARGET_VERSION"; then
    echo "$TARGET_VERSION"
    exit 0
fi

# 2. Try closest match (same major.minor, e.g. 20.04)
BASE_VERSION="${TARGET_VERSION%.*}"

CLOSEST=$(echo "$TAGS" | grep "^$BASE_VERSION" | sort -Vr | head -n1)

if [ -n "$CLOSEST" ]; then
    echo "$CLOSEST"
    exit 0
fi

# 3. Not found
echo "tagnotfound"
exit 0
