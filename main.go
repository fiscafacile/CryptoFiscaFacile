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
	"github.com/fiscafacile/CryptoFiscaFacile/hitbtc"
	"github.com/fiscafacile/CryptoFiscaFacile/kraken"
	"github.com/fiscafacile/CryptoFiscaFacile/ledgerlive"
	"github.com/fiscafacile/CryptoFiscaFacile/localbitcoin"
	"github.com/fiscafacile/CryptoFiscaFacile/mycelium"
	"github.com/fiscafacile/CryptoFiscaFacile/revolut"
	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

func main() {
	// General Options
	pDate := flag.String("date", "2021-01-01T00:00:00", "Date Filter")
	pLocation := flag.String("location", "Europe/Paris", "Date Filter Location")
	pNative := flag.String("native", "EUR", "Native Currency for consolidation")
	pStats := flag.Bool("stats", false, "Display accounts stats")
	// Debug
	pDebug := flag.Bool("debug", false, "Debug Mode (only for devs)")
	pCheck := flag.Bool("check", false, "Check and Display consistency")
	pCurrFilter := flag.String("curr_filter", "", "Currencies to be filtered in Transactions Display (comma separated list)")
	pExact := flag.Bool("exact", false, "Display exact amount (no rounding)")
	pTXsDisplayCat := flag.String("txs_display", "", "Display Transactions By Catergory : Exchanges|Deposits|Withdrawals|CashIn|CashOut|etc")
	// Sources
	pTXsCategCSV := flag.String("txs_categ", "", "Transactions Categories CSV file")
	pCoinAPIKey := flag.String("coinapi_key", "", "CoinAPI Key (https://www.coinapi.io/pricing?apikey)")
	pCoinLayerKey := flag.String("coinlayer_key", "", "CoinLayer Key (https://coinlayer.com/product)")
	pBTCAddressesCSV := flag.String("btc_address", "", "Bitcoin Addresses CSV file")
	pBCD := flag.Bool("bcd", false, "Detect Bitcoin Diamond Fork")
	pBCH := flag.Bool("bch", false, "Detect Bitcoin Cash Fork")
	pBTG := flag.Bool("btg", false, "Detect Bitcoin Gold Fork")
	pLBTC := flag.Bool("lbtc", false, "Detect Lightning Bitcoin Fork")
	pBTGTXsJSON := flag.String("btg_txs", "", "Bitcoin Gold Transactions JSON file")
	pETHAddressesCSV := flag.String("eth_address", "", "Ethereum Addresses CSV file")
	pEtherscanAPIKey := flag.String("etherscan_apikey", "", "Etherscan API Key (https://etherscan.io/myapikey)")
	pBinanceAPIKey := flag.String("binance_api_key", "", "Binance API Key")
	pBinanceSecretKey := flag.String("binance_secret_key", "", "Binance Secret Key")
	pBinanceCSV := flag.String("binance", "", "Binance CSV file")
	pBinanceCSVExtended := flag.Bool("binance_extended", false, "Use Binance CSV file extended format")
	pBitfinexCSV := flag.String("bitfinex", "", "Bitfinex CSV file")
	pBittrexAPIKey := flag.String("bittrex_api_key", "", "Bittrex API key")
	pBittrexAPISecret := flag.String("bittrex_api_secret", "", "Bittrex API secret")
	pBittrexCSV := flag.String("bittrex", "", "Bittrex CSV file")
	pHitBtcCSVTrades := flag.String("hitbtc_trades", "", "HitBTC Trades CSV file")
	pHitBtcCSVTransactions := flag.String("hitbtc_transactions", "", "HitBTC Transactions CSV file")
	pHitBtcAPIKey := flag.String("hitbtc_api_key", "", "HitBTC API Key")
	pHitBtcSecretKey := flag.String("hitbtc_secret_key", "", "HitBTC Secret Key")
	pCoinbaseCSV := flag.String("coinbase", "", "Coinbase CSV file")
	pCdCAppCSVCrypto := flag.String("cdc_app_crypto", "", "Crypto.com App Crypto Wallet CSV file")
	pCdCExAPIKey := flag.String("cdc_ex_api_key", "", "Crypto.com Exchange API Key")
	pCdCExSecretKey := flag.String("cdc_ex_secret_key", "", "Crypto.com Exchange Secret Key")
	pCdCExJSONExportJS := flag.String("cdc_ex_exportjs", "", "Crypto.com Exchange JSON file from json_exporter.js")
	pCdCExCSVTransfer := flag.String("cdc_ex_transfer", "", "Crypto.com Exchange Deposit/Withdrawal CSV file")
	pCdCExCSVStake := flag.String("cdc_ex_stake", "", "Crypto.com Exchange Stake CSV file")
	pCdCExCSVSupercharger := flag.String("cdc_ex_supercharger", "", "Crypto.com Exchange Supercharger CSV file")
	pKrakenAPIKey := flag.String("kraken_api_key", "", "Kraken API key")
	pKrakenAPISecret := flag.String("kraken_api_secret", "", "Kraken API secret")
	pKrakenCSV := flag.String("kraken", "", "Kraken CSV file")
	pLedgerLiveCSV := flag.String("ledgerlive", "", "LedgerLive CSV file")
	pLBCSVTrade := flag.String("lb_trade", "", "Local Bitcoin Trade CSV file")
	pLBCSVTransfer := flag.String("lb_transfer", "", "Local Bitcoin Transfer CSV file")
	pMyCeliumCSV := flag.String("mycelium", "", "MyCelium CSV file")
	pRevoCSV := flag.String("revolut", "", "Revolut CSV file")
	// Output
	p2086Display := flag.Bool("2086_display", false, "Display Cerfa 2086")
	p2086 := flag.Bool("2086", false, "Export Cerfa 2086 in 2086.xlsx")
	p3916 := flag.Bool("3916", false, "Export Cerfa 3916 in 3916.xlsx")
	pStock := flag.Bool("stock", false, "Export stock balances in stock.xlsx")
	flag.Parse()
	if *pCoinAPIKey != "" {
		wallet.CoinAPISetKey(*pCoinAPIKey)
	}
	if *pCoinLayerKey != "" {
		wallet.CoinLayerSetKey(*pCoinLayerKey)
	}
	categ := category.New()
	if *pTXsCategCSV != "" {
		recordFile, err := os.Open(*pTXsCategCSV)
		if err != nil {
			log.Fatal("Error opening Transactions CSV Category file:", err)
		}
		categ.ParseCSVCategory(recordFile)
	}
	loc, err := time.LoadLocation(*pLocation)
	if err != nil {
		log.Fatal("Error parsing Location:", err)
	}
	// Launch APIs access in go routines
	btc := btc.New()
	blkst := blockstream.New()
	if *pBTCAddressesCSV != "" {
		recordFile, err := os.Open(*pBTCAddressesCSV)
		if err != nil {
			log.Fatal("Error opening Bitcoin CSV Addresses file:", err)
		}
		btc.ParseCSVAddresses(recordFile)
		go blkst.GetAllTXs(btc, *categ)
	}
	ethsc := etherscan.New()
	if *pETHAddressesCSV != "" {
		recordFile, err := os.Open(*pETHAddressesCSV)
		if err != nil {
			log.Fatal("Error opening Ethereum CSV Addresses file:", err)
		}
		err = ethsc.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("")
		}
		ethsc.NewAPI(*pEtherscanAPIKey, *pDebug)
		go ethsc.GetAPITXs(*categ)
	}
	cdc := cryptocom.New()
	if *pCdCExAPIKey != "" && *pCdCExSecretKey != "" {
		cdc.NewExchangeAPI(*pCdCExAPIKey, *pCdCExSecretKey, *pDebug)
		fmt.Print("Début de récupération des TXs par l'API CdC Exchange (attention ce processus peut être long la première fois)")
		go cdc.GetAPIExchangeTXs(loc)
	}
	hb := hitbtc.New()
	if *pHitBtcAPIKey != "" && *pHitBtcSecretKey != "" {
		hb.NewAPI(*pHitBtcAPIKey, *pHitBtcSecretKey, *pDebug)
		go hb.GetAPIAllTXs()
	}
	kr := kraken.New()
	if *pKrakenAPIKey != "" && *pKrakenAPISecret != "" {
		kr.NewAPI(*pKrakenAPIKey, *pKrakenAPISecret, *pDebug)
		go kr.GetAPIAllTXs()
	}
	// Now parse local files
	bc := blockchain.New()
	if *pBTGTXsJSON != "" {
		jsonFile, err := os.Open(*pBTGTXsJSON)
		if err != nil {
			log.Fatal("Error opening Bitcoin Gold JSON Transactions file:", err)
		}
		err = bc.ParseTXsJSON(jsonFile, "BTG")
		if err != nil {
			log.Fatal("Error parsing Bitcoin Gold JSON Transactions file:", err)
		}
	}
	b := binance.New()
	if *pBinanceAPIKey != "" && *pBinanceSecretKey != "" {
		b.NewAPI(*pBinanceAPIKey, *pBinanceSecretKey, *pDebug)
		fmt.Println("Début de récupération des TXs par l'API Binance (attention ce processus peut être long la première fois)...")
		go b.GetAPIExchangeTXs(loc)
	}
	if *pBinanceCSV != "" {
		recordFile, err := os.Open(*pBinanceCSV)
		if err != nil {
			log.Fatal("Error opening Binance CSV file:", err)
		}
		err = b.ParseCSV(recordFile, *pBinanceCSVExtended)
		if err != nil {
			log.Fatal("Error parsing Binance CSV file:", err)
		}
	}
	bf := bitfinex.New()
	if *pBitfinexCSV != "" {
		recordFile, err := os.Open(*pBitfinexCSV)
		if err != nil {
			log.Fatal("Error opening Bitfinex CSV file:", err)
		}
		err = bf.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Bitfinex CSV file:", err)
		}
	}
	btrx := bittrex.New()
	if *pBittrexCSV != "" {
		recordFile, err := os.Open(*pBittrexCSV)
		if err != nil {
			log.Fatal("Error opening Bittrex CSV file:", err)
		}
		err = btrx.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Bittrex CSV file:", err)
		}
		if *pBittrexAPIKey != "" && *pBittrexAPISecret != "" {
			go btrx.GetAllTransferTXs(*pBittrexAPIKey, *pBittrexAPISecret, *categ)
			go btrx.GetAllTradeTXs(*pBittrexAPIKey, *pBittrexAPISecret, *categ)
		} else {
			log.Println("Warning, you should provide your API Key/Secret to retrieve Deposits and Withdrawals")
		}
	}
	cb := coinbase.New()
	if *pCoinbaseCSV != "" {
		recordFile, err := os.Open(*pCoinbaseCSV)
		if err != nil {
			log.Fatal("Error opening Coinbase CSV file:", err)
		}
		err = cb.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Coinbase CSV file:", err)
		}
	}
	if *pCdCAppCSVCrypto != "" {
		recordFile, err := os.Open(*pCdCAppCSVCrypto)
		if err != nil {
			log.Fatal("Error opening Crypto.com CSV file:", err)
		}
		err = cdc.ParseCSVAppCrypto(recordFile)
		if err != nil {
			log.Fatal("Error parsing Crypto.com CSV file:", err)
		}
	}
	if *pCdCExJSONExportJS != "" {
		recordFile, err := os.Open(*pCdCExJSONExportJS)
		if err != nil {
			log.Fatal("Error opening Crypto.com Exchange ExportJS JSON file:", err)
		}
		err = cdc.ParseJSONExchangeExportJS(recordFile)
		if err != nil {
			log.Fatal("Error parsing Crypto.com Exchange ExportJS JSON file:", err)
		}
	}
	if *pCdCExCSVTransfer != "" {
		recordFile, err := os.Open(*pCdCExCSVTransfer)
		if err != nil {
			log.Fatal("Error opening Crypto.com Exchange Deposit/Withdrawal CSV file:", err)
		}
		err = cdc.ParseCSVExchangeTransfer(recordFile)
		if err != nil {
			log.Fatal("Error parsing Crypto.com Exchange Deposit/Withdrawal CSV file:", err)
		}
	}
	if *pCdCExCSVStake != "" {
		recordFile, err := os.Open(*pCdCExCSVStake)
		if err != nil {
			log.Fatal("Error opening Crypto.com Exchange Stake CSV file:", err)
		}
		err = cdc.ParseCSVExchangeStake(recordFile)
		if err != nil {
			log.Fatal("Error parsing Crypto.com Exchange Stake CSV file:", err)
		}
	}
	if *pCdCExCSVSupercharger != "" {
		recordFile, err := os.Open(*pCdCExCSVSupercharger)
		if err != nil {
			log.Fatal("Error opening Crypto.com Exchange Supercharger CSV file:", err)
		}
		err = cdc.ParseCSVExchangeSupercharger(recordFile)
		if err != nil {
			log.Fatal("Error parsing Crypto.com Exchange Supercharger CSV file:", err)
		}
	}
	if *pHitBtcCSVTrades != "" {
		recordFile, err := os.Open(*pHitBtcCSVTrades)
		if err != nil {
			log.Fatal("Error opening HitBTC Trades CSV file:", err)
		}
		err = hb.ParseCSVTrades(recordFile)
		if err != nil {
			log.Fatal("Error parsing HitBTC Trades CSV file:", err)
		}
	}
	if *pHitBtcCSVTransactions != "" {
		recordFile, err := os.Open(*pHitBtcCSVTransactions)
		if err != nil {
			log.Fatal("Error opening HitBTC Transactions CSV file:", err)
		}
		err = hb.ParseCSVTransactions(recordFile)
		if err != nil {
			log.Fatal("Error parsing HitBTC Transactions CSV file:", err)
		}
	}
	if *pKrakenCSV != "" {
		recordFile, err := os.Open(*pKrakenCSV)
		if err != nil {
			log.Fatal("Error opening Kraken CSV file:", err)
		}
		err = kr.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Kraken CSV file:", err)
		}
	}
	ll := ledgerlive.New()
	if *pLedgerLiveCSV != "" {
		recordFile, err := os.Open(*pLedgerLiveCSV)
		if err != nil {
			log.Fatal("Error opening LedgerLive CSV file:", err)
		}
		err = ll.ParseCSV(recordFile, *categ)
		if err != nil {
			log.Fatal("Error parsing LedgerLive CSV file:", err)
		}
	}
	lb := localbitcoin.New()
	if *pLBCSVTrade != "" {
		recordFile, err := os.Open(*pLBCSVTrade)
		if err != nil {
			log.Fatal("Error opening Local Bitcoin Trade CSV file:", err)
		}
		err = lb.ParseTradeCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Local Bitcoin Trade CSV file:", err)
		}
	}
	if *pLBCSVTransfer != "" {
		recordFile, err := os.Open(*pLBCSVTransfer)
		if err != nil {
			log.Fatal("Error opening Local Bitcoin Transfer CSV file:", err)
		}
		err = lb.ParseTransferCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Local Bitcoin Transfer CSV file:", err)
		}
	}
	mc := mycelium.New()
	if *pMyCeliumCSV != "" {
		recordFile, err := os.Open(*pMyCeliumCSV)
		if err != nil {
			log.Fatal("Error opening MyCelium CSV file:", err)
		}
		err = mc.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing MyCelium CSV file:", err)
		}
	}
	revo := revolut.New()
	if *pRevoCSV != "" {
		recordFile, err := os.Open(*pRevoCSV)
		if err != nil {
			log.Fatal("Error opening Revolut CSV file:", err)
		}
		err = revo.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Revolut CSV file:", err)
		}
	}
	// Wait for API access to finish
	if *pCdCExAPIKey != "" && *pCdCExSecretKey != "" {
		err := cdc.WaitFinish()
		if err != nil {
			log.Fatal("Error getting Crypto.com Exchange API TXs:", err)
		}
	}
	if *pBinanceAPIKey != "" && *pBinanceSecretKey != "" {
		err := b.WaitFinish()
		if err != nil {
			log.Fatal("Error getting Binance API TXs:", err)
		}
	}
	if *pHitBtcAPIKey != "" && *pHitBtcSecretKey != "" {
		err := hb.WaitFinish()
		if err != nil {
			log.Fatal("Error getting HitBTC API TXs:", err)
		}
	}
	if *pETHAddressesCSV != "" {
		err := ethsc.WaitFinish()
		if err != nil {
			log.Fatal("Error parsing Ethereum CSV file:", err)
		}
	}
	if *pBittrexAPIKey != "" && *pBittrexAPISecret != "" {
		errTransfer := btrx.WaitTransfersFinish()
		if errTransfer != nil {
			log.Fatalln("Error parsing Bittrex API transfers:", errTransfer)
		}
		errTrades := btrx.WaitTradesFinish()
		if errTrades != nil {
			log.Fatalln("Error parsing Bittrex API trades:", errTrades)
		}
	}
	if *pKrakenAPIKey != "" && *pKrakenAPISecret != "" {
		err := kr.WaitFinish()
		if err != nil {
			log.Fatal("Error getting Kraken API TXs:", err)
		}
	}
	if *pBTCAddressesCSV != "" {
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
	if *p3916 {
		sources := make(source.Sources)
		sources.Add(b.Sources)
		sources.Add(bf.Sources)
		sources.Add(cb.Sources)
		sources.Add(cdc.Sources)
		sources.Add(hb.Sources)
		sources.Add(kr.Sources)
		sources.Add(lb.Sources)
		sources.Add(revo.Sources)
		err = sources.ToXlsx("3916.xlsx", loc)
		if err != nil {
			log.Fatal(err)
		}
	}
	// create Global Wallet up to Date
	global := make(wallet.TXsByCategory)
	global.Add(b.TXsByCategory)
	global.Add(bf.TXsByCategory)
	global.Add(btrx.TXsByCategory)
	global.Add(cb.TXsByCategory)
	global.Add(cdc.TXsByCategory)
	global.Add(hb.TXsByCategory)
	global.Add(kr.TXsByCategory)
	global.Add(ll.TXsByCategory)
	global.Add(lb.TXsByCategory)
	global.Add(mc.TXsByCategory)
	global.Add(revo.TXsByCategory)
	global.Add(ethsc.TXsByCategory)
	global.Add(btc.TXsByCategory)
	global.Add(bc.TXsByCategory)
	fmt.Print("Merging Deposits with Withdrawals into Transfers...")
	global.FindTransfers()
	fmt.Println("Finished")
	if *pStock {
		global.StockToXlsx("stock.xlsx")
	}
	var totalCommercialRebates, totalInterests, totalReferrals decimal.Decimal
	if *p2086 || *p2086Display {
		fmt.Print("Look for CashIn and CashOut...")
		totalCommercialRebates, totalInterests, totalReferrals = global.FindCashInOut(*pNative)
		fmt.Println("Finished")
	}
	global.SortTXsByDate(true)
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
	if *p2086 || *p2086Display {
		var c2086 Cerfa2086
		fmt.Print("Début du calcul pour le 2086...")
		err = c2086.CalculatePVMV(global, *pNative, loc)
		fmt.Println("Fini")
		if err != nil {
			log.Fatal(err)
		}
		if *p2086Display {
			c2086.Println()
		}
		if *p2086 {
			c2086.ToXlsx("2086.xlsx", *pNative)
		}
	}
	os.Exit(0)
}
