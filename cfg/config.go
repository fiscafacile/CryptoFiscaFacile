package cfg

import (
	"log"
	"os"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

type API struct {
	Key    string `yaml:"key"`
	Secret string `yaml:"secret"`
}

// Blockchains
type BlockchainConfig struct {
	Addresses []string `yaml:"addresses"`
	CSV       []string `yaml:"csv"`
	JSON      string   `yaml:"json"`
}

type Blockchains struct {
	ADA BlockchainConfig `yaml:"ADA"`
	BTC BlockchainConfig `yaml:"BTC"`
	BTG BlockchainConfig `yaml:"BTG"`
	ETH BlockchainConfig `yaml:"ETH"`
}

// Exchanges
type CSV struct {
	All           []string `yaml:"all"`
	Staking       []string `yaml:"staking"`
	Supercharger  []string `yaml:"supercharger"`
	Trades        []string `yaml:"trades"`
	Transfers     []string `yaml:"transfers"`
	Deposits      []string `yaml:"deposits"`
	Withdrawals   []string `yaml:"withdrawals"`
	Distributions []string `yaml:"distributions"`
}

type ExchangeConfig struct {
	CSV           CSV      `yaml:"csv"`
	API           API      `yaml:"api"`
	JSON          string   `yaml:"json"`
	Account       string   `yaml:"account"`
	DelistedCoins []string `yaml:"delisted-coins"`
}

type Exchanges struct {
	Binance       ExchangeConfig `yaml:"binance"`
	Bitfinex      ExchangeConfig `yaml:"bitfinex"`
	Bitstamp      ExchangeConfig `yaml:"bitstamp"`
	Bittrex       ExchangeConfig `yaml:"bittrex"`
	CdcApp        ExchangeConfig `yaml:"cdc-app"`
	CdcEx         ExchangeConfig `yaml:"cdc-exchange"`
	Coinbase      ExchangeConfig `yaml:"coinbase"`
	CoinbasePro   ExchangeConfig `yaml:"coinbase-pro"`
	HitBTC        ExchangeConfig `yaml:"hitbtc"`
	Kraken        ExchangeConfig `yaml:"kraken"`
	LocalBitcoins ExchangeConfig `yaml:"localbitcoins"`
	Poloniex      ExchangeConfig `yaml:"poloniex"`
	Revolut       ExchangeConfig `yaml:"revolut"`
	Uphold        ExchangeConfig `yaml:"uphold"`
}

type FiscalYear struct {
	Y2019 bool `yaml:"2019"`
	Y2020 bool `yaml:"2020"`
}

// Options
type Options struct {
	Bcd             bool       `yaml:"bcd"`
	Bch             bool       `yaml:"bch"`
	BinanceExtended bool       `yaml:"binance-extended"`
	Btg             bool       `yaml:"btg"`
	CashInBNC       FiscalYear `yaml:"cashin-bnc"`
	Check           bool       `yaml:"check"`
	CurrencyFilter  string     `yaml:"curr-filter"`
	Date            string     `yaml:"date"`
	Debug           bool       `yaml:"debug"`
	Display2086     bool       `yaml:"display-2086"`
	Exact           bool       `yaml:"exact"`
	Export2086      bool       `yaml:"export-2086"`
	Export3916      bool       `yaml:"export-3916"`
	ExportStock     bool       `yaml:"export-stock"`
	Lbtc            bool       `yaml:"lbtc"`
	Location        string     `yaml:"location"`
	LogFile         string     `yaml:"log"`
	Native          string     `yaml:"native"`
	Stats           bool       `yaml:"stats"`
	TxsCategory     string     `yaml:"txs-categ"`
	TxsDisplay      string     `yaml:"txs-display"`
}

// Tools
type Tools struct {
	CoinAPI   API `yaml:"coinapi"`
	CoinLayer API `yaml:"coinlayer"`
	EtherScan API `yaml:"etherscan"`
}

// Wallets
type WalletConfig struct {
	CSV CSV `yaml:"csv"`
}

type Wallets struct {
	LedgerLive WalletConfig `yaml:"ledgerlive"`
	Monero     WalletConfig `yaml:"monero"`
	MyCelium   WalletConfig `yaml:"mycelium"`
}

type Config struct {
	Blockchains Blockchains `yaml:"blockchains"`
	Exchanges   Exchanges   `yaml:"exchanges"`
	Options     Options     `yaml:"options"`
	Tools       Tools       `yaml:"tools"`
	Wallets     Wallets     `yaml:"wallets"`
}

func loadFile() *Config {
	config := &Config{}
	configPath := "config.yml"
	// Check that file exist and is readable
	s, err := os.Stat(configPath)
	if err != nil {
		log.Printf("Error while reading the configuration file, '%s' not loaded", configPath)
		return config
	}
	if s.IsDir() {
		log.Printf("'%s' is a directory, no file loaded", configPath)
		return config
	}
	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return config
	}
	defer file.Close()
	// Load configuration from file
	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return config
	}
	return config
}

func LoadConfig() (*Config, error) {
	// Load configuration from file
	config := loadFile()
	// Configure the default options
	if config.Options.Date == "" {
		config.Options.Date = "2021-01-01T00:00:00"
	}
	if config.Options.Location == "" {
		config.Options.Location = "Europe/Paris"
	}
	if config.Options.Native == "" {
		config.Options.Native = "EUR"
	}

	// Override configuration from CLI
	// General Options
	pflag.StringVar(&config.Options.Date, "date", config.Options.Date, "Date Filter")
	pflag.StringVar(&config.Options.Location, "location", config.Options.Location, "Date Filter Location")
	pflag.StringVar(&config.Options.Native, "native", config.Options.Native, "Native Currency for consolidation")
	pflag.BoolVarP(&config.Options.Stats, "stats", "s", config.Options.Stats, "Display accounts stats")
	// Debug
	pflag.BoolVarP(&config.Options.Debug, "debug", "d", config.Options.Debug, "Debug Mode (only for devs)")
	pflag.BoolVar(&config.Options.CashInBNC.Y2019, "cashin-bnc-2019", config.Options.CashInBNC.Y2019, "Convert AirDrops/CommercialRebates/Interets/Minings/Referrals into CashIn for 2019's Txs in 2086")
	pflag.BoolVar(&config.Options.CashInBNC.Y2020, "cashin-bnc-2020", config.Options.CashInBNC.Y2020, "Convert AirDrops/CommercialRebates/Interets/Minings/Referrals into CashIn for 2020's Txs in 2086")
	pflag.BoolVarP(&config.Options.Check, "check", "c", config.Options.Check, "Check and Display consistency")
	pflag.StringVarP(&config.Options.CurrencyFilter, "currency-filter", "f", config.Options.CurrencyFilter, "Currencies to be filtered in Transactions Display (comma separated list)")
	pflag.StringVar(&config.Options.LogFile, "log", config.Options.LogFile, "Log file")
	pflag.BoolVar(&config.Options.Debug, "exact", config.Options.Debug, "Display exact amount (no rounding)")
	pflag.StringVarP(&config.Options.TxsDisplay, "txs-display", "t", config.Options.TxsDisplay, "Display Transactions By Category : Exchanges|Deposits|Withdrawals|CashIn|CashOut|etc")
	// Sources
	pflag.StringVar(&config.Options.TxsCategory, "txs-categ", config.Options.TxsCategory, "Transactions Categories CSV file")
	pflag.StringVar(&config.Tools.CoinAPI.Key, "coinapi-key", config.Tools.CoinAPI.Key, "CoinAPI Key (https://www.coinapi.io/pricing?apikey)")
	pflag.StringVar(&config.Tools.CoinLayer.Key, "coinlayer-key", config.Tools.CoinLayer.Key, "CoinLayer Key (https://coinlayer.com/product)")
	pflag.StringSliceVar(&config.Blockchains.BTC.CSV, "btc-addresses-csv", config.Blockchains.BTC.CSV, "Bitcoin Addresses CSV files")
	pflag.StringSliceVar(&config.Blockchains.BTC.Addresses, "btc-address", config.Blockchains.BTC.Addresses, "Bitcoin Address")
	pflag.BoolVar(&config.Options.Bcd, "bcd", config.Options.Bcd, "Detect Bitcoin Diamond Fork")
	pflag.BoolVar(&config.Options.Bch, "bch", config.Options.Bch, "Detect Bitcoin Cash Fork")
	pflag.BoolVar(&config.Options.Btg, "btg", config.Options.Btg, "Detect Bitcoin Gold Fork")
	pflag.BoolVar(&config.Options.Lbtc, "lbtc", config.Options.Lbtc, "Detect Lightning Bitcoin Fork")
	pflag.StringVar(&config.Blockchains.BTG.JSON, "btg-txs", config.Blockchains.BTG.JSON, "Bitcoin Gold Transactions JSON file")
	pflag.StringSliceVar(&config.Blockchains.ETH.CSV, "eth-addresses-csv", config.Blockchains.ETH.CSV, "Ethereum Addresses CSV file")
	pflag.StringSliceVar(&config.Blockchains.ETH.Addresses, "eth-address", config.Blockchains.ETH.Addresses, "Ethereum Address")
	pflag.StringVar(&config.Tools.EtherScan.Key, "etherscan-apikey", config.Tools.EtherScan.Key, "Etherscan API Key (https://etherscan.io/myapikey)")
	pflag.StringVar(&config.Exchanges.Binance.API.Key, "binance-api-key", config.Exchanges.Binance.API.Key, "Binance API key")
	pflag.StringVar(&config.Exchanges.Binance.API.Secret, "binance-api-secret", config.Exchanges.Binance.API.Secret, "Binance API secret")
	pflag.StringSliceVar(&config.Exchanges.Binance.CSV.All, "binance", config.Exchanges.Binance.CSV.All, "Binance CSV file")
	pflag.BoolVar(&config.Options.BinanceExtended, "binance-extended", config.Options.BinanceExtended, "Use Binance CSV file extended format")
	pflag.StringSliceVar(&config.Exchanges.Bitfinex.CSV.All, "bitfinex", config.Exchanges.Bitfinex.CSV.All, "Bitfinex CSV file")
	pflag.StringVar(&config.Exchanges.Bitstamp.API.Key, "bitstamp-api-key", config.Exchanges.Bitstamp.API.Key, "Bitstamp API key")
	pflag.StringVar(&config.Exchanges.Bitstamp.API.Secret, "bitstamp-api-secret", config.Exchanges.Bitstamp.API.Secret, "Bitstamp API secret")
	pflag.StringSliceVar(&config.Exchanges.Bitstamp.CSV.All, "bitstamp", config.Exchanges.Bitstamp.CSV.All, "Bitstamp CSV file")
	pflag.StringVar(&config.Exchanges.Bittrex.API.Key, "bittrex-api-key", config.Exchanges.Bittrex.API.Key, "Bittrex API key")
	pflag.StringVar(&config.Exchanges.Bittrex.API.Secret, "bittrex-api-secret", config.Exchanges.Bittrex.API.Secret, "Bittrex API secret")
	pflag.StringSliceVar(&config.Exchanges.Bittrex.CSV.All, "bittrex", config.Exchanges.Bittrex.CSV.All, "Bittrex CSV file")
	pflag.StringSliceVar(&config.Exchanges.HitBTC.CSV.Trades, "hitbtc-trades", config.Exchanges.HitBTC.CSV.Trades, "HitBTC Trades CSV file")
	pflag.StringSliceVar(&config.Exchanges.HitBTC.CSV.Transfers, "hitbtc-transactions", config.Exchanges.HitBTC.CSV.Transfers, "HitBTC Transfers CSV file")
	pflag.StringVar(&config.Exchanges.HitBTC.API.Key, "hitbtc-api-key", config.Exchanges.HitBTC.API.Key, "HitBTC API Key")
	pflag.StringVar(&config.Exchanges.HitBTC.API.Secret, "hitbtc-api-secret", config.Exchanges.HitBTC.API.Secret, "HitBTC API Secret")
	pflag.StringSliceVar(&config.Exchanges.Coinbase.CSV.All, "coinbase", config.Exchanges.Coinbase.CSV.All, "Coinbase CSV file")
	pflag.StringSliceVar(&config.Exchanges.CoinbasePro.CSV.Trades, "coinbase-pro-fills", config.Exchanges.CoinbasePro.CSV.Trades, "CoinbasePro Fills CSV file")
	pflag.StringSliceVar(&config.Exchanges.CoinbasePro.CSV.Transfers, "coinbase-pro-account", config.Exchanges.CoinbasePro.CSV.Transfers, "CoinbasePro Account CSV file")
	pflag.StringSliceVar(&config.Exchanges.CdcApp.CSV.All, "cdc-app-crypto", config.Exchanges.CdcApp.CSV.All, "Crypto.com App Crypto Wallet CSV file")
	pflag.StringVar(&config.Exchanges.CdcEx.API.Key, "cdc-ex-api-key", config.Exchanges.CdcEx.API.Key, "Crypto.com Exchange API Key")
	pflag.StringVar(&config.Exchanges.CdcEx.API.Secret, "cdc-ex-api-secret", config.Exchanges.CdcEx.API.Secret, "Crypto.com Exchange Secret Key")
	pflag.StringVar(&config.Exchanges.CdcEx.JSON, "cdc-ex-exportjs", config.Exchanges.CdcEx.JSON, "Crypto.com Exchange JSON file from json-exporter.js")
	pflag.StringSliceVar(&config.Exchanges.CdcEx.CSV.Transfers, "cdc-ex-transfer", config.Exchanges.CdcEx.CSV.Transfers, "Crypto.com Exchange Deposit/Withdrawal CSV file")
	// pflag.StringSliceVar(&config.Exchanges.CdcEx.CSV.Staking, "cdc-ex-stake", config.Exchanges.CdcEx.CSV.Staking, "Crypto.com Exchange Stake CSV file")
	// pflag.StringSliceVar(&config.Exchanges.CdcEx.CSV.Supercharger, "cdc-ex-supercharger", config.Exchanges.CdcEx.CSV.Supercharger, "Crypto.com Exchange Supercharger CSV file")
	pflag.StringVar(&config.Exchanges.Kraken.API.Key, "kraken-api-key", config.Exchanges.Kraken.API.Key, "Kraken API key")
	pflag.StringVar(&config.Exchanges.Kraken.API.Secret, "kraken-api-secret", config.Exchanges.Kraken.API.Secret, "Kraken API secret")
	pflag.StringSliceVar(&config.Exchanges.Kraken.CSV.All, "kraken", config.Exchanges.Kraken.CSV.All, "Kraken CSV file")
	pflag.StringSliceVar(&config.Wallets.LedgerLive.CSV.All, "ledgerlive", config.Wallets.LedgerLive.CSV.All, "LedgerLive CSV file")
	pflag.StringSliceVar(&config.Exchanges.LocalBitcoins.CSV.Trades, "lb-trade", config.Exchanges.LocalBitcoins.CSV.Trades, "Local Bitcoin Trade CSV file")
	pflag.StringSliceVar(&config.Exchanges.LocalBitcoins.CSV.Transfers, "lb-transfer", config.Exchanges.LocalBitcoins.CSV.Transfers, "Local Bitcoin Transfer CSV file")
	pflag.StringSliceVar(&config.Wallets.Monero.CSV.All, "monero", config.Wallets.Monero.CSV.All, "Monero CSV file")
	pflag.StringSliceVar(&config.Wallets.MyCelium.CSV.All, "mycelium", config.Wallets.MyCelium.CSV.All, "MyCelium CSV file")
	pflag.StringSliceVar(&config.Exchanges.Poloniex.CSV.Trades, "poloniex-trades", config.Exchanges.Poloniex.CSV.Trades, "Poloniex Trades CSV file")
	pflag.StringSliceVar(&config.Exchanges.Poloniex.CSV.Deposits, "poloniex-deposits", config.Exchanges.Poloniex.CSV.Deposits, "Poloniex Deposits CSV file")
	pflag.StringSliceVar(&config.Exchanges.Poloniex.CSV.Withdrawals, "poloniex-withdrawals", config.Exchanges.Poloniex.CSV.Withdrawals, "Poloniex Withdrawals CSV file")
	pflag.StringSliceVar(&config.Exchanges.Poloniex.CSV.Distributions, "poloniex-distributions", config.Exchanges.Poloniex.CSV.Distributions, "Poloniex Distributions CSV file")
	pflag.StringSliceVar(&config.Exchanges.Revolut.CSV.All, "revolut", config.Exchanges.Revolut.CSV.All, "Revolut CSV file")
	pflag.StringSliceVar(&config.Exchanges.Uphold.CSV.All, "uphold", config.Exchanges.Uphold.CSV.All, "Uphold CSV file")
	// Output
	pflag.BoolVar(&config.Options.Display2086, "2086-display", config.Options.Display2086, "Display Cerfa 2086")
	pflag.BoolVar(&config.Options.Export2086, "2086", config.Options.Export2086, "Export Cerfa 2086 to 2086.xlsx")
	pflag.BoolVar(&config.Options.Export3916, "3916", config.Options.Export3916, "Export Cerfa 3916 to 3916.xlsx")
	pflag.BoolVar(&config.Options.ExportStock, "stock", config.Options.ExportStock, "Export stock balances to stock.xlsx")
	pflag.Parse()
	return config, nil
}
