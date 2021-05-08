package binance

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type spotTradeTX struct {
	Timestamp   time.Time
	Description string
	Pair        string
	Side        string
	Price       decimal.Decimal
	Quantity    decimal.Decimal
	Fee         decimal.Decimal
	FeeCurrency string
}

func (api *api) getSpotTradesTXs() {
	limit := 50
	for _, symbol := range api.symbols {
		orderId := 0
		part := 0
		for {
			fmt.Println("Récupération des trades du symbole", symbol.Symbol)
			trades, err := api.getTrades(symbol.Symbol, limit, orderId+1, part)
			if err != nil {
				api.doneSpotTra <- err
				return
			}
			for _, tra := range trades {
				tx := spotTradeTX{}
				tx.Timestamp = time.Unix(tra.Time, 0)
				tx.Description = fmt.Sprintf("%v", tra.Orderid)
				tx.Pair = tra.Symbol
				tx.Side = tra.Side
				tx.Price, _ = decimal.NewFromString(tra.Price)
				tx.Quantity, _ = decimal.NewFromString(tra.Executedqty)
				// tx.Fee = decimal.NewFromFloat(tra.Fee)
				// tx.FeeCurrency = tra.FeeCurrency
				api.spotTradeTXs = append(api.spotTradeTXs, tx)
			}
			if len(trades) == limit {
				part += 1
				orderId = trades[len(trades)-1].Orderid
			} else {
				break
			}
		}
	}
	api.doneSpotTra <- nil
}

type GetTradesResp []struct {
	Symbol              string `json:"symbol"`
	Orderid             int    `json:"orderId"`
	Orderlistid         int    `json:"orderListId"`
	Clientorderid       string `json:"clientOrderId"`
	Price               string `json:"price"`
	Origqty             string `json:"origQty"`
	Executedqty         string `json:"executedQty"`
	Cummulativequoteqty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	Timeinforce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
	Stopprice           string `json:"stopPrice"`
	Icebergqty          string `json:"icebergQty"`
	Time                int64  `json:"time"`
	Updatetime          int64  `json:"updateTime"`
	Isworking           bool   `json:"isWorking"`
	Origquoteorderqty   string `json:"origQuoteOrderQty"`
}

func (api *api) getTrades(symbol string, limit int, orderId int, part int) (trades GetTradesResp, err error) {
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Binance/api/v3/allOrders", fmt.Sprintf("%v_%v-%v", symbol, part*limit, part*limit+limit), &trades)
		if err == nil && len(trades) != limit { // If cached data is incomplete
			useCache = false
		}
	}
	if !useCache || err != nil {
		endpoint := "api/v3/allOrders"
		queryParams := map[string]string{
			"symbol":     symbol,
			"orderId":    fmt.Sprintf("%v", orderId),
			"limit":      fmt.Sprintf("%v", limit),
			"recvWindow": "60000",
			"timestamp":  fmt.Sprintf("%v", time.Now().UTC().UnixNano()/1e6),
		}
		api.sign(queryParams)
		resp, err := api.clientSpotTra.R().
			SetHeader("X-MBX-APIKEY", api.apiKey).
			SetQueryParams(queryParams).
			SetResult(&GetTradesResp{}).
			SetError(&ErrorResp{}).
			Get(api.basePath + endpoint)
		if err != nil {
			return trades, errors.New("Binance API Trades : Error Requesting " + symbol)
		}
		if resp.StatusCode() > 300 {
			return trades, errors.New("Binance API Trades : Error StatusCode " + strconv.Itoa(resp.StatusCode()) + " for " + symbol)
		}
		trades = *resp.Result().(*GetTradesResp)
		if useCache {
			err = db.Write("Binance/api/v3/allOrders/", fmt.Sprintf("%v_%v-%v", symbol, part*limit, part*limit+limit), trades)
			if err != nil {
				return trades, errors.New("Binance API Trades : Error Caching" + fmt.Sprintf("%v_%v-%v", symbol, part*limit, part*limit+limit))
			}
		}
		weightHeader := fmt.Sprintf("X-MBX-USED-WEIGHT-%v%v", api.reqWeightIntervalNum, string(api.reqWeightInterval[0]))
		usedWeight, _ := strconv.Atoi(resp.Header().Get(weightHeader))
		if usedWeight >= api.reqWeightlimit-20 {
			if api.debug {
				fmt.Println(usedWeight, "/", api.reqWeightlimit, "Weight utilisé --> pause pendant", api.reqWeightTimeToWait)
			}
			time.Sleep(api.reqWeightTimeToWait)
		} else {
			if api.debug {
				fmt.Println(usedWeight, "/", api.reqWeightlimit, "Weight utilisé --> pause pendant", api.timeBetweenReqOrder)
			}
			time.Sleep(api.timeBetweenReqOrder)
		}
	}
	return trades, nil
}
