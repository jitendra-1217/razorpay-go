package payment

import (
	"context"
	"errors"
	"testing"

	faker "github.com/bxcodec/faker/v3"
	razorpay "github.com/jitendra-1217/razorpay-go"
	_ "github.com/jitendra-1217/razorpay-go/testutil"
	"github.com/stretchr/testify/assert"
)

var (
	// paymentID holds an existing captured payment id. There exists a refund
	// on this payment as well.
	paymentID = "pay_FtZSrSlgxsJKiQ"
)

func TestClient_Update(t *testing.T) {
	referenceUUID := faker.UUIDDigit()
	params := &razorpay.PaymentUpdateParams{
		Notes: razorpay.Notes{
			"REFERENCE_UUID": referenceUUID,
		},
	}
	payment, err := Update(context.Background(), paymentID, params)
	assert.Nil(t, err)
	assert.Equal(t, paymentID, payment.ID)
	assert.Equal(t, referenceUUID, payment.Notes["REFERENCE_UUID"])
}

func TestClient_Capture(t *testing.T) {
	// Case: Attempts to capture already captured payment, hence expects error response.
	// Actually there is no way to create a new payment to capture either. So this is it.
	params := &razorpay.PaymentCaptureParams{
		Amount:   razorpay.Int64(123),
		Currency: razorpay.String("INR"),
	}
	_, err := Capture(context.Background(), paymentID, params)
	assert.NotNil(t, err)
	var razorpayErr *razorpay.Error
	assert.True(t, errors.As(err, &razorpayErr))
	assert.Equal(t, "BAD_REQUEST_ERROR", razorpayErr.Code)
	assert.Equal(t, "This payment has already been captured", razorpayErr.Description)
}

func TestClient_Get(t *testing.T) {
	payment, err := Get(context.Background(), paymentID, nil)
	assert.Nil(t, err)
	assert.Equal(t, paymentID, payment.ID)
}

func TestClient_List(t *testing.T) {
	_, err := List(context.Background(), nil)
	assert.Nil(t, err)
}

func TestClient_GetCard(t *testing.T) {
	card, err := GetCard(context.Background(), paymentID)
	assert.Nil(t, err)
	assert.Equal(t, "Kalidasa B", card.Name)
	assert.Equal(t, "1111", card.Last4)
	assert.Equal(t, "Visa", card.Network)
	assert.Equal(t, "debit", card.Type)
}

func TestClient_CreateRefund(t *testing.T) {
	// Case: Attempts to refund an existing payment. It must return error
	// response because the payment can not be refunded. There is no way to
	// create fresh payment and assert refund success.
	params := &razorpay.RefundCreateParams{
		Amount: razorpay.Int64(123),
	}
	_, err := CreateRefund(context.Background(), paymentID, params)
	assert.NotNil(t, err)
	var razorpayErr *razorpay.Error
	assert.True(t, errors.As(err, &razorpayErr))
	assert.Equal(t, "BAD_REQUEST_ERROR", razorpayErr.Code)
	assert.Equal(t, "The total refund amount is greater than the refund payment amount", razorpayErr.Description)
}

func TestClient_Refunds(t *testing.T) {
	refundList, err := Refunds(context.Background(), paymentID)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), refundList.Count)
	assert.Equal(t, "rfnd_FtakiH6gO6Wehb", refundList.Refunds[0].ID)
}
