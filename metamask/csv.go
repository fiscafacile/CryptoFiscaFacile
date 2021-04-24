package metamask

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
	BlockNumber uint64
	Hash        string
	Timestamp   time.Time
	Value       decimal.Decimal
	Coin        string
	Type        string
	From        string
	To          string
	Gas         uint64
	GasUsed     uint64
	GasPrice    uint64
}

func (mm *MetaMask) ParseCSV(reader io.Reader) (err error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "blockNumber" {
				tx := CsvTX{}
				tx.BlockNumber, err = strconv.ParseUint(r[0], 10, 64)
				if err != nil {
					log.Println("Error Parsing BlockNumber : ", r[0])
				}
				tx.Hash = r[1]
				ts, err := strconv.ParseInt(r[2], 10, 64)
				if err != nil {
					log.Println("Error Parsing Timestamp : ", r[2])
				} else {
					tx.Timestamp = time.Unix(ts, 0)
				}
				tx.Value, err = decimal.NewFromString(r[3])
				if err != nil {
					log.Println("Error Parsing Value : ", r[3])
				}
				tx.Coin = r[4]
				tx.Type = r[5]
				tx.From = r[6]
				tx.To = r[7]
				tx.Gas, err = strconv.ParseUint(r[8], 10, 64)
				if err != nil {
					log.Println("Error Parsing Gas : ", r[8])
				}
				tx.GasUsed, err = strconv.ParseUint(r[9], 10, 64)
				if err != nil {
					log.Println("Error Parsing GasUsed : ", r[9])
				}
				tx.GasPrice, err = strconv.ParseUint(r[10], 10, 64)
				if err != nil {
					log.Println("Error Parsing GasPrice : ", r[10])
				}
				mm.CsvTXs = append(mm.CsvTXs, tx)
				// Fill TXsByCategory
				if tx.Type == "DEPOSIT" {
					t := wallet.TX{Timestamp: tx.Timestamp, Note: tx.Hash + " : " + tx.From + " -> " + tx.To}
					t.Items = make(map[string][]wallet.Currency)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Coin, Amount: tx.Value})
					// t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: tx.Fees})
					mm.TXsByCategory["Deposits"] = append(mm.TXsByCategory["Deposits"], t)
				} else if tx.Type == "WITHDRAWAL" {
					t := wallet.TX{Timestamp: tx.Timestamp, Note: tx.Hash + " : " + tx.From + " -> " + tx.To}
					t.Items = make(map[string][]wallet.Currency)
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Coin, Amount: tx.Value})
					// t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Currency, Amount: tx.Fees})
					mm.TXsByCategory["Withdrawals"] = append(mm.TXsByCategory["Withdrawals"], t)
				} else {
					log.Println("Unmanaged ", tx)
				}
			}
		}
	}
	return
}
