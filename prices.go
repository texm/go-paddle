package paddle

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

type PricesService service

type CurrencyPrice struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currency_code"`
}

type CurrencyPriceOverride struct {
	CountryCodes []string      `json:"country_codes"`
	UnitPrice    CurrencyPrice `json:"unit_price"`
}

type Price struct {
	Id                 string                  `json:"id"`
	ProductId          string                  `json:"product_id"`
	Description        string                  `json:"description"`
	Name               *string                 `json:"name"`
	BillingCycle       *TimeInterval           `json:"billing_cycle"`
	TrialPeriod        *TimeInterval           `json:"trial_period"`
	TaxMode            string                  `json:"tax_mode"`
	UnitPrice          CurrencyPrice           `json:"unit_price"`
	UnitPriceOverrides []CurrencyPriceOverride `json:"unit_price_overrides"`
	CustomData         *map[string]any         `json:"custom_data"`
	Status             string                  `json:"status"`
	Quantity           MinMax                  `json:"quantity"`
	Product            *Product                `json:"product"`
}

type ListPricesParams struct {
	IncludeProduct        bool
	Ids                   []string
	CustomerIds           []string
	AddressIds            []string
	CollectionMode        string
	ScheduledChangeAction []SubscriptionScheduledChangeAction
	Status                []SubscriptionStatus
}

func toStringSlice[T any](items []T) []string {
	strs := make([]string, len(items))
	for i, item := range items {
		strs[i] = fmt.Sprintf("%v", item)
	}
	return strs
}

func (s *PricesService) List(ctx context.Context, params *ListPricesParams) ([]*Price, error) {
	endpoint := "prices"
	if params != nil {
		q := url.Values{}
		if params.IncludeProduct {
			q.Set("include", "product")
		}
		if len(params.Ids) > 0 {
			q.Set("id", strings.Join(params.Ids, ","))
		}
		if len(params.CustomerIds) > 0 {
			q.Set("customer_id", strings.Join(params.CustomerIds, ","))
		}
		if len(params.AddressIds) > 0 {
			q.Set("address_id", strings.Join(params.AddressIds, ","))
		}
		if len(params.CollectionMode) > 0 {
			q.Set("collection_mode", params.CollectionMode)
		}
		if len(params.ScheduledChangeAction) > 0 {
			q.Set("scheduled_change_action", strings.Join(toStringSlice(params.ScheduledChangeAction), ","))
		}
		if len(params.Status) > 0 {
			q.Set("status", strings.Join(toStringSlice(params.Status), ","))
		}
		endpoint += "?" + q.Encode()
	}
	return listItems[Price](ctx, s.client, endpoint)
}

func (s *PricesService) Get(ctx context.Context, id string, includeProduct bool) (*Price, error) {
	endpoint := "prices/" + id
	if includeProduct {
		endpoint += "?include=product"
	}
	return getItem[Price](ctx, s.client, endpoint)
}
