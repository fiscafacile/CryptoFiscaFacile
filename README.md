# CryptoFiscaFacile

Cet outil veut vous aider à déclarer vos crypto aux impôts !

Gardez en tête que la loi n'étant pas encore définie sur tous les points, cet outil peut différer de votre point de vue, c'est pour cela qu'il est en open-source : à vous de modifier (ou faire modifier) à vos besoin.

Tout comme le fait qu'il ne supporte pas toutes les plateformes existantes, mais un guide vous est fourni pour vous aider à développer votre propre module.

Enfin, tout pull request est le bienvenu, j'essayerai de les intégrer le pus vite possible.

### Installation

```bash
$ go get github.com/fiscafacile/CryptoFiscaFacile
$ cd $GOPATH/src/github.com/fiscafacile/CryptoFiscaFacile
$ go build
```

### Utilisation

```bash
$ CryptoFiscaFacile -h
Usage of CryptoFiscaFacile:
  -2086
        Display Cerfa 2086
  -bcd
        Detect Bitcoin Diamond Fork
  -bch
        Detect Bitcoin Cash Fork
  -binance string
        Binance CSV file
  -binance_extended
        Use Binance CSV file extended format
  -bitfinex string
        Bitfinex CSV file
  -btc_address string
        Bitcoin Addresses CSV file
  -btc_categ string
        Bitcoin Categories CSV file
  -btc_exclude float
        Exclude Bitcoin Amount
  -btg
        Detect Bitcoin Gold Fork
  -btg_txs string
        Bitcoin Gold Transactions JSON file
  -cdc_app string
        Crypto.com App CSV file
  -cdc_ex_stake string
        Crypto.com Exchange Stake CSV file
  -cdc_ex_supercharger string
        Crypto.com Exchange Supercharger CSV file
  -cdc_ex_transfer string
        Crypto.com Exchange Deposit/Withdrawal CSV file
  -check
        Check and Display consistency
  -coinapi_key string
        CoinAPI Key (https://www.coinapi.io/pricing?apikey)
  -coinbase string
        Coinbase CSV file
  -coinlayer_key string
        CoinLayer Key (https://coinlayer.com/product)
  -curr_filter string
        Currencies to be filtered in Transactions Display (comma separated list)
  -date string
        Date Filter (default "2021-01-01T00:00:00")
  -eth_address string
        Ethereum Addresses CSV file
  -etherscan_apikey string
        Etherscan API Key (https://etherscan.io/myapikey)
  -lb_trade string
        Local Bitcoin Trade CSV file
  -lb_transfer string
        Local Bitcoin Transfer CSV file
  -ledgerlive string
        LedgerLive CSV file
  -location string
        Date Filter Location (default "Europe/Paris")
  -metamask string
        MetaMask CSV file
  -mycelium string
        MyCelium CSV file
  -native string
        Native Currency for consolidation (default "EUR")
  -revolut string
        Revolut CSV file
  -stats
        Display accounts stats
  -txscat string
        Display Transactions By Catergory : Exchanges|Deposits|Withdrawals|CashIn|CashOut|etc
```

### Options

#### Options de base

```bash
-date string
        Date Filter (default "2021-01-01T00:00:00")
```
Permet d'afficher votre protefeuille global valorisé en Fiat à une date donnée.
Utile pour vérifier l'état du stock et estimer s'il manque des sources.
