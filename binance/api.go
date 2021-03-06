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
	clientExInfo         *resty.Client
	doneExInfo           chan error
	clientDep            *resty.Client
	doneDep              chan error
	clientWit            *resty.Client
	doneWit              chan error
	clientSpotTra        *resty.Client
	doneSpotTra          chan error
	clientAssDiv         *resty.Client
	doneAssDiv           chan error
	basePath             string
	apiKey               string
	secretKey            string
	firstTimeUsed        time.Time
	lastTimeUsed         time.Time
	withdrawalTXs        []withdrawalTX
	depositTXs           []depositTX
	spotTradeTXs         []spotTradeTX
	assetDividendTXs     []assetDividendTX
	txsByCategory        wallet.TXsByCategory
	symbols              []Symbols
	timeBetweenRequests  time.Duration
	reqWeightlimit       int
	reqWeightInterval    string
	reqWeightIntervalNum int
	reqWeightTimeToWait  time.Duration
	debug                bool
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
	b.api.clientExInfo = resty.New()
	b.api.clientExInfo.SetRetryCount(3)
	b.api.clientExInfo.SetDebug(false)
	b.api.doneExInfo = make(chan error)
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
	b.api.clientAssDiv = resty.New()
	b.api.clientAssDiv.SetRetryCount(3).SetRetryWaitTime(1 * time.Second)
	b.api.clientAssDiv.SetDebug(debug)
	b.api.doneAssDiv = make(chan error)
	b.api.basePath = "https://api.binance.com/"
	b.api.apiKey = apiKey
	b.api.secretKey = secretKey
	b.api.firstTimeUsed = time.Now()
	b.api.lastTimeUsed = time.Date(2019, time.November, 14, 0, 0, 0, 0, time.UTC)
	b.api.debug = debug
}

func (api *api) getAllTXs(loc *time.Location) (err error) {
	api.getExchangeInfo()
	go api.getDepositsTXs(loc)
	go api.getWithdrawalsTXs(loc)
	go api.getSpotTradesTXs()
	go api.getAssetDividendTXs(loc)
	<-api.doneDep
	<-api.doneWit
	<-api.doneSpotTra
	<-api.doneAssDiv
	api.categorize()
	return
}

func (api *api) categorize() {
	for _, tx := range api.withdrawalTXs {
		t := wallet.TX{Timestamp: tx.Timestamp, ID: tx.ID, Note: "Binance API : Withdrawal " + tx.Description}
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
		t := wallet.TX{Timestamp: tx.Timestamp, ID: tx.ID, Note: "Binance API : Deposit " + tx.Description}
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
		found := false
		for i, t := range api.txsByCategory["Exchanges"] {
			if t.ID == tx.ID {
				found = true
				api.txsByCategory["Exchanges"][i].Note += " " + tx.Description
				if tx.Side == "BUY" {
					api.txsByCategory["Exchanges"][i].Items["From"] = append(api.txsByCategory["Exchanges"][i].Items["From"], wallet.Currency{Code: tx.QuoteAsset, Amount: tx.QuoteQty})
					api.txsByCategory["Exchanges"][i].Items["To"] = append(api.txsByCategory["Exchanges"][i].Items["To"], wallet.Currency{Code: tx.BaseAsset, Amount: tx.Qty})
				} else if tx.Side == "SELL" {
					api.txsByCategory["Exchanges"][i].Items["From"] = append(api.txsByCategory["Exchanges"][i].Items["From"], wallet.Currency{Code: tx.BaseAsset, Amount: tx.Qty})
					api.txsByCategory["Exchanges"][i].Items["To"] = append(api.txsByCategory["Exchanges"][i].Items["To"], wallet.Currency{Code: tx.QuoteAsset, Amount: tx.QuoteQty})
				} else {
					fmt.Println("Unknown transaction kind", tx.Side)
				}
				if !tx.Fee.IsZero() {
					api.txsByCategory["Exchanges"][i].Items["Fee"] = append(api.txsByCategory["Exchanges"][i].Items["Fee"], wallet.Currency{Code: tx.FeeCurrency, Amount: tx.Fee})
				}
			}
		}
		if !found {
			t := wallet.TX{Timestamp: tx.Timestamp, ID: tx.ID, Note: "Binance API : Exchange " + tx.Description}
			t.Items = make(map[string]wallet.Currencies)
			if tx.Side == "BUY" {
				t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.QuoteAsset, Amount: tx.QuoteQty})
				t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.BaseAsset, Amount: tx.Qty})
			} else if tx.Side == "SELL" {
				t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.BaseAsset, Amount: tx.Qty})
				t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.QuoteAsset, Amount: tx.QuoteQty})
			} else {
				fmt.Println("Unknown transaction kind", tx.Side)
			}
			if !tx.Fee.IsZero() {
				t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.FeeCurrency, Amount: tx.Fee})
			}
			api.txsByCategory["Exchanges"] = append(api.txsByCategory["Exchanges"], t)
		}
		if tx.Timestamp.Before(api.firstTimeUsed) {
			api.firstTimeUsed = tx.Timestamp
		}
		if tx.Timestamp.After(api.lastTimeUsed) {
			api.lastTimeUsed = tx.Timestamp
		}
	}
	for _, tx := range api.assetDividendTXs {
		t := wallet.TX{Timestamp: tx.Timestamp, ID: tx.ID, Note: "Binance API : Asset Dividend " + tx.Description}
		t.Items = make(map[string]wallet.Currencies)
		t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Asset, Amount: tx.Amount})
		if tx.Description == "BNB Vault" ||
			tx.Description == "Locked Staking" ||
			tx.Description == "DeFi Staking" ||
			tx.Description == "Launchpool" ||
			tx.Description == "XVS distribution" || // maybe should be better to get this list by API of by https://launchpad.binance.com/en/viewall/lp
			tx.Description == "FLM distribution" ||
			tx.Description == "CTK distribution" ||
			tx.Description == "ALPHA distribution" ||
			tx.Description == "HARD distribution" ||
			tx.Description == "UNFI distribution" ||
			tx.Description == "JUV distribution" ||
			tx.Description == "PSG distribution" ||
			tx.Description == "REEF distribution" ||
			tx.Description == "ASR distribution" ||
			tx.Description == "ATM distribution" ||
			tx.Description == "OG distribution" ||
			tx.Description == "BTCST distribution" ||
			tx.Description == "LIT distribution" ||
			tx.Description == "DODO distribution" ||
			tx.Description == "ALICE distribution" ||
			tx.Description == "TLM distribution" {
			api.txsByCategory["Minings"] = append(api.txsByCategory["Minings"], t)
		} else if tx.Description == "Flexible Savings" {
			api.txsByCategory["Interests"] = append(api.txsByCategory["Interests"], t)
		} else {
			api.txsByCategory["AirDrops"] = append(api.txsByCategory["AirDrops"], t)
		}
		if tx.Timestamp.Before(api.firstTimeUsed) {
			api.firstTimeUsed = tx.Timestamp
		}
		if tx.Timestamp.After(api.lastTimeUsed) {
			api.lastTimeUsed = tx.Timestamp
		}
	}
}

func (api *api) sign(params map[string]string) {
	paramString := []string{}
	for _, keySorted := range api.getSortedKeys(params) {
		paramString = append(paramString, keySorted+"="+fmt.Sprintf("%v", params[keySorted]))
	}
	sigPayload := strings.Join(paramString, "&")
	key := []byte(api.secretKey)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(sigPayload))
	params["signature"] = hex.EncodeToString(mac.Sum(nil))
}

func (api *api) getSortedKeys(params map[string]string) []string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
