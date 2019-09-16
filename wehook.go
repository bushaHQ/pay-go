package pay

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// webhook
// verify

// WebhookService service
type WebhookService service

// FormattedCharge ...

// Notification resource
type Notification struct {
	ID    uint  `json:"id"`
	Event Event `json:"event"`
}

// Event resource
type Event struct {
	Resource   string    `json:"resource"`
	Type       string    `json:"Type"`
	APIVersion string    `json:"api_version"`
	CreatedAt  time.Time `json:"created_at"`
	Data       Charge    `json:"data"`
}

// VerifyEvent verifies the X-BP-Webhook-Signature header's value is correct for the body passed
func (w *WebhookService) VerifyEvent(r *http.Request) bool {
	b, err := copyBody(r)
	if err != nil {
		return false
	}
	headerSignature := r.Header.Get("X-BP-Webhook-Signature")

	return (w.genHash(b) == headerSignature)
}

// GetNotification gets the notification from a request made to the webhook
func (w *WebhookService) GetNotification(r *http.Request) (Notification, error) {

	if !w.VerifyEvent(r) {
		return Notification{}, errors.New("'X-BP-Webhook-Signature' passed not a valid signature")
	}

	var notification Notification
	b, err := copyBody(r)
	if err != nil {
		return Notification{}, err
	}

	err = json.Unmarshal(b, &notification)
	return notification, err
}

// genHash generates signature of the body for X-BP-Webhook-Signature header check
func (w *WebhookService) genHash(b []byte) string {
	return genHash(b, []byte(w.client.webhookSecret))
}

func genHash(b, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(b)
	expectedMAC := mac.Sum(nil)
	return fmt.Sprintf("%x", expectedMAC)
}

func copyBody(r *http.Request) ([]byte, error) {
	body := r.Body
	b, err := ioutil.ReadAll(body)
	defer body.Close()

	if err != nil {
		return nil, err
	}

	b1 := make([]byte, len(b))
	copy(b1, b)

	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	return b1, nil
}
