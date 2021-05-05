package hitbtc

import (
	"encoding/csv"
	"io"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type csvTradeTX struct {
	Email      string
	Date       time.Time
	Instrument string
	TradeID    string
	OrderID    string
	Side       string
	Quantity   decimal.Decimal
	Price      string
	Volume     string
	Fee        string
	Rebate     string
	Total      string
	Taker      string
}

func (hb *HitBTC) ParseCSVTrades(reader io.Reader) (err error) {
	const SOURCE = "HitBTC CSV Trades :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		alreadyAsked := []string{}
		for _, r := range records {
			if r[0] != "Email" {
				tx := csvTradeTX{}
				tx.Email = r[0]
				tx.Date, err = time.Parse("2006-01-02 15:04:05", r[1])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Date", r[1])
				}
				tx.Instrument = r[2]
				tx.TradeID = r[3]
				tx.OrderID = r[4]
				tx.Side = r[5]
				tx.Quantity, err = decimal.NewFromString(r[6])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Quantity", r[6])
				}
				tx.Price = r[7]
				tx.Volume = r[8]
				tx.Fee = r[9]
				tx.Rebate = r[10]
				tx.Total = r[11]
				tx.Taker = r[12]
				hb.csvTradeTXs = append(hb.csvTradeTXs, tx)
				// Fill TXsByCategory
				alreadyAsked = wallet.AskForHelp(SOURCE, tx, alreadyAsked)
			}
		}
	}
	return
}
