package btc

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type BTC struct {
	csvCategories []csvCategorie
	CSVAddresses  []CSVAddress
	TXsByCategory wallet.TXsByCategory
}

func New() *BTC {
	btc := &BTC{}
	btc.TXsByCategory = make(map[string]wallet.TXs)
	return btc
}

func (btc BTC) OwnAddress(add string) bool {
	for _, a := range btc.CSVAddresses {
		if a.Address == add {
			return true
		}
	}
	return false
}

func (btc BTC) IsTxCashOut(txid string) (is bool, desc string, val decimal.Decimal, curr string) {
	is = false
	for _, a := range btc.csvCategories {
		if a.txID == txid && a.kind == "OUT" {
			is = true
			desc = a.description
			val = a.value
			curr = a.currency
			return
		}
	}
	return
}

func (btc BTC) IsTxCashIn(txid string) (is bool, desc string, val decimal.Decimal, curr string) {
	is = false
	for _, a := range btc.csvCategories {
		if a.txID == txid && a.kind == "IN" {
			is = true
			desc = a.description
			val = a.value
			curr = a.currency
			return
		}
	}
	return
}

func (btc BTC) IsTxExchange(txid string) (is bool, desc string, val decimal.Decimal, curr string) {
	is = false
	for _, a := range btc.csvCategories {
		if a.txID == txid && a.kind == "EXC" {
			is = true
			desc = a.description
			val = a.value
			curr = a.currency
			return
		}
	}
	return
}

func (btc BTC) HasCustody(txid string) (is bool, desc string, val decimal.Decimal) {
	is = false
	for _, a := range btc.csvCategories {
		if a.txID == txid && a.kind == "CUS" {
			is = true
			desc = a.description
			val = a.value
			return
		}
	}
	return
}

func (btc BTC) IsTxGift(txid string) (is bool, desc string) {
	is = false
	for _, a := range btc.csvCategories {
		if a.txID == txid && a.kind == "GIFT" {
			is = true
			desc = a.description
			return
		}
	}
	return
}
