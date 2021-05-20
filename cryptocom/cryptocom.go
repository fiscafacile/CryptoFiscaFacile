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
	cdc.done <- nil
}

func (cdc *CryptoCom) WaitFinish(account string) error {
	err := <-cdc.done
	// Merge TX without Duplicates
	cdc.TXsByCategory.AddUniq(cdc.apiEx.txsByCategory)
	// Add 3916 Source infos
	if _, ok := cdc.Sources["CdC Exchange"]; ok {
		if cdc.Sources["CdC Exchange"].OpeningDate.After(cdc.apiEx.firstTimeUsed) {
			src := cdc.Sources["CdC Exchange"]
			src.OpeningDate = cdc.apiEx.firstTimeUsed
			cdc.Sources["CdC Exchange"] = src
		}
		if cdc.Sources["CdC Exchange"].ClosingDate.Before(cdc.apiEx.lastTimeUsed) {
			src := cdc.Sources["CdC Exchange"]
			src.ClosingDate = cdc.apiEx.lastTimeUsed
			cdc.Sources["CdC Exchange"] = src
		}
	} else {
		cdc.Sources["CdC Exchange"] = source.Source{
			Crypto:        true,
			AccountNumber: account,
			OpeningDate:   cdc.apiEx.firstTimeUsed,
			ClosingDate:   cdc.apiEx.lastTimeUsed,
			LegalName:     "MCO Malta DAX Limited",
			Address:       "Level 7, Spinola Park, Triq Mikiel Ang Borg,\nSt Julian's SPK 1000,\nMalte",
			URL:           "https://crypto.com/exchange",
		}
	}
	return err
}
