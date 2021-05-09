package category

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

func (cat *Category) ParseCSVCategory(reader io.Reader) {
	const SOURCE = "TXs Categorie CSV :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "TxID" {
				a := csvCategorie{}
				a.txID = r[0]
				a.kind = r[1]
				a.description = r[2]
				if r[3] != "" {
					a.value, err = decimal.NewFromString(r[3])
					if err != nil {
						log.Println(SOURCE, "Error Parsing Value", r[3])
					}
				}
				a.currency = r[4]
				cat.csvCategories = append(cat.csvCategories, a)
			}
		}
	}
}
