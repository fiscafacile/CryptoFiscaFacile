package kraken

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type ledgerTX struct {
	Amount      decimal.Decimal
	Asset       string
	Class       string
	Fee         decimal.Decimal
	FeeCurrency string
	RefId       string
	SubType     string
	Time        time.Time
	TxId        string
	Type        string
}

func (api *api) getAPISpotTrades() {
	trades, err := api.getTrades()
	if err != nil {
		api.doneLedgers <- err
		return
	}
	assets := api.assets.Result.(map[string]interface{})
	for txID, txData := range trades {
		tra := txData.(map[string]interface{})
		tx := ledgerTX{}
		tx.Amount, err = decimal.NewFromString(tra["amount"].(string))
		if err != nil {
			fmt.Println("Error while parsing amount", tra["amount"].(string))
		}
		if val, ok := assets[tra["asset"].(string)]; ok {
			tx.Asset = ReplaceAssets(val.(map[string]interface{})["altname"].(string))
		} else {
			tx.Asset = ReplaceAssets(tra["asset"].(string))
		}
		tx.Class = tra["aclass"].(string)
		tx.Fee, err = decimal.NewFromString(tra["fee"].(string))
		if err != nil {
			fmt.Println("Error while parsing fee", tra["fee"].(string))
		}
		tx.RefId = tra["refid"].(string)
		tx.SubType = tra["subtype"].(string)
		sec, dec := math.Modf(tra["time"].(float64))
		tx.Time = time.Unix(int64(sec), int64(dec*(1e9)))
		tx.TxId = txID
		tx.Type = tra["type"].(string)
		api.ledgerTX = append(api.ledgerTX, tx)
	}
	api.doneLedgers <- nil
}

type TradesHistory struct {
	Error  []string    `json:"error"`
	Result interface{} `json:"result"`
}

func (api *api) getTrades() (fullTradeTx map[string]interface{}, err error) {
	fullTradeTx = make(map[string]interface{})
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Kraken/private", "Ledgers", &fullTradeTx)
	}
	if !useCache || err != nil {
		resource := "/0/private/Ledgers"
		headers := make(map[string]string)
		headers["API-Key"] = api.apiKey
		headers["Content-Type"] = "application/json"
		body := url.Values{}
		body.Add("trades", "true")
		offset := 0
		totalTrades := 1000
		for offset < totalTrades {
			body.Set("nonce", strconv.FormatInt(time.Now().UTC().Unix()*1000, 10))
			body.Set("ofs", fmt.Sprint(offset))
			api.sign(headers, body, resource)
			fmt.Println("Getting trades transactions from", offset, "to", offset+50)
			resp, err := api.clientLedgers.R().
				SetHeaders(headers).
				SetFormDataFromValues(body).
				SetResult(&TradesHistory{}).
				Post(api.basePath + resource)
			if err != nil || len((*resp.Result().(*TradesHistory)).Error) > 0 {
				time.Sleep(6 * time.Second)
				body.Set("nonce", strconv.FormatInt(time.Now().UTC().Unix()*1000, 10))
				api.sign(headers, body, resource)
				resp, err = api.clientLedgers.R().
					SetHeaders(headers).
					SetFormDataFromValues(body).
					SetResult(&TradesHistory{}).
					Post(api.basePath + resource)
				if err != nil || len((*resp.Result().(*TradesHistory)).Error) > 0 {
					fmt.Println("Kraken API Trades : Error Requesting TradesHistory" + strings.Join((*resp.Result().(*TradesHistory)).Error, ""))
					return fullTradeTx, errors.New("Kraken API Trades : Error Requesting TradesHistory" + strings.Join((*resp.Result().(*TradesHistory)).Error, ""))
				}
			}
			result := (*resp.Result().(*TradesHistory)).Result.(map[string]interface{})
			totalTrades = int(result["count"].(float64))
			for k, v := range result["ledger"].(map[string]interface{}) {
				fullTradeTx[k] = v
			}
			offset += 50
			time.Sleep(time.Second)
		}
		if useCache {
			err = db.Write("Kraken/private", "Ledgers", fullTradeTx)
			if err != nil {
				return fullTradeTx, errors.New("Kraken API Trades : Error Caching TradesHistory")
			}
		}
	}
	return fullTradeTx, nil
}
