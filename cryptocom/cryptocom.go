package cryptocom

import (
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type CryptoCom struct {
	CsvTXs               []CsvTX
	CsvTXsExTransfer     []CsvTXExTransfer
	CsvTXsExStake        []CsvTXExStake
	CsvTXsExSupercharger []CsvTXExSupercharger
	TXsByCategory        wallet.TXsByCategory
}

func New() *CryptoCom {
	cdc := &CryptoCom{}
	cdc.TXsByCategory = make(map[string]wallet.TXs)
	return cdc
}
