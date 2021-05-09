package bittrex

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type tradeTX struct {
	ID           string
	MarketSymbol string
	Direction    string
	FillQuantity decimal.Decimal
	Commission   decimal.Decimal
	Proceeds     decimal.Decimal
	Time         time.Time
}

func (api *api) getTradesTXs() {
	const SOURCE = "Bittrex API Trades :"
	// Request API with a page offset until all the records are retrieved
	trades, err := api.getTrades()
	if err != nil {
		api.doneTrades <- err
		return
	}
	// Process transfer transactions
	for _, trd := range trades {
		tx := tradeTX{}
		tx.ID = trd.ID
		tx.MarketSymbol = trd.MarketSymbol
		tx.Direction = trd.Direction
		tx.FillQuantity, err = decimal.NewFromString(trd.FillQuantity)
		if err != nil {
			log.Println(SOURCE, "Error Parsing FillQuantity : ", trd.FillQuantity)
		}
		tx.Commission, err = decimal.NewFromString(trd.Commission)
		if err != nil {
			log.Println(SOURCE, "Error Parsing Commission : ", trd.Commission)
		}
		tx.Proceeds, err = decimal.NewFromString(trd.Proceeds)
		if err != nil {
			log.Println(SOURCE, "Error Parsing Proceeds : ", trd.Proceeds)
		}
		tx.Time = trd.ClosedAt
		api.tradeTXs = append(api.tradeTXs, tx)
	}
	api.doneTrades <- nil
}

type GetTradeResponse []struct {
	ID            string    `json:"id"`
	MarketSymbol  string    `json:"marketSymbol"`
	Direction     string    `json:"direction"`
	Type          string    `json:"type"`
	Quantity      string    `json:"quantity"`
	Limit         string    `json:"limit"`
	Ceiling       string    `json:"ceiling"`
	TimeInForce   string    `json:"timeInForce"`
	ClientOrderID string    `json:"clientOrderId"`
	FillQuantity  string    `json:"fillQuantity"`
	Commission    string    `json:"commission"`
	Proceeds      string    `json:"proceeds"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ClosedAt      time.Time `json:"closedAt"`
	OrderToCancel struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"orderToCancel"`
}

func (api *api) getTrades() (tradesResp GetTradeResponse, err error) {
	const SOURCE = "Bittrex API Trades :"
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Bittrex/orders", "closed", &tradesResp)
	}
	if !useCache || err != nil {
		hash := api.hash("")
		ressource := "orders/closed"
		lastObjectId := ""
		lastResponseCount := 200
		for lastResponseCount == 200 {
			paramsEncoded := "?"
			if lastObjectId != "" {
				paramsEncoded += "nextPageToken=" + lastObjectId + "&"
			}
			paramsEncoded += "pageSize=200"
			timestamp, signature := api.sign("", ressource, "GET", hash, paramsEncoded)
			params := map[string]string{
				"pageSize": "200",
			}
			if lastObjectId != "" {
				params["nextPageToken"] = lastObjectId
			}
			resp, err := api.clientTrades.R().
				SetQueryParams(params).
				SetHeaders(map[string]string{
					"Accept":           "application/json",
					"Content-Type":     "application/json",
					"Api-Content-Hash": hash,
					"Api-Key":          api.apiKey,
					"Api-Signature":    signature,
					"Api-Timestamp":    timestamp,
				}).
				SetResult(&GetTradeResponse{}).
				// SetError(&ErrorResp{}).
				Get(api.basePath + ressource)
			if err != nil {
				return tradesResp, errors.New(SOURCE + " Error Requesting")
			}
			if resp.StatusCode() > 300 {
				return tradesResp, errors.New(SOURCE + " Error StatusCode" + strconv.Itoa(resp.StatusCode()))
			}
			trades := *resp.Result().(*GetTradeResponse)
			lastResponseCount = len(trades)
			lastObjectId = trades[lastResponseCount-1].ID
			tradesResp = append(tradesResp, trades...)
			time.Sleep(api.timeBetweenReq)
		}
		if useCache {
			err = db.Write("Bittrex/orders", "closed", tradesResp)
			if err != nil {
				return tradesResp, errors.New(SOURCE + " Error Caching")
			}
		}
	}
	return tradesResp, nil
}
