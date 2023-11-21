package paddle

import (
	"context"
	"net/url"
	"strings"
	"time"
)

type SubscriptionsService service

type SubscriptionStatus string
type SubscriptionScheduledChangeAction string

const (
	SubscriptionStatusActive   = SubscriptionStatus("active")
	SubscriptionStatusCanceled = SubscriptionStatus("canceled")
	SubscriptionStatusPastDue  = SubscriptionStatus("past_due")
	SubscriptionStatusPaused   = SubscriptionStatus("paused")
	SubscriptionStatusTrialing = SubscriptionStatus("trialing")

	SubscriptionScheduledChangeActionCancel = SubscriptionScheduledChangeAction("cancel")
	SubscriptionScheduledChangeActionPause  = SubscriptionScheduledChangeAction("pause")
	SubscriptionScheduledChangeActionResume = SubscriptionScheduledChangeAction("resume")
)

type SubscriptionDiscount struct {
	Id       string    `json:"id"`
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
}

type SubscriptionBillingDetails struct {
	EnableCheckout        bool         `json:"enable_checkout"`
	PurchaseOrderNumber   string       `json:"purchase_order_number"`
	AdditionalInformation *string      `json:"additional_information"`
	PaymentTerms          TimeInterval `json:"payment_terms"`
}

type SubscriptionScheduledChange struct {
	Action      SubscriptionScheduledChangeAction `json:"action"`
	EffectiveAt time.Time                         `json:"effective_at"`
	ResumeAt    time.Time                         `json:"resume_at"`
}

type SubscriptionManagementUrls struct {
	UpdatePaymentMethod string `json:"update_payment_method"`
	Cancel              string `json:"cancel"`
}

type SubscriptionItem struct {
	Status             string     `json:"status"`
	Quantity           int        `json:"quantity"`
	Recurring          bool       `json:"recurring"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	PreviouslyBilledAt time.Time  `json:"previously_billed_at"`
	NextBilledAt       time.Time  `json:"next_billed_at"`
	TrialDates         TimePeriod `json:"trial_dates"`
	Price              Price      `json:"price"`
}

type Subscription struct {
	Id                   string                       `json:"id"`
	Status               SubscriptionStatus           `json:"status"`
	CustomerId           string                       `json:"customer_id"`
	AddressId            string                       `json:"address_id"`
	BusinessId           *string                      `json:"business_id"`
	CurrencyCode         string                       `json:"currency_code"`
	CreatedAt            time.Time                    `json:"created_at"`
	UpdatedAt            time.Time                    `json:"updated_at"`
	StartedAt            time.Time                    `json:"started_at"`
	FirstBilledAt        time.Time                    `json:"first_billed_at"`
	NextBilledAt         time.Time                    `json:"next_billed_at"`
	PausedAt             time.Time                    `json:"paused_at"`
	CanceledAt           time.Time                    `json:"canceled_at"`
	Discount             *SubscriptionDiscount        `json:"discount"`
	CollectionMode       PaymentCollectionMode        `json:"collection_mode"`
	BillingDetails       *SubscriptionBillingDetails  `json:"billing_details"`
	CurrentBillingPeriod *TimePeriod                  `json:"current_billing_period"`
	BillingCycle         TimeInterval                 `json:"billing_cycle"`
	ScheduledChange      *SubscriptionScheduledChange `json:"scheduled_change"`
	Items                []SubscriptionItem           `json:"items"`
	CustomData           map[string]any               `json:"custom_data"`
	ManagementUrls       SubscriptionManagementUrls   `json:"management_urls"`
}

type ListSubscriptionsParams struct {
	Ids            []string
	CollectionMode string
	Status         []Status
	Search         string
}

func (s *SubscriptionsService) List(ctx context.Context, params *ListSubscriptionsParams) ([]*Subscription, error) {
	endpoint := "subscriptions"
	if params != nil {
		q := url.Values{}
		if len(params.Ids) > 0 {
			q.Set("id", strings.Join(params.Ids, ","))
		}
		if len(params.Status) > 0 {
			q.Set("status", strings.Join(params.Ids, ","))
		}
		if len(params.Search) > 0 {
			q.Set("search", params.Search)
		}
		endpoint += "?" + q.Encode()
	}
	return listItems[Subscription](ctx, s.client, endpoint)
}

func (s *SubscriptionsService) Get(ctx context.Context, id string) (*Subscription, error) {
	return getItem[Subscription](ctx, s.client, "subscriptions/"+id)
}

type SubscriptionEffectFromOption string

const (
	SubscriptionEffectFromOptionImmediately       = SubscriptionEffectFromOption("immediately")
	SubscriptionEffectFromOptionNextBillingPeriod = SubscriptionEffectFromOption("next_billing_period")
)

type CancelSubscriptionParams struct {
	EffectiveFrom SubscriptionEffectFromOption `json:"effective_from"`
}

func (s *SubscriptionsService) Cancel(ctx context.Context, id string, params *CancelSubscriptionParams) (*Subscription, error) {
	return postItem[Subscription](ctx, s.client, "subscriptions/"+id+"/cancel", params)
}

func (s *SubscriptionsService) RemoveScheduledCancellation(ctx context.Context, id string) (*Subscription, error) {
	body := struct {
		ScheduledChange *any `json:"scheduled_change"`
	}{}
	return patchItem[Subscription](ctx, s.client, "subscriptions/"+id, body)
}

func (s *SubscriptionsService) GetUpdatePaymentMethodTransaction(ctx context.Context, id string) (*Transaction, error) {
	return getItem[Transaction](ctx, s.client, "subscriptions/"+id+"/update-payment-method-transaction")
}
