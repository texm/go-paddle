package paddle

import (
	"context"
	"net/url"
	"strings"
	"time"
)

type TransactionsService service

type TransactionStatus string

const (
	TransactionStatusDraft     = TransactionStatus("draft")
	TransactionStatusReady     = TransactionStatus("ready")
	TransactionStatusBilled    = TransactionStatus("billed")
	TransactionStatusPaid      = TransactionStatus("paid")
	TransactionStatusCompleted = TransactionStatus("completed")
	TransactionStatusCanceled  = TransactionStatus("canceled")
	TransactionStatusPastDue   = TransactionStatus("past_due")
)

type Proration struct {
	Rate          string     `json:"rate"`
	BillingPeriod TimePeriod `json:"billing_period"`
}

type TransactionItem struct {
	Price     Price      `json:"price"`
	Quantity  int        `json:"quantity"`
	Proration *Proration `json:"proration"`
}

type TransactionBillingDetails struct {
	EnableCheckout        bool         `json:"enable_checkout"`
	PaymentTerms          TimeInterval `json:"payment_terms"`
	PurchaseOrderNumber   string       `json:"purchase_order_number"`
	AdditionalInformation *string      `json:"additional_information"`
}

type TransactionTotals struct {
	Subtotal     string  `json:"subtotal"`
	Discount     string  `json:"discount"`
	Tax          string  `json:"tax"`
	Total        string  `json:"total"`
	Credit       string  `json:"credit"`
	Balance      string  `json:"balance"`
	GrandTotal   string  `json:"grand_total"`
	Fee          *string `json:"fee"`
	Earnings     *string `json:"earnings"`
	CurrencyCode string  `json:"currency_code"`
}

type TransactionAdjustedTotals struct {
	Subtotal     string  `json:"subtotal"`
	Tax          string  `json:"tax"`
	Total        string  `json:"total"`
	GrandTotal   string  `json:"grand_total"`
	Fee          *string `json:"fee"`
	Earnings     *string `json:"earnings"`
	CurrencyCode string  `json:"currency_code"`
}

type TransactionPayoutTotals struct {
	Subtotal     string `json:"subtotal"`
	Discount     string `json:"discount"`
	Tax          string `json:"tax"`
	Total        string `json:"total"`
	Credit       string `json:"credit"`
	Balance      string `json:"balance"`
	GrandTotal   string `json:"grand_total"`
	Fee          string `json:"fee"`
	Earnings     string `json:"earnings"`
	CurrencyCode string `json:"currency_code"`
}

type TransactionAdjustedPayoutTotalsChargebackFee struct {
	Amount   string         `json:"amount"`
	Original *CurrencyPrice `json:"original"`
}

type TransactionAdjustedPayoutTotals struct {
	Subtotal      string                                       `json:"subtotal"`
	Tax           string                                       `json:"tax"`
	Total         string                                       `json:"total"`
	Fee           string                                       `json:"fee"`
	ChargebackFee TransactionAdjustedPayoutTotalsChargebackFee `json:"chargeback_fee"`
	Earnings      string                                       `json:"earnings"`
	CurrencyCode  string                                       `json:"currency_code"`
}

type TransactionLineItemTotal struct {
	Subtotal string `json:"subtotal"`
	Tax      string `json:"tax"`
	Discount string `json:"discount"`
	Total    string `json:"total"`
}

type TransactionLineItem struct {
	Id         string                   `json:"id"`
	PriceId    string                   `json:"price_id"`
	Quantity   int                      `json:"quantity"`
	Proration  *Proration               `json:"proration"`
	TaxRate    string                   `json:"tax_rate"`
	Totals     TransactionLineItemTotal `json:"totals"`
	UnitTotals TransactionLineItemTotal `json:"unit_totals"`
	Product    Product                  `json:"product"`
}

type TransactionDetails struct {
	TaxRatesUsed         []TaxRate                        `json:"tax_rates_used"`
	Totals               TransactionTotals                `json:"totals"`
	AdjustedTotals       *TransactionAdjustedTotals       `json:"adjusted_totals"`
	PayoutTotals         *TransactionPayoutTotals         `json:"payout_totals"`
	AdjustedPayoutTotals *TransactionAdjustedPayoutTotals `json:"adjusted_payout_totals"`
	LineItems            []TransactionLineItem            `json:"line_items"`
}

type Checkout struct {
	Url *string `json:"url"`
}

type Transaction struct {
	Id     string            `json:"id"`
	Status TransactionStatus `json:"status"`

	CustomerId     *string `json:"customer_id"`
	AddressId      *string `json:"address_id"`
	BusinessId     *string `json:"business_id"`
	SubscriptionId *string `json:"subscription_id"`
	DiscountId     *string `json:"discount_id"`
	InvoiceId      *string `json:"invoice_id"`
	InvoiceNumber  *string `json:"invoice_number"`

	CurrencyCode string `json:"currency_code"`
	Origin       string `json:"origin"`

	CollectionMode PaymentCollectionMode      `json:"collection_mode"`
	BillingDetails *TransactionBillingDetails `json:"billing_details"`
	BillingPeriod  *TimePeriod                `json:"billing_period"`

	Items    []TransactionItem           `json:"items"`
	Details  TransactionDetails          `json:"details"`
	Payments []TransactionPaymentAttempt `json:"payments"`
	Checkout *Checkout                   `json:"checkout"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	BilledAt  time.Time `json:"billed_at"`

	CustomData *map[string]any `json:"custom_data"`

	Adjustments       []any `json:"adjustments"`
	AdjustmentsTotals *any  `json:"adjustments_totals"`

	Address  *any      `json:"address"`
	Business *any      `json:"business"`
	Customer *Customer `json:"customer"`
	Discount *any      `json:"discount"`
}

type TransactionIncludeParam struct {
	Address          bool
	Adjustment       bool
	AdjustmentTotals bool
	Business         bool
	Customer         bool
	Discount         bool
}

func (ti *TransactionIncludeParam) String() string {
	var includes []string
	if ti.Address {
		includes = append(includes, "address")
	}
	if ti.Adjustment {
		includes = append(includes, "adjustment")
	}
	if ti.AdjustmentTotals {
		includes = append(includes, "adjustment_totals")
	}
	if ti.Business {
		includes = append(includes, "business")
	}
	if ti.Customer {
		includes = append(includes, "customer")
	}
	if ti.Discount {
		includes = append(includes, "discount")
	}
	return strings.Join(includes, ",")
}

type ListTransactionsParams struct {
	Ids             []string
	Include         *TransactionIncludeParam
	CollectionMode  PaymentCollectionMode
	CustomerIds     []string
	SubscriptionIds []string
	InvoiceNumber   []string
	Status          []TransactionStatus
	CreatedAt       string
	BilledAt        string
}

func (ltp *ListTransactionsParams) Encode() string {
	q := url.Values{}
	if len(ltp.Ids) > 0 {
		q.Set("id", strings.Join(ltp.Ids, ","))
	}
	if ltp.Include != nil {
		q.Set("include", ltp.Include.String())
	}
	if len(ltp.CollectionMode) > 0 {
		q.Set("collection_mode", string(ltp.CollectionMode))
	}
	if len(ltp.CustomerIds) > 0 {
		q.Set("customer_id", strings.Join(ltp.CustomerIds, ","))
	}
	if len(ltp.SubscriptionIds) > 0 {
		q.Set("subscription_id", strings.Join(ltp.SubscriptionIds, ","))
	}
	if len(ltp.Status) > 0 {
		q.Set("status", strings.Join(toStringSlice(ltp.Status), ","))
	}
	if len(ltp.InvoiceNumber) > 0 {
		q.Set("invoice_number", strings.Join(toStringSlice(ltp.InvoiceNumber), ","))
	}
	if len(ltp.CreatedAt) > 0 {
		q.Set("created_at", ltp.CreatedAt)
	}
	if len(ltp.BilledAt) > 0 {
		q.Set("billed_at", ltp.BilledAt)
	}
	return q.Encode()
}

func (s *TransactionsService) List(ctx context.Context, params *ListTransactionsParams) ([]*Transaction, error) {
	endpoint := "transactions"
	if params != nil {
		endpoint += "?" + params.Encode()
	}
	return listItems[Transaction](ctx, s.client, endpoint)
}

func (s *TransactionsService) Get(ctx context.Context, id string, include *TransactionIncludeParam) (*Transaction, error) {
	endpoint := "transactions/" + id
	if include != nil {
		endpoint += "?" + include.String()
	}
	return getItem[Transaction](ctx, s.client, "transactions/"+id)
}
