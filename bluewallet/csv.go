package bluewallet

import (
	"encoding/csv"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type CsvTX struct {
	ID        string
	Timestamp time.Time
	Type      string
	Amount    decimal.Decimal
	Size      int
	Fees      decimal.Decimal
}

func (bw *BlueWallet) ParseCSV(reader io.Reader) (err error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "txid" {
				tx := CsvTX{}
				tx.ID = r[0]
				ts, err := strconv.ParseInt(r[1], 10, 64)
				if err != nil {
					log.Println("Error Parsing Timestamp : ", r[1])
				} else {
					tx.Timestamp = time.Unix(ts, 0)
				}
				tx.Type = r[2]
				tx.Amount, err = decimal.NewFromString(r[3])
				if err != nil {
					log.Println("Error Parsing Amount : ", r[3])
				}
				tx.Size, err = strconv.Atoi(r[4])
				if err != nil {
					log.Println("Error Parsing Size : ", r[4])
				}
				tx.Fees, err = decimal.NewFromString(r[5])
				if err != nil {
					log.Println("Error Parsing Fees : ", r[5])
				}
				bw.CsvTXs = append(bw.CsvTXs, tx)
				// Fill TXsByCategory
				if tx.Type == "DEPOSIT" {
					bw.Wallets["BTC"] += tx.Amount
					bw.TXsByCategory["Deposits"] = append(bw.TXsByCategory["Deposits"], wallet.TX{Timestamp: tx.Timestamp, Currency: wallet.Currency{Code: "BTC", Amount: tx.Amount}, Note: "BlueWallet CSV : " + tx.ID})
				} else if tx.Type == "WITHDRAWAL" {
					bw.Wallets["BTC"] -= tx.Amount
					bw.TXsByCategory["Withdrawals"] = append(bw.TXsByCategory["Withdrawals"], wallet.TX{Timestamp: tx.Timestamp, Currency: wallet.Currency{Code: "BTC", Amount: tx.Amount}, Note: tx.ID})
				} else {
					log.Println("Unmanaged ", tx)
				}
			}
		}
	}
	bw.Wallets.Round()
	return
}
