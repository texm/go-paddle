package paddle

import (
	"context"
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

func (s *PricesService) List(ctx context.Context) ([]*Price, error) {
	return listItems[Price](ctx, s.client, "prices")
}

func (s *PricesService) Get(ctx context.Context, id string, includeProduct bool) (*Price, error) {
	endpoint := "prices/" + id
	if includeProduct {
		endpoint += "?include=product"
	}
	return getItem[Price](ctx, s.client, endpoint)
}
