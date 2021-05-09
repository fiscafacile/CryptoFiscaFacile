package binance

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type assetDividendTX struct {
	Timestamp   time.Time
	Description string
	Asset       string
	Price       decimal.Decimal
	Amount      decimal.Decimal
	ID          string
}

func (api *api) getAssetDividendTXs(loc *time.Location) {
	today := time.Now()
	thisYear := today.Year()
	for y := thisYear; y > 2017; y-- {
		for t := 6; t > 0; t-- {
			assDiv, err := api.getAssetDividend(y, t, loc)
			if err != nil {
				api.doneSpotTra <- err
				return
			}
			for _, div := range assDiv.Rows {
				tx := assetDividendTX{}
				tx.Timestamp = time.Unix(div.Divtime/1e3, 0)
				tx.ID = strconv.FormatInt(div.Tranid, 10)
				tx.Description = div.Eninfo
				tx.Amount, _ = decimal.NewFromString(div.Amount)
				tx.Asset = div.Asset
				api.assetDividendTXs = append(api.assetDividendTXs, tx)
			}
		}
	}
	api.doneAssDiv <- nil
}

type Rows []struct {
	Amount  string `json:"amount"`
	Asset   string `json:"asset"`
	Divtime int64  `json:"divTime"`
	Eninfo  string `json:"enInfo"`
	Tranid  int64  `json:"tranId"`
}

type GetAssetDividendResp struct {
	Rows  Rows `json:"rows"`
	Total int  `json:"total"`
}

func (api *api) getAssetDividend(year, trimester int, loc *time.Location) (assDiv GetAssetDividendResp, err error) {
	var start_month time.Month
	var end_month time.Month
	end_year := year
	period := strconv.Itoa(year) + "-T" + strconv.Itoa(trimester)
	if trimester == 1 {
		start_month = time.January
		end_month = time.March
	} else if trimester == 2 {
		start_month = time.March
		end_month = time.May
	} else if trimester == 3 {
		start_month = time.May
		end_month = time.July
	} else if trimester == 4 {
		start_month = time.July
		end_month = time.September
	} else if trimester == 5 {
		start_month = time.September
		end_month = time.November
	} else if trimester == 6 {
		start_month = time.November
		end_month = time.January
		end_year = year + 1
	} else {
		err = errors.New("Binance API Deposits : Invalid trimester" + period)
		return
	}
	start_ts := time.Date(year, start_month, 1, 0, 0, 0, 0, loc)
	end_ts := time.Date(end_year, end_month, 1, 0, 0, 0, 0, loc)
	now := time.Now()
	if start_ts.After(now) {
		return
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
		err = db.Read("Binance/sapi/v1/asset/assetDividend/", period, &assDiv)
	}
	if !useCache || err != nil {
		endpoint := "sapi/v1/asset/assetDividend"
		queryParams := map[string]string{
			"limit":      "500",
			"startTime":  fmt.Sprintf("%v", start_ts.UTC().UnixNano()/1e6),
			"endTime":    fmt.Sprintf("%v", end_ts.UTC().UnixNano()/1e6),
			"recvWindow": "60000",
			"timestamp":  fmt.Sprintf("%v", time.Now().UTC().UnixNano()/1e6),
		}
		api.sign(queryParams)
		resp, err := api.clientAssDiv.R().
			SetHeader("X-MBX-APIKEY", api.apiKey).
			SetQueryParams(queryParams).
			SetResult(&GetAssetDividendResp{}).
			SetError(&ErrorResp{}).
			Get(api.basePath + endpoint)
		if err != nil {
			return assDiv, errors.New("Binance API Asset Dividend : Error Requesting " + period)
		}
		if resp.StatusCode() > 300 {
			return assDiv, errors.New("Binance API Asset Dividend : Error StatusCode " + strconv.Itoa(resp.StatusCode()) + " for " + period)
		}
		assDiv = *resp.Result().(*GetAssetDividendResp)
		if useCache {
			err = db.Write("Binance/sapi/v1/asset/assetDividend/", period, assDiv)
			if err != nil {
				return assDiv, errors.New("Binance API Asset Dividend : Error Caching" + period)
			}
		}
		time.Sleep(api.timeBetweenRequests)
	}
	return assDiv, nil
}
