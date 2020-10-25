package razorpay

// Payment is a Razorpay entity representation.
type Payment struct {
	Response
	Entity
	Amount           int64  `json:"amount"`
	Currency         string `json:"currency"`
	Status           string `json:"status"`
	Method           string `json:"method"`
	OrderID          string `json:"order_id"`
	Description      string `json:"description"`
	AmountRefunded   int64  `json:"amount_refunded"`
	RefundStatus     string `json:"refund_status"`
	Email            string `json:"email"`
	Contact          string `json:"contact"`
	Notes            Notes  `json:"notes"`
	Fee              int64  `json:"fee"`
	Tax              int64  `json:"tax"`
	ErrorCode        string `json:"error_code"`
	ErrorDescription string `json:"error_description"`
}

type Card struct {
	Response
	Entity
	Name          string `json:"name"`
	Last4         string `json:"last4"`
	Network       string `json:"network"`
	Type          string `json:"type"`
	Issuer        string `json:"issuer"`
	International bool   `json:"international"`
	Emi           bool   `json:"emi"`
}

// PaymentList is collection of payments.
type PaymentList struct {
	Response
	EntityList
	Payments []*Payment `json:"items"`
}

// PaymentUpdateParams is list of params that can be used when updating existing payment.
type PaymentUpdateParams struct {
	Params
	Notes Notes `json:"notes,omitempty"`
}

// PaymentCaptureParams is list of params that can be used when capturing existing payment.
type PaymentCaptureParams struct {
	Params
	Amount   *int64  `json:"amount,omitempty"`
	Currency *string `json:"currency,omitempty"`
}

// PaymentListParams is list of params that can be used when listing payments.
type PaymentListParams struct {
	ListParams
}
