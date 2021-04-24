package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fiscafacile/CryptoFiscaFacile/binance"
	"github.com/fiscafacile/CryptoFiscaFacile/bitfinex"
	"github.com/fiscafacile/CryptoFiscaFacile/blockstream"
	"github.com/fiscafacile/CryptoFiscaFacile/coinbase"
	"github.com/fiscafacile/CryptoFiscaFacile/cryptocom"
	"github.com/fiscafacile/CryptoFiscaFacile/etherscan"
	"github.com/fiscafacile/CryptoFiscaFacile/ledgerlive"
	"github.com/fiscafacile/CryptoFiscaFacile/localbitcoin"
	"github.com/fiscafacile/CryptoFiscaFacile/metamask"
	"github.com/fiscafacile/CryptoFiscaFacile/mycelium"
	"github.com/fiscafacile/CryptoFiscaFacile/revolut"
	"github.com/fiscafacile/CryptoFiscaFacile/wallet"
	"github.com/shopspring/decimal"
)

func main() {
	// Parse args
	pDate := flag.String("date", "2021-01-01T00:00:00", "Date Filter")
	pLocation := flag.String("location", "Europe/Paris", "Date Filter Location")
	pNative := flag.String("native", "EUR", "Native Currency for consolidation")
	pAccount := flag.String("acc", "", "Display account : Exchanges|Deposits|Withdrawals|CashIn|CashOut|etc")
	pStats := flag.Bool("stats", false, "Display accounts stats")
	pUnjustifWithdrawals := flag.Bool("unjustifwithdrawals", false, "Display Unjustified Withdrawals")
	p2086 := flag.Bool("2086", false, "Display Cerfa 2086")
	pCoinAPIKey := flag.String("coinapi_key", "", "CoinAPI Key (https://www.coinapi.io/pricing?apikey)")
	pCSVBtcAddress := flag.String("btc_address", "", "Bitcoin Addresses CSV file")
	pCSVBtcPayment := flag.String("btc_payment", "", "Bitcoin Payments CSV file")
	pFloatBtcExclude := flag.Float64("btc_exclude", 0.0, "Exclude Bitcoin Amount")
	pCSVEthAddress := flag.String("eth_address", "", "Ethereum Addresses CSV file")
	pEtherscanAPIKey := flag.String("etherscan_apikey", "", "Etherscan API Key (https://etherscan.io/myapikey)")
	pCSVBinance := flag.String("binance", "", "Binance CSV file")
	pCSVBinanceExtended := flag.Bool("binance_extended", false, "Use Binance CSV file extended format")
	pCSVBitfinex := flag.String("bitfinex", "", "Bitfinex CSV file")
	pCSVCoinbase := flag.String("coinbase", "", "Coinbase CSV file")
	pCSVCdC := flag.String("cdc_app", "", "Crypto.com App CSV file")
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
	blkst := blockstream.New()
	if *pCSVBtcPayment != "" {
		recordFile, err := os.Open(*pCSVBtcPayment)
		if err != nil {
			log.Fatal("Error opening Bitcoin CSV Payments file:", err)
		}
		blkst.ParseCSVPayments(recordFile)
	}
	if *pCSVBtcAddress != "" {
		recordFile, err := os.Open(*pCSVBtcAddress)
		if err != nil {
			log.Fatal("Error opening Bitcoin CSV Addresses file:", err)
		}
		go blkst.ParseCSVAddresses(recordFile)
	}
	ethsc := etherscan.New()
	if *pCSVEthAddress != "" {
		recordFile, err := os.Open(*pCSVEthAddress)
		if err != nil {
			log.Fatal("Error opening Ethereum CSV Addresses file:", err)
		}
		ethsc.APIConnect(*pEtherscanAPIKey)
		go ethsc.ParseCSV(recordFile)
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
	if *pCSVCdC != "" {
		recordFile, err := os.Open(*pCSVCdC)
		if err != nil {
			log.Fatal("Error opening Crypto.com CSV file:", err)
		}
		err = cdc.ParseCSV(recordFile)
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
		err = ll.ParseCSV(recordFile)
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
	if *pCSVBtcAddress != "" {
		err := blkst.WaitFinish()
		if err != nil {
			log.Fatal("Error parsing Bitcoin CSV file:", err)
		}
	}
	// create Global Wallet up to Date
	global := make(wallet.Accounts)
	if *pFloatBtcExclude != 0.0 {
		t := wallet.TX{Timestamp: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC), Note: "Manual Exclusion"}
		t.Items = make(map[string][]wallet.Currency)
		t.Items["From"] = append(t.Items["From"], wallet.Currency{Code: "BTC", Amount: decimal.NewFromFloat(*pFloatBtcExclude)})
		global["Excludes"] = append(global["Excludes"], t)
	}
	global.Add(b.Accounts)
	global.Add(bf.Accounts)
	global.Add(cb.Accounts)
	global.Add(cdc.Accounts)
	global.Add(ll.Accounts)
	global.Add(lb.Accounts)
	global.Add(mm.Accounts)
	global.Add(mc.Accounts)
	global.Add(revo.Accounts)
	global.Add(ethsc.Accounts)
	global.Add(blkst.Accounts)
	global.FindTransfers()
	global.FindCashInOut()
	global.SortTXsByDate(true)
	loc, err := time.LoadLocation(*pLocation)
	if err != nil {
		log.Fatal("Error parsing Location:", err)
	}
	if *pStats {
		global.PrintStats()
	}
	if *pUnjustifWithdrawals {
		global.PrintUnjustifiedWithdrawals(loc)
	}
	// Debug
	if *pAccount != "" {
		if *pAccount == "Alls" {
			spew.Dump(global)
		} else {
			spew.Dump(global[*pAccount])
		}
	}
	// Construct global wallet up to date
	filterDate, err := time.ParseInLocation("2006-01-02T15:04:05", *pDate, loc)
	if err != nil {
		log.Fatal("Error parsing Date:", err)
	}
	globalWallet := global.GetWallets(filterDate, false)
	globalWallet.Println("Global Crypto")
	globalWalletTotalValue, err := globalWallet.CalculateTotalValue(*pNative)
	if err != nil {
		log.Fatal("Error Calculating Global Wallet:", err)
	} else {
		globalWalletTotalValue.Amount = globalWalletTotalValue.Amount.RoundBank(0)
		fmt.Print("Total Value :")
		globalWalletTotalValue.Println()
	}
	if *p2086 {
		var cessions Cessions
		var fractionCapital decimal.Decimal
		// source Bofip
		// RPPM - Plus-values sur biens meubles et taxe forfaitaire sur les objets précieux
		// - Cession d'actifs numériques à titre occasionnel - Base d'imposition
		// Cas des cessions antérieures au 1er Janvier 2019
		// 130
		// L'article 41 de la loi n° 2018-1317 du 28 décembre 2018 de finances pour 2019,
		// codifié à l'article 150 VH bis du CGI, s'appliquant aux cessions réalisées à
		// compter du 1er janvier 2019, il convient, pour la détermination du prix total
		// d'acquisition, de n'inclure dans ce dernier que les prix effectifs d'acquisition
		// des actifs détenus à cette date.
		// Ainsi, en cas de cessions réalisées antérieurement au 1er janvier 2019, il
		// convient notamment de ne pas inclure dans le prix total d'acquisition déclaré à
		// l'occasion de la première cession réalisée postérieurement à cette date, les
		// prix d'acquisition :
		// - mentionnés dans les déclarations de plus-values de cessions déclarées en
		//   application du droit en vigueur avant le 1er janvier 2019 ;
		// - n'ayant pas été déclarés, conformément au droit en vigueur avant le 1er
		//   janvier 2019 (cessions dont le prix de cession était inférieur à 5 000 € et
		//   ayant bénéficié de l'exonération prévue au 2° du II de l'article 150 UA du CGI
		//   par exemple) ;
		// - n'ayant pas été déclarés en contravention avec le droit en vigueur avant le
		//   1er janvier 2019.
		// Il est rappelé que les éventuelles plus-values réalisées antérieurement au
		// 1er janvier 2019 relèvent du droit de reprise de l'administration.
		// Remarque : En cas d'échange entre actifs numériques réalisé, même sans soulte,
		// antérieurement au 1er janvier 2019, le prix total d'acquisition à retenir à
		// compter du 1er janvier 2019 est constitué de la valeur de l'actif numérique
		// remis lors de cet échange (valeur à la date de cet échange). Corrélativement,
		// le prix d'acquisition retenu à l'occasion de cette cession n'est pas inclus dans
		// le prix total d'acquisition déclaré à compter du 1er janvier 2019.
		// Par ailleurs, les moins-values constatées lors des cessions réalisées
		// antérieurement au 1er janvier 2019 ne peuvent être imputées sur d'éventuelles
		// plus-values réalisées, quelle que soit leur date de réalisation.
		var prixTotalAcquisition decimal.Decimal
		date2019Jan1 := time.Date(2019, time.January, 1, 0, 0, 0, 0, loc)
		globalWallet2019Jan1 := global.GetWallets(date2019Jan1, false)
		// Consolidate all knowns TXs
		var allTXs wallet.TXs
		for k := range global {
			if k != "Transfers" { // Do not consider Transfers for initial prixTotalAcquisition
				allTXs = append(allTXs, global[k]...)
			}
		}
		allTXs.SortByDate(false)
		for crypto, quantity := range globalWallet2019Jan1.Currencies {
			if quantity.IsNegative() {
				globalWallet2019Jan1.Println("2019 Jan 1st Global")
				log.Fatal("Error Initial stock have a negative stock, some TXs are missing !")
			}
			var amountToFind decimal.Decimal
			amountToFind = quantity
			// spew.Dump(crypto)
			var fifoValue decimal.Decimal
			for _, tx := range allTXs {
				// Find all Tx before 2019 Jan 1st ...
				if tx.Timestamp.Before(date2019Jan1) {
					// ... that have the wanted crypto into Items["To"]
					for _, c := range tx.Items["To"] {
						if c.Code == crypto {
							// spew.Dump(tx)
							rate, err := c.GetExchangeRate(tx.Timestamp, *pNative)
							if err != nil {
								log.Println(err)
							} else {
								if amountToFind.LessThan(c.Amount) {
									fifoValue = fifoValue.Add(rate.Mul(amountToFind))
								} else {
									fifoValue = fifoValue.Add(rate.Mul(c.Amount))
								}
							}
							amountToFind = amountToFind.Sub(c.Amount)
						}
					}
					// ... and the ones consoming the wanted crypto
					for _, c := range tx.Items["From"] {
						if c.Code == crypto {
							amountToFind = amountToFind.Add(c.Amount)
						}
					}
					for _, c := range tx.Items["Fee"] {
						if c.Code == crypto {
							amountToFind = amountToFind.Add(c.Amount)
						}
					}
					if !amountToFind.IsPositive() {
						prixTotalAcquisition = prixTotalAcquisition.Add(fifoValue)
						// spew.Dump(prixTotalAcquisition)
						break
					}
				}
			}
			if amountToFind.IsPositive() {
				log.Println("Could not find enough TXs to calculate FIFO value for", crypto, "missing", amountToFind)
			}
		}
		// Consolidate all CashIn/CashOut TXs
		var cashInOut wallet.TXs
		cashInOut = append(cashInOut, global["CashIn"]...)
		cashInOut = append(cashInOut, global["CashOut"]...)
		cashInOut.SortByDate(true)
		// Calculate PV starting on 2019 Jan 1st
		for _, tx := range cashInOut {
			if tx.Timestamp.After(time.Date(2018, time.December, 31, 23, 59, 59, 999, loc)) {
				if tx.Items["To"][0].IsFiat() { // CashOut
					c := Cession{Date211: tx.Timestamp}
					// Valeur globale du portefeuille au moment de la cession
					// Il s’agit de la somme des valeurs, évaluées au moment de la cession
					// imposable, des différents actifs numériques et droits s'y rapportant,
					// détenus par le cédant avant de procéder à la cession, quel que soit
					// leur support de conservation (plateformes d’échanges, y compris
					// étrangères, serveurs personnels, dispositif de stockage hors-ligne,
					// etc.). Cette valorisation doit s’effectuer au moment de chaque cession
					// imposable en application de l’article 150 VH bis du CGI.
					globalWallet := global.GetWallets(tx.Timestamp, false)
					globalWalletTotalValue, err := globalWallet.CalculateTotalValue(*pNative)
					if err != nil {
						log.Println("Error Calculating Global Wallet at", tx.Timestamp, err)
					}
					// spew.Dump(globalWallet)
					c.ValeurPortefeuille212 = globalWalletTotalValue.Amount
					// Prix de cession
					// Il correspond au prix réel perçu ou à la valeur de la contrepartie
					// obtenue par le cédant lors de la cession.
					if tx.Items["To"][0].Code == *pNative {
						c.Prix213 = tx.Items["To"][0].Amount
					} else {
						var api wallet.CoinAPI
						rates, err := api.GetExchangeRates(tx.Timestamp, *pNative)
						if err != nil {
							log.Println("Error Getting Rates for", tx.Timestamp, err)
						} else {
							found := false
							for _, r := range rates.Rates {
								if r.Quote == tx.Items["To"][0].Code {
									c.Prix213 = tx.Items["To"][0].Amount.Mul(r.Rate)
									found = true
									break
								}
							}
							if !found {
								log.Println("Rate missing : CashOut integration into Prix213")
								spew.Dump(tx, c)
							}
						}
					}
					// Prix de cession - Frais
					// Il est réduit, sur justificatifs, des frais supportés par le cédant à
					// l’occasion de cette cession. Ces frais s'entendent, notamment, de
					// ceux perçus à l’occasion de l’opération imposable par les plateformes
					// où s'effectuent les cessions à titre onéreux d'actifs numériques ou
					// de droits s'y rapportant ainsi que de ceux perçus par les membres du
					// réseau (appelés "mineurs") chargés de vérifier et valider les
					// transactions qui s'y opèrent. Le paiement de ces frais de transaction
					// perçus par les plateformes ou les "mineurs" peut s'effectuer au moyen
					// d'actifs numériques. Or, dans ce cas, ce paiement est la contrepartie
					// d'un service fourni au cédant et constitue une opération imposable au
					// sens du I de l'article 150 VH bis du CGI. A titre de mesure de
					// simplification, il est toutefois admis que la cession en tant que
					// telle et les différentes prestations de services rendues en
					// contrepartie des frais perçus par les plateformes et les "mineurs"
					// soient assimilées à une seule et même opération de cession pour
					// l'application de l'article 150 VH bis du CGI, pour laquelle le
					// contribuable détermine une seule plus ou moins-value, en déduisant
					// ces frais du prix de cession.
					for _, f := range tx.Items["Fee"] {
						if f.Code == *pNative {
							c.Frais214 = c.Frais214.Add(f.Amount)
						} else {
							var api wallet.CoinAPI
							rates, err := api.GetExchangeRates(tx.Timestamp, *pNative)
							if err != nil {
								log.Println("Error Getting Rates for", tx.Timestamp, err)
							} else {
								found := false
								for _, r := range rates.Rates {
									if r.Quote == f.Code {
										c.Frais214 = c.Frais214.Add(f.Amount.Mul(r.Rate))
										found = true
										break
									}
								}
								if !found {
									log.Println("Rate missing : CashOut integration into Frais214")
									spew.Dump(tx, c)
								}
							}
						}
					}
					// Prix de cession - Soultes
					// Le prix de cession doit être majoré de la soulte que le cédant a
					// reçue lors de la cession ou minoré de la soulte qu’il a versée lors
					// de cette même cession.
					// c.SoulteRecueOuVersee216 = ???
					c.PrixTotalAcquisition220 = prixTotalAcquisition
					// Fractions de capital initial
					// Il s’agit de la fraction de capital contenue dans la valeur ou le
					// prix de chacune des différentes cessions d'actifs numériques à titre
					// gratuit ou onéreux réalisées antérieurement, hors opérations d’échange
					// ayant bénéficié du sursis d’imposition sans soulte.
					c.FractionDeCapital221 = fractionCapital
					// Soulte reçue en cas d’échanges antérieurs à la cession
					// Lorsqu’un ou plusieurs échanges avec soulte reçue par le cédant ont été
					// réalisés antérieurement à la cession imposable, le prix total d’acquisition
					// est minoré du montant des soultes. Indiquez donc les montants reçus.
					// c.SoulteRecueEnCasDechangeAnterieur222 = ???
					c.Calculate() // to have 217 and 223
					// spew.Dump(c)
					cessions = append(cessions, c)
					// Les frais déductibles, quels qu'ils soient, ne viennent pas en
					// diminution du prix de cession pour la détermination du quotient du
					// prix de cession sur la valeur globale du portefeuille (ils doivent
					// seulement être déduits du prix de cession qui constitue le premier
					// terme de la différence prévue dans la formule de calcul mentionnée
					// ci-dessus).
					coefCession := c.PrixNetDeSoulte217.Div(c.ValeurPortefeuille212)
					fractionAcquisition := coefCession.Mul(c.PrixTotalAcquisitionNet223)
					fractionCapital = fractionCapital.Add(fractionAcquisition)
				} else { // CashIn
					// Prix total d’acquisition du portefeuille
					// Le prix total d'acquisition du portefeuille d'actifs numériques est
					// égal à la somme de tous les prix acquittés en monnaie ayant cours
					// légal à l'occasion de l'ensemble des acquisitions d’actifs numériques
					// (sauf opérations d'échange ayant bénéficié du sursis d'imposition)
					// réalisées avant la cession, et de la valeur des biens ou services,
					// comprenant le cas échéant les soultes versées, fournis en
					// contrepartie de ces acquisitions.
					if tx.Items["From"][0].Code == *pNative {
						prixTotalAcquisition = prixTotalAcquisition.Add(tx.Items["From"][0].Amount)
					} else {
						var api wallet.CoinAPI
						rates, err := api.GetExchangeRates(tx.Timestamp, *pNative)
						if err != nil {
							log.Println("Error Getting Rates for", tx.Timestamp, err)
						} else {
							found := false
							for _, r := range rates.Rates {
								if r.Quote == tx.Items["From"][0].Code {
									prixTotalAcquisition = prixTotalAcquisition.Add(r.Rate.Mul(tx.Items["From"][0].Amount))
									found = true
									break
								}
							}
							if !found {
								log.Println("Rate missing : CashIn integration into prixTotalAcquisition")
								spew.Dump(tx)
							}
						}
					}
				}
			}
		}
		cessions.Println()
	}
	os.Exit(0)
}
