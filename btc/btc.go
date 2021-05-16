package btc

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Address struct {
	Address     string
	Description string
}

type BTC struct {
	Addresses     []Address
	TXsByCategory wallet.TXsByCategory
}

func New() *BTC {
	btc := &BTC{}
	btc.TXsByCategory = make(map[string]wallet.TXs)
	return btc
}

func (btc *BTC) AddListAddresses(list []string) {
	for _, add := range list {
		btc.Addresses = append(btc.Addresses, Address{Address: add})
	}
}

func (btc BTC) OwnAddress(add string) bool {
	for _, a := range btc.Addresses {
		if a.Address == add {
			return true
		}
	}
	return false
}
