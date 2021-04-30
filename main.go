package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/binance"
	"github.com/fiscafacile/CryptoFiscaFacile/bitfinex"
	"github.com/fiscafacile/CryptoFiscaFacile/bittrex"
	"github.com/fiscafacile/CryptoFiscaFacile/blockchain"
	"github.com/fiscafacile/CryptoFiscaFacile/blockstream"
	"github.com/fiscafacile/CryptoFiscaFacile/btc"
	"github.com/fiscafacile/CryptoFiscaFacile/category"
	"github.com/fiscafacile/CryptoFiscaFacile/coinbase"
	"github.com/fiscafacile/CryptoFiscaFacile/cryptocom"
	"github.com/fiscafacile/CryptoFiscaFacile/etherscan"
	"github.com/fiscafacile/CryptoFiscaFacile/ledgerlive"
	"github.com/fiscafacile/CryptoFiscaFacile/localbitcoin"
	"github.com/fiscafacile/CryptoFiscaFacile/metamask"
	"github.com/fiscafacile/CryptoFiscaFacile/mycelium"
	"github.com/fiscafacile/CryptoFiscaFacile/revolut"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

func main() {
	// Parse args
	pDate := flag.String("date", "2021-01-01T00:00:00", "Date Filter")
	pLocation := flag.String("location", "Europe/Paris", "Date Filter Location")
	pNative := flag.String("native", "EUR", "Native Currency for consolidation")
	pTXsDisplayCat := flag.String("txs_display", "", "Display Transactions By Catergory : Exchanges|Deposits|Withdrawals|CashIn|CashOut|etc")
	pExact := flag.Bool("exact", false, "Display exact amount (no rounding)")
	pCurrFilter := flag.String("curr_filter", "", "Currencies to be filtered in Transactions Display (comma separated list)")
	pStats := flag.Bool("stats", false, "Display accounts stats")
	pCheck := flag.Bool("check", false, "Check and Display consistency")
	p2086 := flag.Bool("2086", false, "Display Cerfa 2086")
	pCSVTXsCateg := flag.String("txs_categ", "", "Transactions Categories CSV file")
	pCoinAPIKey := flag.String("coinapi_key", "", "CoinAPI Key (https://www.coinapi.io/pricing?apikey)")
	pCoinLayerKey := flag.String("coinlayer_key", "", "CoinLayer Key (https://coinlayer.com/product)")
	pCSVBtcAddress := flag.String("btc_address", "", "Bitcoin Addresses CSV file")
	pBCD := flag.Bool("bcd", false, "Detect Bitcoin Diamond Fork")
	pBCH := flag.Bool("bch", false, "Detect Bitcoin Cash Fork")
	pBTG := flag.Bool("btg", false, "Detect Bitcoin Gold Fork")
	pLBTC := flag.Bool("lbtc", false, "Detect Lightning Bitcoin Fork")
	pJsonBtgTXs := flag.String("btg_txs", "", "Bitcoin Gold Transactions JSON file")
	pCSVEthAddress := flag.String("eth_address", "", "Ethereum Addresses CSV file")
	pEtherscanAPIKey := flag.String("etherscan_apikey", "", "Etherscan API Key (https://etherscan.io/myapikey)")
	pCSVBinance := flag.String("binance", "", "Binance CSV file")
	pCSVBinanceExtended := flag.Bool("binance_extended", false, "Use Binance CSV file extended format")
	pCSVBitfinex := flag.String("bitfinex", "", "Bitfinex CSV file")
	pAPIBittrexKey := flag.String("bittrex_api_key", "", "Bittrex API key")
	pAPIBittrexSecret := flag.String("bittrex_api_secret", "", "Bittrex API secret")
	pCSVCoinbase := flag.String("coinbase", "", "Coinbase CSV file")
	pCSVCdCAppCrypto := flag.String("cdc_app_crypto", "", "Crypto.com App Crypto Wallet CSV file")
	pCSVCdCExTransfer := flag.String("cdc_ex_transfer", "", "Crypto.com Exchange Deposit/Withdrawal CSV file")
	pCSVCdCExStake := flag.String("cdc_ex_stake", "", "Crypto.com Exchange Stake CSV file")
	pCSVCdCExSupercharger := flag.String("cdc_ex_supercharger", "", "Crypto.com Exchange Supercharger CSV file")
	pCSVLedgerLive := flag.String("ledgerlive", "", "LedgerLive CSV file")
	pCSVLBTrade := flag.String("lb_trade", "", "Local Bitcoin Trade CSV file")
	pCSVLBTransfer := flag.String("lb_transfer", "", "Local Bitcoin Transfer CSV file")
	pCSVMetaMask := flag.String("metamask", "", "MetaMask CSV file")
	pCSVMyCelium := flag.String("mycelium", "", "MyCelium CSV file")
	pCSVRevo := flag.String("revolut", "", "Revolut CSV file")
	flag.Parse()
	if *pCoinAPIKey != "" {
		wallet.CoinAPISetKey(*pCoinAPIKey)
	}
	if *pCoinLayerKey != "" {
		wallet.CoinLayerSetKey(*pCoinLayerKey)
	}
	categ := category.New()
	if *pCSVTXsCateg != "" {
		recordFile, err := os.Open(*pCSVTXsCateg)
		if err != nil {
			log.Fatal("Error opening Transactions CSV Category file:", err)
		}
		categ.ParseCSVCategory(recordFile)
	}
	btc := btc.New()
	blkst := blockstream.New()
	if *pCSVBtcAddress != "" {
		recordFile, err := os.Open(*pCSVBtcAddress)
		if err != nil {
			log.Fatal("Error opening Bitcoin CSV Addresses file:", err)
		}
		btc.ParseCSVAddresses(recordFile)
		go blkst.GetAllTXs(btc, *categ)
	}
	ethsc := etherscan.New()
	if *pCSVEthAddress != "" {
		recordFile, err := os.Open(*pCSVEthAddress)
		if err != nil {
			log.Fatal("Error opening Ethereum CSV Addresses file:", err)
		}
		ethsc.APIConnect(*pEtherscanAPIKey)
		go ethsc.ParseCSV(recordFile, *categ)
	}
	bc := blockchain.New()
	if *pJsonBtgTXs != "" {
		jsonFile, err := os.Open(*pJsonBtgTXs)
		if err != nil {
			log.Fatal("Error opening Bitcoin Gold JSON Transactions file:", err)
		}
		err = bc.ParseTXsJSON(jsonFile, "BTG")
		if err != nil {
			log.Fatal("Error parsing Bitcoin Gold JSON Transactions file:", err)
		}
	}
	b := binance.New()
	if *pCSVBinance != "" {
		recordFile, err := os.Open(*pCSVBinance)
		if err != nil {
			log.Fatal("Error opening Binance CSV file:", err)
		}
		if *pCSVBinanceExtended {
			err = b.ParseCSVExtended(recordFile)
		} else {
			err = b.ParseCSV(recordFile)
		}
		if err != nil {
			log.Fatal("Error parsing Binance CSV file:", err)
		}
	}
	bf := bitfinex.New()
	if *pCSVBitfinex != "" {
		recordFile, err := os.Open(*pCSVBitfinex)
		if err != nil {
			log.Fatal("Error opening Bitfinex CSV file:", err)
		}
		err = bf.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Bitfinex CSV file:", err)
		}
	}
	btrx := bittrex.New()
	if *pAPIBittrexKey != "" && *pAPIBittrexSecret != "" {
		go btrx.GetAllTransferTXs(*pAPIBittrexKey, *pAPIBittrexSecret, *categ)
		go btrx.GetAllTradeTXs(*pAPIBittrexKey, *pAPIBittrexSecret, *categ)
	}
	cb := coinbase.New()
	if *pCSVCoinbase != "" {
		recordFile, err := os.Open(*pCSVCoinbase)
		if err != nil {
			log.Fatal("Error opening Coinbase CSV file:", err)
		}
		err = cb.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Coinbase CSV file:", err)
		}
	}
	cdc := cryptocom.New()
	if *pCSVCdCAppCrypto != "" {
		recordFile, err := os.Open(*pCSVCdCAppCrypto)
		if err != nil {
			log.Fatal("Error opening Crypto.com CSV file:", err)
		}
		err = cdc.ParseCSVCrypto(recordFile)
		if err != nil {
			log.Fatal("Error parsing Crypto.com CSV file:", err)
		}
	}
	if *pCSVCdCExTransfer != "" {
		recordFile, err := os.Open(*pCSVCdCExTransfer)
		if err != nil {
			log.Fatal("Error opening Crypto.com Exchange Deposit/Withdrawal CSV file:", err)
		}
		err = cdc.ParseCSVExTransfer(recordFile)
		if err != nil {
			log.Fatal("Error parsing Crypto.com Exchange Deposit/Withdrawal CSV file:", err)
		}
	}
	if *pCSVCdCExStake != "" {
		recordFile, err := os.Open(*pCSVCdCExStake)
		if err != nil {
			log.Fatal("Error opening Crypto.com Exchange Stake CSV file:", err)
		}
		err = cdc.ParseCSVExStake(recordFile)
		if err != nil {
			log.Fatal("Error parsing Crypto.com Exchange Stake CSV file:", err)
		}
	}
	if *pCSVCdCExSupercharger != "" {
		recordFile, err := os.Open(*pCSVCdCExSupercharger)
		if err != nil {
			log.Fatal("Error opening Crypto.com Exchange Supercharger CSV file:", err)
		}
		err = cdc.ParseCSVExSupercharger(recordFile)
		if err != nil {
			log.Fatal("Error parsing Crypto.com Exchange Supercharger CSV file:", err)
		}
	}
	ll := ledgerlive.New()
	if *pCSVLedgerLive != "" {
		recordFile, err := os.Open(*pCSVLedgerLive)
		if err != nil {
			log.Fatal("Error opening LedgerLive CSV file:", err)
		}
		err = ll.ParseCSV(recordFile, *categ)
		if err != nil {
			log.Fatal("Error parsing LedgerLive CSV file:", err)
		}
	}
	lb := localbitcoin.New()
	if *pCSVLBTrade != "" {
		recordFile, err := os.Open(*pCSVLBTrade)
		if err != nil {
			log.Fatal("Error opening Local Bitcoin Trade CSV file:", err)
		}
		err = lb.ParseTradeCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Local Bitcoin Trade CSV file:", err)
		}
	}
	if *pCSVLBTransfer != "" {
		recordFile, err := os.Open(*pCSVLBTransfer)
		if err != nil {
			log.Fatal("Error opening Local Bitcoin Transfer CSV file:", err)
		}
		err = lb.ParseTransferCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Local Bitcoin Transfer CSV file:", err)
		}
	}
	mm := metamask.New()
	if *pCSVMetaMask != "" {
		recordFile, err := os.Open(*pCSVMetaMask)
		if err != nil {
			log.Fatal("Error opening MetaMask CSV file:", err)
		}
		err = mm.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing MetaMask CSV file:", err)
		}
	}
	mc := mycelium.New()
	if *pCSVMyCelium != "" {
		recordFile, err := os.Open(*pCSVMyCelium)
		if err != nil {
			log.Fatal("Error opening MyCelium CSV file:", err)
		}
		err = mc.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing MyCelium CSV file:", err)
		}
	}
	revo := revolut.New()
	if *pCSVRevo != "" {
		recordFile, err := os.Open(*pCSVRevo)
		if err != nil {
			log.Fatal("Error opening Revolut CSV file:", err)
		}
		err = revo.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Revolut CSV file:", err)
		}
	}
	if *pCSVEthAddress != "" {
		err := ethsc.WaitFinish()
		if err != nil {
			log.Fatal("Error parsing Ethereum CSV file:", err)
		}
	}
	if *pAPIBittrexKey != "" && *pAPIBittrexSecret != "" {
		errTransfer := btrx.WaitTransfersFinish()
		if errTransfer != nil {
			log.Fatalln("Error parsing Bittrex API transfers:", errTransfer)
		}
		errTrades := btrx.WaitTradesFinish()
		if errTrades != nil {
			log.Fatalln("Error parsing Bittrex API trades:", errTrades)
		}
	}
	if *pCSVBtcAddress != "" {
		err := blkst.WaitFinish()
		if err != nil {
			log.Fatal("Error parsing Bitcoin CSV file:", err)
		}
		if *pBCD {
			blkst.DetectBCD(btc)
		}
		if *pBCH {
			blkst.DetectBCH(btc)
		}
		if *pBTG {
			blkst.DetectBTG(btc)
		}
		if *pLBTC {
			blkst.DetectLBTC(btc)
		}
	}
	// create Global Wallet up to Date
	global := make(wallet.TXsByCategory)
	global.Add(b.TXsByCategory)
	global.Add(bf.TXsByCategory)
	global.Add(btrx.TXsByCategory)
	global.Add(cb.TXsByCategory)
	global.Add(cdc.TXsByCategory)
	global.Add(ll.TXsByCategory)
	global.Add(lb.TXsByCategory)
	global.Add(mm.TXsByCategory)
	global.Add(mc.TXsByCategory)
	global.Add(revo.TXsByCategory)
	global.Add(ethsc.TXsByCategory)
	global.Add(btc.TXsByCategory)
	global.Add(bc.TXsByCategory)
	global.FindTransfers()
	totalCommercialRebates, totalInterests, totalReferrals := global.FindCashInOut(*pNative)
	global.SortTXsByDate(true)
	loc, err := time.LoadLocation(*pLocation)
	if err != nil {
		log.Fatal("Error parsing Location:", err)
	}
	if *pStats {
		global.PrintStats(*pNative, totalCommercialRebates, totalInterests, totalReferrals)
	}
	if *pCheck {
		global.CheckConsistency(loc)
	}
	// Debug
	if *pTXsDisplayCat != "" {
		if *pTXsDisplayCat == "Alls" {
			global.Println(*pCurrFilter)
		} else {
			global[*pTXsDisplayCat].Println("Category "+*pTXsDisplayCat, *pCurrFilter)
		}
	}
	// Construct global wallet up to date
	filterDate, err := time.ParseInLocation("2006-01-02T15:04:05", *pDate, loc)
	if err != nil {
		log.Fatal("Error parsing Date:", err)
	}
	globalWallet := global.GetWallets(filterDate, false, !*pExact)
	globalWallet.Println("Global Crypto", "")
	globalWalletTotalValue, err := globalWallet.CalculateTotalValue(*pNative)
	if err != nil {
		log.Fatal("Error Calculating Global Wallet:", err)
	} else {
		globalWalletTotalValue.Amount = globalWalletTotalValue.Amount.RoundBank(0)
		fmt.Print("Total Value : ")
		globalWalletTotalValue.Println("")
	}
	if *p2086 {
		var cessions Cessions
		err = cessions.CalculatePVMV(global, *pNative, loc)
		if err != nil {
			log.Fatal(err)
		}
		cessions.Println()
	}
	os.Exit(0)
}
