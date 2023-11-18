package paddle

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	apiBaseURL        = "https://api.paddle.com/"
	sandboxApiBaseURL = "https://sandbox-api.paddle.com/"
)

type Config struct {
	Sandbox    bool
	HttpClient *http.Client

	APIKey           string
	WebhookSecretKey string
}

type Client struct {
	client     *http.Client
	cfg        *Config
	baseURL    string
	apiKey     string
	webhookKey []byte

	Customers     *CustomersService
	Subscriptions *SubscriptionsService
	Products      *ProductsService
	Prices        *PricesService
	Transactions  *TransactionsService
}

type service struct {
	client *Client
}

func NewClient(cfg *Config) *Client {
	if cfg.HttpClient == nil {
		cfg.HttpClient = http.DefaultClient
	}

	c := &Client{
		client:     cfg.HttpClient,
		cfg:        cfg,
		apiKey:     cfg.APIKey,
		webhookKey: []byte(cfg.WebhookSecretKey),
	}

	if cfg.Sandbox {
		c.baseURL = sandboxApiBaseURL
	} else {
		c.baseURL = apiBaseURL
	}

	s := &service{client: c}

	c.Customers = (*CustomersService)(s)
	c.Subscriptions = (*SubscriptionsService)(s)
	c.Products = (*ProductsService)(s)
	c.Prices = (*PricesService)(s)
	c.Transactions = (*TransactionsService)(s)

	return c
}

func (c *Client) TestAuthentication(ctx context.Context) error {
	req, reqErr := c.NewRequest(http.MethodGet, "event-types", nil)
	if reqErr != nil {
		return reqErr
	}
	_, resErr := c.Do(ctx, req)
	if resErr != nil {
		return resErr
	}
	return nil
}

func (c *Client) NewRequest(method string, path string, body any) (*http.Request, error) {
	endpoint, parseErr := url.Parse(c.baseURL + path)
	if parseErr != nil {
		return nil, parseErr
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if encodeErr := enc.Encode(body); encodeErr != nil {
			return nil, encodeErr
		}
	}

	req, reqErr := http.NewRequest(method, endpoint.String(), buf)
	if reqErr != nil {
		return nil, reqErr
	}

	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

type ApiResponseMetaPagination struct {
	PerPage        int    `json:"per_page"`
	Next           string `json:"next"`
	HasMore        bool   `json:"has_more"`
	EstimatedTotal int    `json:"estimated_total"`
}

type ApiResponseMeta struct {
	RequestId  string                    `json:"request_id"`
	Pagination ApiResponseMetaPagination `json:"pagination,omitempty"`
}

type ErrorType string

const (
	ErrorTypeRequest = ErrorType("request_error")
	ErrorTypeApi     = ErrorType("api_error")
)

type ApiError struct {
	res *http.Response

	Type             ErrorType `json:"type"`
	Code             string    `json:"code"`
	Detail           string    `json:"detail"`
	DocumentationUrl string    `json:"documentation_url"`
}

type ApiResponse struct {
	Data  json.RawMessage `json:"data"`
	Error *ApiError       `json:"error,omitempty"`
	Meta  ApiResponseMeta `json:"meta"`
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*ApiResponse, error) {
	resp, respErr := c.client.Do(req)

	// HTTP status codes do not contribute to a response error
	if respErr != nil {
		var err error
		select {
		case <-ctx.Done():
			err = errors.Join(err, ctx.Err())
		default:
		}

		// If the error type is *url.Error
		var urlErr *url.Error
		if errors.As(respErr, &urlErr) {
			if parsedUrl, parseErr := url.Parse(urlErr.URL); parseErr == nil {
				urlErr.URL = parsedUrl.String()
				err = errors.Join(err, urlErr)
			}
		}

		return nil, err
	}
	defer resp.Body.Close()

	data, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	res := &ApiResponse{}
	if jsonErr := json.Unmarshal(data, res); jsonErr != nil {
		return nil, fmt.Errorf("http %d: failed to read response: %w", resp.StatusCode, jsonErr)
	}

	if res.Error != nil {
		res.Error.res = resp
		return nil, res.Error
	}

	return res, nil
}

func (e *ApiError) Error() string {
	res := e.res
	return fmt.Sprintf("[%v %v] HTTP %d '%s': %s",
		res.Request.Method, res.Request.URL.Path, res.StatusCode, e.Code, e.Detail)
}

/*
func checkError(r *http.Response, data []byte) error {
	if r.StatusCode < http.StatusBadRequest {
		return nil
	}
	errRes := &ErrorResponse{response: r}
	if data != nil {
		if jsonErr := json.Unmarshal(data, errRes); jsonErr != nil {
			return jsonErr
		}
	}
	return errRes
}
*/
