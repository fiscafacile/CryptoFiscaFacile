package coinbase

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type CsvTX struct {
	Timestamp time.Time
	Type      string
	Asset     string
	Quantity  decimal.Decimal
	SpotPrice decimal.Decimal
	Subtotal  decimal.Decimal
	Total     decimal.Decimal
	Fees      decimal.Decimal
	Notes     string
}

func (cb *Coinbase) ParseCSV(reader io.ReadSeeker) (err error) {
	buf := make([]byte, 1)
	var found int64
	for found != 5 {
		reader.Read(buf)
		if buf[0] == 'T' {
			found = 1
		} else if buf[0] == 'i' && found == 1 {
			found = 2
		} else if buf[0] == 'm' && found == 2 {
			found = 3
		} else if buf[0] == 'e' && found == 3 {
			found = 4
		} else if buf[0] == 's' && found == 4 {
			found = 5
		} else {
			found = 0
		}
	}
	reader.Seek(-found, os.SEEK_CUR)
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		var fiat string
		for _, r := range records {
			if r[0] == "Timestamp" {
				fiat = strings.Split(r[4], " ")[0]
			} else {
				tx := CsvTX{}
				tx.Timestamp, err = time.Parse("2006-01-02T15:04:05Z", r[0])
				if err != nil {
					log.Println("Error Parsing Timestamp : ", r[0])
				}
				tx.Type = r[1]
				tx.Asset = r[2]
				tx.Quantity, err = decimal.NewFromString(r[3])
				if err != nil {
					log.Println("Error Parsing Quantity : ", r[3])
				}
				tx.SpotPrice, err = decimal.NewFromString(r[4])
				if err != nil {
					log.Println("Error Parsing SpotPrice : ", r[4])
				}
				if r[5] != "" {
					tx.Subtotal, err = decimal.NewFromString(r[5])
					if err != nil {
						log.Println("Error Parsing Subtotal : ", r[5])
					}
				}
				if r[6] != "" {
					tx.Total, err = decimal.NewFromString(r[6])
					if err != nil {
						log.Println("Error Parsing Total : ", r[6])
					}
				}
				if r[7] != "" {
					tx.Fees, err = decimal.NewFromString(r[7])
					if err != nil {
						log.Println("Error Parsing Fees : ", r[7])
					}
				}
				tx.Notes = r[8]
				cb.CsvTXs = append(cb.CsvTXs, tx)
				// Fill TXsByCategory
				if tx.Type == "Receive" {
					t := wallet.TX{Timestamp: tx.Timestamp, Note: "Coinbase CSV : " + tx.Notes}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Asset, Amount: tx.Quantity})
					cb.TXsByCategory["Deposits"] = append(cb.TXsByCategory["Deposits"], t)
				} else if tx.Type == "Send" {
					t := wallet.TX{Timestamp: tx.Timestamp, Note: "Coinbase CSV : " + tx.Notes}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Asset, Amount: tx.Quantity})
					cb.TXsByCategory["Withdrawals"] = append(cb.TXsByCategory["Withdrawals"], t)
				} else if tx.Type == "Sell" {
					t := wallet.TX{Timestamp: tx.Timestamp, Note: "Coinbase CSV : " + tx.Notes}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Asset, Amount: tx.Quantity})
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: fiat, Amount: tx.Subtotal})
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: fiat, Amount: tx.Fees})
					cb.TXsByCategory["Exchanges"] = append(cb.TXsByCategory["Exchanges"], t)
				} else if tx.Type == "Buy" {
					t := wallet.TX{Timestamp: tx.Timestamp, Note: "Coinbase CSV : " + tx.Notes}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Asset, Amount: tx.Quantity})
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: fiat, Amount: tx.Subtotal})
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: fiat, Amount: tx.Fees})
					cb.TXsByCategory["Exchanges"] = append(cb.TXsByCategory["Exchanges"], t)
				} else {
					log.Println("Coinbase : Unmanaged ", tx)
				}
			}
		}
	}
	return
}
