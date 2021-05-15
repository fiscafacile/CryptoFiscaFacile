package poloniex

import (
	"encoding/csv"
	"io"
	"log"
	"strings"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type csvTradesTX struct {
	Date              time.Time
	Market            string          // DOGE/BTC
	FirstCurrency     string          // DOGE
	SecondCurrency    string          // BTC
	Category          string          // Exchange
	Type              string          // Sell
	Price             decimal.Decimal // 0.00000107
	Amount            decimal.Decimal // 24234.19729729
	Total             decimal.Decimal // 0.02593059
	Fee               string          // 0.125%
	OrderNumber       string          // 39236960790
	BaseTotalLessFee  decimal.Decimal // 0.02589818
	QuoteTotalLessFee decimal.Decimal // -24234.19729729
	FeeCurrency       string          // BTC
	FeeTotal          decimal.Decimal // 0.00003241
}

func (pl *Poloniex) ParseTradesCSV(reader io.Reader) (err error) {
	firstTimeUsed := time.Now()
	lastTimeUsed := time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC)
	const SOURCE = "Poloniex Trades CSV :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		alreadyAsked := []string{}
		for _, r := range records {
			if r[0] != "Date" {
				tx := csvTradesTX{}
				tx.Date, err = time.Parse("2006-01-02 15:04:05", r[0])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Date", r[0])
				}
				tx.Market = r[1]
				curr := strings.Split(r[1], "/")
				tx.FirstCurrency = curr[0]
				tx.SecondCurrency = curr[1]
				tx.Category = r[2]
				tx.Type = r[3]
				tx.Price, err = decimal.NewFromString(r[4])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Price", r[4])
				}
				tx.Amount, err = decimal.NewFromString(r[5])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Amount", r[5])
				}
				tx.Total, err = decimal.NewFromString(r[6])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Total", r[6])
				}
				tx.Fee = r[7]
				tx.OrderNumber = r[8]
				tx.BaseTotalLessFee, err = decimal.NewFromString(r[9])
				if err != nil {
					log.Println(SOURCE, "Error Parsing BaseTotalLessFee", r[9])
				}
				tx.QuoteTotalLessFee, err = decimal.NewFromString(r[10])
				if err != nil {
					log.Println(SOURCE, "Error Parsing QuoteTotalLessFee", r[10])
				}
				tx.FeeCurrency = r[11]
				tx.FeeTotal, err = decimal.NewFromString(r[12])
				if err != nil {
					log.Println(SOURCE, "Error Parsing FeeTotal", r[12])
				}
				pl.csvTradesTXs = append(pl.csvTradesTXs, tx)
				if tx.Date.Before(firstTimeUsed) {
					firstTimeUsed = tx.Date
				}
				if tx.Date.After(lastTimeUsed) {
					lastTimeUsed = tx.Date
				}
				// Fill TXsByCategory
				t := wallet.TX{Timestamp: tx.Date, ID: tx.OrderNumber, Note: SOURCE + " " + tx.Market + " " + tx.Type}
				t.Items = make(map[string]wallet.Currencies)
				if tx.Type == "Buy" {
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.FeeCurrency, Amount: tx.Amount.Sub(tx.QuoteTotalLessFee)})
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.FirstCurrency, Amount: tx.Amount})
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.SecondCurrency, Amount: tx.Total})
					pl.TXsByCategory["Exchanges"] = append(pl.TXsByCategory["Exchanges"], t)
				} else if tx.Type == "Sell" {
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.FeeCurrency, Amount: tx.FeeTotal})
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.FirstCurrency, Amount: tx.Amount})
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.SecondCurrency, Amount: tx.Total})
					pl.TXsByCategory["Exchanges"] = append(pl.TXsByCategory["Exchanges"], t)
				} else {
					alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Type, tx, alreadyAsked)
				}
			}
		}
	}
	if _, ok := pl.Sources["Poloniex"]; !ok {
		pl.Sources["Poloniex"] = source.Source{
			Crypto:        true,
			AccountNumber: "emailAROBASEdomainPOINTcom",
			OpeningDate:   firstTimeUsed,
			ClosingDate:   lastTimeUsed,
			LegalName:     "Polo Digital Assets Ltd",
			Address:       "F20, 1st Floor, Eden Plaza,\nEden Island,\nSeychelles",
			URL:           "https://poloniex.com/",
		}
	}
	return
}
