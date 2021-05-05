package kraken

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/go-resty/resty/v2"
)

type api struct {
	clientAssets  *resty.Client
	doneAssets    chan error
	clientLedgers *resty.Client
	doneLedgers   chan error
	basePath      string
	apiKey        string
	secretKey     string
	firstTimeUsed time.Time
	ledgerTX      []ledgerTX
	assets        AssetsInfo
	txsByCategory wallet.TXsByCategory
}

type AssetsInfo struct {
	Error  []string    `json:"error"`
	Result interface{} `json:"result"`
}

func (kr *Kraken) NewAPI(apiKey, secretKey string, debug bool) {
	kr.api.txsByCategory = make(map[string]wallet.TXs)
	kr.api.clientAssets = resty.New()
	kr.api.clientAssets.SetRetryCount(3)
	kr.api.clientAssets.SetDebug(debug)
	kr.api.doneAssets = make(chan error)
	kr.api.clientLedgers = resty.New()
	kr.api.clientLedgers.SetRetryCount(3).SetRetryWaitTime(1 * time.Second)
	kr.api.clientLedgers.SetDebug(debug)
	kr.api.doneLedgers = make(chan error)
	kr.api.basePath = "https://api.kraken.com"
	kr.api.apiKey = apiKey
	kr.api.secretKey = secretKey
	kr.api.firstTimeUsed = time.Now()
}

func (api *api) getAPITxs() (err error) {
	api.getAPIAssets()
	go api.getAPISpotTrades()
	<-api.doneLedgers
	api.categorize()
	return
}

func (api *api) GetExchangeFirstUsedTime() time.Time {
	return api.firstTimeUsed
}

func (api *api) categorize() {
	for _, tx := range api.ledgerTX {
		if tx.Type == "trade" || tx.Type == "margin" || tx.Type == "rollover" || tx.Type == "transfer" || tx.Type == "settled" {
			t := wallet.TX{Timestamp: tx.Time, Note: "Kraken API : " + strings.Title(tx.Type) + "  " + tx.TxId, ID: tx.TxId}
			t.Items = make(map[string]wallet.Currencies)
			if tx.Amount.IsPositive() {
				t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Asset, Amount: tx.Amount})
			} else {
				t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Asset, Amount: tx.Amount.Neg()})
			}
			if !tx.Fee.IsZero() {
				t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Asset, Amount: tx.Fee})
			}
			api.txsByCategory["Exchanges"] = append(api.txsByCategory["Exchanges"], t)
		} else if tx.Type == "deposit" {
			t := wallet.TX{Timestamp: tx.Time}
			t.Items = make(map[string]wallet.Currencies)
			if !tx.Fee.IsZero() {
				t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Asset, Amount: tx.Fee})
			}
			t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Asset, Amount: tx.Amount})
			api.txsByCategory["Deposits"] = append(api.txsByCategory["Deposits"], t)
		} else if tx.Type == "withdrawal" {
			t := wallet.TX{Timestamp: tx.Time}
			t.Items = make(map[string]wallet.Currencies)
			if !tx.Fee.IsZero() {
				t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Asset, Amount: tx.Fee})
			}
			t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Asset, Amount: tx.Amount.Neg()})
			api.txsByCategory["Withdrawals"] = append(api.txsByCategory["Withdrawals"], t)
		} else {
			log.Println("Kraken : Unmanaged ", tx.Type)
			fmt.Println(tx)
		}
	}
}

func (api *api) sign(headers map[string]string, body url.Values, resource string) {
	sha := sha256.New()
	sha.Write([]byte(body.Get("nonce") + body.Encode()))
	shasum := sha.Sum(nil)
	b64DecodedSecret, _ := base64.StdEncoding.DecodeString(api.secretKey)
	mac := hmac.New(sha512.New, b64DecodedSecret)
	mac.Write(append([]byte(resource), shasum...))
	macsum := mac.Sum(nil)
	headers["API-Sign"] = base64.StdEncoding.EncodeToString(macsum)
}
