package kraken

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Kraken struct {
	csvTXs         []csvTX
	TXsByCategory  wallet.TXsByCategory
}

func New() *Kraken {
	kr := &Kraken{}
	kr.TXsByCategory = make(map[string]wallet.TXs)
	return kr
}
