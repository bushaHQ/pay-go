package pay

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
)

const (
	baseURL           = "https://api.pay.busha.co"
	defautHttpTimeout = 60 * time.Second
)

type service struct {
	client *Client
}

type logger interface {
	Println(...interface{})
}

type Client struct {
	common     service
	httpClient *http.Client
	key        string
	BaseURL    *url.URL

	Charge  *ChargeService
	Webhook *WebhookService

	LogEnabled bool
	Logger     logger
}

// Response response from busha-pay
type Response map[string]interface{}

// NewClient creates a new busha-pay
func NewClient(key string, httpClient *http.Client) *Client {

	if httpClient == nil {
		httpClient = &http.Client{Timeout: defautHttpTimeout}
	}
	u, _ := url.Parse(baseURL)

	client := &Client{
		httpClient: httpClient,
		key:        key,
		BaseURL:    u,
		Logger:     log.New(os.Stdout, "", log.LstdFlags),
	}

	client.common.client = client

	client.Charge = (*ChargeService)(&client.common)
	client.Webhook = (*WebhookService)(&client.common)

	return client
}

// Call make the http request to busha pay.method, the http method. path is the url path.
// body is the http request body. v is the interface the response will be written to
func (c *Client) Call(method, path string, body, v interface{}) error {
	var buf bytes.Buffer

	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		if err != nil {
			return err
		}
	}

	u, _ := c.BaseURL.Parse(path)
	req, err := http.NewRequest(method, u.String(), &buf)
	if err != nil {
		if c.LogEnabled {
			c.Logger.Println("Request error:", err)
		}
		return err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("X-BP-Api-Key", c.key)
	// req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if c.LogEnabled {
			c.Logger.Println("Response error:", err)
		}
		return err
	}
	defer resp.Body.Close()

	return c.decodeResponse(resp, v)
}

func (c *Client) decodeResponse(r *http.Response, v interface{}) error {

	var resp Response
	respBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if c.LogEnabled {
			c.Logger.Println("Response error:", err)
		}
		return err
	}

	json.Unmarshal(respBody, &resp)

	if r.StatusCode/100 != 2 {
		message := r.Status
		if v, ok := resp["error"]; ok {
			var errResp Error
			mapstruct(v, &errResp)
			message = errResp.Message
		}

		if c.LogEnabled {
			c.Logger.Println("Busha Pay error:", message)
		}
		return errors.New(message)
	}

	if data, ok := resp["data"]; ok {
		buf, _ := json.Marshal(data)
		return json.Unmarshal(buf, v)
	}
	// if response data does not contain data key, map entire response to v
	buf, _ := json.Marshal(resp)
	return json.Unmarshal(buf, v)
}

func mapstruct(data interface{}, v interface{}) error {
	config := &mapstructure.DecoderConfig{
		Result:           v,
		TagName:          "json",
		WeaklyTypedInput: true,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	err = decoder.Decode(data)
	return err
}
