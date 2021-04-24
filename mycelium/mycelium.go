package mycelium

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type MyCelium struct {
	CsvTXs        []CsvTX
	TXsByCategory wallet.TXsByCategory
}

func New() *MyCelium {
	mc := &MyCelium{}
	mc.TXsByCategory = make(map[string]wallet.TXs)
	return mc
}
