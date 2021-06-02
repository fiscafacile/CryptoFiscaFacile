package cryptocom

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

type apiEx struct {
	clientDep          *resty.Client
	doneDep            chan error
	clientWit          *resty.Client
	doneWit            chan error
	clientSpotTra      *resty.Client
	doneSpotTra        chan error
	basePath           string
	apiKey             string
	secretKey          string
	startTime          time.Time
	firstTimeUsed      time.Time
	lastTimeUsed       time.Time
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

func (cdc *CryptoCom) NewExchangeAPI(apiKey, secretKey string, debug bool) {
	cdc.apiEx.txsByCategory = make(map[string]wallet.TXs)
	cdc.apiEx.clientDep = resty.New()
	cdc.apiEx.clientDep.SetRetryCount(3)
	cdc.apiEx.clientDep.SetDebug(debug)
	cdc.apiEx.doneDep = make(chan error)
	cdc.apiEx.clientWit = resty.New()
	cdc.apiEx.clientWit.SetRetryCount(3)
	cdc.apiEx.clientWit.SetDebug(debug)
	cdc.apiEx.doneWit = make(chan error)
	cdc.apiEx.clientSpotTra = resty.New()
	cdc.apiEx.clientSpotTra.SetRetryCount(3).SetRetryWaitTime(1 * time.Second)
	cdc.apiEx.clientSpotTra.SetDebug(debug)
	cdc.apiEx.doneSpotTra = make(chan error)
	cdc.apiEx.basePath = "https://api.crypto.com/v2/"
	cdc.apiEx.apiKey = apiKey
	cdc.apiEx.secretKey = secretKey
	cdc.apiEx.firstTimeUsed = time.Now()
	cdc.apiEx.startTime = time.Date(2019, time.November, 14, 0, 0, 0, 0, time.UTC)
	cdc.apiEx.lastTimeUsed = cdc.apiEx.startTime
	cdc.apiEx.timeBetweenReq = 100 * time.Millisecond
	cdc.apiEx.timeBetweenReqSpot = time.Second
}

func (api *apiEx) getAllTXs(loc *time.Location) (err error) {
	go api.getDepositsTXs(loc)
	go api.getWithdrawalsTXs(loc)
	go api.getSpotTradesTXs(loc)
	<-api.doneDep
	<-api.doneWit
	<-api.doneSpotTra
	fmt.Println("Finished")
	api.categorize()
	return
}

func (api *apiEx) categorize() {
	const SOURCE = "Crypto.com Exchange API :"
	alreadyAsked := []string{}
	for _, tx := range api.withdrawalTXs {
		t := wallet.TX{Timestamp: tx.Timestamp, ID: tx.ID, Note: SOURCE + " Withdrawal " + tx.Description}
		t.Items = make(map[string]wallet.Currencies)
		t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
		if !tx.Fee.IsZero() {
			t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: tx.Fee})
		}
		api.txsByCategory["Withdrawals"] = append(api.txsByCategory["Withdrawals"], t)
		if tx.Timestamp.Before(api.firstTimeUsed) {
			api.firstTimeUsed = tx.Timestamp
		}
		if tx.Timestamp.After(api.lastTimeUsed) {
			api.lastTimeUsed = tx.Timestamp
		}
	}
	for _, tx := range api.depositTXs {
		t := wallet.TX{Timestamp: tx.Timestamp, ID: tx.ID, Note: SOURCE + " Deposit " + tx.Description}
		t.Items = make(map[string]wallet.Currencies)
		t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
		if !tx.Fee.IsZero() {
			t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: tx.Fee})
		}
		api.txsByCategory["Deposits"] = append(api.txsByCategory["Deposits"], t)
		if tx.Timestamp.Before(api.firstTimeUsed) {
			api.firstTimeUsed = tx.Timestamp
		}
		if tx.Timestamp.After(api.lastTimeUsed) {
			api.lastTimeUsed = tx.Timestamp
		}
	}
	for _, tx := range api.spotTradeTXs {
		t := wallet.TX{Timestamp: tx.Timestamp, ID: tx.ID, Note: SOURCE + " Exchange " + tx.Description}
		t.Items = make(map[string]wallet.Currencies)
		curr := strings.Split(tx.Pair, "_")
		if !tx.Fee.IsZero() {
			t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.FeeCurrency, Amount: tx.Fee})
		}
		if tx.Side == "BUY" {
			t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: curr[1], Amount: tx.Quantity.Mul(tx.Price)})
			t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr[0], Amount: tx.Quantity})
			api.txsByCategory["Exchanges"] = append(api.txsByCategory["Exchanges"], t)
		} else if tx.Side == "SELL" {
			t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: curr[0], Amount: tx.Quantity})
			t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr[1], Amount: tx.Quantity.Mul(tx.Price)})
			api.txsByCategory["Exchanges"] = append(api.txsByCategory["Exchanges"], t)
		} else {
			alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Side, tx, alreadyAsked)
		}
		if tx.Timestamp.Before(api.firstTimeUsed) {
			api.firstTimeUsed = tx.Timestamp
		}
		if tx.Timestamp.After(api.lastTimeUsed) {
			api.lastTimeUsed = tx.Timestamp
		}
	}
}

func (api *apiEx) sign(request map[string]interface{}) {
	if _, ok := request["id"]; !ok {
		request["id"] = api.nextReqID
		api.nextReqID += 1
	}
	if _, ok := request["api_key"]; !ok {
		request["api_key"] = api.apiKey
	}
	if _, ok := request["nonce"]; !ok {
		request["nonce"] = time.Now().UTC().UnixNano() / 1e6
	}
	params := request["params"].(map[string]interface{})
	paramString := ""
	for _, keySorted := range api.getSortedKeys(params) {
		paramString += keySorted + fmt.Sprintf("%v", params[keySorted])
	}
	sigPayload := fmt.Sprintf("%v%v%s%s%v", request["method"], request["id"], api.apiKey, paramString, request["nonce"])
	key := []byte(api.secretKey)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(sigPayload))
	request["sig"] = hex.EncodeToString(mac.Sum(nil))
}

func (api *apiEx) getSortedKeys(params map[string]interface{}) []string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
