package etherscan

import (
	"encoding/csv"
	"io"
	"strings"
)

func (ethsc *Etherscan) ParseCSVAddresses(reader io.Reader) (err error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "Address" {
				a := address{}
				a.address = strings.ToLower(r[0])
				a.description = r[1]
				ethsc.addresses = append(ethsc.addresses, a)
			}
		}
	}
	return
}
