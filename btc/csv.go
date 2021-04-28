package btc

import (
	"encoding/csv"
	"io"
)

type CSVAddress struct {
	Address     string
	Description string
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
