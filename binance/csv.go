package binance

import (
	"encoding/csv"
	"io"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/utils"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type csvTX struct {
	Time      time.Time
	ID        string
	Account   string
	Operation string
	Coin      string
	Change    decimal.Decimal
	Fee       decimal.Decimal
	Remark    string
}

func (b *Binance) ParseCSV(reader io.Reader, extended bool, account string) (err error) {
	firstTimeUsed := time.Now()
	lastTimeUsed := time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC)
	const SOURCE = "Binance CSV :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		alreadyAsked := []string{}
		loc, _ := time.LoadLocation("Europe/Paris")
		for _, r := range records {
			if r[0] != "UTC_Time" {
				tx := csvTX{}
				tx.Time, err = time.ParseInLocation("2006-01-02 15:04:05", r[0], loc)
				if err != nil {
					log.Println(SOURCE, "Error Parsing Time", r[0])
				}
				tx.ID = utils.GetUniqueID(SOURCE + tx.Time.String())
				tx.Account = r[1]
				tx.Operation = r[2]
				tx.Coin = r[3]
				tx.Change, err = decimal.NewFromString(r[4])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Amount", r[4])
				}
				if extended {
					if r[5] != "" {
						tx.Fee, err = decimal.NewFromString(r[5])
						if err != nil {
							log.Println(SOURCE, "Error Parsing Fee", r[5])
						} else {
							if tx.Fee.IsNegative() {
								tx.Fee = tx.Fee.Neg()
							}
						}
					}
					tx.Remark = r[6]
				} else {
					tx.Remark = r[5]
				}
				b.csvTXs = append(b.csvTXs, tx)
				if tx.Time.Before(firstTimeUsed) {
					firstTimeUsed = tx.Time
				}
				if tx.Time.After(lastTimeUsed) {
					lastTimeUsed = tx.Time
				}
				// Fill TXsByCategory
				if tx.Operation == "Buy" ||
					tx.Operation == "Sell" ||
					tx.Operation == "Fee" ||
					tx.Operation == "Transaction Related" ||
					tx.Operation == "Small assets exchange BNB" ||
					tx.Operation == "The Easiest Way to Trade" {
					found := false
					for i, ex := range b.TXsByCategory["Exchanges"] {
						if ex.SimilarDate(time.Minute, tx.Time) {
							symbolsMatch := true
							if tx.Change.IsPositive() {
								if len(ex.Items["To"]) > 0 {
									if ex.Items["To"][0].Code != tx.Coin {
										symbolsMatch = false
									}
								}
							} else {
								if tx.Operation != "Fee" {
									if len(ex.Items["From"]) > 0 {
										if ex.Items["From"][0].Code != tx.Coin {
											symbolsMatch = false
										}
									}
								}
							}
							if symbolsMatch {
								found = true
								if b.TXsByCategory["Exchanges"][i].Items == nil {
									b.TXsByCategory["Exchanges"][i].Items = make(map[string]wallet.Currencies)
								}
								if tx.Change.IsPositive() {
									b.TXsByCategory["Exchanges"][i].Items["To"] = append(b.TXsByCategory["Exchanges"][i].Items["To"], wallet.Currency{Code: tx.Coin, Amount: tx.Change})
								} else {
									if tx.Operation == "Fee" {
										b.TXsByCategory["Exchanges"][i].Items["Fee"] = append(b.TXsByCategory["Exchanges"][i].Items["Fee"], wallet.Currency{Code: tx.Coin, Amount: tx.Change.Neg()})
									} else {
										b.TXsByCategory["Exchanges"][i].Items["From"] = append(b.TXsByCategory["Exchanges"][i].Items["From"], wallet.Currency{Code: tx.Coin, Amount: tx.Change.Neg()})
									}
								}
								if !tx.Fee.IsZero() {
									b.TXsByCategory["Exchanges"][i].Items["Fee"] = append(b.TXsByCategory["Exchanges"][i].Items["Fee"], wallet.Currency{Code: tx.Coin, Amount: tx.Fee})
								}
							}
						}
					}
					if !found {
						t := wallet.TX{Timestamp: tx.Time, ID: tx.ID, Note: "Binance CSV : Buy Sell Fee " + tx.Remark}
						t.Items = make(map[string]wallet.Currencies)
						if !tx.Fee.IsZero() {
							t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Coin, Amount: tx.Fee})
						}
						if tx.Change.IsPositive() {
							t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Coin, Amount: tx.Change})
							b.TXsByCategory["Exchanges"] = append(b.TXsByCategory["Exchanges"], t)
						} else {
							if tx.Operation == "Fee" {
								t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Coin, Amount: tx.Change.Neg()})
							} else {
								t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Coin, Amount: tx.Change.Neg()})
							}
							b.TXsByCategory["Exchanges"] = append(b.TXsByCategory["Exchanges"], t)
						}
					}
				} else if tx.Operation == "Deposit" ||
					tx.Operation == "transfer_in" ||
					tx.Operation == "Distribution" ||
					tx.Operation == "Super BNB Mining" ||
					tx.Operation == "POS savings interest" ||
					tx.Operation == "DeFi Staking Interest" ||
					tx.Operation == "Savings Interest" ||
					tx.Operation == "Launchpool Interest" ||
					tx.Operation == "Commission History" ||
					tx.Operation == "Commission Fee Shared With You" {
					t := wallet.TX{Timestamp: tx.Time, ID: tx.ID, Note: "Binance CSV : " + tx.Operation + " " + tx.Remark}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.Coin, Amount: tx.Change})
					if !tx.Fee.IsZero() {
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Coin, Amount: tx.Fee})
					}
					if tx.Operation == "Distribution" {
						b.TXsByCategory["AirDrops"] = append(b.TXsByCategory["AirDrops"], t)
					} else if tx.Operation == "POS savings interest" ||
						tx.Operation == "Super BNB Mining" ||
						tx.Operation == "DeFi Staking Interest" ||
						tx.Operation == "Launchpool Interest" {
						b.TXsByCategory["Minings"] = append(b.TXsByCategory["Minings"], t)
					} else if tx.Operation == "Savings Interest" {
						b.TXsByCategory["Interests"] = append(b.TXsByCategory["Interests"], t)
					} else if tx.Operation == "Commission History" ||
						tx.Operation == "Commission Fee Shared With You" {
						b.TXsByCategory["Referrals"] = append(b.TXsByCategory["Referrals"], t)
					} else {
						b.TXsByCategory["Deposits"] = append(b.TXsByCategory["Deposits"], t)
					}
				} else if tx.Operation == "Withdraw" ||
					tx.Operation == "transfer_out" {
					t := wallet.TX{Timestamp: tx.Time, ID: tx.ID, Note: "Binance CSV : " + tx.Operation + " " + tx.Remark}
					t.Items = make(map[string]wallet.Currencies)
					if tx.Fee.IsZero() {
						t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Coin, Amount: tx.Change.Neg()})
					} else {
						t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.Coin, Amount: tx.Change.Neg().Sub(tx.Fee)})
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: tx.Coin, Amount: tx.Fee})
					}
					b.TXsByCategory["Withdrawals"] = append(b.TXsByCategory["Withdrawals"], t)
				} else if tx.Operation == "POS savings purchase" ||
					tx.Operation == "POS savings redemption" ||
					tx.Operation == "Savings purchase" ||
					tx.Operation == "Savings Principal redemption" ||
					tx.Operation == "DeFi Staking purchase" ||
					tx.Operation == "DeFi Staking redemption" ||
					tx.Operation == "Liquid Swap add" ||
					tx.Operation == "Liquid Swap remove" {
					// Don't care
				} else {
					alreadyAsked = wallet.AskForHelp(SOURCE+" "+tx.Operation, tx, alreadyAsked)
				}
			}
		}
		b.Sources["Binance"] = source.Source{
			Crypto:        true,
			AccountNumber: account,
			OpeningDate:   firstTimeUsed,
			ClosingDate:   lastTimeUsed,
			LegalName:     "Binance Europe Services Limited",
			Address:       "LEVEL G (OFFICE 1/1235), QUANTUM HOUSE,75 ABATE RIGORD STREET, TA' XBIEXXBX 1120\nMalta",
			URL:           "https://www.binance.com/fr",
		}
	}
	return
}
