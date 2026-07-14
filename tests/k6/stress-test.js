import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import { BASE_URL, HEADERS } from './config.js';

export const options = {
  stages: [
    { duration: '2m', target: 50 },
    { duration: '5m', target: 200 },
    { duration: '2m', target: 300 },
    { duration: '3m', target: 300 },
    { duration: '2m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95) < 1000', 'p(99) < 3000'],
    http_req_failed: ['rate < 0.05'],
  },
};

export default function () {
  const payload = JSON.stringify({
    amount: randomIntBetween(100, 50000),
    currency: 'BDT',
    payment_method: 'card',
    confirm: true,
  });
  const res = http.post(`${BASE_URL}/payments`, payload, { headers: HEADERS });
  check(res, { 'payment created under stress': (r) => r.status === 201 });
  sleep(randomIntBetween(0.1, 0.5));
}
