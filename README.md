# hipay-go

Unofficial Go SDK for the [HiPay](https://hipay.mn) payment API. Supports checkout (E-Commerce), payment get/cancel/correction, and statement (POS) operations.

**Official API docs:** [https://developers.hipay.mn](https://developers.hipay.mn)

## Installation

```bash
go get github.com/techpartners-asia/hipay-go
```

## Initialize

```go
import "github.com/techpartners-asia/hipay-go/hipay"

client := hipay.New(
    "https://api.hipay.mn",  // HiPay API base URL (sandbox or production)
    "YOUR_TOKEN",            // API auth token
    "YOUR_ENTITY_ID",        // Merchant / entity ID
)
```

### Custom HTTP client

Use `WithClient` to inject your own `resty.Client` (e.g. for custom timeouts or TLS). Nil options are ignored.

```go
import (
    "time"
    "resty.dev/v3"
    "github.com/techpartners-asia/hipay-go/hipay"
)

httpClient := resty.New().SetTimeout(15 * time.Second)
client := hipay.New("https://api.hipay.mn", "TOKEN", "ENTITY_ID", hipay.WithClient(httpClient))
```

## Checkout

### Create checkout

```go
resp, err := client.Checkout(1000.0) // amount in MNT
if err != nil {
    log.Fatal(err)
}
if resp.Code != 1 {
    // API returned business error (handled via err from Checkout)
    return
}
fmt.Println(resp.CheckoutID)
fmt.Println(resp.QrData)
```

### Get checkout status

```go
resp, err := client.CheckoutGet("CHECKOUT_ID")
```

### Cancel checkout

```go
resp, err := client.CheckoutCancel("CHECKOUT_ID")
```

## Payment

### Get payment

```go
resp, err := client.PaymentGet("PAYMENT_ID")
```

### Payment correction

```go
resp, err := client.PaymentCorrection("PAYMENT_ID")
```

### Cancel payment

```go
resp, err := client.PaymentCancel("PAYMENT_ID")
```

## Statement

```go
resp, err := client.Statement("2024-01-15") // date format: YYYY-MM-DD
if err != nil {
    log.Fatal(err)
}
for _, row := range resp.Data.List {
    fmt.Println(row.PaymentID, row.Amount, row.PaymentDate)
}
```

## Error handling

All methods return a standard Go `error`. On HTTP or API errors (e.g. `Code != 1`), the error message includes the API description/message and optional detail (field/issue). Success responses with `Code == 1` have `Data` populated where applicable; for `Statement`, if the API returns success but no data, the SDK returns an error instead of a nil `Data` to avoid caller panics.

## HiPay API reference

- [HiPay developers — API docs](https://developers.hipay.mn) (Үндсэн API, E-Commerce, Pos, etc.)

## License

MIT
