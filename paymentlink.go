package razorpay

import "encoding/json"

// PaymentLink is a Razorpay entity representation.
type PaymentLink struct {
	Response
	Entity
	Amount                int64                       `json:"amount"`
	Currency              string                      `json:"currency"`
	AcceptPartial         bool                        `json:"accept_partial"`
	FirstMinPartialAmount int64                       `json:"first_min_partial_amount"`
	AmountPaid            int64                       `json:"amount_paid"`
	Description           string                      `json:"description"`
	Customer              Customer                    `json:"customer"`
	CallbackMethod        string                      `json:"callback_method"`
	CallbackUrl           string                      `json:"callback_url"`
	CancelledAt           int64                       `json:"cancelled_at"`
	ExpireBy              int64                       `json:"expire_by"`
	ExpiredAt             int64                       `json:"expired_at"`
	Notify                PaymentLinkNotify           `json:"notify"`
	ReferenceID           string                      `json:"reference_id"`
	ReminderEnable        bool                        `json:"reminder_enable"`
	Reminders             PaymentLinkRemindersWrapper `json:"reminders"`
	ShortUrl              string                      `json:"short_url"`
	Status                string                      `json:"status"`
	Notes                 Notes                       `json:"notes"`

	// TODO: To add `Payments` field in the struct. Need to handle `[]` i.e. empty list as value.
}

type PaymentLinkNotify struct {
	Email bool `json:"email"`
	SMS   bool `json:"sms"`
}

type PaymentLinkReminders struct {
	Status string `json:"status"`
}

// PaymentLinkRemindersWrapper wraps PaymentLinkReminders and overrides unmarshalling.
type PaymentLinkRemindersWrapper struct {
	*PaymentLinkReminders
}

// UnmarshalJSON unmarshals raw `reminders` into the type.
func (w *PaymentLinkRemindersWrapper) UnmarshalJSON(data []byte) error {
	// In json response, `reminders` will appear slice i.e. `[]` when it is empty,
	// and so only attempts unmarshal when it appears object.
	if len(data) > 0 && data[0] == '{' {
		return json.Unmarshal(data, w.PaymentLinkReminders)
	}
	return nil
}

// PaymentLinkList is collection of payment links.
type PaymentLinkList struct {
	Response
	EntityList
	PaymentLinks []*PaymentLink `json:"items"`
}

// PaymentLinkParams is list of params that can be used when creating or updating payment link.
type PaymentLinkParams struct {
	Params
	Amount                *int64                   `json:"amount,omitempty"`
	Currency              *string                  `json:"currency,omitempty"`
	AcceptPartial         *bool                    `json:"accept_partial,omitempty"`
	FirstMinPartialAmount *int64                   `json:"first_min_partial_amount,omitempty"`
	Description           *string                  `json:"description,omitempty"`
	Customer              *CustomerParams          `json:"customer,omitempty"`
	CallbackMethod        *string                  `json:"callback_method,omitempty"`
	CallbackUrl           *string                  `json:"callback_url,omitempty"`
	ExpireBy              *int64                   `json:"expire_by,omitempty"`
	Notify                *PaymentLinkNotifyParams `json:"notify,omitempty"`
	ReferenceID           *string                  `json:"reference_id,omitempty"`
	ReminderEnable        *bool                    `json:"reminder_enable,omitempty"`
	Notes                 Notes                    `json:"notes,omitempty"`
}

// PaymentLinkNotifyParams is part of PaymentLinkParams.
type PaymentLinkNotifyParams struct {
	Email *bool `json:"email,omitempty"`
	SMS   *bool `json:"sms,omitempty"`
}

// PaymentLinkListParams is list params that can be used when listing payment links.
type PaymentLinkListParams struct {
	ListParams
	PaymentID   *string `url:"payment_id,omitempty"`
	ReferenceID *string `url:"reference_id,omitempty"`
}
