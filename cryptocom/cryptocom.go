package cryptocom

import (
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/source"
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
	Sources              source.Sources
}

func New() *CryptoCom {
	cdc := &CryptoCom{}
	cdc.done = make(chan error)
	cdc.TXsByCategory = make(wallet.TXsByCategory)
	cdc.Sources = make(source.Sources)
	return cdc
}

func (cdc *CryptoCom) GetAPIExchangeTXs(loc *time.Location) {
	err := cdc.apiEx.getAllTXs(loc)
	if err != nil {
		cdc.done <- err
		return
	}
	cdc.TXsByCategory.Add(cdc.apiEx.txsByCategory)
	cdc.Sources["CdC Exchange"] = source.Source{
		Crypto:        true,
		AccountNumber: "emailAROBASEdomainPOINTcom",
		OpeningDate:   cdc.apiEx.firstTimeUsed,
		ClosingDate:   cdc.apiEx.lastTimeUsed,
		LegalName:     "MCO Malta DAX Limited",
		Address:       "Level 7, Spinola Park, Triq Mikiel Ang Borg,\nSt Julian's SPK 1000,\nMalte",
		URL:           "https://crypto.com/exchange",
	}
	cdc.done <- nil
}

func (cdc *CryptoCom) WaitFinish() error {
	return <-cdc.done
}
