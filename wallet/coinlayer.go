package wallet

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/nanobox-io/golang-scribble"
	// "github.com/shopspring/decimal"
	"gopkg.in/resty.v1"
)

type HistoricalData struct {
	Success    bool               `json:"success"`
	Terms      string             `json:"terms"`
	Privacy    string             `json:"privacy"`
	Timestamp  int                `json:"timestamp"`
	Target     string             `json:"target"`
	Historical bool               `json:"historical"`
	Date       string             `json:"date"`
	Rates      map[string]float64 `json:"rates"`
}

type CoinLayer struct {
}

func CoinLayerSetKey(key string) error {
	return os.Setenv("COINLAYER_KEY", key)
}

func (api CoinLayer) GetExchangeRates(date time.Time, native string) (rates HistoricalData, err error) {
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		return
	}
	err = db.Read("CoinLayer", native+"-"+date.UTC().Format("2006-01-02"), &rates)
	if err != nil {
		if os.Getenv("COINLAYER_KEY") == "" {
			return rates, errors.New("Need CoinLayer Key")
		}
		url := "http://api.coinlayer.com/" + date.UTC().Format("2006-01-02") + "?access_key=" + os.Getenv("COINLAYER_KEY") + "&target=" + native
		resp, err := resty.R().SetHeaders(map[string]string{
			"Accept": "application/json",
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
		err = db.Write("CoinLayer", native+"-"+date.UTC().Format("2006-01-02"), rates)
		return rates, err
	}
	return
}
