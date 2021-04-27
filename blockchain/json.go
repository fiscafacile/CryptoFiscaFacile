package blockchain

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type Wallet struct {
	Address string          `json:address`
	Amount  decimal.Decimal `json:amount`
}

type Wallets []Wallet

type JsonTX struct {
	TxID string          `json:txid`
	Date string          `json:date`
	Fee  decimal.Decimal `json:fee,omitempty`
	From Wallets         `json:from,omitempty`
	To   Wallets         `json:to,omitempty`
}

func (bc *BlockChain) ParseTXsJSON(reader io.Reader, currency string) (err error) {
	var txs []JsonTX
	jsonDecoder := json.NewDecoder(reader)
	err = jsonDecoder.Decode(&txs)
	if err == nil {
		bc.jsonTXs = append(bc.jsonTXs, txs...)
		// Fill TXsByCategory
		for _, tx := range txs {
			date, err := time.Parse("Jan 2, 2006 15:04:05 PM", tx.Date)
			if err != nil {
				log.Println("BlockChain JSON :", err)
			}
			t := wallet.TX{Timestamp: date, Note: "BlockChain JSON : " + tx.TxID}
			t.Items = make(map[string]wallet.Currencies)
			if !tx.Fee.IsZero() {
				t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: currency, Amount: tx.Fee})
			}
			haveFrom := false
			for _, w := range tx.From {
				haveFrom = true
				t.Note += " " + w.Address
				t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: currency, Amount: w.Amount})
			}
			t.Note += " ->"
			haveTo := false
			for _, w := range tx.To {
				haveTo = true
				t.Note += " " + w.Address
				t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: currency, Amount: w.Amount})
			}
			if haveFrom && haveTo {
				bc.TXsByCategory["Transfers"] = append(bc.TXsByCategory["Transfers"], t)
			} else if haveFrom {
				bc.TXsByCategory["Withdrawals"] = append(bc.TXsByCategory["Withdrawals"], t)
			} else if haveTo {
				bc.TXsByCategory["Deposits"] = append(bc.TXsByCategory["Deposits"], t)
			} else {
				bc.TXsByCategory["Fees"] = append(bc.TXsByCategory["Fees"], t)
			}
		}
	}
	return
}
