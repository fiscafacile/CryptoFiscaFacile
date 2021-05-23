package cryptocom

import (
	"encoding/json"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/source"
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
		HistoryList []struct {
			CreateTime       string  `json:"createTime"`
			RebateAmount     string  `json:"rebateAmount"`
			RebatePercentage float64 `json:"rebatePercentage,string"`
			FeePaid          string  `json:"feePaid"`
			CreatedAtTime    int64   `json:"createdAtTime"`
			CoinSymbol       string  `json:"coinSymbol"`
			Destination      string  `json:"destination"`
			StatusText       string  `json:"status_text"`
			Status           int     `json:"status"`
			Extra            string  `json:"extra"`
		} `json:"historyList"`
		Count    int `json:"count"`
		Pagesize int `json:"pageSize"`
	} `json:"rebs"`
	Syn struct {
		Activities []struct {
			DeliveredSize          string `json:"deliveredSize"`
			ActivityCROn           string `json:"activityCron"`
			UserStatus             string `json:"userStatus"`
			ID                     string `json:"id"`
			AllocationTime         string `json:"allocationTime"`
			MinCommittedCRO        string `json:"minCommittedCRO"`
			RefundedCRO            string `json:"refundedCRO"`
			PoolSize               string `json:"poolSize"`
			DiscountedPrice        string `json:"discountedPrice"`
			MinPurchased           string `json:"minPurchased"`
			DiscountRate           string `json:"discountRate"`
			AnnouncementTime       string `json:"announcementTime"`
			ActivityStatus         string `json:"activityStatus"`
			PriceDeterminationTime string `json:"priceDeterminationTime"`
			AllocatedPriceCRO      string `json:"allocatedPriceCRO"`
			UserID                 string `json:"userId"`
			ActivityModifyTime     string `json:"activityModifyTime"`
			EndTime                string `json:"endTime"`
			DeliveryTime           int64  `json:"deliveryTime,string"`
			SyndicateCoin          string `json:"syndicateCoin"`
			TotalCommittedCRO      string `json:"totalCommittedCro"`
			ActivityCreateTime     string `json:"activityCreateTime"`
			UserEmailStatus        string `json:"userEmailStatus"`
			UserCreateTime         int64  `json:"userCreateTime,string"`
			PoolSizeCapUSD         string `json:"poolSizeCapUSD"`
			StartTime              string `json:"startTime"`
			AllocatedVolume        string `json:"allocatedVolume"`
			CommittedCRO           string `json:"committedCRO"`
			AllocatedSize          string `json:"allocatedSize"`
			AllocatedPriceUSD      string `json:"allocatedPriceUSD"`
			UserModifyTime         string `json:"userModifyTime"`
		} `json:"activities"`
	} `json:"syn"`
	Sup struct {
		HistoryList []struct {
			// CreatedAt string `json:"createdAt"`
			CreatedAt    int64  `json:"createdAt,string"`
			CoinSymbol   string `json:"coinSymbol"`
			Extra        string `json:"extra"`
			RewardAmount string `json:"rewardAmount"`
		} `json:"historyList"`
		Count    int `json:"count"`
		Pagesize int `json:"pageSize"`
	} `json:"sup"`
	Tcom struct {
		Total int `json:"total"`
		Data  []struct {
			Commission             string `json:"commission"`
			ID                     string `json:"id"`
			MTime                  int64  `json:"mtime,string"`
			Status                 int    `json:"status,string"`
			NetTradingFee          string `json:"netTradingFee"`
			ReferralRelationshipID string `json:"referralRelationshipId"`
			CTime                  int64  `json:"ctime,string"`
			TradingFeeRebate       string `json:"tradingFeeRebate"`
		} `json:"data"`
	} `json:"tcom"`
	Bon struct {
		Total int `json:"total"`
		Data  []struct {
			ReferralRelationshipID string `json:"referralRelationshipId"`
			ReferralBonusInCRO     string `json:"referralBonusInCro"`
			CTime                  int64  `json:"ctime,string"`
			ID                     string `json:"id"`
			MTime                  int64  `json:"mtime,string"`
			ReferralBonusTierID    string `json:"referralBonusTierId"`
			Status                 int    `json:"status,string"`
			ReferralBonusInUSD     string `json:"referralBonusInUsd"`
		} `json:"data"`
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

type jsonEx struct {
	txsByCategory wallet.TXsByCategory
}

func (cdc *CryptoCom) ParseJSONExchangeExportJS(reader io.Reader, account string) (err error) {
	cdc.jsonEx.txsByCategory = make(map[string]wallet.TXs)
	firstTimeUsed := time.Now()
	lastTimeUsed := time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC)
	const SOURCE = "Crypto.com Exchange JSON ExportJS :"
	var exch ExchangeJson
	jsonDecoder := json.NewDecoder(reader)
	err = jsonDecoder.Decode(&exch)
	if err == nil {
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
				cdc.jsonEx.txsByCategory["Withdrawals"] = append(cdc.jsonEx.txsByCategory["Withdrawals"], t)
			}
			if time.Unix(w.UpdateAtTime/1000, 0).Before(firstTimeUsed) {
				firstTimeUsed = time.Unix(w.UpdateAtTime/1000, 0)
			}
			if time.Unix(w.UpdateAtTime/1000, 0).After(lastTimeUsed) {
				lastTimeUsed = time.Unix(w.UpdateAtTime/1000, 0)
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
					cdc.jsonEx.txsByCategory["Deposits"] = append(cdc.jsonEx.txsByCategory["Deposits"], t)
				}
			}
			if time.Unix(d.UpdateAtTime/1000, 0).Before(firstTimeUsed) {
				firstTimeUsed = time.Unix(d.UpdateAtTime/1000, 0)
			}
			if time.Unix(d.UpdateAtTime/1000, 0).After(lastTimeUsed) {
				lastTimeUsed = time.Unix(d.UpdateAtTime/1000, 0)
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
					cdc.jsonEx.txsByCategory["Interests"] = append(cdc.jsonEx.txsByCategory["Interests"], t)
				}
			}
			if time.Unix(cs.CreatedAtTime/1000, 0).Before(firstTimeUsed) {
				firstTimeUsed = time.Unix(cs.CreatedAtTime/1000, 0)
			}
			if time.Unix(cs.CreatedAtTime/1000, 0).After(lastTimeUsed) {
				lastTimeUsed = time.Unix(cs.CreatedAtTime/1000, 0)
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
					cdc.jsonEx.txsByCategory["Interests"] = append(cdc.jsonEx.txsByCategory["Interests"], t)
				}
			}
			if time.Unix(ss.CalculateDate/1000, 0).Before(firstTimeUsed) {
				firstTimeUsed = time.Unix(ss.CalculateDate/1000, 0)
			}
			if time.Unix(ss.CalculateDate/1000, 0).After(lastTimeUsed) {
				lastTimeUsed = time.Unix(ss.CalculateDate/1000, 0)
			}
		}
		for _, r := range exch.Rebs.HistoryList {
			if r.StatusText == "Completed" {
				t := wallet.TX{Timestamp: time.Unix(r.CreatedAtTime/1000, 0), Note: SOURCE + " Rebate on Fee paid " + r.FeePaid + " " + r.CoinSymbol + " at " + strconv.FormatFloat(r.RebatePercentage*100, 'f', 1, 64) + "%"}
				t.Items = make(map[string]wallet.Currencies)
				amount, err := decimal.NewFromString(r.RebateAmount)
				if err != nil {
					log.Println(SOURCE, "Error Parsing RebateAmount", r.RebateAmount)
				} else {
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: r.CoinSymbol, Amount: amount})
					cdc.jsonEx.txsByCategory["CommercialRebates"] = append(cdc.jsonEx.txsByCategory["CommercialRebates"], t)
				}
			}
			if time.Unix(r.CreatedAtTime/1000, 0).Before(firstTimeUsed) {
				firstTimeUsed = time.Unix(r.CreatedAtTime/1000, 0)
			}
			if time.Unix(r.CreatedAtTime/1000, 0).After(lastTimeUsed) {
				lastTimeUsed = time.Unix(r.CreatedAtTime/1000, 0)
			}
		}
		for _, s := range exch.Syn.Activities {
			t := wallet.TX{Timestamp: time.Unix(s.DeliveryTime/1000, 0), ID: s.ID, Note: SOURCE + " Syndicate"}
			t.Items = make(map[string]wallet.Currencies)
			allocatedVolume, err1 := decimal.NewFromString(s.AllocatedVolume)
			if err1 != nil {
				log.Println(SOURCE, "Error Parsing AllocatedVolume", s.AllocatedVolume)
			} else {
				t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: s.SyndicateCoin, Amount: allocatedVolume})
			}
			committedCRO, err2 := decimal.NewFromString(s.CommittedCRO)
			var err3 error
			if err2 != nil {
				log.Println(SOURCE, "Error Parsing CommittedCRO", s.CommittedCRO)
			} else {
				refundedCRO, err3 := decimal.NewFromString(s.RefundedCRO)
				if err3 != nil {
					log.Println(SOURCE, "Error Parsing RefundedCRO", s.RefundedCRO)
				} else {
					t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: "CRO", Amount: committedCRO.Sub(refundedCRO)})
				}
			}
			if err1 == nil && err2 == nil && err3 == nil {
				cdc.jsonEx.txsByCategory["Exchanges"] = append(cdc.jsonEx.txsByCategory["Exchanges"], t)
			}
			if time.Unix(s.UserCreateTime/1000, 0).Before(firstTimeUsed) {
				firstTimeUsed = time.Unix(s.UserCreateTime/1000, 0)
			}
			if time.Unix(s.DeliveryTime/1000, 0).After(lastTimeUsed) {
				lastTimeUsed = time.Unix(s.DeliveryTime/1000, 0)
			}
		}
		for _, s := range exch.Sup.HistoryList {
			t := wallet.TX{Timestamp: time.Unix(s.CreatedAt/1000, 0), Note: SOURCE + " Supercharger Reward"}
			t.Items = make(map[string]wallet.Currencies)
			amount, err := decimal.NewFromString(s.RewardAmount)
			if err != nil {
				log.Println(SOURCE, "Error Parsing RewardAmount", s.RewardAmount)
			} else {
				t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: s.CoinSymbol, Amount: amount})
				cdc.jsonEx.txsByCategory["Minings"] = append(cdc.jsonEx.txsByCategory["Minings"], t)
			}
			if time.Unix(s.CreatedAt/1000, 0).Before(firstTimeUsed) {
				firstTimeUsed = time.Unix(s.CreatedAt/1000, 0)
			}
			if time.Unix(s.CreatedAt/1000, 0).After(lastTimeUsed) {
				lastTimeUsed = time.Unix(s.CreatedAt/1000, 0)
			}
		}
		for _, tc := range exch.Tcom.Data {
			if tc.Status == 1 {
				t := wallet.TX{Timestamp: time.Unix(tc.MTime/1000, 0), ID: tc.ID, Note: SOURCE + " Trade Commission"}
				t.Items = make(map[string]wallet.Currencies)
				amount, err := decimal.NewFromString(tc.Commission)
				if err != nil {
					log.Println(SOURCE, "Error Parsing Commission", tc.Commission)
				} else {
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "CRO", Amount: amount})
					cdc.jsonEx.txsByCategory["Referrals"] = append(cdc.jsonEx.txsByCategory["Referrals"], t)
				}
			}
			if time.Unix(tc.MTime/1000, 0).Before(firstTimeUsed) {
				firstTimeUsed = time.Unix(tc.MTime/1000, 0)
			}
			if time.Unix(tc.MTime/1000, 0).After(lastTimeUsed) {
				lastTimeUsed = time.Unix(tc.MTime/1000, 0)
			}
		}
		for _, b := range exch.Bon.Data {
			if b.Status == 2 {
				t := wallet.TX{Timestamp: time.Unix(b.MTime/1000, 0), ID: b.ID, Note: SOURCE + " Referral Bonus"}
				t.Items = make(map[string]wallet.Currencies)
				amount, err := decimal.NewFromString(b.ReferralBonusInCRO)
				if err != nil {
					log.Println(SOURCE, "Error Parsing ReferralBonusInCRO", b.ReferralBonusInCRO)
				} else {
					t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "CRO", Amount: amount})
					cdc.jsonEx.txsByCategory["Referrals"] = append(cdc.jsonEx.txsByCategory["Referrals"], t)
				}
			}
			if time.Unix(b.MTime/1000, 0).Before(firstTimeUsed) {
				firstTimeUsed = time.Unix(b.MTime/1000, 0)
			}
			if time.Unix(b.MTime/1000, 0).After(lastTimeUsed) {
				lastTimeUsed = time.Unix(b.MTime/1000, 0)
			}
		}
		if exch.Rew.SignupBonus != "0" {
			signupBonusCreatedAt, err := strconv.ParseInt(exch.Rew.SignupBonusCreatedAt, 10, 64)
			if err != nil {
				log.Println(SOURCE, "Error Parsing SignupBonusCreatedAt", exch.Rew.SignupBonusCreatedAt)
			}
			t := wallet.TX{Timestamp: time.Unix(signupBonusCreatedAt/1000, 0), Note: SOURCE + " Signup Bonus"}
			t.Items = make(map[string]wallet.Currencies)
			amount, err := decimal.NewFromString(exch.Rew.SignupBonus)
			if err != nil {
				log.Println(SOURCE, "Error Parsing SignupBonus", exch.Rew.SignupBonus)
			} else {
				t.Items["To"] = append(t.Items["To"], wallet.Currency{Code: "CRO", Amount: amount})
				cdc.jsonEx.txsByCategory["CommercialRebates"] = append(cdc.jsonEx.txsByCategory["CommercialRebates"], t)
			}
		}
	}
	if _, ok := cdc.Sources["CdC Exchange"]; !ok {
		cdc.Sources["CdC Exchange"] = source.Source{
			Crypto:        true,
			AccountNumber: account,
			OpeningDate:   firstTimeUsed,
			ClosingDate:   lastTimeUsed,
			LegalName:     "MCO Malta DAX Limited",
			Address:       "Level 7, Spinola Park, Triq Mikiel Ang Borg,\nSt Julian's SPK 1000,\nMalte",
			URL:           "https://crypto.com/exchange",
		}
	}
	return
}
