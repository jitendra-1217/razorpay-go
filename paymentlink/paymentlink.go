package paymentlink

import (
	"context"
	"net/http"

	razorpay "github.com/jitendra-1217/razorpay-go"
)

// Client is used to access /payment-links apis.
type Client struct {
	*razorpay.Client
}

// Create creates new payment link.
func (c *Client) Create(ctx context.Context, params *razorpay.PaymentLinkParams) (*razorpay.PaymentLink, error) {
	paymentLink := &razorpay.PaymentLink{}
	err := c.Call(ctx, http.MethodPost, "/payment_links", params, paymentLink)
	return paymentLink, err
}

// Update updates existing payment link.
func (c *Client) Update(ctx context.Context, id string, params *razorpay.PaymentLinkParams) (*razorpay.PaymentLink, error) {
	paymentLink := &razorpay.PaymentLink{}
	err := c.Call(ctx, http.MethodPatch, "/payment_links/"+id, params, paymentLink)
	return paymentLink, err
}

// Get returns payment link for id.
func (c *Client) Get(ctx context.Context, id string, params *razorpay.GetParams) (*razorpay.PaymentLink, error) {
	if params == nil {
		params = &razorpay.GetParams{}
	}

	paymentLink := &razorpay.PaymentLink{}
	err := c.Call(ctx, http.MethodGet, "/payment_links/"+id, params, paymentLink)
	return paymentLink, err
}

// // List returns list of payment links for params.
// TODO: Uncomment after handling custom unmarshalling. The response in this case is not same as in common contract.
// func (c *Client) List(ctx context.Context, params *razorpay.PaymentLinkListParams) (*razorpay.PaymentLinkList, error) {
// 	if params == nil {
// 		params = &razorpay.PaymentLinkListParams{}
// 	}

// 	paymentLinkList := &razorpay.PaymentLinkList{}
// 	err := c.Call(ctx, http.MethodGet, "/payment_links", params, paymentLinkList)
// 	return paymentLinkList, err
// }

// Notify sends or resends notifications for payment link.
func (c *Client) Notify(ctx context.Context, id string, medium string) error {
	return c.Call(ctx, http.MethodPost, "/payment_links/"+id+"/notify_by/"+medium, nil, nil)
}

// Cancel cancels payment link.
func (c *Client) Cancel(ctx context.Context, id string) (*razorpay.PaymentLink, error) {
	paymentLink := &razorpay.PaymentLink{}
	err := c.Call(ctx, http.MethodPost, "/payment_links/"+id+"/cancel", nil, paymentLink)
	return paymentLink, err
}

// Create creates new payment link.
func Create(ctx context.Context, params *razorpay.PaymentLinkParams) (*razorpay.PaymentLink, error) {
	return getDefaultClient().Create(ctx, params)
}

// Update updates existing payment link.
func Update(ctx context.Context, id string, params *razorpay.PaymentLinkParams) (*razorpay.PaymentLink, error) {
	return getDefaultClient().Update(ctx, id, params)
}

// Get returns payment link for id.
func Get(ctx context.Context, id string, params *razorpay.GetParams) (*razorpay.PaymentLink, error) {
	return getDefaultClient().Get(ctx, id, params)
}

// // List returns list of payment links for params.
// func List(ctx context.Context, params *razorpay.PaymentLinkListParams) (*razorpay.PaymentLinkList, error) {
// 	return getDefaultClient().List(ctx, params)
// }

// Notify sends or resends notifications for payment link.
func Notify(ctx context.Context, id string, medium string) error {
	return getDefaultClient().Notify(ctx, id, medium)
}

// Cancel cancels payment link.
func Cancel(ctx context.Context, id string) (*razorpay.PaymentLink, error) {
	return getDefaultClient().Cancel(ctx, id)
}

// NewClient returns new client.
func NewClient(apiKey string, apiSecret string, apiBackend razorpay.Backend) *Client {
	return &Client{razorpay.NewClient(apiKey, apiSecret, apiBackend)}
}

func getDefaultClient() *Client {
	return &Client{razorpay.GetDefaultClient()}
}
