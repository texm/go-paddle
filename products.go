package paddle

import (
	"context"
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
	Prices      *[]Price        `json:"prices"`
}

func (p *ProductsService) List(ctx context.Context) ([]*Product, error) {
	return listItems[Product](ctx, p.client, "products")
}

func (p *ProductsService) Get(ctx context.Context, id string, includePrices bool) (*Product, error) {
	return getItem[Product](ctx, p.client, "products/"+id)
}
