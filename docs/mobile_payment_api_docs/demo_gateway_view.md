Payment Gateway UI Architecture
1. Main Container
Full Screen/Page
Background
Desktop: Branded background image
Mobile: Solid dark background
Centered Payment Card
2. Header
Left Section
Back Button
Merchant Logo
Merchant Name
Transaction ID
Center (Desktop)
Countdown Timer
Right Section
Welcome Text (Desktop)
Offers Button
Login Button
3. Main Content

Desktop Layout

---------------------------------------------------------
| Left Panel               | Right Payment Panel         |
---------------------------------------------------------

Mobile Layout

Header
↓

Amount Card

↓

Payment Methods

↓

Selected Payment Method

↓

Payment Form

↓

Pay Button

↓

Footer
LEFT PANEL (Desktop)
A. Payment Summary Card

Contains:

"You are paying"
Total Amount
QR Payment Button
B. Cost Breakdown
Subtotal
Convenience Charge
Divider
Total Amount
C. Offers Section

Title

Special Offers & Savings

Description

Eligible offers

Content

Empty State Illustration
"No offers available"
D. Footer Links
Support
FAQ
Language Selector
RIGHT PANEL (Payment Area)
A. Welcome / User Area
Welcome Text
Login
B. Payment Method Tabs

Horizontal Tab Menu

Card
Mobile Banking
Net Banking
More

Selected tab highlighted.

C. Dynamic Payment Content

Changes according to selected payment option.

CARD PAYMENT STRUCTURE
Section Header
Pay with New Cards
Other Cards Link
Card Form
Card Number

Input Field

Supported card icons

Expiry + CVV Row

Left

MM/YY

Right

CVC/CVV

Card Preview Icon

Card Holder Name

Text Input

Remember Card

Checkbox

Info Icon

MOBILE BANKING STRUCTURE
Section Header
Pay with Mobile Banking
Mobile Banking Grid

Grid Layout

Each Item contains:

Provider Logo
Provider Name

Example structure:

Row 1
-------
Rocket
Nagad
Upay
CellFin

Row 2
-------
MCash
Pocket
Pathao Pay
Rainbow

Row 3
-------
Tap
OK Wallet
Islami Wallet
NET BANKING STRUCTURE
Section Header
Pay with Net Banking

Grid/List of Banks

Each Item

Bank Logo
Bank Name
MORE PAYMENT METHODS

Grouped Categories

Examples

EMI
Wallets
International Cards
Digital Wallets
Others
PAYMENT ACTION AREA

Large Primary Button

Pay ৳1,390.00

Disabled until

Payment method selected
Required fields completed
TERMS SECTION

Small Text

Agreement Statement
Terms of Service Link
FOOTER

Powered By

SSLCommerz Logo

Security Logos

PCI DSS
Security Certifications

Bottom Links

Support
FAQ
Language Selector
Mobile Page Flow
Header
    ↓
Payment Amount Card
    ↓
QR Payment Option
    ↓
Payment Method Tabs
    ↓
Selected Payment Method
    ↓
Payment Form / Provider Grid
    ↓
Pay Button
    ↓
Terms & Conditions
    ↓
Powered By
    ↓
Support + Language
Desktop Page Flow
Page
│
├── Header
│
├── Main Content
│   │
│   ├── Left Panel
│   │   ├── Payment Summary
│   │   ├── Cost Breakdown
│   │   ├── Offers
│   │   └── Footer Links
│   │
│   └── Right Panel
│       ├── Welcome/Login
│       ├── Payment Method Tabs
│       ├── Dynamic Payment Section
│       │   ├── Card Form
│       │   ├── Mobile Banking Grid
│       │   ├── Net Banking
│       │   └── More Methods
│       ├── Pay Button
│       ├── Terms & Conditions
│       └── Powered By / Security
│
└── Bottom Footer

This structure captures the overall information architecture shown in the screenshots while remaining implementation-agnostic, making it suitable as a wireframe or UI specification before moving into design or development.

You've reached today's image upload limit as a guest. Log in to continue uploading for free.
Log in
Sign up for free
draw those

I can do that, but image generation isn't available in this chat unless you're logged in. Once you're logged in, I can generate a clean wireframe image of the UI.

Meanwhile, here's an ASCII wireframe that represents the layout.

Desktop Wireframe
┌────────────────────────────────────────────────────────────────────────────────────────────┐
│ ← Logo Merchant Name                  04:37                     Welcome     Login          │
├───────────────────────────────────────┬────────────────────────────────────────────────────┤
│                                       │                                                    │
│ ┌───────────────────────────────┐     │  ┌──────┐ ┌────────────┐ ┌──────────┐ ┌────────┐ │
│ │ You are paying                │     │  │ Card │ │ MobileBank │ │ NetBank  │ │ More   │ │
│ │                               │     │  └──────┘ └────────────┘ └──────────┘ └────────┘ │
│ │        ৳1,390.00              │     │                                                    │
│ │                     QR Pay    │     │  Pay with New Cards                Other Cards     │
│ └───────────────────────────────┘     │                                                    │
│                                       │ ┌───────────────────────────────────────────────┐ │
│ Subtotal                 ৳1,390.00    │ │ Card Number                              VISA │ │
│ Convenience Charge        ৳0.00       │ ├───────────────────────┬───────────────────────┤ │
│ ───────────────────────────────────   │ │ MM/YY                 │ CVC/CVV              │ │
│ Total                    ৳1,390.00    │ ├───────────────────────────────────────────────┤ │
│                                       │ │ Card Holder Name                              │ │
│                                       │ └───────────────────────────────────────────────┘ │
│ Special Offers & Savings              │                                                  │
│ ───────────────────────────────────   │ □ Remember this card                            │
│                                       │                                                  │
│      [ Empty Illustration ]           │ ┌──────────────────────────────────────────────┐ │
│  No offers available                  │ │              PAY ৳1,390.00                   │ │
│                                       │ └──────────────────────────────────────────────┘ │
│                                       │ Terms & Conditions                              │
│ Support   FAQ        Language         │ Powered by SSLCommerz | Security Logos          │
└───────────────────────────────────────┴────────────────────────────────────────────────────┘
Mobile Wireframe
┌──────────────────────────────┐
│ ← Logo      Offers   Login   │
│ Merchant                  ID │
├──────────────────────────────┤
│ You are paying              │
│                              │
│        ৳1,390.00      QR     │
└──────────────────────────────┘

┌──────┬────────┬────────┬─────┐
│ Card │ Mobile │ NetBk  │More │
└──────┴────────┴────────┴─────┘

──────── Selected Content ────────

(Card)

┌──────────────────────────────┐
│ Card Number            VISA  │
├──────────────┬───────────────┤
│ MM/YY        │ CVC/CVV       │
├──────────────────────────────┤
│ Card Holder Name             │
└──────────────────────────────┘

□ Remember this card

┌──────────────────────────────┐
│      PAY ৳1,390.00           │
└──────────────────────────────┘

Terms & Conditions

Powered by SSLCommerz

Support     FAQ     Language
Mobile Banking Screen
┌──────────────────────────────┐
│ ← Logo      Offers   Login   │
├──────────────────────────────┤
│        ৳1,390.00      QR      │
└──────────────────────────────┘

┌──────┬────────┬────────┬─────┐
│ Card │ Mobile │ NetBk  │More │
└──────┴────────┴────────┴─────┘

Pay with Mobile Banking

┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐
│Logo │ │Logo │ │Logo │ │Logo │
│Bkash│ │Nagad│ │Upay │ │Cell │
└─────┘ └─────┘ └─────┘ └─────┘

┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐
│Logo │ │Logo │ │Logo │ │Logo │
│MCash│ │Pocket│Pathao│Rainbow│
└─────┘ └─────┘ └─────┘ └─────┘

┌─────┐ ┌─────┐ ┌─────┐
│Logo │ │Logo │ │Logo │
│ Tap │ │ OK  │ │ Isl │
└─────┘ └─────┘ └─────┘

┌──────────────────────────────┐
│      PAY ৳1,390.00           │
└──────────────────────────────┘

Terms & Conditions

Powered by SSLCommerz