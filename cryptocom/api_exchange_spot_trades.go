package cryptocom

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/nanobox-io/golang-scribble"
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

func (api *apiEx) getSpotTradesTXs(loc *time.Location) {
	date := time.Now().Add(-24 * time.Hour)
	for date.After(api.startTime) {
		fmt.Print(".")
		trades, err := api.getTrades(date.Year(), date.Month(), date.Day(), loc)
		if err != nil {
			api.doneSpotTra <- err
			return
		}
		for _, tra := range trades.Result.TradeList {
			tx := spotTradeTX{}
			tx.Timestamp = time.Unix(tra.CreateTime/1000, 0)
			tx.Description = tra.TradeID + " " + tra.LiquidityIndicator
			tx.Pair = tra.InstrumentName
			tx.Side = tra.Side
			tx.Price = decimal.NewFromFloat(tra.TradedPrice)
			tx.Quantity = decimal.NewFromFloat(tra.TradedQuantity)
			tx.Fee = decimal.NewFromFloat(tra.Fee)
			tx.FeeCurrency = tra.FeeCurrency
			api.spotTradeTXs = append(api.spotTradeTXs, tx)
		}
		date = date.Add(-24 * time.Hour)
	}
	api.doneSpotTra <- nil
}

type ResultTrade struct {
	Side               string  `json:"side"`
	InstrumentName     string  `json:"instrument_name"`
	Fee                float64 `json:"fee"`
	TradeID            string  `json:"trade_id"`
	CreateTime         int64   `json:"create_time"`
	TradedPrice        float64 `json:"traded_price"`
	TradedQuantity     float64 `json:"traded_quantity"`
	LiquidityIndicator string  `json:"liquidity_indicator"`
	FeeCurrency        string  `json:"fee_currency"`
	OrderID            string  `json:"order_id"`
}

type TradeList struct {
	TradeList []ResultTrade `json:"trade_list"`
}

type GetTradesResp struct {
	ID     int64     `json:"id"`
	Method string    `json:"method"`
	Code   int       `json:"code"`
	Result TradeList `json:"result"`
}

func (api *apiEx) getTrades(year int, month time.Month, day int, loc *time.Location) (trades GetTradesResp, err error) {
	start_ts := time.Date(year, month, day, 0, 0, 0, 0, loc)
	end_ts := start_ts.Add(24 * time.Hour).Add(-time.Millisecond)
	period := start_ts.Format("2006-01-02")
	now := time.Now()
	if start_ts.After(now) {
		return // without error
	}
	if end_ts.After(now) {
		end_ts = now
		period += "-" + strconv.FormatInt(end_ts.Unix(), 10)
	}
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Crypto.com/Exchange/private/get-trades", period, &trades)
	}
	if !useCache || err != nil {
		method := "private/get-trades"
		body := make(map[string]interface{})
		body["method"] = method
		body["params"] = map[string]interface{}{
			"start_ts":  start_ts.UnixNano() / 1e6,
			"end_ts":    end_ts.UnixNano() / 1e6,
			"page_size": 200,
			"page":      0,
		}
		api.sign(body)
		resp, err := api.clientSpotTra.R().
			SetBody(body).
			SetResult(&GetTradesResp{}).
			SetError(&ErrorResp{}).
			Post(api.basePath + method)
		if err != nil {
			return trades, errors.New("Crypto.com Exchange API Trades : Error Requesting" + period)
		}
		if resp.StatusCode() > 300 {
			return trades, errors.New("Crypto.com Exchange API Trades : Error StatusCode" + strconv.Itoa(resp.StatusCode()) + " for " + period)
		}
		trades = *resp.Result().(*GetTradesResp)
		if useCache {
			err = db.Write("Crypto.com/Exchange/private/get-trades", period, trades)
			if err != nil {
				return trades, errors.New("Crypto.com Exchange API Trades : Error Caching" + period)
			}
		}
		time.Sleep(api.timeBetweenReqSpot)
	}
	return trades, nil
}
