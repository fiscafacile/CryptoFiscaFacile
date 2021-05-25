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
	"github.com/fiscafacile/CryptoFiscaFacile/poloniex"
	"github.com/fiscafacile/CryptoFiscaFacile/revolut"
	"github.com/fiscafacile/CryptoFiscaFacile/source"
	"github.com/fiscafacile/CryptoFiscaFacile/uphold"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
)

func main() {
	// Configuration
	config, err := cfg.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	if config.Options.LogFile != "" {
		f, err := os.OpenFile(config.Options.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
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
	btc.AddListAddresses(config.Blockchains.BTC.Addresses)
	for _, file := range config.Blockchains.BTC.CSV {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Bitcoin CSV Addresses file:", err)
		}
		err = btc.ParseCSVAddresses(recordFile)
		if err != nil {
			log.Fatal("")
		}
	}
	blkst := blockstream.New()
	if len(config.Blockchains.BTC.CSV)+len(config.Blockchains.BTC.Addresses) > 0 {
		go blkst.GetAllTXs(btc, *categ)
	}
	ethsc := etherscan.New()
	ethsc.AddListAddresses(config.Blockchains.ETH.Addresses)
	for _, file := range config.Blockchains.ETH.CSV {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Ethereum CSV Addresses file:", err)
		}
		err = ethsc.ParseCSVAddresses(recordFile)
		if err != nil {
			log.Fatal("")
		}
	}
	if len(config.Blockchains.ETH.CSV)+len(config.Blockchains.ETH.Addresses) > 0 {
		ethsc.NewAPI(config.Tools.EtherScan.Key, config.Options.Debug)
		go ethsc.GetAPITXs(*categ)
	}
	b := binance.New()
	if config.Exchanges.Binance.API.Key != "" && config.Exchanges.Binance.API.Secret != "" {
		b.NewAPI(config.Exchanges.Binance.API.Key, config.Exchanges.Binance.API.Secret, config.Options.Debug)
		fmt.Print("Début de récupération des TXs par l'API Binance (attention ce processus peut être long la première fois)...")
		go b.GetAPIAllTXs(loc)
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
	for _, file := range config.Exchanges.Binance.CSV.All {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Binance CSV file:", err)
		}
		err = b.ParseCSV(recordFile, config.Options.BinanceExtended, config.Exchanges.Binance.Account)
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
		err = bf.ParseCSV(recordFile, config.Exchanges.Bitfinex.Account)
		if err != nil {
			log.Fatal("Error parsing Bitfinex CSV file:", err)
		}
	}
	for _, file := range config.Exchanges.Bitstamp.CSV.All {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Bitstamp CSV file:", err)
		}
		err = bs.ParseCSV(recordFile, config.Exchanges.Bitstamp.Account)
		if err != nil {
			log.Fatal("Error parsing Bitstamp CSV file:", err)
		}
	}
	for _, file := range config.Exchanges.Bittrex.CSV.All {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Bittrex CSV file:", err)
		}
		err = btrx.ParseCSV(recordFile, *categ, config.Exchanges.Bittrex.Account)
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
		err = cb.ParseCSV(recordFile, config.Exchanges.Coinbase.Account)
		if err != nil {
			log.Fatal("Error parsing Coinbase CSV file:", err)
		}
	}
	for _, file := range config.Exchanges.CdcApp.CSV.All {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Crypto.com CSV file:", err)
		}
		err = cdc.ParseCSVAppCrypto(recordFile, *categ, config.Exchanges.CdcApp.Account)
		if err != nil {
			log.Fatal("Error parsing Crypto.com CSV file:", err)
		}
	}
	if config.Exchanges.CdcEx.JSON != "" {
		recordFile, err := os.Open(config.Exchanges.CdcEx.JSON)
		if err != nil {
			log.Fatal("Error opening Crypto.com Exchange ExportJS JSON file:", err)
		}
		err = cdc.ParseJSONExchangeExportJS(recordFile, config.Exchanges.CdcEx.Account)
		if err != nil {
			log.Fatal("Error parsing Crypto.com Exchange ExportJS JSON file:", err)
		}
	}
	for _, file := range config.Exchanges.CdcEx.CSV.Transfers {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Crypto.com Exchange Deposit/Withdrawal CSV file:", err)
		}
		err = cdc.ParseCSVExchangeTransfer(recordFile)
		if err != nil {
			log.Fatal("Error parsing Crypto.com Exchange Deposit/Withdrawal CSV file:", err)
		}
	}
	for _, file := range config.Exchanges.CdcEx.CSV.Staking {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Crypto.com Exchange Stake CSV file:", err)
		}
		err = cdc.ParseCSVExchangeStake(recordFile)
		if err != nil {
			log.Fatal("Error parsing Crypto.com Exchange Stake CSV file:", err)
		}
	}
	for _, file := range config.Exchanges.CdcEx.CSV.Supercharger {
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
	for _, file := range config.Exchanges.Kraken.CSV.All {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Kraken CSV file:", err)
		}
		err = kr.ParseCSV(recordFile, *categ, config.Exchanges.Kraken.Account)
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
		err = lb.ParseTradeCSV(recordFile, config.Exchanges.LocalBitcoins.Account)
		if err != nil {
			log.Fatal("Error parsing Local Bitcoin Trade CSV file:", err)
		}
	}
	for _, file := range config.Exchanges.LocalBitcoins.CSV.Transfers {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Local Bitcoin Transfer CSV file:", err)
		}
		err = lb.ParseTransferCSV(recordFile, config.Exchanges.LocalBitcoins.Account)
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
	pl := poloniex.New()
	for _, file := range config.Exchanges.Poloniex.CSV.Deposits {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Poloniex Deposits CSV file:", err)
		}
		err = pl.ParseDepositsCSV(recordFile, config.Exchanges.Poloniex.Account)
		if err != nil {
			log.Fatal("Error parsing Poloniex Deposits CSV file:", err)
		}
	}
	for _, file := range config.Exchanges.Poloniex.CSV.Distributions {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Poloniex Distributions CSV file:", err)
		}
		err = pl.ParseDistributionsCSV(recordFile, config.Exchanges.Poloniex.Account)
		if err != nil {
			log.Fatal("Error parsing Poloniex Distributions CSV file:", err)
		}
	}
	for _, file := range config.Exchanges.Poloniex.CSV.Trades {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Poloniex Trades CSV file:", err)
		}
		err = pl.ParseTradesCSV(recordFile, *categ, config.Exchanges.Poloniex.Account)
		if err != nil {
			log.Fatal("Error parsing Poloniex Trades CSV file:", err)
		}
	}
	for _, file := range config.Exchanges.Poloniex.CSV.Withdrawals {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Poloniex Withdrawals CSV file:", err)
		}
		err = pl.ParseWithdrawalsCSV(recordFile, *categ, config.Exchanges.Poloniex.Account)
		if err != nil {
			log.Fatal("Error parsing Poloniex Withdrawals CSV file:", err)
		}
	}
	revo := revolut.New()
	for _, file := range config.Exchanges.Revolut.CSV.All {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Revolut CSV file:", err)
		}
		err = revo.ParseCSV(recordFile, config.Exchanges.Revolut.Account)
		if err != nil {
			log.Fatal("Error parsing Revolut CSV file:", err)
		}
	}
	uh := uphold.New()
	for _, file := range config.Exchanges.Uphold.CSV.All {
		recordFile, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening Uphold CSV file:", err)
		}
		err = uh.ParseCSV(recordFile, config.Exchanges.Uphold.Account)
		if err != nil {
			log.Fatal("Error parsing Uphold CSV file:", err)
		}
	}
	// Wait for API access to finish
	if config.Exchanges.Binance.API.Key != "" && config.Exchanges.Binance.API.Secret != "" {
		err := b.WaitFinish(config.Exchanges.Bitstamp.Account)
		if err != nil {
			log.Fatal("Error getting Binance API TXs:", err)
		}
	}
	if config.Exchanges.Bitstamp.API.Key != "" && config.Exchanges.Bitstamp.API.Secret != "" {
		err := bs.WaitFinish(config.Exchanges.Bitstamp.Account)
		if err != nil {
			log.Fatal("Error getting BiTstamp API TXs:", err)
		}
	}
	if config.Exchanges.Bittrex.API.Key != "" && config.Exchanges.Bittrex.API.Secret != "" {
		err := btrx.WaitFinish(config.Exchanges.Bittrex.Account)
		if err != nil {
			log.Fatal("Error getting Bittrex API TXs:", err)
		}
	}
	if config.Exchanges.CdcEx.API.Key != "" && config.Exchanges.CdcEx.API.Secret != "" {
		err := cdc.WaitFinish(config.Exchanges.CdcEx.Account)
		if err != nil {
			log.Fatal("Error getting Crypto.com Exchange API TXs:", err)
		}
	}
	if config.Exchanges.HitBTC.API.Key != "" && config.Exchanges.HitBTC.API.Secret != "" {
		err := hb.WaitFinish(config.Exchanges.HitBTC.Account)
		if err != nil {
			log.Fatal("Error getting HitBTC API TXs:", err)
		}
	}
	if len(config.Blockchains.ETH.CSV)+len(config.Blockchains.ETH.Addresses) > 0 {
		err := ethsc.WaitFinish()
		if err != nil {
			log.Fatal("Error parsing Ethereum CSV file:", err)
		}
	}
	if config.Exchanges.Kraken.API.Key != "" && config.Exchanges.Kraken.API.Secret != "" {
		err := kr.WaitFinish(config.Exchanges.Kraken.Account)
		if err != nil {
			log.Fatal("Error getting Kraken API TXs:", err)
		}
	}
	if len(config.Blockchains.BTC.CSV)+len(config.Blockchains.BTC.Addresses) > 0 {
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
	// Set delisted coins balances to zero
	if len(config.Exchanges.Binance.DelistedCoins) > 0 {
		for _, dc := range config.Exchanges.Binance.DelistedCoins {
			b.TXsByCategory.RemoveDelistedCoins(dc)
		}
	}
	if len(config.Exchanges.Bitstamp.DelistedCoins) > 0 {
		for _, dc := range config.Exchanges.Bitstamp.DelistedCoins {
			bs.TXsByCategory.RemoveDelistedCoins(dc)
		}
	}
	if len(config.Exchanges.Bittrex.DelistedCoins) > 0 {
		for _, dc := range config.Exchanges.Bittrex.DelistedCoins {
			btrx.TXsByCategory.RemoveDelistedCoins(dc)
		}
	}
	if len(config.Exchanges.CdcEx.DelistedCoins) > 0 {
		for _, dc := range config.Exchanges.CdcEx.DelistedCoins {
			cdc.TXsByCategory.RemoveDelistedCoins(dc)
		}
	}
	if len(config.Exchanges.CdcApp.DelistedCoins) > 0 {
		for _, dc := range config.Exchanges.CdcApp.DelistedCoins {
			cdc.TXsByCategory.RemoveDelistedCoins(dc)
		}
	}
	if len(config.Exchanges.HitBTC.DelistedCoins) > 0 {
		for _, dc := range config.Exchanges.HitBTC.DelistedCoins {
			hb.TXsByCategory.RemoveDelistedCoins(dc)
		}
	}
	if len(config.Exchanges.Kraken.DelistedCoins) > 0 {
		for _, dc := range config.Exchanges.Kraken.DelistedCoins {
			kr.TXsByCategory.RemoveDelistedCoins(dc)
		}
	}
	if len(config.Exchanges.Bitfinex.DelistedCoins) > 0 {
		for _, dc := range config.Exchanges.Bitfinex.DelistedCoins {
			bf.TXsByCategory.RemoveDelistedCoins(dc)
		}
	}
	if len(config.Exchanges.Coinbase.DelistedCoins) > 0 {
		for _, dc := range config.Exchanges.Coinbase.DelistedCoins {
			cb.TXsByCategory.RemoveDelistedCoins(dc)
		}
	}
	if len(config.Exchanges.Poloniex.DelistedCoins) > 0 {
		for _, dc := range config.Exchanges.Poloniex.DelistedCoins {
			pl.TXsByCategory.RemoveDelistedCoins(dc)
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
		sources.Add(pl.Sources)
		sources.Add(revo.Sources)
		sources.Add(uh.Sources)
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
	global.Add(pl.TXsByCategory)
	global.Add(revo.TXsByCategory)
	global.Add(uh.TXsByCategory)
	global.Add(ethsc.TXsByCategory)
	global.Add(btc.TXsByCategory)
	global.Add(bc.TXsByCategory)
	fmt.Print("Merging Deposits with Withdrawals into Transfers...")
	global.FindTransfers()
	fmt.Println("Finished")
	if config.Options.ExportStock {
		global.StockToXlsx("stock.xlsx")
	}
	if config.Options.Export2086 || config.Options.Display2086 {
		fmt.Print("Look for CashIn and CashOut...")
		global.FindCashInOut(config.Options.Native)
		fmt.Println("Finished")
	}
	global.SortTXsByDate(true)
	if config.Options.Stats {
		global.PrintStats(config.Options.Native)
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
	fmt.Print("Calculating Total native value...")
	globalWalletTotalValue, err := globalWallet.CalculateTotalValue(config.Options.Native)
	if err != nil {
		fmt.Println("Error")
		log.Fatal("Error Calculating Global Wallet:", err)
	} else {
		fmt.Println("Finished")
		globalWalletTotalValue.Amount = globalWalletTotalValue.Amount.RoundBank(0)
		fmt.Print("Total Value : ")
		globalWalletTotalValue.Println("")
	}
	if config.Options.Export2086 || config.Options.Display2086 {
		c2086 := New2086()
		fmt.Print("Début du calcul pour le 2086...")
		err = c2086.CalculatePVMV(global, config.Options.Native, loc, config.Options.CashInBNC)
		fmt.Println("Fini")
		if err != nil {
			log.Fatal(err)
		}
		if config.Options.Display2086 {
			c2086.Println(config.Options.Native)
		}
		if config.Options.Export2086 {
			c2086.ToXlsx("2086.xlsx", config.Options.Native)
		}
	}
	os.Exit(0)
}
