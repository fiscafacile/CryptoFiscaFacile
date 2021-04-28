package etherscan

import (
	"encoding/csv"
	"io"
	"strings"

	"github.com/fiscafacile/CryptoFiscaFacile/category"
)

type csvAddress struct {
	address     string
	description string
}

func (ethsc *Etherscan) ParseCSV(reader io.Reader, cat category.Category) {
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
		err = ethsc.apiGetAllTXs(cat)
	}
	ethsc.done <- err
}
