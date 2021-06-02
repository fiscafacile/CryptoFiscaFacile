package poloniex

import (
	"encoding/csv"
	"io"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/utils"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type csvDepositsTX struct {
	Date     time.Time
	ID       string
	Currency string
	Amount   decimal.Decimal
	Address  string
	Status   string
}

func (pl *Poloniex) ParseDepositsCSV(reader io.Reader, account string) (err error) {
	firstTimeUsed := time.Now()
	lastTimeUsed := time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC)
	const SOURCE = "Poloniex Deposits CSV :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "Date" {
				tx := csvDepositsTX{}
				tx.Date, err = time.Parse("2006-01-02 15:04:05", r[0])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Date", r[0])
				}
				tx.ID = utils.GetUniqueID(SOURCE + tx.Date.String())
				tx.Currency = r[1]
				tx.Amount, err = decimal.NewFromString(r[2])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Amount", r[2])
				}
				tx.Address = r[3]
				tx.Status = r[4]
				pl.csvDepositsTXs = append(pl.csvDepositsTXs, tx)
				if tx.Date.Before(firstTimeUsed) {
					firstTimeUsed = tx.Date
				}
				if tx.Date.After(lastTimeUsed) {
					lastTimeUsed = tx.Date
				}
				// Fill TXsByCategory
				t := wallet.TX{Timestamp: tx.Date, ID: tx.ID, Note: SOURCE + " " + tx.Address + " " + tx.Status}
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
