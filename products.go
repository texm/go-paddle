package paddle

import (
	"context"
	"net/url"
	"strings"
	"time"
)

type ProductsService service

type Product struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	TaxCategory string          `json:"tax_category"`
	ImageUrl    *string         `json:"image_url"`
	CustomData  *map[string]any `json:"custom_data"`
	Status      Status          `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	Prices      *[]*Price       `json:"prices"`
}

type ListProductsParams struct {
	IncludePrices bool
	Ids           []string
	Status        []Status
	TaxCategory   []string
}

func (p *ProductsService) List(ctx context.Context, params *ListProductsParams) ([]*Product, error) {
	endpoint := "products"
	if params != nil {
		q := url.Values{}
		if len(params.Ids) > 0 {
			q.Set("id", strings.Join(params.Ids, ","))
		}
		if params.IncludePrices {
			q.Set("include", "prices")
		}
		if len(params.Status) > 0 {
			q.Set("status", strings.Join(params.Ids, ","))
		}
		if len(params.TaxCategory) > 0 {
			q.Set("tax_category", strings.Join(params.TaxCategory, ","))
		}
		endpoint += "?" + q.Encode()
	}
	return listItems[Product](ctx, p.client, endpoint)
}

func (p *ProductsService) Get(ctx context.Context, id string, includePrices bool) (*Product, error) {
	endpoint := "products/" + id
	if includePrices {
		endpoint += "?include=prices"
	}
	return getItem[Product](ctx, p.client, endpoint)
}
