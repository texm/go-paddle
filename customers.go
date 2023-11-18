package paddle

import (
	"context"
	"net/url"
	"strings"
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

type ListCustomersParams struct {
	Ids    []string
	Status []Status
	Search string
}

func (c *CustomersService) List(ctx context.Context, params *ListCustomersParams) ([]*Customer, error) {
	endpoint := "customers"
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
	return listItems[Customer](ctx, c.client, endpoint)
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

func (c *CustomersService) Update(ctx context.Context, id string, params *UpdateCustomerParams) (*Customer, error) {
	return patchItem[Customer](ctx, c.client, "customers/"+id, params)
}
