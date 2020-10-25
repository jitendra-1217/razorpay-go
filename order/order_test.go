package order

import (
	"context"
	"testing"

	razorpay "github.com/jitendra-1217/razorpay-go"
	"github.com/jitendra-1217/razorpay-go/testutil"
	"github.com/stretchr/testify/assert"
)

var (
	// orderID holds new order id created in Create test.
	orderID string
)

func TestClient_Create(t *testing.T) {
	params := &razorpay.OrderParams{
		Amount:   razorpay.Int64(123),
		Currency: razorpay.String("INR"),
	}
	order, err := Create(context.Background(), params)
	// For use in later tests.
	orderID = order.ID
	assert.Nil(t, err)
	assert.True(t, testutil.IsAnyID(order.ID))
	assert.Equal(t, int64(123), order.Amount)
	assert.Equal(t, "INR", order.Currency)
}

func TestClient_Update(t *testing.T) {
	params := &razorpay.OrderParams{
		Notes: razorpay.Notes{
			"key-1": "value-1",
			"key-2": "value-2",
		},
	}
	order, err := Update(context.Background(), orderID, params)
	assert.Nil(t, err)
	assert.Equal(t, "value-1", order.Notes["key-1"])
	assert.Equal(t, "value-2", order.Notes["key-2"])
}

func TestClient_Get(t *testing.T) {
	order, err := Get(context.Background(), orderID, nil)
	assert.Nil(t, err)
	assert.Equal(t, orderID, order.ID)
}

func TestClient_List(t *testing.T) {
	params := &razorpay.OrderListParams{}
	params.Authorized = razorpay.String("1")
	params.Expand = []string{"payments", "transfers"}
	_, err := List(context.Background(), params)
	assert.Nil(t, err)
}

func TestClient_Payments(t *testing.T) {
	paymentList, err := Payments(context.Background(), "order_FtZ56sg2NgG0tX")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), paymentList.Count)
	assert.Equal(t, "pay_FtZSrSlgxsJKiQ", paymentList.Payments[0].ID)
}
