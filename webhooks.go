package paddle

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	ErrInvalidHeader    = errors.New("invalid header")
	ErrInvalidSignature = errors.New("invalid webhook signature")

	MaxWebhookBodyBytes = int64(65536)
)

func (c *Client) VerifyWebhookRequest(req *http.Request) error {
	sig, providedErr := getWebhookSignature(req.Header.Get("Paddle-Signature"))
	if providedErr != nil {
		return providedErr
	}

	body, bodyErr := getBody(req)
	if bodyErr != nil {
		return fmt.Errorf("failed to read body: %w", bodyErr)
	}

	return sig.validate(c.webhookKey, body)
}

type signature struct {
	timestamp         string
	providedSignature string
}

func (w *signature) validate(key []byte, body []byte) error {
	hash := hmac.New(sha256.New, key)
	prefix := []byte(w.timestamp + ":")
	if _, wErr := hash.Write(append(prefix, body...)); wErr != nil {
		return fmt.Errorf("failed to create hash: %w", wErr)
	}
	sum := hash.Sum(nil)
	if hex.EncodeToString(sum) != w.providedSignature {
		return ErrInvalidSignature
	}
	return nil
}

func getWebhookSignature(raw string) (*signature, error) {
	elements := strings.Split(raw, ";")
	if len(elements) != 2 {
		return nil, ErrInvalidHeader
	}
	ts := strings.Split(elements[0], "=")
	h1 := strings.Split(elements[1], "=")
	if len(ts) != 2 || len(h1) != 2 {
		return nil, ErrInvalidHeader
	}
	return &signature{
		timestamp:         ts[1],
		providedSignature: h1[1],
	}, nil
}

func getBody(req *http.Request) ([]byte, error) {
	reqBody, bodyErr := req.GetBody()
	if bodyErr != nil {
		return nil, fmt.Errorf("failed to create copy: %w", bodyErr)
	}
	body, readErr := io.ReadAll(io.LimitReader(reqBody, MaxWebhookBodyBytes))
	return body, errors.Join(readErr, reqBody.Close())
}
