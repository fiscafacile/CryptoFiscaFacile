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

type csvWithdrawalsTX struct {
	Date        time.Time
	Currency    string
	Amount      decimal.Decimal
	FeeDeducted decimal.Decimal
	AmountFee   decimal.Decimal
	Address     string
	Status      string
}

func (pl *Poloniex) ParseWithdrawalsCSV(reader io.Reader, account string) (err error) {
	firstTimeUsed := time.Now()
	lastTimeUsed := time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC)
	const SOURCE = "Poloniex Withdrawals CSV :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "Date" {
				tx := csvWithdrawalsTX{}
				tx.Date, err = time.Parse("2006-01-02 15:04:05", r[0])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Date", r[0])
				}
				tx.Currency = r[1]
				tx.Amount, err = decimal.NewFromString(r[2])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Amount", r[2])
				}
				tx.FeeDeducted, err = decimal.NewFromString(r[3])
				if err != nil {
					log.Println(SOURCE, "Error Parsing FeeDeducted", r[3])
				}
				tx.AmountFee, err = decimal.NewFromString(r[4])
				if err != nil {
					log.Println(SOURCE, "Error Parsing AmountFee", r[4])
				}
				tx.Address = r[5]
				tx.Status = r[6]
				pl.csvWithdrawalsTXs = append(pl.csvWithdrawalsTXs, tx)
				if tx.Date.Before(firstTimeUsed) {
					firstTimeUsed = tx.Date
				}
				if tx.Date.After(lastTimeUsed) {
					lastTimeUsed = tx.Date
				}
				// Fill TXsByCategory
				t := wallet.TX{Timestamp: tx.Date, Note: SOURCE + " " + tx.Address + " " + tx.Status}
				t.Items = make(map[string]wallet.Currencies)
				if !tx.FeeDeducted.IsZero() {
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: tx.FeeDeducted})
				}
				t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.AmountFee})
				pl.TXsByCategory["Withdrawals"] = append(pl.TXsByCategory["Withdrawals"], t)
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
