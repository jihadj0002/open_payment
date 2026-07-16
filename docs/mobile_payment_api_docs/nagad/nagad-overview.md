# Nagad Payment Gateway Overview

Nagad is a mobile financial service (MFS) and digital wallet operated in Bangladesh. It is one of Bangladesh's most widely used mobile wallets alongside bKash and Rocket.

## What is Nagad?

Nagad allows users to send money, pay bills, and make online purchases through the Nagad mobile app without needing a traditional bank card. For businesses, accepting Nagad at checkout means reaching customers who prefer mobile wallet payments over cards.

### Key Features

- State digital wallet of Bangladesh
- Large user base — tens of millions of citizens have accounts
- No bank card required for transactions
- Real-time transaction confirmation
- Secure data transmission via HTTPS
- Works even with weak internet connections
- QR code payment support
- PIN-based transaction confirmation (takes seconds)

## Benefits for Merchants

- **Increased conversion**: Payment through a familiar brand increases checkout completion
- **Reduced cart abandonment**: Offering Nagad reduces payment friction for Bangladeshi customers
- **Easy integration**: Simple REST API and ready-to-use plugins (WooCommerce, etc.)
- **No complex programming**: Integration through API or plugin does not require deep technical knowledge
- **Real-time processing**: Transactions are confirmed instantly
- **Marketing channel**: Access to in-app marketing opportunities

## Supported Platforms

Nagad can be integrated into:
- Websites
- Mobile apps (Flutter, native)
- E-commerce platforms (WooCommerce, Magento, OpenCart)
- Custom platforms via REST API

## Comparison

| Feature | Nagad | bKash | Rocket |
|---------|-------|-------|--------|
| Type | MFS / Digital Wallet | MFS | MFS |
| Online Payment API | Yes | Yes | Yes |
| User Base | Tens of millions | Large | Moderate |
| Sandbox Testing | Yes | Yes | Yes |

## Integration Options

1. **Flutter SDK** (`nagad_payment_gateway` package) — for Flutter mobile apps
2. **REST API** — full control over create/activate payment flow
3. **Payment plugin (iFrame)** — embed Nagad checkout via iFrame
4. **Postback** — configure callback URLs for transaction status notifications

## References

- [Nagad Flutter Package](https://pub.dev/packages/nagad_payment_gateway)
- [Nagad Official Website](https://www.nagad.com.bd/)
- [Nagad Online Payment API Integration Guide v3.3](https://github.com/muhibbin-munna/nagad_pg_php/blob/master/resource/Nagad%20Online%20Payment%20API%20Integration%20Guide%20v3.3.pdf)
