# Pay-Go

# Introduction

`pay-go` is a golang API wrapper for `busha pay` service. check docs here [https://docs.api.pay.busha.co](https://docs.api.pay.busha.co)

```go

package main

import (
	"fmt"
	"net/url"

	"github.com/bushaHQ/pay-go"
)

func main() {

	apiKey := "DEV_VZJjh28yYnGB"
	c := pay.NewClient(apiKey, nil)
	c.BaseURL, _ = url.Parse("https://api.staging.pay.busha.co")

	fmt.Println(c.Charge.Create(pay.ChargeCreate{
		LocalPrice: pay.Amount{
			Amount:   1000,
			Currency: "NGN",
		},
		RedirectURL: "http://Loling.com",
		CancelURL:   "http://Loling.com",
		MetaData: map[string]interface{}{
			"hello": "world",
		},
	}))
	ch, err := c.Charge.List(1, 2)
	fmt.Println(len(ch), err)
	fmt.Println(c.Charge.Get(ch[0].Code))
	fmt.Println(c.Charge.Cancel(ch[0].ID))
	fmt.Println(c.Charge.Dispense(pay.DispenseReq{
		Amount:   0.0001,
		Address:  "Shalla",
		Currency: "BTC",
	}))
}

```