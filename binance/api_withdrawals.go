package binance

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type withdrawalTX struct {
	Timestamp   time.Time
	ID          string
	Description string
	Currency    string
	Amount      decimal.Decimal
	Fee         decimal.Decimal
}

func (api *api) getWithdrawalsTXs(loc *time.Location) {
	today := time.Now()
	thisYear := today.Year()
	for y := thisYear; y > 2017; y-- {
		for t := 6; t > 0; t-- {
			fmt.Print(".")
			withHist, err := api.getWithdrawalHistory(y, t, loc)
			if err != nil {
				// api.doneWit <- err
				// return
				log.Println(err)
				break
			}
			for _, wit := range withHist {
				tx := withdrawalTX{}
				tx.Timestamp, err = time.Parse("2006-01-02 15:04:05", wit.ApplyTime)
				if err != nil {
					log.Println("Error Parsing Time : ", wit.ApplyTime)
				}
				tx.ID = wit.TxID
				tx.Description = "to " + wit.Address
				tx.Currency = wit.Coin
				tx.Amount, _ = decimal.NewFromString(wit.Amount)
				api.withdrawalTXs = append(api.withdrawalTXs, tx)
			}
		}
	}
	api.doneWit <- nil
}

type GetWithdrawalHistoryResp []struct {
	ApplyTime       string `json:"applyTime"`
	Amount          string `json:"amount"`
	Coin            string `json:"coin"`
	Network         string `json:"network"`
	Status          int    `json:"status"`
	Address         string `json:"address"`
	ID              string `json:"id"`
	TxID            string `json:"txId"`
	Transfertype    int    `json:"transferType"`
	WithdrawOrderId string `json:"withdrawOrderId"`
}

func (api *api) getWithdrawalHistory(year, trimester int, loc *time.Location) (withHist GetWithdrawalHistoryResp, err error) {
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
		err = errors.New("Binance API Withdrawals : Invalid trimester" + period)
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
		// period += "-" + strconv.FormatInt(end_ts.Unix(), 10)
	}
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Binance/sapi/v1/capital/withdraw/history", period, &withHist)
	}
	if !useCache || err != nil {
		endpoint := "sapi/v1/capital/withdraw/history"
		queryParams := map[string]string{
			"status":     "6",
			"startTime":  fmt.Sprintf("%v", start_ts.UTC().UnixNano()/1e6),
			"endTime":    fmt.Sprintf("%v", end_ts.UTC().UnixNano()/1e6),
			"recvWindow": "60000",
			"timestamp":  fmt.Sprintf("%v", time.Now().UTC().UnixNano()/1e6),
		}
		api.sign(queryParams)
		resp, err := api.clientWit.R().
			SetHeader("X-MBX-APIKEY", api.apiKey).
			SetQueryParams(queryParams).
			SetResult(&GetWithdrawalHistoryResp{}).
			SetError(&ErrorResp{}).
			Get(api.basePath + endpoint)
		if err != nil {
			return withHist, errors.New("Binance API Withdrawals : Error Requesting" + period)
		}
		if resp.StatusCode() > 300 {
			return withHist, errors.New("Binance API Withdrawals : Error StatusCode" + strconv.Itoa(resp.StatusCode()) + " for " + period)
		}
		withHist = *resp.Result().(*GetWithdrawalHistoryResp)
		if useCache {
			err = db.Write("Binance/sapi/v1/capital/withdraw/history", period, withHist)
			if err != nil {
				return withHist, errors.New("Binance API Withdrawals : Error Caching" + period)
			}
		}
		time.Sleep(api.timeBetweenRequests)
	}
	return withHist, nil
}
