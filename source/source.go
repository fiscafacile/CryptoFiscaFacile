package source

import (
	"errors"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const URL_3916 = "https://www.impots.gouv.fr/portail/files/formulaires/3916/2021/3916_3425.pdf"

type Source struct {
	Crypto        bool
	AccountNumber string
	OpeningDate   time.Time
	ClosingDate   time.Time
	LegalName     string
	Address       string
	URL           string
}

type Sources map[string]Source

func (ss Sources) Add(srcs Sources) {
	for k, v := range srcs {
		ss[k] = v
	}
}

func (ss Sources) ToXlsx(filename string, loc *time.Location) error {
	if len(ss) > 0 {
		f := excelize.NewFile()
		for src, s := range ss {
			f.NewSheet(src)
			if s.Crypto {
				f.SetCellValue(src, "A1", "4.1 Désignation du compte d'actifs numériques ouvert, détenu, utilisé ou clos à l'étranger")
				f.SetCellValue(src, "A2", "Numéro de compte")
				f.SetCellValue(src, "B2", s.AccountNumber)
				f.SetCellValue(src, "A3", "Date d'ouverture*")
				f.SetCellValue(src, "B3", s.OpeningDate.In(loc).Format("02-01-2006"))
				f.SetCellValue(src, "A4", "Date de clôture*")
				f.SetCellValue(src, "B4", s.ClosingDate.In(loc).Format("02-01-2006"))
				f.SetCellValue(src, "A5", "Designation de l'organisme gestionnaire du compte")
				f.SetCellValue(src, "B5", s.LegalName)
				f.SetCellValue(src, "A6", "Adresse de l'organisme gestionnaire du compte")
				f.SetCellValue(src, "B6", s.Address)
				f.SetCellValue(src, "A7", "URL du site internet de l'organisme gestionnaire du compte")
				f.SetCellValue(src, "B7", s.URL)
				f.SetCellValue(src, "A8", "* Ce ne sont pas des dates formelles, juste des estimations en fonctions de vos transactions")
			} else {
				f.SetCellValue(src, "A1", "3.1 Désignation du compte bancaire ouvert, détenu, utilisé ou clos à l'étranger")
				f.SetCellValue(src, "A2", "Numéro de compte")
				f.SetCellValue(src, "B2", s.AccountNumber)
				f.SetCellValue(src, "A3", "Caractéristiques du compte")
				f.SetCellValue(src, "B3", "[x] Compte courant")
				f.SetCellValue(src, "A4", "Date d'ouverture*")
				f.SetCellValue(src, "B4", s.OpeningDate.In(loc).Format("02-01-2006"))
				f.SetCellValue(src, "A5", "Date de clôture*")
				f.SetCellValue(src, "B5", s.ClosingDate.In(loc).Format("02-01-2006"))
				f.SetCellValue(src, "A6", "Designation de l'organisme gestionnaire du compte")
				f.SetCellValue(src, "B6", s.LegalName)
				f.SetCellValue(src, "A7", "Adresse de l'organisme gestionnaire du compte")
				f.SetCellValue(src, "B7", s.Address)
				f.SetCellValue(src, "A8", "* Ce ne sont pas des dates formelles, juste des estimations en fonctions de vos transactions")
			}
			f.SetColWidth(src, "A", "B", 83)
		}
		f.DeleteSheet("Sheet1")
		return f.SaveAs(filename)
	} else {
		return errors.New("No Sources Available")
	}
}
