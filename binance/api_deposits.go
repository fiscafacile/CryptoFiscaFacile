package binance

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type depositTX struct {
	Timestamp   time.Time
	Description string
	Currency    string
	Amount      decimal.Decimal
	Fee         decimal.Decimal
}

func (api *api) getDepositsTXs(loc *time.Location) {
	today := time.Now()
	thisYear := today.Year()
	for y := thisYear; y > 2017; y-- {
		for q := 4; q > 0; q-- {
			depoHist, err := api.getDepositHistory(y, q, loc)
			if err != nil {
				api.doneDep <- err
				return
			}
			for _, dep := range depoHist.Depositlist {
				tx := depositTX{}
				tx.Timestamp = time.Unix(dep.Inserttime, 0)
				tx.Description = "from " + dep.Address
				tx.Currency = dep.Asset
				tx.Amount = decimal.NewFromFloat(dep.Amount)
				api.depositTXs = append(api.depositTXs, tx)
			}
		}
	}
	api.doneDep <- nil
}

type Depositlist struct {
	Inserttime int64   `json:"insertTime"`
	Amount     float64 `json:"amount"`
	Asset      string  `json:"asset"`
	Address    string  `json:"address"`
	Txid       string  `json:"txId"`
	Status     int     `json:"status"`
	Addresstag string  `json:"addressTag,omitempty"`
}
type GetDepositHistoryResp struct {
	Depositlist []Depositlist
	Success     bool `json:"success"`
}

func (api *api) getDepositHistory(year, trimester int, loc *time.Location) (depoHist GetDepositHistoryResp, err error) {
	var start_month time.Month
	var end_month time.Month
	end_year := year
	period := strconv.Itoa(year) + "-T" + strconv.Itoa(trimester)
	if trimester == 1 {
		start_month = time.January
		end_month = time.March
	} else if trimester == 2 {
		start_month = time.April
		end_month = time.June
	} else if trimester == 3 {
		start_month = time.July
		end_month = time.September
	} else if trimester == 4 {
		start_month = time.October
		end_month = time.December
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
		err = db.Read("Binance/wapi/v3/depositHistory", period, &depoHist)
	}
	if !useCache || err != nil {
		endpoint := "wapi/v3/depositHistory.html"
		queryParams := map[string]string{
			"status":     "1",
			"startTime":  strconv.FormatInt(start_ts.Unix(), 10),
			"endTime":    strconv.FormatInt(end_ts.Unix(), 10),
			"recvWindow": "60000",
			"timestamp":  fmt.Sprintf("%v", time.Now().UTC().UnixNano()/1e6),
		}
		api.sign(queryParams)
		resp, err := api.clientDep.R().
			SetHeader("X-MBX-APIKEY", api.apiKey).
			SetQueryParams(queryParams).
			SetResult(&GetDepositHistoryResp{}).
			SetError(&ErrorResp{}).
			Get(api.basePath + endpoint)
		if err != nil {
			return depoHist, errors.New("Binance API Deposits : Error Requesting" + period)
		}
		if resp.StatusCode() > 300 {
			return depoHist, errors.New("Binance API Deposits : Error StatusCode" + strconv.Itoa(resp.StatusCode()) + " for " + period)
		}
		depoHist = *resp.Result().(*GetDepositHistoryResp)
		if useCache {
			err = db.Write("Binance/wapi/v3/depositHistory", period, depoHist)
			if err != nil {
				return depoHist, errors.New("Binance API Deposits : Error Caching" + period)
			}
		}
		time.Sleep(api.timeBetweenRequests)
	}
	return depoHist, nil
}
