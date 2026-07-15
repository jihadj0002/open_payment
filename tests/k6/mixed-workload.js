import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const baseUrl = __ENV.BASE_URL || 'http://localhost:8080/v1';
const token = __ENV.AUTH_TOKEN || 'test_token';

const errorRate = new Rate('errors');

export const options = {
  stages: [
    { duration: '1m', target: 100 },
    { duration: '3m', target: 100 },
    { duration: '1m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    errors: ['rate<0.05'],
  },
};

const params = {
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  },
};

const paymentIds = [];

export default function () {
  const scenario = __ITER % 4;

  switch (scenario) {
    case 0: {
      const idempotencyKey = `k6_create_${Date.now()}_${__VU}_${__ITER}`;
      const res = http.post(`${baseUrl}/payments`, JSON.stringify({
        amount: 2000 + (__VU * 50),
        currency: 'USD',
        payment_method: 'card',
        capture_method: 'manual',
      }), { ...params, headers: { ...params.headers, 'Idempotency-Key': idempotencyKey } });

      if (res.status === 201 || res.status === 200) {
        paymentIds.push(res.json('data.id'));
      }
      errorRate.add(res.status >= 400);
      break;
    }

    case 1: {
      const res = http.get(`${baseUrl}/payments?limit=10&offset=0`, params);
      errorRate.add(res.status >= 400);
      break;
    }

    case 2: {
      if (paymentIds.length > 0) {
        const id = paymentIds[__VU % paymentIds.length];
        const res = http.get(`${baseUrl}/payments/${id}`, params);
        errorRate.add(res.status >= 400);
      }
      break;
    }

    case 3: {
      const res = http.get(`${baseUrl}/balance`, params);
      errorRate.add(res.status >= 400);
      break;
    }
  }

  sleep(0.5);
}
