package kraken

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

const (
	API_BASE = "https://api.kraken.com"
)

type TradesHistory struct {
	Error  []string    `json:"error"`
	Result interface{} `json:"result"`
}

func (krkn *Kraken) sendRequest(apiKey, apiSecret, resource string, payload url.Values) (response *resty.Response, err error) {
	// Generate signature
	sha := sha256.New()
	sha.Write([]byte(payload.Get("nonce") + payload.Encode()))
	shasum := sha.Sum(nil)
	b64DecodedSecret, _ := base64.StdEncoding.DecodeString(apiSecret)
	mac := hmac.New(sha512.New, b64DecodedSecret)
	mac.Write(append([]byte(resource), shasum...))
	macsum := mac.Sum(nil)
	signature := base64.StdEncoding.EncodeToString(macsum)

	client := resty.New()
	response, err = client.R().
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"API-Key":      apiKey,
			"API-Sign":     signature,
		}).
		SetFormDataFromValues(payload).
		SetResult(&TradesHistory{}).
		Post(API_BASE + resource)
	errorMsg := ""
	if err != nil {
		fmt.Println("Error while fetching trade history", err)
	} else if len(response.Result().(*TradesHistory).Error) > 0 {

		for _, msg := range response.Result().(*TradesHistory).Error {
			errorMsg += msg + "\n"
		}
	}
	if errorMsg != "" {
		err = errors.New(errorMsg)
	}
	return response, err
}
