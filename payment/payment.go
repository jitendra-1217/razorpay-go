package payment

import (
	"context"
	"net/http"

	razorpay "github.com/jitendra-1217/razorpay-go"
)

// Client is used to access /payments apis.
type Client struct {
	*razorpay.Client
}

// Update updates existing payment.
func (c *Client) Update(ctx context.Context, id string, params *razorpay.PaymentUpdateParams) (*razorpay.Payment, error) {
	payment := &razorpay.Payment{}
	err := c.Call(ctx, http.MethodPatch, "/payments/"+id, params, payment)
	return payment, err
}

// Capture captures existing payment.
func (c *Client) Capture(ctx context.Context, id string, params *razorpay.PaymentCaptureParams) (*razorpay.Payment, error) {
	payment := &razorpay.Payment{}
	err := c.Call(ctx, http.MethodPost, "/payments/"+id+"/capture", params, payment)
	return payment, err
}

// Get returns payment for id.
func (c *Client) Get(ctx context.Context, id string, params *razorpay.GetParams) (*razorpay.Payment, error) {
	if params == nil {
		params = &razorpay.GetParams{}
	}

	payment := &razorpay.Payment{}
	err := c.Call(ctx, http.MethodGet, "/payments/"+id, params, payment)
	return payment, err
}

// List returns list of payments for params.
func (c *Client) List(ctx context.Context, params *razorpay.PaymentListParams) (*razorpay.PaymentList, error) {
	if params == nil {
		params = &razorpay.PaymentListParams{}
	}

	paymentList := &razorpay.PaymentList{}
	err := c.Call(ctx, http.MethodGet, "/payments", params, paymentList)
	return paymentList, err
}

// GetCard returns card details of payment.
func (c *Client) GetCard(ctx context.Context, paymentID string) (*razorpay.Card, error) {
	card := &razorpay.Card{}
	err := c.Call(ctx, http.MethodGet, "/payments/"+paymentID+"/card", nil, card)
	return card, err
}

// CreateRefund creates new refund for the payment.
func (c *Client) CreateRefund(ctx context.Context, paymentID string, params *razorpay.RefundCreateParams) (*razorpay.Refund, error) {
	refund := &razorpay.Refund{}
	err := c.Call(ctx, http.MethodPost, "/payments/"+paymentID+"/refund", params, refund)
	return refund, err
}

// Refunds returns list of refunds for payment.
func (c *Client) Refunds(ctx context.Context, paymentID string) (*razorpay.RefundList, error) {
	refundList := &razorpay.RefundList{}
	err := c.Call(ctx, http.MethodGet, "/payments/"+paymentID+"/refunds", nil, refundList)
	return refundList, err
}

// Update updates existing payment.
func Update(ctx context.Context, id string, params *razorpay.PaymentUpdateParams) (*razorpay.Payment, error) {
	return getDefaultClient().Update(ctx, id, params)
}

// Capture captures existing payment.
func Capture(ctx context.Context, id string, params *razorpay.PaymentCaptureParams) (*razorpay.Payment, error) {
	return getDefaultClient().Capture(ctx, id, params)
}

// Get returns payment for id.
func Get(ctx context.Context, id string, params *razorpay.GetParams) (*razorpay.Payment, error) {
	return getDefaultClient().Get(ctx, id, params)
}

// List returns list of payments for params.
func List(ctx context.Context, params *razorpay.PaymentListParams) (*razorpay.PaymentList, error) {
	return getDefaultClient().List(ctx, params)
}

// GetCard returns card details of payment.
func GetCard(ctx context.Context, paymentID string) (*razorpay.Card, error) {
	return getDefaultClient().GetCard(ctx, paymentID)
}

// CreateRefund creates new refund for the payment.
func CreateRefund(ctx context.Context, paymentID string, params *razorpay.RefundCreateParams) (*razorpay.Refund, error) {
	return getDefaultClient().CreateRefund(ctx, paymentID, params)
}

// Refunds returns list of refunds for payment.
func Refunds(ctx context.Context, paymentID string) (*razorpay.RefundList, error) {
	return getDefaultClient().Refunds(ctx, paymentID)
}

// NewClient returns new client.
func NewClient(apiKey string, apiSecret string, apiBackend razorpay.Backend) *Client {
	return &Client{razorpay.NewClient(apiKey, apiSecret, apiBackend)}
}

func getDefaultClient() *Client {
	return &Client{razorpay.GetDefaultClient()}
}
