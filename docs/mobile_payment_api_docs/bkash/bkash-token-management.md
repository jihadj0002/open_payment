# bKash - Token Management

bKash API access tokens are credentials required to access all payment APIs. Tokens inform the API that the bearer has been permitted to access the API with specific scopes.

## Token Lifecycle

1. Call **Grant Token API** to get an authorization token
2. Use the token in the `Authorization` header for all API calls
3. Before the token expires (at the 50th/55th minute), call **Refresh Token API**
4. Use the new token for subsequent API calls

## Token Specifications

- Default lifetime: **3600 seconds** (1 hour)
- Do not call Refresh Token more than **2 times within an hour**
- Exceeding the limit returns an error and blocks the merchant for **1 hour**

## Grant Token

Gets a new authorization token.

**Request URL:** `{base_URL}/tokenized/checkout/token/grant`

**Method:** POST

### Request Headers

| Header | Value |
|--------|-------|
| Content-Type | application/json |
| Accept | application/json |
| username | Username shared by bKash during onboarding |
| password | Password shared by bKash during onboarding |

### Request Body

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| app_key | string | Yes | App key shared during onboarding |
| app_secret | string | Yes | App secret shared during onboarding |

### Success Response

| Parameter | Type | Description |
|-----------|------|-------------|
| id_token | string | The authorization token (JWT) |
| token_type | string | Token type (e.g. Bearer) |
| expires_in | int | Token lifetime in seconds |

## Refresh Token

Refreshes an existing token before it expires.

**Request URL:** `{base_URL}/tokenized/checkout/token/refresh`

**Method:** POST

### Request Headers

| Header | Value |
|--------|-------|
| Content-Type | application/json |
| Accept | application/json |

### Request Body

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| app_key | string | Yes | App key shared during onboarding |
| app_secret | string | Yes | App secret shared during onboarding |
| refresh_token | string | Yes | Existing refresh token |

### Success Response

Returns new id_token and refresh_token values.
