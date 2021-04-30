package bittrex

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"strconv"
	"time"

	"gopkg.in/resty.v1"
)

const (
	API_BASE    = "https://api.bittrex.com/"
	API_VERSION = "v3/"
)

func (btrx *Bittrex) sendRequest(apiKey string, apiSecret string, resource string, method string, request *resty.Request) (response *resty.Response, err error) {
	// Generate signature
	sha_512 := sha512.New()
	hmac512 := hmac.New(sha512.New, []byte(apiSecret))
	params := ""
	if len(request.QueryParam.Encode()) > 0 {
		params += "?" + request.QueryParam.Encode()
	}
	url := API_BASE + API_VERSION + resource
	timestamp := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
	payload := ""
	if method == "POST" {
		payload = "json.dumps(query)"
	}
	sha_512.Write([]byte(payload))
	hash := hex.EncodeToString(sha_512.Sum(nil))
	pre_signature := timestamp + url + params + method + hash
	hmac512.Write([]byte(pre_signature))
	signature := hex.EncodeToString(hmac512.Sum(nil))

	// Send request
	response, err = request.SetHeaders(map[string]string{
		"Accept":           "application/json",
		"Content-Type":     "application/json",
		"Api-Content-Hash": hash,
		"Api-Key":          apiKey,
		"Api-Signature":    signature,
		"Api-Timestamp":    timestamp,
	}).Get(url)
	return response, err
}
