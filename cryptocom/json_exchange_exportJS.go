package cryptocom

import (
	"encoding/json"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

type ExchangeJson struct {
	Withs struct {
		FinanceList []struct {
			Symbol        string      `json:"symbol"`
			Reason        string      `json:"reason"`
			Amount        string      `json:"amount"`
			Fee           float64     `json:"fee"`
			UpdateAt      string      `json:"updateAt"`
			TxID          string      `json:"txid"`
			Label         string      `json:"label"`
			AddressTo     string      `json:"addressTo"`
			Network       interface{} `json:"network"`
			CreatedAt     string      `json:"createdAt"`
			TxIDAddr      string      `json:"txidAddr"`
			UpdateAtTime  int64       `json:"updateAtTime"`
			CreatedAtTime int64       `json:"createdAtTime"`
			ID            int         `json:"id"`
			StatusText    string      `json:"status_text"`
			Status        int         `json:"status"`
		} `json:"financeList"`
		Count    int `json:"count"`
		Pagesize int `json:"pageSize"`
	} `json:"withs"`
	Deps struct {
		FinanceList []struct {
			Symbol        string      `json:"symbol"`
			Amount        string      `json:"amount"`
			UpdateAt      string      `json:"updateAt"`
			TxID          string      `json:"txid"`
			NoteStatus    interface{} `json:"note_status"`
			Confirmdesc   string      `json:"confirmDesc"`
			AddressTo     string      `json:"addressTo"`
			Network       interface{} `json:"network"`
			CreatedAt     string      `json:"createdAt"`
			TxIDAddr      string      `json:"txidAddr"`
			UpdateAtTime  int64       `json:"updateAtTime"`
			CreatedAtTime int64       `json:"createdAtTime"`
			StatusText    string      `json:"status_text"`
			Status        int         `json:"status"`
		} `json:"financeList"`
		Count    int `json:"count"`
		Pagesize int `json:"pageSize"`
	} `json:"deps"`
	Cros struct {
		HistoryList []struct {
			StakeAmount    string      `json:"stakeAmount"`
			Apr            float64     `json:"apr,string"`
			CoinSymbol     string      `json:"coinSymbol"`
			CreateTime     string      `json:"createTime"`
			Extra          interface{} `json:"extra"`
			Destination    string      `json:"destination"`
			InterestAmount string      `json:"interestAmount"`
			CreatedAtTime  int64       `json:"createdAtTime"`
			StatusText     string      `json:"status_text"`
			Status         int         `json:"status"`
		} `json:"historyList"`
		Count    int `json:"count"`
		Pagesize int `json:"pageSize"`
	} `json:"cros"`
	Sstake struct {
		Count                   int `json:"count"`
		Pagesize                int `json:"pageSize"`
		SoftStakingInterestList []struct {
			Principal       string      `json:"principal"`
			Reason          interface{} `json:"reason"`
			Amount          string      `json:"amount"`
			Apr             float64     `json:"apr,string"`
			CoinSymbol      string      `json:"coinSymbol"`
			CalculateDate   int64       `json:"calculateDate"`
			Ctime           int64       `json:"ctime"`
			ID              int         `json:"id"`
			StakedCROAmount string      `json:"stakedCroAmount"`
			Mtime           int64       `json:"mtime"`
			UserID          int         `json:"userId"`
			Status          int         `json:"status"`
		} `json:"softStakingInterestList"`
	} `json:"sstake"`
	Rebs struct {
		HistoryList []interface{} `json:"historyList"`
		Count       int           `json:"count"`
		Pagesize    int           `json:"pageSize"`
	} `json:"rebs"`
	Syn struct {
		Activities []interface{} `json:"activities"`
	} `json:"syn"`
	Sup struct {
		HistoryList []struct {
			CreatedAt    int64  `json:"createdAt,string"`
			CoinSymbol   string `json:"coinSymbol"`
			Extra        string `json:"extra"`
			RewardAmount string `json:"rewardAmount"`
		} `json:"historyList"`
		Count    int `json:"count"`
		Pagesize int `json:"pageSize"`
	} `json:"sup"`
	Tcom struct {
		Total int           `json:"total"`
		Data  []interface{} `json:"data"`
	} `json:"tcom"`
	Bon struct {
		Total int           `json:"total"`
		Data  []interface{} `json:"data"`
	} `json:"bon"`
	Rew struct {
		SignupBonusCreatedAt         string `json:"signUpBonusCreatedAt"`
		TotalEarnFromReferral        string `json:"totalEarnFromReferral"`
		TotalNumberOfUsersBeReferred string `json:"totalNumberOfUsersBeReferred"`
		TotalTradeCommission         string `json:"totalTradeCommission"`
		SignupBonus                  string `json:"signUpBonus"`
		TotalReferralBonus           string `json:"totalReferralBonus"`
		TotalNumberOfUsersSignedUp   string `json:"totalNumberOfUsersSignedUp"`
	} `json:"rew"`
}

func (cdc *CryptoCom) ParseJSONExchangeExportJS(reader io.Reader) (err error) {
	const SOURCE = "Crypto.com Exchange JSON ExportJS :"
	var exch ExchangeJson
	jsonDecoder := json.NewDecoder(reader)
	err = jsonDecoder.Decode(&exch)
	if err == nil {
		alreadyAsked := []string{}
		// Fill TXsByCategory
		for _, w := range exch.Withs.FinanceList {
			if w.StatusText == "Completed" {
				t := wallet.TX{Timestamp: time.Unix(w.UpdateAtTime/1000, 0), ID: w.TxID, Note: SOURCE + " Withdrawal " + w.AddressTo}
				t.Items = make(map[string]wallet.Currencies)
				amount, err := decimal.NewFromString(w.Amount)
				if err != nil {
					log.Println(SOURCE, "Error Parsing Amount", w.Amount)
				} else {
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: w.Symbol, Amount: amount})
				}
				if w.Fee != 0 {
					t.Items["Fee"] = append(t.Items["Fee"], wallet.Currency{Code: w.Symbol, Amount: decimal.NewFromFloat(w.Fee)})
				}
				cdc.TXsByCategory["Withdrawals"] = append(cdc.TXsByCategory["Withdrawals"], t)
			}
		}
		for _, d := range exch.Deps.FinanceList {
			if d.StatusText == "Payment received" {
				t := wallet.TX{Timestamp: time.Unix(d.UpdateAtTime/1000, 0), ID: d.TxID, Note: SOURCE + " Deposit " + d.AddressTo}
				t.Items = make(map[string]wallet.Currencies)
				amount, err := decimal.NewFromString(d.Amount)
				if err != nil {
					log.Println(SOURCE, "Error Parsing Amount", d.Amount)
				} else {
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: d.Symbol, Amount: amount})
				}
				cdc.TXsByCategory["Deposits"] = append(cdc.TXsByCategory["Deposits"], t)
			}
		}
		for _, cs := range exch.Cros.HistoryList {
			if cs.StatusText == "Completed" {
				t := wallet.TX{Timestamp: time.Unix(cs.CreatedAtTime/1000, 0), Note: SOURCE + " CRO Stake Interest " + cs.StakeAmount + " " + cs.CoinSymbol + " at " + strconv.FormatFloat(cs.Apr*100, 'f', 1, 64) + "%"}
				t.Items = make(map[string]wallet.Currencies)
				amount, err := decimal.NewFromString(cs.InterestAmount)
				if err != nil {
					log.Println(SOURCE, "Error Parsing InterestAmount", cs.InterestAmount)
				} else {
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: cs.CoinSymbol, Amount: amount})
				}
				cdc.TXsByCategory["Interests"] = append(cdc.TXsByCategory["Interests"], t)
			}
		}
		for _, ss := range exch.Sstake.SoftStakingInterestList {
			if ss.Status == 2 {
				t := wallet.TX{Timestamp: time.Unix(ss.CalculateDate/1000, 0), ID: strconv.Itoa(ss.ID), Note: SOURCE + " Soft Stake Interest " + ss.StakedCROAmount + " " + ss.CoinSymbol + " at " + strconv.FormatFloat(ss.Apr*100, 'f', 1, 64) + "%"}
				t.Items = make(map[string]wallet.Currencies)
				amount, err := decimal.NewFromString(ss.Amount)
				if err != nil {
					log.Println(SOURCE, "Error Parsing Amount", ss.Amount)
				} else {
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: ss.CoinSymbol, Amount: amount})
				}
				cdc.TXsByCategory["Interests"] = append(cdc.TXsByCategory["Interests"], t)
			}
		}
		for _, r := range exch.Rebs.HistoryList {
			// t := wallet.TX{Timestamp: time.Unix(r.CalculateDate/1000, 0), ID: strconv.Itoa(r.ID), Note: SOURCE + " Rebate"}
			// t.Items = make(map[string]wallet.Currencies)
			// amount, err := decimal.NewFromString(r.Amount)
			// if err != nil {
			// 	log.Println(SOURCE, "Error Parsing Amount", r.Amount)
			// } else {
			// 	t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: r.CoinSymbol, Amount: amount})
			// }
			// cdc.TXsByCategory["CommercialRebates"] = append(cdc.TXsByCategory["CommercialRebates"], t)
			alreadyAsked = wallet.AskForHelp(SOURCE+" Rebate", r, alreadyAsked)
		}
		for _, s := range exch.Syn.Activities {
			// t := wallet.TX{Timestamp: time.Unix(s.CalculateDate/1000, 0), ID: strconv.Itoa(s.ID), Note: SOURCE + " Syndicate"}
			// t.Items = make(map[string]wallet.Currencies)
			// amount, err := decimal.NewFromString(s.Amount)
			// if err != nil {
			// 	log.Println(SOURCE, "Error Parsing Amount", s.Amount)
			// } else {
			// 	t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: s.CoinSymbol, Amount: amount})
			// }
			// cdc.TXsByCategory["CommercialRebates"] = append(cdc.TXsByCategory["CommercialRebates"], t)
			alreadyAsked = wallet.AskForHelp(SOURCE+" Syndicate", s, alreadyAsked)
		}
		for _, s := range exch.Sup.HistoryList {
			t := wallet.TX{Timestamp: time.Unix(s.CreatedAt/1000, 0), Note: SOURCE + " Supercharger Reward"}
			t.Items = make(map[string]wallet.Currencies)
			amount, err := decimal.NewFromString(s.RewardAmount)
			if err != nil {
				log.Println(SOURCE, "Error Parsing RewardAmount", s.RewardAmount)
			} else {
				t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: s.CoinSymbol, Amount: amount})
			}
			cdc.TXsByCategory["Minings"] = append(cdc.TXsByCategory["Minings"], t)
		}
		for _, t := range exch.Tcom.Data {
			// t := wallet.TX{Timestamp: time.Unix(t.CalculateDate/1000, 0), ID: strconv.Itoa(t.ID), Note: SOURCE + " Trade Commission"}
			// t.Items = make(map[string]wallet.Currencies)
			// amount, err := decimal.NewFromString(t.Amount)
			// if err != nil {
			// 	log.Println(SOURCE, "Error Parsing Amount", t.Amount)
			// } else {
			// 	t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: t.CoinSymbol, Amount: amount})
			// }
			// cdc.TXsByCategory["Referrals"] = append(cdc.TXsByCategory["Referrals"], t)
			alreadyAsked = wallet.AskForHelp(SOURCE+" Trade Commission", t, alreadyAsked)
		}
		for _, b := range exch.Bon.Data {
			// b := wallet.TX{Timestamp: time.Unix(b.CalculateDate/1000, 0), ID: strconv.Itoa(b.ID), Note: SOURCE + " Referral Bonus"}
			// b.Items = make(map[string]wallet.Currencies)
			// amount, err := decimal.NewFromString(b.Amount)
			// if err != nil {
			// 	log.Println(SOURCE, "Error Parsing Amount", b.Amount)
			// } else {
			// 	b.Items["To"] = append(b.Items["To"], wallet.Currency{Code: b.CoinSymbol, Amount: amount})
			// }
			// cdc.TXsByCategory["Referrals"] = append(cdc.TXsByCategory["Referrals"], b)
			alreadyAsked = wallet.AskForHelp(SOURCE+" Referral Bonus", b, alreadyAsked)
		}
		if exch.Rew.SignupBonus != "0" {
			alreadyAsked = wallet.AskForHelp(SOURCE+" Referral Reward", exch.Rew, alreadyAsked)
			/*
				Rew struct {
					SignupBonusCreatedAt         string `json:"signUpBonusCreatedAt"`
					TotalEarnFromReferral        string `json:"totalEarnFromReferral"`
					TotalNumberOfUsersBeReferred string `json:"totalNumberOfUsersBeReferred"`
					TotalTradeCommission         string `json:"totalTradeCommission"`
					SignupBonus                  string `json:"signUpBonus"`
					TotalReferralBonus           string `json:"totalReferralBonus"`
					TotalNumberOfUsersSignedUp   string `json:"totalNumberOfUsersSignedUp"`
				} `json:"rew"`*/
		}
	}
	return
}
