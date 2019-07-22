package pay

// ChargeCreateReq the request body when creating new charges
type ChargeCreateReq struct {
	LocalPrice  Amount                 `json:"local_price,string"`
	RedirectURL string                 `json:"redirect_url"`
	CancelURL   string                 `json:"cancel_url"`
	MetaData    map[string]interface{} `json:"meta_data"`
}

// DispenseReq the request body when making dispense
type DispenseReq struct {
	Amount   float64 `json:"amount,string"`
	Address  string  `json:"address"`
	Currency string  `json:"currency"`
}

// Error the error response from busha pay
type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
