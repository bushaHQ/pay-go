package pay

import (
	"net/url"
	"os"
	"reflect"
	"testing"
)

var (
	payClient *Client
	// webhook
)

func init() {
	key := os.Getenv("PAY_TEST_KEY")
	if key == "" {
		panic("PAY_TEST_KEY must be passed")
	}
	payClient = NewClient(key, nil)
	// set base url to staging
	payClient.BaseURL, _ = url.Parse("https://api.staging.pay.busha.co")

	// Webhooks Request

}

func TestChargeService_Cancel(t *testing.T) {
	type fields struct {
		client *Client
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Charge
		wantErr bool
	}{
		{
			"Fail: ref does not exist",
			fields{client: payClient},
			args{id: "cool"},
			Charge{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ChargeService{
				client: tt.fields.client,
			}
			got, err := c.Cancel(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChargeService.Cancel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChargeService.Cancel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChargeService_List(t *testing.T) {
	type fields struct {
		client *Client
	}
	type args struct {
		page  int
		limit int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Charge
		wantErr bool
	}{
		{
			"Return Charges",
			fields{client: payClient},
			args{page: 1, limit: 10},
			[]Charge{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ChargeService{
				client: tt.fields.client,
			}
			got, err := c.List(tt.args.page, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChargeService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("ChargeService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
