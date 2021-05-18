package kraken

import (
	"encoding/csv"
	"io"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type csvTX struct {
	TxId    string
	RefId   string
	Time    time.Time
	Type    string
	SubType string
	Class   string
	Asset   string
	Amount  decimal.Decimal
	Fee     decimal.Decimal
	Balance decimal.Decimal
}

func (kr *Kraken) ParseCSV(reader io.Reader) (err error) {
	firstTimeUsed := time.Now()
	lastTimeUsed := time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC)
	const SOURCE = "Kraken CSV :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		alreadyAsked := []string{}
		for _, r := range records {
			// Ignore duplicate when no TxId
			if r[0] != "" && r[0] != "txid" {
				tx := csvTX{}
				tx.Time, err = time.Parse("2006-01-02 15:04:05", r[2])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Time", r[2])
				}
				tx.TxId = r[0]
				tx.RefId = r[1]
				tx.Type = r[3]
				tx.SubType = r[4]
				tx.Class = r[5]
				tx.Asset = ReplaceAssets(r[6])
				tx.Amount, err = decimal.NewFromString(r[7])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Amount", r[7])
				}
				tx.Fee, err = decimal.NewFromString(r[8])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Fee", r[8])
				}
				if tx.TxId == "" {
					tx.Balance, err = decimal.NewFromString(r[9])
					if err != nil {
						log.Println(SOURCE, "Error Parsing Balance", r[9])
					}
				} else {
					tx.Balance = decimal.NewFromInt(0)
				}
				kr.csvTXs = append(kr.csvTXs, tx)
				if tx.Time.Before(firstTimeUsed) {
					firstTimeUsed = tx.Time
				}
				if tx.Time.After(lastTimeUsed) {
					lastTimeUsed = tx.Time
				}
				// Fill TXsByCategory
				if tx.Type == "trade" {
					found := false
					for i, ex := range kr.TXsByCategory["Exchanges"] {
						if ex.SimilarDate(2*time.Second, tx.Time) {
							found = true
							if kr.TXsByCategory["Exchanges"][i].Items == nil {
								kr.TXsByCategory["Exchanges"][i].Items = make(map[string]wallet.Currencies)
							}
							if tx.Amount.IsPositive() {
								kr.TXsByCategory["Exchanges"][i].Items["To"] = append(kr.TXsByCategory["Exchanges"][i].Items["To"], wallet.Currency{Code: tx.Asset, Amount: tx.Amount})
							} else {
								kr.TXsByCategory["Exchanges"][i].Items["From"] = append(kr.TXsByCategory["Exchanges"][i].Items["From"], wallet.Currency{Code: tx.Asset, Amount: tx.Amount.Neg()})
							}
							if !tx.Fee.IsZero() {
								kr.TXsByCategory["Exchanges"][i].Items["Fee"] = append(kr.TXsByCategory["Exchanges"][i].Items["Fee"], wallet.Currency{Code: tx.Asset, Amount: tx.Fee})
							}
						}
					}
					if !found {
						t := wallet.TX{Timestamp: tx.Time}
						t.Items = make(map[string]wallet.Currencies)
						if !tx.Fee.IsZero() {
							t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Asset, Amount: tx.Fee})
						}
						if tx.Amount.IsPositive() {
							t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Asset, Amount: tx.Amount})
							kr.TXsByCategory["Exchanges"] = append(kr.TXsByCategory["Exchanges"], t)
						} else {
							t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Asset, Amount: tx.Amount.Neg()})
							kr.TXsByCategory["Exchanges"] = append(kr.TXsByCategory["Exchanges"], t)
						}
					}
				} else if tx.Type == "deposit" {
					t := wallet.TX{Timestamp: tx.Time}
					t.Items = make(map[string]wallet.Currencies)
					if !tx.Fee.IsZero() {
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Asset, Amount: tx.Fee})
					}
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Asset, Amount: tx.Amount})
					kr.TXsByCategory["Deposits"] = append(kr.TXsByCategory["Deposits"], t)
				} else if tx.Type == "withdrawal" {
					t := wallet.TX{Timestamp: tx.Time}
					t.Items = make(map[string]wallet.Currencies)
					if !tx.Fee.IsZero() {
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Asset, Amount: tx.Fee})
					}
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Asset, Amount: tx.Amount.Neg()})
					kr.TXsByCategory["Withdrawals"] = append(kr.TXsByCategory["Withdrawals"], t)
				} else if tx.Type == "staking" {
					t := wallet.TX{Timestamp: tx.Time}
					t.Items = make(map[string]wallet.Currencies)
					if !tx.Fee.IsZero() {
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Asset, Amount: tx.Fee})
					}
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Asset, Amount: tx.Amount})
					kr.TXsByCategory["Interests"] = append(kr.TXsByCategory["Interests"], t)
				} else if tx.Type == "transfer" {
					// Ignore transfer because it's a intra-account transfert
					// is there some Fees to consider ?
				} else {
					alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Type, tx, alreadyAsked)
				}
			}
		}
	}
	return
}
