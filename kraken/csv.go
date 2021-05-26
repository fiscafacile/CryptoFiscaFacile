package kraken

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

func (kr *Kraken) ParseCSV(reader io.Reader, cat category.Category, account string) (err error) {
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
				if tx.Type == "trade" ||
					tx.Type == "spend" ||
					tx.Type == "margin" ||
					tx.Type == "settled" ||
					tx.Type == "receive" {
					found := false
					for i, ex := range kr.TXsByCategory["Exchanges"] {
						if strings.Contains(ex.ID, tx.RefId) {
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
						t := wallet.TX{Timestamp: tx.Time, ID: tx.TxId + "-" + tx.RefId, Note: SOURCE + " " + tx.Type}
						t.Items = make(map[string]wallet.Currencies)
						if is, desc, val, curr := cat.IsTxShit(t.ID); is {
							t.Note += " " + desc
							t.Items["Lost"] = append(t.Items["Lost"], wallet.Currency{Code: curr, Amount: val})
						}
						if !tx.Fee.IsZero() {
							t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Asset, Amount: tx.Fee})
						}
						if tx.Amount.IsPositive() {
							t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Asset, Amount: tx.Amount})
						} else {
							t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Asset, Amount: tx.Amount.Neg()})
						}
						kr.TXsByCategory["Exchanges"] = append(kr.TXsByCategory["Exchanges"], t)
					}
				} else if tx.Type == "deposit" {
					t := wallet.TX{Timestamp: tx.Time, ID: tx.TxId + "-" + tx.RefId, Note: SOURCE + " " + tx.Type}
					t.Items = make(map[string]wallet.Currencies)
					if is, desc, val, curr := cat.IsTxShit(t.ID); is {
						t.Note += " " + desc
						t.Items["Lost"] = append(t.Items["Lost"], wallet.Currency{Code: curr, Amount: val})
					}
					if !tx.Fee.IsZero() {
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Asset, Amount: tx.Fee})
					}
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Asset, Amount: tx.Amount})
					kr.TXsByCategory["Deposits"] = append(kr.TXsByCategory["Deposits"], t)
				} else if tx.Type == "withdrawal" {
					t := wallet.TX{Timestamp: tx.Time, ID: tx.TxId + "-" + tx.RefId, Note: SOURCE + " " + tx.Type}
					t.Items = make(map[string]wallet.Currencies)
					if is, desc, val, curr := cat.IsTxShit(t.ID); is {
						t.Note += " " + desc
						t.Items["Lost"] = append(t.Items["Lost"], wallet.Currency{Code: curr, Amount: val})
					}
					if !tx.Fee.IsZero() {
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Asset, Amount: tx.Fee})
					}
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Asset, Amount: tx.Amount.Neg()})
					if is, desc := cat.IsTxGift(t.ID); is {
						t.Note += " gift " + desc
						kr.TXsByCategory["Gifts"] = append(kr.TXsByCategory["Gifts"], t)
					} else if is, desc, val, curr := cat.IsTxExchange(t.ID); is {
						t.Note += " " + desc
						t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr, Amount: val})
						kr.TXsByCategory["Exchanges"] = append(kr.TXsByCategory["Exchanges"], t)
					} else if is, desc, val, curr := cat.IsTxCashOut(t.ID); is {
						t.Note += " crypto_payment " + desc
						t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr, Amount: val})
						kr.TXsByCategory["CashOut"] = append(kr.TXsByCategory["CashOut"], t)
					} else {
						kr.TXsByCategory["Withdrawals"] = append(kr.TXsByCategory["Withdrawals"], t)
					}
				} else if tx.Type == "staking" {
					t := wallet.TX{Timestamp: tx.Time, ID: tx.TxId + "-" + tx.RefId, Note: SOURCE + " " + tx.Type}
					t.Items = make(map[string]wallet.Currencies)
					if is, desc, val, curr := cat.IsTxShit(t.ID); is {
						t.Note += " " + desc
						t.Items["Lost"] = append(t.Items["Lost"], wallet.Currency{Code: curr, Amount: val})
					}
					if !tx.Fee.IsZero() {
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Asset, Amount: tx.Fee})
					}
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Asset, Amount: tx.Amount})
					kr.TXsByCategory["Interests"] = append(kr.TXsByCategory["Interests"], t)
				} else if tx.Type == "rollover" {
					fee := wallet.Currency{Code: tx.Asset, Amount: tx.Fee}
					if !fee.IsFiat() {
						t := wallet.TX{Timestamp: tx.Time, ID: tx.TxId + "-" + tx.RefId, Note: SOURCE + " " + tx.Type}
						t.Items = make(map[string]wallet.Currencies)
						if is, desc, val, curr := cat.IsTxShit(tx.TxId); is {
							t.Note += " " + desc
							t.Items["Lost"] = append(t.Items["Lost"], wallet.Currency{Code: curr, Amount: val})
						}
						t.Items["Fee"] = append(t.Items["Fee"], fee)
						kr.TXsByCategory["Fees"] = append(kr.TXsByCategory["Fees"], t)
					}
				} else if tx.Type == "transfer" {
					t := wallet.TX{Timestamp: tx.Time, ID: tx.TxId + "-" + tx.RefId, Note: SOURCE + " " + tx.Type}
					t.Items = make(map[string]wallet.Currencies)
					if is, desc, val, curr := cat.IsTxShit(tx.TxId); is {
						t.Note += " " + desc
						t.Items["Lost"] = append(t.Items["Lost"], wallet.Currency{Code: curr, Amount: val})
					}
					if tx.SubType == "" {
						if tx.Amount.IsPositive() {
							t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Asset, Amount: tx.Amount})
						} else {
							t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Asset, Amount: tx.Amount.Neg()})
						}
						kr.TXsByCategory["AirDrops"] = append(kr.TXsByCategory["AirDrops"], t)
					} else {
						// Ignore non void subType transfer because it's a intra-account transfert
						if !tx.Fee.IsZero() {
							t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Asset, Amount: tx.Fee})
							kr.TXsByCategory["Fees"] = append(kr.TXsByCategory["Fees"], t)
						}
					}
				} else {
					alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Type, tx, alreadyAsked)
				}
			}
		}
		kr.Sources["Kraken"] = source.Source{
			Crypto:        true,
			AccountNumber: account,
			OpeningDate:   kr.api.firstTimeUsed,
			ClosingDate:   kr.api.lastTimeUsed,
			LegalName:     "Payward Ltd.",
			Address:       "6th Floor,\nOne London Wall,\nLondon, EC2Y 5EB,\nRoyaume-Uni",
			URL:           "https://www.kraken.com",
		}
	}
	return
}
