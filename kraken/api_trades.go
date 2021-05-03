package kraken

import (
	"fmt"
	"log"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/category"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/go-resty/resty/v2"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type tradeResponse struct {
	Error  []interface{} `json:"error"`
	Result struct {
		Trades struct {
			Thvrqm33VkhUci7Bs struct {
				Ordertxid string  `json:"ordertxid"`
				Postxid   string  `json:"postxid"`
				Pair      string  `json:"pair"`
				Time      float64 `json:"time"`
				Type      string  `json:"type"`
				Ordertype string  `json:"ordertype"`
				Price     string  `json:"price"`
				Cost      string  `json:"cost"`
				Fee       string  `json:"fee"`
				Vol       string  `json:"vol"`
				Margin    string  `json:"margin"`
				Misc      string  `json:"misc"`
			} `json:"THVRQM-33VKH-UCI7BS"`
		} `json:"trades"`
		Count int `json:"count"`
	} `json:"result"`
}

type apiTradeTX struct {
	Time        time.Time
	Operation   string
	FromSymbol  string
	ToSymbol    string
	FromAmount  decimal.Decimal
	ToAmount    decimal.Decimal
	Fee         decimal.Decimal
	FeeCurrency string
	ID          string
}

func (krkn *Kraken) getTrades(apiKey, apiSecret string, offset int) (tradeTx *resty.Response, err error) {
	// Prepare body payload
	payload := url.Values{}
	payload.Add("trades", "true")
	payload.Add("nonce", strconv.FormatInt(time.Now().UTC().Unix()*1000, 10))
	payload.Add("ofs", fmt.Sprint(offset))
	response, err := krkn.sendRequest(apiKey, apiSecret, "/0/private/TradesHistory", payload)
	return response, err
}

func (krkn *Kraken) GetAllTradeTXs(apiKey, apiSecret string, cat category.Category) {
	useCache := true
	fullTradeTx := make(map[string]interface{})
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Kraken", "trades", &fullTradeTx)
	}
	if !useCache || err != nil {
		// Request API with an offset until all the records are retrieved
		offset := 0
		lastResponseCount := 50
		for lastResponseCount == 50 {
			fmt.Println("Getting trades transactions from", offset, "to", offset+50)
			response, err := krkn.getTrades(apiKey, apiSecret, offset)
			if err != nil {
				time.Sleep(6 * time.Second)
				response, err = krkn.getTrades(apiKey, apiSecret, offset)
				if err != nil {
					log.Println("Kraken API : Error while fetching trades", err)
				}
			}
			trades := response.Result().(*TradesHistory).Result.(map[string]interface{})["trades"].(map[string]interface{})
			lastResponseCount = len(trades)
			offset += 50
			for k, v := range trades {
				fullTradeTx[k] = v
			}
			time.Sleep(time.Second)
		}
		if useCache {
			err = db.Write("Kraken", "trades", fullTradeTx)
			if err != nil {
				log.Println("Kraken API : Error while caching trades", err)
			}
		}
	}
	// Process trade transactions
	for txID, transaction := range fullTradeTx {
		trd := transaction.(map[string]interface{})
		tx := apiTradeTX{}
		sec, dec := math.Modf(trd["time"].(float64))
		tx.Time = time.Unix(int64(sec), int64(dec*(1e9)))
		tx.Operation = trd["type"].(string)
		tx.Fee, err = decimal.NewFromString(trd["fee"].(string))
		if err != nil {
			fmt.Println("Error while parsing fee", trd["fee"].(string))
		}
		tx.ID = txID
		vol, err := decimal.NewFromString(trd["vol"].(string))
		if err != nil {
			fmt.Println("Error while parsing vol", trd["vol"].(string))
		}
		cost, err := decimal.NewFromString(trd["cost"].(string))
		if err != nil {
			fmt.Println("Error while parsing cost", trd["cost"].(string))
		}
		firstSymbol := strings.ReplaceAll(trd["pair"].(string)[len(trd["pair"].(string))/2-3:len(trd["pair"].(string))/2], "XBT", "BTC")
		secondSymbol := strings.ReplaceAll(trd["pair"].(string)[len(trd["pair"].(string))/2+1:], "XBT", "BTC")
		tx.FeeCurrency = secondSymbol
		if tx.Operation == "buy" || tx.Operation == "sell" {
			if tx.Operation == "buy" {
				tx.FromSymbol = secondSymbol
				tx.FromAmount = cost
				tx.ToSymbol = firstSymbol
				tx.ToAmount = vol
			} else if tx.Operation == "sell" {
				tx.FromSymbol = firstSymbol
				tx.FromAmount = vol
				tx.ToSymbol = secondSymbol
				tx.ToAmount = cost
			}
			found := false
			for i := range krkn.TXsByCategory["Exchanges"] {
				if tx.ID == krkn.TXsByCategory["Exchanges"][i].ID {
					found = true
				}
			}
			if !found {
				// fmt.Println("Nouvelle transaction :", tx)
				fmt.Println(tx.Time, "\t", tx.Operation, "\t", "FROM", tx.FromAmount, tx.FromSymbol, "TO", tx.ToAmount, tx.ToSymbol)
				t := wallet.TX{Timestamp: tx.Time, Note: "Kraken API : " + tx.Operation + " TxID " + tx.ID, ID: tx.ID}
				t.Items = make(map[string]wallet.Currencies)
				t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.FromSymbol, Amount: tx.FromAmount})
				t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.ToSymbol, Amount: tx.ToAmount})
				t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.FeeCurrency, Amount: tx.Fee})
				krkn.TXsByCategory["Exchanges"] = append(krkn.TXsByCategory["Exchanges"], t)
			} else {
				// fmt.Println("Transaction déjà enregistrée : ", tx.ID)
			}
		} else {
			log.Println("Kraken API : Unmanaged operation -> ", tx.Operation)
		}
	}
	krkn.tradesDone <- err
}
