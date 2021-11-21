# CryptoFiscaFacile

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)
[![Open Source? Yes!](https://badgen.net/badge/Open%20Source%20%3F/Yes%21/blue?icon=github)](https://github.com/fiscafacile/CryptoFiscaFacile/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-blue.svg "PRs Welcome")](https://github.com/fiscafacile/CryptoFiscaFacile/pulls)
[![Telegram CryptoFiscaFacile](https://img.shields.io/badge/Telegram-CryptoFiscaFacile%20(support%20utilisation)-32AFED?logo=telegram)](https://telegram.me/cryptofiscafacile)
[![Discord CryptoFiscaFacile](https://img.shields.io/badge/Discord-CryptoFiscaFacile%20(développent%20v2)-5865f2?logo=discord)](https://discord.gg/PndJKReU9E)

Cet outil veut vous aider à déclarer vos cryptos aux impôts !

Attention néanmoins : les développeurs de CryptoFiscaFacile ne peuvent pas être tenus pour responsables des éventuelles erreurs ou imprécisions qui pourraient survenir dans vos déclarations fiscales suite à d'éventuels bugs de l'outil. En cas d'erreurs, l'entière responsabilité de ces erreurs vous incombe.

Gardez en tête que la loi n'étant pas encore définie sur tous les points, cet outil peut différer de votre point de vue, c'est pour cela qu'il est en open-source : à vous de modifier (ou faire modifier) à vos besoins.

Gardez aussi en tête le fait qu'il ne supporte pas toutes les plateformes existantes, mais un guide vous est fourni pour vous aider à développer votre propre module.

Tout pull request est le bienvenu, j'essayerai de les intégrer le plus vite possible.

Enfin, le code actuel est en constante évolution, il se peut donc que la documentation ci dessous ne soit pas précise, mais elle vous fournira une bonne base pour utiliser cet outil.

## Installation / Mise à jour

Vous aurez besoin de Go dont voici la [doc officelle d'installation](https://golang.org/doc/install).

Une fois Go installé sur votre système, ouvrez un terminal (cmd.exe/PowerShell sous Windows, Terminal sous MacOS, un shell sous Linux) et tappez cette commande :

```bash
$ go install github.com/fiscafacile/CryptoFiscaFacile
```

Le binaire de l'outil sera généré sur votre PC, vous pourrez le lancer en ligne de commande (donc dans un terminal) avec les [Options](#configuration) nécessaires à vos besoins.

```bash
$ CryptoFiscaFacile --help
```

Pour mettre à jour, il suffit de relancer la commande (pas trop sur de cela) :

```bash
$ go get -u github.com/fiscafacile/CryptoFiscaFacile
```

## Utilisation

### Principe de fonctionnement

Cet outil a besoin de "Sources" pour établir une liste de transactions qui constituent votre protefeuille global.

Ces "Sources" peuvent être :

- des fichiers CSV (souvent exportés depuis une plateforme ou établis manuellement)

- des fichiers JSON (autre formalisme de données structurées)

- des API de plateforme

Toutes les APIs utilisées par cet outil sont mis en cache dans des fichiers JSON rangés dans le répertoire `Cache` créé à côté de l'exécutable. Vous pouvez donc vérifier/exporter/modifier ces informations pour rendre votre utilisation cohérente. Pensez aussi à supprimer/déplacer/renommer les fichiers de cache si vous voulez récupérer les dernières informations de la plateforme.

Chaque transaction est composée d'une `Date`, d'une `Note` (donnant des informations pour la comprendre), optionellement d'une liste de frais `Fee`, optionellement d'une liste de sources `From` et optionellement d'une liste de destinations `To`.

Les `Fee`, `From` et `To` sont des "Actifs" composés d'un `Code` et d'un montant `Amount`.

Tous les montants dans l'outil sont des chiffres décimaux avec précision arbitraire : aucun arrondi n'est fait dans les calculs, seulement à l'affichage à la fin pour plus de clareté.

Une fois toutes les transactions récupérées de toutes les "Sources" que vous avez fournies à l'outil, il va essayer de catégoriser ces TXs.

#### Catégories de TXs relatives à une "Source"

- "Dépôts" `Deposits` : ce sont des TXs qui ont un ou plusieurs `To` mais n'ont pas de `From` et possiblement des `Fee`.

- "Retraits" `Withdrawals` : c'est l'inverse des "Dépôts".

- "Frais" `Fees` : les TXs qui n'ont que des `Fee`.

- "Echanges" `Exchanges` : des TXs qui ont des `From` et des `To`, et possiblement des `Fee`.

##### Catégories manuelles et semi-automatiques

Vous pourrez fournir une [Source](#catégorisation-manuelle-) particulière pour rediriger certaines TXs dans des catégories manuelles comme des "Dons" `Gifts` et autres `AirDrops`.

Vous pourrez aussi activer la détection de `Forks` sur certaines cryptos.

##### Catégories spécifiques à certaines plateformes

Sur certaines plateformes comme Crypto.com il existe aussi des `CommercialRebates` (cashback de carte, remboursement Netflix, Pay Checkout Reward et Gift Card Reward), `Interests` (intérêts du programme Earn, intérêts de "Stacking") et autres `Referrals`. Certaines TXs sont directement catégorisées en `CashOut` comme les paiements en crypto.

##### Catégories spécifiques ETH

Pour les sources ETH, il y a d'autres catégories spécifiques : `Burns`, `Claims`, `Selfs` et `Swaps`.

#### Catégories de TXs relatives au portefeuille global

Une fois toutes les TXs rangées dans des catégories, l'outil va essayer de rapprocher des TXs de différentes "Sources" pour synthétiser et recatégoriser au niveau du portefeuille global :

- "Transferts" `Transfers` : par fusion d'un `Deposits` avec un `Withdrawals` si les `Date` et `Amount` correspondent.

- `CashIn` et `CashOut` : ce sont respectivement des `Deposits` et `Withdrawals` ou des `Exchanges` dont l'"Actif" source ou destination sont des Fiats.

- les `Interests` sont transformés en `CashIn` et leurs montant global est affiché.

- les `CommercialRebates` sont transformés en `CashIn` si aucun "reversal" n'est venu les annuler et leurs montant global est affiché.

- les `Referrals` sont transformés en `CashIn` et leurs montant global est affiché.

## Configuration

Vous pouvez utiliser le fichier de configuration `config.exemple.yaml`, copiez le en `config.yaml` puis mofidiez le a votre guise. Il sera utilisé pour vos options par défaut (c'est à dire que si vous spécifiez une autre valeur d'une option dans la ligne de commande, elle sera prioritaire sur les valeurs du fichier de configuration).

### Options de base

#### Help

```
  --help
        Display all available arguments
```
Permet d'afficher toutes les options possibles.

#### Native Currency

```
  --native
        Native Currency for consolidation (default "EUR")
```
Choix de la Fiat pour consolidation. Si vous voulez déclarer aux impôts français, il faut laisser "EUR".

#### Location

```
  --location
        Date Filter Location (default "Europe/Paris")
```
Permet de choisir le fuseau horaire pour calculer les dates. Si vous voulez déclarer aux impôts français, il faut laisser "Europe/Paris".

#### Date

```
  --date
        Date Filter (default "2021-01-01T00:00:00")
```
Permet d'afficher votre protefeuille global valorisé en Fiat à une date donnée.
Utile pour vérifier l'état du stock et estimer s'il manque des sources.

### Options d'aide à l'établissement d'un portefeuille global cohérent

#### Stats

```
  --stats
        Display accounts stats
```
Permet d'afficher le nombre de transactions par catégorie (toutes cryptos confondues).

#### Check

```
  --check
        Check and Display consistency
```
Lance des vérifications d'intégrité sur les TXs du portefeuille globale et affiche les TXs KO. Les vérifications sont :

- tous les `Withdrawals` postérieurs au 1 Janvier 2019 doivent être justifiés, donc catégorisés ailleurs (`CashOut`, `Gifts`,...).

- tous les `Transfers` doivent avoir une balance nulle (la balance est la somme des `To` moins la somme des `From` moins la somme des `Fee`). Note pour pouvoir aditioner ces montant, ils faut qu'ils soient dans la même devise, ce qui est le cas pour les `Transfers` (normalement).

- toutes les TXs doivent avoir des montants positifs. Les montants de `From` et de `Fee` seront consédérés négativement par l'outil mais ils doivent être enregistré positivement dans leur TX par la "Source" qui les a produites.

#### Display

```
  --txs-display
        Display Transactions By Catergory : Exchanges|Deposits|Withdrawals|CashIn|CashOut|etc
  --currency-filter
        Currencies to be filtered in Transactions Display (comma separated list)
```
Affiche toutes les TXs d'une Catégorie (attention ceci peut être très long...).

Vous pouvez afficher toutes les Catégories avec `--txs-display Alls`.

Vous pouvez aussi afficher que les TXs concernant certaines cryptos, par exemple pour n'afficher que le BTC et le BCH : `--currency-filter BTC,BCH`.

### Options de "Sources"

Pour chaque Source, je vous indique le taux de support fourni par l'outil (l'exactitude de l'analyse pour cette Source). Si ce taux de support n'est pas bon, c'est sûrement parce que je n'ai pas assez d'exemples de transactions pour bien les analyser. Vous pouvez ouvrir un Ticket Github pour ajouter votre cas qui ne fontionne pas, j'essayerai de faire évoluer l'outil pour le rendre compatible.

#### Catégorisation Manuelle [![Support manuel](https://img.shields.io/badge/support-manuel-red)](#catégorisation-manuelle-)

```
  --txs-categ
        Transactions Categories CSV file
```
Il faut fournir un CSV à faire manuellement contenant toutes les transactions que vous voulez catégoriser manuellement (attention les champs dans le CSV doivent être séparés par des virgules, pas des points virgules comme le fait Excel en Français, le plus simple est de le faire dans un editeur de texte simple comme Notepad). Un CSV d'exemple est disponible, essayez `--txs-categ Inputs/TXS_Categ_exemple.csv --btc-address Inputs/BTC_Addresses_exemple.csv`.

Ce CSV identifie une TX par son `TxID` (identifiant dans la blockchain BTC, ETH, ou autre) et donne un `Type`. Les différents `Type` supportés sont :

- IN : va transformer la TX en `CashIn` même si ses `From` ne sont pas en Fiat. Utile pour simuler des plateformes qui ne proposent pas de CSV (comme DigyCode).

- OUT : va transformer la TX en `CashOut` même si ses `To` ne sont pas en Fiat. Utile pour les achats de bien ou service en crypto. Cela va transformer le `To` avec les infos `Value` et `Currency` de ce CSV.

- GIFT : va catégoriser la TX en don `Gifts`. Utile si vous offrez des cryptos à un ami pour lui montrer comment cela fonctionne lors de son anniversaire.

- INT : va catégoriser la TX en intérêt `Interests`.

- AIR : va catégoriser la TX en `AirDrops`.

- FEE : va associer toutes les TXs dont les Hash sont concaténées entre eux avec un point virgule ";" et fournis dans `Description` à la TX dont le Hash est donné dans `TXID`. Utile pour faire le ménage dans la catégorie `Fees`.

- TRANS : va associer la TX dont l'ID est fournis dans `Description` à la TX dont l'ID est donné dans `TXID`. Utile pour force l'association d'un dépôt avec un retrait même si l'un des deux a des frais non dissocié du montant. Dans le cas d'un forcage de `Deposits` et `Withdrawals` en `Transfers` avec TRANS, pour respecter la balance nulle hors frais, la différence de montant entre les deux TX initiales sera déduite du montant le plus grand et ajouté en tant que `Fee` dans la TX de `Transfers` résultante.

- SHIT : va ignorer la TX donc aucune catégorisation. Utile si vous avez des Shitcoins dont vous ne voulez pas.

- CUS : va retrancher une partie du montant de `From` ou `To` comme si vous en aviez la gestion mais qu'ils ne vous appartenaient pas (Custody), ils ne seront donc pas consiédérés dans votre portefeuille global. Utile si vous avez acheté des cryptos pour votre grand-père, mais attention, il devra lui aussi les déclarer.

Les colones du CSV doivent être : `TxID,Type,Description,Value,Currency`

#### Binance [![Support léger](https://img.shields.io/badge/support-bon-blue)](#binance-)

Par API :
```
  --binance-api-key
        Binance API key
  --binance-api-secret
        Binance API secret
```

Par CSV :
```
  --binance
        Binance CSV file
  --binance-extended
        Use Binance CSV file extended format
```
Il faut fournir le fichier CSV récupéré dans Binance (https://www.binance.com/fr/my/wallet/history puis "Générer un relevé complet").
Vous pouvez modifier ce fichier CSV pour ajouter une colone `Fee` entre `Change` et `Remark`, et donc reseigner la part de frais dans les `Withdraw` qui ont un `Remark` avec `Withdraw fee is included`, cela permet de bien fusioner ce `Withdrawals` avec un autre `Deposits` pour en faire un `Transfers` lors de l'analyse des TXs. Dans ce cas, n'oubliez pas de rajouter l'option `--binance-extended`. Ces frais seront automatiquement déduits du montant du retrait, veuillez donc ne pas toucher à la valeur `Change`.

Les colones du CSV d'origine doivent être : `UTC_Time,Account,Operation,Coin,Change,Remark`
Les colones du CSV étendu doivent être : `UTC_Time,Account,Operation,Coin,Change,Fee,Remark`

#### Bitfinex [![Support bon](https://img.shields.io/badge/support-bon-blue)](#bitfinex-)

```
  --bitfinex
        Bitfinex CSV file
```
Il faut fournir le fichier CSV récupéré dans Bitfinex (https://report.bitfinex.com/ledgers puis choisissez les dates et "Export", choisissez Date Format : DD-MM-YY).

Les colones du CSV d'origine doivent être : `#,DESCRIPTION,CURRENCY,AMOUNT,BALANCE,DATE,WALLET`

#### Bitstamp [![Support bon](https://img.shields.io/badge/support-bon-blue)](#bitstamp-)

```
  --bitstamp
        Bitstamp CSV file
  --bitstamp-api-key
        Bitstamp API key
  --bitstamp-api-secret
        Bitstamp API secret
```
Il faut fournir les fichiers CSV récupérés dans Bitstamp.

Les colones du CSV d'origine doivent être : `Type,Datetime,Account,Amount,Value,Rate,Fee,Sub Type`

#### Bittrex [![Support bon](https://img.shields.io/badge/support-bon-blue)](#bittrex-)

```
  --bittrex
        Bittrex CSV file
  --bittrex-api-key
        Bittrex API key
  --bittrex-api-secret
        Bittrex API secret
```
Il faut fournir les fichiers CSV récupérés dans Bittrex (https://global.bittrex.com/history puis "Download Order History").

Les colones du CSV d'origine doivent être : `Uuid,Exchange,TimeStamp,OrderType,Limit,Quantity,QuantityRemaining,Commission,Price,PricePerUnit,IsConditional,Condition,ConditionTarget,ImmediateOrCancel,Closed,TimeInForceTypeId,TimeInForce`

Il est nécessaire de fournir l'API et le CSV car chaque support a son défaut :
- l'API ne retourne pas les transactions liées à des assets délistés.
- le CSV ne comprend pas l'historique de dépot/retrait.

#### BTC [![Support avancé](https://img.shields.io/badge/support-avanc%C3%A9-green)](#btc-)

```
  --btc-address
        Bitcoin Address
  --btc-addresses-csv
        Bitcoin Addresses CSV file
  --bcd
        Detect Bitcoin Diamond Fork
  --bch
        Detect Bitcoin Cash Fork
  --btg
        Detect Bitcoin Gold Fork
  --lbtc
        Detect Lightning Bitcoin Fork
```
Il faut fournir :

- soit une ou plusieurs adresses directement dans la ligne de commande : `--btc-address 36BTpmPbZaG2e5DyMpjEfDeEaiwjR8jGUM --btc-address bc1qlmsx8vtk03jwcuafe7vzvddjzg4nsfvflgs4k9`

- soit un CSV à faire manuellement contenant toutes les addresses BTC que vous possédez (attention les champs dans le CSV doivent être séparés par des virgules, pas des points virgules comme le fait Excel en Français, le plus simple est de le faire dans un editeur de texte simple comme Notepad). Un CSV d'exemple est disponible, essayez `--btc-addresses-csv Inputs/BTC_Addresses_exemple.csv`.

L'outil se chargera de récupérer la liste des transactions associées sur Blockstream (pas besoin de API Key).

Vous pouvez aussi demander la detection d'un des Fork de BTC, l'outil vous dira dans quel wallet vous avez un montant dû au Fork et intègrera ces montants à votre portefeuille global.

Les colones du CSV doivent être : `Address,Description`

#### BTG [![Support manuel](https://img.shields.io/badge/support-manuel-red)](#btg-)

```
  --btg-txs
        Bitcoin Gold Transactions JSON file
```
Expériemental.

#### Crypto.com [![Support avancé](https://img.shields.io/badge/support-avanc%C3%A9-green)](#crypto.com-)

- App avec CSV:
```
  --cdc-app-crypto
        Crypto.com App Crypto Wallet CSV file
```
Il faut fournir les CSV récupérés dans l'App (celui des Transactions du Portefeuille Crypto).

Les colones du CSV du portefeuille Crypto de l'APP doivent être : `Timestamp (UTC),Transaction Description,Currency,Amount,To Currency,To Amount,Native Currency,Native Amount,Native Amount (in USD),Transaction Kind`

- Exchange avec JS et JSON:
```
  --cdc-ex-exportjs
        Crypto.com Exchange JSON file from json_exporter.js
```
Il faut fournir le JSON récupéré dans l'Exchange Crypto.com avec la méthode décrite [ici](https://github.com/fiscafacile/CryptoFiscaFacile/wiki/Crypto.com-Exchange-JSON-method)

Cette méthode vous permet de récupérer les `Deposits` et `Withdrawals`, les `Interests` du Staking de CRO et Soft Staking, les `CommercialRebates` du bonus d'inscription et des Syndicates, les `Referrals` du programme de parainage, les `Minings` des Superchargers.

- Exchange avec CSV:
```
  --cdc-ex-spot-trade
        Crypto.com Exchange Spot Trade CSV file
  --cdc-ex-transfer
        Crypto.com Exchange Deposit/Withdrawal CSV file
```
Il faut fournir les CSV récupérés dans l'Exchange Crypto.com.

Préférez la methode JS+JSON ci dessus, elle est plus complète.

Les colones du CSV de l'Exchange Spot Trade doivent être : `account_type,order_id,trade_id,create_time_utc,symbol,side,liquditiy_indicator,traded_price,traded_quantity,fee,fee_currency`

Les colones du CSV de l'Exchange Transfer doivent être : `create_time_utc,currency,amount,fee,address,status`

- Exchange avec API:
```
  --cdc-ex-api-key
        Crypto.com Exchange API Key
  --cdc-ex-secret-key
        Crypto.com Exchange Secret Key
```
Il faut donner le api-key et secret-key que vous pouvez créer dans votre compte.

Il faut activer le droit de "Withdrawal" (si disponible pour vous) si vous voulez récupérer les `Withdrawals` et `Deposits` (je ne l'ai pas donc je n'ai pas pu tester). Dans le cas contraire, le CSV Transfers  ou le JSON permet de les mettre dans l'outil sans l'API.

Par contre les `Exchanges` sur le Spot Market seront bien récupérés sans droit particulier (attention tout de même c'est assez long, on ne peut faire qu'une requête par seconde pour récupérer les Trades d'une seule journée, il faut donc de nombreuses requêtes pour remonter au jour du lancement de l'Exchange le 14 Nov 2019).

#### Coinbase [![Support bon](https://img.shields.io/badge/support-bon-blue)](#coinbase-)

```
  --coinbase
        Coinbase CSV file
```
Il faut fournir le CSV récupéré sur Coinbase.

Le CSV contient une entête qui sera ignorée par l'outil.

Pour les "Transaction Type" "Send" du CSV, les frais ne sont pas renseignés, l'outil ne pourra donc pas agréger ce `Withdrawals` avec le `Deposits` d'une autre Source. Vous pouvez l'y aider en retrouvant le `Depostis` correspondant à la main et calculant les frais (la différence entre les deux montants) puis en le rajoutant dans la colone `EUR Fees` de ce CSV.

Les colones du CSV doivent être : `Timestamp,Transaction Type,Asset,Quantity Transacted,EUR Spot Price at Transaction,EUR Subtotal,EUR Total (inclusive of fees),EUR Fees,Notes`

#### Coinbase Pro [![Support bon](https://img.shields.io/badge/support-bon-blue)](#coinbase-pro-)

```
  --coinbase-pro-account
        Coinbase Pro Account CSV file
  --coinbase-pro-fills
        Coinbase Pro Fills CSV file
```
Il faut fournir les CSV récupérés sur Coinbase.

Les colones du CSV Account doivent être : `portfolio,type,time,amount,balance,amount/balance unit,transfer id,trade id,order id`

Les colones du CSV Fills doivent être : `portfolio,trade id,product,side,created at,size,size unit,price,fee,total,price/fee/total unit`

#### ETH [![Support avancé](https://img.shields.io/badge/support-avanc%C3%A9-green)](#eth-)

```
  --eth-address
        Ethereum Address
  --eth-addresses-csv
        Ethereum Addresses CSV file
```
Il faut fournir :

- soit une ou plusieurs adresses directement dans la ligne de commande : `--eth-address 0x9302F624d2C35fe880BFce22A36917b5dB5FAFeD --eth-address 0x9302F624d2C35fe880BFce22A36917b5dB5FAFeD`

- soit un CSV à faire manuellement contenant toutes les addresses ETH que vous possédez (attention les champs dans le CSV doivent être séparés par des virgules, pas des points virgules comme le fait Excel en Français, le plus simple est de le faire dans un editeur de texte simple comme Notepad). Un CSV d'exemple est disponible, essayez `--eth-addresses-csv Inputs/ETH_Addresses_exemple.csv`.

L'outil se chargera de récupérer la liste des transactions associées sur [Etherscan.io](#etherscan.io) (à une vitesse limitée de 5 requêtes par secondes si vous ne fournissez pas une API Key).

Il détectera aussi les Token ERC20 et ERC721 (NFT) associés.

Les colones du CSV doivent être : `Address,Description`

#### HitBTC [![Support bon](https://img.shields.io/badge/support-bon-blue)](#hitbtc-)

- Via API
```
  --hitbtc-api-key
        HitBTC API Key
  --hitbtc-secret-key
        HitBTC Secret Key
```
L'API est utilisée pour récupérer les transactions du compte `Deposits` et `Withdrawals` afin de les ajouter à votre portefeuille global.

La clé API doit avoir les droits :

- `Payment information` pour les `Deposits` et `Withdrawals` (équivalent du CSV Transactions).

- `Orderbook, History, Trading balance` pour les `Exchanges` (équivalent du CSV Trades).

- Via CSV
```
  --hitbtc-trades
        HitBTC Trades CSV file
  --hitbtc-transactions
        HitBTC Transactions CSV file
```
Il faut fournir les fichiers CSV récupérés dans HitBTC (https://hitbtc.com/reports).

Le CSV Transactions ne fournit pas les infos de `Fee`, il vaut donc mieux utiliser l'API.

Les colones du CSV Trades doivent être : `Email,Date (UTC),Instrument,Trade ID,Order ID,Side,Quantity,Price,Volume,Fee,Rebate,Total,Taker`

Les colones du CSV Transactions doivent être : `Email,Date (UTC),Operation id,Type,Amount,Transaction hash,Main account balance,Currency`

#### Kraken [![Support bon](https://img.shields.io/badge/support-bon-blue)](#kraken-)

- Via API
```
  --kraken-api-key
        Kraken API key
  --kraken-api-secret
        Kraken API secret
```
L'API est utilisée pour récupérer l'ensemble des transactions du ledger afin de les ajouter à votre portefeuille global.

- Via CSV
```
  --kraken
        Kraken CSV file
```
Il faut fournir le fichier CSV récupéré dans Kraken (https://www.kraken.com/u/history/export puis sélectionner "Ledgers" et "All fields").

Les colones du CSV d'origine doivent être : `txid,refid,time,type,subtype,aclass,asset,amount,fee,balance`

#### Local Bitcoin [![Support bon](https://img.shields.io/badge/support-bon-blue)](#local-bitcoin-)

```
  --lb-trade
        Local Bitcoin Trade CSV file
  --lb-transfer
        Local Bitcoin Transfer CSV file
```

Les colones du CSV de Trade doivent être : `id,created_at,buyer,seller,trade_type,btc_amount,btc_traded,fee_btc,btc_amount_less_fee,btc_final,fiat_amount,fiat_fee,fiat_per_btc,currency,exchange_rate,transaction_released_at,online_provider,reference`

Les colones du CSV de Transfer doivent être : `TXID, Created, Received, Sent, TXtype, TXdesc, TXNotes`

#### Ledger Live [![Support bon](https://img.shields.io/badge/support-bon-blue)](#ledger-live-)

```
  --ledgerlive
        LedgerLive CSV file
```

Les colones du CSV doivent être : `Operation Date,Currency Ticker,Operation Type,Operation Amount,Operation Fees,Operation Hash,Account Name,Account xpub`

#### Monero Wallet [![Support bon](https://img.shields.io/badge/support-bon-blue)](#monero-wallet-)

```
  --monero
        Monero Wallet CSV file
```

Les colones du CSV doivent être : `blockHeight,epoch,date,direction,amount,atomicAmount,fee,txid,label,subaddrAccount,paymentId`

#### MyCelium [![Support déprécié](https://img.shields.io/badge/support-d%C3%A9pr%C3%A9ci%C3%A9-red)](#mycelium-)

Vous devriez exporter les clés publiques de votre wallet et utiliser la "Source" [BTC](#btc-).

```
  --mycelium
        MyCelium CSV file
```

Les colones du CSV doivent être : `Account,Transaction ID,Destination Address,Timestamp,Value,Currency,Transaction Label`

#### Poloniex [![Support bon](https://img.shields.io/badge/support-bon-blue)](#poloniex-)

```
  --poloniex-trades
        Poloniex Trades CSV file
  --poloniex-deposits
        Poloniex Deposits CSV file
  --poloniex-withdrawals
        Poloniex Withdrawals CSV file
  --poloniex-distributions
        Poloniex Distributions CSV file
```

Les colones du CSV de Trades doivent être : `Date,Market,Category,Type,Price,Amount,Total,Fee,Order Number,Base Total Less Fee,Quote Total Less Fee,Fee Currency,Fee Total`

Les colones du CSV de Deposits doivent être : `Date,Currency,Amount,Address,Status`

Les colones du CSV de Withdrawals doivent être : `Date,Currency,Amount,Fee Deducted,Amount - Fee,Address,Status`

Les colones du CSV de Distributions doivent être : `date,currency,amount,wallet`

#### Revolut [![Support bon](https://img.shields.io/badge/support-bon-blue)](#revolut-)

```
  --revolut
        Revolut CSV file
```

Les colones du CSV doivent être : `Completed Date,Description,Paid Out (BTC),Paid In (BTC),Exchange Out, Exchange In, Balance (BTC), Category, Notes`

#### Uphold [![Support bon](https://img.shields.io/badge/support-bon-blue)](#uphold-)

```
  --uphold
        Uphold CSV file
```

Les colones du CSV doivent être : `Date,Destination,Destination Amount,Destination Currency,Fee Amount,Fee Currency,Id,Origin,Origin Amount,Origin Currency,Status,Type`

### Options de "Providers"

Cet outil utilise plusieurs APIs de plateformes pour récupérer soit des taux de changes (CoinGecko, CoinLayer et CoinAPI), soit des transactions sur une blockchain particulière (Blockstream pour BTC et Etherscan pour ETH). Certaines de ces APIs ont besoins d'une clé.

#### CoinAPI.io

```
  --coinapi-key
        CoinAPI Key (https://www.coinapi.io/pricing?apikey)
```

#### CoinLayer.com

```
  --coinlayer-key
        CoinLayer Key (https://coinlayer.com/product)
```

#### Etherscan.io

```
  --etherscan-apikey
        Etherscan API Key (https://etherscan.io/myapikey)
```
Utilisé pour la Source [ETH](#eth-), si vous ne la fournissez pas les requêtes seront limitées à 5 par secondes.

### Options de sortie

```
  --2086-display
        Display Cerfa 2086
  --2086
        Export Cerfa 2086 in 2086.xlsx
  --cashin-bnc-2019
        Convert AirDrops/CommercialRebates/Interests/Minings/Referrals into CashIn for 2019's Txs in 2086
  --cashin-bnc-2020
        Convert AirDrops/CommercialRebates/Interests/Minings/Referrals into CashIn for 2020's Txs in 2086
```
Cela vous génère automatiquement le formulaire 2086 !

Il y a deux façons de considérer les AirDrops/CommercialRebates/Interests/Minings/Referrals :

- soit ils sont ajoutés simplement au portefeuille avec une valeur de 0€

- soit ils sont ajoutés au portefeuille avec leur valeur du jour en EUR (donc convertis en CashIn), ce qui va accroitre votre Prix Total d'Acquisition et faire baisser votre Plus-Value, mais en contrepartie, il convient de les déclarer en revenus non commerciaux non professionnels (régime déclaratif spécial ou micro BNC) dans la case 5KU/5LU/5MU de votre 2042-C-PRO.

Pour activer la seconde méthode d'intégration, vous pouvez utiliser les options `--cashin-bnc-2019` et/ou `--cashin-bnc-2020`.

En attendant que la loi soit plus claire à ce sujet, nous vous laissons le choix. Vous pouvez venir demander de l'aide à ce sujet sur le groupe [![Fiscalité crypto FR](https://img.shields.io/badge/Telegram-Fiscalité%20crypto%20FR-32AFED?style=for-the-badge&logo=telegram)](https://telegram.me/fiscalitecryptofr).

```
  --3916
        Export Cerfa 3916 in 3916.xlsx
```
Cela vous génère automatiquement le formulaire 3916 !

```
  --stock
        Export stock balances in stock.xlsx
```
Cela vous génère automatiquement une fiche de stock de tous vos coins !

## Donation

Si vous voulez faire un don à l'outil (pas à moi), cela permettra d'acheter un nom de domaine et payer un hébergement par exemple :

![Donate with Bitcoin](https://img.shields.io/static/v1?label=BTC&message=bc1q760y8ynwuf6ckfndtuasnhatcv6ewfrk3wy8wm&color=f2a900&logo=bitcoin)

![Donate with Ethereum](https://img.shields.io/static/v1?label=ETH&message=0x5e69c2cf20f80a6622f83f779010d077ea3f3c52&color=c99d66&logo=ethereum)

![Donate with Crytpo.org Chain](https://img.shields.io/static/v1?label=CRO&message=cro1e0ch6yemdletz44k74wvf9pmg3yv0lkrs9dlrv&color=002d74&logo=data:image/svg%2bxml;base64,PHN2ZyByb2xlPSJpbWciIHZpZXdCb3g9IjAgMCA0NSA0NSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cGF0aCBkPSJNMTkuMyAwTDAgMTEuMlYzMy41TDE5LjMgNDQuN0wzOC42IDMzLjVWMTEuMkwxOS4zIDBaTTMyLjkgMzAuMkwxOS4zIDM4TDUuNyAzMC4yVjE0LjVMMTkuMyA2LjZMMzIuOSAxNC41VjMwLjJaIiBmaWxsPSIjMDAyZDc0Ii8+PHBhdGggZD0iTTI4LjQwMDggMjcuNDk4NEwxOS4zMDA4IDMyLjY5ODRMMTAuMzAwOCAyNy40OTg0VjE3LjA5ODRMMTkuMzAwOCAxMS44OTg0TDI4LjQwMDggMTcuMDk4NFYyNy40OTg0WiIgZmlsbD0iIzAwMmQ3NCIvPjwvc3ZnPg==)

![Donate with Litecoin](https://img.shields.io/static/v1?label=LTC&message=ltc1qyk50t590ztuzclut3tm2s9ktne4kal8qjdyk7u&color=lightgrey&logo=litecoin)

## Support

Si vous avez un problème d'utilisation ou pour le développement d'un module et que cette doc ne vous apporte pas de réponse, venez me la poser dans les groupes officiels de support sur :

[![Telegram CryptoFiscaFacile](https://img.shields.io/badge/Telegram-CryptoFiscaFacile-32AFED?style=for-the-badge&logo=telegram&logoColor=white)](https://telegram.me/cryptofiscafacile)


## Contribution

Pour les développeurs de la v2 :

[![Discord CryptoFiscaFacile](https://img.shields.io/badge/Discord-CryptoFiscaFacile-5865f2?style=for-the-badge&logo=discord&logoColor=white)](https://discord.gg/PndJKReU9E)

## Remerciements

Merci au groupe [![Fiscalité crypto FR](https://img.shields.io/badge/Telegram-Fiscalité%20crypto%20FR-32AFED?style=for-the-badge&logo=telegram)](https://telegram.me/fiscalitecryptofr) qui est une mine d'or d'informations pour essayer de comprendre comment cela fonctionne.

## Copyright & License

Copyright (c) 2021-present FiscaFacile.

Released under the terms of the MIT license. See LICENSE for details.
