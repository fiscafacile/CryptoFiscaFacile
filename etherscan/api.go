package etherscan

import (
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/category"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
)

type api struct {
	clientNor      *resty.Client
	doneNor        chan error
	clientInt      *resty.Client
	doneInt        chan error
	clientTok      *resty.Client
	doneTok        chan error
	clientNft      *resty.Client
	doneNft        chan error
	basePath       string
	apiKey         string
	firstTimeUsed  time.Time
	timeBetweenReq time.Duration
	normalTXs      []normalTX
	internalTXs    []internalTX
	tokenTXs       []tokenTX
	nftTXs         []nftTX
	txsByCategory  wallet.TXsByCategory
}

func (ethsc *Etherscan) NewAPI(apiKey string, debug bool) {
	ethsc.api.txsByCategory = make(map[string]wallet.TXs)
	ethsc.api.clientNor = resty.New()
	ethsc.api.clientNor.SetRetryCount(3)
	ethsc.api.clientNor.SetDebug(debug)
	ethsc.api.doneNor = make(chan error)
	ethsc.api.clientInt = resty.New()
	ethsc.api.clientInt.SetRetryCount(3)
	ethsc.api.clientInt.SetDebug(debug)
	ethsc.api.doneInt = make(chan error)
	ethsc.api.clientTok = resty.New()
	ethsc.api.clientTok.SetRetryCount(3).SetRetryWaitTime(1 * time.Second)
	ethsc.api.clientTok.SetDebug(debug)
	ethsc.api.doneTok = make(chan error)
	ethsc.api.clientNft = resty.New()
	ethsc.api.clientNft.SetRetryCount(3).SetRetryWaitTime(1 * time.Second)
	ethsc.api.clientNft.SetDebug(debug)
	ethsc.api.doneNft = make(chan error)
	ethsc.api.basePath = "https://api.etherscan.io/api"
	ethsc.api.apiKey = apiKey
	ethsc.api.firstTimeUsed = time.Now()
	ethsc.api.timeBetweenReq = 200 * time.Millisecond
}

func (api *api) getAllTXs(addresses []string, cat category.Category) (err error) {
	go api.getNormalTXs(addresses)
	go api.getInternalTXs(addresses)
	go api.getTokenTXs(addresses)
	go api.getNftTXs(addresses)
	<-api.doneNor
	<-api.doneInt
	<-api.doneTok
	<-api.doneNft
	api.categorize(addresses, cat)
	return
}

func (api *api) GetExchangeFirstUsedTime() time.Time {
	return api.firstTimeUsed
}

func (api *api) categorize(addresses []string, cat category.Category) {
	alreadyAsked := []string{}
	for i, tx := range api.nftTXs {
		if !tx.used {
			if is, _, _, _ := cat.IsTxShit(tx.Hash); is {
				api.nftTXs[i].used = true
			} else {
				if api.ownAddress(tx.To, addresses) && api.ownAddress(tx.From, addresses) {
					alreadyAsked = wallet.AskForHelp("Etherscan API ERC721 Self TX", tx, alreadyAsked)
				} else if api.ownAddress(tx.To, addresses) || api.ownAddress(tx.From, addresses) {
					t := wallet.TX{Timestamp: tx.TimeStamp, ID: tx.Hash, Note: "Etherscan API : " + strconv.Itoa(tx.BlockNumber) + " " + tx.To}
					t.Items = make(map[string]wallet.Currencies)
					t.Nfts = make(map[string]wallet.Nfts)
					if api.ownAddress(tx.From, addresses) {
						if !tx.GasPrice.IsZero() && !tx.GasUsed.IsZero() {
							t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: tx.GasPrice.Mul(tx.GasUsed)})
						}
						t.Nfts["From"] = append(t.Nfts["From"], wallet.Nft{ID: tx.TokenID, Name: tx.TokenName, Symbol: tx.TokenSymbol})
					} else {
						t.Nfts["To"] = append(t.Nfts["To"], wallet.Nft{ID: tx.TokenID, Name: tx.TokenName, Symbol: tx.TokenSymbol})
					}
					for j, ntx := range api.normalTXs {
						if ntx.Hash == tx.Hash {
							if !ntx.GasPrice.IsZero() && !ntx.GasUsed.IsZero() {
								t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: ntx.GasPrice.Mul(ntx.GasUsed)})
							}
							api.normalTXs[j].used = true
							break
						}
					}
					if is, feeHash := cat.IsTxFee(tx.Hash); is {
						for k, ntx2 := range api.normalTXs {
							for _, fee := range strings.Split(feeHash, ";") {
								if ntx2.Hash == fee {
									if !ntx2.GasPrice.IsZero() && !ntx2.GasUsed.IsZero() &&
										(!ntx2.GasPrice.Equal(tx.GasPrice) || !ntx2.GasUsed.Equal(tx.GasUsed)) {
										t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: ntx2.GasPrice.Mul(ntx2.GasUsed)})
									}
									api.normalTXs[k].used = true
								}
							}
						}
					}
					api.txsByCategory["NFTs"] = append(api.txsByCategory["NFTs"], t)
					api.nftTXs[i].used = true
				} else {
					alreadyAsked = wallet.AskForHelp("Etherscan API ERC721 TX", tx, alreadyAsked)
				}
			}
		}
	}
	for i, tx := range api.tokenTXs {
		if !tx.used {
			if is, _, _, _ := cat.IsTxShit(tx.Hash); is {
				api.tokenTXs[i].used = true
			} else {
				if api.ownAddress(tx.To, addresses) && api.ownAddress(tx.From, addresses) {
					found := false
					for _, tr := range api.txsByCategory["Transfers"] {
						if tx.Hash == tr.ID {
							found = true
						}
					}
					if !found {
						t := wallet.TX{Timestamp: tx.TimeStamp, ID: tx.Hash, Note: "Etherscan API : " + strconv.Itoa(tx.BlockNumber) + " " + tx.To}
						t.Items = make(map[string]wallet.Currencies)
						t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.TokenSymbol, Amount: tx.Value})
						t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.TokenSymbol, Amount: tx.Value})
						t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: tx.GasPrice.Mul(tx.GasUsed)})
						api.txsByCategory["Transfers"] = append(api.txsByCategory["Transfers"], t)
					}
					api.tokenTXs[i].used = true
				} else if api.ownAddress(tx.To, addresses) {
					t := wallet.TX{Timestamp: tx.TimeStamp, ID: tx.Hash, Note: "Etherscan API : " + strconv.Itoa(tx.BlockNumber) + " " + tx.To}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: tx.TokenSymbol, Amount: tx.Value})
					api.tokenTXs[i].used = true
					// Add declared Fee if any
					if is, feeHash := cat.IsTxFee(tx.Hash); is {
						for k, ntx2 := range api.normalTXs {
							for _, fee := range strings.Split(feeHash, ";") {
								if ntx2.Hash == fee {
									api.normalTXs[k].used = true
									t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: ntx2.GasPrice.Mul(ntx2.GasUsed)})
								}
							}
						}
					}
					// Add normal Fee if any
					found := false
					for j, ntx := range api.normalTXs {
						if ntx.Hash == tx.Hash {
							found = true
							t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: ntx.GasPrice.Mul(ntx.GasUsed)})
							if tx.From == "0x0000000000000000000000000000000000000000" {
								found2 := false
								if is, buyHash := cat.IsTxTokenSale(tx.Hash); is {
									for k, ntx2 := range api.normalTXs {
										if ntx2.Hash == buyHash {
											found2 = true
											api.normalTXs[k].used = true
											t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: "ETH", Amount: ntx2.Value})
											t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: ntx2.GasPrice.Mul(ntx2.GasUsed)})
										}
									}
								}
								if found2 {
									api.txsByCategory["TokenBuys"] = append(api.txsByCategory["TokenBuys"], t)
								} else {
									api.txsByCategory["Claims"] = append(api.txsByCategory["Claims"], t)
								}
							} else {
								if ntx.Value.IsZero() {
									if is, desc := cat.IsTxAirDrop(tx.Hash); is {
										t.Note += " " + desc
										api.txsByCategory["AirDrops"] = append(api.txsByCategory["AirDrops"], t)
									} else {
										api.txsByCategory["Deposits"] = append(api.txsByCategory["Deposits"], t)
									}
								} else {
									t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: "ETH", Amount: ntx.Value})
									api.txsByCategory["Swaps"] = append(api.txsByCategory["Swaps"], t)
								}
							}
							api.normalTXs[j].used = true
							api.tokenTXs[i].used = true
							break
						}
					}
					if !found {
						// Look for other tokenTX with same Hash
						for j, ttx := range api.tokenTXs {
							if !ttx.used {
								if ttx.Hash == tx.Hash {
									api.tokenTXs[j].used = true
									if api.ownAddress(ttx.To, addresses) && api.ownAddress(ttx.From, addresses) {
										alreadyAsked = wallet.AskForHelp("Etherscan API ERC20 Self1 TX", ttx, alreadyAsked)
									} else if api.ownAddress(ttx.To, addresses) {
										t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: ttx.TokenSymbol, Amount: ttx.Value})
									} else if api.ownAddress(ttx.From, addresses) {
										found = true
										t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: ttx.TokenSymbol, Amount: ttx.Value})
									}
								}
							}
						}
						if found {
							api.txsByCategory["Swaps"] = append(api.txsByCategory["Swaps"], t)
						} else {
							api.txsByCategory["Deposits"] = append(api.txsByCategory["Deposits"], t)
						}
					}
					for k, ntx := range api.normalTXs {
						if ntx.Hash == tx.Hash {
							api.normalTXs[k].used = true
							break
						}
					}
				} else if api.ownAddress(tx.From, addresses) {
					t := wallet.TX{Timestamp: tx.TimeStamp, ID: tx.Hash, Note: "Etherscan API : " + strconv.Itoa(tx.BlockNumber) + " " + tx.To}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: tx.TokenSymbol, Amount: tx.Value})
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: tx.GasPrice.Mul(tx.GasUsed)})
					api.tokenTXs[i].used = true
					// Add declared Fee if any
					if is, feeHash := cat.IsTxFee(tx.Hash); is {
						for k, ntx2 := range api.normalTXs {
							for _, fee := range strings.Split(feeHash, ";") {
								if ntx2.Hash == fee {
									api.normalTXs[k].used = true
									t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: ntx2.GasPrice.Mul(ntx2.GasUsed)})
								}
							}
						}
					}
					if tx.To == "0x0000000000000000000000000000000000000000" {
						api.txsByCategory["Burns"] = append(api.txsByCategory["Burns"], t)
						api.tokenTXs[i].used = true
					} else {
						found := false
						for j, itx := range api.internalTXs {
							if itx.TimeStamp.Equal(tx.TimeStamp) &&
								itx.BlockNumber == tx.BlockNumber &&
								itx.Hash == tx.Hash {
								found = true
								t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "ETH", Amount: itx.Value})
								api.txsByCategory["Swaps"] = append(api.txsByCategory["Swaps"], t)
								api.internalTXs[j].used = true
								api.tokenTXs[i].used = true
								break
							}
						}
						if !found {
							// Look for other tokenTX with same Hash
							for j, ttx := range api.tokenTXs {
								if !ttx.used && ttx.Hash == tx.Hash {
									api.tokenTXs[j].used = true
									if api.ownAddress(ttx.To, addresses) && api.ownAddress(ttx.From, addresses) {
										alreadyAsked = wallet.AskForHelp("Etherscan API ERC20 Self2 TX", ttx, alreadyAsked)
									} else if api.ownAddress(ttx.To, addresses) {
										found = true
										t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: ttx.TokenSymbol, Amount: ttx.Value})
									} else if api.ownAddress(ttx.From, addresses) {
										t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: ttx.TokenSymbol, Amount: ttx.Value})
									}
								}
							}
							if found {
								api.txsByCategory["Swaps"] = append(api.txsByCategory["Swaps"], t)
							} else {
								api.txsByCategory["Withdrawals"] = append(api.txsByCategory["Withdrawals"], t)
							}
						}
					}
					for k, ntx := range api.normalTXs {
						if ntx.Hash == tx.Hash {
							api.normalTXs[k].used = true
							break
						}
					}
				} else {
					alreadyAsked = wallet.AskForHelp("Etherscan API ERC20 TX", tx, alreadyAsked)
				}
			}
		}
	}
	for i, tx := range api.internalTXs {
		if !tx.used {
			if is, _, _, _ := cat.IsTxShit(tx.Hash); is {
				api.internalTXs[i].used = true
			} else {
				if api.ownAddress(tx.To, addresses) && api.ownAddress(tx.From, addresses) {
					alreadyAsked = wallet.AskForHelp("Etherscan API Internal Self TX", tx, alreadyAsked)
				} else if api.ownAddress(tx.To, addresses) {
					t := wallet.TX{Timestamp: tx.TimeStamp, ID: tx.Hash, Note: "Etherscan API : " + strconv.Itoa(tx.BlockNumber) + " " + tx.From}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "ETH", Amount: tx.Value})
					if is, feeHash := cat.IsTxFee(tx.Hash); is {
						for k, ntx2 := range api.normalTXs {
							for _, fee := range strings.Split(feeHash, ";") {
								if ntx2.Hash == fee {
									api.normalTXs[k].used = true
									t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: ntx2.GasPrice.Mul(ntx2.GasUsed)})
								}
							}
						}
					}
					isExchange := false
					for j, ntx := range api.normalTXs {
						if ntx.Hash == tx.Hash {
							t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: ntx.GasPrice.Mul(ntx.GasUsed)})
							if !ntx.Value.IsZero() {
								if api.ownAddress(ntx.From, addresses) {
									t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: "ETH", Amount: ntx.Value})
									isExchange = true
								} else {
									alreadyAsked = wallet.AskForHelp("Etherscan API Internal Deposits TX with Normal Deposits TX associated", tx, alreadyAsked)
								}
							}
							api.normalTXs[j].used = true
							break
						}
					}
					if isExchange {
						api.txsByCategory["Exchanges"] = append(api.txsByCategory["Exchanges"], t)
					} else if is, desc, val, curr := cat.IsTxExchange(tx.Hash); is {
						t.Note += " crypto_exchange " + desc
						t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: curr, Amount: val})
						api.txsByCategory["Exchanges"] = append(api.txsByCategory["Exchanges"], t)
					} else {
						api.txsByCategory["Deposits"] = append(api.txsByCategory["Deposits"], t)
					}
					api.internalTXs[i].used = true
				} else if api.ownAddress(tx.From, addresses) {
					alreadyAsked = wallet.AskForHelp("Etherscan API Internal Withdrawal TX", tx, alreadyAsked)
				} else {
					alreadyAsked = wallet.AskForHelp("Etherscan API Internal TX", tx, alreadyAsked)
				}
			}
		}
		if tx.TimeStamp.Before(api.firstTimeUsed) {
			api.firstTimeUsed = tx.TimeStamp
		}
	}
	for i, tx := range api.normalTXs {
		if !tx.used {
			if is, _, _, _ := cat.IsTxShit(tx.Hash); is {
				api.normalTXs[i].used = true
			} else {
				if api.ownAddress(tx.To, addresses) && api.ownAddress(tx.From, addresses) {
					t := wallet.TX{Timestamp: tx.TimeStamp, ID: tx.Hash, Note: "Etherscan API : " + strconv.Itoa(tx.BlockNumber) + " "}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: tx.GasPrice.Mul(tx.GasUsed)})
					if tx.To == tx.From || tx.IsError != 0 {
						api.txsByCategory["Fees"] = append(api.txsByCategory["Fees"], t)
						api.normalTXs[i].used = true
					} else {
						t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "ETH", Amount: tx.Value})
						t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: "ETH", Amount: tx.Value})
						api.txsByCategory["Transfers"] = append(api.txsByCategory["Transfers"], t)
						api.normalTXs[i].used = true
						for j, ntx := range api.normalTXs {
							if ntx.Hash == tx.Hash &&
								!ntx.used {
								api.normalTXs[j].used = true
								break
							}
						}
					}
				} else if api.ownAddress(tx.To, addresses) {
					if !tx.Value.IsZero() && tx.IsError == 0 {
						t := wallet.TX{Timestamp: tx.TimeStamp, ID: tx.Hash, Note: "Etherscan API : " + strconv.Itoa(tx.BlockNumber) + " " + tx.From}
						t.Items = make(map[string]wallet.Currencies)
						t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "ETH", Amount: tx.Value})
						if is, desc, val, curr := cat.IsTxExchange(tx.Hash); is {
							t.Note += " crypto_exchange " + desc
							t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: curr, Amount: val})
							api.txsByCategory["Exchanges"] = append(api.txsByCategory["Exchanges"], t)
						} else {
							api.txsByCategory["Deposits"] = append(api.txsByCategory["Deposits"], t)
						}
						api.normalTXs[i].used = true
					}
				} else if api.ownAddress(tx.From, addresses) {
					t := wallet.TX{Timestamp: tx.TimeStamp, ID: tx.Hash, Note: "Etherscan API : " + strconv.Itoa(tx.BlockNumber) + " " + tx.To}
					t.Items = make(map[string]wallet.Currencies)
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: "ETH", Amount: tx.GasPrice.Mul(tx.GasUsed)})
					if !tx.Value.IsZero() && tx.IsError == 0 {
						t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: "ETH", Amount: tx.Value})
						// Is declared Exchanges
						if is, desc, val, curr := cat.IsTxExchange(tx.Hash); is {
							t.Note += " crypto_exchange " + desc
							t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: curr, Amount: val})
							api.txsByCategory["Exchanges"] = append(api.txsByCategory["Exchanges"], t)
						} else {
							if tx.To == "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2" { // Special Case WETH
								t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "WETH", Amount: tx.Value})
								api.txsByCategory["Wraps"] = append(api.txsByCategory["Wraps"], t)
							} else if tx.To == "0xf786c34106762ab4eeb45a51b42a62470e9d5332" { // Special Case fETH
								t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "fETH", Amount: tx.Value.Mul(decimal.New(99, -2))})
								api.txsByCategory["Wraps"] = append(api.txsByCategory["Wraps"], t)
							} else {
								api.txsByCategory["Withdrawals"] = append(api.txsByCategory["Withdrawals"], t)
							}
						}
						api.normalTXs[i].used = true
					} else {
						api.txsByCategory["Fees"] = append(api.txsByCategory["Fees"], t)
						api.normalTXs[i].used = true
					}
				} else {
					alreadyAsked = wallet.AskForHelp("Etherscan API TX", tx, alreadyAsked)
				}
			}
		}
		if tx.TimeStamp.Before(api.firstTimeUsed) {
			api.firstTimeUsed = tx.TimeStamp
		}
	}
}

func (api *api) ownAddress(add string, addresses []string) bool {
	for _, a := range addresses {
		if a == add {
			return true
		}
	}
	return false
}

// BigInt is a wrapper over big.Int to implement only unmarshalText
// for json decoding.
type bigInt big.Int

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (b *bigInt) UnmarshalText(text []byte) (err error) {
	var bi = new(big.Int)
	err = bi.UnmarshalText(text)
	if err != nil {
		return
	}

	*b = bigInt(*bi)
	return nil
}

// MarshalText implements the encoding.TextMarshaler
func (b *bigInt) MarshalText() (text []byte, err error) {
	return []byte(b.Int().String()), nil
}

// Int returns b's *big.Int form
func (b *bigInt) Int() *big.Int {
	return (*big.Int)(b)
}
