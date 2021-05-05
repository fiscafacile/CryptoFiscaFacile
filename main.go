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
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
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
	p2086 := flag.Bool("2086", false, "Display Cerfa 2086")
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
		fmt.Println("Début de récupération des TXs par l'API CdC Exchange (attention ce processus peut être long la première fois)...")
		go cdc.GetAPIExchangeTXs(loc)
	}
	hb := hitbtc.New()
	if *pHitBtcAPIKey != "" && *pHitBtcSecretKey != "" {
		hb.NewAPI(*pHitBtcAPIKey, *pHitBtcSecretKey, *pDebug)
		go hb.GetAPIAllTXs()
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
	if *pBinanceCSV != "" {
		recordFile, err := os.Open(*pBinanceCSV)
		if err != nil {
			log.Fatal("Error opening Binance CSV file:", err)
		}
		if *pBinanceCSVExtended {
			err = b.ParseCSVExtended(recordFile)
		} else {
			err = b.ParseCSV(recordFile)
		}
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
		err = hb.ParseCSVTransactions(recordFile)
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
	kr := kraken.New()
	if *pKrakenAPIKey != "" && *pKrakenAPISecret != "" {
		kr.NewAPI(*pKrakenAPIKey, *pKrakenAPISecret, *pDebug)
		kr.GetAPITxs()
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
	global.FindTransfers()
	totalCommercialRebates, totalInterests, totalReferrals := global.FindCashInOut(*pNative)
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
