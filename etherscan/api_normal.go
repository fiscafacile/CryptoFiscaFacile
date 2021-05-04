package etherscan

import (
	"errors"
	"time"

	"github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type normalTX struct {
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

func (api *api) getNormalTXs(addresses []string) {
	for _, eth := range addresses {
		accTXList, err := api.getAccountTXList(eth, false)
		if err != nil {
			api.doneNor <- err
			return
		}
		for _, norTX := range accTXList.Result {
			tx := normalTX{
				used:              false,
				BlockNumber:       norTX.BlockNumber,
				TimeStamp:         time.Unix(norTX.TimeStamp, 0),
				Hash:              norTX.Hash,
				Nonce:             norTX.Nonce,
				BlockHash:         norTX.BlockHash,
				TransactionIndex:  norTX.TransactionIndex,
				From:              norTX.From,
				To:                norTX.To,
				Value:             decimal.NewFromBigInt(norTX.Value.Int(), -18),
				Gas:               norTX.Gas,
				GasPrice:          decimal.NewFromBigInt(norTX.GasPrice.Int(), -18),
				IsError:           norTX.IsError,
				TxReceiptStatus:   norTX.TxReceiptStatus,
				Input:             norTX.Input,
				ContractAddress:   norTX.ContractAddress,
				CumulativeGasUsed: norTX.CumulativeGasUsed,
				GasUsed:           decimal.NewFromInt(norTX.GasUsed),
				Confirmations:     norTX.Confirmations,
			}
			api.normalTXs = append(api.normalTXs, tx)
		}
	}
	api.doneNor <- nil
}

type ResultTXList struct {
	BlockNumber       int     `json:"blockNumber,string"`
	TimeStamp         int64   `json:"timeStamp,string"`
	Hash              string  `json:"hash"`
	Nonce             int     `json:"nonce,string"`
	BlockHash         string  `json:"blockHash"`
	TransactionIndex  int     `json:"transactionIndex,string"`
	From              string  `json:"from"`
	To                string  `json:"to"`
	Value             *bigInt `json:"value"`
	Gas               int     `json:"gas,string"`
	GasPrice          *bigInt `json:"gasPrice"`
	IsError           int     `json:"isError,string"`
	TxReceiptStatus   string  `json:"txreceipt_status"`
	Input             string  `json:"input"`
	ContractAddress   string  `json:"contractAddress"`
	CumulativeGasUsed int     `json:"cumulativeGasUsed,string"`
	GasUsed           int64   `json:"gasUsed,string"`
	Confirmations     int     `json:"confirmations,string"`
}

type GetAccountTXListResp struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Result  []ResultTXList `json:"result"`
}

func (api *api) getAccountTXList(address string, desc bool) (accTXList GetAccountTXListResp, err error) {
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Etherscan.io/account/txlist", address, &accTXList)
	}
	if !useCache || err != nil {
		params := map[string]string{
			"module":  "account",
			"action":  "txlist",
			"address": address,
			"apikey":  api.apiKey,
		}
		if desc {
			params["sort"] = "desc"
		} else {
			params["sort"] = "asc"
		}
		resp, err := api.clientNor.R().
			SetQueryParams(params).
			SetHeader("Accept", "application/json").
			SetResult(&GetAccountTXListResp{}).
			Get(api.basePath)
		if err != nil {
			return accTXList, errors.New("Etherscan API TX List : Error Requesting" + address)
		}
		accTXList = *resp.Result().(*GetAccountTXListResp)
		if useCache {
			err = db.Write("Etherscan.io/account/txlist", address, accTXList)
			if err != nil {
				return accTXList, errors.New("Etherscan API TX List : Error Caching" + address)
			}
		}
		if accTXList.Message == "OK" {
			time.Sleep(api.timeBetweenReq)
		}
	}
	return accTXList, nil
}
