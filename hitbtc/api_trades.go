package hitbtc

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type tradeTX struct {
	ID            int             // 9535486,
	OrderID       int             // 816088377,
	ClientOrderID string          // "f8dbaab336d44d5ba3ff578098a68454",
	Symbol        string          // "ETHBTC",
	Side          string          // "sell",
	Quantity      decimal.Decimal // "0.061",
	Price         decimal.Decimal // "0.045487",
	Fee           decimal.Decimal // "0.000002775",
	Timestamp     time.Time       // "2017-05-17T12:32:57.848Z"
}

func (api *api) getTradesTXs() {
	const SOURCE = "HitBTC API Trades :"
	trades, err := api.getTrades()
	if err != nil {
		api.doneTrade <- err
		return
	}
	for _, tra := range trades {
		tx := tradeTX{}
		tx.ID = tra.ID
		tx.OrderID = tra.OrderID
		tx.ClientOrderID = tra.ClientOrderID
		tx.Symbol = tra.Symbol
		tx.Side = tra.Side
		tx.Quantity, err = decimal.NewFromString(tra.Quantity)
		if err != nil {
			log.Println(SOURCE, "Error Parsing Quantity : ", tra.Quantity)
		}
		tx.Price, err = decimal.NewFromString(tra.Price)
		if err != nil {
			log.Println(SOURCE, "Error Parsing Price : ", tra.Price)
		}
		tx.Fee, err = decimal.NewFromString(tra.Fee)
		if err != nil {
			log.Println(SOURCE, "Error Parsing Fee : ", tra.Fee)
		}
		tx.Timestamp = tra.Timestamp
		api.tradeTXs = append(api.tradeTXs, tx)
	}
	api.doneTrade <- nil
}

type GetTradesResp []struct {
	ID            int       `json:"id"`
	OrderID       int       `json:"orderId"`
	ClientOrderID string    `json:"clientOrderId"`
	Symbol        string    `json:"symbol"`
	Side          string    `json:"side"`
	Quantity      string    `json:"quantity"`
	Price         string    `json:"price"`
	Fee           string    `json:"fee"`
	Timestamp     time.Time `json:"timestamp"`
}

func (api *api) getTrades() (trades GetTradesResp, err error) {
	const SOURCE = "HitBTC API Trades :"
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("HitBTC/history", "trades", &trades)
	}
	if !useCache || err != nil {
		method := "history/trades"
		resp, err := api.clientTrade.R().
			SetBasicAuth(api.apiKey, api.secretKey).
			SetResult(&GetTradesResp{}).
			SetError(&ErrorResp{}).
			Get(api.basePath + method)
		if err != nil {
			return trades, errors.New(SOURCE + " Error Requesting")
		}
		if resp.StatusCode() > 300 {
			return trades, errors.New(SOURCE + " Error StatusCode" + strconv.Itoa(resp.StatusCode()))
		}
		trades = *resp.Result().(*GetTradesResp)
		if useCache {
			err = db.Write("HitBTC/history", "trades", trades)
			if err != nil {
				return trades, errors.New(SOURCE + " Error Caching")
			}
		}
		time.Sleep(api.timeBetweenReq)
	}
	return trades, nil
}
