package cryptocom

import (
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type CryptoCom struct {
	apiEx                apiEx
	jsonEx               jsonEx
	csvStake             csvStake
	csvSupercharger      csvSupercharger
	csvTransfer          csvTransfer
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

func (cdc *CryptoCom) MergeTXs() {
	// Merge TX without Duplicates
	cdc.TXsByCategory.Add(cdc.jsonEx.txsByCategory)
	cdc.TXsByCategory.AddUniq(cdc.apiEx.txsByCategory)
	cdc.TXsByCategory.AddUniq(cdc.csvStake.txsByCategory)
	cdc.TXsByCategory.AddUniq(cdc.csvSupercharger.txsByCategory)
	cdc.TXsByCategory.AddUniq(cdc.csvTransfer.txsByCategory)
}

func (cdc *CryptoCom) WaitFinish(account string) error {
	err := <-cdc.done
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
