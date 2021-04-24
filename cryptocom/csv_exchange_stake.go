package cryptocom

import (
	"encoding/csv"
	"io"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type CsvTXExStake struct {
	Time     time.Time
	Stake    wallet.Currency
	Apr      string
	Interest wallet.Currency
	Status   string
}

func (cdc *CryptoCom) ParseCSVExStake(reader io.Reader) (err error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "create_time_utc" {
				tx := CsvTXExStake{}
				tx.Time, err = time.Parse("2006-01-02 15:04:05.000", r[0])
				if err != nil {
					log.Println("Error Parsing Time : ", r[0])
				}
				tx.Stake.Code = r[1]
				tx.Stake.Amount, err = decimal.NewFromString(r[2])
				if err != nil {
					log.Println("Error Parsing Stake.Amount : ", r[2])
				}
				tx.Apr = r[3]
				tx.Interest.Code = r[4]
				tx.Interest.Amount, err = decimal.NewFromString(r[5])
				if err != nil {
					log.Println("Error Parsing Interest.Amount : ", r[5])
				}
				tx.Status = r[6]
				cdc.CsvTXsExStake = append(cdc.CsvTXsExStake, tx)
				t := wallet.TX{Timestamp: tx.Time, Note: "Crypto.com Exchange Stake CSV : " + tx.Stake.Amount.String() + " " + tx.Stake.Code + " " + tx.Apr}
				t.Items = make(map[string][]wallet.Currency)
				t.Items["To"] = append(t.Items["To"], tx.Interest)
				cdc.TXsByCategory["Deposits"] = append(cdc.TXsByCategory["Deposits"], t)
			}
		}
	}
	return
}
