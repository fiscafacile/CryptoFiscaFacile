package ledgerlive

import (
	"encoding/csv"
	"io"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/btc"
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

func (ll *LedgerLive) ParseCSV(reader io.Reader, b *btc.BTC) (err error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "Operation Date" {
				tx := CsvTX{}
				tx.Date, err = time.Parse("2006-01-02T15:04:05.000Z", r[0])
				if err != nil {
					log.Println("Error Parsing Date : ", r[0])
				}
				tx.Currency = r[1]
				tx.Type = r[2]
				tx.Amount, err = decimal.NewFromString(r[3])
				if err != nil {
					log.Println("Error Parsing Amount : ", r[3])
				}
				if r[4] != "" {
					tx.Fees, err = decimal.NewFromString(r[4])
					if err != nil {
						log.Println("Error Parsing Fees : ", r[4])
					}
				}
				tx.Hash = r[5]
				tx.AccountName = r[6]
				tx.AccountXpub = r[7]
				ll.CsvTXs = append(ll.CsvTXs, tx)
				// Fill TXsByCategory
				if tx.Type == "IN" {
					t := wallet.TX{Timestamp: tx.Date, Note: "LedgerLive CSV " + tx.AccountName + " : " + tx.Hash + " -> " + tx.AccountXpub}
					t.Items = make(map[string][]wallet.Currency)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: tx.Fees})
					if is, desc, val, curr := b.IsTxCashIn(tx.Hash); is {
						t.Note += " crypto_purchase " + desc
						t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: curr, Amount: val})
						ll.TXsByCategory["CashIn"] = append(ll.TXsByCategory["CashIn"], t)
					} else {
						ll.TXsByCategory["Deposits"] = append(ll.TXsByCategory["Deposits"], t)
					}
				} else if tx.Type == "OUT" {
					t := wallet.TX{Timestamp: tx.Date, Note: "LedgerLive CSV " + tx.AccountName + " : " + tx.AccountXpub + " -> " + tx.Hash}
					t.Items = make(map[string][]wallet.Currency)
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: tx.Fees})
					if tx.Amount.Sub(tx.Fees).IsZero() {
						ll.TXsByCategory["Fees"] = append(ll.TXsByCategory["Fees"], t)
					} else {
						t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount.Sub(tx.Fees)})
						if is, desc, val, curr := b.IsTxCashOut(tx.Hash); is {
							t.Note += " crypto_payment " + desc
							t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr, Amount: val})
							ll.TXsByCategory["CashOut"] = append(ll.TXsByCategory["CashOut"], t)
						} else {
							ll.TXsByCategory["Withdrawals"] = append(ll.TXsByCategory["Withdrawals"], t)
						}
					}
				} else {
					log.Println("Unmanaged ", tx)
				}
			}
		}
	}
	return
}
