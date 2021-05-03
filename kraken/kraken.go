package kraken

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Kraken struct {
	api           api
	csvTXs        []csvTX
	TXsByCategory wallet.TXsByCategory
}

func New() *Kraken {
	kr := &Kraken{}
	kr.TXsByCategory = make(map[string]wallet.TXs)
	return kr
}

func (kr *Kraken) GetAPITxs() (err error) {
	err = kr.api.getAPITxs()
	if err != nil {
		return
	}
	kr.TXsByCategory.Add(kr.api.txsByCategory)
	return
}
