package bitstamp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type api struct {
	clientCryptoTrans *resty.Client
	doneCryptoTrans   chan error
	clientUserTrans   *resty.Client
	doneUserTrans     chan error
	basePath          string
	apiKey            string
	secretKey         string
	firstTimeUsed     time.Time
	lastTimeUsed      time.Time
	timeBetweenReq    time.Duration
	depositTXs        []cryptoTX
	withdrawalTXs     []cryptoTX
	userTXs           []userTX
	txsByCategory     wallet.TXsByCategory
}

type ErrorResp struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
	Code   string `json:"code"`
}

func (bs *Bitstamp) NewAPI(apiKey, secretKey string, debug bool) {
	bs.api.txsByCategory = make(map[string]wallet.TXs)
	bs.api.clientCryptoTrans = resty.New()
	bs.api.clientCryptoTrans.SetRetryCount(3)
	bs.api.clientCryptoTrans.SetDebug(debug)
	bs.api.doneCryptoTrans = make(chan error)
	bs.api.clientUserTrans = resty.New()
	bs.api.clientUserTrans.SetRetryCount(3)
	bs.api.clientUserTrans.SetDebug(debug)
	bs.api.doneUserTrans = make(chan error)
	bs.api.basePath = "https://www.bitstamp.net/api/v2/"
	bs.api.apiKey = apiKey
	bs.api.secretKey = secretKey
	bs.api.firstTimeUsed = time.Now()
	bs.api.lastTimeUsed = time.Date(2019, time.November, 14, 0, 0, 0, 0, time.UTC)
	bs.api.timeBetweenReq = 100 * time.Millisecond
}

func (api *api) getAllTXs() (err error) {
	go api.getCryptoTXs()
	go api.getUserTXs()
	<-api.doneCryptoTrans
	<-api.doneUserTrans
	api.categorize()
	return
}

func (api *api) categorize() {
	const SOURCE = "Bitstamp API :"
	for _, tx := range api.depositTXs {
		t := wallet.TX{Timestamp: tx.DateTime, ID: tx.ID, Note: SOURCE + " Deposit from " + tx.DestinationAddress}
		t.Items = make(map[string]wallet.Currencies)
		t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
		api.txsByCategory["Deposits"] = append(api.txsByCategory["Deposits"], t)
		if tx.DateTime.Before(api.firstTimeUsed) {
			api.firstTimeUsed = tx.DateTime
		}
		if tx.DateTime.After(api.lastTimeUsed) {
			api.lastTimeUsed = tx.DateTime
		}
	}
	for _, tx := range api.withdrawalTXs {
		t := wallet.TX{Timestamp: tx.DateTime, ID: tx.ID, Note: SOURCE + " Withdrawal to " + tx.DestinationAddress}
		t.Items = make(map[string]wallet.Currencies)
		t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
		api.txsByCategory["Withdrawals"] = append(api.txsByCategory["Withdrawals"], t)
		if tx.DateTime.Before(api.firstTimeUsed) {
			api.firstTimeUsed = tx.DateTime
		}
		if tx.DateTime.After(api.lastTimeUsed) {
			api.lastTimeUsed = tx.DateTime
		}
	}
	for _, tx := range api.userTXs {
		t := wallet.TX{Timestamp: tx.DateTime, ID: strconv.Itoa(tx.ID), Note: SOURCE + " " + tx.Type}
		t.Items = make(map[string]wallet.Currencies)
		if tx.Type == "withdrawal" {
			exist := false
			for i, w := range api.txsByCategory["Withdrawals"] {
				if w.SimilarDate(time.Hour, tx.DateTime) {
					for k, v := range tx.Currencies {
						if w.Items["From"][0].Amount.Equal(v.Neg().Add(tx.Fee)) {
							exist = true
							if !tx.Fee.IsZero() {
								delete(api.txsByCategory["Withdrawals"][i].Items, "From")
								api.txsByCategory["Withdrawals"][i].Items["From"] = append(api.txsByCategory["Withdrawals"][i].Items["From"], wallet.Currency{Code: k, Amount: v.Neg()})
								api.txsByCategory["Withdrawals"][i].Items["Fee"] = append(api.txsByCategory["Withdrawals"][i].Items["Fee"], wallet.Currency{Code: k, Amount: tx.Fee})
							}
						}
					}
				}
			}
			if !exist {
				for k, v := range tx.Currencies {
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: k, Amount: v.Neg()})
					if !tx.Fee.IsZero() {
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: k, Amount: tx.Fee})
					}
				}
				api.txsByCategory["Withdrawals"] = append(api.txsByCategory["Withdrawals"], t)
			}
		} else if tx.Type == "deposit" {
			// We don't care about Deposit, nothing we already know (except Fiat deposits but not used in this tool)
		} else if tx.Type == "market trade" {
			for k, v := range tx.Currencies {
				if v.IsPositive() {
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: k, Amount: v})
				} else {
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: k, Amount: v.Neg()})
					if !tx.Fee.IsZero() {
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: k, Amount: tx.Fee})
					}
				}
			}
			api.txsByCategory["Exchanges"] = append(api.txsByCategory["Exchanges"], t)
		}
		if tx.DateTime.Before(api.firstTimeUsed) {
			api.firstTimeUsed = tx.DateTime
		}
		if tx.DateTime.After(api.lastTimeUsed) {
			api.lastTimeUsed = tx.DateTime
		}
	}
}

func (api *api) sign(req *resty.Request, method, url string) {
	header := make(map[string]string)
	header["X-Auth"] = "BITSTAMP" + " " + api.apiKey
	header["X-Auth-Nonce"] = uuid.NewString()
	header["X-Auth-Timestamp"] = strconv.FormatInt(time.Now().Add(2*time.Minute).UnixNano()/1e6, 10)
	header["X-Auth-Version"] = "v2"
	stringToSign := fmt.Sprintf("%v%s%s", header["X-Auth"], method, strings.TrimPrefix(url, "https://"))
	if req.FormData != nil {
		stringToSign += "application/x-www-form-urlencoded"
	}
	stringToSign += fmt.Sprintf("%v%v%v", header["X-Auth-Nonce"], header["X-Auth-Timestamp"], header["X-Auth-Version"])
	if req.FormData != nil {
		stringToSign += req.FormData.Encode()
	}
	key := []byte(api.secretKey)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(stringToSign))
	header["X-Auth-Signature"] = hex.EncodeToString(mac.Sum(nil))
	req.SetHeaders(header)
}
