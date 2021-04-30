package bittrex

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/category"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
	"gopkg.in/resty.v1"
)

type API struct {
	client *resty.Client
}

type ApiTXTransfer struct {
	Time     time.Time
	Currency string
	Amount   decimal.Decimal
	Fee      decimal.Decimal
	Address  string
	Status   string
}

type TransferRequestParams struct {
	Status   string `json:"status,omitempty"`
	PageSize int    `json:"pageSize,omitempty"`
}

type TransferResponse struct {
	ID               string `json:"id"`
	Currencysymbol   string `json:"currencySymbol"`
	Quantity         string `json:"quantity"`
	Cryptoaddress    string `json:"cryptoAddress"`
	Cryptoaddresstag string `json:"cryptoAddressTag"`
	Txid             string `json:"txId"`
	Confirmations    string `json:"confirmations"`
	Updatedat        string `json:"updatedAt"`
	Completedat      string `json:"completedAt"`
	Status           string `json:"status"`
	Source           string `json:"source"`
	Accountid        string `json:"accountId"`
}

func (btrx *Bittrex) getDeposits(apiKey string, apiSecret string) (depositTx *resty.Response, err error) {
	btrx.api.client = resty.New()
	requestParams := &TransferRequestParams{
		Status:   "COMPLETED",
		PageSize: 200,
	}
	// Convert params struct to json
	jsonParams, _ := json.Marshal(requestParams)
	// fmt.Print(string(jsonParams))
	// Convert json to map
	map_data := make(map[string]string)
	json.Unmarshal([]byte(jsonParams), &map_data)
	// Generate signature
	request := btrx.api.client.R().SetQueryParams(map_data)
	response, err := btrx.sendRequest(apiKey, apiSecret, "deposits/closed", "GET", request)
	return response, err
}

func (btrx *Bittrex) getWithdrawals(apiKey string, apiSecret string) (withdrawalTx *resty.Response, err error) {
	btrx.api.client = resty.New()
	requestParams := &TransferRequestParams{
		Status:   "COMPLETED",
		PageSize: 200,
	}
	// Convert params struct to json
	jsonParams, _ := json.Marshal(requestParams)
	fmt.Print(string(jsonParams))
	// Convert json to map
	map_data := make(map[string]string)
	json.Unmarshal([]byte(jsonParams), &map_data)
	// Generate signature
	request := btrx.api.client.R().SetQueryParams(map_data)
	response, err := btrx.sendRequest(apiKey, apiSecret, "withdrawals/closed", "GET", request)
	return response, err
}

func (btrx *Bittrex) GetAllTransferTXs(apiKey string, apiSecret string, cat category.Category) {
	useCache := true
	var depositTx []TransferResponse
	var withdrawalTx []TransferResponse
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Bittrex", "deposit", &depositTx)
	}
	if !useCache || err != nil {
		// Retrieve and cache deposit transfers
		deposit, err := btrx.getDeposits(apiKey, apiSecret)
		if err != nil {
			time.Sleep(6 * time.Second)
			deposit, err = btrx.getDeposits(apiKey, apiSecret)
			if err != nil {
				log.Println("Bittrex API : Error while fetching deposits", err)
			}
		}
		json.Unmarshal(deposit.Body(), &depositTx)
		if useCache {
			err = db.Write("Bittrex", "deposit", depositTx)
			if err != nil {
				log.Println("Bittrex API : Error while caching deposits", err)
			}
		}
		// Retrieve and cache withdrawals transfers
		withdrawal, err := btrx.getWithdrawals(apiKey, apiSecret)
		if err != nil {
			time.Sleep(6 * time.Second)
			withdrawal, err = btrx.getWithdrawals(apiKey, apiSecret)
			if err != nil {
				log.Println("Bittrex API : Error while fetching withdrawals", err)
			}
		}
		json.Unmarshal(withdrawal.Body(), &withdrawalTx)
		if useCache {
			err = db.Write("Bittrex", "deposit", withdrawalTx)
			if err != nil {
				log.Println("Bittrex API : Error while caching withdrawals", err)
			}
		}
	}
	// Process transfer transactions
	for _, trf := range append(depositTx, withdrawalTx...) {
		tx := ApiTXTransfer{}
		tx.Time, err = time.Parse("2006-01-02T15:04:05.99Z", trf.Completedat)
		if err != nil {
			log.Println("Error Parsing Time : ", trf.Completedat)
		}
		tx.Currency = trf.Currencysymbol
		tx.Amount, err = decimal.NewFromString(trf.Quantity)
		if err != nil {
			log.Println("Error Parsing Amount : ", trf.Quantity)
		}
		tx.Address = trf.Cryptoaddress
		tx.Status = trf.Status
		t := wallet.TX{Timestamp: tx.Time, Note: "Bittrex Transfer API : " + tx.Address}
		t.Items = make(map[string]wallet.Currencies)
		if trf.Source == "" {
			t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
			btrx.TXsByCategory["Withdrawals"] = append(btrx.TXsByCategory["Withdrawals"], t)
		} else {
			t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
			btrx.TXsByCategory["Deposits"] = append(btrx.TXsByCategory["Deposits"], t)
		}
	}
	btrx.done <- err
}
