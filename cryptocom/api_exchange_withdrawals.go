package cryptocom

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type withdrawalTX struct {
	Timestamp   time.Time
	Description string
	Currency    string
	Amount      decimal.Decimal
	Fee         decimal.Decimal
}

func (api *apiEx) getWithdrawalsTXs(loc *time.Location) {
	today := time.Now()
	thisYear := today.Year()
	for y := thisYear; y > 2019; y-- {
		for q := 4; q > 0; q-- {
			fmt.Print(".")
			withHist, err := api.getWithdrawalHistory(y, q, loc)
			if err != nil {
				api.doneWit <- err
				return
			}
			for _, wit := range withHist.Result.WithdrawalList {
				tx := withdrawalTX{}
				tx.Timestamp = time.Unix(wit.UpdateTime/1000, 0)
				tx.Description = "to " + wit.Address
				tx.Currency = wit.Currency
				tx.Amount = decimal.NewFromFloat(wit.Amount)
				tx.Fee = decimal.NewFromFloat(wit.Fee)
				api.withdrawalTXs = append(api.withdrawalTXs, tx)
			}
		}
	}
	api.doneWit <- nil
}

type ResultWithdrawal struct {
	Currency   string  `json:"currency"`
	ClientWid  string  `json:"client_wid"`
	Fee        float64 `json:"fee"`
	CreateTime int64   `json:"create_time"`
	ID         string  `json:"id"`
	UpdateTime int64   `json:"update_time"`
	Amount     float64 `json:"amount"`
	Address    string  `json:"address"`
	Status     string  `json:"status"`
}

type WithdrawalList struct {
	WithdrawalList []ResultWithdrawal `json:"withdrawal_list"`
}

type GetWithdrawalHistoryResp struct {
	ID     int64          `json:"id"`
	Method string         `json:"method"`
	Code   int            `json:"code"`
	Result WithdrawalList `json:"result"`
}

func (api *apiEx) getWithdrawalHistory(year, quarter int, loc *time.Location) (withHist GetWithdrawalHistoryResp, err error) {
	var start_month time.Month
	var end_month time.Month
	end_year := year
	period := strconv.Itoa(year) + "-Q" + strconv.Itoa(quarter)
	if quarter == 1 {
		start_month = time.January
		end_month = time.April
	} else if quarter == 2 {
		start_month = time.April
		end_month = time.July
	} else if quarter == 3 {
		start_month = time.July
		end_month = time.October
	} else if quarter == 4 {
		start_month = time.October
		end_month = time.January
		end_year = year + 1
	} else {
		err = errors.New("Crypto.com Exchange API Withdrawals : Invalid Quarter" + period)
		return
	}
	start_ts := time.Date(year, start_month, 1, 0, 0, 0, 0, loc)
	end_ts := time.Date(end_year, end_month, 1, 0, 0, 0, 0, loc)
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
		err = db.Read("Crypto.com/Exchange/private/get-withdrawal-history", period, &withHist)
	}
	if !useCache || err != nil {
		method := "private/get-withdrawal-history"
		body := make(map[string]interface{})
		body["method"] = method
		body["params"] = map[string]interface{}{
			"start_ts":  start_ts.UnixNano() / 1e6,
			"end_ts":    end_ts.UnixNano() / 1e6,
			"page_size": 200,
			"page":      0,
			"status":    "5",
		}
		api.sign(body)
		resp, err := api.clientWit.R().
			SetBody(body).
			SetResult(&GetWithdrawalHistoryResp{}).
			SetError(&ErrorResp{}).
			Post(api.basePath + method)
		if err != nil {
			return withHist, errors.New("Crypto.com Exchange API Withdrawals : Error Requesting" + period)
		}
		if resp.StatusCode() > 300 {
			return withHist, errors.New("Crypto.com Exchange API Withdrawals : Error StatusCode" + strconv.Itoa(resp.StatusCode()) + " for " + period)
		}
		withHist = *resp.Result().(*GetWithdrawalHistoryResp)
		if useCache {
			err = db.Write("Crypto.com/Exchange/private/get-withdrawal-history", period, withHist)
			if err != nil {
				return withHist, errors.New("Crypto.com Exchange API Withdrawals : Error Caching" + period)
			}
		}
		time.Sleep(api.timeBetweenReq)
	}
	return withHist, nil
}
