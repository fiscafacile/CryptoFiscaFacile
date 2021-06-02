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

type csvTransfer struct {
	txsByCategory wallet.TXsByCategory
}

type csvExTransferTX struct {
	Time     time.Time
	ID       string
	Currency string
	Amount   decimal.Decimal
	Fee      decimal.Decimal
	Address  string
	Status   string
}

func (cdc *CryptoCom) ParseCSVExchangeTransfer(reader io.Reader) (err error) {
	const SOURCE = "Crypto.com Exchange Transfer CSV :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "create_time_utc" {
				tx := csvExTransferTX{}
				tx.Time, err = time.Parse("2006-01-02 15:04:05.000", r[0])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Time", r[0])
				}
				tx.ID = utils.GetUniqueID(SOURCE + tx.Time.String())
				tx.Currency = r[1]
				tx.Amount, err = decimal.NewFromString(r[2])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Amount", r[2])
				}
				tx.Fee, err = decimal.NewFromString(r[3])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Fee", r[3])
				}
				tx.Address = r[4]
				tx.Status = r[5]
				if tx.Address == "EARLY_SWAP_BONUS_DEPOSIT" ||
					tx.Address == "INTERNAL_DEPOSIT" {
					cdc.csvExTransferTXs = append(cdc.csvExTransferTXs, tx)
					t := wallet.TX{Timestamp: tx.Time, ID: tx.ID, Note: SOURCE + " " + tx.Address}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
					cdc.csvTransfer.txsByCategory["Deposits"] = append(cdc.csvTransfer.txsByCategory["Deposits"], t)
				} else {
					cdc.csvExTransferTXs = append(cdc.csvExTransferTXs, tx)
					t := wallet.TX{Timestamp: tx.Time, ID: tx.ID, Note: SOURCE + " " + tx.Address}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Currency, Amount: tx.Amount})
					cdc.csvTransfer.txsByCategory["Withdrawals"] = append(cdc.csvTransfer.txsByCategory["Withdrawals"], t)
				}
			}
		}
	}
	return
}
