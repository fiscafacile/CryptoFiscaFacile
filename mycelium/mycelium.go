package mycelium

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type MyCelium struct {
	CsvTXs   []CsvTX
	Accounts wallet.Accounts
}

func New() *MyCelium {
	mc := &MyCelium{}
	mc.Accounts = make(map[string]wallet.TXs)
	return mc
}
