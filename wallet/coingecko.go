package wallet

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
	"github.com/superoo7/go-gecko/v3"
)

type CoinGeckoAPI struct {
	httpClient *http.Client
	client     *coingecko.Client
}

type CoinsListItem struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type CoinList []CoinsListItem

func NewCoinGeckoAPI() (*CoinGeckoAPI, error) {
	cg := &CoinGeckoAPI{}
	cg.httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
	cg.client = coingecko.NewClient(cg.httpClient)
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		return nil, err
	}
	var coinsList CoinList
	err = db.Read("CoinGecko/coins", "list", &coinsList)
	if err != nil {
		cgCoinsList, err := cg.client.CoinsList()
		if err != nil {
			return nil, err
		}
		err = db.Write("CoinGecko/coins", "list", cgCoinsList)
		return cg, err
	}
	return cg, err
}

func (api CoinGeckoAPI) GetExchangeRates(date time.Time, coin string) (rates ExchangeRates, err error) {
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		return rates, err
	}
	coin = strings.ReplaceAll(coin, "DSH", "dash")
	coin = strings.ReplaceAll(coin, "IOT", "miota")
	coin = strings.ReplaceAll(coin, "MEET.ONE", "meetone")
	err = db.Read("CoinGecko/coins/history", coin+"-"+date.UTC().Format("2006-01-02"), &rates)
	if err != nil {
		err = nil
		coinID := ""
		var coinsList CoinList
		err = db.Read("CoinGecko/coins", "list", &coinsList)
		if err != nil {
			return rates, err
		}
		for _, c := range coinsList {
			if c.Symbol == strings.ToLower(coin) {
				coinID = c.ID
				break
			}
		}
		if coinID != "" {
			hist, err := api.client.CoinsIDHistory(coinID, date.UTC().Format("02-01-2006"), false)
			if err != nil {
				return rates, err
			}
			rates.Base = coin
			if hist.MarketData == nil {
				err = errors.New("CoinGecko API replied a null MarketData")
			} else {
				for k, v := range hist.MarketData.CurrentPrice {
					r := Rate{Time: date, Quote: strings.ToUpper(k), Rate: decimal.NewFromFloat(v)}
					rates.Rates = append(rates.Rates, r)
				}
				db.Write("CoinGecko/coins/history", coin+"-"+date.UTC().Format("2006-01-02"), rates)
				if err != nil {
					return rates, err
				}
			}
		}
	}
	return rates, nil
}
