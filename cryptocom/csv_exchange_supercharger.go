package cryptocom

import (
	"encoding/csv"
	"io"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/utils"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type csvSupercharger struct {
	txsByCategory wallet.TXsByCategory
}

type csvExSuperchargerTX struct {
	Time        time.Time
	ID          string
	Currency    string
	Amount      decimal.Decimal
	Description string
}

func (cdc *CryptoCom) ParseCSVExchangeSupercharger(reader io.Reader) (err error) {
	const SOURCE = "Crypto.com Exchange SuperCharger CSV :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "create_time_utc" {
				tx := csvExSuperchargerTX{}
				tx.Time, err = time.Parse("2006-01-02 15:04:05", r[0])
				if err != nil {
					log.Println("Error Parsing Time : ", r[0])
				}
				tx.ID = utils.GetUniqueID(SOURCE + tx.Time.String())
				tx.Currency = r[1]
				tx.Amount, err = decimal.NewFromString(r[2])
				if err != nil {
					log.Println("Error Parsing Amount : ", r[2])
				}
				tx.Description = r[3]
				cdc.csvExSuperchargerTXs = append(cdc.csvExSuperchargerTXs, tx)
				t := wallet.TX{Timestamp: tx.Time, ID: tx.ID, Note: SOURCE + " " + tx.Description}
				t.Items = make(map[string]wallet.Currencies)
				t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
				cdc.csvSupercharger.txsByCategory["Minings"] = append(cdc.csvSupercharger.txsByCategory["Minings"], t)
			}
		}
	}
	return
}
