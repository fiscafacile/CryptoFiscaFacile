package btc

import (
	"encoding/csv"
	"io"
)

func (btc *BTC) ParseCSVAddresses(reader io.Reader) (err error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "Address" {
				a := Address{}
				a.Address = r[0]
				a.Description = r[1]
				btc.Addresses = append(btc.Addresses, a)
			}
		}
	}
	return
}
