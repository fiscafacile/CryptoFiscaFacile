package main

import (
	"fmt"
	"strconv"
	"time"

	// "github.com/fiscafacile/CryptoFiscaFacile/wallet"
	// "github.com/davecgh/go-spew/spew"
	"github.com/shopspring/decimal"
)

type Cession struct {
	Date211                              time.Time
	ValeurPortefeuille212                decimal.Decimal
	Prix213                              decimal.Decimal
	Frais214                             decimal.Decimal
	PrixNetDeFrais215                    decimal.Decimal
	SoulteRecueOuVersee216               decimal.Decimal
	PrixNetDeSoulte217                   decimal.Decimal
	PrixNet218                           decimal.Decimal
	PrixTotalAcquisition220              decimal.Decimal
	FractionDeCapital221                 decimal.Decimal
	SoulteRecueEnCasDechangeAnterieur222 decimal.Decimal
	PrixTotalAcquisitionNet223           decimal.Decimal
	PlusMoinsValue                       decimal.Decimal
}

type Cessions []Cession

func (c *Cession) Calculate() {
	c.PrixNetDeFrais215 = c.Prix213.Sub(c.Frais214)
	c.PrixNetDeSoulte217 = c.Prix213.Sub(c.SoulteRecueOuVersee216)
	// c.PrixNetDeSoulte217 = c.Prix213.Add(c.SoulteRecueOuVersee216)
	c.PrixNet218 = c.Prix213.Sub(c.Frais214).Sub(c.SoulteRecueOuVersee216)
	// c.PrixNet218 = c.Prix213.Sub(c.Frais214).Add(c.SoulteRecueOuVersee216)
	c.PrixTotalAcquisitionNet223 = c.PrixTotalAcquisition220.Sub(c.FractionDeCapital221).Sub(c.SoulteRecueEnCasDechangeAnterieur222)
	if !c.ValeurPortefeuille212.IsZero() {
		c.PlusMoinsValue = c.PrixNet218.Sub(c.PrixTotalAcquisitionNet223.Mul(c.PrixNetDeSoulte217).Div(c.ValeurPortefeuille212))
	}
}

func (c Cession) Println() {
	fmt.Println("211 Date de la cession :", c.Date211.Format("02/01/2006"))
	fmt.Println("212 Valeur globale du portefeuille au moment de la cession :", c.ValeurPortefeuille212.RoundBank(0))
	fmt.Println("213 Prix de cession :", c.Prix213.RoundBank(0))
	fmt.Println("214 Frais de cession :", c.Frais214.RoundBank(0))
	fmt.Println("215 Prix de cession net des frais :", c.PrixNetDeFrais215.RoundBank(0))
	fmt.Println("216 Soulte reçue ou versée lors de la cession :", c.SoulteRecueOuVersee216.RoundBank(0))
	fmt.Println("217 Prix de cession net des soultes :", c.PrixNetDeSoulte217.RoundBank(0))
	fmt.Println("218 Prix de cession net des frais et soultes :", c.PrixNet218.RoundBank(0))
	fmt.Println("220 Prix total d’acquisition :", c.PrixTotalAcquisition220.RoundBank(0))
	fmt.Println("221 Fractions de capital initial contenues dans le prix total d’acquisition :", c.FractionDeCapital221.RoundBank(0))
	fmt.Println("222 Soultes reçues en cas d’échanges antérieurs à la cession :", c.SoulteRecueEnCasDechangeAnterieur222.RoundBank(0))
	fmt.Println("223 Prix total d’acquisition net :", c.PrixTotalAcquisitionNet223.RoundBank(0))
	fmt.Println("Plus-values et moins-values :", c.PlusMoinsValue.RoundBank(0))
}

func (ac Cessions) Println() {
	for year := 2019; year < time.Now().Year(); year++ {
		var plusMoinsValueGlobale decimal.Decimal
		fmt.Println("-------------------------")
		fmt.Println("| Cerfa 2086 année " + strconv.Itoa(year) + " |")
		fmt.Println("-------------------------")
		for _, c := range ac {
			if c.Date211.After(time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)) {
				if c.Date211.Before(time.Date(year, time.December, 31, 23, 59, 59, 999, time.UTC)) {
					c.Println()
					fmt.Println("-------------------------")
					plusMoinsValueGlobale = plusMoinsValueGlobale.Add(c.PlusMoinsValue)
				} else {
					break
				}
			}
		}
		fmt.Println("224 Plus-value ou moins-value globale :", plusMoinsValueGlobale.RoundBank(0))
		fmt.Println("-------------------------")
	}
}
