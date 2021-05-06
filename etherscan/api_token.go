package etherscan

import (
	"errors"
	"time"

	"github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type tokenTX struct {
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

func (api *api) getTokenTXs(addresses []string) {
	for _, eth := range addresses {
		accTokTX, err := api.getAccountTokenTX(eth, "", false)
		if err != nil {
			api.doneTok <- err
			return
		}
		for _, tokTX := range accTokTX.Result {
			tx := tokenTX{
				used:              false,
				BlockNumber:       tokTX.BlockNumber,
				TimeStamp:         time.Unix(tokTX.TimeStamp, 0),
				Hash:              tokTX.Hash,
				Nonce:             tokTX.Nonce,
				BlockHash:         tokTX.BlockHash,
				From:              tokTX.From,
				ContractAddress:   tokTX.ContractAddress,
				To:                tokTX.To,
				Value:             decimal.NewFromBigInt(tokTX.Value.Int(), -int32(tokTX.TokenDecimal)),
				TokenName:         tokTX.TokenName,
				TokenSymbol:       tokTX.TokenSymbol,
				TokenDecimal:      tokTX.TokenDecimal,
				TransactionIndex:  tokTX.TransactionIndex,
				Gas:               tokTX.Gas,
				GasPrice:          decimal.NewFromBigInt(tokTX.GasPrice.Int(), -18),
				GasUsed:           decimal.NewFromInt(tokTX.GasUsed),
				CumulativeGasUsed: tokTX.CumulativeGasUsed,
				Input:             tokTX.Input,
				Confirmations:     tokTX.Confirmations,
			}
			api.tokenTXs = append(api.tokenTXs, tx)
		}
	}
	api.doneTok <- nil
}

type ResultTokenTX struct {
	BlockNumber       int     `json:"blockNumber,string"`
	TimeStamp         int64   `json:"timeStamp,string"`
	Hash              string  `json:"hash"`
	Nonce             int     `json:"nonce,string"`
	BlockHash         string  `json:"blockHash"`
	From              string  `json:"from"`
	ContractAddress   string  `json:"contractAddress"`
	To                string  `json:"to"`
	Value             *bigInt `json:"value"`
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

type GetAccountTokenTXResp struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Result  []ResultTokenTX `json:"result"`
}

func (api *api) getAccountTokenTX(address, contractAddress string, desc bool) (accTokTX GetAccountTokenTXResp, err error) {
	const SOURCE = "Etherscan API TokenTX :"
	ident := "a" + address + "-c" + contractAddress
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Etherscan.io/account/tokentx", ident, &accTokTX)
	}
	if !useCache || err != nil {
		params := map[string]string{
			"module": "account",
			"action": "tokentx",
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
		resp, err := api.clientTok.R().
			SetQueryParams(params).
			SetHeader("Accept", "application/json").
			SetResult(&GetAccountTokenTXResp{}).
			Get(api.basePath)
		if err != nil {
			return accTokTX, errors.New(SOURCE + " Error Requesting " + ident)
		}
		accTokTX = *resp.Result().(*GetAccountTokenTXResp)
		if useCache {
			err = db.Write("Etherscan.io/account/tokentx", ident, accTokTX)
			if err != nil {
				return accTokTX, errors.New(SOURCE + " Error Caching " + ident)
			}
		}
		if accTokTX.Message == "OK" {
			time.Sleep(api.timeBetweenReq)
		}
	}
	return accTokTX, nil
}
