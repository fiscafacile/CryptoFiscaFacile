package blockstream

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fiscafacile/CryptoFiscaFacile/btc"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
	"gopkg.in/resty.v0"
)

type apiTXStatus struct {
	Confirmed   bool   `json:"confirmed"`
	BlockHeight int    `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	BlockTime   int    `json:"block_time"`
}

type apiTX struct {
	used     bool   `json:-`
	Txid     string `json:"txid"`
	Version  int    `json:"version"`
	Locktime int    `json:"locktime"`
	Vin      []struct {
		Txid    string `json:"txid"`
		Vout    int    `json:"vout"`
		Prevout struct {
			Scriptpubkey        string `json:"scriptpubkey"`
			ScriptpubkeyAsm     string `json:"scriptpubkey_asm"`
			ScriptpubkeyType    string `json:"scriptpubkey_type"`
			ScriptpubkeyAddress string `json:"scriptpubkey_address"`
			Value               int    `json:"value"`
		} `json:"prevout"`
		Scriptsig    string   `json:"scriptsig"`
		ScriptsigAsm string   `json:"scriptsig_asm"`
		Witness      []string `json:"witness"`
		IsCoinbase   bool     `json:"is_coinbase"`
		Sequence     int64    `json:"sequence"`
	} `json:"vin"`
	Vout []struct {
		Scriptpubkey        string `json:"scriptpubkey"`
		ScriptpubkeyAsm     string `json:"scriptpubkey_asm"`
		ScriptpubkeyType    string `json:"scriptpubkey_type"`
		ScriptpubkeyAddress string `json:"scriptpubkey_address"`
		Value               int    `json:"value"`
	} `json:"vout"`
	Size   int         `json:"size"`
	Weight int         `json:"weight"`
	Fee    int         `json:"fee"`
	Status apiTXStatus `json:"status"`
}

type apiTXs []apiTX

func (txs apiTXs) SortByHeight(asc bool) {
	if asc {
		sort.Slice(txs, func(i, j int) bool {
			return txs[i].Status.BlockHeight < txs[j].Status.BlockHeight
		})
	} else {
		sort.Slice(txs, func(i, j int) bool {
			return txs[i].Status.BlockHeight > txs[j].Status.BlockHeight
		})
	}
}

func (blkst *Blockstream) GetAddressTXs(add string) (txs apiTXs, err error) {
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("BlockStream/address/txs", add, &txs)
	}
	if !useCache || err != nil {
		resp, err := resty.R().SetHeaders(map[string]string{
			"Accept": "application/json",
		}).Get("https://blockstream.info/api/address/" + add + "/txs")
		if err != nil || resp.StatusCode() != http.StatusOK {
			time.Sleep(6 * time.Second)
			resp, err = resty.R().SetHeaders(map[string]string{
				"Accept": "application/json",
			}).Get("https://blockstream.info/api/address/" + add + "/txs")
		}
		if err != nil || resp.StatusCode() != http.StatusOK {
			return txs, errors.New("Blockstream API : Error Getting BTC TX for " + add)
		}
		err = json.Unmarshal(resp.Body(), &txs)
		if err != nil {
			return txs, errors.New("Blockstream API : Error Unmarshaling BTC TX for " + add)
		}
		if useCache {
			err = db.Write("BlockStream/address/txs", add, txs)
			if err != nil {
				return txs, errors.New("Blockstream API : Error Caching" + add)
			}
		}
	}
	return txs, nil
}

func (blkst *Blockstream) GetAddressBalanceAtDate(add string, date time.Time) (bal decimal.Decimal, err error) {
	apiTXs, err := blkst.GetAddressTXs(add)
	if err != nil {
		return
	}
	apiTXs.SortByHeight(true)
	for _, tx := range apiTXs {
		if time.Unix(int64(tx.Status.BlockTime), 0).After(date) {
			break
		}
		for _, vin := range tx.Vin {
			if add == vin.Prevout.ScriptpubkeyAddress {
				bal = bal.Sub(decimal.New(int64(vin.Prevout.Value), -8))
			}
		}
		for _, vout := range tx.Vout {
			if add == vout.ScriptpubkeyAddress {
				bal = bal.Add(decimal.New(int64(vout.Value), -8))
			}
		}
	}
	return
}

func (blkst *Blockstream) GetAllTXs(b *btc.BTC) {
	for _, btc := range b.CSVAddresses {
		apiTXs, err := blkst.GetAddressTXs(btc.Address)
		if err != nil {
			log.Println(err)
		}
		for _, tx := range apiTXs {
			found := false
			for _, have := range blkst.apiTXs {
				if tx.Txid == have.Txid {
					found = true
					break
				}
			}
			if !found {
				blkst.apiTXs = append(blkst.apiTXs, tx)
			}
		}
	}
	for i, tx := range blkst.apiTXs {
		if !tx.used {
			valueIn := 0
			isInVinPrevVout := false
			missing := ""
			with := ""
			for _, vin := range tx.Vin {
				if b.OwnAddress(vin.Prevout.ScriptpubkeyAddress) {
					valueIn -= vin.Prevout.Value
					isInVinPrevVout = true
					with += " " + vin.Prevout.ScriptpubkeyAddress
				} else {
					if missing == "" {
						missing = " missing :"
					}
					missing += " " + vin.Prevout.ScriptpubkeyAddress
				}
			}
			if isInVinPrevVout && missing != "" {
				log.Println("Blockstream API : found co-signed address", missing[11:], "with", with)
			}
			valueOut := 0
			isInVout := false
			dest := ""
			for _, vout := range tx.Vout {
				if b.OwnAddress(vout.ScriptpubkeyAddress) {
					valueOut += vout.Value
					isInVout = true
				} else {
					if dest == "" {
						dest = " destination"
					}
					dest += " " + vout.ScriptpubkeyAddress
				}
			}
			if valueIn+valueOut == 0 {
				log.Println("Blockstream API : Detected zero Value TX")
				// spew.Dump(tx)
			}
			if isInVinPrevVout {
				t := wallet.TX{Timestamp: time.Unix(int64(tx.Status.BlockTime), 0), Note: "Blockstream API : " + strconv.Itoa(tx.Status.BlockHeight) + " " + tx.Txid + dest + missing}
				t.Items = make(map[string][]wallet.Currency)
				t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "BTC", Amount: decimal.New(int64(tx.Fee), -8)})
				if isInVout && dest == "" {
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: "BTC", Amount: decimal.New(int64(-valueIn-tx.Fee), -8)})
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "BTC", Amount: decimal.New(int64(valueOut), -8)})
					b.TXsByCategory["Transfers"] = append(b.TXsByCategory["Transfers"], t)
				} else if is, desc, val, curr := b.IsTxCashOut(tx.Txid); is {
					t.Note += " crypto_payment " + desc
					from := wallet.Currency{Code: "BTC", Amount: decimal.New(int64(-valueOut-valueIn-tx.Fee), -8)}
					t.Items["From"] = append(t.Items["From"], from)
					if val.IsZero() {
						rate, err := from.GetExchangeRate(t.Timestamp, "EUR")
						if err == nil {
							t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "EUR", Amount: from.Amount.Mul(rate)})
						}
					} else {
						t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr, Amount: val})
					}
					b.TXsByCategory["CashOut"] = append(b.TXsByCategory["CashOut"], t)
				} else if is, desc, val, curr := b.IsTxExchange(tx.Txid); is {
					t.Note += " crypto_exchange " + desc
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: "BTC", Amount: decimal.New(int64(-valueOut-valueIn-tx.Fee), -8)})
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr, Amount: val})
					b.TXsByCategory["Exchanges"] = append(b.TXsByCategory["Exchanges"], t)
				} else if is, desc := b.IsTxGift(tx.Txid); is {
					t.Note += " gift " + desc
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: "BTC", Amount: decimal.New(int64(-valueOut-valueIn-tx.Fee), -8)})
					b.TXsByCategory["Gifts"] = append(b.TXsByCategory["Gifts"], t)
				} else {
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: "BTC", Amount: decimal.New(int64(-valueOut-valueIn-tx.Fee), -8)})
					b.TXsByCategory["Withdrawals"] = append(b.TXsByCategory["Withdrawals"], t)
				}
				blkst.apiTXs[i].used = true
			} else if isInVout {
				t := wallet.TX{Timestamp: time.Unix(int64(tx.Status.BlockTime), 0), Note: "Blockstream API : " + strconv.Itoa(tx.Status.BlockHeight) + " " + tx.Txid}
				t.Items = make(map[string][]wallet.Currency)
				t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "BTC", Amount: decimal.New(int64(tx.Fee), -8)})
				if is, desc, val := b.HasCustody(tx.Txid); is {
					t.Note += " crypto_custody " + desc
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "BTC", Amount: decimal.New(int64(valueOut), -8).Sub(val)})
				} else {
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "BTC", Amount: decimal.New(int64(valueOut), -8)})
				}
				if is, desc, val, curr := b.IsTxCashIn(tx.Txid); is {
					t.Note += " crypto_purchase " + desc
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: curr, Amount: val})
					b.TXsByCategory["CashIn"] = append(b.TXsByCategory["CashIn"], t)
				} else if is, desc, val, curr := b.IsTxExchange(tx.Txid); is {
					t.Note += " crypto_exchange " + desc
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: curr, Amount: val})
					b.TXsByCategory["Exchanges"] = append(b.TXsByCategory["Exchanges"], t)
				} else {
					b.TXsByCategory["Deposits"] = append(b.TXsByCategory["Deposits"], t)
				}
				blkst.apiTXs[i].used = true
			} else {
				log.Println("Blockstream API : Unmanaged TX")
				spew.Dump(tx)
			}
		}
	}
	blkst.done <- nil
}
