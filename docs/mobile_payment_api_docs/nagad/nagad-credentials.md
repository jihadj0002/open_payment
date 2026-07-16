# NagadCredentials

The `NagadCredentials` class holds the merchant credentials required to authenticate with the Nagad Payment Gateway.

## Class Definition

```dart
class NagadCredentials {
  final String merchantID;
  final String merchantPrivateKey;
  final String pgPublicKey;
  final bool isSandbox;

  const NagadCredentials({
    required this.merchantID,
    required this.merchantPrivateKey,
    required this.pgPublicKey,
    required this.isSandbox,
  });
}
```

## Fields

| Field | Type | Description |
|-------|------|-------------|
| `merchantID` | `String` | Your unique merchant ID provided by Nagad |
| `merchantPrivateKey` | `String` | Your private key for signing requests |
| `pgPublicKey` | `String` | Nagad Payment Gateway's public key for verification |
| `isSandbox` | `bool` | Set to `true` for sandbox/testing, `false` for production |

## Obtaining Credentials

1. Register on the Nagad merchant portal
2. Receive your `merchantID`, `merchantPrivateKey`, and `pgPublicKey`
3. Provide your server IP address to Nagad for whitelisting
4. Configure your callback URL for payment notifications

## Example

```dart
final credentials = NagadCredentials(
  merchantID: "YOUR_MERCHANT_ID",
  merchantPrivateKey: "YOUR_MERCHANT_PRIVATE_KEY",
  pgPublicKey: "NAGAD_PAYMENT_GATEWAY_PUBLIC_KEY",
  isSandbox: true, // Set to false for production
);
```

## Security Warning

**Do not store merchant credentials directly in your Flutter application.** Use a secure backend server to handle payments securely and prevent exposure of sensitive data. Never store or expose your `merchantPrivateKey` or other sensitive data in the frontend.
