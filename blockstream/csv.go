package blockstream

import (
	"encoding/csv"
	"io"
	"log"

	"github.com/shopspring/decimal"
)

type csvPayment struct {
	txID        string
	description string
	value       decimal.Decimal
	currency    string
}

type csvAddress struct {
	address     string
	description string
}

func (blkst *Blockstream) ParseCSVPayments(reader io.Reader) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "TxID" {
				a := csvPayment{}
				a.txID = r[0]
				a.description = r[1]
				a.value, err = decimal.NewFromString(r[2])
				if err != nil {
					log.Println("BTC Payments CSV Error Parsing Value : ", r[2])
				}
				a.currency = r[3]
				blkst.csvPayments = append(blkst.csvPayments, a)
			}
		}
	}
}

func (blkst *Blockstream) ParseCSVAddresses(reader io.Reader) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		for _, r := range records {
			if r[0] != "Address" {
				a := csvAddress{}
				a.address = r[0]
				a.description = r[1]
				blkst.csvAddresses = append(blkst.csvAddresses, a)
			}
		}
		// Fill TXsByCategory
		err = blkst.apiGetAllTXs()
	}
	blkst.done <- err
}
