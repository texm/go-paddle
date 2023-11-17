package paddle

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

func makeApiRequest[T any](ctx context.Context, c *Client, method string, endpoint string, body any) (*T, *ApiResponse, error) {
	req, reqErr := c.NewRequest(method, endpoint, body)
	if reqErr != nil {
		return nil, nil, reqErr
	}
	res, resErr := c.Do(ctx, req)
	if resErr != nil {
		return nil, res, resErr
	}
	var item *T
	if jsonErr := json.Unmarshal(res.Data, &item); jsonErr != nil {
		return nil, res, jsonErr
	}
	return item, res, nil
}

func postItem[T any](ctx context.Context, c *Client, endpoint string, body any) (*T, error) {
	item, _, err := makeApiRequest[T](ctx, c, http.MethodPost, endpoint, body)
	return item, err
}

func patchItem[T any](ctx context.Context, c *Client, endpoint string, body any) (*T, error) {
	item, _, err := makeApiRequest[T](ctx, c, http.MethodPatch, endpoint, body)
	return item, err
}

func getItem[T any](ctx context.Context, c *Client, endpoint string) (*T, error) {
	item, _, err := makeApiRequest[T](ctx, c, http.MethodGet, endpoint, nil)
	return item, err
}

func listItems[T any](ctx context.Context, c *Client, basePath string) ([]*T, error) {
	curPath := basePath
	hasMore := true
	var items []*T
	for hasMore {
		resItems, res, resErr := makeApiRequest[[]T](ctx, c, http.MethodGet, curPath, nil)
		if resErr != nil {
			return nil, resErr
		}
		if resItems == nil {
			continue
		}
		for _, item := range *resItems {
			it := item
			items = append(items, &it)
		}

		hasMore = res.Meta.Pagination.HasMore
		if hasMore {
			curPath = strings.TrimPrefix(res.Meta.Pagination.Next, c.baseURL)
		}
	}
	return items, nil
}
