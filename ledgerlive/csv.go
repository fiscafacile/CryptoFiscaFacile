package ledgerlive

import (
	"encoding/csv"
	"io"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/category"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type CsvTX struct {
	Date        time.Time
	Currency    string
	Type        string
	Amount      decimal.Decimal
	Fees        decimal.Decimal
	Hash        string
	AccountName string
	AccountXpub string
}

func (ll *LedgerLive) ParseCSV(reader io.Reader, cat category.Category) (err error) {
	const SOURCE = "LedgerLive CSV"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		alreadyAsked := []string{}
		for _, r := range records {
			if r[0] != "Operation Date" {
				tx := CsvTX{}
				tx.Date, err = time.Parse("2006-01-02T15:04:05.000Z", r[0])
				if err != nil {
					log.Println(SOURCE, ": Error Parsing Date", r[0])
				}
				tx.Currency = r[1]
				tx.Type = r[2]
				tx.Amount, err = decimal.NewFromString(r[3])
				if err != nil {
					log.Println(SOURCE, ": Error Parsing Amount", r[3])
				}
				if r[4] != "" {
					tx.Fees, err = decimal.NewFromString(r[4])
					if err != nil {
						log.Println(SOURCE, ": Error Parsing Fees", r[4])
					}
				}
				tx.Hash = r[5]
				tx.AccountName = r[6]
				tx.AccountXpub = r[7]
				ll.CsvTXs = append(ll.CsvTXs, tx)
			}
		}
		for _, tx := range ll.CsvTXs {
			// Fill TXsByCategory
			if tx.Type == "IN" {
				t := wallet.TX{Timestamp: tx.Date, Note: SOURCE + " " + tx.AccountName + " : " + tx.Hash + " -> " + tx.AccountXpub}
				t.Items = make(map[string]wallet.Currencies)
				t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
				if !tx.Fees.IsZero() {
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: tx.Fees})
				}
				if is, desc, val, curr := cat.IsTxCashIn(tx.Hash); is {
					t.Note += " crypto_purchase " + desc
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: curr, Amount: val})
					ll.TXsByCategory["CashIn"] = append(ll.TXsByCategory["CashIn"], t)
				} else {
					ll.TXsByCategory["Deposits"] = append(ll.TXsByCategory["Deposits"], t)
				}
			} else if tx.Type == "OUT" {
				if !tx.Fees.Equal(tx.Amount) { // ignore Fee associated to other OUT, will be found later
					t := wallet.TX{Timestamp: tx.Date, Note: SOURCE + " " + tx.AccountName + " : " + tx.AccountXpub + " -> " + tx.Hash}
					t.Items = make(map[string]wallet.Currencies)
					if !tx.Fees.IsZero() {
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: tx.Fees})
					}
					for _, tx2 := range ll.CsvTXs {
						if tx2.Hash == tx.Hash &&
							tx2.Type == "OUT" &&
							tx2.Fees.Equal(tx2.Amount) {
							t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx2.Currency, Amount: tx2.Fees})
						}
					}
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount.Sub(tx.Fees)})
					if is, desc, val, curr := cat.IsTxCashOut(tx.Hash); is {
						t.Note += " crypto_payment " + desc
						t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr, Amount: val})
						ll.TXsByCategory["CashOut"] = append(ll.TXsByCategory["CashOut"], t)
					} else {
						ll.TXsByCategory["Withdrawals"] = append(ll.TXsByCategory["Withdrawals"], t)
					}
				}
			} else if tx.Type == "FEES" {
				t := wallet.TX{Timestamp: tx.Date, Note: SOURCE + " " + tx.AccountName + " : " + tx.AccountXpub + " -> " + tx.Hash}
				t.Items = make(map[string]wallet.Currencies)
				t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: tx.Fees})
				ll.TXsByCategory["Fees"] = append(ll.TXsByCategory["Fees"], t)
			} else {
				alreadyAsked = wallet.AskForHelp(SOURCE+" : "+tx.Type, tx, alreadyAsked)
			}
		}
	}
	return
}
