package binance

import (
	"encoding/csv"
	"io"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type csvTX struct {
	Time      time.Time
	Account   string
	Operation string
	Coin      string
	Change    decimal.Decimal
	Remark    string
}

func (b *Binance) ParseCSV(reader io.Reader) (err error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "UTC_Time" {
				tx := csvTX{}
				tx.Time, err = time.Parse("2006-01-02 15:04:05", r[0])
				if err != nil {
					log.Println("Error Parsing Time : ", r[0])
				}
				tx.Account = r[1]
				tx.Operation = r[2]
				tx.Coin = r[3]
				tx.Change, err = decimal.NewFromString(r[4])
				if err != nil {
					log.Println("Error Parsing Amount : ", r[4])
				}
				tx.Remark = r[5]
				b.csvTXs = append(b.csvTXs, tx)
				// Fill TXsByCategory
				if tx.Operation == "Buy" ||
					tx.Operation == "Sell" ||
					tx.Operation == "Fee" {
					found := false
					for i, ex := range b.TXsByCategory["Exchanges"] {
						if ex.SimilarDate(2*time.Second, tx.Time) {
							found = true
							if b.TXsByCategory["Exchanges"][i].Items == nil {
								b.TXsByCategory["Exchanges"][i].Items = make(map[string][]wallet.Currency)
							}
							if tx.Change.IsPositive() {
								b.TXsByCategory["Exchanges"][i].Items["To"] = append(b.TXsByCategory["Exchanges"][i].Items["To"], wallet.Currency{Code: tx.Coin, Amount: tx.Change})
							} else {
								b.TXsByCategory["Exchanges"][i].Items["From"] = append(b.TXsByCategory["Exchanges"][i].Items["From"], wallet.Currency{Code: tx.Coin, Amount: tx.Change.Neg()})
							}
						}
					}
					if !found {
						t := wallet.TX{Timestamp: tx.Time, Note: "Binance CSV : Buy Sell Fee " + tx.Remark}
						t.Items = make(map[string][]wallet.Currency)
						if tx.Change.IsPositive() {
							t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Coin, Amount: tx.Change})
							b.TXsByCategory["Exchanges"] = append(b.TXsByCategory["Exchanges"], t)
						} else {
							t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Coin, Amount: tx.Change.Neg()})
							b.TXsByCategory["Exchanges"] = append(b.TXsByCategory["Exchanges"], t)
						}
					}
				} else if tx.Operation == "Deposit" ||
					tx.Operation == "Distribution" {
					t := wallet.TX{Timestamp: tx.Time, Note: "Binance CSV : " + tx.Operation + " " + tx.Remark}
					t.Items = make(map[string][]wallet.Currency)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Coin, Amount: tx.Change})
					b.TXsByCategory["Deposits"] = append(b.TXsByCategory["Deposits"], t)
				} else if tx.Operation == "Withdraw" {
					t := wallet.TX{Timestamp: tx.Time, Note: "Binance CSV : " + tx.Operation + " " + tx.Remark}
					t.Items = make(map[string][]wallet.Currency)
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Coin, Amount: tx.Change.Neg()})
					b.TXsByCategory["Withdrawals"] = append(b.TXsByCategory["Withdrawals"], t)
				} else {
					log.Println("Binance : Unmanaged ", tx.Operation)
				}
			}
		}
	}
	return
}
