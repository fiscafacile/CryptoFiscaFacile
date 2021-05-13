package bitstamp

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type userTX struct {
	ID         int
	DateTime   time.Time
	Type       string
	OrderID    int64
	Fee        decimal.Decimal
	Currencies map[string]decimal.Decimal
}

func (api *api) getUserTXs() {
	const SOURCE = "Bitstamp API User Transactions :"
	cryptoTXs, err := api.getUserTransactions()
	if err != nil {
		api.doneUserTrans <- err
		return
	}
	alreadyAsked := []string{}
	for _, t := range cryptoTXs {
		tx := userTX{}
		tx.ID = t.ID
		tx.DateTime, err = time.Parse("2006-01-02 15:04:05.999999", t.DateTime)
		if err != nil {
			log.Println(SOURCE, "Error Parsing DateTime : ", t.DateTime)
		}
		if t.Type == "0" {
			tx.Type = "deposit"
		} else if t.Type == "1" {
			tx.Type = "withdrawal"
		} else if t.Type == "2" {
			tx.Type = "market trade"
		} else if t.Type == "14" {
			tx.Type = "sub account transfer"
			alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Type, t, alreadyAsked)
		} else if t.Type == "25" {
			tx.Type = "credited with staked assets"
			alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Type, t, alreadyAsked)
		} else if t.Type == "26" {
			tx.Type = "sent assets to staking"
			alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Type, t, alreadyAsked)
		} else if t.Type == "27" {
			tx.Type = "staking reward"
			alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Type, t, alreadyAsked)
		} else if t.Type == "32" {
			tx.Type = "referral reward"
			alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Type, t, alreadyAsked)
		} else {
			alreadyAsked = wallet.AskForHelp(SOURCE, t, alreadyAsked)
			log.Println(SOURCE, "Error Parsing Type : ", t.Type)
		}
		tx.OrderID = t.OrderID
		if t.Fee != "" {
			tx.Fee, err = decimal.NewFromString(t.Fee)
			if err != nil {
				log.Println(SOURCE, "Error Parsing Fee : ", t.Fee)
			}
		}
		tx.Currencies = make(map[string]decimal.Decimal)
		if t.Xrp != "" {
			tx.Currencies["XRP"], err = decimal.NewFromString(t.Xrp)
			if err != nil {
				log.Println(SOURCE, "Error Parsing XRP : ", t.Xrp)
			}
		}
		if t.Ltc != "" {
			tx.Currencies["LTC"], err = decimal.NewFromString(t.Ltc)
			if err != nil {
				log.Println(SOURCE, "Error Parsing LTC : ", t.Ltc)
			}
		}
		if t.Btc != "0" {
			tx.Currencies["BTC"], err = decimal.NewFromString(string(t.Btc))
			if err != nil {
				log.Println(SOURCE, "Error Parsing BTC : ", t.Btc)
			}
		}
		if t.Eth != "" {
			tx.Currencies["ETH"], err = decimal.NewFromString(t.Eth)
			if err != nil {
				log.Println(SOURCE, "Error Parsing ETH : ", t.Eth)
			}
		}
		if t.Eur != "0" {
			tx.Currencies["EUR"], err = decimal.NewFromString(string(t.Eur))
			if err != nil {
				log.Println(SOURCE, "Error Parsing EUR : ", t.Eur)
			}
		}
		if t.Usd != 0 {
			tx.Currencies["USD"] = decimal.NewFromFloat(t.Usd)
		}
		api.userTXs = append(api.userTXs, tx)
	}
	api.doneUserTrans <- nil
}

type FloatOrString string

func (f *FloatOrString) UnmarshalJSON(d []byte) error {
	var v string
	err := json.Unmarshal(append([]byte(`"`), append(bytes.Trim(d, `"`), []byte(`"`)...)...), &v)
	*f = FloatOrString(v)
	return err
}

type GetUserTransactionsResp []struct {
	ID       int           `json:"id"`
	DateTime string        `json:"datetime"`
	Type     string        `json:"type"`
	OrderID  int64         `json:"order_id,omitempty"`
	Fee      string        `json:"fee"`
	Btc      FloatOrString `json:"btc"`
	Eth      string        `json:"eth,omitempty"`
	Eur      FloatOrString `json:"eur"`
	Ltc      string        `json:"ltc,omitempty"`
	Usd      float64       `json:"usd"`
	Xrp      string        `json:"xrp,omitempty"`
	BtcEur   float64       `json:"btc_eur,omitempty"`
	BtcUsd   string        `json:"btc_usd,omitempty"`
	EthBtc   float64       `json:"eth_btc,omitempty"`
	LtcBtc   float64       `json:"ltc_btc,omitempty"`
	XrpBtc   float64       `json:"xrp_btc,omitempty"`
}

func (api *api) getUserTransactions() (cryptoTXs GetUserTransactionsResp, err error) {
	const SOURCE = "Bitstamp API User Transactions :"
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Bitstamp", "user_transactions", &cryptoTXs)
	}
	if !useCache || err != nil {
		url := api.basePath + "user_transactions/"
		req := api.clientUserTrans.R().
			SetFormData(map[string]string{
				"limit": "1000",
			})
		api.sign(req, "POST", url)
		resp, err := req.SetResult(&GetUserTransactionsResp{}).
			SetError(&ErrorResp{}).
			Post(url)
		if err != nil {
			return cryptoTXs, errors.New(SOURCE + " Error Requesting")
		}
		if resp.StatusCode() > 300 {
			return cryptoTXs, errors.New(SOURCE + " Error StatusCode" + strconv.Itoa(resp.StatusCode()))
		}
		cryptoTXs = *resp.Result().(*GetUserTransactionsResp)
		if useCache {
			err = db.Write("Bitstamp", "user_transactions", cryptoTXs)
			if err != nil {
				return cryptoTXs, errors.New(SOURCE + " Error Caching")
			}
		}
		time.Sleep(api.timeBetweenReq)
	}
	return cryptoTXs, nil
}
