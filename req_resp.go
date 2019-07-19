package pay

type ChargeCreate struct {
	LocalPrice  Amount                 `json:"local_price,string"`
	RedirectURL string                 `json:"redirect_url"`
	CancelURL   string                 `json:"cancel_url"`
	MetaData    map[string]interface{} `json:"meta_data"`
}

type DispenseReq struct {
	Amount   float64 `json:"amount,string"`
	Address  string  `json:"address"`
	Currency string  `json:"currency"`
}

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
