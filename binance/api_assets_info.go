package binance

import (
	"errors"
	"strconv"
	"time"

	scribble "github.com/nanobox-io/golang-scribble"
)

type Symbols struct {
	Symbol                     string   `json:"symbol"`
	Status                     string   `json:"status"`
	Baseasset                  string   `json:"baseAsset"`
	Baseassetprecision         int      `json:"baseAssetPrecision"`
	Quoteasset                 string   `json:"quoteAsset"`
	Quoteprecision             int      `json:"quotePrecision"`
	Quoteassetprecision        int      `json:"quoteAssetPrecision"`
	Basecommissionprecision    int      `json:"baseCommissionPrecision"`
	Quotecommissionprecision   int      `json:"quoteCommissionPrecision"`
	Ordertypes                 []string `json:"orderTypes"`
	Icebergallowed             bool     `json:"icebergAllowed"`
	Ocoallowed                 bool     `json:"ocoAllowed"`
	Quoteorderqtymarketallowed bool     `json:"quoteOrderQtyMarketAllowed"`
	Isspottradingallowed       bool     `json:"isSpotTradingAllowed"`
	Ismargintradingallowed     bool     `json:"isMarginTradingAllowed"`
	Filters                    []struct {
		Filtertype       string `json:"filterType"`
		Minprice         string `json:"minPrice,omitempty"`
		Maxprice         string `json:"maxPrice,omitempty"`
		Ticksize         string `json:"tickSize,omitempty"`
		Multiplierup     string `json:"multiplierUp,omitempty"`
		Multiplierdown   string `json:"multiplierDown,omitempty"`
		Avgpricemins     int    `json:"avgPriceMins,omitempty"`
		Minqty           string `json:"minQty,omitempty"`
		Maxqty           string `json:"maxQty,omitempty"`
		Stepsize         string `json:"stepSize,omitempty"`
		Minnotional      string `json:"minNotional,omitempty"`
		Applytomarket    bool   `json:"applyToMarket,omitempty"`
		Limit            int    `json:"limit,omitempty"`
		Maxnumorders     int    `json:"maxNumOrders,omitempty"`
		Maxnumalgoorders int    `json:"maxNumAlgoOrders,omitempty"`
	} `json:"filters"`
	Permissions []string `json:"permissions"`
}

type Ratelimits struct {
	Ratelimittype string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	Intervalnum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

type GetExchangeInfoResp struct {
	Timezone        string        `json:"timezone"`
	Servertime      int64         `json:"serverTime"`
	Ratelimits      []Ratelimits  `json:"rateLimits"`
	Exchangefilters []interface{} `json:"exchangeFilters"`
	Symbols         []Symbols     `json:"symbols"`
}

func (api *api) getExchangeInfo() (exchangeInfo GetExchangeInfoResp, err error) {
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Binance/api/v3/", "exchangeInfo", &exchangeInfo)
	}
	if !useCache || err != nil {
		method := "api/v3/exchangeInfo"
		resp, err := api.clientDep.R().
			SetResult(&GetExchangeInfoResp{}).
			SetError(&ErrorResp{}).
			Get(api.basePath + method)
		if err != nil {
			return exchangeInfo, errors.New("Binance API Deposits : Error Requesting exchangeInfo")
		}
		if resp.StatusCode() > 300 {
			return exchangeInfo, errors.New("Binance API Deposits : Error StatusCode" + strconv.Itoa(resp.StatusCode()))
		}
		exchangeInfo = *resp.Result().(*GetExchangeInfoResp)
		if useCache {
			err = db.Write("Binance/api/v3/", "exchangeInfo", exchangeInfo)
			if err != nil {
				return exchangeInfo, errors.New("Binance API Deposits : Error Caching exchangeInfo")
			}
		}
		time.Sleep(api.timeBetweenReq)
	}
	return exchangeInfo, nil
}
