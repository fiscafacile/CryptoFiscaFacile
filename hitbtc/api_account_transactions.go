package hitbtc

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type accountTX struct {
	ID            string
	Index         int64
	Type          string
	Subtype       string
	Status        string
	Currency      string
	Amount        decimal.Decimal
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Fee           decimal.Decimal
	Hash          string
	Address       string
	Confirmations int64
}

func (api *api) getAccountTXs() {
	const SOURCE = "HitBTC API Account Transactions :"
	accTXs, err := api.getAccountTransactions()
	if err != nil {
		api.doneAccTrans <- err
		return
	}
	for _, t := range accTXs {
		tx := accountTX{}
		tx.ID = t.ID
		tx.Index = t.Index
		tx.Type = t.Type
		tx.Subtype = t.Subtype
		tx.Status = t.Status
		tx.Currency = apiCurrencyCure(t.Currency)
		tx.Amount, err = decimal.NewFromString(t.Amount)
		if err != nil {
			log.Println(SOURCE, "Error Parsing Amount : ", t.Amount)
		}
		tx.CreatedAt = t.CreatedAt
		tx.UpdatedAt = t.UpdatedAt
		if t.Fee != "" {
			tx.Fee, err = decimal.NewFromString(t.Fee)
			if err != nil {
				log.Println(SOURCE, "Error Parsing Fee : ", t.Fee)
			}
		}
		tx.Hash = t.Hash
		tx.Address = t.Address
		tx.Confirmations = t.Confirmations
		api.accountTXs = append(api.accountTXs, tx)
	}
	api.doneAccTrans <- nil
}

type GetAccountTransactionsResp []struct {
	ID            string    `json:"id"`
	Index         int64     `json:"index"`
	Type          string    `json:"type"`
	Subtype       string    `json:"subType,omitempty"`
	Status        string    `json:"status"`
	Currency      string    `json:"currency"`
	Amount        string    `json:"amount"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Fee           string    `json:"fee,omitempty"`
	Hash          string    `json:"hash,omitempty"`
	Address       string    `json:"address,omitempty"`
	Confirmations int64     `json:"confirmations,omitempty"`
}

func (api *api) getAccountTransactions() (accTXs GetAccountTransactionsResp, err error) {
	const SOURCE = "HitBTC API Account Transactions :"
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("HitBTC/account", "transactions", &accTXs)
	}
	if !useCache || err != nil {
		method := "account/transactions"
		resp, err := api.clientAccTrans.R().
			SetBasicAuth(api.apiKey, api.secretKey).
			SetResult(&GetAccountTransactionsResp{}).
			SetError(&ErrorResp{}).
			Get(api.basePath + method)
		if err != nil {
			return accTXs, errors.New(SOURCE + " Error Requesting")
		}
		if resp.StatusCode() > 300 {
			return accTXs, errors.New(SOURCE + " Error StatusCode" + strconv.Itoa(resp.StatusCode()))
		}
		accTXs = *resp.Result().(*GetAccountTransactionsResp)
		if useCache {
			err = db.Write("HitBTC/account", "transactions", accTXs)
			if err != nil {
				return accTXs, errors.New(SOURCE + " Error Caching")
			}
		}
		time.Sleep(api.timeBetweenReq)
	}
	return accTXs, nil
}
