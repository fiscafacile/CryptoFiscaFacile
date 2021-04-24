package metamask

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type MetaMask struct {
	CsvTXs   []CsvTX
	Accounts wallet.Accounts
}

func New() *MetaMask {
	mm := &MetaMask{}
	mm.Accounts = make(map[string]wallet.TXs)
	return mm
}
