package pay

import (
	"fmt"
	"time"
)

// ChargeService service user by charge
type ChargeService service

type Charge struct {
	ID          string                 `json:"id"`
	Addresses   map[string]string      `json:"addresses"`
	Pricing     map[string]Pricing     `json:"pricing"`
	Payments    []Payment              `json:"payments"`
	CreatedAt   time.Time              `json:"created_at"`
	Code        string                 `json:"code"`
	HostedURL   string                 `json:"hosted_url"`
	RedirectURL string                 `json:"redirect_url"`
	CancelURL   string                 `json:"cancel_url"`
	ExpiresAt   time.Time              `json:"expires_at"`
	ConfirmedAt time.Time              `json:"confirmed_at"`
	Timeline    []ChargeStatus         `json:"timelines"`
	Meta        map[string]interface{} `json:"metadata"`
}

// Pricing ...
type Pricing struct {
	ID        string    `json:"-" db:"id"`
	ChargeID  string    `json:"-" db:"charge_id"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
	Currency  string    `json:"currency" db:"currency"`
	Value     float64   `json:"value,string" db:"value"`
	Rate      float64   `json:"rate,string" db:"rate"`
}

// ChargeStatus ...
type ChargeStatus struct {
	CreatedAt time.Time `json:"time"`
	Status    string    `json:"status"`
	Context   string    `json:"context"`
}

type Payment struct {
	CreatedAt       time.Time `json:"created_at"`
	Network         string    `json:"network"`
	TransactionID   string    `json:"transaction_id"`
	TransactionHash string    `json:"transaction_hash"`
	Status          string    `json:"status"`
	DetectedAt      time.Time `json:"detected_at"`
	Traded          bool      `json:"traded"`
	Value           Balance   `json:"value" bson:"value"`
}

type Amount struct {
	Amount   float64 `json:"amount,string" bson:"amount"`
	Currency string  `json:"currency" bson:"currency"`
}

type Balance struct {
	Local  *Amount `json:"local" bson:"local"`
	Crypto *Amount `json:"crypto,omitempty"`
}

func (c *ChargeService) List(page, limit int) ([]Charge, error) {
	var cc = []Charge{}

	err := c.client.Call("GET", fmt.Sprintf("/charges?page=%d&limit=%d", page, limit), nil, &cc)
	return cc, err
}

func (c *ChargeService) Get(id string) (Charge, error) {
	var ch Charge

	err := c.client.Call("GET", fmt.Sprintf("/charges/%s", id), nil, &ch)
	return ch, err
}

func (c *ChargeService) Create(cc ChargeCreate) (Charge, error) {
	var ch Charge

	err := c.client.Call("POST", fmt.Sprintf("/charges"), cc, &ch)
	return ch, err
}

// Cancel
func (c *ChargeService) Cancel(id string) (Charge, error) {
	var ch Charge

	err := c.client.Call("POST", fmt.Sprintf("/charges/%s/cancel", id), nil, &ch)
	return ch, err
}

// Dispense
func (c *ChargeService) Dispense(d DispenseReq) (Response, error) {
	var resp Response

	err := c.client.Call("POST", fmt.Sprintf("/charges/faucet/dispense"), d, &resp)
	return resp, err
}
