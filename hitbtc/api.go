package hitbtc

import (
	// "strings"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/go-resty/resty/v2"
)

type api struct {
	clientAccTrans *resty.Client
	doneAccTrans   chan error
	clientTrade    *resty.Client
	doneTrade      chan error
	basePath       string
	apiKey         string
	secretKey      string
	firstTimeUsed  time.Time
	timeBetweenReq time.Duration
	accountTXs     []accountTX
	tradeTXs       []tradeTX
	txsByCategory  wallet.TXsByCategory
}

type ErrorResp struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (hb *HitBTC) NewAPI(apiKey, secretKey string, debug bool) {
	hb.api.txsByCategory = make(map[string]wallet.TXs)
	hb.api.clientAccTrans = resty.New()
	hb.api.clientAccTrans.SetRetryCount(3)
	hb.api.clientAccTrans.SetDebug(debug)
	hb.api.doneAccTrans = make(chan error)
	hb.api.clientTrade = resty.New()
	hb.api.clientTrade.SetRetryCount(3).SetRetryWaitTime(1 * time.Second)
	hb.api.clientTrade.SetDebug(debug)
	hb.api.doneTrade = make(chan error)
	hb.api.basePath = "https://api.hitbtc.com/api/2/"
	hb.api.apiKey = apiKey
	hb.api.secretKey = secretKey
	hb.api.firstTimeUsed = time.Now()
	hb.api.timeBetweenReq = 100 * time.Millisecond
}

func (api *api) getAllTXs() (err error) {
	go api.getAccountTXs()
	go api.getTradesTXs()
	<-api.doneAccTrans
	<-api.doneTrade
	api.categorize()
	return
}

func (api *api) GetExchangeFirstUsedTime() time.Time {
	return api.firstTimeUsed
}

func (api *api) categorize() {
	const SOURCE = "HitBTC API :"
	alreadyAsked := []string{}
	for _, tx := range api.accountTXs {
		t := wallet.TX{Timestamp: tx.UpdatedAt, ID: tx.ID, Note: SOURCE + " " + tx.Type}
		t.Items = make(map[string]wallet.Currencies)
		if !tx.Fee.IsZero() {
			t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: tx.Fee})
		}
		if tx.Type == "deposit" ||
			tx.Type == "payin" {
			t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
			api.txsByCategory["Deposits"] = append(api.txsByCategory["Deposits"], t)
		} else if tx.Type == "withdraw" ||
			tx.Type == "payout" {
			if tx.Type == "payout" {
				t.Note += " " + tx.Hash + " -> " + tx.Address
			}
			t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
			api.txsByCategory["Withdrawals"] = append(api.txsByCategory["Withdrawals"], t)
		} else if tx.Type == "bankToExchange" ||
			tx.Type == "exchangeToBank" {
			// Ignore Source internal transfer
		} else {
			alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Type, tx, alreadyAsked)
		}
		if tx.UpdatedAt.Before(api.firstTimeUsed) {
			api.firstTimeUsed = tx.UpdatedAt
		}
	}
	// for _, tx := range api.tradeTXs {
	// 	t := wallet.TX{Timestamp: tx.Timestamp, Note: "Crypto.com Exchange API : Exchange " + tx.Description}
	// 	t.Items = make(map[string]wallet.Currencies)
	// 	curr := strings.Split(tx.Pair, "_")
	// 	if tx.Side == "BUY" {
	// 		t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: curr[1], Amount: tx.Quantity})
	// 		t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr[0], Amount: tx.Quantity.Mul(tx.Price)})
	// 	} else { // if tx.Side == "SELL"
	// 		t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: curr[0], Amount: tx.Quantity})
	// 		t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr[1], Amount: tx.Quantity.Mul(tx.Price)})
	// 	}
	// 	if !tx.Fee.IsZero() {
	// 		t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.FeeCurrency, Amount: tx.Fee})
	// 	}
	// 	api.txsByCategory["Exchanges"] = append(api.txsByCategory["Exchanges"], t)
	// 	if tx.Timestamp.Before(api.firstTimeUsed) {
	// 		api.firstTimeUsed = tx.Timestamp
	// 	}
	// }
}
