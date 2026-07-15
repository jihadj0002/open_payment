import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

const baseUrl = __ENV.BASE_URL || 'http://localhost:8080/v1';
const token = __ENV.AUTH_TOKEN || 'test_token';

const paymentSuccessRate = new Rate('payment_success_rate');
const paymentDuration = new Trend('payment_duration');

export const options = {
  stages: [
    { duration: '30s', target: 50 },
    { duration: '1m', target: 50 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    payment_success_rate: ['rate>0.99'],
  },
};

const params = {
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
    'Idempotency-Key': '',
  },
};

export default function () {
  const idempotencyKey = `k6_${Date.now()}_${__VU}_${__ITER}`;
  params.headers['Idempotency-Key'] = idempotencyKey;

  const paymentPayload = JSON.stringify({
    amount: 1000 + (__VU * 100),
    currency: 'USD',
    payment_method: 'card',
    confirm: true,
    capture_method: 'automatic',
  });

  const createRes = http.post(`${baseUrl}/payments`, paymentPayload, params);
  const created = check(createRes, {
    'payment created': (r) => r.status === 201 || r.status === 200,
  });
  paymentSuccessRate.add(created);
  paymentDuration.add(createRes.timings.duration);
  sleep(0.5);

  if (createRes.status === 201 || createRes.status === 200) {
    const paymentId = createRes.json('data.id');

    const captureRes = http.post(`${baseUrl}/payments/${paymentId}/capture`, '{}', params);
    check(captureRes, {
      'payment captured': (r) => r.status === 200,
    });

    sleep(0.3);

    if (captureRes.status === 200) {
      const refundRes = http.post(`${baseUrl}/payments/${paymentId}/refund`, JSON.stringify({
        amount: 500,
        reason: 'k6 test refund',
      }), params);
      check(refundRes, {
        'payment refunded': (r) => r.status === 200,
      });
    }
  }

  sleep(1);
}
