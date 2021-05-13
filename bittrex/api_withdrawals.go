package bittrex

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type withdrawalTX struct {
	Time           time.Time
	ID             string
	CurrencySymbol string
	Quantity       decimal.Decimal
	Fee            decimal.Decimal
	Address        string
	Status         string
}

func (api *api) getWithdrawalsTXs() {
	const SOURCE = "Bittrex API Withdrawals :"
	withdrawals, err := api.getWithdrawals()
	if err != nil {
		api.doneWithdrawals <- err
		return
	}
	for _, wit := range withdrawals {
		tx := withdrawalTX{}
		tx.Time = wit.CompletedAt
		tx.ID = wit.ID
		tx.CurrencySymbol = wit.CurrencySymbol
		tx.Quantity, err = decimal.NewFromString(wit.Quantity)
		if err != nil {
			log.Println(SOURCE, "Error Parsing Quantity : ", wit.Quantity)
		}
		tx.Address = wit.CryptoAddress
		tx.Status = wit.Status
		if wit.TxCost != "" {
			tx.Fee, err = decimal.NewFromString(wit.TxCost)
			if err != nil {
				log.Println(SOURCE, "Error Parsing Fee : ", wit.TxCost)
			}
		}
		api.withdrawalTXs = append(api.withdrawalTXs, tx)
	}
	api.doneWithdrawals <- nil
}

type GetTransferResponse []struct {
	ID               string    `json:"id"`
	CurrencySymbol   string    `json:"currencySymbol"`
	Quantity         string    `json:"quantity"`
	CryptoAddress    string    `json:"cryptoAddress"`
	UpdatedAt        time.Time `json:"updatedAt"`
	CompletedAt      time.Time `json:"completedAt"`
	Status           string    `json:"status"`
	TxCost           string    `json:"txCost"`
	TxID             string    `json:"txId"`
	CreatedAt        time.Time `json:"createdAt"`
	CryptoAddressTag string    `json:"cryptoAddressTag,omitempty"`
}

func (api *api) getWithdrawals() (withdrawalResp GetTransferResponse, err error) {
	const SOURCE = "Bittrex API Withdrawals :"
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Bittrex/withdrawals", "closed", &withdrawalResp)
	}
	if !useCache || err != nil {
		hash := api.hash("")
		ressource := "withdrawals/closed"
		timestamp, signature := api.sign("", ressource, "GET", hash, "?pageSize=200&status=COMPLETED")
		resp, err := api.clientWithdrawals.R().
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
			SetResult(&GetTransferResponse{}).
			// SetError(&ErrorResp{}).
			Get(api.basePath + ressource)
		if err != nil {
			return withdrawalResp, errors.New(SOURCE + " Error Requesting")
		}
		if resp.StatusCode() > 300 {
			return withdrawalResp, errors.New(SOURCE + " Error StatusCode" + strconv.Itoa(resp.StatusCode()))
		}
		withdrawalResp = *resp.Result().(*GetTransferResponse)
		if useCache {
			err = db.Write("Bittrex/withdrawals", "closed", withdrawalResp)
			if err != nil {
				return withdrawalResp, errors.New(SOURCE + " Error Caching")
			}
		}
		time.Sleep(api.timeBetweenReq)
	}
	return withdrawalResp, nil
}
