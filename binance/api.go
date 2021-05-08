package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/go-resty/resty/v2"
)

type api struct {
	clientExInf        *resty.Client
	doneExInf          chan error
	clientDep          *resty.Client
	doneDep            chan error
	clientWit          *resty.Client
	doneWit            chan error
	clientSpotTra      *resty.Client
	doneSpotTra        chan error
	basePath           string
	apiKey             string
	secretKey          string
	firstTimeUsed      time.Time
	startTime          time.Time
	timeBetweenReq     time.Duration
	timeBetweenReqSpot time.Duration
	nextReqID          int64
	withdrawalTXs      []withdrawalTX
	depositTXs         []depositTX
	spotTradeTXs       []spotTradeTX
	txsByCategory      wallet.TXsByCategory
}

type ErrorResp struct {
	ID       int64  `json:"id"`
	Method   string `json:"method"`
	Code     int    `json:"code"`
	Message  string `json:"message,omitempty"`
	Original string `json:"original,omitempty"`
}

func (b *Binance) NewAPI(apiKey, secretKey string, debug bool) {
	b.api.txsByCategory = make(map[string]wallet.TXs)
	b.api.clientExInf = resty.New()
	b.api.clientExInf.SetRetryCount(3)
	b.api.clientExInf.SetDebug(debug)
	b.api.doneExInf = make(chan error)
	b.api.clientDep = resty.New()
	b.api.clientDep.SetRetryCount(3)
	b.api.clientDep.SetDebug(debug)
	b.api.doneDep = make(chan error)
	b.api.clientWit = resty.New()
	b.api.clientWit.SetRetryCount(3)
	b.api.clientWit.SetDebug(debug)
	b.api.doneWit = make(chan error)
	b.api.clientSpotTra = resty.New()
	b.api.clientSpotTra.SetRetryCount(3).SetRetryWaitTime(1 * time.Second)
	b.api.clientSpotTra.SetDebug(debug)
	b.api.doneSpotTra = make(chan error)
	b.api.basePath = "https://api.binance.com/"
	b.api.apiKey = apiKey
	b.api.secretKey = secretKey
	b.api.firstTimeUsed = time.Now()
	b.api.startTime = time.Date(2019, time.November, 14, 0, 0, 0, 0, time.UTC)
	b.api.timeBetweenReq = 100 * time.Millisecond
	b.api.timeBetweenReqSpot = time.Second
}

func (api *api) getAllTXs(loc *time.Location) (err error) {
	// 1. Récupérer le liste des paires dispo :
	// https://api.binance.com/api/v3/exchangeInfo
	// 2. Boucler sur les paires en récupérant l'historique
	// 3. Mettre en cache après chaque call
	api.getExchangeInfo()
	// go api.getDepositsTXs(loc)
	// go api.getWithdrawalsTXs(loc)
	// go api.getSpotTradesTXs(loc)
	// <-api.doneDep
	// <-api.doneWit
	// <-api.doneSpotTra
	// api.categorize()
	return
}

func (api *api) GetFirstUsedTime() time.Time {
	return api.firstTimeUsed
}

func (api *api) categorize() {
	for _, tx := range api.withdrawalTXs {
		t := wallet.TX{Timestamp: tx.Timestamp, Note: "Binance API : Withdrawal " + tx.Description}
		t.Items = make(map[string]wallet.Currencies)
		t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
		if !tx.Fee.IsZero() {
			t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: tx.Fee})
		}
		api.txsByCategory["Withdrawals"] = append(api.txsByCategory["Withdrawals"], t)
		if tx.Timestamp.Before(api.firstTimeUsed) {
			api.firstTimeUsed = tx.Timestamp
		}
	}
	for _, tx := range api.depositTXs {
		t := wallet.TX{Timestamp: tx.Timestamp, Note: "Binance API : Deposit " + tx.Description}
		t.Items = make(map[string]wallet.Currencies)
		t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
		if !tx.Fee.IsZero() {
			t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: tx.Fee})
		}
		api.txsByCategory["Deposits"] = append(api.txsByCategory["Deposits"], t)
		if tx.Timestamp.Before(api.firstTimeUsed) {
			api.firstTimeUsed = tx.Timestamp
		}
	}
	for _, tx := range api.spotTradeTXs {
		t := wallet.TX{Timestamp: tx.Timestamp, Note: "Binance API : Exchange " + tx.Description}
		t.Items = make(map[string]wallet.Currencies)
		curr := strings.Split(tx.Pair, "_")
		if tx.Side == "BUY" {
			t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: curr[1], Amount: tx.Quantity})
			t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr[0], Amount: tx.Quantity.Mul(tx.Price)})
		} else { // if tx.Side == "SELL"
			t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: curr[0], Amount: tx.Quantity})
			t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr[1], Amount: tx.Quantity.Mul(tx.Price)})
		}
		if !tx.Fee.IsZero() {
			t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.FeeCurrency, Amount: tx.Fee})
		}
		api.txsByCategory["Exchanges"] = append(api.txsByCategory["Exchanges"], t)
		if tx.Timestamp.Before(api.firstTimeUsed) {
			api.firstTimeUsed = tx.Timestamp
		}
	}
}

func (api *api) sign(request map[string]interface{}) {
	if _, ok := request["api_key"]; !ok {
		request["api_key"] = api.apiKey
	}
	if _, ok := request["timestamp"]; !ok {
		request["timestamp"] = time.Now().UTC().UnixNano() / 1e6
	}
	params := request["params"].(map[string]interface{})
	paramString := []string{}
	for _, keySorted := range api.getSortedKeys(params) {
		paramString = append(paramString, keySorted+"="+fmt.Sprintf("%v", params[keySorted]))
	}
	sigPayload := strings.Join(paramString, "&")
	key := []byte(api.secretKey)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(sigPayload))
	request["sig"] = hex.EncodeToString(mac.Sum(nil))
}

func (api *api) getSortedKeys(params map[string]interface{}) []string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
