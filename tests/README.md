# Load Testing Setup

This directory contains all the components needed to run load tests against the API service using k6 and visualize the results with InfluxDB and Grafana.

## Prerequisites

- Docker and Docker Compose
- k6 (optional, can run through Docker)

## Components

- `loadtest.js`: The k6 script that defines the load test scenarios and metrics
- `docker-compose.yml`: Configuration to run InfluxDB and Grafana
- `run-loadtest.sh`: Shell script to automate the setup and execution of the load test
- `grafana/provisioning/`: Auto-configuration for Grafana data sources and dashboards

## Running Load Tests

### Quick Start

1. Make sure your API is running (typically on port 8080)
2. From the project root, run:
   ```
   cd tests
   ./run-loadtest.sh
   ```
   
   Or specify a custom API URL:
   ```
   ./run-loadtest.sh http://localhost:8080
   ```

### Viewing Results

- InfluxDB: http://localhost:8086 (admin/admin)
- Grafana Dashboard: http://localhost:3000 (pre-configured with data source)

The Grafana dashboard is pre-configured with:
- InfluxDB data source
- Folder for k6 dashboards
- A custom k6 dashboard that shows:
  - HTTP request duration by endpoint
  - Total requests
  - Failed request rate
  - Virtual users
  - Request rate
  - Requests by endpoint
  - HTTP status codes

You can also import additional dashboards for k6:
- Dashboard ID: 2587 (k6 Load Testing Results)
- Dashboard ID: 4411 (k6 Load Testing Results by Endpoint)

## Key Metrics

The load test captures several important metrics:

- `http_reqs`: Counter for total number of requests
- `http_req_duration`: Trend metric for request duration
- `http_req_failed`: Rate metric for failed requests
- Custom metrics defined in the loadtest.js file:
  - `http_requests_success`: Counter for successful requests
  - `http_requests_failed`: Counter for failed requests
  - `http_request_duration`: Trend metric for request duration by endpoint

The metrics are automatically stored in InfluxDB with tags for:
- scenario: The test scenario name
- status: HTTP status code
- method: HTTP method (GET, POST, etc.)
- name: Request name or URL path
- group: Test group name
- endpoint: Custom tag to identify specific API endpoints

## Test Scenarios

The load test script includes the following test scenarios:

1. **Health Check** (GET /api/health)
   - Verifies the API health endpoint
   - Expects 200 OK and `{"status":"healthy"}` response
   - Target p99 response time: 100ms

2. **Version Check** (GET /api/version)
   - Checks the version information endpoint
   - Expects 200 OK and a response with version information
   - Target p99 response time: 200ms

3. **Home Page** (GET /)
   - Tests the main welcome page
   - Expects 200 OK and a welcome message
   - Target p99 response time: 200ms

4. **List Users** (GET /api/users)
   - Retrieves the list of all users
   - Expects 200 OK and a list of user objects
   - Target p99 response time: 250ms

5. **Get User by ID** (GET /api/users/{id})
   - Tests both valid and invalid user IDs
   - For valid IDs: expects 200 OK and user data
   - For invalid IDs: expects 400 Bad Request 
   - Target p99 response time: 250ms

6. **Create User** (POST /api/users)
   - Creates new users with random names
   - Expects 201 Created and confirmation message
   - Target p99 response time: 300ms

## Load Profile

The tests execute with the following load profile:

- 30 seconds: Ramp up from 0 to 20 virtual users
- 1 minute: Steady load with 20 virtual users
- 30 seconds: Ramp down from 20 to 0 virtual users

This profile tests both the API's ability to handle gradually increasing load and sustained traffic.

## Customizing Tests

To modify the test scenarios or thresholds:

1. Edit `loadtest.js` to change the test stages, thresholds, or endpoints
2. For more complex scenarios, refer to the [k6 documentation](https://k6.io/docs/)

## Cleaning Up

When you're done with testing, stop the containers:

```
docker-compose down
```

## Advantages of InfluxDB for k6

- Native integration with k6 without additional components
- Time-series database optimized for metrics and monitoring data
- Simple configuration with minimal setup required
- Well-established integration with Grafana for dashboarding
- Built-in data retention policies and continuous queries
- Easy to use with k6 through the `--out=influxdb` flag

## Troubleshooting

- **macOS**: The script uses `host.docker.internal` to access the API from containers
- **Linux**: You may need to adjust the BASE_URL parameter or network configuration 