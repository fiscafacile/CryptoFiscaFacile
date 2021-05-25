package kraken

import (
	"strings"

	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

type Kraken struct {
	api           api
	csvTXs        []csvTX
	done          chan error
	TXsByCategory wallet.TXsByCategory
	Sources       source.Sources
}

func New() *Kraken {
	kr := &Kraken{}
	kr.done = make(chan error)
	kr.TXsByCategory = make(map[string]wallet.TXs)
	kr.Sources = make(source.Sources)
	return kr
}

func (kr *Kraken) GetAPIAllTXs() {
	err := kr.api.getAPITxs()
	if err != nil {
		kr.done <- err
		return
	}
	kr.done <- nil
}

func (kr *Kraken) MergeTXs() {
	// Merge TX without Duplicates
	kr.TXsByCategory.AddUniq(kr.api.txsByCategory)
}

func (kr *Kraken) WaitFinish(account string) error {
	err := <-kr.done
	// Add 3916 Source infos
	if _, ok := kr.Sources["Kraken"]; ok {
		if kr.Sources["Kraken"].OpeningDate.After(kr.api.firstTimeUsed) {
			src := kr.Sources["Kraken"]
			src.OpeningDate = kr.api.firstTimeUsed
			kr.Sources["Kraken"] = src
		}
		if kr.Sources["Kraken"].ClosingDate.Before(kr.api.lastTimeUsed) {
			src := kr.Sources["Kraken"]
			src.ClosingDate = kr.api.lastTimeUsed
			kr.Sources["Kraken"] = src
		}
	} else {
		kr.Sources["Kraken"] = source.Source{
			Crypto:        true,
			AccountNumber: account,
			OpeningDate:   kr.api.firstTimeUsed,
			ClosingDate:   kr.api.lastTimeUsed,
			LegalName:     "Payward Ltd.",
			Address:       "6th Floor,\nOne London Wall,\nLondon, EC2Y 5EB,\nRoyaume-Uni",
			URL:           "https://www.kraken.com",
		}
	}
	return err
}

func ReplaceAssets(assetToReplace string) string {
	assetRplcr := strings.NewReplacer(
		// "ADA.S", "ADA",
		// "ATOM.S", "ATOM",
		// "DOT.S", "DOT",
		// "ETH2.S", "ETH",
		"ETH2", "ETH",
		"EUR.HOLD", "EUR",
		"EUR.M", "EUR",
		// "FLOW.S", "FLOW",
		"FLOWH", "FLOW",
		// "FLOWH.S", "FLOW",
		// "KAVA.S", "KAVA",
		"KFEE", "FEE",
		// "KSM.S", "KSM",
		"USD.HOLD", "USD",
		"USD.M", "USD",
		"XBT", "BTC",
		"XBT.M", "BTC",
		"XETC", "ETC",
		"XETH", "ETH",
		"XLTC", "LTC",
		"XMLN", "MLN",
		"XREP", "REP",
		"XTZ", "XTZ",
		// "XTZ.S", "XTZ",
		"XXBT", "BTC",
		"XXDG", "DOGE",
		"XXLM", "XLM",
		"XXMR", "XMR",
		"XXRP", "XRP",
		"XZEC", "ZEC",
		"ZAUD", "AUD",
		"ZCAD", "CAD",
		"ZEUR", "EUR",
		"ZGBP", "GBP",
		"ZJPY", "JPY",
		"ZRX", "ZRX",
		"ZUSD", "USD",
	)
	return assetRplcr.Replace(assetToReplace)
}
