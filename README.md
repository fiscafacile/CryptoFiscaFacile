# CryptoFiscaFacile

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)
[![Open Source? Yes!](https://badgen.net/badge/Open%20Source%20%3F/Yes%21/blue?icon=github)](https://github.com/fiscafacile/CryptoFiscaFacile/)

Cet outil veut vous aider à déclarer vos cryptos aux impôts !

Gardez en tête que la loi n'étant pas encore définie sur tous les points, cet outil peut différer de votre point de vue, c'est pour cela qu'il est en open-source : à vous de modifier (ou faire modifier) à vos besoins.

Gardez aussi en tête le fait qu'il ne supporte pas toutes les plateformes existantes, mais un guide vous est fourni pour vous aider à développer votre propre module.

Tout pull request est le bienvenu, j'essayerai de les intégrer le plus vite possible.

Enfin, le code actuel est en constante évolution, il se peut donc que la documentation ci dessous ne soit pas précise, mais elle vous fournira une bonne base pour utiliser cet outil.

## Installation

```bash
$ go get github.com/fiscafacile/CryptoFiscaFacile
$ cd $GOPATH/src/github.com/fiscafacile/CryptoFiscaFacile
$ go build
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

Vous pourrez fournir une "Source" particulière pour rediriger certaines TXs dans des catégories manuelles comme des "Dons" `Gifts` et autres `AirDrops`.

Vous pourrez aussi activer la détection de `Forks` sur certaines cryptos.

##### Catégories spécifiques à certaines plateformes

Sur certaines plateformes comme Crypto.com il existe aussi des `Cashbacks`, `Earns` et autres `Rewards`. Certaines TXs sont directement catégorisées en `CashOut` comme les paiements en crypto.

##### Catégories spécifiques ETH

Pour les sources ETH, il y a d'autres catégories spécifiques : `Burns`, `Claims`, `Selfs` et `Swaps`.

#### Catégories de TXs relatives au portefeuille global

Une fois toutes les TXs rangées dans des catégories, l'outil va essayer de rapprocher des TXs de différentes "Sources" pour synthétiser et recatégoriser au niveau du portefeuille global :

- "Transferts" `Transfers` : par fusion d'un `Deposits` avec un `Withdrawals` si les `Date` et `Amount` correspondent.

- `CashIn` et `CashOut` : ce sont respectivement des `Deposits` et `Withdrawals` ou des `Exchanges` dont l'"Actif" source ou destination sont des Fiats.

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
Permet d'afficher le nombre de transactions par catégorie (toutes cryptos confondues).

```bash
  -native string
        Native Currency for consolidation (default "EUR")
```
Choix de la Fiat pour consolidation. Si vous voulez déclarer aux impôts français, il faut laisser "EUR".

### Options de "Sources"

Pour chaque Source, je vous indique le taux de support fourni par l'outil (l'exactitude de l'analyse pour cette Source). Si ce taux de support n'est pas bon, c'est sûrement parce que je n'ai pas assez d'exemples de transactions pour bien les analyser. Vous pouvez ouvrir un Ticket Github pour ajouter votre cas qui ne fontionne pas, j'essayerai de faire évoluer l'outil pour le rendre compatible.

#### Binance ![Support léger](https://img.shields.io/badge/support-l%C3%A9ger-yellow)

```bash
  -binance string
        Binance CSV file
  -binance_extended
        Use Binance CSV file extended format
```
Il faut fournir le fichier CSV récupéré dans Binance (https://www.binance.com/fr/my/wallet/history puis "Générer un relevé complet").
Vous pouvez modifier ce fichier CSV pour ajouter une colone `Fee` entre `Change` et `Remark`, et donc reseigner la part de frais dans les `Withdraw` qui ont un `Remark` avec `Withdraw fee is included`, cela permet de bien fusioner ce `Withdrawals` avec un autre `Deposits` pour en faire un `Transfers` lors de l'analyse des TXs. Dans ce cas, n'oubliez pas de rajouter l'option `-binance_extended`.

#### Bitfinex ![Support bon](https://img.shields.io/badge/support-bon-blue)

```bash
  -bitfinex string
        Bitfinex CSV file
```


## Donation

Si vous voulez faire un don à l'outil (pas à moi), cela permettra d'acheter un nom de domaine et payer un hébergement par exemple :

[![Donate with Bitcoin](https://en.cryptobadges.io/badge/small/36BTpmPbZaG2e5DyMpjEfDeEaiwjR8jGUM)](https://en.cryptobadges.io/donate/36BTpmPbZaG2e5DyMpjEfDeEaiwjR8jGUM)

[![Donate with Ethereum](https://en.cryptobadges.io/badge/small/0x9302F624d2C35fe880BFce22A36917b5dB5FAFeD)](https://en.cryptobadges.io/donate/0x9302F624d2C35fe880BFce22A36917b5dB5FAFeD)

## Remerciements

Merci au canal [![Fiscalité crypto FR](https://img.shields.io/badge/Join-Telegram-blue?style=for-the-badge&logo=data:image/svg%2bxml;base64,PHN2ZyBlbmFibGUtYmFja2dyb3VuZD0ibmV3IDAgMCAyNCAyNCIgaGVpZ2h0PSI1MTIiIHZpZXdCb3g9IjAgMCAyNCAyNCIgd2lkdGg9IjUxMiIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cGF0aCBkPSJtOS40MTcgMTUuMTgxLS4zOTcgNS41ODRjLjU2OCAwIC44MTQtLjI0NCAxLjEwOS0uNTM3bDIuNjYzLTIuNTQ1IDUuNTE4IDQuMDQxYzEuMDEyLjU2NCAxLjcyNS4yNjcgMS45OTgtLjkzMWwzLjYyMi0xNi45NzIuMDAxLS4wMDFjLjMyMS0xLjQ5Ni0uNTQxLTIuMDgxLTEuNTI3LTEuNzE0bC0yMS4yOSA4LjE1MWMtMS40NTMuNTY0LTEuNDMxIDEuMzc0LS4yNDcgMS43NDFsNS40NDMgMS42OTMgMTIuNjQzLTcuOTExYy41OTUtLjM5NCAxLjEzNi0uMTc2LjY5MS4yMTh6IiBmaWxsPSIjMDM5YmU1Ii8+PC9zdmc+)](https://telegram.me/fiscalitecryptofr) sur Telegram qui est une mine d'or d'informations pour essayer de comprendre comment cela fonctionne.

## Copyright & License

Copyright (c) 2021-present FiscaFacile.

Released under the terms of the MIT license. See LICENSE for details.
