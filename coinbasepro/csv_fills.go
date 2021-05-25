package coinbasepro

import (
	"encoding/csv"
	"io"
	"log"
	"strings"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type CsvFillsTX struct {
	Portfolio         string
	TradeID           string
	ProductLeft       string
	ProductRight      string
	Side              string
	CreatedAt         time.Time
	Size              decimal.Decimal
	SizeUnit          string
	Price             decimal.Decimal
	Fee               decimal.Decimal
	Total             decimal.Decimal
	PriceFeeTotalUnit string
}

func (cbp *CoinbasePro) ParseFillsCSV(reader io.ReadSeeker, account string) (err error) {
	firstTimeUsed := time.Now()
	lastTimeUsed := time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC)
	const SOURCE = "CoinbasePro Fills CSV :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		alreadyAsked := []string{}
		for _, r := range records {
			if r[0] != "portfolio" {
				tx := CsvFillsTX{}
				tx.Portfolio = r[0]
				tx.TradeID = r[1]
				products := strings.Split(r[2], "-")
				tx.ProductLeft = products[0]
				tx.ProductRight = products[1]
				tx.Side = r[3]
				tx.CreatedAt, err = time.Parse("2006-01-02T15:04:05.999Z", r[4])
				if err != nil {
					log.Println(SOURCE, "Error Parsing CreatedAt : ", r[4])
				}
				tx.Size, err = decimal.NewFromString(r[5])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Size : ", r[5])
				}
				tx.SizeUnit = r[6]
				tx.Price, err = decimal.NewFromString(r[7])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Price : ", r[7])
				}
				if r[8] != "" {
					tx.Fee, err = decimal.NewFromString(r[8])
					if err != nil {
						log.Println(SOURCE, "Error Parsing Fee : ", r[8])
					}
				}
				if r[9] != "" {
					tx.Total, err = decimal.NewFromString(r[9])
					if err != nil {
						log.Println(SOURCE, "Error Parsing Total : ", r[9])
					}
				}
				tx.PriceFeeTotalUnit = r[10]
				cbp.CsvFillsTXs = append(cbp.CsvFillsTXs, tx)
				if tx.CreatedAt.Before(firstTimeUsed) {
					firstTimeUsed = tx.CreatedAt
				}
				if tx.CreatedAt.After(lastTimeUsed) {
					lastTimeUsed = tx.CreatedAt
				}
				// Fill TXsByCategory
				if tx.Side == "BUY" {
					t := wallet.TX{Timestamp: tx.CreatedAt, ID: tx.TradeID, Note: SOURCE + " portfolio " + tx.Portfolio}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.ProductLeft, Amount: tx.Size})
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.ProductRight, Amount: tx.Size.Mul(tx.Price)})
					if !tx.Fee.IsZero() {
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.PriceFeeTotalUnit, Amount: tx.Fee})
					}
					cbp.TXsByCategory["Exchanges"] = append(cbp.TXsByCategory["Exchanges"], t)
				} else {
					alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Side, tx, alreadyAsked)
				}
			}
		}
	}
	if _, ok := cbp.Sources["CoinbasePro"]; ok {
		if cbp.Sources["CoinbasePro"].OpeningDate.After(firstTimeUsed) {
			src := cbp.Sources["CoinbasePro"]
			src.OpeningDate = firstTimeUsed
			cbp.Sources["CoinbasePro"] = src
		}
		if cbp.Sources["CoinbasePro"].ClosingDate.Before(lastTimeUsed) {
			src := cbp.Sources["CoinbasePro"]
			src.ClosingDate = lastTimeUsed
			cbp.Sources["CoinbasePro"] = src
		}
	} else {
		cbp.Sources["CoinbasePro"] = source.Source{
			Crypto:        true,
			AccountNumber: account,
			OpeningDate:   firstTimeUsed,
			ClosingDate:   lastTimeUsed,
			LegalName:     "Coinbase Europe Limited",
			Address:       "70 Sir John Rogerson’s Quay,\nDublin D02 R296\nIrlande",
			URL:           "https://pro.coinbase.com",
		}
	}
	return
}
