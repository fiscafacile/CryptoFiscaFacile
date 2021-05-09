package hitbtc

import (
	"encoding/csv"
	"io"
	"log"
	"strings"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/utils"
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
	Price      decimal.Decimal
	Volume     decimal.Decimal
	Fee        decimal.Decimal
	Rebate     string
	Total      string
	Taker      string
}

func (hb *HitBTC) ParseCSVTrades(reader io.Reader) (err error) {
	firstTimeUsed := time.Now()
	lastTimeUsed := time.Date(2019, time.November, 14, 0, 0, 0, 0, time.UTC)
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
				if r[7] != "" {
					tx.Price, err = decimal.NewFromString(r[7])
					if err != nil {
						log.Println(SOURCE, "Error Parsing Price", r[7])
					}
				}
				if r[8] != "" {
					tx.Volume, err = decimal.NewFromString(r[8])
					if err != nil {
						log.Println(SOURCE, "Error Parsing Volume", r[8])
					}
				}
				if r[9] != "" {
					tx.Fee, err = decimal.NewFromString(r[9])
					if err != nil {
						log.Println(SOURCE, "Error Parsing Fee", r[9])
					}
				}
				tx.Rebate = r[10]
				tx.Total = r[11]
				tx.Taker = r[12]
				hb.csvTradeTXs = append(hb.csvTradeTXs, tx)
				hb.emails = utils.AppendUniq(hb.emails, tx.Email)
				// Fill TXsByCategory
				t := wallet.TX{Timestamp: tx.Date, ID: tx.TradeID, Note: SOURCE + " " + tx.Instrument + " " + tx.OrderID}
				t.Items = make(map[string]wallet.Currencies)
				curr := strings.Split(tx.Instrument, "_")
				for i, c := range curr {
					curr[i] = csvCurrencyCure(c)
				}
				if tx.Volume.IsZero() {
					tx.Volume = tx.Quantity.Mul(tx.Price)
				}
				if !tx.Fee.IsZero() {
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: curr[0], Amount: tx.Fee})
				}
				if tx.Side == "sell" {
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: curr[0], Amount: tx.Quantity})
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr[1], Amount: tx.Volume})
					hb.TXsByCategory["Exchanges"] = append(hb.TXsByCategory["Exchanges"], t)
				} else if tx.Side == "buy" {
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: curr[1], Amount: tx.Volume})
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr[0], Amount: tx.Quantity})
					hb.TXsByCategory["Exchanges"] = append(hb.TXsByCategory["Exchanges"], t)
				} else {
					alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Side, tx, alreadyAsked)
				}
				if tx.Date.Before(firstTimeUsed) {
					firstTimeUsed = tx.Date
				}
				if tx.Date.After(lastTimeUsed) {
					lastTimeUsed = tx.Date
				}
			}
		}
	}
	for _, e := range hb.emails {
		if _, ok := hb.Sources["HitBTC_"+e]; !ok {
			hb.Sources["HitBTC_"+e] = source.Source{
				Crypto:        true,
				AccountNumber: utils.RemoveSymbol(e),
				OpeningDate:   firstTimeUsed,
				ClosingDate:   lastTimeUsed,
				LegalName:     "Hit Tech Solutions Development Ltd.",
				Address:       "Suite 15, Oliaji Trade Centre, Francis Rachel Street,\nVictoria, Mahe,\nSeychelles",
				URL:           "https://hitbtc.com",
			}
		}
	}
	return
}
