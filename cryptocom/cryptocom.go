package cryptocom

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type CryptoCom struct {
	CsvTXs               []CsvTX
	CsvTXsExTransfer     []CsvTXExTransfer
	CsvTXsExStake        []CsvTXExStake
	CsvTXsExSupercharger []CsvTXExSupercharger
	Accounts             wallet.Accounts
}

func New() *CryptoCom {
	cdc := &CryptoCom{}
	cdc.Accounts = make(map[string]wallet.TXs)
	return cdc
}
