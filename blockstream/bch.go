package blockstream

import (
	"fmt"
	"log"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/btc"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

func (blkst *Blockstream) DetectBCH(b *btc.BTC) {
	bchForkDate := time.Date(2017, time.August, 1, 15, 16, 0, 0, time.UTC)
	w := b.TXsByCategory.GetWallets(bchForkDate, false)
	w.Println("BTC (at time of BCH Fork)", "BTC")
	fmt.Println("Addresses :")
	for _, a := range b.CSVAddresses {
		bal, err := blkst.GetAddressBalanceAtDate(a.Address, bchForkDate)
		if err != nil {
			log.Println("")
			break
		}
		if !bal.IsZero() {
			fmt.Println("  ", a.Address, "balance", bal)
			t := wallet.TX{Timestamp: bchForkDate, Note: "Blockstream API : 478558 Bitcoin Cash Fork on " + a.Address}
			t.Items = make(map[string]wallet.Currencies)
			t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "BCH", Amount: bal})
			b.TXsByCategory["Forks"] = append(b.TXsByCategory["Forks"], t)
		}
	}
}
