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
	ID          string
	Description string
	Currency    string
	Amount      decimal.Decimal
	Fee         decimal.Decimal
}

func (api *api) getDepositsTXs(loc *time.Location) {
	today := time.Now()
	thisYear := today.Year()
	for y := thisYear; y > 2017; y-- {
		for t := 6; t > 0; t-- {
			fmt.Print(".")
			depoHist, err := api.getDepositHistory(y, t, loc)
			if err != nil {
				api.doneDep <- err
				return
			}
			for _, dep := range depoHist {
				tx := depositTX{}
				tx.Timestamp = time.Unix(dep.Inserttime/1e3, 0)
				tx.ID = dep.Txid
				tx.Description = "from " + dep.Address
				tx.Currency = dep.Coin
				tx.Amount, _ = decimal.NewFromString(dep.Amount)
				api.depositTXs = append(api.depositTXs, tx)
			}
		}
	}
	api.doneDep <- nil
}

type GetDepositHistoryResp []struct {
	Amount       string `json:"amount"`
	Coin         string `json:"coin"`
	Network      string `json:"network"`
	Status       int    `json:"status"`
	Address      string `json:"address"`
	Addresstag   string `json:"addressTag"`
	Txid         string `json:"txId"`
	Inserttime   int64  `json:"insertTime"`
	Transfertype int    `json:"transferType"`
	Confirmtimes string `json:"confirmTimes"`
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
		err = db.Read("Binance/sapi/v1/capital/deposit/hisrec", period, &depoHist)
	}
	if !useCache || err != nil {
		endpoint := "sapi/v1/capital/deposit/hisrec"
		queryParams := map[string]string{
			"status":     "1",
			"startTime":  fmt.Sprintf("%v", start_ts.UTC().UnixNano()/1e6),
			"endTime":    fmt.Sprintf("%v", end_ts.UTC().UnixNano()/1e6),
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
			err = db.Write("Binance/sapi/v1/capital/deposit/hisrec", period, depoHist)
			if err != nil {
				return depoHist, errors.New("Binance API Deposits : Error Caching" + period)
			}
		}
		time.Sleep(api.timeBetweenRequests)
	}
	return depoHist, nil
}
