package paddle

import (
	"context"
	"time"
)

type SubscriptionsService service

type SubscriptionStatus string
type SubscriptionCollectionMode string
type SubscriptionScheduledChangeAction string

const (
	SubscriptionStatusActive   = SubscriptionStatus("active")
	SubscriptionStatusCanceled = SubscriptionStatus("canceled")
	SubscriptionStatusPastDue  = SubscriptionStatus("past_due")
	SubscriptionStatusPaused   = SubscriptionStatus("paused")
	SubscriptionStatusTrialing = SubscriptionStatus("trialing")

	SubscriptionCollectionModeAutomatic = SubscriptionCollectionMode("automatic")
	SubscriptionCollectionModeManual    = SubscriptionCollectionMode("manual")

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
	CollectionMode       string                       `json:"collection_mode"`
	BillingDetails       *SubscriptionBillingDetails  `json:"billing_details"`
	CurrentBillingPeriod *TimePeriod                  `json:"current_billing_period"`
	BillingCycle         TimeInterval                 `json:"billing_cycle"`
	ScheduledChange      *SubscriptionScheduledChange `json:"scheduled_change"`
	Items                []SubscriptionItem           `json:"items"`
	CustomData           map[string]any               `json:"custom_data"`
	ManagementUrls       SubscriptionManagementUrls   `json:"management_urls"`
}

func (s *SubscriptionsService) List(ctx context.Context) ([]*Subscription, error) {
	return listItems[Subscription](ctx, s.client, "subscriptions")
}

func (s *SubscriptionsService) Get(ctx context.Context, id string) (*Subscription, error) {
	return getItem[Subscription](ctx, s.client, "subscriptions/"+id)
}
