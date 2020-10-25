package razorpay

// Order is a Razorpay entity representation.
type Order struct {
	Response
	Entity
	Amount     int64  `json:"amount"`
	AmountPaid int64  `json:"amount_paid"`
	AmountDue  int64  `json:"amount_due"`
	Currency   string `json:"currency"`
	Receipt    string `json:"receipt"`
	Status     string `json:"status"`
	Attempts   int64  `json:"attempts"`
	Notes      Notes  `json:"notes"`
}

// OrderList is collection of orders.
type OrderList struct {
	Response
	EntityList
	Orders []*Order `json:"items"`
}

// OrderParams is list of params that can be used when creating or updating order.
type OrderParams struct {
	Params
	Amount   *int64  `json:"amount,omitempty"`
	Currency *string `json:"currency,omitempty"`
	Receipt  *string `json:"receipt,omitempty"`
	Notes    Notes   `json:"notes,omitempty"`
}

// OrderListParams is list params that can be used when listing orders.
type OrderListParams struct {
	ListParams
	Authorized *string `url:"authorized,omitempty"`
	Receipt    *string `url:"receipt,omitempty"`
}
