package bittrex

import (
	"encoding/csv"
	"io"
	"log"
	"strings"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type csvTX struct {
	ID         string
	FromSymbol string
	ToSymbol   string
	Time       time.Time
	Operation  string
	FromAmount decimal.Decimal
	ToAmount   decimal.Decimal
	Fee        decimal.Decimal
	Remark     string
}

func (btrx *Bittrex) ParseCSV(reader io.Reader) (err error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "Uuid" {
				tx := csvTX{}
				tx.ID = r[0]
				symbolSlice := strings.Split(r[1], "-")
				tx.Time, err = time.Parse("1/2/2006 3:04:05 PM", r[14])
				if err != nil {
					log.Println("Error Parsing Time : ", r[14])
				}
				rplcr := strings.NewReplacer(
					"LIMIT_SELL", "SELL",
					"LIMIT_BUY", "BUY",
					"CEILING_MARKET_BUY", "BUY",
				)
				tx.Operation = rplcr.Replace(r[3])
				quantity, err := decimal.NewFromString(r[5])
				if err != nil {
					log.Println("Error Parsing Amount : ", r[5])
				}
				quantityRemaining, err := decimal.NewFromString(r[6])
				if err != nil {
					log.Println("Error Parsing Amount : ", r[6])
				}
				tx.Fee, err = decimal.NewFromString(r[7])
				if err != nil {
					log.Println("Error Parsing Amount : ", r[7])
				}
				price, err := decimal.NewFromString(r[8])
				if err != nil {
					log.Println("Error Parsing Amount : ", r[8])
				}
				// Fill TXsByCategory
				if tx.Operation == "BUY" || tx.Operation == "SELL" {
					if tx.Operation == "BUY" {
						tx.FromSymbol = symbolSlice[0]
						tx.FromAmount = price
						tx.ToSymbol = symbolSlice[1]
						tx.ToAmount = quantity.Sub(quantityRemaining)
					} else if tx.Operation == "SELL" {
						tx.FromSymbol = symbolSlice[1]
						tx.FromAmount = quantity.Sub(quantityRemaining)
						tx.ToSymbol = symbolSlice[0]
						tx.ToAmount = price
					}
					found := false
					for i := range btrx.TXsByCategory["Exchanges"] {
						if tx.ID == btrx.TXsByCategory["Exchanges"][i].ID {
							found = true
						}
					}
					if !found {
						// fmt.Println("Nouvelle transaction :", tx)
						// fmt.Println(tx.Time, "\t", tx.Operation, "\t", "FROM", tx.FromAmount, tx.FromSymbol, "TO", tx.ToAmount, tx.ToSymbol)
						t := wallet.TX{Timestamp: tx.Time, Note: "Bittrex API : " + tx.Operation + " TxID " + tx.ID, ID: tx.ID}
						t.Items = make(map[string]wallet.Currencies)
						t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.FromSymbol, Amount: tx.FromAmount})
						t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.ToSymbol, Amount: tx.ToAmount})
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.FromSymbol, Amount: tx.Fee})
						btrx.TXsByCategory["Exchanges"] = append(btrx.TXsByCategory["Exchanges"], t)
					} else {
						// fmt.Println("Transaction déjà enregistrée : ", tx.ID)

					}
				} else {
					log.Println("Bittrex API : Unmanaged operation -> ", tx.Operation)
				}
			}
		}
	}
	return
}
