package etherscan

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fiscafacile/CryptoFiscaFacile/category"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/nanmu42/etherscan-api"
	"github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type API struct {
	client *etherscan.Client
}

type apiNormalTX struct {
	used              bool
	BlockNumber       int
	TimeStamp         time.Time
	Hash              string
	Nonce             int
	BlockHash         string
	TransactionIndex  int
	From              string
	To                string
	Value             decimal.Decimal
	Gas               int
	GasPrice          decimal.Decimal
	IsError           int
	TxReceiptStatus   string
	Input             string
	ContractAddress   string
	CumulativeGasUsed int
	GasUsed           decimal.Decimal
	Confirmations     int
}

type apiInternalTX struct {
	used            bool
	BlockNumber     int
	TimeStamp       time.Time
	Hash            string
	From            string
	To              string
	Value           decimal.Decimal
	ContractAddress string
	Input           string
	Type            string
	Gas             int
	GasUsed         decimal.Decimal
	TraceID         string
	IsError         int
	ErrCode         string
}

type apiERC20TX struct {
	used              bool
	BlockNumber       int
	TimeStamp         time.Time
	Hash              string
	Nonce             int
	BlockHash         string
	From              string
	ContractAddress   string
	To                string
	Value             decimal.Decimal
	TokenName         string
	TokenSymbol       string
	TokenDecimal      uint8
	TransactionIndex  int
	Gas               int
	GasPrice          decimal.Decimal
	GasUsed           decimal.Decimal
	CumulativeGasUsed int
	Input             string
	Confirmations     int
}

func (ethsc *Etherscan) APIConnect(apikey string) {
	ethsc.api.client = etherscan.New(etherscan.Mainnet, apikey)
}

func (ethsc *Etherscan) apiGetAllTXs(cat category.Category) (err error) {
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	for _, eth := range ethsc.csvAddresses {
		var nTXs []etherscan.NormalTx
		if useCache {
			err = db.Read("EtherScan/accounts/txlist", eth.address, &nTXs)
		}
		if !useCache || err != nil {
			nTXs, err = ethsc.api.client.NormalTxByAddress(eth.address, nil, nil, 1, 500, false)
			if err != nil {
				time.Sleep(6 * time.Second)
				nTXs, err = ethsc.api.client.NormalTxByAddress(eth.address, nil, nil, 1, 500, false)
				if err != nil {
					log.Println("Etherscan API : Error Getting by API Normal TX for", eth.address)
				}
			}
			if useCache {
				err = db.Write("EtherScan/accounts/txlist", eth.address, nTXs)
				if err != nil {
					log.Println("Etherscan API : Error Caching Normal TX for", eth.address)
				}
			}
		}
		for _, t := range nTXs {
			tx := apiNormalTX{
				used:              false,
				BlockNumber:       t.BlockNumber,
				TimeStamp:         t.TimeStamp.Time(),
				Hash:              t.Hash,
				Nonce:             t.Nonce,
				BlockHash:         t.BlockHash,
				TransactionIndex:  t.TransactionIndex,
				From:              t.From,
				To:                t.To,
				Value:             decimal.NewFromBigInt(t.Value.Int(), -18),
				Gas:               t.Gas,
				GasPrice:          decimal.NewFromBigInt(t.GasPrice.Int(), -18),
				IsError:           t.IsError,
				TxReceiptStatus:   t.TxReceiptStatus,
				Input:             t.Input,
				ContractAddress:   t.ContractAddress,
				CumulativeGasUsed: t.CumulativeGasUsed,
				GasUsed:           decimal.NewFromInt(int64(t.GasUsed)),
				Confirmations:     t.Confirmations,
			}
			ethsc.apiNormalTXs = append(ethsc.apiNormalTXs, tx)
		}
		var iTXs []etherscan.InternalTx
		if useCache {
			err = db.Read("EtherScan/accounts/txlistinternal", eth.address, &iTXs)
		}
		if !useCache || err != nil {
			iTXs, err = ethsc.api.client.InternalTxByAddress(eth.address, nil, nil, 1, 500, false)
			if err != nil {
				time.Sleep(6 * time.Second)
				iTXs, err = ethsc.api.client.InternalTxByAddress(eth.address, nil, nil, 1, 500, false)
				if err != nil {
					log.Println("Etherscan API : Error Getting by API Internal TX for", eth.address)
				}
			}
			if useCache {
				err = db.Write("EtherScan/accounts/txlistinternal", eth.address, iTXs)
				if err != nil {
					log.Println("Etherscan API : Error Caching Internal TX for", eth.address)
				}
			}
		}
		for _, t := range iTXs {
			tx := apiInternalTX{
				used:            false,
				BlockNumber:     t.BlockNumber,
				TimeStamp:       t.TimeStamp.Time(),
				Hash:            t.Hash,
				From:            t.From,
				To:              t.To,
				Value:           decimal.NewFromBigInt(t.Value.Int(), -18),
				ContractAddress: t.ContractAddress,
				Input:           t.Input,
				Type:            t.Type,
				Gas:             t.Gas,
				GasUsed:         decimal.NewFromInt(int64(t.GasUsed)),
				TraceID:         t.TraceID,
				IsError:         t.IsError,
				ErrCode:         t.ErrCode,
			}
			ethsc.apiInternalTXs = append(ethsc.apiInternalTXs, tx)
		}
		var tTXs []etherscan.ERC20Transfer
		if useCache {
			err = db.Read("EtherScan/accounts/tokentx", eth.address, &tTXs)
		}
		if !useCache || err != nil {
			tTXs, err = ethsc.api.client.ERC20Transfers(nil, &eth.address, nil, nil, 1, 500, false)
			if err != nil {
				time.Sleep(6 * time.Second)
				tTXs, err = ethsc.api.client.ERC20Transfers(nil, &eth.address, nil, nil, 1, 500, false)
				if err != nil {
					log.Println("Etherscan API : Error Getting by API ERC20 TX for", eth.address)
				}
			}
			if useCache {
				err = db.Write("EtherScan/accounts/tokentx", eth.address, tTXs)
				if err != nil {
					log.Println("Etherscan API : Error Caching ERC20 TX for", eth.address)
				}
			}
		}
		for _, t := range tTXs {
			tx := apiERC20TX{
				used:              false,
				BlockNumber:       t.BlockNumber,
				TimeStamp:         t.TimeStamp.Time(),
				Hash:              t.Hash,
				Nonce:             t.Nonce,
				BlockHash:         t.BlockHash,
				From:              t.From,
				ContractAddress:   t.ContractAddress,
				To:                t.To,
				Value:             decimal.NewFromBigInt(t.Value.Int(), -18),
				TokenName:         t.TokenName,
				TokenSymbol:       t.TokenSymbol,
				TokenDecimal:      t.TokenDecimal,
				TransactionIndex:  t.TransactionIndex,
				Gas:               t.Gas,
				GasPrice:          decimal.NewFromBigInt(t.GasPrice.Int(), -18),
				GasUsed:           decimal.NewFromInt(int64(t.GasUsed)),
				CumulativeGasUsed: t.CumulativeGasUsed,
				Input:             t.Input,
				Confirmations:     t.Confirmations,
			}
			ethsc.apiERC20TXs = append(ethsc.apiERC20TXs, tx)
		}
	}
	ethsc.apiFillTXsByCategory(cat)
	return
}

func (ethsc *Etherscan) apiFillTXsByCategory(cat category.Category) {
	for i, tx := range ethsc.apiERC20TXs {
		if !tx.used {
			if is, _ := cat.IsTxShit(tx.Hash); is {
				ethsc.apiERC20TXs[i].used = true
			} else {
				if ethsc.ownAddress(tx.To) && ethsc.ownAddress(tx.From) {
					log.Println("Detected Self ERC20 TX", tx)
				} else if ethsc.ownAddress(tx.To) {
					t := wallet.TX{Timestamp: tx.TimeStamp, Note: "Etherscan API : " + strconv.Itoa(tx.BlockNumber) + " " + tx.Hash + " " + tx.To}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.TokenSymbol, Amount: tx.Value})
					if is, feeHash := cat.IsTxFee(tx.Hash); is {
						for k, ntx2 := range ethsc.apiNormalTXs {
							for _, fee := range strings.Split(feeHash, ";") {
								if ntx2.Hash == fee {
									ethsc.apiNormalTXs[k].used = true
									t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: ntx2.GasPrice.Mul(ntx2.GasUsed)})
								}
							}
						}
					}
					found := false
					for j, ntx := range ethsc.apiNormalTXs {
						if ntx.Hash == tx.Hash {
							found = true
							t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: tx.GasPrice.Mul(tx.GasUsed)})
							if tx.From == "0x0000000000000000000000000000000000000000" {
								found2 := false
								if is, buyHash := cat.IsTxTokenSale(tx.Hash); is {
									for k, ntx2 := range ethsc.apiNormalTXs {
										if ntx2.Hash == buyHash {
											found2 = true
											ethsc.apiNormalTXs[k].used = true
											t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: "ETH", Amount: ntx2.Value})
											t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: ntx2.GasPrice.Mul(ntx2.GasUsed)})
										}
									}
								}
								if found2 {
									ethsc.TXsByCategory["TokenBuys"] = append(ethsc.TXsByCategory["TokenBuys"], t)
								} else {
									ethsc.TXsByCategory["Claims"] = append(ethsc.TXsByCategory["Claims"], t)
								}
							} else {
								if ntx.Value.IsZero() {
									ethsc.TXsByCategory["Deposits"] = append(ethsc.TXsByCategory["Deposits"], t)
								} else {
									t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: "ETH", Amount: ntx.Value})
									ethsc.TXsByCategory["Swaps"] = append(ethsc.TXsByCategory["Swaps"], t)
								}
							}
							ethsc.apiNormalTXs[j].used = true
							ethsc.apiERC20TXs[i].used = true
							break
						}
					}
					if !found {
						ethsc.TXsByCategory["Deposits"] = append(ethsc.TXsByCategory["Deposits"], t)
						ethsc.apiERC20TXs[i].used = true
					}
					for k, ntx := range ethsc.apiNormalTXs {
						if ntx.TimeStamp.Equal(tx.TimeStamp) &&
							ntx.BlockNumber == tx.BlockNumber &&
							ntx.Hash == tx.Hash {
							ethsc.apiNormalTXs[k].used = true
							break
						}
					}
				} else if ethsc.ownAddress(tx.From) {
					t := wallet.TX{Timestamp: tx.TimeStamp, Note: "Etherscan API : " + strconv.Itoa(tx.BlockNumber) + " " + tx.Hash + " " + tx.To}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.TokenSymbol, Amount: tx.Value})
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: tx.GasPrice.Mul(tx.GasUsed)})
					if is, feeHash := cat.IsTxFee(tx.Hash); is {
						for k, ntx2 := range ethsc.apiNormalTXs {
							for _, fee := range strings.Split(feeHash, ";") {
								if ntx2.Hash == fee {
									ethsc.apiNormalTXs[k].used = true
									t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: ntx2.GasPrice.Mul(ntx2.GasUsed)})
								}
							}
						}
					}
					if tx.To == "0x0000000000000000000000000000000000000000" {
						ethsc.TXsByCategory["Burns"] = append(ethsc.TXsByCategory["Burns"], t)
						ethsc.apiERC20TXs[i].used = true
					} else {
						found := false
						for j, itx := range ethsc.apiInternalTXs {
							if itx.TimeStamp.Equal(tx.TimeStamp) &&
								itx.BlockNumber == tx.BlockNumber &&
								itx.Hash == tx.Hash {
								found = true
								t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "ETH", Amount: itx.Value})
								ethsc.TXsByCategory["Swaps"] = append(ethsc.TXsByCategory["Swaps"], t)
								ethsc.apiInternalTXs[j].used = true
								ethsc.apiERC20TXs[i].used = true
								break
							}
						}
						if !found {
							ethsc.TXsByCategory["Withdrawals"] = append(ethsc.TXsByCategory["Withdrawals"], t)
							ethsc.apiERC20TXs[i].used = true
						}
					}
					for k, ntx := range ethsc.apiNormalTXs {
						if ntx.TimeStamp.Equal(tx.TimeStamp) &&
							ntx.BlockNumber == tx.BlockNumber &&
							ntx.Hash == tx.Hash {
							ethsc.apiNormalTXs[k].used = true
							break
						}
					}
				} else {
					log.Println("Unmanaged ERC20 TX")
					spew.Dump(tx)
				}
			}
		}
	}
	for i, tx := range ethsc.apiInternalTXs {
		if !tx.used {
			if is, _ := cat.IsTxShit(tx.Hash); is {
				ethsc.apiInternalTXs[i].used = true
			} else {
				if ethsc.ownAddress(tx.To) && ethsc.ownAddress(tx.From) {
					log.Println("Detected Self Internal TX", tx)
				} else if ethsc.ownAddress(tx.To) {
					t := wallet.TX{Timestamp: tx.TimeStamp, Note: "Etherscan API : " + strconv.Itoa(tx.BlockNumber) + " " + tx.Hash + " " + tx.From}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "ETH", Amount: tx.Value})
					if is, feeHash := cat.IsTxFee(tx.Hash); is {
						for k, ntx2 := range ethsc.apiNormalTXs {
							for _, fee := range strings.Split(feeHash, ";") {
								if ntx2.Hash == fee {
									ethsc.apiNormalTXs[k].used = true
									t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: ntx2.GasPrice.Mul(ntx2.GasUsed)})
								}
							}
						}
					}
					for j, ntx := range ethsc.apiNormalTXs {
						if ntx.TimeStamp.Equal(tx.TimeStamp) &&
							ntx.BlockNumber == tx.BlockNumber &&
							ntx.Hash == tx.Hash {
							t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: ntx.GasPrice.Mul(ntx.GasUsed)})
							if !ntx.Value.IsZero() {
								log.Println("Detected Deposits Internal TX with non Zero Normal Value", tx)
							}
							ethsc.apiNormalTXs[j].used = true
							break
						}
					}
					ethsc.TXsByCategory["Deposits"] = append(ethsc.TXsByCategory["Deposits"], t)
					ethsc.apiInternalTXs[i].used = true
				} else if ethsc.ownAddress(tx.From) {
					log.Println("Detected Withdrawal Internal TX", tx)
				} else {
					log.Println("Unmanaged Internal TX")
					spew.Dump(tx)
				}
			}
		}
	}
	for i, tx := range ethsc.apiNormalTXs {
		if !tx.used {
			if is, _ := cat.IsTxShit(tx.Hash); is {
				ethsc.apiNormalTXs[i].used = true
			} else {
				if ethsc.ownAddress(tx.To) && ethsc.ownAddress(tx.From) {
					t := wallet.TX{Timestamp: tx.TimeStamp, Note: "Etherscan API : " + strconv.Itoa(tx.BlockNumber) + " " + tx.Hash + " "}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: tx.GasPrice.Mul(tx.GasUsed)})
					if tx.To == tx.From {
						if !tx.Value.IsZero() {
							log.Println("Detected non zero Value Self TX", tx)
						}
						ethsc.TXsByCategory["Fees"] = append(ethsc.TXsByCategory["Fees"], t)
						ethsc.apiNormalTXs[i].used = true
					} else {
						t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "ETH", Amount: tx.Value})
						t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: "ETH", Amount: tx.Value})
						ethsc.TXsByCategory["Transfers"] = append(ethsc.TXsByCategory["Transfers"], t)
						ethsc.apiNormalTXs[i].used = true
					}
				} else if ethsc.ownAddress(tx.To) {
					if !tx.Value.IsZero() {
						t := wallet.TX{Timestamp: tx.TimeStamp, Note: "Etherscan API : " + strconv.Itoa(tx.BlockNumber) + " " + tx.Hash + " " + tx.From}
						t.Items = make(map[string]wallet.Currencies)
						t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "ETH", Amount: tx.Value})
						ethsc.TXsByCategory["Deposits"] = append(ethsc.TXsByCategory["Deposits"], t)
						ethsc.apiNormalTXs[i].used = true
					}
				} else if ethsc.ownAddress(tx.From) {
					t := wallet.TX{Timestamp: tx.TimeStamp, Note: "Etherscan API : " + strconv.Itoa(tx.BlockNumber) + " " + tx.Hash + " " + tx.To}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: tx.GasPrice.Mul(tx.GasUsed)})
					if !tx.Value.IsZero() {
						t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: "ETH", Amount: tx.Value})
						ethsc.TXsByCategory["Withdrawals"] = append(ethsc.TXsByCategory["Withdrawals"], t)
						ethsc.apiNormalTXs[i].used = true
					} else {
						ethsc.TXsByCategory["Fees"] = append(ethsc.TXsByCategory["Fees"], t)
						ethsc.apiNormalTXs[i].used = true
					}
				} else {
					log.Println("Unmanaged Normal TX")
					spew.Dump(tx)
				}
			}
		}
	}
}
