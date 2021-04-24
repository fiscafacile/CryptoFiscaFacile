package blockstream

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
	"gopkg.in/resty.v0"
)

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
	Size   int `json:"size"`
	Weight int `json:"weight"`
	Fee    int `json:"fee"`
	Status struct {
		Confirmed   bool   `json:"confirmed"`
		BlockHeight int    `json:"block_height"`
		BlockHash   string `json:"block_hash"`
		BlockTime   int    `json:"block_time"`
	} `json:"status"`
}

func (blkst *Blockstream) apiGetAllTXs() (err error) {
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	for _, btc := range blkst.csvAddresses {
		var apiTXs []apiTX
		if useCache {
			err = db.Read("BlockStream/address/txs", btc.address, &apiTXs)
		}
		if !useCache || err != nil {
			resp, err := resty.R().SetHeaders(map[string]string{
				"Accept": "application/json",
			}).Get("https://blockstream.info/api/address/" + btc.address + "/txs")
			if err != nil || resp.StatusCode() != http.StatusOK {
				time.Sleep(6 * time.Second)
				resp, err = resty.R().SetHeaders(map[string]string{
					"Accept": "application/json",
				}).Get("https://blockstream.info/api/address/" + btc.address + "/txs")
			}
			if err != nil || resp.StatusCode() != http.StatusOK {
				log.Println("Blockstream API : Error Getting BTC TX for", btc.address)
				break
			}
			err = json.Unmarshal(resp.Body(), &apiTXs)
			if err != nil {
				log.Println("Blockstream API : Error Unmarshaling BTC TX for", btc.address)
				break
			}
			if useCache {
				err = db.Write("BlockStream/address/txs", btc.address, apiTXs)
				if err != nil {
					log.Println("Blockstream API : Error Caching", btc.address)
				}
			}
		}
		err = nil
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
			for _, vin := range tx.Vin {
				if blkst.ownAddress(vin.Prevout.ScriptpubkeyAddress) {
					valueIn -= vin.Prevout.Value
					isInVinPrevVout = true
				} else {
					if missing == "" {
						missing = " missing :"
					}
					missing += " " + vin.Prevout.ScriptpubkeyAddress
				}
			}
			valueOut := 0
			isInVout := false
			dest := ""
			for _, vout := range tx.Vout {
				if blkst.ownAddress(vout.ScriptpubkeyAddress) {
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
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: "BTC", Amount: decimal.New(int64(-valueIn), -8)})
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "BTC", Amount: decimal.New(int64(valueOut), -8)})
					blkst.TXsByCategory["Transfers"] = append(blkst.TXsByCategory["Transfers"], t)
				} else if is, desc, val, curr := blkst.isTxPayment(tx.Txid); is {
					t.Note += " payment " + desc
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: "BTC", Amount: decimal.New(int64(-valueOut-valueIn-tx.Fee), -8)})
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr, Amount: val})
					blkst.TXsByCategory["CashOut"] = append(blkst.TXsByCategory["CashOut"], t)
				} else {
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: "BTC", Amount: decimal.New(int64(-valueOut-valueIn-tx.Fee), -8)})
					blkst.TXsByCategory["Withdrawals"] = append(blkst.TXsByCategory["Withdrawals"], t)
				}
				blkst.apiTXs[i].used = true
			} else if isInVout {
				t := wallet.TX{Timestamp: time.Unix(int64(tx.Status.BlockTime), 0), Note: "Blockstream API : " + strconv.Itoa(tx.Status.BlockHeight) + " " + tx.Txid}
				t.Items = make(map[string][]wallet.Currency)
				t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "BTC", Amount: decimal.New(int64(tx.Fee), -8)})
				t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "BTC", Amount: decimal.New(int64(valueOut), -8)})
				blkst.TXsByCategory["Deposits"] = append(blkst.TXsByCategory["Deposits"], t)
				blkst.apiTXs[i].used = true
			} else {
				log.Println("Blockstream API : Unmanaged TX")
				spew.Dump(tx)
			}
		}
	}
	return
}
