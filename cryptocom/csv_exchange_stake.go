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

type csvStake struct {
	txsByCategory wallet.TXsByCategory
}

type csvExStakeTX struct {
	Time     time.Time
	ID       string
	Stake    wallet.Currency
	Apr      string
	Interest wallet.Currency
	Status   string
}

func (cdc *CryptoCom) ParseCSVExchangeStake(reader io.Reader) (err error) {
	const SOURCE = "Crypto.com Exchange Stake CSV :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "create_time_utc" {
				tx := csvExStakeTX{}
				tx.Time, err = time.Parse("2006-01-02 15:04:05.000", r[0])
				if err != nil {
					log.Println("Error Parsing Time : ", r[0])
				}
				tx.ID = utils.GetUniqueID(SOURCE + tx.Time.String())
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
				cdc.csvExStakeTXs = append(cdc.csvExStakeTXs, tx)
				t := wallet.TX{Timestamp: tx.Time, ID: tx.ID, Note: SOURCE + " " + tx.Stake.Amount.String() + " " + tx.Stake.Code + " " + tx.Apr}
				t.Items = make(map[string]wallet.Currencies)
				t.Items["To"] = append(t.Items["To"], tx.Interest)
				cdc.csvStake.txsByCategory["Interests"] = append(cdc.csvStake.txsByCategory["Interests"], t)
			}
		}
	}
	return
}
