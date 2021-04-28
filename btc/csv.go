package btc

import (
	"encoding/csv"
	"io"
	"log"

	"github.com/shopspring/decimal"
)

type csvCategorie struct {
	txID        string
	kind        string
	description string
	value       decimal.Decimal
	currency    string
}

type CSVAddress struct {
	Address     string
	Description string
}

func (btc *BTC) ParseCSVCategorie(reader io.Reader) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "TxID" {
				a := csvCategorie{}
				a.txID = r[0]
				a.kind = r[1]
				a.description = r[2]
				a.value, err = decimal.NewFromString(r[3])
				if err != nil {
					log.Println("BTC Categorie CSV Error Parsing Value : ", r[3])
				}
				a.currency = r[4]
				btc.csvCategories = append(btc.csvCategories, a)
			}
		}
	}
}

func (btc *BTC) ParseCSVAddresses(reader io.Reader) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "Address" {
				a := CSVAddress{}
				a.Address = r[0]
				a.Description = r[1]
				btc.CSVAddresses = append(btc.CSVAddresses, a)
			}
		}
	}
}