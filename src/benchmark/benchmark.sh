#!/bin/bash

# Start Docker containers
echo "Starting Docker containers..."
docker-compose down
docker-compose up -d

# Define the MAX_INSERTS values
declare -a max_inserts_values=(10 100 1000 10000)

# Run the Go program 10 times for each MAX_INSERTS value
for max_inserts in "${max_inserts_values[@]}"; do
    echo "Benchmarking with MAX_INSERTS=$max_inserts..."

    for i in {1..10}; do
        # Export MAX_INSERTS environment variable
        export MAX_INSERTS=$max_inserts

        # Run the Go program and capture the output
        output=$(go run $(pwd)/main.go)

        # Extract MySQL and MongoDB times using grep and awk
        echo "$output"
    done
done

# Stop Docker containers
echo "Stopping Docker containers..."
docker-compose down
