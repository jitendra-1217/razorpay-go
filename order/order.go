package order

import (
	"context"
	"net/http"

	razorpay "github.com/jitendra-1217/razorpay-go"
)

// Client is used to access /orders apis.
type Client struct {
	*razorpay.Client
}

// Create creates new order.
func (c *Client) Create(ctx context.Context, params *razorpay.OrderParams) (*razorpay.Order, error) {
	order := &razorpay.Order{}
	err := c.Call(ctx, http.MethodPost, "/orders", params, order)
	return order, err
}

// Update updates existing order.
func (c *Client) Update(ctx context.Context, id string, params *razorpay.OrderParams) (*razorpay.Order, error) {
	order := &razorpay.Order{}
	err := c.Call(ctx, http.MethodPatch, "/orders/"+id, params, order)
	return order, err
}

// Get returns order for id.
func (c *Client) Get(ctx context.Context, id string, params *razorpay.GetParams) (*razorpay.Order, error) {
	if params == nil {
		params = &razorpay.GetParams{}
	}

	order := &razorpay.Order{}
	err := c.Call(ctx, http.MethodGet, "/orders/"+id, params, order)
	return order, err
}

// List returns list of orders for params.
func (c *Client) List(ctx context.Context, params *razorpay.OrderListParams) (*razorpay.OrderList, error) {
	if params == nil {
		params = &razorpay.OrderListParams{}
	}

	orderList := &razorpay.OrderList{}
	err := c.Call(ctx, http.MethodGet, "/orders", params, orderList)
	return orderList, err
}

// Payments returns list of payments for order.
func (c *Client) Payments(ctx context.Context, orderID string) (*razorpay.PaymentList, error) {
	paymentList := &razorpay.PaymentList{}
	err := c.Call(ctx, http.MethodGet, "/orders/"+orderID+"/payments", nil, paymentList)
	return paymentList, err
}

// Create creates new order.
func Create(ctx context.Context, params *razorpay.OrderParams) (*razorpay.Order, error) {
	return getDefaultClient().Create(ctx, params)
}

// Update updates existing order.
func Update(ctx context.Context, id string, params *razorpay.OrderParams) (*razorpay.Order, error) {
	return getDefaultClient().Update(ctx, id, params)
}

// Get returns order for id.
func Get(ctx context.Context, id string, params *razorpay.GetParams) (*razorpay.Order, error) {
	return getDefaultClient().Get(ctx, id, params)
}

// List returns list of orders for params.
func List(ctx context.Context, params *razorpay.OrderListParams) (*razorpay.OrderList, error) {
	return getDefaultClient().List(ctx, params)
}

// Payments returns list of payments for order.
func Payments(ctx context.Context, orderID string) (*razorpay.PaymentList, error) {
	return getDefaultClient().Payments(ctx, orderID)
}

// NewClient returns new client.
func NewClient(apiKey string, apiSecret string, apiBackend razorpay.Backend) *Client {
	return &Client{razorpay.NewClient(apiKey, apiSecret, apiBackend)}
}

func getDefaultClient() *Client {
	return &Client{razorpay.GetDefaultClient()}
}
