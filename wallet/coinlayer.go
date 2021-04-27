package wallet

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/nanobox-io/golang-scribble"
	// "github.com/shopspring/decimal"
	"gopkg.in/resty.v0"
)

type HistoricalData struct {
	Success    bool               `json:"success"`
	Terms      string             `json:"terms"`
	Privacy    string             `json:"privacy"`
	Timestamp  int                `json:"timestamp"`
	Target     string             `json:"target"`
	Historical bool               `json:"historical"`
	Date       string             `json:"date"`
	Rates      map[string]float64 `json:"rates"`
	/*	struct {
			Num611 float64 `json:"611,omitempty"`
			Abc    float64 `json:"ABC,omitempty"`
			Acp    float64 `json:"ACP,omitempty"`
			Act    float64 `json:"ACT,omitempty"`
			ActS   float64 `json:"ACT*,omitempty"`
			Ada    float64 `json:"ADA,omitempty"`
			Adcn   float64 `json:"ADCN,omitempty"`
			Adl    float64 `json:"ADL,omitempty"`
			Adx    float64 `json:"ADX,omitempty"`
			Adz    float64 `json:"ADZ,omitempty"`
			Ae     float64 `json:"AE,omitempty"`
			Agi    float64 `json:"AGI,omitempty"`
			Aib    float64 `json:"AIB,omitempty"`
			Aidoc  float64 `json:"AIDOC,omitempty"`
			Aion   float64 `json:"AION,omitempty"`
			Air    float64 `json:"AIR,omitempty"`
			Alt    float64 `json:"ALT,omitempty"`
			Amb    float64 `json:"AMB,omitempty"`
			Amm    float64 `json:"AMM,omitempty"`
			Ant    float64 `json:"ANT,omitempty"`
			Apc    float64 `json:"APC,omitempty"`
			Appc   float64 `json:"APPC,omitempty"`
			Arc    float64 `json:"ARC,omitempty"`
			Arct   float64 `json:"ARCT,omitempty"`
			Ardr   float64 `json:"ARDR,omitempty"`
			Ark    float64 `json:"ARK,omitempty"`
			Arn    float64 `json:"ARN,omitempty"`
			Asafe2 float64 `json:"ASAFE2,omitempty"`
			Ast    float64 `json:"AST,omitempty"`
			Atb    float64 `json:"ATB,omitempty"`
			Atm    float64 `json:"ATM,omitempty"`
			Aurs   float64 `json:"AURS,omitempty"`
			Avt    float64 `json:"AVT,omitempty"`
			Bar    float64 `json:"BAR,omitempty"`
			Bash   float64 `json:"BASH,omitempty"`
			Bat    float64 `json:"BAT,omitempty"`
			Bay    float64 `json:"BAY,omitempty"`
			Bbp    float64 `json:"BBP,omitempty"`
			Bcd    float64 `json:"BCD,omitempty"`
			Bch    float64 `json:"BCH,omitempty"`
			Bcn    float64 `json:"BCN,omitempty"`
			Bcpt   float64 `json:"BCPT,omitempty"`
			Bee    float64 `json:"BEE,omitempty"`
			Bio    float64 `json:"BIO,omitempty"`
			Blc    float64 `json:"BLC,omitempty"`
			Block  float64 `json:"BLOCK,omitempty"`
			Blu    float64 `json:"BLU,omitempty"`
			Blz    float64 `json:"BLZ,omitempty"`
			Bmc    float64 `json:"BMC,omitempty"`
			Bnb    float64 `json:"BNB,omitempty"`
			Bnt    float64 `json:"BNT,omitempty"`
			Bost   float64 `json:"BOST,omitempty"`
			Bq     float64 `json:"BQ,omitempty"`
			Bqx    float64 `json:"BQX,omitempty"`
			Brd    float64 `json:"BRD,omitempty"`
			Brit   float64 `json:"BRIT,omitempty"`
			Bt1    float64 `json:"BT1,omitempty"`
			Bt2    float64 `json:"BT2,omitempty"`
			Btc    float64 `json:"BTC,omitempty"`
			Btca   float64 `json:"BTCA,omitempty"`
			Btcs   float64 `json:"BTCS,omitempty"`
			Btcz   float64 `json:"BTCZ,omitempty"`
			Btg    float64 `json:"BTG,omitempty"`
			Btlc   float64 `json:"BTLC,omitempty"`
			Btm    float64 `json:"BTM,omitempty"`
			BtmS   float64 `json:"BTM*,omitempty"`
			Btq    float64 `json:"BTQ,omitempty"`
			Bts    float64 `json:"BTS,omitempty"`
			Btx    float64 `json:"BTX,omitempty"`
			Burst  float64 `json:"BURST,omitempty"`
			Calc   float64 `json:"CALC,omitempty"`
			Cas    float64 `json:"CAS,omitempty"`
			Cat    float64 `json:"CAT,omitempty"`
			Ccrb   float64 `json:"CCRB,omitempty"`
			Cdt    float64 `json:"CDT,omitempty"`
			Cesc   float64 `json:"CESC,omitempty"`
			Chat   float64 `json:"CHAT,omitempty"`
			Cj     float64 `json:"CJ,omitempty"`
			Cl     float64 `json:"CL,omitempty"`
			Cld    float64 `json:"CLD,omitempty"`
			Cloak  float64 `json:"CLOAK,omitempty"`
			Cmt    float64 `json:"CMT*,omitempty"`
			Cnd    float64 `json:"CND,omitempty"`
			Cnx    float64 `json:"CNX,omitempty"`
			Cpc    float64 `json:"CPC,omitempty"`
			Crave  float64 `json:"CRAVE,omitempty"`
			Crc    float64 `json:"CRC,omitempty"`
			Cre    float64 `json:"CRE,omitempty"`
			Crw    float64 `json:"CRW,omitempty"`
			Cto    float64 `json:"CTO,omitempty"`
			Ctr    float64 `json:"CTR,omitempty"`
			Cvc    float64 `json:"CVC,omitempty"`
			Das    float64 `json:"DAS,omitempty"`
			Dash   float64 `json:"DASH,omitempty"`
			Dat    float64 `json:"DAT,omitempty"`
			Data   float64 `json:"DATA,omitempty"`
			Dbc    float64 `json:"DBC,omitempty"`
			Dbet   float64 `json:"DBET,omitempty"`
			Dcn    float64 `json:"DCN,omitempty"`
			Dcr    float64 `json:"DCR,omitempty"`
			Dct    float64 `json:"DCT,omitempty"`
			Deep   float64 `json:"DEEP,omitempty"`
			Dent   float64 `json:"DENT,omitempty"`
			Dgb    float64 `json:"DGB,omitempty"`
			Dgd    float64 `json:"DGD,omitempty"`
			Dim    float64 `json:"DIM,omitempty"`
			Dime   float64 `json:"DIME,omitempty"`
			Dmd    float64 `json:"DMD,omitempty"`
			Dnt    float64 `json:"DNT,omitempty"`
			Doge   float64 `json:"DOGE,omitempty"`
			Drgn   float64 `json:"DRGN,omitempty"`
			Drz    float64 `json:"DRZ,omitempty"`
			Dsh    float64 `json:"DSH,omitempty"`
			Dta    float64 `json:"DTA,omitempty"`
			Ec     float64 `json:"EC,omitempty"`
			Edg    float64 `json:"EDG,omitempty"`
			Edo    float64 `json:"EDO,omitempty"`
			Edr    float64 `json:"EDR,omitempty"`
			Eko    float64 `json:"EKO,omitempty"`
			Ela    float64 `json:"ELA,omitempty"`
			Elf    float64 `json:"ELF,omitempty"`
			Emc    float64 `json:"EMC,omitempty"`
			Emgo   float64 `json:"EMGO,omitempty"`
			Eng    float64 `json:"ENG,omitempty"`
			Enj    float64 `json:"ENJ,omitempty"`
			Eos    float64 `json:"EOS,omitempty"`
			Ert    float64 `json:"ERT,omitempty"`
			Etc    float64 `json:"ETC,omitempty"`
			Eth    float64 `json:"ETH,omitempty"`
			Etn    float64 `json:"ETN,omitempty"`
			Etp    float64 `json:"ETP,omitempty"`
			Ett    float64 `json:"ETT,omitempty"`
			Evr    float64 `json:"EVR,omitempty"`
			Evx    float64 `json:"EVX,omitempty"`
			Fct    float64 `json:"FCT,omitempty"`
			Flp    float64 `json:"FLP,omitempty"`
			Fota   float64 `json:"FOTA,omitempty"`
			Frst   float64 `json:"FRST,omitempty"`
			Fuel   float64 `json:"FUEL,omitempty"`
			Fun    float64 `json:"FUN,omitempty"`
			Func   float64 `json:"FUNC,omitempty"`
			Futc   float64 `json:"FUTC,omitempty"`
			Game   float64 `json:"GAME,omitempty"`
			Gas    float64 `json:"GAS,omitempty"`
			Gbyte  float64 `json:"GBYTE,omitempty"`
			Gmx    float64 `json:"GMX,omitempty"`
			Gno    float64 `json:"GNO,omitempty"`
			Gnt    float64 `json:"GNT,omitempty"`
			Gnx    float64 `json:"GNX,omitempty"`
			Grc    float64 `json:"GRC,omitempty"`
			Grs    float64 `json:"GRS,omitempty"`
			Grwi   float64 `json:"GRWI,omitempty"`
			Gtc    float64 `json:"GTC,omitempty"`
			Gto    float64 `json:"GTO,omitempty"`
			Gup    float64 `json:"GUP,omitempty"`
			Gvt    float64 `json:"GVT,omitempty"`
			Gxs    float64 `json:"GXS,omitempty"`
			Hac    float64 `json:"HAC,omitempty"`
			Hnc    float64 `json:"HNC,omitempty"`
			Hsr    float64 `json:"HSR,omitempty"`
			Hst    float64 `json:"HST,omitempty"`
			Hvn    float64 `json:"HVN,omitempty"`
			Icn    float64 `json:"ICN,omitempty"`
			Icos   float64 `json:"ICOS,omitempty"`
			Icx    float64 `json:"ICX,omitempty"`
			Ignis  float64 `json:"IGNIS,omitempty"`
			Ilc    float64 `json:"ILC,omitempty"`
			Ink    float64 `json:"INK,omitempty"`
			Ins    float64 `json:"INS,omitempty"`
			Insn   float64 `json:"INSN,omitempty"`
			Int    float64 `json:"INT,omitempty"`
			Iop    float64 `json:"IOP,omitempty"`
			Iost   float64 `json:"IOST,omitempty"`
			Itc    float64 `json:"ITC,omitempty"`
			Kcs    float64 `json:"KCS,omitempty"`
			Kick   float64 `json:"KICK,omitempty"`
			Kin    float64 `json:"KIN,omitempty"`
			Klc    float64 `json:"KLC,omitempty"`
			Kmd    float64 `json:"KMD,omitempty"`
			Knc    float64 `json:"KNC,omitempty"`
			Krb    float64 `json:"KRB,omitempty"`
			La     float64 `json:"LA,omitempty"`
			Lend   float64 `json:"LEND,omitempty"`
			Leo    float64 `json:"LEO,omitempty"`
			Linda  float64 `json:"LINDA,omitempty"`
			Link   float64 `json:"LINK,omitempty"`
			Loc    float64 `json:"LOC,omitempty"`
			Log    float64 `json:"LOG,omitempty"`
			Lrc    float64 `json:"LRC,omitempty"`
			Lsk    float64 `json:"LSK,omitempty"`
			Ltc    float64 `json:"LTC,omitempty"`
			Lun    float64 `json:"LUN,omitempty"`
			Lux    float64 `json:"LUX,omitempty"`
			Maid   float64 `json:"MAID,omitempty"`
			Mana   float64 `json:"MANA,omitempty"`
			Mcap   float64 `json:"MCAP,omitempty"`
			Mco    float64 `json:"MCO,omitempty"`
			Mda    float64 `json:"MDA,omitempty"`
			Mds    float64 `json:"MDS,omitempty"`
			Miota  float64 `json:"MIOTA,omitempty"`
			Mkr    float64 `json:"MKR,omitempty"`
			Mln    float64 `json:"MLN,omitempty"`
			Mnx    float64 `json:"MNX,omitempty"`
			Mod    float64 `json:"MOD,omitempty"`
			Moin   float64 `json:"MOIN,omitempty"`
			Mona   float64 `json:"MONA,omitempty"`
			Mtl    float64 `json:"MTL,omitempty"`
			Mtn    float64 `json:"MTN*,omitempty"`
			Mtx    float64 `json:"MTX,omitempty"`
			Nas    float64 `json:"NAS,omitempty"`
			Nav    float64 `json:"NAV,omitempty"`
			Nbt    float64 `json:"NBT,omitempty"`
			Ndc    float64 `json:"NDC,omitempty"`
			Nebl   float64 `json:"NEBL,omitempty"`
			Neo    float64 `json:"NEO,omitempty"`
			Neu    float64 `json:"NEU,omitempty"`
			Newb   float64 `json:"NEWB,omitempty"`
			Ngc    float64 `json:"NGC,omitempty"`
			Nkc    float64 `json:"NKC,omitempty"`
			Nlc2   float64 `json:"NLC2,omitempty"`
			Nmc    float64 `json:"NMC,omitempty"`
			Nmr    float64 `json:"NMR,omitempty"`
			Nuls   float64 `json:"NULS,omitempty"`
			Nvc    float64 `json:"NVC,omitempty"`
			Nxt    float64 `json:"NXT,omitempty"`
			Oax    float64 `json:"OAX,omitempty"`
			Obits  float64 `json:"OBITS,omitempty"`
			Oc     float64 `json:"OC,omitempty"`
			Ocn    float64 `json:"OCN,omitempty"`
			Odn    float64 `json:"ODN,omitempty"`
			Ok     float64 `json:"OK,omitempty"`
			Omg    float64 `json:"OMG,omitempty"`
			Omni   float64 `json:"OMNI,omitempty"`
			Ore    float64 `json:"ORE,omitempty"`
			Orme   float64 `json:"ORME,omitempty"`
			Ost    float64 `json:"OST,omitempty"`
			Otn    float64 `json:"OTN,omitempty"`
			Otx    float64 `json:"OTX,omitempty"`
			Oxy    float64 `json:"OXY,omitempty"`
			Part   float64 `json:"PART,omitempty"`
			Pay    float64 `json:"PAY,omitempty"`
			Pbt    float64 `json:"PBT,omitempty"`
			Pcs    float64 `json:"PCS,omitempty"`
			Pivx   float64 `json:"PIVX,omitempty"`
			Pizza  float64 `json:"PIZZA,omitempty"`
			Plbt   float64 `json:"PLBT,omitempty"`
			Plr    float64 `json:"PLR,omitempty"`
			Poe    float64 `json:"POE,omitempty"`
			Poly   float64 `json:"POLY,omitempty"`
			Posw   float64 `json:"POSW,omitempty"`
			Powr   float64 `json:"POWR,omitempty"`
			Ppc    float64 `json:"PPC,omitempty"`
			Ppt    float64 `json:"PPT,omitempty"`
			Ppy    float64 `json:"PPY,omitempty"`
			Prc    float64 `json:"PRC,omitempty"`
			Pres   float64 `json:"PRES,omitempty"`
			Prg    float64 `json:"PRG,omitempty"`
			Prl    float64 `json:"PRL,omitempty"`
			Pro    float64 `json:"PRO,omitempty"`
			Pura   float64 `json:"PURA,omitempty"`
			Put    float64 `json:"PUT,omitempty"`
			Qash   float64 `json:"QASH,omitempty"`
			Qau    float64 `json:"QAU,omitempty"`
			Qsp    float64 `json:"QSP,omitempty"`
			Qtum   float64 `json:"QTUM,omitempty"`
			Qun    float64 `json:"QUN,omitempty"`
			R      float64 `json:"R,omitempty"`
			Rbies  float64 `json:"RBIES,omitempty"`
			Rcn    float64 `json:"RCN,omitempty"`
			Rdd    float64 `json:"RDD,omitempty"`
			Rdn    float64 `json:"RDN,omitempty"`
			RdnS   float64 `json:"RDN*,omitempty"`
			Rebl   float64 `json:"REBL,omitempty"`
			Ree    float64 `json:"REE,omitempty"`
			Rep    float64 `json:"REP,omitempty"`
			Req    float64 `json:"REQ,omitempty"`
			Rev    float64 `json:"REV,omitempty"`
			Rgc    float64 `json:"RGC,omitempty"`
			Rhoc   float64 `json:"RHOC,omitempty"`
			Riya   float64 `json:"RIYA,omitempty"`
			Rkc    float64 `json:"RKC,omitempty"`
			Rlc    float64 `json:"RLC,omitempty"`
			Rpx    float64 `json:"RPX,omitempty"`
			Ruff   float64 `json:"RUFF,omitempty"`
			Salt   float64 `json:"SALT,omitempty"`
			San    float64 `json:"SAN,omitempty"`
			Sbc    float64 `json:"SBC,omitempty"`
			Sc     float64 `json:"SC,omitempty"`
			Sent   float64 `json:"SENT,omitempty"`
			Shift  float64 `json:"SHIFT,omitempty"`
			Sib    float64 `json:"SIB,omitempty"`
			Smart  float64 `json:"SMART,omitempty"`
			Smly   float64 `json:"SMLY,omitempty"`
			Smt    float64 `json:"SMT*,omitempty"`
			Snc    float64 `json:"SNC,omitempty"`
			Sngls  float64 `json:"SNGLS,omitempty"`
			Snm    float64 `json:"SNM,omitempty"`
			Snt    float64 `json:"SNT,omitempty"`
			Spk    float64 `json:"SPK,omitempty"`
			Srn    float64 `json:"SRN,omitempty"`
			Steem  float64 `json:"STEEM,omitempty"`
			Stk    float64 `json:"STK,omitempty"`
			Storj  float64 `json:"STORJ,omitempty"`
			Strat  float64 `json:"STRAT,omitempty"`
			Stu    float64 `json:"STU,omitempty"`
			Stx    float64 `json:"STX,omitempty"`
			Sub    float64 `json:"SUB,omitempty"`
			Sur    float64 `json:"SUR,omitempty"`
			Swftc  float64 `json:"SWFTC,omitempty"`
			Sys    float64 `json:"SYS,omitempty"`
			Taas   float64 `json:"TAAS,omitempty"`
			Tesla  float64 `json:"TESLA,omitempty"`
			Thc    float64 `json:"THC,omitempty"`
			Theta  float64 `json:"THETA,omitempty"`
			Ths    float64 `json:"THS,omitempty"`
			Tio    float64 `json:"TIO,omitempty"`
			Tkn    float64 `json:"TKN,omitempty"`
			Tky    float64 `json:"TKY,omitempty"`
			Tnb    float64 `json:"TNB,omitempty"`
			Tnt    float64 `json:"TNT,omitempty"`
			Toa    float64 `json:"TOA,omitempty"`
			Trc    float64 `json:"TRC,omitempty"`
			Trig   float64 `json:"TRIG,omitempty"`
			Trst   float64 `json:"TRST,omitempty"`
			Trump  float64 `json:"TRUMP,omitempty"`
			Trx    float64 `json:"TRX,omitempty"`
			Ubq    float64 `json:"UBQ,omitempty"`
			Uno    float64 `json:"UNO,omitempty"`
			Unrc   float64 `json:"UNRC,omitempty"`
			Uqc    float64 `json:"UQC,omitempty"`
			Usdt   float64 `json:"USDT,omitempty"`
			Utk    float64 `json:"UTK,omitempty"`
			Utt    float64 `json:"UTT,omitempty"`
			Vee    float64 `json:"VEE,omitempty"`
			Ven    float64 `json:"VEN,omitempty"`
			Veri   float64 `json:"VERI,omitempty"`
			Via    float64 `json:"VIA,omitempty"`
			Vib    float64 `json:"VIB,omitempty"`
			Vibe   float64 `json:"VIBE,omitempty"`
			Voise  float64 `json:"VOISE,omitempty"`
			Vox    float64 `json:"VOX,omitempty"`
			Vrs    float64 `json:"VRS,omitempty"`
			Vtc    float64 `json:"VTC,omitempty"`
			Vuc    float64 `json:"VUC,omitempty"`
			Wabi   float64 `json:"WABI,omitempty"`
			Waves  float64 `json:"WAVES,omitempty"`
			Wax    float64 `json:"WAX,omitempty"`
			Wc     float64 `json:"WC,omitempty"`
			Wgr    float64 `json:"WGR,omitempty"`
			Wings  float64 `json:"WINGS,omitempty"`
			Wlk    float64 `json:"WLK,omitempty"`
			Wop    float64 `json:"WOP,omitempty"`
			Wpr    float64 `json:"WPR,omitempty"`
			Wrc    float64 `json:"WRC,omitempty"`
			Wtc    float64 `json:"WTC,omitempty"`
			Xaur   float64 `json:"XAUR,omitempty"`
			Xbp    float64 `json:"XBP,omitempty"`
			Xby    float64 `json:"XBY,omitempty"`
			Xcp    float64 `json:"XCP,omitempty"`
			Xcxt   float64 `json:"XCXT,omitempty"`
			Xdn    float64 `json:"XDN,omitempty"`
			Xem    float64 `json:"XEM,omitempty"`
			Xgb    float64 `json:"XGB,omitempty"`
			Xhi    float64 `json:"XHI,omitempty"`
			Xid    float64 `json:"XID,omitempty"`
			Xlm    float64 `json:"XLM,omitempty"`
			Xmr    float64 `json:"XMR,omitempty"`
			Xnc    float64 `json:"XNC,omitempty"`
			Xrb    float64 `json:"XRB,omitempty"`
			Xrp    float64 `json:"XRP,omitempty"`
			Xto    float64 `json:"XTO,omitempty"`
			Xtz    float64 `json:"XTZ,omitempty"`
			Xuc    float64 `json:"XUC,omitempty"`
			Xvg    float64 `json:"XVG,omitempty"`
			Xzc    float64 `json:"XZC,omitempty"`
			Yee    float64 `json:"YEE,omitempty"`
			Yoc    float64 `json:"YOC,omitempty"`
			Yoyow  float64 `json:"YOYOW,omitempty"`
			Zbc    float64 `json:"ZBC,omitempty"`
			Zcl    float64 `json:"ZCL,omitempty"`
			Zec    float64 `json:"ZEC,omitempty"`
			Zen    float64 `json:"ZEN,omitempty"`
			Zil    float64 `json:"ZIL,omitempty"`
			Zny    float64 `json:"ZNY,omitempty"`
			Zrx    float64 `json:"ZRX,omitempty"`
			Zsc    float64 `json:"ZSC,omitempty"`
		} `json:"rates"`
	*/
}

type CoinLayer struct {
}

func CoinLayerSetKey(key string) error {
	return os.Setenv("COINLAYER_KEY", key)
}

func (api CoinLayer) GetExchangeRates(date time.Time, native string) (rates HistoricalData, err error) {
	db, err := scribble.New("./Cache", nil)
	if err != nil {
		return
	}
	err = db.Read("CoinLayer", native+"-"+date.UTC().Format("2006-01-02"), &rates)
	if err != nil {
		if os.Getenv("COINLAYER_KEY") == "" {
			return rates, errors.New("Need CoinLayer Key")
		}
		url := "http://api.coinlayer.com/" + date.UTC().Format("2006-01-02") + "?access_key=" + os.Getenv("COINLAYER_KEY") + "&target=" + native
		resp, err := resty.R().SetHeaders(map[string]string{
			"Accept": "application/json",
		}).Get(url)
		if err != nil {
			return rates, err
		}
		if resp.StatusCode() != http.StatusOK {
			err = errors.New("Error Status : " + strconv.Itoa(resp.StatusCode()))
			return rates, err
		}
		err = json.Unmarshal(resp.Body(), &rates)
		if err != nil {
			return rates, err
		}
		err = db.Write("CoinLayer", native+"-"+date.UTC().Format("2006-01-02"), rates)
		return rates, err
	}
	return
}
