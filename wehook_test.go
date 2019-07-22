package pay

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestWebhookService_GetNotification(t *testing.T) {

	type fields struct {
		client *Client
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Notification
		wantErr bool
	}{
		{
			"Fail: key does not match",
			fields{payClient},
			args{newWebhookRequest(payClient.key, false)},
			Notification{},
			true,
		},
		{
			"Fail: key should match",
			fields{payClient},
			args{newWebhookRequest(payClient.key, true)},
			Notification{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WebhookService{
				client: tt.fields.client,
			}
			got, err := w.GetNotification(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("WebhookService.GetNotification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WebhookService.GetNotification() = %v, want %v", got, tt.want)
			}
		})
	}
}

// newWebhookRequest creates a new request webhook
func newWebhookRequest(key string, pass bool) *http.Request {
	body := `{"message":"test"}`
	buf := strings.NewReader(body)
	req, err := http.NewRequest("GET", "", buf)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if !pass {
		// invalidate the key
		key = "lalalala"
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-BP-Webhook-Signature", genHash([]byte(body), []byte(key)))

	return req
}
