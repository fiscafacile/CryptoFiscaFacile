package etherscan

import (
	"errors"
	"time"

	"github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type nftTX struct {
	used              bool
	BlockNumber       int
	TimeStamp         time.Time
	Hash              string
	Nonce             int
	BlockHash         string
	From              string
	ContractAddress   string
	To                string
	TokenID           string
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

func (api *api) getNftTXs(addresses []string) {
	for _, eth := range addresses {
		accNftTX, err := api.getAccountNftTX(eth, "", false)
		if err != nil {
			api.doneNft <- err
			return
		}
		for _, nfTX := range accNftTX.Result {
			tx := nftTX{
				used:              false,
				BlockNumber:       nfTX.BlockNumber,
				TimeStamp:         time.Unix(nfTX.TimeStamp, 0),
				Hash:              nfTX.Hash,
				Nonce:             nfTX.Nonce,
				BlockHash:         nfTX.BlockHash,
				From:              nfTX.From,
				ContractAddress:   nfTX.ContractAddress,
				To:                nfTX.To,
				TokenID:           nfTX.TokenID,
				TokenName:         nfTX.TokenName,
				TokenSymbol:       nfTX.TokenSymbol,
				TokenDecimal:      nfTX.TokenDecimal,
				TransactionIndex:  nfTX.TransactionIndex,
				Gas:               nfTX.Gas,
				GasPrice:          decimal.NewFromBigInt(nfTX.GasPrice.Int(), -18),
				GasUsed:           decimal.NewFromInt(nfTX.GasUsed),
				CumulativeGasUsed: nfTX.CumulativeGasUsed,
				Input:             nfTX.Input,
				Confirmations:     nfTX.Confirmations,
			}
			api.nftTXs = append(api.nftTXs, tx)
		}
	}
	api.doneNft <- nil
}

type ResultNftTX struct {
	BlockNumber       int     `json:"blockNumber,string"`
	TimeStamp         int64   `json:"timeStamp,string"`
	Hash              string  `json:"hash"`
	Nonce             int     `json:"nonce,string"`
	BlockHash         string  `json:"blockHash"`
	From              string  `json:"from"`
	ContractAddress   string  `json:"contractAddress"`
	To                string  `json:"to"`
	TokenID           string  `json:"tokenID"`
	TokenName         string  `json:"tokenName"`
	TokenSymbol       string  `json:"tokenSymbol"`
	TokenDecimal      uint8   `json:"tokenDecimal,string"`
	TransactionIndex  int     `json:"transactionIndex,string"`
	Gas               int     `json:"gas,string"`
	GasPrice          *bigInt `json:"gasPrice"`
	GasUsed           int64   `json:"gasUsed,string"`
	CumulativeGasUsed int     `json:"cumulativeGasUsed,string"`
	Input             string  `json:"input"`
	Confirmations     int     `json:"confirmations,string"`
}

type GetAccountNftTXResp struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Result  []ResultNftTX `json:"result"`
}

func (api *api) getAccountNftTX(address, contractAddress string, desc bool) (accNftTX GetAccountNftTXResp, err error) {
	ident := "a" + address + "-c" + contractAddress
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Etherscan.io/account/tokennfttx", ident, &accNftTX)
	}
	if !useCache || err != nil {
		params := map[string]string{
			"module": "account",
			"action": "tokennfttx",
			"apikey": api.apiKey,
		}
		if address != "" {
			params["address"] = address
		}
		if contractAddress != "" {
			params["contractaddress"] = contractAddress
		}
		if desc {
			params["sort"] = "desc"
		} else {
			params["sort"] = "asc"
		}
		resp, err := api.clientNft.R().
			SetQueryParams(params).
			SetHeader("Accept", "application/json").
			SetResult(&GetAccountNftTXResp{}).
			Get(api.basePath)
		if err != nil {
			return accNftTX, errors.New("Etherscan API Nft TX : Error Requesting" + ident)
		}
		accNftTX = *resp.Result().(*GetAccountNftTXResp)
		if useCache {
			err = db.Write("Etherscan.io/account/tokennfttx", ident, accNftTX)
			if err != nil {
				return accNftTX, errors.New("Etherscan API Nft TX : Error Caching" + ident)
			}
		}
		if accNftTX.Message == "OK" {
			time.Sleep(api.timeBetweenReq)
		}
	}
	return accNftTX, nil
}
