package razorpay

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Sets api key, secret and host for tests, reading from env values.
func init() {
	APIKey = os.Getenv("TEST_API_KEY")
	APISecret = os.Getenv("TEST_API_SECRET")
	APIHost = os.Getenv("TEST_API_HOST")
}

func TestIsValidPaymentSignature(t *testing.T) {
	// Case: Positive.
	params := map[string]string{
		"razorpay_order_id":   "order_FtH5WfFwLE2Rsu",
		"razorpay_payment_id": "pay_FtHnZWyFUcWO2p",
		"razorpay_signature":  "1d82f20d3c967ba3d13837b9853cad26b0616c34cfeb42beba22dde46651f875",
	}
	isValid, err := IsValidPaymentSignature(context.Background(), params)
	assert.True(t, isValid)
	assert.Nil(t, err)

	// Case: When api secret is set to incorrect value.
	client := NewClient("rzp_test_1DP5mmOlF5G5ag", "INCORRECT_SECRET", nil)
	isValid, err = client.IsValidPaymentSignature(context.Background(), params)
	assert.False(t, isValid)
	assert.Nil(t, err)

	// Case: When params are malformed.
	params["razorpay_order_id"] = "order_INCORRECT_ID"
	isValid, err = IsValidPaymentSignature(context.Background(), params)
	assert.False(t, isValid)
	assert.Nil(t, err)
}

func TestIsValidWebhookRequest(t *testing.T) {
	// Builds mocked webhook request.
	var body io.Reader = bytes.NewBuffer([]byte("{\"entity\":\"event\",\"event\":\"order.paid\"}"))
	req, err := http.NewRequest(http.MethodPost, "example.com/webhook", body)
	assert.Nil(t, err)

	// Case: When request does not contain signature.
	isValid, err := IsValidWebhookRequest(context.Background(), req, "WEBHOOK_SECRET")
	assert.False(t, isValid)
	assert.Nil(t, err)

	req.Header.Set("X-Razorpay-Signature", "f95146f045329385bc5710008367443e2fa390cfc79024621233281cb6612cd7")

	// Case: Positive.
	isValid, err = IsValidWebhookRequest(context.Background(), req, "WEBHOOK_SECRET")
	assert.True(t, isValid)
	assert.Nil(t, err)

	// Case: When incorrect webhook secret is used.
	isValid, err = IsValidWebhookRequest(context.Background(), req, "WEBHOOK_INCORRECT_SECRET")
	assert.False(t, isValid)
	assert.Nil(t, err)
}
