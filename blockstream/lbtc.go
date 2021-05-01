package blockstream

import (
	"fmt"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/btc"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

func (blkst *Blockstream) DetectLBTC(b *btc.BTC) {
	lbtcForkDate := time.Date(2017, time.December, 18, 18, 34, 0, 0, time.UTC)
	w := b.TXsByCategory.GetWallets(lbtcForkDate, false, false)
	w.Println("BTC (at time of LBTC Fork)", "BTC")
	fmt.Println("Addresses :")
	for _, a := range b.CSVAddresses {
		bal, err := blkst.GetAddressBalanceAtDate(a.Address, lbtcForkDate)
		if err != nil {
			log.Println("")
			break
		}
		if !bal.IsZero() {
			fmt.Println("  ", a.Address, "balance", bal)
			t := wallet.TX{Timestamp: lbtcForkDate, Note: "Blockstream API : 499999 Lightning Bitcoin Fork on " + a.Address}
			t.Items = make(map[string]wallet.Currencies)
			t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "LBTC", Amount: bal})
			b.TXsByCategory["Forks"] = append(b.TXsByCategory["Forks"], t)
		}
	}
}
