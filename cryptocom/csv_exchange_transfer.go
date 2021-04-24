package cryptocom

import (
	"encoding/csv"
	"io"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type CsvTXExTransfer struct {
	Time     time.Time
	Currency string
	Amount   decimal.Decimal
	Fee      decimal.Decimal
	Address  string
	Status   string
}

func (cdc *CryptoCom) ParseCSVExTransfer(reader io.Reader) (err error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "create_time_utc" {
				tx := CsvTXExTransfer{}
				tx.Time, err = time.Parse("2006-01-02 15:04:05.000", r[0])
				if err != nil {
					log.Println("Error Parsing Time : ", r[0])
				}
				tx.Currency = r[1]
				tx.Amount, err = decimal.NewFromString(r[2])
				if err != nil {
					log.Println("Error Parsing Amount : ", r[2])
				}
				tx.Fee, err = decimal.NewFromString(r[3])
				if err != nil {
					log.Println("Error Parsing Fee : ", r[3])
				}
				tx.Address = r[4]
				tx.Status = r[5]
				if tx.Address == "EARLY_SWAP_BONUS_DEPOSIT" ||
					tx.Address == "INTERNAL_DEPOSIT" {
					cdc.CsvTXsExTransfer = append(cdc.CsvTXsExTransfer, tx)
					t := wallet.TX{Timestamp: tx.Time, Note: "Crypto.com Exchange Transfer CSV : " + tx.Address}
					t.Items = make(map[string][]wallet.Currency)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
					cdc.TXsByCategory["Deposits"] = append(cdc.TXsByCategory["Deposits"], t)
				} else {
					cdc.CsvTXsExTransfer = append(cdc.CsvTXsExTransfer, tx)
					t := wallet.TX{Timestamp: tx.Time, Note: "Crypto.com Exchange Transfer CSV : " + tx.Address}
					t.Items = make(map[string][]wallet.Currency)
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
					cdc.TXsByCategory["Withdrawals"] = append(cdc.TXsByCategory["Withdrawals"], t)
				}
			}
		}
	}
	return
}
