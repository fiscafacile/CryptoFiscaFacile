package bitstamp

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

type csvTX struct {
	Type      string
	DateTime  time.Time
	Account   string
	Amount    decimal.Decimal
	Symbol    string
	ToAmount  decimal.Decimal
	ToSymbol  string
	Rate      string
	Fee       decimal.Decimal
	FeeSymbol string
	SubType   string
}

func (bs *Bitstamp) ParseCSV(reader io.Reader, account string) (err error) {
	firstTimeUsed := time.Now()
	lastTimeUsed := time.Date(2019, time.November, 14, 0, 0, 0, 0, time.UTC)
	const SOURCE = "Bitstamp CSV :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		alreadyAsked := []string{}
		for _, r := range records {
			if r[0] != "Type" {
				tx := csvTX{}
				tx.Type = r[0]
				tx.DateTime, err = time.Parse("Jan. 02, 2006, 03:04 PM", r[1])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Date", r[1])
				}
				tx.Account = r[2]
				curr := strings.Split(r[3], " ")
				tx.Amount, err = decimal.NewFromString(curr[0])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Amount", curr[0])
				}
				tx.Symbol = curr[1]
				if r[4] != "" {
					toCurr := strings.Split(r[4], " ")
					tx.ToAmount, err = decimal.NewFromString(toCurr[0])
					if err != nil {
						log.Println(SOURCE, "Error Parsing ToAmount", toCurr[0])
					}
					tx.ToSymbol = toCurr[1]
				}
				tx.Rate = r[5]
				if r[6] != "" {
					fee := strings.Split(r[6], " ")
					tx.Fee, err = decimal.NewFromString(fee[0])
					if err != nil {
						log.Println(SOURCE, "Error Parsing Fee", fee[0])
					}
					tx.FeeSymbol = fee[1]
				}
				tx.SubType = r[7]
				bs.csvTXs = append(bs.csvTXs, tx)
				// Fill TXsByCategory
				t := wallet.TX{Timestamp: tx.DateTime, Note: SOURCE + " " + tx.Type + " " + tx.SubType}
				t.Items = make(map[string]wallet.Currencies)
				if !tx.Fee.IsZero() {
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.FeeSymbol, Amount: tx.Fee})
				}
				if tx.Type == "Deposit" {
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Symbol, Amount: tx.Amount})
					bs.TXsByCategory["Deposits"] = append(bs.TXsByCategory["Deposits"], t)
				} else if tx.Type == "Withdrawal" {
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Symbol, Amount: tx.Amount})
					bs.TXsByCategory["Withdrawals"] = append(bs.TXsByCategory["Withdrawals"], t)
				} else if tx.Type == "Market" {
					if tx.SubType == "Buy" {
						t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Symbol, Amount: tx.Amount})
						t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.ToSymbol, Amount: tx.ToAmount})
						if !tx.Fee.IsZero() && tx.FeeSymbol != "" {
							t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.FeeSymbol, Amount: tx.Fee})
						}
						bs.TXsByCategory["Exchanges"] = append(bs.TXsByCategory["Exchanges"], t)
					} else if tx.SubType == "Sell" {
						t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Symbol, Amount: tx.Amount})
						t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.ToSymbol, Amount: tx.ToAmount})
						if !tx.Fee.IsZero() && tx.FeeSymbol != "" {
							t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.FeeSymbol, Amount: tx.Fee})
						}
						bs.TXsByCategory["Exchanges"] = append(bs.TXsByCategory["Exchanges"], t)
					} else {
						alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Type+" "+tx.SubType, tx, alreadyAsked)
					}
				} else if tx.Type == "Crypto currency purchase" {
					// ignore
				} else {
					alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Type+" "+tx.SubType, tx, alreadyAsked)
				}
				if tx.DateTime.Before(firstTimeUsed) {
					firstTimeUsed = tx.DateTime
				}
				if tx.DateTime.After(lastTimeUsed) {
					lastTimeUsed = tx.DateTime
				}
			}
		}
		if _, ok := bs.Sources["Bitstamp"]; !ok {
			bs.Sources["Bitstamp"] = source.Source{
				Crypto:        true,
				AccountNumber: account,
				OpeningDate:   firstTimeUsed,
				ClosingDate:   lastTimeUsed,
				LegalName:     "Bitstamp Ltd",
				Address:       "5 New Street Square,\nLondon EC4A 3TW,\nRoyaume-Uni",
				URL:           "https://bitstamp.com",
			}
		}
	}
	return
}
