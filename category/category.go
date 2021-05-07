package category

import (
	"github.com/shopspring/decimal"
)

type Category struct {
	csvCategories []csvCategorie
}

func New() *Category {
	cat := &Category{}
	return cat
}

func (cat Category) IsTxCashOut(txid string) (is bool, desc string, val decimal.Decimal, curr string) {
	is = false
	for _, a := range cat.csvCategories {
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

func (cat Category) IsTxCashIn(txid string) (is bool, desc string, val decimal.Decimal, curr string) {
	is = false
	for _, a := range cat.csvCategories {
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

func (cat Category) IsTxExchange(txid string) (is bool, desc string, val decimal.Decimal, curr string) {
	is = false
	for _, a := range cat.csvCategories {
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

func (cat Category) HasCustody(txid string) (is bool, desc string, val decimal.Decimal) {
	is = false
	for _, a := range cat.csvCategories {
		if a.txID == txid && a.kind == "CUS" {
			is = true
			desc = a.description
			val = a.value
			return
		}
	}
	return
}

func (cat Category) IsTxGift(txid string) (is bool, desc string) {
	is = false
	for _, a := range cat.csvCategories {
		if a.txID == txid && a.kind == "GIFT" {
			is = true
			desc = a.description
			return
		}
	}
	return
}

func (cat Category) IsTxAirDrop(txid string) (is bool, desc string) {
	is = false
	for _, a := range cat.csvCategories {
		if a.txID == txid && a.kind == "AIR" {
			is = true
			desc = a.description
			return
		}
	}
	return
}

func (cat Category) IsTxShit(txid string) (is bool, desc string) {
	is = false
	for _, a := range cat.csvCategories {
		if a.txID == txid && a.kind == "SHIT" {
			is = true
			desc = a.description
			return
		}
	}
	return
}

func (cat Category) IsTxTokenSale(txid string) (is bool, buy string) {
	is = false
	for _, a := range cat.csvCategories {
		if a.txID == txid && a.kind == "TOK" {
			is = true
			buy = a.description
			return
		}
	}
	return
}

func (cat Category) IsTxFee(txid string) (is bool, fee string) {
	is = false
	for _, a := range cat.csvCategories {
		if a.txID == txid && a.kind == "FEE" {
			is = true
			fee = a.description
			return
		}
	}
	return
}
