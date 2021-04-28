package wallet

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
	"gopkg.in/resty.v1"
)

type Rate struct {
	Time  time.Time       `json:"time"`
	Quote string          `json:"asset_id_quote"`
	Rate  decimal.Decimal `json:"rate"`
}

type ExchangeRates struct {
	Base  string `json:"asset_id_base"`
	Rates []Rate `json:"rates"`
}

type CoinAPI struct {
}

func CoinAPISetKey(key string) error {
	return os.Setenv("COINAPI_KEY", key)
}

func (api CoinAPI) GetExchangeRates(date time.Time, native string) (rates ExchangeRates, err error) {
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		return
	}
	err = db.Read("CoinAPI/exchangerate", native+"-"+date.UTC().Format("2006-01-02-15-04-05"), &rates)
	if err != nil {
		if os.Getenv("COINAPI_KEY") == "" {
			return rates, errors.New("Need CoinAPI Key")
		}
		hour := 0
		for len(rates.Rates) == 0 && hour < 15 {
			url := "http://rest.coinapi.io/v1/exchangerate/" + native + "?invert=true&time=" + date.Add(time.Duration(hour)*time.Hour).UTC().Format(time.RFC3339)
			resp, err := resty.R().SetHeaders(map[string]string{
				"Accept":        "application/json",
				"X-CoinAPI-Key": os.Getenv("COINAPI_KEY"),
			}).Get(url)
			if err != nil {
				return rates, err
			}
			if resp.StatusCode() != http.StatusOK {
				err = errors.New("Error Status : " + strconv.Itoa(resp.StatusCode()))
				return rates, err
			}
			err = json.Unmarshal(resp.Body(), &rates)
			if err != nil {
				return rates, err
			}
			if len(rates.Rates) == 0 {
				hour += 1
				if hour == 15 {
					log.Println("CoinAPI Get void Rates:", url)
				}
			}
		}
		if len(rates.Rates) != 0 {
			err = db.Write("CoinAPI/exchangerate", native+"-"+date.UTC().Format("2006-01-02-15-04-05"), rates)
		}
		return rates, err
	}
	return
}
