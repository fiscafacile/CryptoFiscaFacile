package monero

import (
	"encoding/csv"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/category"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type CsvTX struct {
	BlockHeight    string
	Epoch          time.Time
	Date           string
	Direction      string
	Amount         decimal.Decimal
	AtomicAmount   decimal.Decimal
	Fee            decimal.Decimal
	TxID           string
	Label          string
	SubaddrAccount string
	PaymentId      string
}

func (xmr *Monero) ParseCSV(reader io.Reader, cat category.Category) (err error) {
	const SOURCE = "Monero CSV :"
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err == nil {
		alreadyAsked := []string{}
		for _, r := range records {
			if r[0] != "blockHeight" {
				tx := CsvTX{}
				tx.BlockHeight = r[0]
				epoch, err := strconv.ParseInt(r[1], 10, 64)
				if err != nil {
					log.Println(SOURCE, "Error Parsing Epoch", r[1])
				} else {
					tx.Epoch = time.Unix(epoch, 0)
				}
				tx.Date = r[2]
				tx.Direction = r[3]
				tx.Amount, err = decimal.NewFromString(r[4])
				if err != nil {
					log.Println(SOURCE, "Error Parsing Amount", r[4])
				}
				atomic, err := strconv.ParseInt(r[5], 10, 64)
				if err != nil {
					log.Println(SOURCE, "Error Parsing AtomicAmount", r[5])
				} else {
					tx.AtomicAmount = decimal.New(atomic, -12)
				}
				if r[6] != "" {
					tx.Fee, err = decimal.NewFromString(r[6])
					if err != nil {
						log.Println(SOURCE, "Error Parsing Fee", r[6])
					}
				}
				tx.TxID = r[7]
				tx.SubaddrAccount = r[8]
				tx.PaymentId = r[9]
				xmr.CsvTXs = append(xmr.CsvTXs, tx)
			}
		}
		for _, tx := range xmr.CsvTXs {
			// Fixmr TXsByCategory
			if tx.Direction == "in" {
				t := wallet.TX{Timestamp: tx.Epoch, ID: tx.TxID, Note: SOURCE + " " + tx.BlockHeight + " " + tx.Label}
				t.Items = make(map[string]wallet.Currencies)
				t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "XMR", Amount: tx.AtomicAmount})
				if !tx.Fee.IsZero() {
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "XMR", Amount: tx.Fee})
				}
				if is, desc, val, curr := cat.IsTxCashIn(tx.TxID); is {
					t.Note += " crypto_purchase " + desc
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: curr, Amount: val})
					xmr.TXsByCategory["CashIn"] = append(xmr.TXsByCategory["CashIn"], t)
				} else {
					xmr.TXsByCategory["Deposits"] = append(xmr.TXsByCategory["Deposits"], t)
				}
			} else if tx.Direction == "out" {
				t := wallet.TX{Timestamp: tx.Epoch, ID: tx.TxID, Note: SOURCE + " " + tx.BlockHeight + " " + tx.Label}
				t.Items = make(map[string]wallet.Currencies)
				if !tx.Fee.IsZero() {
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "XMR", Amount: tx.Fee})
				}
				t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: "XMR", Amount: tx.AtomicAmount})
				if is, desc, val, curr := cat.IsTxCashOut(tx.TxID); is {
					t.Note += " crypto_payment " + desc
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr, Amount: val})
					xmr.TXsByCategory["CashOut"] = append(xmr.TXsByCategory["CashOut"], t)
				} else {
					xmr.TXsByCategory["Withdrawals"] = append(xmr.TXsByCategory["Withdrawals"], t)
				}
			} else {
				alreadyAsked = wallet.AskForHelp(SOURCE+" : "+tx.Direction, tx, alreadyAsked)
			}
		}
	}
	return
}
