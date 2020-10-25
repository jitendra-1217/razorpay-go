package refund

import (
	"context"
	"net/http"

	razorpay "github.com/jitendra-1217/razorpay-go"
)

// Client is used to access /refunds apis.
type Client struct {
	*razorpay.Client
}

// Update updates existing refund.
func (c *Client) Update(ctx context.Context, id string, params *razorpay.RefundUpdateParams) (*razorpay.Refund, error) {
	refund := &razorpay.Refund{}
	err := c.Call(ctx, http.MethodPatch, "/refunds/"+id, params, refund)
	return refund, err
}

// Get returns refund for id.
func (c *Client) Get(ctx context.Context, id string, params *razorpay.GetParams) (*razorpay.Refund, error) {
	if params == nil {
		params = &razorpay.GetParams{}
	}

	refund := &razorpay.Refund{}
	err := c.Call(ctx, http.MethodGet, "/refunds/"+id, params, refund)
	return refund, err
}

// List returns list of refunds for params.
func (c *Client) List(ctx context.Context, params *razorpay.RefundListParams) (*razorpay.RefundList, error) {
	if params == nil {
		params = &razorpay.RefundListParams{}
	}

	refundList := &razorpay.RefundList{}
	err := c.Call(ctx, http.MethodGet, "/refunds", params, refundList)
	return refundList, err
}

// Update updates existing refund.
func Update(ctx context.Context, id string, params *razorpay.RefundUpdateParams) (*razorpay.Refund, error) {
	return getDefaultClient().Update(ctx, id, params)
}

// Get returns refund for id.
func Get(ctx context.Context, id string, params *razorpay.GetParams) (*razorpay.Refund, error) {
	return getDefaultClient().Get(ctx, id, params)
}

// List returns list of refunds for params.
func List(ctx context.Context, params *razorpay.RefundListParams) (*razorpay.RefundList, error) {
	return getDefaultClient().List(ctx, params)
}

// NewClient returns new client.
func NewClient(apiKey string, apiSecret string, apiBackend razorpay.Backend) *Client {
	return &Client{razorpay.NewClient(apiKey, apiSecret, apiBackend)}
}

func getDefaultClient() *Client {
	return &Client{razorpay.GetDefaultClient()}
}
