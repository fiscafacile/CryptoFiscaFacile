package hitbtc

import (
	"encoding/csv"
	"io"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type csvTransactionTX struct {
	Email              string
	Date               time.Time
	OperationID        string
	Type               string
	Amount             decimal.Decimal
	Hash               string
	MainAccountBalance decimal.Decimal
	Currency           string
}

func (hb *HitBTC) ParseCSVTransactions(reader io.Reader) (err error) {
	const SOURCE = "HitBTC CSV Transactions :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		alreadyAsked := []string{}
		for _, r := range records {
			if r[0] != "Email" {
				tx := csvTransactionTX{}
				tx.Email = r[0]
				tx.Date, err = time.Parse("2006-01-02 15:04:05", r[1])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Date", r[1])
				}
				tx.OperationID = r[2]
				tx.Type = r[3]
				tx.Amount, err = decimal.NewFromString(r[4])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Amount", r[4])
				}
				tx.Hash = r[5]
				tx.MainAccountBalance, err = decimal.NewFromString(r[6])
				if err != nil {
					log.Println(SOURCE, "Error Parsing MainAccountBalance", r[6])
				}
				tx.Currency = r[7]
				hb.csvTransactionTXs = append(hb.csvTransactionTXs, tx)
				// Fill TXsByCategory
				if tx.Type == "Deposit" {
					t := wallet.TX{Timestamp: tx.Date, ID: tx.Hash, Note: SOURCE + " " + tx.OperationID}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
					hb.TXsByCategory["Deposits"] = append(hb.TXsByCategory["Deposits"], t)
				} else if tx.Type == "Withdrawal" {
					t := wallet.TX{Timestamp: tx.Date, ID: tx.Hash, Note: SOURCE + " " + tx.OperationID}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
					hb.TXsByCategory["Withdrawals"] = append(hb.TXsByCategory["Withdrawals"], t)
				} else if tx.Type == "Transfer to main account" ||
					tx.Type == "Transfer to trading account" {
					// Do not use, it is Source internal transfers
				} else {
					alreadyAsked = wallet.AskForHelp(SOURCE+tx.Type, tx, alreadyAsked)
				}
			}
		}
	}
	return
}
