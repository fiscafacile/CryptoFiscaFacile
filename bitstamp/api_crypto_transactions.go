package bitstamp

import (
	"errors"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type cryptoTX struct {
	DateTime           time.Time
	ID                 string
	Amount             decimal.Decimal
	Currency           string
	DestinationAddress string
}

func (api *api) getCryptoTXs() {
	const SOURCE = "Bitstamp API Crypto Transactions :"
	cryptoTXs, err := api.getCryptoTransactions()
	if err != nil {
		api.doneCryptoTrans <- err
		return
	}
	for _, t := range cryptoTXs.Deposits {
		tx := cryptoTX{}
		tx.ID = t.TxID
		tx.DateTime = time.Unix(t.DateTime, 0)
		tx.Amount = decimal.NewFromFloat(t.Amount)
		tx.Currency = t.Currency
		tx.DestinationAddress = t.DestinationAddress
		api.depositTXs = append(api.depositTXs, tx)
	}
	for _, t := range cryptoTXs.Withdrawals {
		tx := cryptoTX{}
		tx.ID = t.TxID
		tx.DateTime = time.Unix(t.DateTime, 0)
		tx.Amount = decimal.NewFromFloat(t.Amount)
		tx.Currency = t.Currency
		tx.DestinationAddress = t.DestinationAddress
		api.withdrawalTXs = append(api.withdrawalTXs, tx)
	}
	api.doneCryptoTrans <- nil
}

type CryptoTransaction struct {
	Currency           string  `json:"currency"`
	DestinationAddress string  `json:"destinationAddress"`
	TxID               string  `json:"txid"`
	Amount             float64 `json:"amount"`
	DateTime           int64   `json:"datetime"`
}

type GetCryptoTransactionsResp struct {
	Deposits    []CryptoTransaction `json:"deposits"`
	Withdrawals []CryptoTransaction `json:"withdrawals"`
}

func (api *api) getCryptoTransactions() (cryptoTXs GetCryptoTransactionsResp, err error) {
	const SOURCE = "Bitstamp API Crypto Transactions :"
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Bitstamp", "crypto-transactions", &cryptoTXs)
	}
	if !useCache || err != nil {
		url := api.basePath + "crypto-transactions/"
		req := api.clientCryptoTrans.R().
			SetFormData(map[string]string{
				"limit": "1000",
			})
		api.sign(req, "POST", url)
		resp, err := req.SetResult(&GetCryptoTransactionsResp{}).
			SetError(&ErrorResp{}).
			Post(url)
		if err != nil {
			return cryptoTXs, errors.New(SOURCE + " Error Requesting")
		}
		if resp.StatusCode() > 300 {
			return cryptoTXs, errors.New(SOURCE + " Error StatusCode" + strconv.Itoa(resp.StatusCode()))
		}
		spew.Dump(*resp)
		cryptoTXs = *resp.Result().(*GetCryptoTransactionsResp)
		if useCache {
			err = db.Write("Bitstamp", "crypto-transactions", cryptoTXs)
			if err != nil {
				return cryptoTXs, errors.New(SOURCE + " Error Caching")
			}
		}
		time.Sleep(api.timeBetweenReq)
	}
	return cryptoTXs, nil
}
