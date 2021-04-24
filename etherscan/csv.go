package etherscan

import (
	"encoding/csv"
	"io"
	"strings"
)

type csvAddress struct {
	address     string
	description string
}

func (ethsc *Etherscan) ParseCSV(reader io.Reader) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "Address" {
				a := csvAddress{}
				a.address = strings.ToLower(r[0])
				a.description = r[1]
				ethsc.csvAddresses = append(ethsc.csvAddresses, a)
			}
		}
		// Fill TXsByCategory
		err = ethsc.apiGetAllTXs()
	}
	ethsc.done <- err
}
