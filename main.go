package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fiscafacile/CryptoFiscaFacile/binance"
	"github.com/fiscafacile/CryptoFiscaFacile/bitfinex"
	"github.com/fiscafacile/CryptoFiscaFacile/bitstamp"
	"github.com/fiscafacile/CryptoFiscaFacile/bittrex"
	"github.com/fiscafacile/CryptoFiscaFacile/blockchain"
	"github.com/fiscafacile/CryptoFiscaFacile/blockstream"
	"github.com/fiscafacile/CryptoFiscaFacile/btc"
	"github.com/fiscafacile/CryptoFiscaFacile/category"
	"github.com/fiscafacile/CryptoFiscaFacile/cfg"
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
	// Configuration
	config, err := cfg.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	if config.Tools.CoinAPI.Key != "" {
		wallet.CoinAPISetKey(config.Tools.CoinAPI.Key)
	}
	if config.Tools.CoinLayer.Key != "" {
		wallet.CoinLayerSetKey(config.Tools.CoinLayer.Key)
	}
	categ := category.New()
	if config.Options.TxsCategory != "" {
		recordFile, err := os.Open(config.Options.TxsCategory)
		if err != nil {
			log.Fatal("Error opening Transactions CSV Category file:", err)
		}
		categ.ParseCSVCategory(recordFile)
	}
	loc, err := time.LoadLocation(config.Options.Location)
	if err != nil {
		log.Fatal("Error parsing Location:", err)
	}
	// Launch APIs access in go routines
	btc := btc.New()
	blkst := blockstream.New()
	for _, file := range config.Blockchains.BTC.CSV {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Bitcoin CSV Addresses file:", err)
		}
		btc.ParseCSVAddresses(recordFile)
		go blkst.GetAllTXs(btc, *categ)
	}
	ethsc := etherscan.New()
	for _, file := range config.Blockchains.ETH.CSV {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Ethereum CSV Addresses file:", err)
		}
		err = ethsc.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("")
		}
	}
	if len(config.Blockchains.ETH.CSV) > 0 {
		ethsc.NewAPI(config.Tools.EtherScan.Key, config.Options.Debug)
		go ethsc.GetAPITXs(*categ)
	}
	bs := bitstamp.New()
	if config.Exchanges.Bitstamp.API.Key != "" && config.Exchanges.Bitstamp.API.Secret != "" {
		bs.NewAPI(config.Exchanges.Bitstamp.API.Key, config.Exchanges.Bitstamp.API.Secret, config.Options.Debug)
		go bs.GetAPIAllTXs()
	}
	btrx := bittrex.New()
	if config.Exchanges.Bittrex.API.Key != "" && config.Exchanges.Bittrex.API.Secret != "" {
		btrx.NewAPI(config.Exchanges.Bittrex.API.Key, config.Exchanges.Bittrex.API.Secret, config.Options.Debug)
		go btrx.GetAPIAllTXs()
	}
	cdc := cryptocom.New()
	if config.Exchanges.CdcEx.API.Key != "" && config.Exchanges.CdcEx.API.Secret != "" {
		cdc.NewExchangeAPI(config.Exchanges.CdcEx.API.Key, config.Exchanges.CdcEx.API.Secret, config.Options.Debug)
		fmt.Print("Début de récupération des TXs par l'API CdC Exchange (attention ce processus peut être long la première fois)")
		go cdc.GetAPIExchangeTXs(loc)
	}
	hb := hitbtc.New()
	if config.Exchanges.HitBTC.API.Key != "" && config.Exchanges.HitBTC.API.Secret != "" {
		hb.NewAPI(config.Exchanges.HitBTC.API.Key, config.Exchanges.HitBTC.API.Secret, config.Options.Debug)
		go hb.GetAPIAllTXs()
	}
	kr := kraken.New()
	if config.Exchanges.Kraken.API.Key != "" && config.Exchanges.Kraken.API.Secret != "" {
		kr.NewAPI(config.Exchanges.Kraken.API.Key, config.Exchanges.Kraken.API.Secret, config.Options.Debug)
		go kr.GetAPIAllTXs()
	}
	// Now parse local files
	bc := blockchain.New()
	if config.Blockchains.BTG.JSON != "" {
		jsonFile, err := os.Open(config.Blockchains.BTG.JSON)
		if err != nil {
			log.Fatal("Error opening Bitcoin Gold JSON Transactions file:", err)
		}
		err = bc.ParseTXsJSON(jsonFile, "BTG")
		if err != nil {
			log.Fatal("Error parsing Bitcoin Gold JSON Transactions file:", err)
		}
	}
	b := binance.New()
	for _, file := range config.Exchanges.Binance.CSV.All {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Binance CSV file:", err)
		}
		err = b.ParseCSV(recordFile, config.Options.BinanceExtended)
		if err != nil {
			log.Fatal("Error parsing Binance CSV file:", err)
		}
	}
	bf := bitfinex.New()
	for _, file := range config.Exchanges.Bitfinex.CSV.All {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Bitfinex CSV file:", err)
		}
		err = bf.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Bitfinex CSV file:", err)
		}
	}
	for _, file := range config.Exchanges.Bitstamp.CSV.All {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Bitstamp CSV file:", err)
		}
		err = bs.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Bitstamp CSV file:", err)
		}
	}
	for _, file := range config.Exchanges.Bittrex.CSV.All {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Bittrex CSV file:", err)
		}
		err = btrx.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Bittrex CSV file:", err)
		}
	}
	cb := coinbase.New()
	for _, file := range config.Exchanges.Coinbase.CSV.All {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Coinbase CSV file:", err)
		}
		err = cb.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Coinbase CSV file:", err)
		}
	}
	for _, file := range config.Exchanges.CdcApp.CSV.All {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Crypto.com CSV file:", err)
		}
		err = cdc.ParseCSVAppCrypto(recordFile, *categ)
		if err != nil {
			log.Fatal("Error parsing Crypto.com CSV file:", err)
		}
	}
	if config.Exchanges.CdcEx.JSON != "" {
		recordFile, err := os.Open(config.Exchanges.CdcEx.JSON)
		if err != nil {
			log.Fatal("Error opening Crypto.com Exchange ExportJS JSON file:", err)
		}
		err = cdc.ParseJSONExchangeExportJS(recordFile)
		if err != nil {
			log.Fatal("Error parsing Crypto.com Exchange ExportJS JSON file:", err)
		}
	}
	for _, file := range config.Exchanges.CdcApp.CSV.Transfers {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Crypto.com Exchange Deposit/Withdrawal CSV file:", err)
		}
		err = cdc.ParseCSVExchangeTransfer(recordFile)
		if err != nil {
			log.Fatal("Error parsing Crypto.com Exchange Deposit/Withdrawal CSV file:", err)
		}
	}
	for _, file := range config.Exchanges.CdcApp.CSV.Staking {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Crypto.com Exchange Stake CSV file:", err)
		}
		err = cdc.ParseCSVExchangeStake(recordFile)
		if err != nil {
			log.Fatal("Error parsing Crypto.com Exchange Stake CSV file:", err)
		}
	}
	for _, file := range config.Exchanges.CdcApp.CSV.Supercharger {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Crypto.com Exchange Supercharger CSV file:", err)
		}
		err = cdc.ParseCSVExchangeSupercharger(recordFile)
		if err != nil {
			log.Fatal("Error parsing Crypto.com Exchange Supercharger CSV file:", err)
		}
	}
	for _, file := range config.Exchanges.HitBTC.CSV.Trades {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening HitBTC Trades CSV file:", err)
		}
		err = hb.ParseCSVTrades(recordFile)
		if err != nil {
			log.Fatal("Error parsing HitBTC Trades CSV file:", err)
		}
	}
	for _, file := range config.Exchanges.HitBTC.CSV.Transfers {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening HitBTC Transactions CSV file:", err)
		}
		err = hb.ParseCSVTransactions(recordFile)
		if err != nil {
			log.Fatal("Error parsing HitBTC Transactions CSV file:", err)
		}
	}
	for _, file := range config.Exchanges.Kraken.CSV.Transfers {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Kraken CSV file:", err)
		}
		err = kr.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Kraken CSV file:", err)
		}
	}
	ll := ledgerlive.New()
	for _, file := range config.Wallets.LedgerLive.CSV.All {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening LedgerLive CSV file:", err)
		}
		err = ll.ParseCSV(recordFile, *categ)
		if err != nil {
			log.Fatal("Error parsing LedgerLive CSV file:", err)
		}
	}
	lb := localbitcoin.New()
	for _, file := range config.Exchanges.LocalBitcoins.CSV.Trades {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Local Bitcoin Trade CSV file:", err)
		}
		err = lb.ParseTradeCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Local Bitcoin Trade CSV file:", err)
		}
	}
	for _, file := range config.Exchanges.LocalBitcoins.CSV.Transfers {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Local Bitcoin Transfer CSV file:", err)
		}
		err = lb.ParseTransferCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Local Bitcoin Transfer CSV file:", err)
		}
	}
	mc := mycelium.New()
	for _, file := range config.Wallets.MyCelium.CSV.All {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening MyCelium CSV file:", err)
		}
		err = mc.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing MyCelium CSV file:", err)
		}
	}
	revo := revolut.New()
	for _, file := range config.Exchanges.Revolut.CSV.All {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Revolut CSV file:", err)
		}
		err = revo.ParseCSV(recordFile)
		if err != nil {
			log.Fatal("Error parsing Revolut CSV file:", err)
		}
	}
	// Wait for API access to finish
	if config.Exchanges.Bitstamp.API.Key != "" && config.Exchanges.Bitstamp.API.Secret != "" {
		err := bs.WaitFinish()
		if err != nil {
			log.Fatal("Error getting BiTstamp API TXs:", err)
		}
	}
	if config.Exchanges.CdcEx.API.Key != "" && config.Exchanges.CdcEx.API.Secret != "" {
		err := cdc.WaitFinish()
		if err != nil {
			log.Fatal("Error getting Crypto.com Exchange API TXs:", err)
		}
	}
	if config.Exchanges.HitBTC.API.Key != "" && config.Exchanges.HitBTC.API.Secret != "" {
		err := hb.WaitFinish()
		if err != nil {
			log.Fatal("Error getting HitBTC API TXs:", err)
		}
	}
	if len(config.Blockchains.ETH.CSV) > 0 {
		err := ethsc.WaitFinish()
		if err != nil {
			log.Fatal("Error parsing Ethereum CSV file:", err)
		}
	}
	if config.Exchanges.Bittrex.API.Key != "" && config.Exchanges.Bittrex.API.Secret != "" {
		err := btrx.WaitFinish()
		if err != nil {
			log.Fatal("Error getting Bittrex API TXs:", err)
		}
	}
	if config.Exchanges.Kraken.API.Key != "" && config.Exchanges.Kraken.API.Secret != "" {
		err := kr.WaitFinish()
		if err != nil {
			log.Fatal("Error getting Kraken API TXs:", err)
		}
	}
	if len(config.Blockchains.BTC.CSV) > 0 {
		err := blkst.WaitFinish()
		if err != nil {
			log.Fatal("Error parsing Bitcoin CSV file:", err)
		}
		if config.Options.Bcd {
			blkst.DetectBCD(btc)
		}
		if config.Options.Bch {
			blkst.DetectBCH(btc)
		}
		if config.Options.Btg {
			blkst.DetectBTG(btc)
		}
		if config.Options.Lbtc {
			blkst.DetectLBTC(btc)
		}
	}
	if config.Options.Export3916 {
		sources := make(source.Sources)
		sources.Add(b.Sources)
		sources.Add(bf.Sources)
		sources.Add(bs.Sources)
		sources.Add(btrx.Sources)
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
	global.Add(bs.TXsByCategory)
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
	if config.Options.ExportStock {
		global.StockToXlsx("stock.xlsx")
	}
	var totalCommercialRebates, totalInterests, totalReferrals decimal.Decimal
	if config.Options.Export2086 || config.Options.Display2086 {
		fmt.Print("Look for CashIn and CashOut...")
		totalCommercialRebates, totalInterests, totalReferrals = global.FindCashInOut(config.Options.Native)
		fmt.Println("Finished")
	}
	global.SortTXsByDate(true)
	if config.Options.Stats {
		global.PrintStats(config.Options.Native, totalCommercialRebates, totalInterests, totalReferrals)
	}
	if config.Options.Check {
		global.CheckConsistency(loc)
	}
	// Debug
	if config.Options.TxsDisplay != "" {
		if config.Options.TxsDisplay == "Alls" {
			global.Println(config.Options.CurrencyFilter)
		} else {
			global[config.Options.TxsDisplay].Println("Category "+config.Options.TxsDisplay, config.Options.CurrencyFilter)
		}
	}
	// Construct global wallet up to date
	filterDate, err := time.ParseInLocation("2006-01-02T15:04:05", config.Options.Date, loc)
	if err != nil {
		log.Fatal("Error parsing Date:", err)
	}
	globalWallet := global.GetWallets(filterDate, false, !config.Options.Exact)
	globalWallet.Println("Global Crypto", "")
	globalWalletTotalValue, err := globalWallet.CalculateTotalValue(config.Options.Native)
	if err != nil {
		log.Fatal("Error Calculating Global Wallet:", err)
	} else {
		globalWalletTotalValue.Amount = globalWalletTotalValue.Amount.RoundBank(0)
		fmt.Print("Total Value : ")
		globalWalletTotalValue.Println("")
	}
	if config.Options.Export2086 || config.Options.Display2086 {
		var c2086 Cerfa2086
		fmt.Print("Début du calcul pour le 2086...")
		err = c2086.CalculatePVMV(global, config.Options.Native, loc)
		fmt.Println("Fini")
		if err != nil {
			log.Fatal(err)
		}
		if config.Options.Display2086 {
			c2086.Println()
		}
		if config.Options.Export2086 {
			c2086.ToXlsx("2086.xlsx", config.Options.Native)
		}
	}
	os.Exit(0)
}
