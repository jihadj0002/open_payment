import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { BASE_URL, HEADERS } from './config.js';

export const options = {
  vus: 1,
  duration: '30s',
  thresholds: {
    http_req_duration: ['p(95) < 500'],
    http_req_failed: ['rate < 0.01'],
  },
};

export default function () {
  group('health check', () => {
    const res = http.get(`${BASE_URL}/health`);
    check(res, { 'health status is 200': (r) => r.status === 200 });
  });

  group('create payment', () => {
    const payload = JSON.stringify({
      amount: 1000,
      currency: 'BDT',
      payment_method: 'card',
      description: 'k6 smoke test payment',
    });
    const res = http.post(`${BASE_URL}/payments`, payload, { headers: HEADERS });
    check(res, {
      'payment created status is 201': (r) => r.status === 201,
      'payment has id': (r) => JSON.parse(r.body).data.id !== undefined,
    });
  });

  sleep(1);
}
