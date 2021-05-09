package bittrex

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type depositTX struct {
	Time           time.Time
	ID             string
	CurrencySymbol string
	Quantity       decimal.Decimal
	Fee            decimal.Decimal
	Address        string
	Status         string
}

func (api *api) getDepositsTXs() {
	const SOURCE = "Bittrex API Deposits :"
	deposits, err := api.getDeposits()
	if err != nil {
		api.doneDeposits <- err
		return
	}
	for _, dep := range deposits {
		tx := depositTX{}
		tx.Time = dep.CompletedAt
		tx.ID = dep.ID
		tx.CurrencySymbol = dep.CurrencySymbol
		tx.Quantity, err = decimal.NewFromString(dep.Quantity)
		if err != nil {
			log.Println(SOURCE, "Error Parsing Quantity : ", dep.Quantity)
		}
		tx.Address = dep.CryptoAddress
		tx.Status = dep.Status
		if dep.TxCost != "" {
			tx.Fee, err = decimal.NewFromString(dep.TxCost)
			if err != nil {
				log.Println(SOURCE, "Error Parsing Fee : ", dep.TxCost)
			}
		}
		api.depositTXs = append(api.depositTXs, tx)
	}
	api.doneDeposits <- nil
}

type GetDepositResponse []struct {
	ID               string    `json:"id"`
	CurrencySymbol   string    `json:"currencySymbol"`
	Quantity         string    `json:"quantity"`
	CryptoAddress    string    `json:"cryptoAddress"`
	Confirmations    int       `json:"confirmations"`
	UpdatedAt        time.Time `json:"updatedAt"`
	CompletedAt      time.Time `json:"completedAt"`
	Status           string    `json:"status"`
	Source           string    `json:"source"`
	TxCost           string    `json:"txCost"`
	TxID             string    `json:"txId"`
	CreatedAt        time.Time `json:"createdAt"`
	CryptoAddressTag string    `json:"cryptoAddressTag,omitempty"`
}

func (api *api) getDeposits() (depositResp GetDepositResponse, err error) {
	const SOURCE = "Bittrex API Deposits :"
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Bittrex/deposits", "closed", &depositResp)
	}
	if !useCache || err != nil {
		hash := api.hash("")
		ressource := "deposits/closed"
		timestamp, signature := api.sign("", ressource, "GET", hash, "?pageSize=200&status=COMPLETED")
		resp, err := api.clientDeposits.R().
			SetQueryParams(map[string]string{
				"pageSize": "200",
				"status":   "COMPLETED",
			}).
			SetHeaders(map[string]string{
				"Accept":           "application/json",
				"Content-Type":     "application/json",
				"Api-Content-Hash": hash,
				"Api-Key":          api.apiKey,
				"Api-Signature":    signature,
				"Api-Timestamp":    timestamp,
			}).
			SetResult(&GetDepositResponse{}).
			// SetError(&ErrorResp{}).
			Get(api.basePath + ressource)
		if err != nil {
			return depositResp, errors.New(SOURCE + " Error Requesting")
		}
		if resp.StatusCode() > 300 {
			return depositResp, errors.New(SOURCE + " Error StatusCode" + strconv.Itoa(resp.StatusCode()))
		}
		depositResp = *resp.Result().(*GetDepositResponse)
		if useCache {
			err = db.Write("Bittrex/deposits", "closed", depositResp)
			if err != nil {
				return depositResp, errors.New(SOURCE + " Error Caching")
			}
		}
		time.Sleep(api.timeBetweenReq)
	}
	return depositResp, nil
}
