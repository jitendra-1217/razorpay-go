# razorpay-go
[![Test](https://github.com/jitendra-1217/razorpay-go/workflows/Test/badge.svg)](https://github.com/jitendra-1217/razorpay-go/actions?query=workflow%3ATest)

Go client for Razorpay APIs, refer https://razorpay.com/docs/api for reference.

## Todos
- [ ] Supports for APIs
    - [x] ~Customer~
    - [x] ~Order~
    - [x] ~Payment~
    - [x] ~Payment link~
    - [x] ~Refund~
    - [ ] Item
    - [ ] Invoice
    - [ ] Subscription
    - [ ] Settlement
    - [ ] Route
    - [ ] Smart collect
- [ ] Support for logging?

## Usage

### Using client

Import package and set api credentials:

```golang
import razorpay "github.com/jitendra-1217/razorpay-go"

razorpay.APIKey = "<KEY>"
razorpay.APISecret = "<SECRET>"
```

Import client per need. All clients e.g. customer, order, and payment etc
follow similar pattern.

For example, to capture a payment:

```golang
import razorpay_payment "github.com/jitendra-1217/razorpay-go/payment"

paymentID := "pay_00000000000001"
params := &razorpay.PaymentCaptureParams{
    Amount:   razorpay.Int64(123),
    Currency: razorpay.String("INR"),
}

payment, err := razorpay_payment.Capture(context.Background(), paymentID, params)
fmt.Println(payment.ID) // pay_00000000000001
fmt.Println(payment.Status) // captured
```

All param value are pointer so that only set values are sent in remote request
body.

### Handling errors

```golang
var razorpayErr *razorpay.Error
if errors.As(err, &razorpayErr) {
    // This could happen when an error response was received from remote.
    fmt.Println(razorpayErr.Code) // BAD_REQUEST_ERROR
    fmt.Println(razorpayErr.Description) // This payment has already been captured
} else {
    // This could happen for e.g. when error occurred in network connection etc.
    fmt.Println(err)
}
```

### Using other HTTP client

Any HTTP client satisfying Doer interface can be used instead of default one.

```golang
import "github.com/gojektech/heimdall/v6/httpclient"

timeout := 1000 * time.Millisecond
heimdallHTTPclient := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

razorpay.HTTPClient = heimdallHTTPclient
```

### Writing unit tests for integration code

The backend can be mocked in unit tests to assert request args and to receive
mocked response.

```golang
// First, to generate mock file out of Backend interface:
// mockgen -destination=backend_mock.go -package=<...> github.com/jitendra-1217/razorpay-go Backend

// Get mocked implementation and set default with the same.
mockAPIBackend := NewMockBackend(ctrl)
razorpay.DefaultAPIBackend = mockAPIBackend

// Sets expectation, ref https://github.com/golang/mock.
mockAPIBackend.
    EXPECT().
    Call(gomock.Eq(context.Background()), gomock.Eq("GET"), gomock.Eq("v1/payments/pay_00000000000001"), gomock.Any(), gomock.Any()).
    Return(nil).
    SetArg(4, razorpay.Payment{Entity: razorpay.Entity{ID: "pay_00000000000001"}})

// Follows.. existing code call and assertions.
payment, err := razorpay_payment.Get(context.Background(), paymentID, nil)
assert.Nil(t, err)
assert.Equal(t, "pay_00000000000001", payment.ID)
```

### Unmarshalling response into own struct

Every response struct has the corresponding raw body set in Body field that can
be unmarshalled into own  struct.

```golang
payment, err := razorpay_payment.Get(context.Background(), "pay_00000000000001", nil)

// Defines customPaymentResponse with specific fields and their json mapping.
customPaymentResponse := &struct {
    ID     string `json:"id"`
    Status string `json:"status"`
}{}
_ = json.Unmarshal(payment.Body, customPaymentResponse)
```

### Using multiple clients with separate api credentials

```golang
paymentClient1 := razorpay_payment.NewClient("<KEY-1>", "<SECRET-1>", nil)
payment1, err := paymentClient1.Get(context.Background(), "pay_00000000000001", nil)

paymentClient2 := razorpay_payment.NewClient("<KEY-2>", "<SECRET-2>", nil)
payment2, err := paymentClient2.Get(context.Background(), "pay_00000000000002", nil)
```

### Metrics instrumentation with Prometheus

A prometheus collector for default http client exists for use. When using own
http client, a helper function i.e. NewPrometheusCollector exists for use.

```golang
prometheus.MustRegister(razorpay.HTTPClientPrometheusCollector)

// Or when using own http client...
collector := razorpay.NewPrometheusCollector(heimdallHTTPclient)
prometheus.MustRegister(collector)

// And to serve /metrics, for example...
http.Handle("/metrics", promhttp.Handler())
http.ListenAndServe(":8080", nil)
```
