package paddle

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	ErrInvalidHeader    = errors.New("invalid header")
	ErrInvalidSignature = errors.New("invalid webhook signature")
)

const (
	MaxWebhookBodyBytes = int64(65536)
)

type WebhookEvent struct {
	Id             string          `json:"event_id"`
	Type           string          `json:"event_type"`
	OccurredAt     time.Time       `json:"occurred_at"`
	NotificationId string          `json:"notification_id"`
	Data           json.RawMessage `json:"data"`
}

func (c *Client) ParseWebhook(req *http.Request) (*WebhookEvent, error) {
	sigHeader := req.Header.Get("Paddle-Signature")
	sig, providedErr := getWebhookSignature(sigHeader)
	if providedErr != nil {
		return nil, providedErr
	}

	body, bodyErr := getBody(req)
	if bodyErr != nil {
		return nil, fmt.Errorf("failed to read body: %w", bodyErr)
	}

	if validationErr := sig.validate(c.webhookKey, body); validationErr != nil {
		return nil, validationErr
	}

	var event *WebhookEvent
	if jsonErr := json.Unmarshal(body, event); jsonErr != nil {
		return nil, jsonErr
	}

	return event, nil
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
