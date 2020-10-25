package paymentlink

import (
	"context"
	"strings"
	"testing"

	faker "github.com/bxcodec/faker/v3"
	razorpay "github.com/jitendra-1217/razorpay-go"
	"github.com/jitendra-1217/razorpay-go/testutil"
	"github.com/stretchr/testify/assert"
)

var (
	// paymentLinkID holds new payment link id created in Create test.
	paymentLinkID string
)

func TestClient_Create(t *testing.T) {
	description := faker.Sentence()
	customerName := faker.Name()
	customerContact := faker.E164PhoneNumber()
	customerEmail := strings.ToLower(faker.Email())
	params := &razorpay.PaymentLinkParams{
		Amount:      razorpay.Int64(123),
		Currency:    razorpay.String("INR"),
		Description: &description,
		Customer: &razorpay.CustomerParams{
			Name:    &customerName,
			Contact: &customerContact,
			Email:   &customerEmail,
		},
		Notify: &razorpay.PaymentLinkNotifyParams{
			Email: razorpay.Bool(true),
		},
	}
	paymentLink, err := Create(context.Background(), params)
	// For use in later tests.
	paymentLinkID = paymentLink.ID
	assert.Nil(t, err)
	assert.True(t, testutil.IsAnyID(paymentLink.ID))
	assert.Equal(t, int64(123), paymentLink.Amount)
	assert.Equal(t, customerName, paymentLink.Customer.Name)
	assert.True(t, paymentLink.Notify.Email)
	assert.False(t, paymentLink.Notify.SMS)
}

func TestClient_Update(t *testing.T) {
	referenceID := faker.UUIDDigit()
	params := &razorpay.PaymentLinkParams{
		ReferenceID: &referenceID,
		Notes: razorpay.Notes{
			"key-1": "value-1",
			"key-2": "value-2",
		},
	}
	paymentLink, err := Update(context.Background(), paymentLinkID, params)
	assert.Nil(t, err)
	assert.Equal(t, "value-1", paymentLink.Notes["key-1"])
	assert.Equal(t, "value-2", paymentLink.Notes["key-2"])
	assert.Equal(t, referenceID, paymentLink.ReferenceID)
}

func TestClient_Get(t *testing.T) {
	paymentLink, err := Get(context.Background(), paymentLinkID, nil)
	assert.Nil(t, err)
	assert.Equal(t, paymentLinkID, paymentLink.ID)
}

// func TestClient_List(t *testing.T) {
// 	params := &razorpay.PaymentLinkListParams{}
// 	paymentLinkList, err := List(context.Background(), params)
// 	assert.Nil(t, err)
// 	assert.True(t, paymentLinkList.Count > 0)
// }

func TestClient_Notify(t *testing.T) {
	err := Notify(context.Background(), paymentLinkID, "email")
	assert.Nil(t, err)
}

func TestClient_Cancel(t *testing.T) {
	paymentLink, err := Cancel(context.Background(), paymentLinkID)
	assert.Nil(t, err)
	assert.Equal(t, "cancelled", paymentLink.Status)
}
