import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Trend } from 'k6/metrics';

// Custom metrics
const httpRequestsSuccess = new Counter('http_requests_success');
const httpRequestsFailed = new Counter('http_requests_failed');
const httpRequestDuration = new Trend('http_request_duration');

// Configuration
export const options = {
  // Stages for the test (e.g., ramp-up, steady, ramp-down)
  stages: [
    { duration: '30s', target: 20 }, // Ramp up to 20 users over 30 seconds
    { duration: '1m', target: 20 },  // Stay at 20 users for 1 minute
    { duration: '30s', target: 0 },  // Ramp down to 0 users over 30 seconds
  ],
  // Thresholds
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
    'http_request_duration{endpoint:health}': ['p(99)<100'], // 100ms
    'http_request_duration{endpoint:version}': ['p(99)<200'], // 200ms
    'http_request_duration{endpoint:home}': ['p(99)<200'], // 200ms
    'http_request_duration{endpoint:users}': ['p(99)<250'], // 250ms
    'http_request_duration{endpoint:user_by_id}': ['p(99)<250'], // 250ms
    'http_request_duration{endpoint:create_user}': ['p(99)<300'], // 300ms
    'http_requests_failed': ['count<10'], // Fewer than 10 failed requests
  },
};

// Environment variables
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function() {
  // Test case 1: Health check endpoint
  const healthRes = http.get(`${BASE_URL}/api/health`, {
    tags: { endpoint: 'health' }
  });
  
  // Record custom metrics
  httpRequestDuration.add(healthRes.timings.duration, { endpoint: 'health' });
  
  // Check health response
  const healthSuccess = check(healthRes, {
    'health status is 200': (r) => r.status === 200,
    'health response has correct format': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.status === 'healthy';
      } catch (e) {
        return false;
      }
    },
  });
  
  if (healthSuccess) {
    httpRequestsSuccess.add(1, { endpoint: 'health' });
  } else {
    httpRequestsFailed.add(1, { endpoint: 'health' });
  }
  
  // Short sleep between requests
  sleep(1);
  
  // Test case 2: Version endpoint
  const versionRes = http.get(`${BASE_URL}/api/version`, {
    tags: { endpoint: 'version' }
  });
  
  // Record custom metrics
  httpRequestDuration.add(versionRes.timings.duration, { endpoint: 'version' });
  
  // Check version response
  const versionSuccess = check(versionRes, {
    'version status is 200': (r) => r.status === 200,
    'version response has correct format': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.version !== undefined;
      } catch (e) {
        return false;
      }
    },
  });
  
  if (versionSuccess) {
    httpRequestsSuccess.add(1, { endpoint: 'version' });
  } else {
    httpRequestsFailed.add(1, { endpoint: 'version' });
  }
  
  // Short sleep between requests
  sleep(1);
  
  // Test case 3: Home endpoint
  const homeRes = http.get(`${BASE_URL}/`, {
    tags: { endpoint: 'home' }
  });
  
  // Record custom metrics
  httpRequestDuration.add(homeRes.timings.duration, { endpoint: 'home' });
  
  // Check home response
  const homeSuccess = check(homeRes, {
    'home status is 200': (r) => r.status === 200,
    'home response has correct format': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.message === 'Welcome to the API';
      } catch (e) {
        return false;
      }
    },
  });
  
  if (homeSuccess) {
    httpRequestsSuccess.add(1, { endpoint: 'home' });
  } else {
    httpRequestsFailed.add(1, { endpoint: 'home' });
  }
  
  // Short sleep between requests
  sleep(1);
  
  // Test case 4: Get all users
  const usersRes = http.get(`${BASE_URL}/api/users`, {
    tags: { endpoint: 'users' }
  });
  
  // Record custom metrics
  httpRequestDuration.add(usersRes.timings.duration, { endpoint: 'users' });
  
  // Check users response
  const usersSuccess = check(usersRes, {
    'users status is 200': (r) => r.status === 200,
    'users response has correct format': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.status === 'success' && Array.isArray(body.users);
      } catch (e) {
        return false;
      }
    },
  });
  
  if (usersSuccess) {
    httpRequestsSuccess.add(1, { endpoint: 'users' });
  } else {
    httpRequestsFailed.add(1, { endpoint: 'users' });
  }
  
  // Short sleep between requests
  sleep(1);
  
  // Test case 5: Get user by ID
  // Try with both valid and invalid IDs
  const userIds = [1, 2, 'invalid'];
  const randomIndex = Math.floor(Math.random() * userIds.length);
  const userId = userIds[randomIndex];
  
  const userByIdRes = http.get(`${BASE_URL}/api/users/${userId}`, {
    tags: { endpoint: 'user_by_id' }
  });
  
  // Record custom metrics
  httpRequestDuration.add(userByIdRes.timings.duration, { endpoint: 'user_by_id' });
  
  // Expected status depends on whether we're testing a valid or invalid ID
  const expectedStatus = userId === 'invalid' ? 400 : 200;
  
  // Check user by ID response
  const userByIdSuccess = check(userByIdRes, {
    [`user_by_id status is ${expectedStatus}`]: (r) => r.status === expectedStatus,
    'user_by_id response has correct format': (r) => {
      try {
        const body = JSON.parse(r.body);
        if (userId === 'invalid') {
          return body.status === 'error';
        } else {
          return body.status === 'success' && body.user && body.user.id === userId;
        }
      } catch (e) {
        return false;
      }
    },
  });
  
  if (userByIdSuccess) {
    httpRequestsSuccess.add(1, { endpoint: 'user_by_id' });
  } else {
    httpRequestsFailed.add(1, { endpoint: 'user_by_id' });
  }
  
  // Short sleep between requests
  sleep(1);
  
  // Test case 6: Create new user
  const payload = JSON.stringify({
    name: `Load Test User ${Math.floor(Math.random() * 1000)}`
  });
  
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
    tags: { endpoint: 'create_user' }
  };
  
  const createUserRes = http.post(
    `${BASE_URL}/api/users`,
    payload,
    params
  );
  
  // Record custom metrics
  httpRequestDuration.add(createUserRes.timings.duration, { endpoint: 'create_user' });
  
  // Check create user response
  const createUserSuccess = check(createUserRes, {
    'create_user status is 201': (r) => r.status === 201,
    'create_user response has correct format': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.status === 'success' && 
               body.message === 'User created successfully' && 
               body.user && 
               body.user.name;
      } catch (e) {
        return false;
      }
    },
  });
  
  if (createUserSuccess) {
    httpRequestsSuccess.add(1, { endpoint: 'create_user' });
  } else {
    httpRequestsFailed.add(1, { endpoint: 'create_user' });
  }
  
  // Sleep between iterations
  sleep(2);
} 