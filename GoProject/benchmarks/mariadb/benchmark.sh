#!/bin/bash

MYSQL_USER="root"
MYSQL_PASS="root"
MYSQL_DB="sys"
THREADS=(1 2 4 8)
DURATION=5
REPS=3
OUTPUT_FILE="output.txt"

# Cleanup after test
sysbench \
  --db-driver=mysql \
  --mysql-user=$MYSQL_USER \
  --mysql-password=$MYSQL_PASS \
  --mysql-db=$MYSQL_DB \
  oltp_read_write \
  cleanup

# Prepare test data
sysbench \
  --db-driver=mysql \
  --mysql-user=$MYSQL_USER \
  --mysql-password=$MYSQL_PASS \
  --mysql-db=$MYSQL_DB \
  oltp_read_write \
  prepare

for t in "${THREADS[@]}"; do

  total_tps=0
  total_latency=0
  total_errors=0

  # echo "Running test with $t threads, $REPS times for $DURATION seconds each..."
  for ((i=1; i<=$REPS; i++)); do
    # echo "Run #$i"
    output=$(sysbench \
      --db-driver=mysql \
      --mysql-user=$MYSQL_USER \
      --mysql-password=$MYSQL_PASS \
      --mysql-db=$MYSQL_DB \
      --threads=$t \
      --time=$DURATION \
      oltp_read_write \
      run)

    # Extract metrics
    tps=$(echo "$output" | grep "transactions:" | awk '{print $2}' | tr -d ',' | tr -d ':')
    latency=$(echo "$output" | grep "avg:" | awk '{print $2}' | tr -d ',' | tr -d ':')
    errors=$(echo "$output" | grep "errors:" | awk '{print $2}' | tr -d ',' | tr -d ':')

    # Default to 0 if parsing fails
    tps=${tps:-0}
    latency=${latency:-0}
    errors=${errors:-0}

    # Sum up
    total_tps=$(echo "$total_tps + $tps" | bc)
    total_latency=$(echo "$total_latency + $latency" | bc)
    total_errors=$(echo "$total_errors + $errors" | bc)
  done

  # Calculate averages using bc with scale=2 for decimals
  avg_tps=$(echo "scale=2; $total_tps / $REPS" | bc)
  avg_latency=$(echo "scale=2; $total_latency / $REPS" | bc)
  avg_errors=$(echo "scale=2; $total_errors / $REPS" | bc)

  # Output to console and file
  echo "Threads: $t | Avg TPS: $avg_tps | Avg Latency: $avg_latency ms | Avg Errors: $avg_errors"

done

# Cleanup after test
sysbench \
  --db-driver=mysql \
  --mysql-user=$MYSQL_USER \
  --mysql-password=$MYSQL_PASS \
  --mysql-db=$MYSQL_DB \
  oltp_read_write \
  cleanup
