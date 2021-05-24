package bittrex

import (
	"encoding/csv"
	"io"
	"log"
	"strings"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/category"
	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type csvTX struct {
	ID          string
	FromSymbol  string
	ToSymbol    string
	Time        time.Time
	Operation   string
	FromAmount  decimal.Decimal
	ToAmount    decimal.Decimal
	Fee         decimal.Decimal
	FeeCurrency string
	Remark      string
}

func (btrx *Bittrex) ParseCSV(reader io.Reader, cat category.Category, account string) (err error) {
	firstTimeUsed := time.Now()
	lastTimeUsed := time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC)
	const SOURCE = "Bittrex CSV :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	opeRplcr := strings.NewReplacer(
		"LIMIT_SELL", "SELL",
		"MARKET_SELL", "SELL",
		"CEILING_MARKET_BUY", "BUY",
		"LIMIT_BUY", "BUY",
		"MARKET_BUY", "BUY",
	)
	symRplcr := strings.NewReplacer(
		"REPV2", "REP",
	)
	if err == nil {
		alreadyAsked := []string{}
		for _, r := range records {
			if r[0] != "Uuid" {
				tx := csvTX{}
				tx.ID = r[0]
				symbolSlice := strings.Split(r[1], "-")
				tx.Time, err = time.Parse("1/2/2006 3:04:05 PM", r[14])
				if err != nil {
					log.Println("Error Parsing Time : ", r[14])
				}
				tx.Operation = opeRplcr.Replace(r[3])
				quantity, err := decimal.NewFromString(r[5])
				if err != nil {
					log.Println(SOURCE, "Error Parsing quantity", r[5])
				}
				quantityRemaining, err := decimal.NewFromString(r[6])
				if err != nil {
					log.Println(SOURCE, "Error Parsing quantityRemaining", r[6])
				}
				tx.Fee, err = decimal.NewFromString(r[7])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Fee", r[7])
				}
				tx.FeeCurrency = symRplcr.Replace(symbolSlice[0])
				price, err := decimal.NewFromString(r[8])
				if err != nil {
					log.Println(SOURCE, "Error Parsing price", r[8])
				}
				if tx.Time.Before(firstTimeUsed) {
					firstTimeUsed = tx.Time
				}
				if tx.Time.After(lastTimeUsed) {
					lastTimeUsed = tx.Time
				}
				// Fill TXsByCategory
				if tx.Operation == "BUY" || tx.Operation == "SELL" {
					if tx.Operation == "BUY" {
						tx.FromSymbol = symRplcr.Replace(symbolSlice[0])
						tx.FromAmount = price
						tx.ToSymbol = symRplcr.Replace(symbolSlice[1])
						tx.ToAmount = quantity.Sub(quantityRemaining)
					} else if tx.Operation == "SELL" {
						tx.FromSymbol = symRplcr.Replace(symbolSlice[1])
						tx.FromAmount = quantity.Sub(quantityRemaining)
						tx.ToSymbol = symRplcr.Replace(symbolSlice[0])
						tx.ToAmount = price
					}
					to := wallet.Currency{Code: tx.ToSymbol, Amount: tx.ToAmount}
					from := wallet.Currency{Code: tx.FromSymbol, Amount: tx.FromAmount}
					t := wallet.TX{Timestamp: tx.Time, ID: tx.ID, Note: "Bittrex CSV : " + tx.Operation}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["From"] = append(t.Items["From"], from)
					t.Items["To"] = append(t.Items["To"], to)
					if !tx.Fee.IsZero() {
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.FeeCurrency, Amount: tx.Fee})
					}
					if is, desc, val, curr := cat.IsTxShit(tx.ID); is {
						t.Note += " " + desc
						t.Items["Lost"] = append(t.Items["Lost"], wallet.Currency{Code: curr, Amount: val})
					}
					if to.IsFiat() && from.IsFiat() {
						//ignore
					} else if to.IsFiat() {
						btrx.TXsByCategory["CashOut"] = append(btrx.TXsByCategory["CashOut"], t)
					} else if from.IsFiat() {
						btrx.TXsByCategory["CashIn"] = append(btrx.TXsByCategory["CashIn"], t)
					} else {
						btrx.TXsByCategory["Exchanges"] = append(btrx.TXsByCategory["Exchanges"], t)
					}
				} else {
					alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Operation, tx, alreadyAsked)
				}
			}
		}
	}
	if _, ok := btrx.Sources["Bittrex"]; ok {
		if btrx.Sources["Bittrex"].OpeningDate.After(firstTimeUsed) {
			src := btrx.Sources["Bittrex"]
			src.OpeningDate = firstTimeUsed
			btrx.Sources["Bittrex"] = src
		}
		if btrx.Sources["Bittrex"].ClosingDate.Before(lastTimeUsed) {
			src := btrx.Sources["Bittrex"]
			src.ClosingDate = lastTimeUsed
			btrx.Sources["Bittrex"] = src
		}
	} else {
		btrx.Sources["Bittrex"] = source.Source{
			Crypto:        true,
			AccountNumber: account,
			OpeningDate:   firstTimeUsed,
			ClosingDate:   lastTimeUsed,
			LegalName:     "Bittrex International GmbH",
			Address:       "Dr. Grass-Strasse 12, 9490 Vaduz,\nPrincipality of Liechtenstein",
			URL:           "https://global.bittrex.com",
		}
	}
	return
}
