import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import { BASE_URL, HEADERS } from './config.js';

export const options = {
  stages: [
    { duration: '2m', target: 10 },   // ramp up
    { duration: '5m', target: 50 },   // stay at 50
    { duration: '2m', target: 100 },  // ramp to 100
    { duration: '3m', target: 100 },  // stay at 100
    { duration: '2m', target: 0 },    // ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95) < 500', 'p(99) < 1500'],
    http_req_failed: ['rate < 0.01'],
  },
};

export default function () {
  group('payment lifecycle', () => {
    // Create payment
    const createPayload = JSON.stringify({
      amount: randomIntBetween(100, 10000),
      currency: 'BDT',
      payment_method: 'card',
      confirm: true,
    });
    const createRes = http.post(`${BASE_URL}/payments`, createPayload, { headers: HEADERS });
    check(createRes, { 'payment created': (r) => r.status === 201 });

    // Get payment
    if (createRes.status === 201) {
      const paymentId = JSON.parse(createRes.body).data.id;
      const getRes = http.get(`${BASE_URL}/payments/${paymentId}`, { headers: HEADERS });
      check(getRes, { 'payment retrieved': (r) => r.status === 200 });
    }
  });

  group('balance check', () => {
    const res = http.get(`${BASE_URL}/balance`, { headers: HEADERS });
    check(res, { 'balance retrieved': (r) => r.status === 200 });
  });

  sleep(randomIntBetween(0.5, 2));
}
