package bittrex

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/category"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
	"gopkg.in/resty.v1"
)

const (
	API_BASE    = "https://api.bittrex.com/"
	API_VERSION = "v3/"
)

type API struct {
	client *resty.Client
}

type DepositRequestParams struct {
	Status   string `json:"status,omitempty"`
	PageSize int    `json:"pageSize,omitempty"`
}

type DepositResponse struct {
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

type ApiTXTransfer struct {
	Time     time.Time
	Currency string
	Amount   decimal.Decimal
	Fee      decimal.Decimal
	Address  string
	Status   string
}

func (btrx *Bittrex) SendRequest(apiKey string, apiSecret string, resource string, method string, request *resty.Request) (response *resty.Response, err error) {
	// Generate signature
	sha_512 := sha512.New()
	hmac512 := hmac.New(sha512.New, []byte(apiSecret))
	params := ""
	if len(request.QueryParam.Encode()) > 0 {
		params += "?" + request.QueryParam.Encode()
	}
	url := API_BASE + API_VERSION + resource
	timestamp := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
	payload := ""
	if method == "POST" {
		payload = "json.dumps(query)"
	}
	sha_512.Write([]byte(payload))
	hash := hex.EncodeToString(sha_512.Sum(nil))
	pre_signature := timestamp + url + params + method + hash
	hmac512.Write([]byte(pre_signature))
	signature := hex.EncodeToString(hmac512.Sum(nil))

	// Send request
	response, err = request.SetHeaders(map[string]string{
		"Accept":           "application/json",
		"Content-Type":     "application/json",
		"Api-Content-Hash": hash,
		"Api-Key":          apiKey,
		"Api-Signature":    signature,
		"Api-Timestamp":    timestamp,
	}).Get(url)
	return response, err
}

func (btrx *Bittrex) GetDeposits(apiKey string, apiSecret string) (depositTx *resty.Response, err error) {
	btrx.api.client = resty.New()
	// Prepare body payload
	requestParams := &DepositRequestParams{
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
	response, err := btrx.SendRequest(apiKey, apiSecret, "deposits/closed", "GET", request)
	return response, err
}

func (btrx *Bittrex) GetAllTXs(apiKey string, apiSecret string, cat category.Category) {
	useCache := true
	var depositTx []DepositResponse
	// var depositTx []map[string]interface{}
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Bittrex", "deposit", &depositTx)
	}
	if !useCache || err != nil {
		deposit, err := btrx.GetDeposits(apiKey, apiSecret)
		if err != nil {
			time.Sleep(6 * time.Second)
			deposit, err = btrx.GetDeposits(apiKey, apiSecret)
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
	}
	// Process transfer transactions
	for _, dep := range depositTx {
		tx := ApiTXTransfer{}
		tx.Time, err = time.Parse("2006-01-02T15:04:05.000Z", dep.Completedat)
		if err != nil {
			log.Println("Error Parsing Time : ", dep.Completedat)
		}
		tx.Currency = dep.Currencysymbol
		tx.Amount, err = decimal.NewFromString(dep.Quantity)
		if err != nil {
			log.Println("Error Parsing Amount : ", dep.Quantity)
		}
		tx.Address = dep.Cryptoaddress
		tx.Status = dep.Status
		t := wallet.TX{Timestamp: tx.Time, Note: "Bittrex Transfer API : " + tx.Address}
		t.Items = make(map[string]wallet.Currencies)
		t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
		btrx.TXsByCategory["Deposits"] = append(btrx.TXsByCategory["Deposits"], t)
	}
	btrx.done <- err
}
