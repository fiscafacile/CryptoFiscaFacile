package bittrex

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/go-resty/resty/v2"
)

type api struct {
	clientDeposits    *resty.Client
	doneDeposits      chan error
	clientWithdrawals *resty.Client
	doneWithdrawals   chan error
	clientTrades      *resty.Client
	doneTrades        chan error
	basePath          string
	apiKey            string
	secretKey         string
	firstTimeUsed     time.Time
	lastTimeUsed      time.Time
	timeBetweenReq    time.Duration
	depositTXs        []depositTX
	withdrawalTXs     []withdrawalTX
	tradeTXs          []tradeTX
	txsByCategory     wallet.TXsByCategory
}

func (btrx *Bittrex) NewAPI(apiKey, secretKey string, debug bool) {
	btrx.api.txsByCategory = make(map[string]wallet.TXs)
	btrx.api.clientDeposits = resty.New()
	btrx.api.clientDeposits.SetRetryCount(3)
	btrx.api.clientDeposits.SetDebug(debug)
	btrx.api.doneDeposits = make(chan error)
	btrx.api.clientWithdrawals = resty.New()
	btrx.api.clientWithdrawals.SetRetryCount(3)
	btrx.api.clientWithdrawals.SetDebug(debug)
	btrx.api.doneWithdrawals = make(chan error)
	btrx.api.clientTrades = resty.New()
	btrx.api.clientTrades.SetRetryCount(3)
	btrx.api.clientTrades.SetDebug(debug)
	btrx.api.doneTrades = make(chan error)
	btrx.api.basePath = "https://api.bittrex.com/v3/"
	btrx.api.apiKey = apiKey
	btrx.api.secretKey = secretKey
	btrx.api.firstTimeUsed = time.Now()
	btrx.api.lastTimeUsed = time.Date(2019, time.November, 14, 0, 0, 0, 0, time.UTC)
	btrx.api.timeBetweenReq = 100 * time.Millisecond
}

func (api *api) getAllTXs() (err error) {
	go api.getDepositsTXs()
	go api.getWithdrawalsTXs()
	go api.getTradesTXs()
	<-api.doneDeposits
	<-api.doneWithdrawals
	<-api.doneTrades
	api.categorize()
	return
}

func (api *api) categorize() {
	const SOURCE = "Bittrex API :"
	alreadyAsked := []string{}
	symRplcr := strings.NewReplacer(
		"REPV2", "REP",
	)
	for _, tx := range api.tradeTXs {
		t := wallet.TX{Timestamp: tx.Time, ID: tx.ID, Note: SOURCE + " " + tx.Direction}
		symbolSlice := strings.Split(tx.MarketSymbol, "-")
		t.Items = make(map[string]wallet.Currencies)
		if tx.Direction == "BUY" {
			t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: symRplcr.Replace(symbolSlice[0]), Amount: tx.FillQuantity})
			t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: symRplcr.Replace(symbolSlice[1]), Amount: tx.Proceeds})
			if !tx.Commission.IsZero() {
				t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: symRplcr.Replace(symbolSlice[1]), Amount: tx.Commission})
			}
		} else if tx.Direction == "SELL" {
			t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: symRplcr.Replace(symbolSlice[0]), Amount: tx.FillQuantity})
			t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: symRplcr.Replace(symbolSlice[1]), Amount: tx.Proceeds})
			if !tx.Commission.IsZero() {
				t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: symRplcr.Replace(symbolSlice[1]), Amount: tx.Commission})
			}
		} else {
			alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Direction, tx, alreadyAsked)
		}
		api.txsByCategory["Exchanges"] = append(api.txsByCategory["Exchanges"], t)
		if tx.Time.Before(api.firstTimeUsed) {
			api.firstTimeUsed = tx.Time
		}
		if tx.Time.After(api.lastTimeUsed) {
			api.lastTimeUsed = tx.Time
		}
	}
	// Process transfer transactions
	for _, tx := range api.depositTXs {
		t := wallet.TX{Timestamp: tx.Time, ID: tx.ID, Note: SOURCE + " " + tx.Address}
		t.Items = make(map[string]wallet.Currencies)
		t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.CurrencySymbol, Amount: tx.Quantity})
		api.txsByCategory["Deposits"] = append(api.txsByCategory["Deposits"], t)
	}
	for _, tx := range api.withdrawalTXs {
		t := wallet.TX{Timestamp: tx.Time, ID: tx.ID, Note: SOURCE + " " + tx.Address}
		t.Items = make(map[string]wallet.Currencies)
		t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.CurrencySymbol, Amount: tx.Quantity})
		if !tx.Fee.IsZero() {
			t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.CurrencySymbol, Amount: tx.Fee})
		}
		api.txsByCategory["Withdrawals"] = append(api.txsByCategory["Withdrawals"], t)
	}
}

func (api *api) hash(payload string) string {
	sha_512 := sha512.New()
	sha_512.Write([]byte(payload))
	return hex.EncodeToString(sha_512.Sum(nil))
}

func (api *api) sign(timestamp, ressource, method, hash, queryParamEncoded string) (string, string) {
	hmac512 := hmac.New(sha512.New, []byte(api.secretKey))
	url := api.basePath + ressource
	if timestamp == "" {
		timestamp = strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
	}
	pre_signature := timestamp + url + queryParamEncoded + method + hash
	hmac512.Write([]byte(pre_signature))
	return timestamp, hex.EncodeToString(hmac512.Sum(nil))
}
