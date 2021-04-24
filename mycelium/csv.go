package mycelium

import (
	"encoding/csv"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type CsvTX struct {
	Account     string
	ID          string
	DestAddress string
	Timestamp   time.Time
	Value       decimal.Decimal
	Currency    string
	Label       string
}

func (mc *MyCelium) ParseCSV(reader io.Reader) (err error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "Account" {
				tx := CsvTX{}
				tx.Account = r[0]
				tx.ID = r[1]
				tx.DestAddress = r[2]
				tx.Timestamp, err = time.Parse("2006-01-02T15:04Z", r[3])
				if err != nil {
					log.Println("Error Parsing Timestamp : ", r[3])
				}
				tx.Value, err = decimal.NewFromString(r[4])
				if err != nil {
					log.Println("Error Parsing Value : ", r[4])
				}
				tx.Currency = r[5]
				tx.Label = r[6]
				mc.CsvTXs = append(mc.CsvTXs, tx)
				// Fill TXsByCategory
				if tx.Value.IsPositive() {
					t := wallet.TX{Timestamp: tx.Timestamp, Note: "MyCelium CSV : " + tx.ID + " from " + tx.DestAddress + " " + tx.Label}
					t.Items = make(map[string][]wallet.Currency)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: ticker(tx.Currency), Amount: tx.Value})
					// t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: tx.Fees})
					mc.TXsByCategory["Deposits"] = append(mc.TXsByCategory["Deposits"], t)
				} else if tx.Value.IsNegative() {
					t := wallet.TX{Timestamp: tx.Timestamp, Note: "MyCelium CSV : " + tx.ID + " to " + tx.DestAddress + " " + tx.Label}
					t.Items = make(map[string][]wallet.Currency)
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: ticker(tx.Currency), Amount: tx.Value.Neg()})
					if strings.Contains(tx.Label, "crypto_payment") {
						r := regexp.MustCompile(`\(([+-]?([0-9]*[.,])?[0-9]+)€\)`)
						if r.MatchString(tx.Label) {
							to, err := decimal.NewFromString(r.FindStringSubmatch(tx.Label)[1])
							if err == nil {
								t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "EUR", Amount: to})
							}
						}
						mc.TXsByCategory["CashOut"] = append(mc.TXsByCategory["CashOut"], t)
					} else {
						// t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: tx.Fees})
						mc.TXsByCategory["Withdrawals"] = append(mc.TXsByCategory["Withdrawals"], t)
					}
				} else {
					log.Println("Unmanaged ", tx)
				}
			}
		}
	}
	return
}

func ticker(currency string) string {
	if currency == "Bitcoin" {
		return "BTC"
	}
	return "???"
}
