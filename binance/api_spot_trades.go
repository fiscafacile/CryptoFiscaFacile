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
	BaseAsset   string
	QuoteAsset  string
	Side        string
	Price       decimal.Decimal
	Qty         decimal.Decimal
	QuoteQty    decimal.Decimal
	Fee         decimal.Decimal
	FeeCurrency string
	ID          string
}

func (api *api) getSpotTradesTXs() {
	limit := 10
	totalSymbols := len(api.symbols)
	for i, symbol := range api.symbols {
		orderId := 0
		part := 0
		if api.debug {
			fmt.Printf("[%v/%v] Récupération des trades du symbole %v\n", i+1, totalSymbols, symbol.Symbol)
		}
		for {
			trades, err := api.getTrades(symbol.Symbol, limit, orderId+1, part)
			if err != nil {
				api.doneSpotTra <- err
				return
			}
			for _, tra := range trades {
				tx := spotTradeTX{}
				tx.Timestamp = time.Unix(tra.Time, 0)
				tx.ID = strconv.Itoa(tra.ID)
				tx.Description = fmt.Sprintf("Order ID: %v, Trade ID: %v", tra.Orderid, tra.ID)
				tx.BaseAsset = symbol.Baseasset
				tx.QuoteAsset = symbol.Quoteasset
				if tra.Isbuyer {
					tx.Side = "BUY"
				} else {
					tx.Side = "SELL"
				}
				tx.Price, _ = decimal.NewFromString(tra.Price)
				tx.Qty, _ = decimal.NewFromString(tra.Qty)
				tx.QuoteQty, _ = decimal.NewFromString(tra.Quoteqty)
				tx.Fee, _ = decimal.NewFromString(tra.Commission)
				tx.FeeCurrency = tra.Commissionasset
				api.spotTradeTXs = append(api.spotTradeTXs, tx)
			}
			if len(trades) == limit {
				part += 1
				orderId = trades[len(trades)-1].ID
			} else {
				break
			}
		}
	}
	api.doneSpotTra <- nil
}

type GetTradesResp []struct {
	Symbol          string `json:"symbol"`
	ID              int    `json:"id"`
	Orderid         int    `json:"orderId"`
	Orderlistid     int    `json:"orderListId"`
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	Quoteqty        string `json:"quoteQty"`
	Commission      string `json:"commission"`
	Commissionasset string `json:"commissionAsset"`
	Time            int64  `json:"time"`
	Isbuyer         bool   `json:"isBuyer"`
	Ismaker         bool   `json:"isMaker"`
	Isbestmatch     bool   `json:"isBestMatch"`
}

func (api *api) getTrades(symbol string, limit int, orderId int, part int) (trades GetTradesResp, err error) {
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Binance/api/v3/myTrades", fmt.Sprintf("%v_%v-%v", symbol, part*limit, part*limit+limit), &trades)
		if err == nil && len(trades) != limit { // If cached data is incomplete
			useCache = false
		}
	}
	if !useCache || err != nil {
		endpoint := "api/v3/myTrades"
		queryParams := map[string]string{
			"symbol":     symbol,
			"fromId":     fmt.Sprintf("%v", orderId),
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
			err = db.Write("Binance/api/v3/myTrades/", fmt.Sprintf("%v_%v-%v", symbol, part*limit, part*limit+limit), trades)
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
				fmt.Println(usedWeight, "/", api.reqWeightlimit, "Weight utilisé --> pause pendant", api.timeBetweenRequests)
			}
			time.Sleep(api.timeBetweenRequests)
		}
	}
	return trades, nil
}
