package metamask

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type MetaMask struct {
	CsvTXs        []CsvTX
	TXsByCategory wallet.TXsByCategory
}

func New() *MetaMask {
	mm := &MetaMask{}
	mm.TXsByCategory = make(map[string]wallet.TXs)
	return mm
}
