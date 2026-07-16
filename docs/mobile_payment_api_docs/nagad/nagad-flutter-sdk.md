# Nagad Flutter SDK

The `nagad_payment_gateway` package (v0.3.1) allows Flutter developers to integrate Nagad Online Payment into their applications.

## Installation

Add the dependency to your `pubspec.yaml`:

```yaml
dependencies:
  nagad_payment_gateway: ^0.3.1
```

## Initialization

### 1. Create Credentials

```dart
final credentials = NagadCredentials(
  merchantID: "YOUR_MERCHANT_ID",
  merchantPrivateKey: "YOUR_MERCHANT_PRIVATE_KEY",
  pgPublicKey: "NAGAD_PAYMENT_GATEWAY_PUBLIC_KEY",
  isSandbox: true, // Set to false for production
);
```

### 2. Create Nagad Instance

```dart
Nagad nagad = Nagad(
  credentials: credentials,
);
```

## Additional Merchant Info

You can provide additional details about your service to be displayed on the payment page.

### Merchant Info Fields

| Field Name | Max Length | Description |
|------------|-----------|-------------|
| `serviceName` | 25 | Service name provided by merchant |
| `serviceLogoURL` | 1~1024 | Publicly accessible logo URL |
| `additionalFieldNameEN` | 20 | Additional field name (English) |
| `additionalFieldNameBN` | 20 | Additional field name (Bangla) |
| `additionalFieldValue` | 20 | Value of additional field in English |

```dart
Map<String, dynamic> additionalMerchantInfo = {
  "serviceName": "T Shirt",
  "serviceLogoURL": "tinyurl.com/sampleLogoUrl",
  "additionalFieldNameEN": "Color",
  "additionalFieldNameBN": "রং",
  "additionalFieldValue": "White",
};

nagad.setAdditionalMerchantInfo(additionalMerchantInfo);
```

> **Note:** `additionalMerchantInfo` must be `Map<String, dynamic>`. Additional Merchant Info can be anything and will be saved for further usage. However only these fields will be shown in the payment page.

## Initiating a Payment

### Regular Payment

```dart
NagadResponse nagadResponse = await nagad.regularPayment(
  context,
  amount: 10.25,
  orderId: orderId,
);
```

### Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `context` | `BuildContext` | Flutter build context |
| `amount` | `double` | Payment amount |
| `orderId` | `dynamic` | Unique identifier for each payment |

#### Generating a Unique Order ID

```dart
DateTime now = DateTime.now();
String orderId = 'order${now.millisecondsSinceEpoch}';
```

## Nagad Class Reference

### Properties

| Property | Type | Description |
|----------|------|-------------|
| `credentials` | `NagadCredentials` | Merchant credentials |
| `client` | `Client` | HTTP client |
| `additionalMerchantInfo` | `Map<String, dynamic>` | Additional merchant info |

### Methods

| Method | Returns | Description |
|--------|---------|-------------|
| `regularPayment(BuildContext context, {required dynamic orderId, required double amount})` | `Future<NagadResponse>` | Initiate a regular payment |
| `setAdditionalMerchantInfo(Map<String, dynamic> info)` | `void` | Set additional merchant info |
| `generateRandomString(int size, Uint8List seed)` | `String` | Generate random string for signing |

## Testing in Sandbox Mode

To test payments in sandbox mode, ensure `isSandbox: true` when initializing `Nagad`. Use the test credentials provided by Nagad.

## Production Deployment

When moving to production:

1. Change `isSandbox: false`
2. Replace sandbox credentials with production credentials
3. Ensure your merchant account is verified by Nagad

## Security Best Practices

- **Use a backend server** to securely process payments and store credentials
- **Never store or expose** your `merchantPrivateKey` or other sensitive data in the frontend
- Always use HTTPS in production
