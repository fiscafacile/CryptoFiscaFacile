# CryptoFiscaFacile

Cet outil veut vous aider à déclarer vos crypto aux impôts !

Gardez en tête que la loi n'étant pas encore définie sur tous les points, cet outil peut différer de votre point de vue, c'est pour cela qu'il est en open-source : à vous de modifier (ou faire modifier) à vos besoin.

Tout comme le fait qu'il ne supporte pas toutes les plateformes existantes, mais un guide vous est fourni pour vous aider à développer votre propre module.

Enfin, tout pull request est le bienvenu, j'essayerai de les intégrer le plus vite possible.

Aussi, le code actuel est en constante évolution, il se peut donc que la documentation ci dessous ne soit pas précise, mais il vous fournira une bonne base pour utiliser cet outil.

## Installation

```bash
$ go get github.com/fiscafacile/CryptoFiscaFacile
$ cd $GOPATH/src/github.com/fiscafacile/CryptoFiscaFacile
$ go build
```

## Utilisation

### Principe de focntionnement

Cet outil a besoin de "Sources" pour établir une liste de transactions qui constituent votre protefeuille global.

Ces "Sources" peuvent être :

- des fichiers CSV (souvent exportés depuis une plateforme ou établis manuellement)

- des fichiers JSON (autre formalisme de données structuées)

- des API de plateforme

Toutes les APIs utilisés par ce outils sont mis en cache dans des fichiers JSON rangés dans le répertoire `Cache` créé à coté de l'exécutable. Vous pouvez donc vérifier/exporter/modifier ces informations pour rendre votre utilisation cohérente. Pensez aussi à supprimer/déplacer/renommer les fichiers de cache si vous voulez récupérer les dernières informations de la plateforme.

Chaque transaction est composé d'une `Date`, d'une `Note` (donnant des informations pour la comprendre), optionellement d'une liste de frais `Fee`, optionellement d'une liste de sources `From` et optionellement d'une liste de destinations `To`.

Les `Fee`, `From` et `To` sont des "Actifs" composé d'un `Code` et d'un montant `Amount`.

Tous les montant dans l'outils sont des chiffres décimaux avec précision arbitraire : aucun arrondis n'est fait dans les calculs, seuls à l'affichage à la fin pour plus de claretée.

Une fois toutes les transactions récupérées de toutes les "Sources" que vous avez fournies à l'outil, il va essayer de catégoriser ces TXs

#### Catégories de TXs relatifs à une "Source"

- "Dépots" `Deposits` : ce sont des TXs qui ont un ou plusieurs `To` mais n'ont pas de `From` et possiblement des `Fee`.

- "Retraits" `Withdrawals` : c'est l'inverse des "Dépots".

- "Frais" `Fees` : les TXs qui n'ont que des `Fee`.

- "Echanges" `Exchanges` : des TXs qui ont des `From` et des `To`, et possiblement des `Fee`.

##### Catégories manuelles et semi-automatiques

Vous pourrez fournir une "Source" particulière pour rediriger certaines TXs dans des catégories manuelles comme des "Dons" `Gifts` et autres `AirDrops`.

Vous pourrez aussi activer la détection de `Forks` sur certaines crypto.

##### Catégories spécifiques à certaines plateforme

Sur certaines plateformes comme Crypto.com il existe aussi des `Cashbacks`, `Earns` et autres `Rewards`. Certaines TXs sont directement catégorisées en `CashOut` comme les paiements en crypto.

##### Catégories spécifiques ETH

Pour les sources ETH, il y a d'autre catégories spécificques : `Burns`, `Claims`, `Selfs` et `Swaps`.

#### Catégories de TXs relatifs au portefeuille global

Une fois toutes les TXs rangées dans des catégories, l'outils va essayer de rapprocher des TXs de différentes "Sources" pour synthétiser et recatégoriser au niveau du portefeuille global :

- "Transferts" `Transfers` : par fusion d'un `Deposits` avec un `Withdrawals` si les `Date` et `Amount` correspondent.

- `CashIn` et `CashOut` : ce sont respectivement des `Deposits` et `Withdrawals` ou des `Exchanges` donc l'"Actif" source ou destination sont des Fiats.

- les `Cashbacks` sont transformés en `CashIn` si aucun "reversal" n'est venu les annuler.

## Options

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

### Options de base

```bash
  -date string
        Date Filter (default "2021-01-01T00:00:00")
```
Permet d'afficher votre protefeuille global valorisé en Fiat à une date donnée.
Utile pour vérifier l'état du stock et estimer s'il manque des sources.

```bash
  -location string
        Date Filter Location (default "Europe/Paris")
```
Permet de choisir le fuseau horaire pour calculer les dates. Si vous voulez déclarer aux impôts français, il faut laisser "Europe/Paris".

```bash
  -stats
        Display accounts stats
```
Permet d'afficher le nombre de transactions par catégorie (toutes crypto confondues).

```bash
  -native string
        Native Currency for consolidation (default "EUR")
```
Choix de la Fiat pour consolidataion. Si vous voulez déclarer aux impôts français, il faut laisser "EUR".

### Options de "Sources"


## Donation

Si vous voulez faire un don a l'outil (pas à moi), cela permettra d'acheter un nom de domaine et payer un hébergement par exemple :

[![Donate with Bitcoin](https://en.cryptobadges.io/badge/small/36BTpmPbZaG2e5DyMpjEfDeEaiwjR8jGUM)](https://en.cryptobadges.io/donate/36BTpmPbZaG2e5DyMpjEfDeEaiwjR8jGUM)

[![Donate with Ethereum](https://en.cryptobadges.io/badge/small/0x9302F624d2C35fe880BFce22A36917b5dB5FAFeD)](https://en.cryptobadges.io/donate/0x9302F624d2C35fe880BFce22A36917b5dB5FAFeD)

## Remerciements

Merci au channel [Fiscalité crypto FR](t.me/fiscalitecryptofr) sur Telegram qui est une mine d'or d'information pour essayer de comprendre comment cela fonctionne.

## Copyright & License

Copyright (c) 2021-present FiscaFacile.

Released under the terms of the MIT license. See LICENSE for details.
