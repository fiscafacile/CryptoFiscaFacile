package bitfinex

import (
	"encoding/csv"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type CsvTX struct {
	ID          int
	Description string
	Currency    string
	Amount      decimal.Decimal
	Balance     decimal.Decimal
	Date        time.Time
	Wallet      string
}

func (bf *Bitfinex) ParseCSV(reader io.Reader) (err error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "#" {
				tx := CsvTX{}
				id, err := strconv.Atoi(r[0])
				if err != nil {
					log.Println("Error Parsing ID : ", r[0])
				} else {
					tx.ID = id
				}
				tx.Description = r[1]
				tx.Currency = strings.ReplaceAll(r[2], "BAB", "BCH")
				tx.Amount, err = decimal.NewFromString(r[3])
				if err != nil {
					log.Println("Error Parsing Amount : ", r[3])
				}
				tx.Balance, err = decimal.NewFromString(r[4])
				if err != nil {
					log.Println("Error Parsing Balance : ", r[4])
				}
				tx.Date, err = time.Parse("02-01-06 15:04:05", r[5])
				if err != nil {
					log.Println("Error Parsing Date : ", r[5])
				}
				tx.Wallet = r[6]
				bf.CsvTXs = append(bf.CsvTXs, tx)
				// Fill TXsByCategory
				if strings.Contains(tx.Description, "Exchange") ||
					strings.Contains(tx.Description, "Transfer") ||
					strings.Contains(tx.Description, "Settlement") {
					found := false
					for i, ex := range bf.TXsByCategory["Exchanges"] {
						// log.Println(strings.Split(tx.Description, " ")[1], strings.Split(ex.Note, " ")[4])
						if ex.Note == "Bitfinex CSV : "+tx.Description ||
							(ex.SimilarDate(2*time.Second, tx.Date) &&
								strings.Split(strings.Split(ex.Note, " ")[4], ".")[0] == strings.Split(strings.Split(tx.Description, " ")[1], ".")[0] &&
								strings.Split(strings.Split(ex.Note, " ")[4], ".")[1][:1] == strings.Split(strings.Split(tx.Description, " ")[1], ".")[1][:1]) {
							found = true
							if ex.Note != "Bitfinex CSV : "+tx.Description {
								bf.TXsByCategory["Exchanges"][i].Note = "Bitfinex CSV : " + tx.Description
							}
							if bf.TXsByCategory["Exchanges"][i].Items == nil {
								bf.TXsByCategory["Exchanges"][i].Items = make(map[string][]wallet.Currency)
							}
							if tx.Amount.IsPositive() {
								bf.TXsByCategory["Exchanges"][i].Items["To"] = append(bf.TXsByCategory["Exchanges"][i].Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
							} else {
								bf.TXsByCategory["Exchanges"][i].Items["From"] = append(bf.TXsByCategory["Exchanges"][i].Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount.Neg()})
							}
						}
					}
					if !found {
						t := wallet.TX{Timestamp: tx.Date, Note: "Bitfinex CSV : " + tx.Description}
						t.Items = make(map[string][]wallet.Currency)
						if tx.Amount.IsPositive() {
							t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
							bf.TXsByCategory["Exchanges"] = append(bf.TXsByCategory["Exchanges"], t)
						} else {
							t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount.Neg()})
							bf.TXsByCategory["Exchanges"] = append(bf.TXsByCategory["Exchanges"], t)
						}
					}
				} else if strings.Contains(tx.Description, "Trading fees") {
					found := false
					// log.Println(tx.Description)
					for i, ex := range bf.TXsByCategory["Exchanges"] {
						// log.Println(strings.Split(tx.Description, " ")[3], strings.Split(ex.Note, " ")[4])
						if ex.SimilarDate(2*time.Second, tx.Date) &&
							strings.Split(strings.Split(ex.Note, " ")[4], ".")[0] == strings.Split(strings.Split(tx.Description, " ")[3], ".")[0] {
							found = true
							if bf.TXsByCategory["Exchanges"][i].Items == nil {
								bf.TXsByCategory["Exchanges"][i].Items = make(map[string][]wallet.Currency)
							}
							bf.TXsByCategory["Exchanges"][i].Items["Fee"] = append(bf.TXsByCategory["Exchanges"][i].Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount.Neg()})
						}
					}
					if !found {
						t := wallet.TX{Timestamp: tx.Date, Note: "Bitfinex CSV : " + tx.Description}
						t.Items = make(map[string][]wallet.Currency)
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount.Neg()})
						bf.TXsByCategory["Exchanges"] = append(bf.TXsByCategory["Exchanges"], t)
					}
				} else if strings.Contains(tx.Description, "Deposit") ||
					strings.Contains(tx.Description, "fork credit") {
					t := wallet.TX{Timestamp: tx.Date, Note: "Bitfinex CSV : " + tx.Description}
					t.Items = make(map[string][]wallet.Currency)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
					bf.TXsByCategory["Deposits"] = append(bf.TXsByCategory["Deposits"], t)
				} else if strings.Contains(tx.Description, "Withdrawal") ||
					strings.Contains(tx.Description, "fork clear") {
					if strings.Contains(tx.Description, "fee") {
						found := false
						for i, ex := range bf.TXsByCategory["Withdrawals"] {
							if ex.SimilarDate(2*time.Second, tx.Date) {
								found = true
								if bf.TXsByCategory["Withdrawals"][i].Items == nil {
									bf.TXsByCategory["Withdrawals"][i].Items = make(map[string][]wallet.Currency)
								}
								bf.TXsByCategory["Withdrawals"][i].Items["Fee"] = append(bf.TXsByCategory["Withdrawals"][i].Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount.Neg()})
							}
						}
						if !found {
							t := wallet.TX{Timestamp: tx.Date, Note: "Bitfinex CSV : " + tx.Description}
							t.Items = make(map[string][]wallet.Currency)
							t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount.Neg()})
							bf.TXsByCategory["Withdrawals"] = append(bf.TXsByCategory["Withdrawals"], t)
						}
					} else {
						found := false
						for i, ex := range bf.TXsByCategory["Withdrawals"] {
							if ex.SimilarDate(2*time.Second, tx.Date) {
								found = true
								if bf.TXsByCategory["Withdrawals"][i].Items == nil {
									bf.TXsByCategory["Withdrawals"][i].Items = make(map[string][]wallet.Currency)
								}
								bf.TXsByCategory["Withdrawals"][i].Items["From"] = append(bf.TXsByCategory["Withdrawals"][i].Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount.Neg()})
								bf.TXsByCategory["Withdrawals"][i].Note = "Bitfinex CSV : " + tx.Description
							}
						}
						if !found {
							t := wallet.TX{Timestamp: tx.Date, Note: "Bitfinex CSV : " + tx.Description}
							t.Items = make(map[string][]wallet.Currency)
							t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount.Neg()})
							bf.TXsByCategory["Withdrawals"] = append(bf.TXsByCategory["Withdrawals"], t)
						}
					}
				} else {
					log.Println("Bitfinex : Unmanaged ", tx.Description)
				}
			}
		}
	}
	return
}
