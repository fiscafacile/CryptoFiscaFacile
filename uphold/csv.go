package uphold

import (
	"encoding/csv"
	"io"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/category"
	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type CsvTX struct {
	Date                time.Time
	Destination         string
	DestinationAmount   decimal.Decimal
	DestinationCurrency string
	FeeAmount           decimal.Decimal
	FeeCurrency         string
	ID                  string
	Origin              string
	OriginAmount        decimal.Decimal
	OriginCurrency      string
	Status              string
	Type                string
}

func (uh *Uphold) ParseCSV(reader io.Reader, cat category.Category, account string) (err error) {
	firstTimeUsed := time.Now()
	lastTimeUsed := time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC)
	const SOURCE = "Uphold CSV :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		alreadyAsked := []string{}
		for _, r := range records {
			if r[0] != "Date" {
				tx := CsvTX{}
				tx.Date, err = time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700", r[0])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Date :", r[0])
				}
				tx.Destination = r[1]
				tx.DestinationAmount, err = decimal.NewFromString(r[2])
				if err != nil {
					log.Println(SOURCE, "Error Parsing DestinationAmount :", r[2])
				}
				tx.DestinationCurrency = r[3]
				if r[4] != "" {
					tx.FeeAmount, err = decimal.NewFromString(r[4])
					if err != nil {
						log.Println(SOURCE, "Error Parsing FeeAmount :", r[4])
					}
				}
				tx.FeeCurrency = r[5]
				tx.ID = r[6]
				tx.Origin = r[7]
				tx.OriginAmount, err = decimal.NewFromString(r[8])
				if err != nil {
					log.Println(SOURCE, "Error Parsing OriginAmount :", r[8])
				}
				tx.OriginCurrency = r[9]
				tx.Status = r[10]
				tx.Type = r[11]
				uh.CsvTXs = append(uh.CsvTXs, tx)
				if tx.Date.Before(firstTimeUsed) {
					firstTimeUsed = tx.Date
				}
				if tx.Date.After(lastTimeUsed) {
					lastTimeUsed = tx.Date
				}
				// Fill TXsByCategory
				if tx.OriginCurrency != tx.DestinationCurrency {
					t := wallet.TX{Timestamp: tx.Date, ID: "e" + tx.ID, Note: SOURCE + " exchange"}
					t.Items = make(map[string]wallet.Currencies)
					if tx.FeeCurrency == tx.DestinationCurrency &&
						!tx.FeeAmount.IsZero() {
						t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.DestinationCurrency, Amount: tx.DestinationAmount.Add(tx.FeeAmount)})
					} else {
						t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.DestinationCurrency, Amount: tx.DestinationAmount})
					}
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.OriginCurrency, Amount: tx.OriginAmount})
					uh.TXsByCategory["Exchanges"] = append(uh.TXsByCategory["Exchanges"], t)
				}
				if tx.Type == "in" {
					t := wallet.TX{Timestamp: tx.Date, ID: tx.ID, Note: SOURCE + " " + tx.Type + " " + tx.Status}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.DestinationCurrency, Amount: tx.DestinationAmount})
					if !tx.FeeAmount.IsZero() {
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.FeeCurrency, Amount: tx.FeeAmount})
					}
					if is, desc := cat.IsTxInterest(t.ID); is {
						t.Note += " interest " + desc
						uh.TXsByCategory["Interests"] = append(uh.TXsByCategory["Interests"], t)
					} else {
						uh.TXsByCategory["Deposits"] = append(uh.TXsByCategory["Deposits"], t)
					}
				} else if tx.Type == "out" {
					t := wallet.TX{Timestamp: tx.Date, ID: tx.ID, Note: SOURCE + " " + tx.Type + " " + tx.Status}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.DestinationCurrency, Amount: tx.DestinationAmount})
					if !tx.FeeAmount.IsZero() {
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.FeeCurrency, Amount: tx.FeeAmount})
					}
					if is, desc := cat.IsTxGift(t.ID); is {
						t.Note += " gift " + desc
						uh.TXsByCategory["Gifts"] = append(uh.TXsByCategory["Gifts"], t)
					} else {
						uh.TXsByCategory["Withdrawals"] = append(uh.TXsByCategory["Withdrawals"], t)
					}
				} else {
					alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Type, tx, alreadyAsked)
				}
			}
		}
	}
	uh.Sources["Uphold"] = source.Source{
		Crypto:        true,
		AccountNumber: account,
		OpeningDate:   firstTimeUsed,
		ClosingDate:   lastTimeUsed,
		LegalName:     "Uphold Europe Limited",
		Address:       "Suite A, 6 Honduras Street, London, England, EC1Y 0TH\nRoyaume-Uni",
		URL:           "https://uphold.com",
	}
	return
}
