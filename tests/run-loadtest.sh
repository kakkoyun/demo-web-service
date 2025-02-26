#!/bin/bash
set -e

# Function to check if command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Check if docker and docker-compose are installed
if ! command_exists docker || ! command_exists docker-compose; then
  echo "Error: Docker and docker-compose are required to run this script."
  echo "Please install Docker Desktop (macOS) or Docker Engine and docker-compose (Linux)."
  exit 1
fi

# Check if k6 is installed
if ! command_exists k6; then
  echo "k6 is not installed. You can install it with:"
  echo "macOS: brew install k6"
  echo "Linux: follow instructions at https://k6.io/docs/getting-started/installation/"
  echo "Or we can use the Docker version."
  USE_DOCKER_K6=true
else
  USE_DOCKER_K6=false
fi

# Change to the tests directory
cd "$(dirname "$0")"

# Start InfluxDB and Grafana
echo "Starting InfluxDB and Grafana..."
docker-compose up -d
echo "Waiting for services to initialize..."
sleep 5

# Get the base URL from command line argument or use default
BASE_URL=${1:-"http://host.docker.internal:8080"}

echo "Using API base URL: $BASE_URL"
echo "Note: If you're running on Linux, make sure your API is accessible from Docker containers"
echo "      or adjust the BASE_URL parameter."

# Run the load test
echo "Starting k6 load test..."
if [ "$USE_DOCKER_K6" = true ]; then
  # Use k6 in Docker
  docker run --rm -i \
    -v "${PWD}:/tests" \
    -e BASE_URL="$BASE_URL" \
    --network="host" \
    grafana/k6:latest run \
    --out=influxdb=http://localhost:8086/k6 \
    /tests/loadtest.js
else
  # Use locally installed k6
  k6 run \
    --out=influxdb=http://localhost:8086/k6 \
    -e BASE_URL="$BASE_URL" \
    loadtest.js
fi

echo ""
echo "Load test completed!"
echo ""
echo "View the Grafana dashboard at: http://localhost:3000"
echo ""
echo "Grafana is pre-configured with InfluxDB data source and a custom k6 dashboard."
echo "You should see the 'k6 Load Testing Results' dashboard automatically in the K6 folder."
echo "You can also import additional dashboards for k6:"
echo "1. Dashboard ID: 2587 (k6 Load Testing Results)"
echo "2. Dashboard ID: 4411 (k6 Load Testing Results by Endpoint)"
echo ""
echo "To stop the visualization services when finished:"
echo "docker-compose down" 