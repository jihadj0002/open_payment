INSERT INTO merchants (id, name, email, secret_key, public_key, webhook_url, status)
VALUES
    ('a1b2c3d4-e5f6-7890-abcd-ef1234567890', 'Demo Merchant', 'merchant@demo.com',
     'sk_test_demo_secret_key_1234567890', 'pk_test_demo_public_key_1234567890',
     'https://webhook.demo.com/hooks', 'active'),
    ('b2c3d4e5-f6a7-8901-bcde-f12345678901', 'Test Shop', 'admin@testshop.com',
     'sk_test_testshop_secret_abcdef123456', 'pk_test_testshop_public_abcdef123456',
     'https://testshop.com/webhook', 'active');

INSERT INTO customers (merchant_id, email, phone, name)
VALUES
    ('a1b2c3d4-e5f6-7890-abcd-ef1234567890', 'customer1@example.com', '+8801712345678', 'Rahim Ahmed'),
    ('a1b2c3d4-e5f6-7890-abcd-ef1234567890', 'customer2@example.com', '+8801812345678', 'Karim Hasan'),
    ('b2c3d4e5-f6a7-8901-bcde-f12345678901', 'shop_customer@example.com', '+8801912345678', 'Sultana Begum');

INSERT INTO api_keys (merchant_id, key_prefix, key_hash, name, permissions)
VALUES
    ('a1b2c3d4-e5f6-7890-abcd-ef1234567890', 'sk_live_', 'hashed_live_key_1', 'Production Key', '["read","write"]'),
    ('a1b2c3d4-e5f6-7890-abcd-ef1234567890', 'sk_test_', 'hashed_test_key_1', 'Test Key', '["read","write"]');
