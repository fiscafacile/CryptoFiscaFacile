package poloniex

import (
	"encoding/csv"
	"io"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type csvDistributionsTX struct {
	Date     time.Time       // 2020-09-22
	Currency string          // REPV2
	Amount   decimal.Decimal // 82.50028940
	Wallet   string          // exchange
}

func (pl *Poloniex) ParseDistributionsCSV(reader io.Reader, account string) (err error) {
	firstTimeUsed := time.Now()
	lastTimeUsed := time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC)
	const SOURCE = "Poloniex Distributions CSV :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "date" {
				tx := csvDistributionsTX{}
				tx.Date, err = time.Parse("2006-01-02", r[0])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Date", r[0])
				}
				tx.Currency = r[1]
				tx.Amount, err = decimal.NewFromString(r[2])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Amount", r[2])
				}
				tx.Wallet = r[3]
				pl.csvDistributionsTXs = append(pl.csvDistributionsTXs, tx)
				if tx.Date.Before(firstTimeUsed) {
					firstTimeUsed = tx.Date
				}
				if tx.Date.After(lastTimeUsed) {
					lastTimeUsed = tx.Date
				}
				// Fill TXsByCategory
				t := wallet.TX{Timestamp: tx.Date, Note: SOURCE + " " + tx.Wallet}
				t.Items = make(map[string]wallet.Currencies)
				t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
				pl.TXsByCategory["Deposits"] = append(pl.TXsByCategory["Deposits"], t)
			}
		}
	}
	if _, ok := pl.Sources["Poloniex"]; !ok {
		pl.Sources["Poloniex"] = source.Source{
			Crypto:        true,
			AccountNumber: account,
			OpeningDate:   firstTimeUsed,
			ClosingDate:   lastTimeUsed,
			LegalName:     "Polo Digital Assets Ltd",
			Address:       "F20, 1st Floor, Eden Plaza,\nEden Island,\nSeychelles",
			URL:           "https://poloniex.com/",
		}
	}
	return
}
