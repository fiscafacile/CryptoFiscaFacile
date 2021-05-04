package etherscan

import (
	"errors"
	"time"

	"github.com/nanobox-io/golang-scribble"
	"github.com/shopspring/decimal"
)

type internalTX struct {
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

func (api *api) getInternalTXs(addresses []string) {
	for _, eth := range addresses {
		accTXListInternal, err := api.getAccountTXListInternal(eth, false)
		if err != nil {
			api.doneInt <- err
			return
		}
		for _, intTX := range accTXListInternal.Result {
			tx := internalTX{
				used:            false,
				BlockNumber:     intTX.BlockNumber,
				TimeStamp:       time.Unix(intTX.TimeStamp, 0),
				Hash:            intTX.Hash,
				From:            intTX.From,
				To:              intTX.To,
				Value:           decimal.NewFromBigInt(intTX.Value.Int(), -18),
				ContractAddress: intTX.ContractAddress,
				Input:           intTX.Input,
				Type:            intTX.Type,
				Gas:             intTX.Gas,
				GasUsed:         decimal.NewFromInt(intTX.GasUsed),
				TraceID:         intTX.TraceID,
				IsError:         intTX.IsError,
				ErrCode:         intTX.ErrCode,
			}
			api.internalTXs = append(api.internalTXs, tx)
		}
	}
	api.doneInt <- nil
}

type ResultTXListInternal struct {
	BlockNumber     int     `json:"blockNumber,string"`
	TimeStamp       int64   `json:"timeStamp,string"`
	Hash            string  `json:"hash"`
	From            string  `json:"from"`
	To              string  `json:"to"`
	Value           *bigInt `json:"value"`
	ContractAddress string  `json:"contractAddress"`
	Input           string  `json:"input"`
	Type            string  `json:"type"`
	Gas             int     `json:"gas,string"`
	GasUsed         int64   `json:"gasUsed,string"`
	TraceID         string  `json:"traceId"`
	IsError         int     `json:"isError,string"`
	ErrCode         string  `json:"errCode"`
}

type GetAccountTXListInternalResp struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Result  []ResultTXListInternal `json:"result"`
}

func (api *api) getAccountTXListInternal(address string, desc bool) (accTXListInternal GetAccountTXListInternalResp, err error) {
	useCache := true
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		useCache = false
	}
	if useCache {
		err = db.Read("Etherscan.io/account/txlistinternal", address, &accTXListInternal)
	}
	if !useCache || err != nil {
		params := map[string]string{
			"module":  "account",
			"action":  "txlistinternal",
			"address": address,
			"apikey":  api.apiKey,
		}
		if desc {
			params["sort"] = "desc"
		} else {
			params["sort"] = "asc"
		}
		resp, err := api.clientInt.R().
			SetQueryParams(params).
			SetHeader("Accept", "application/json").
			SetResult(&GetAccountTXListInternalResp{}).
			Get(api.basePath)
		if err != nil {
			return accTXListInternal, errors.New("Etherscan API TX List Internal : Error Requesting" + address)
		}
		accTXListInternal = *resp.Result().(*GetAccountTXListInternalResp)
		if useCache {
			err = db.Write("Etherscan.io/account/txlistinternal", address, accTXListInternal)
			if err != nil {
				return accTXListInternal, errors.New("Etherscan API TX List Internal : Error Caching" + address)
			}
		}
		if accTXListInternal.Message == "OK" {
			time.Sleep(api.timeBetweenReq)
		}
	}
	return accTXListInternal, nil
}
