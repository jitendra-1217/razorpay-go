package refund

import (
	"context"
	"testing"

	faker "github.com/bxcodec/faker/v3"
	razorpay "github.com/jitendra-1217/razorpay-go"
	_ "github.com/jitendra-1217/razorpay-go/testutil"
	"github.com/stretchr/testify/assert"
)

var (
	// refundID holds an existing captured refund id.
	refundID = "rfnd_FtakiH6gO6Wehb"
)

func TestClient_Update(t *testing.T) {
	referenceUUID := faker.UUIDDigit()
	params := &razorpay.RefundUpdateParams{
		Notes: razorpay.Notes{
			"REFERENCE_UUID": referenceUUID,
		},
	}
	refund, err := Update(context.Background(), refundID, params)
	assert.Nil(t, err)
	assert.Equal(t, refundID, refund.ID)
	assert.Equal(t, referenceUUID, refund.Notes["REFERENCE_UUID"])
}

func TestClient_Get(t *testing.T) {
	refund, err := Get(context.Background(), refundID, nil)
	assert.Nil(t, err)
	assert.Equal(t, refundID, refund.ID)
}

func TestClient_List(t *testing.T) {
	refundList, err := List(context.Background(), nil)
	assert.Nil(t, err)
	assert.True(t, refundList.Count > 0)
}
