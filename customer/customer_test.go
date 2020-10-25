package customer

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
	// customerID holds new customer id created in Create test.
	customerID string
)

func TestClient_Create(t *testing.T) {
	name := faker.Name()
	contact := faker.E164PhoneNumber()
	email := strings.ToLower(faker.Email())
	username := faker.Username()
	params := &razorpay.CustomerParams{
		Name:    &name,
		Contact: &contact,
		Email:   &email,
		Notes: razorpay.Notes{
			"USERNAME": username,
		},
	}
	customer, err := Create(context.Background(), params)
	// For use in later tests.
	customerID = customer.ID
	assert.Nil(t, err)
	assert.True(t, testutil.IsAnyID(customer.ID))
	assert.Equal(t, name, customer.Name)
	assert.Equal(t, contact, customer.Contact)
	assert.Equal(t, email, customer.Email)
	assert.Equal(t, username, customer.Notes["USERNAME"])
}

func TestClient_Update(t *testing.T) {
	email := strings.ToLower(faker.Email())
	params := &razorpay.CustomerParams{
		Email: &email,
	}
	customer, err := Update(context.Background(), customerID, params)
	assert.Nil(t, err)
	assert.Equal(t, email, customer.Email)
}

func TestClient_Get(t *testing.T) {
	customer, err := Get(context.Background(), customerID, nil)
	assert.Nil(t, err)
	assert.Equal(t, customerID, customer.ID)
}

func TestClient_List(t *testing.T) {
	_, err := List(context.Background(), nil)
	assert.Nil(t, err)
}
