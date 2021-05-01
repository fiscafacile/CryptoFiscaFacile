package bittrex

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/category"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
	"gopkg.in/resty.v1"
)

type tradeRequestParams struct {
	PageSize      string `json:"pageSize,omitempty"`
	NextPageToken string `json:"nextPageToken,omitempty"`
}

type tradeResponse struct {
	ID            string `json:"id"`
	Marketsymbol  string `json:"marketSymbol"`
	Direction     string `json:"direction"`
	Type          string `json:"type"`
	Quantity      string `json:"quantity"`
	Limit         string `json:"limit"`
	Ceiling       string `json:"ceiling"`
	Timeinforce   string `json:"timeInForce"`
	Clientorderid string `json:"clientOrderId"`
	Fillquantity  string `json:"fillQuantity"`
	Commission    string `json:"commission"`
	Proceeds      string `json:"proceeds"`
	Status        string `json:"status"`
	Createdat     string `json:"createdAt"`
	Updatedat     string `json:"updatedAt"`
	Closedat      string `json:"closedAt"`
	Ordertocancel struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"orderToCancel"`
}

type apiTradeTX struct {
	Time       time.Time
	Operation  string
	FromSymbol string
	ToSymbol   string
	FromAmount decimal.Decimal
	ToAmount   decimal.Decimal
	Fee        decimal.Decimal
	ID         string
}

func (btrx *Bittrex) getTrades(apiKey, apiSecret string, pageSize int, lastObjectId string) (tradeTx *resty.Response, err error) {
	btrx.api.client = resty.New()
	// Prepare body payload
	requestParams := &tradeRequestParams{
		PageSize:      fmt.Sprint(pageSize),
		NextPageToken: lastObjectId,
	}
	// Convert params struct to json
	jsonParams, _ := json.Marshal(requestParams)
	// Convert json to map
	map_data := make(map[string]string)
	json.Unmarshal([]byte(jsonParams), &map_data)
	// Generate signature
	request := btrx.api.client.R().SetQueryParams(map_data)
	response, err := btrx.sendRequest(apiKey, apiSecret, "orders/closed", "GET", request)
	return response, err
}

func (btrx *Bittrex) GetAllTradeTXs(apiKey, apiSecret string, cat category.Category) {
	useCache := true
	var tradeTx, fullTradeTx []tradeResponse
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Bittrex", "trades", &fullTradeTx)
	}
	if !useCache || err != nil {
		// Request API with a page offset until all the records are retrieved
		lastObjectId := ""
		lastResponseCount := 200
		pageLength := 200
		for lastResponseCount == pageLength {
			trades, err := btrx.getTrades(apiKey, apiSecret, pageLength, lastObjectId)
			if err != nil {
				time.Sleep(6 * time.Second)
				trades, err = btrx.getTrades(apiKey, apiSecret, pageLength, lastObjectId)
				if err != nil {
					log.Println("Bittrex API : Error while fetching trades", err)
				}
			}
			json.Unmarshal(trades.Body(), &tradeTx)
			lastResponseCount = len(tradeTx)
			lastObjectId = tradeTx[len(tradeTx)-1].ID
			fullTradeTx = append(fullTradeTx, tradeTx...)
		}
		if useCache {
			err = db.Write("Bittrex", "trades", fullTradeTx)
			if err != nil {
				log.Println("Bittrex API : Error while caching trades", err)
			}
		}
	}
	// Process transfer transactions
	for _, trd := range fullTradeTx {
		tx := apiTradeTX{}
		tx.Time, err = time.Parse("2006-01-02T15:04:05.99Z", trd.Closedat)
		if err != nil {
			log.Println("Error Parsing Time : ", trd.Closedat)
		}
		symbolSlice := strings.Split(trd.Marketsymbol, "-")
		tx.Operation = trd.Direction
		tx.Fee, err = decimal.NewFromString(trd.Commission)
		if err != nil {
			log.Println("Error Parsing Amount : ", trd.Commission)
		}
		tx.ID = trd.ID
		if tx.Operation == "BUY" || tx.Operation == "SELL" {
			if tx.Operation == "BUY" {
				tx.FromSymbol = symbolSlice[1]
				tx.FromAmount, err = decimal.NewFromString(trd.Proceeds)
				if err != nil {
					log.Println("Error Parsing Amount : ", trd.Proceeds)
				}
				tx.ToSymbol = symbolSlice[0]
				tx.ToAmount, err = decimal.NewFromString(trd.Fillquantity)
				if err != nil {
					log.Println("Error Parsing Amount : ", trd.Fillquantity)
				}
			} else if tx.Operation == "SELL" {
				tx.FromSymbol = symbolSlice[0]
				tx.FromAmount, err = decimal.NewFromString(trd.Fillquantity)
				if err != nil {
					log.Println("Error Parsing Amount : ", trd.Fillquantity)
				}
				tx.ToSymbol = symbolSlice[1]
				tx.ToAmount, err = decimal.NewFromString(trd.Proceeds)
				if err != nil {
					log.Println("Error Parsing Amount : ", trd.Proceeds)
				}
			}
			found := false
			for i := range btrx.TXsByCategory["Exchanges"] {
				if tx.ID == btrx.TXsByCategory["Exchanges"][i].ID {
					found = true
				}
			}
			if !found {
				// fmt.Println("Nouvelle transaction :", tx)
				// fmt.Println(tx.Time, "\t", tx.Operation, "\t", "FROM", tx.FromAmount, tx.FromSymbol, "TO", tx.ToAmount, tx.ToSymbol)
				t := wallet.TX{Timestamp: tx.Time, Note: "Bittrex API : " + tx.Operation + " TxID " + tx.ID, ID: tx.ID}
				t.Items = make(map[string]wallet.Currencies)
				t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.FromSymbol, Amount: tx.FromAmount})
				t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.ToSymbol, Amount: tx.ToAmount})
				t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.FromSymbol, Amount: tx.Fee})
				btrx.TXsByCategory["Exchanges"] = append(btrx.TXsByCategory["Exchanges"], t)
			} else {
				// fmt.Println("Transaction déjà enregistrée : ", tx.ID)

			}
		} else {
			log.Println("Bittrex API : Unmanaged operation -> ", tx.Operation)
		}
	}
	btrx.tradesDone <- err
}
