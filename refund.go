package razorpay

// Refund is a Razorpay entity representation.
type Refund struct {
	Response
	Entity
	Amount         int64             `json:"amount"`
	Currency       string            `json:"currency"`
	PaymentId      string            `json:"payment_id"`
	Receipt        string            `json:"receipt"`
	AcquirerData   map[string]string `json:"acquirer_data"`
	Status         string            `json:"status"`
	SpeedProcessed string            `json:"speed_processed"`
	SpeedRequested string            `json:"speed_requested"`
	Notes          Notes             `json:"notes"`
}

// RefundList is collection of refunds.
type RefundList struct {
	Response
	EntityList
	Refunds []*Refund `json:"items"`
}

// RefundCreateParams is list of params that can be used when creating refund.
type RefundCreateParams struct {
	Params
	Amount  *int64  `json:"amount,omitempty"`
	Receipt *string `json:"receipt,omitempty"`
	Speed   *string `json:"speed,omitempty"`
	Notes   Notes   `json:"notes,omitempty"`
}

// RefundListParams is list of params that can be used when listing refunds.
type RefundListParams struct {
	ListParams
}

// RefundUpdateParams is list of params that can be used when updating refund.
type RefundUpdateParams struct {
	Params
	Notes Notes `json:"notes,omitempty"`
}
