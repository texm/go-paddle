package paddle

import (
	"context"
	"time"
)

type CustomersService service

type Customer struct {
	Id               string          `json:"id"`
	Status           Status          `json:"status"`
	CustomData       *map[string]any `json:"custom_data"`
	Name             *string         `json:"name"`
	Email            string          `json:"email"`
	MarketingConsent bool            `json:"marketing_consent"`
	Locale           string          `json:"locale"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

func (c *CustomersService) List(ctx context.Context) ([]*Customer, error) {
	return listItems[Customer](ctx, c.client, "customers")
}

func (c *CustomersService) Get(ctx context.Context, id string) (*Customer, error) {
	return getItem[Customer](ctx, c.client, "customers/"+id)
}

type CreateCustomerParams struct {
	Email      string          `json:"email"`
	Name       *string         `json:"name,omitempty"`
	CustomData *map[string]any `json:"custom_data,omitempty"`
	Locale     *string         `json:"locale,omitempty"`
}

func (c *CustomersService) Create(ctx context.Context, params *CreateCustomerParams) (*Customer, error) {
	return postItem[Customer](ctx, c.client, "customers", params)
}

type UpdateCustomerParams struct {
	Email      *string         `json:"email,omitempty"`
	Name       *string         `json:"name,omitempty"`
	CustomData *map[string]any `json:"custom_data,omitempty"`
	Locale     *string         `json:"locale,omitempty"`
}

func (c *CustomersService) Update(ctx context.Context, params *UpdateCustomerParams) (*Customer, error) {
	return patchItem[Customer](ctx, c.client, "customers", params)
}
