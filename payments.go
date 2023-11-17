package paddle

import "time"

type CardDetails struct {
	Type           string `json:"type"`
	Last4          string `json:"last4"`
	ExpiryMonth    int    `json:"expiry_month"`
	ExpiryYear     int    `json:"expiry_year"`
	CardholderName string `json:"cardholder_name"`
}

type PaymentMethod struct {
	Type string       `json:"type"`
	Card *CardDetails `json:"card"`
}

type TransactionPaymentAttempt struct {
	PaymentAttemptId      string        `json:"payment_attempt_id"`
	StoredPaymentMethodId string        `json:"stored_payment_method_id"`
	Amount                string        `json:"amount"`
	Status                string        `json:"status"`
	ErrorCode             *string       `json:"error_code"`
	MethodDetails         PaymentMethod `json:"method_details"`
	CreatedAt             time.Time     `json:"created_at"`
	CapturedAt            time.Time     `json:"captured_at"`
}

type TaxRateTotals struct {
	Subtotal string `json:"subtotal"`
	Discount string `json:"discount"`
	Tax      string `json:"tax"`
	Total    string `json:"total"`
}

type TaxRate struct {
	TaxRate string        `json:"tax_rate"`
	Totals  TaxRateTotals `json:"totals"`
}
