package cryptocom

import (
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type CryptoCom struct {
	apiEx                apiEx
	csvAppCryptoTXs      []csvAppCryptoTX
	csvExTransferTXs     []csvExTransferTX
	csvExStakeTXs        []csvExStakeTX
	csvExSuperchargerTXs []csvExSuperchargerTX
	TXsByCategory        wallet.TXsByCategory
}

func New() *CryptoCom {
	cdc := &CryptoCom{}
	cdc.TXsByCategory = make(map[string]wallet.TXs)
	return cdc
}

func (cdc *CryptoCom) GetAPIExchangeTxs(loc *time.Location) (err error) {
	err = cdc.apiEx.getAPIExchangeTxs(loc)
	if err != nil {
		return
	}
	cdc.TXsByCategory.Add(cdc.apiEx.txsByCategory)
	return
}
