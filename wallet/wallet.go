package wallet

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/davecgh/go-spew/spew"
	"github.com/shopspring/decimal"
)

type Nft struct {
	ID     string
	Name   string
	Symbol string
}

type Nfts []Nft

type Currency struct {
	Code   string
	Amount decimal.Decimal
}

type Currencies []Currency

type WalletCurrencies map[string]decimal.Decimal

type Wallets struct {
	Date       time.Time
	Currencies WalletCurrencies
	Nfts       Nfts
}

type TX struct {
	Timestamp time.Time
	ID        string
	Source    string
	Category  string
	Items     map[string]Currencies
	Nfts      map[string]Nfts
	Note      string
	used      bool
}

type TXs []TX

type TXsByCategory map[string]TXs

func AskForHelp(id string, tx interface{}, alreadyAsked []string) []string {
	found := false
	for _, i := range alreadyAsked {
		if i == id {
			found = true
			break
		}
	}
	if !found {
		log.Println("Nouveau type de transaction détecté", id, "merci de copier ce texte dans t.me/cryptofiscafacile pour que nous puissions ajouter son support :", base64.StdEncoding.EncodeToString([]byte(spew.Sdump(tx))))
		alreadyAsked = append(alreadyAsked, id)
	}
	return alreadyAsked
}

func (c *Currency) IsFiat() bool {
	if c.Code == "EUR" ||
		c.Code == "USD" ||
		c.Code == "HKD" {
		return true
	}
	return false
}

func (c *Currency) Println(filter string) {
	if strings.Contains(filter, c.Code) ||
		filter == "" {
		fmt.Println(c.Amount, c.Code)
	}
}

func (c Currency) GetExchangeRate(date time.Time, to string) (rate decimal.Decimal, err error) {
	if !c.IsFiat() {
		gecko, err := NewCoinGeckoAPI()
		if err == nil {
			ratesCG, err := gecko.GetExchangeRates(date, c.Code)
			if err == nil {
				for _, r := range ratesCG.Rates {
					if r.Quote == to && !r.Rate.IsZero() {
						return r.Rate, nil
					}
				}
			}
		}
	}
	var layer CoinLayer
	ratesCL, err := layer.GetExchangeRates(date, to)
	if err == nil {
		for k, v := range ratesCL.Rates {
			if k == c.Code && v != 0 {
				return decimal.NewFromFloat(ratesCL.Rates[c.Code]), nil
			}
		}
	}
	var api CoinAPI
	rates, err := api.GetExchangeRates(date, to)
	if err == nil {
		for _, r := range rates.Rates {
			if r.Quote == c.Code && !r.Rate.IsZero() {
				return r.Rate, nil
			}
		}
	}
	return rate, errors.New("Cannot find rate for " + c.Code + " at " + date.String())
}

func (wc WalletCurrencies) Add(a WalletCurrencies) {
	for k, v := range a {
		wc[k] = wc[k].Add(v)
	}
}

func (w Wallets) CalculateTotalValue(native string) (totalValue Currency, err error) {
	totalValue.Code = native
	for k, v := range w.Currencies {
		fmt.Print(".")
		if k == native {
			totalValue.Amount = totalValue.Amount.Add(v)
		} else {
			c := Currency{Code: k, Amount: v}
			rate, err := c.GetExchangeRate(w.Date, native)
			if err != nil {
				log.Println(err)
			} else {
				totalValue.Amount = totalValue.Amount.Add(rate.Mul(v))
			}
		}
	}
	return
}

func (w Wallets) Round(rounding bool) {
	for k, v := range w.Currencies {
		if rounding {
			if k == "BAB" {
				if v.Abs().LessThan(decimal.NewFromInt(1)) {
					delete(w.Currencies, k)
				}
			} else if k == "EUR" || k == "USD" || k == "HKD" || k == "CRO" || k == "USDC" || k == "USDT" || k == "sUSD" || k == "XRP" || k == "IOT" {
				if v.Abs().LessThan(decimal.NewFromFloat(0.5)) {
					delete(w.Currencies, k)
				}
			} else if k == "LPT" {
				if v.Abs().LessThan(decimal.NewFromInt(100)) {
					delete(w.Currencies, k)
				}
			} else {
				if v.Abs().LessThan(decimal.NewFromFloat(0.01)) {
					delete(w.Currencies, k)
				}
			}
		} else {
			if v.IsZero() {
				delete(w.Currencies, k)
			}
		}
	}
}

func (w Wallets) Println(name string, filter string) {
	fmt.Println(strings.Repeat("-", len(name)+11))
	fmt.Println("| " + name + " Wallet |")
	fmt.Println(strings.Repeat("-", len(name)+11))
	fmt.Println("Time :", w.Date.UTC())
	fmt.Println("Amounts :")
	keys := make([]string, 0, len(w.Currencies))
	for k := range w.Currencies {
		if filter == "" || strings.Contains(filter, k) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println("  ", w.Currencies[k], k)
	}
}

func (tx *TX) SimilarDate(delta time.Duration, t time.Time) bool {
	if tx.Timestamp.Sub(t) < delta &&
		t.Sub(tx.Timestamp) < delta {
		return true
	}
	return false
}

func (tx TX) Println(filter string) (printed bool) {
	toDisplay := strings.Split(filter, ",")
	if filter == "" {
		printed = true
	} else {
		printed = false
	}
	toPrint := fmt.Sprintln("Time :", tx.Timestamp.UTC())
	if tx.ID != "" {
		toPrint += fmt.Sprintln("ID :", tx.ID)
	}
	if tx.Source != "" {
		toPrint += fmt.Sprintln("Source :", tx.Source)
	}
	if tx.Category != "" {
		toPrint += fmt.Sprintln("Category :", tx.Category)
	}
	for k, v := range tx.Items {
		toPrint += fmt.Sprintln(k, ":")
		for _, i := range v {
			toPrint += fmt.Sprintln("  ", i.Amount.String(), i.Code)
			for _, d := range toDisplay {
				if d == i.Code {
					printed = true
				}
			}
		}
	}
	for k, v := range tx.Nfts {
		toPrint += fmt.Sprintln("NFT", k, ":")
		for _, i := range v {
			toPrint += fmt.Sprintln("  ", i.Name, i.Symbol, i.ID)
		}
	}
	toPrint += fmt.Sprintln("Note :", tx.Note)
	if printed {
		fmt.Print(toPrint)
	}
	return
}

func (tx TX) GetBalances(includeFiat, includeFee bool) (cs WalletCurrencies) {
	cs = make(WalletCurrencies)
	for k, i := range tx.Items {
		for _, c := range i {
			if c.Code != "" &&
				(includeFiat || !c.IsFiat()) {
				if (k == "Fee" && includeFee) || k == "From" || k == "Lost" {
					cs[c.Code] = cs[c.Code].Sub(c.Amount)
				} else if k == "To" {
					cs[c.Code] = cs[c.Code].Add(c.Amount)
				}
			}
		}
	}
	return
}

func (tx TX) SameBalances(t TX) bool {
	cbs := tx.GetBalances(false, false)
	b := t.GetBalances(false, false)
	if len(cbs) != len(b) {
		return false
	}
	for k, v := range cbs {
		if !b[k].Equal(v) {
			return false
		}
	}
	return true
}

func (txs TXs) GetBalances(includeFiat, includeFee bool) (cs WalletCurrencies) {
	cs = make(WalletCurrencies)
	for _, tx := range txs {
		cs.Add(tx.GetBalances(includeFiat, includeFee))
	}
	return
}

func (txs TXs) Println(name string, filter string) {
	fmt.Println(strings.Repeat("-", len(name)+11))
	fmt.Println("| TXs in " + name + " |")
	fmt.Println(strings.Repeat("-", len(name)+11))
	printed := false
	for _, tx := range txs {
		if printed {
			fmt.Println(strings.Repeat("-", len(name)+11))
		}
		printed = tx.Println(filter)
	}
}

func (txs TXs) SortByDate(chrono bool) {
	if chrono {
		sort.Slice(txs, func(i, j int) bool {
			return txs[i].Timestamp.Before(txs[j].Timestamp)
		})
	} else {
		sort.Slice(txs, func(i, j int) bool {
			return txs[i].Timestamp.After(txs[j].Timestamp)
		})
	}
}

func (txs TXs) After(date time.Time) TXs {
	var filteredTXs TXs
	for _, t := range txs {
		if t.Timestamp.After(date) {
			filteredTXs = append(filteredTXs, t)
		}
	}
	return filteredTXs
}

func (txs TXs) Before(date time.Time) TXs {
	var filteredTXs TXs
	for _, t := range txs {
		if t.Timestamp.Before(date) {
			filteredTXs = append(filteredTXs, t)
		}
	}
	return filteredTXs
}

func (txs TXs) ApplyFromReversal() TXs {
	var filteredTXs TXs
	var toCancel Currencies
	txs.SortByDate(false)
	for _, tx := range txs {
		isReversal := false
		for _, c := range tx.Items["From"] {
			toCancel = append(toCancel, c)
			isReversal = true
		}
		if !isReversal {
			for _, c := range tx.Items["To"] {
				for i, tc := range toCancel {
					if c.Code == tc.Code &&
						c.Amount.Equal(tc.Amount) {
						if len(toCancel) > 1 {
							toCancel[i] = toCancel[len(toCancel)-1]
						}
						toCancel = toCancel[:len(toCancel)-1]
						break
					}
				}
				filteredTXs = append(filteredTXs, tx)
			}
		}
	}
	if len(toCancel) > 0 {
		log.Println("Couldn't apply", len(toCancel), "reversals")
	}
	return filteredTXs
}

func (txs TXs) AddFromNativeValue(native string) TXs {
	for i, t := range txs {
		for _, c := range t.Items["To"] {
			fmt.Print(".")
			rate, err := c.GetExchangeRate(t.Timestamp, native)
			if err != nil {
				log.Println(err)
			} else {
				txs[i].Items["From"] = append(txs[i].Items["From"], Currency{Code: native, Amount: c.Amount.Mul(rate)})
			}
		}
	}
	return txs
}

func (txs TXsByCategory) Println(filter string) {
	for k, v := range txs {
		v.Println("Category "+k, filter)
	}
}

func (txs TXsByCategory) GetCoinsList(includeFiat bool) (coins []string) {
	for _, v := range txs {
		for _, tx := range v {
			for _, i := range tx.Items {
				for _, c := range i {
					found := false
					for _, coin := range coins {
						if coin == c.Code {
							found = true
						}
					}
					if !found &&
						(includeFiat || !c.IsFiat()) {
						coins = append(coins, c.Code)
					}
				}
			}
		}
	}
	sort.Strings(coins)
	return
}

func (txs TXsByCategory) StockToXlsx(filename string) {
	f := excelize.NewFile()
	var allTXs TXs
	for cat, list := range txs {
		for _, tx := range list {
			if cat == "Withdrawals" {
				tx.Category = "Retrait"
			} else if cat == "Deposits" {
				tx.Category = "Dépot"
			} else if cat == "Exchanges" {
				tx.Category = "Echange"
			} else if cat == "Fees" {
				tx.Category = "Frais"
			} else if cat == "Gifts" {
				tx.Category = "Don"
			} else if cat == "Transfers" {
				tx.Category = "Transfert"
			} else {
				tx.Category = cat
			}
			allTXs = append(allTXs, tx)
		}
	}
	allTXs.SortByDate(true)
	coins := txs.GetCoinsList(false)
	for _, coin := range coins {
		f.NewSheet(coin)
		f.SetCellValue(coin, "A1", "Date (UTC)")
		f.SetCellValue(coin, "B1", "Type d'opération")
		f.SetCellValue(coin, "C1", "Entrée")
		f.SetCellValue(coin, "D1", "Sortie")
		f.SetCellValue(coin, "E1", "Balance")
		f.SetCellValue(coin, "F1", "Note")
		row := 2
		var balance decimal.Decimal
		for _, t := range allTXs {
			var input decimal.Decimal
			var output decimal.Decimal
			for k, i := range t.Items {
				for _, c := range i {
					if c.Code == coin {
						if k == "To" {
							input = input.Add(c.Amount)
							balance = balance.Add(c.Amount)
						} else if k == "From" || k == "Fee" || k == "Lost" {
							output = output.Add(c.Amount)
							balance = balance.Sub(c.Amount)
						}
					}
				}
			}
			if !input.IsZero() || !output.IsZero() {
				f.SetCellValue(coin, "A"+strconv.Itoa(row), t.Timestamp.Format("02/01/2006 15:04:05"))
				f.SetCellValue(coin, "B"+strconv.Itoa(row), t.Category)
				if !input.IsZero() {
					in, _ := input.Float64()
					f.SetCellValue(coin, "C"+strconv.Itoa(row), in)
				}
				if !output.IsZero() {
					out, _ := output.Float64()
					f.SetCellValue(coin, "D"+strconv.Itoa(row), out)
				}
				bal, _ := balance.Float64()
				f.SetCellValue(coin, "E"+strconv.Itoa(row), bal)
				f.SetCellValue(coin, "F"+strconv.Itoa(row), t.Note)
				row += 1
			}
		}
		f.SetColWidth(coin, "A", "A", 18)
		f.SetColWidth(coin, "B", "B", 16)
		f.SetColWidth(coin, "F", "F", 50)
	}
	f.DeleteSheet("Sheet1")
	if err := f.SaveAs(filename); err != nil {
		log.Fatal(err)
	}
}

func (txs TXsByCategory) GetWallets(date time.Time, includeFiat bool, rounding bool) (w Wallets) {
	w.Date = date
	w.Currencies = make(WalletCurrencies)
	for _, v := range txs {
		w.Currencies.Add(v.Before(date).GetBalances(includeFiat, true))
	}
	w.Round(rounding)
	return
}

func (txs TXsByCategory) Add(a TXsByCategory) {
	for k, v := range a {
		txs[k] = append(txs[k], v...)
	}
}

func (txs TXsByCategory) AddUniq(a TXsByCategory) {
	for k, v := range a {
		for _, tx := range v {
			found := false
			for _, t := range txs[k] {
				if (t.SimilarDate(2*time.Hour+time.Second, tx.Timestamp) && t.SameBalances(tx)) ||
					(t.ID == tx.ID && tx.ID != "") {
					found = true
					break
				}
			}
			if !found {
				txs[k] = append(txs[k], tx)
			}
		}
	}
}

func (txs TXsByCategory) FindTransfers() TXsByCategory {
	similarTimeDelta := 12 * time.Hour
	txs["Deposits"].SortByDate(true)
	txs["Withdrawals"].SortByDate(true)
	for di, depTX := range txs["Deposits"] {
		if !depTX.used && len(depTX.Items["To"]) > 0 {
			var depFees decimal.Decimal
			if _, ok := depTX.Items["Fee"]; ok {
				for _, f := range depTX.Items["Fee"] {
					depFees = depFees.Add(f.Amount)
				}
			}
			for wi, witTX := range txs["Withdrawals"] {
				if !witTX.used && len(witTX.Items["From"]) > 0 {
					if depTX.Items["To"][0].Code == witTX.Items["From"][0].Code &&
						depTX.SimilarDate(similarTimeDelta, witTX.Timestamp) &&
						strings.Split(depTX.Note, ":")[0] != strings.Split(witTX.Note, ":")[0] {
						var witFees decimal.Decimal
						if _, ok := witTX.Items["Fee"]; ok {
							for _, f := range witTX.Items["Fee"] {
								witFees = witFees.Add(f.Amount)
							}
						}
						// log.Println("Here")
						// depTX.Println("")
						// witTX.Println("")
						if depTX.Items["To"][0].Amount.Equal(witTX.Items["From"][0].Amount) ||
							depTX.Items["To"][0].Amount.Equal(witTX.Items["From"][0].Amount.Sub(witFees)) ||
							depTX.Items["To"][0].Amount.Equal(witTX.Items["From"][0].Amount.Sub(depFees)) {
							txs["Deposits"][di].used = true
							txs["Withdrawals"][wi].used = true
							t := TX{Timestamp: witTX.Timestamp, Note: witTX.Note + " => " + depTX.Note}
							t.ID = witTX.ID + "-" + depTX.ID
							t.Items = make(map[string]Currencies)
							t.Items["To"] = append(t.Items["To"], depTX.Items["To"]...)
							t.Items["From"] = append(t.Items["From"], witTX.Items["From"]...)
							if _, ok := witTX.Items["Fee"]; ok {
								t.Items["Fee"] = append(t.Items["Fee"], witTX.Items["Fee"]...)
							}
							if _, ok := depTX.Items["Fee"]; ok {
								for _, df := range depTX.Items["Fee"] {
									missing := true
									for _, f := range t.Items["Fee"] {
										if f.Code == df.Code &&
											f.Amount.Equal(df.Amount) {
											missing = false
										}
									}
									if missing {
										t.Items["Fee"] = append(t.Items["Fee"], df)
									}
								}
							}
							if _, ok := witTX.Items["Lost"]; ok {
								t.Items["Lost"] = append(t.Items["Lost"], witTX.Items["Lost"]...)
							}
							if _, ok := depTX.Items["Lost"]; ok {
								t.Items["Lost"] = append(t.Items["Lost"], depTX.Items["Lost"]...)
							}
							txs["Transfers"] = append(txs["Transfers"], t)
							break
						}
					}
				}
			}
		}
	}
	// Purge used TXs
	var realDeposits TXs
	for _, depTX := range txs["Deposits"] {
		if !depTX.used {
			realDeposits = append(realDeposits, depTX)
		}
	}
	txs["Deposits"] = realDeposits
	var realWithdrawals TXs
	for _, witTX := range txs["Withdrawals"] {
		if !witTX.used {
			realWithdrawals = append(realWithdrawals, witTX)
		}
	}
	txs["Withdrawals"] = realWithdrawals
	return txs
}

func (txs TXsByCategory) FindCashInOut(native string) {
	var realExchanges TXs
	for _, exTX := range txs["Exchanges"] {
		fmt.Print(".")
		fromHasFiat := false
		for _, i := range exTX.Items["From"] {
			if i.IsFiat() {
				fromHasFiat = true
			}
		}
		toHasFiat := false
		for _, c := range exTX.Items["To"] {
			if c.IsFiat() {
				toHasFiat = true
			}
		}
		if fromHasFiat {
			txs["CashIn"] = append(txs["CashIn"], exTX)
		} else if toHasFiat {
			txs["CashOut"] = append(txs["CashOut"], exTX)
		} else {
			realExchanges = append(realExchanges, exTX)
		}
	}
	if len(realExchanges) > 0 {
		txs["Exchanges"] = realExchanges
	} else {
		delete(txs, "Exchanges")
	}
}

func (txs TXsByCategory) RemoveDelistedCoins(coin string) {
	var coinBalance decimal.Decimal
	var lastTx *TX
	for k, v := range txs {
		for tk, tv := range v {
			for ik, iv := range tv.Items {
				for id := range iv {
					if iv[id].Code == coin {
						lastTx = &txs[k][tk]
						if ik == "From" || ik == "Fee" {
							coinBalance = coinBalance.Sub(iv[id].Amount)
						} else if ik == "To" {
							coinBalance = coinBalance.Add(iv[id].Amount)
						}
					}
				}
			}
		}
	}
	lastTx.Items["Lost"] = Currencies{
		Currency{
			Amount: coinBalance,
			Code:   coin,
		},
	}
	lastTx.Note += " Force balance to 0 for " + coin + " as it has been delisted"
}

func (txs TXsByCategory) SortByDate(chrono bool) {
	for k := range txs {
		txs[k].SortByDate(chrono)
	}
}

func (txs TXsByCategory) PrintStats(native string) {
	fmt.Println("-------------------------------")
	fmt.Println("| Quantity of TXs By Category |")
	fmt.Println("-------------------------------")
	keys := make([]string, 0, len(txs))
	for k := range txs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println(k, ":", len(txs[k]), "TXs")
	}
}

func (txs TXsByCategory) CheckConsistency(loc *time.Location) {
	fmt.Println("--------------------------------------------------------")
	fmt.Println("| List of Unjustified Withdrawals (after 2019 Jan 1st) |")
	for _, tx := range txs["Withdrawals"] {
		if tx.Timestamp.After(time.Date(2018, time.December, 31, 23, 59, 59, 999, loc)) &&
			len(tx.Items["From"]) > 0 {
			fmt.Println("--------------------------------------------------------")
			tx.Println("")
		}
	}
	fmt.Println("--------------------------------------------------------")
	fmt.Println("| List of Non-Zero balance Transfers                   |")
	for _, tx := range txs["Transfers"] {
		txcs := tx.GetBalances(false, false)
		for _, v := range txcs {
			if !v.IsZero() {
				fmt.Println("--------------------------------------------------------")
				tx.Println("")
			}
		}
	}
	fmt.Println("--------------------------------------------------------")
	fmt.Println("| List of Negative Amounts TXs                         |")
	for _, v := range txs {
		for _, tx := range v {
			for _, i := range tx.Items {
				for _, c := range i {
					if c.Amount.IsNegative() {
						fmt.Println("--------------------------------------------------------")
						tx.Println("")
					}
				}
			}
		}
	}
	for _, cat := range []string{"Deposits", "AirDrops", "CommercialRebates", "Interests", "Minings", "Referrals"} {
		fmt.Println("--------------------------------------------------------")
		fmt.Println("| List of " + cat + " with some From" + strings.Repeat(" ", 30-len(cat)) + "|")
		txsCat := txs[cat]
		if cat == "CommercialRebates" {
			txsCat = txsCat.ApplyFromReversal()
		}
		for _, tx := range txsCat {
			for k := range tx.Items {
				if k == "From" {
					fmt.Println("--------------------------------------------------------")
					tx.Println("")
				}
			}
		}
	}
	fmt.Println("--------------------------------------------------------")
	fmt.Println("| List of Withdrawals with some To                     |")
	for _, tx := range txs["Withdrawals"] {
		for k := range tx.Items {
			if k == "To" {
				fmt.Println("--------------------------------------------------------")
				tx.Println("")
			}
		}
	}
	fmt.Println("--------------------------------------------------------")
}
