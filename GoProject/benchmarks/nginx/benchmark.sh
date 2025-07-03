#!/bin/bash

# Define parameters
URL="http://localhost:80"
THREADS=(1)
CONNECTIONS=(10 50 100 200 400 800)
DURATION="10s"

# Loop over combinations
for t in "${THREADS[@]}"; do
  for c in "${CONNECTIONS[@]}"; do
    echo "Running wrk with $t threads and $c connections for $DURATION..."
    wrk -t$t -c$c -d$DURATION "$URL"
    echo ""
  done
done
