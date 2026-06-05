import http from 'k6/http';
import { check, sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

// Determine the attack phase from environment variable. Default is 'smoke'.
const phase = __ENV.PHASE || 'smoke';

const scenarios = {
    smoke: {
        executor: 'constant-vus',
        vus: 5,
        duration: '1m',
    },
    average: {
        executor: 'ramping-vus',
        startVUs: 0,
        stages: [
            { duration: '1m', target: 500 }, // Ramp up
            { duration: '5m', target: 500 }, // Hold (Simulate busy day)
            { duration: '1m', target: 0 },   // Ramp down
        ],
    },
    stress: {
        executor: 'ramping-vus',
        startVUs: 0,
        stages: [
            { duration: '2m', target: 5000 }, // Aggressive ramp up
            { duration: '10m', target: 5000 },// Force backpressure (HTTP 429)
            { duration: '2m', target: 0 },
        ],
    },
    spike: {
        executor: 'ramping-vus',
        startVUs: 0,
        stages: [
            { duration: '10s', target: 10000 }, // Instant viral event
            { duration: '1m', target: 10000 },  // Hold
            { duration: '10s', target: 0 },     // Recover
        ],
    },
};

export const options = {
    scenarios: {
        [phase]: scenarios[phase],
    },
    thresholds: {
        // Goal: 95% of API requests must complete in under 50ms
        http_req_duration: ['p(95)<50'],
    }
};

export default function () {
    const url = 'http://localhost:8080/jobs'; // Adjust port if necessary

    // Dynamic Payload Generation
    // Bypasses Idempotency checks to ensure real throughput
    const payload = JSON.stringify({
        type: 'email_campaign',
        payload: {
            user_id: Math.floor(Math.random() * 1000000),
            message: 'Hello from k6 load test!'
        },
        priority: 'high'
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
            'Idempotency-Key': uuidv4(), // Unique key per request
            'X-Correlation-ID': uuidv4(),
            'Connection': 'keep-alive',
        },
    };

    const res = http.post(url, payload, params);

    // The Assertions
    check(res, {
        'is status 201 (Created)': (r) => r.status === 201,
        'is status 429 (Too Many Requests)': (r) => r.status === 429,
        'is status 500 (Internal Server Error)': (r) => r.status === 500,
    });
    
    // Simulates realistic user pacing. 1 VU = ~1 Request per second.
    sleep(1); 
}