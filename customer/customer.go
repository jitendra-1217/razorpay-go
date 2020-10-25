package customer

import (
	"context"
	"net/http"

	razorpay "github.com/jitendra-1217/razorpay-go"
)

// Client is used to access /customers apis.
type Client struct {
	*razorpay.Client
}

// Create creates new customer.
func (c *Client) Create(ctx context.Context, params *razorpay.CustomerParams) (*razorpay.Customer, error) {
	customer := &razorpay.Customer{}
	err := c.Call(ctx, http.MethodPost, "/customers", params, customer)
	return customer, err
}

// Update updates existing customer.
func (c *Client) Update(ctx context.Context, id string, params *razorpay.CustomerParams) (*razorpay.Customer, error) {
	customer := &razorpay.Customer{}
	err := c.Call(ctx, http.MethodPut, "/customers/"+id, params, customer)
	return customer, err
}

// Get returns customer for id.
func (c *Client) Get(ctx context.Context, id string, params *razorpay.GetParams) (*razorpay.Customer, error) {
	if params == nil {
		params = &razorpay.GetParams{}
	}

	customer := &razorpay.Customer{}
	err := c.Call(ctx, http.MethodGet, "/customers/"+id, params, customer)
	return customer, err
}

// List returns list of customers for params.
func (c *Client) List(ctx context.Context, params *razorpay.CustomerListParams) (*razorpay.CustomerList, error) {
	if params == nil {
		params = &razorpay.CustomerListParams{}
	}

	customerList := &razorpay.CustomerList{}
	err := c.Call(ctx, http.MethodGet, "/customers", params, customerList)
	return customerList, err
}

// Create creates new customer.
func Create(ctx context.Context, params *razorpay.CustomerParams) (*razorpay.Customer, error) {
	return getDefaultClient().Create(ctx, params)
}

// Update updates existing customer.
func Update(ctx context.Context, id string, params *razorpay.CustomerParams) (*razorpay.Customer, error) {
	return getDefaultClient().Update(ctx, id, params)
}

// Get returns customer for id.
func Get(ctx context.Context, id string, params *razorpay.GetParams) (*razorpay.Customer, error) {
	return getDefaultClient().Get(ctx, id, params)
}

// List returns list of customers for params.
func List(ctx context.Context, params *razorpay.CustomerListParams) (*razorpay.CustomerList, error) {
	return getDefaultClient().List(ctx, params)
}

// NewClient returns new client.
func NewClient(apiKey string, apiSecret string, apiBackend razorpay.Backend) *Client {
	return &Client{razorpay.NewClient(apiKey, apiSecret, apiBackend)}
}

func getDefaultClient() *Client {
	return &Client{razorpay.GetDefaultClient()}
}
