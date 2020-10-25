package razorpay

// Customer is a Razorpay entity representation.
type Customer struct {
	Response
	Entity
	Name    string `json:"name"`
	Email   string `json:"email"`
	Contact string `json:"contact"`
	Gstin   string `json:"gstin"`
	Notes   Notes  `json:"notes"`
}

// CustomerList is collection of customers.
type CustomerList struct {
	Response
	EntityList
	Customers []*Customer `json:"items"`
}

// CustomerParams is list of params that can be used when creating or updating customer.
type CustomerParams struct {
	Params
	Name    *string `json:"name,omitempty"`
	Contact *string `json:"contact,omitempty"`
	Email   *string `json:"email,omitempty"`
	Gstin   *string `json:"gstin,omitempty"`
	Notes   Notes   `json:"notes,omitempty"`
}

// CustomerListParams is list params that can be used when listing customers.
type CustomerListParams struct {
	ListParams
}
