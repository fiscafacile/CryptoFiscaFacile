package btc

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type BTC struct {
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
