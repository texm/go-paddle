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
	Sandbox bool

	APIKey           string
	WebhookSecretKey string
}

type ProductService service
type SubscriptionService service

type Client struct {
	client  *http.Client
	cfg     *Config
	baseURL *url.URL

	Subscriptions *SubscriptionService
	Products      *ProductService
}

type service struct {
	client *Client
}

func (conf *Config) NewClient(client *http.Client) *Client {
	c := &Client{
		client: client,
		cfg:    conf,
	}

	if conf.Sandbox {
		c.baseURL, _ = url.Parse(sandboxApiBaseURL)
	} else {
		c.baseURL, _ = url.Parse(apiBaseURL)
	}

	s := &service{client: c}

	c.Subscriptions = (*SubscriptionService)(s)
	c.Products = (*ProductService)(s)

	return c
}

func (c *Client) TestAuthentication(ctx context.Context) error {
	req, reqErr := c.NewRequest(http.MethodGet, "event-types", nil)
	if reqErr != nil {
		return reqErr
	}
	_, resErr := c.Do(ctx, req, nil)
	if resErr != nil {
		return resErr
	}
	return nil
}

func (c *Client) NewRequest(method, path string, body any) (*http.Request, error) {
	endpoint, parseErr := c.baseURL.Parse(path)
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

func (c *Client) Do(ctx context.Context, req *http.Request, v any) (*http.Response, error) {
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

	if err := checkError(resp, data); err != nil {
		return resp, err
	}

	if v != nil {
		if err := json.Unmarshal(data, v); err != nil {
			return resp, fmt.Errorf("err=%v, data=%v", err, string(data))
		}
	}

	return resp, nil
}

type ErrorType string

const (
	ErrorTypeRequest = ErrorType("request_error")
	ErrorTypeApi     = ErrorType("api_error")
)

type Error struct {
	Type             ErrorType `json:"type"`
	Code             string    `json:"code"`
	Detail           string    `json:"detail"`
	DocumentationUrl string    `json:"documentation_url"`
}

type ErrorResponse struct {
	response *http.Response

	ErrorDetails *Error `json:"error"`
	Meta         struct {
		RequestId string `json:"request_id"`
	} `json:"meta"`
}

func (r *ErrorResponse) Error() string {
	if r.ErrorDetails == nil {
		res := r.response
		return fmt.Sprintf("[%v %v] HTTP %d: unknown api error", res.Request.Method, res.Request.URL.Path, res.StatusCode)
	}
	return fmt.Sprintf("[%v %v] HTTP %d '%s': %s",
		r.response.Request.Method, r.response.Request.URL.Path,
		r.response.StatusCode, r.ErrorDetails.Code, r.ErrorDetails.Detail)
}

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
