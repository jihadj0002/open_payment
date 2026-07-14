import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import { BASE_URL, HEADERS } from './config.js';

export const options = {
  stages: [
    { duration: '5m', target: 100 },
    { duration: '60m', target: 100 },
    { duration: '5m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95) < 600'],
    http_req_failed: ['rate < 0.02'],
  },
};

export default function () {
  const payload = JSON.stringify({
    amount: 1000,
    currency: 'BDT',
    payment_method: 'card',
    confirm: true,
  });
  const res = http.post(`${BASE_URL}/payments`, payload, { headers: HEADERS });
  check(res, { 'payment processed in soak': (r) => r.status === 201 });
  sleep(randomIntBetween(0.5, 1));
}
