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
	done                 chan error
	TXsByCategory        wallet.TXsByCategory
}

func New() *CryptoCom {
	cdc := &CryptoCom{}
	cdc.done = make(chan error)
	cdc.TXsByCategory = make(map[string]wallet.TXs)
	return cdc
}

func (cdc *CryptoCom) GetAPIExchangeTXs(loc *time.Location) {
	err := cdc.apiEx.getAllTXs(loc)
	if err != nil {
		cdc.done <- err
		return
	}
	cdc.TXsByCategory.Add(cdc.apiEx.txsByCategory)
	cdc.done <- nil
}

func (cdc *CryptoCom) WaitFinish() error {
	return <-cdc.done
}
