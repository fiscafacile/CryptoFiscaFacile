package coinbasepro

import (
	"encoding/csv"
	"io"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type CsvAccountTX struct {
	Portfolio         string
	Type              string
	Time              time.Time
	Amount            decimal.Decimal
	Balance           decimal.Decimal
	AmountBalanceUnit string
	TransferID        string
	TradeID           string
	OrderID           string
}

func (cbp *CoinbasePro) ParseAccountCSV(reader io.ReadSeeker, account string) (err error) {
	firstTimeUsed := time.Now()
	lastTimeUsed := time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC)
	const SOURCE = "CoinbasePro Account CSV :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		alreadyAsked := []string{}
		for _, r := range records {
			if r[0] != "portfolio" {
				tx := CsvAccountTX{}
				tx.Portfolio = r[0]
				tx.Type = r[1]
				tx.Time, err = time.Parse("2006-01-02T15:04:05.999Z", r[2])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Time : ", r[2])
				}
				tx.Amount, err = decimal.NewFromString(r[3])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Amount : ", r[3])
				}
				tx.Balance, err = decimal.NewFromString(r[4])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Balance : ", r[4])
				}
				tx.AmountBalanceUnit = r[5]
				tx.TransferID = r[6]
				tx.TradeID = r[7]
				tx.OrderID = r[8]
				cbp.CsvAccountTXs = append(cbp.CsvAccountTXs, tx)
				if tx.Time.Before(firstTimeUsed) {
					firstTimeUsed = tx.Time
				}
				if tx.Time.After(lastTimeUsed) {
					lastTimeUsed = tx.Time
				}
				// Fill TXsByCategory
				if tx.Type == "deposit" {
					t := wallet.TX{Timestamp: tx.Time, ID: tx.TransferID, Note: SOURCE + " portfolio " + tx.Portfolio}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.AmountBalanceUnit, Amount: tx.Amount})
					cbp.TXsByCategory["Deposits"] = append(cbp.TXsByCategory["Deposits"], t)
				} else if tx.Type == "withdrawal" {
					t := wallet.TX{Timestamp: tx.Time, ID: tx.TransferID, Note: SOURCE + " portfolio " + tx.Portfolio}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.AmountBalanceUnit, Amount: tx.Amount.Neg()})
					cbp.TXsByCategory["Withdrawals"] = append(cbp.TXsByCategory["Withdrawals"], t)
				} else if tx.Type == "match" ||
					tx.Type == "fee" {
					found := false
					for i, ex := range cbp.TXsByCategory["Exchanges"] {
						if ex.ID == tx.OrderID+"-"+tx.TradeID {
							found = true
							if tx.Type == "fee" {
								cbp.TXsByCategory["Exchanges"][i].Items["Fee"] = append(cbp.TXsByCategory["Exchanges"][i].Items["Fee"], wallet.Currency{Code: tx.AmountBalanceUnit, Amount: tx.Amount.Neg()})
							} else {
								if tx.Amount.IsPositive() {
									cbp.TXsByCategory["Exchanges"][i].Items["To"] = append(cbp.TXsByCategory["Exchanges"][i].Items["To"], wallet.Currency{Code: tx.AmountBalanceUnit, Amount: tx.Amount})
								} else {
									cbp.TXsByCategory["Exchanges"][i].Items["From"] = append(cbp.TXsByCategory["Exchanges"][i].Items["From"], wallet.Currency{Code: tx.AmountBalanceUnit, Amount: tx.Amount.Neg()})
								}
							}
						}
					}
					if !found {
						t := wallet.TX{Timestamp: tx.Time, ID: tx.OrderID + "-" + tx.TradeID, Note: SOURCE + " portfolio " + tx.Portfolio}
						t.Items = make(map[string]wallet.Currencies)
						if tx.Type == "fee" {
							t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.AmountBalanceUnit, Amount: tx.Amount.Neg()})
						} else {
							if tx.Amount.IsPositive() {
								t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.AmountBalanceUnit, Amount: tx.Amount})
							} else {
								t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.AmountBalanceUnit, Amount: tx.Amount.Neg()})
							}
						}
						cbp.TXsByCategory["Exchanges"] = append(cbp.TXsByCategory["Exchanges"], t)
					}
				} else {
					alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Type, tx, alreadyAsked)
				}
			}
		}
	}
	if _, ok := cbp.Sources["CoinbasePro"]; ok {
		if cbp.Sources["CoinbasePro"].OpeningDate.After(firstTimeUsed) {
			src := cbp.Sources["CoinbasePro"]
			src.OpeningDate = firstTimeUsed
			cbp.Sources["CoinbasePro"] = src
		}
		if cbp.Sources["CoinbasePro"].ClosingDate.Before(lastTimeUsed) {
			src := cbp.Sources["CoinbasePro"]
			src.ClosingDate = lastTimeUsed
			cbp.Sources["CoinbasePro"] = src
		}
	} else {
		cbp.Sources["CoinbasePro"] = source.Source{
			Crypto:        true,
			AccountNumber: account,
			OpeningDate:   firstTimeUsed,
			ClosingDate:   lastTimeUsed,
			LegalName:     "Coinbase Europe Limited",
			Address:       "70 Sir John Rogerson’s Quay,\nDublin D02 R296\nIrlande",
			URL:           "https://pro.coinbase.com",
		}
	}
	return
}
