package binance

import (
	"errors"
	"strconv"
	"time"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type exchangeInfo struct {
	Timestamp   time.Time
	Description string
	Currency    string
	Amount      decimal.Decimal
	Fee         decimal.Decimal
}

type ResultExchange struct {
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

type InfoList struct {
	DepositList []ResultExchange `json:"deposit_list"`
}

type GetExchangeInfoResp struct {
	ID     int64    `json:"id"`
	Method string   `json:"method"`
	Code   int      `json:"code"`
	Result InfoList `json:"result"`
}

func (api *api) getExchangeInfo() (exchangeInfo GetDepositHistoryResp, err error) {
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
		body := make(map[string]interface{})
		body["method"] = method
		api.sign(body)
		resp, err := api.clientDep.R().
			SetBody(body).
			SetResult(&GetDepositHistoryResp{}).
			SetError(&ErrorResp{}).
			Post(api.basePath + method)
		if err != nil {
			return exchangeInfo, errors.New("Binance API Deposits : Error Requesting exchangeInfo")
		}
		if resp.StatusCode() > 300 {
			return exchangeInfo, errors.New("Binance API Deposits : Error StatusCode" + strconv.Itoa(resp.StatusCode()))
		}
		exchangeInfo = *resp.Result().(*GetDepositHistoryResp)
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
