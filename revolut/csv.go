package revolut

import (
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type CsvTX struct {
	Timestamp   time.Time
	Description string
	Rate        decimal.Decimal
	PaidOut     decimal.Decimal
	PaidIn      decimal.Decimal
	ExchangeOut wallet.Currency
	ExchangeIn  string
	Balance     decimal.Decimal
	Category    string
	Notes       string
}

func (revo *Revolut) ParseCSV(reader io.Reader, account string) (err error) {
	firstTimeUsed := time.Now()
	lastTimeUsed := time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC)
	const SOURCE = "Revolut CSV :"
	raw, _ := ioutil.ReadAll(reader)
	raw = bytes.ReplaceAll(raw, []byte(" , "), []byte(","))
	raw = bytes.ReplaceAll(raw, []byte(", "), []byte(","))
	csvReader := csv.NewReader(bytes.NewReader(raw))
	records, err := csvReader.ReadAll()
	if err == nil {
		alreadyAsked := []string{}
		var curr string
		for _, r := range records {
			if r[0] == "Completed Date" {
				curr = strings.Split(r[2], "(")[1]
				curr = strings.Split(curr, ")")[0]
			} else {
				tx := CsvTX{}
				tx.Timestamp, err = time.Parse("2 Jan 2006", f2e(r[0]))
				if err != nil {
					tx.Timestamp, err = time.Parse("Jan 2, 2006", r[0])
					if err != nil {
						log.Println(SOURCE, "Error Parsing Timestamp :", r[0])
					}
				}
				tx.Description = strings.ReplaceAll(r[1], "\u00a0", "")
				fields := strings.Split(tx.Description, " ")
				for i := 0; i < len(fields); i++ {
					if strings.Contains(fields[i], "€") {
						tx.Rate, err = decimal.NewFromString(strings.ReplaceAll(fields[i], "€", ""))
						if err != nil {
							log.Println(SOURCE, "Error Parsing Rate :", strings.ReplaceAll(fields[i], "€", ""))
						}
						break
					}
				}
				if r[2] != "" {
					tx.PaidOut, err = decimal.NewFromString(r[2])
					if err != nil {
						log.Println(SOURCE, "Error Parsing PaidOut :", r[2])
					}
				} else {
					tx.PaidIn, err = decimal.NewFromString(r[3])
					if err != nil {
						log.Println(SOURCE, "Error Parsing PaidIn :", r[3])
					}
				}
				s := strings.Split(r[4], " ")
				tx.ExchangeOut.Code = s[0]
				tx.ExchangeOut.Amount, err = decimal.NewFromString(s[1])
				if err != nil {
					log.Println(SOURCE, "Error Parsing ExchangeOut.Amount :", s[1])
				}
				tx.ExchangeIn = r[5]
				tx.Balance, err = decimal.NewFromString(r[6])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Balance :", r[6])
				}
				tx.Category = r[7]
				tx.Notes = r[8]
				revo.CsvTXs = append(revo.CsvTXs, tx)
				if tx.Timestamp.Before(firstTimeUsed) {
					firstTimeUsed = tx.Timestamp
				}
				if tx.Timestamp.After(lastTimeUsed) {
					lastTimeUsed = tx.Timestamp
				}
				// Fill TXsByCategory
				if !tx.PaidIn.IsZero() {
					t := wallet.TX{Timestamp: tx.Timestamp, Note: SOURCE + " " + tx.Description}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr, Amount: tx.PaidIn})
					t.Items["From"] = append(t.Items["From"], tx.ExchangeOut)
					revo.TXsByCategory["Exchanges"] = append(revo.TXsByCategory["Exchanges"], t)
				} else if !tx.PaidOut.IsZero() {
					t := wallet.TX{Timestamp: tx.Timestamp, Note: SOURCE + " " + tx.Description}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: curr, Amount: tx.PaidOut})
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "EUR", Amount: tx.PaidOut.Mul(tx.Rate)})
					revo.TXsByCategory["Exchanges"] = append(revo.TXsByCategory["Exchanges"], t)
				} else {
					alreadyAsked = wallet.AskForHelp(SOURCE, tx, alreadyAsked)
				}
			}
		}
	}
	revo.Sources["Revolut"] = source.Source{
		Crypto:        true,
		AccountNumber: account,
		OpeningDate:   firstTimeUsed,
		ClosingDate:   lastTimeUsed,
		LegalName:     "Revolut Limited",
		Address:       "4th Floor, 7 Westferry Circus\nE14 4HD Londres, Royaume-Uni",
		URL:           "https://www.revolut.com",
	}
	return
}

func f2e(french string) (english string) {
	english = strings.ReplaceAll(french, "jan.", "Jan")
	english = strings.ReplaceAll(english, "févr.", "Feb")
	english = strings.ReplaceAll(english, "mars", "Mar")
	english = strings.ReplaceAll(english, "avr.", "Apr")
	english = strings.ReplaceAll(english, "mai", "May")
	english = strings.ReplaceAll(english, "juin", "June")
	english = strings.ReplaceAll(english, "juil.", "July")
	english = strings.ReplaceAll(english, "août", "Aug")
	english = strings.ReplaceAll(english, "sept.", "Sep")
	english = strings.ReplaceAll(english, "oct.", "Oct")
	english = strings.ReplaceAll(english, "nov.", "Nov")
	english = strings.ReplaceAll(english, "déc.", "Dec")
	return
}
