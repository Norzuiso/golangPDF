package main

import (
	"bufio"
	"fmt"
	"main/src/model"
	"os"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	err := generatePdf("Hello.pdf")
	if err != nil {
		panic(err)
	}
}

func generatePdf(fileName string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.SetTopMargin(30)

	countryList := make([]model.Invoice, 0, 8)
	loadData := func() {

		fl, err := os.Open("pagos.txt")
		if err == nil {
			scanner := bufio.NewScanner(fl)
			var invoiceTemp model.Invoice

			for scanner.Scan() {

				lineStr := scanner.Text()
				list := strings.Split(lineStr, ";")
				if len(list) == 4 {
					invoiceTemp.Date = list[0]
					invoiceTemp.Invoice = list[1]
					invoiceTemp.Provider = list[2]
					invoiceTemp.Amount = list[3]
					countryList = append(countryList, invoiceTemp)
				} else {
					err = fmt.Errorf("error tokenizing %s", lineStr)
				}
			}
			fl.Close()
			if len(countryList) == 0 {
				err = fmt.Errorf("error loading data from %s", "pagos.txt")
			}
		}
		if err != nil {
			pdf.SetError(err)
		}
	}

	basicTable := func() {
		left := 5.0
		pdf.SetX(left)
		for i, c := range countryList {
			pdf.SetX(left)

			if i == 0 {
				pdf.CellFormat(30, 6, "Fecha", "1", 0, "C", false, 0, "")
				pdf.CellFormat(35, 6, "Folio", "1", 0, "C", false, 0, "")
				pdf.CellFormat(105, 6, "Provedor", "1", 0, "C", false, 0, "")
				pdf.CellFormat(30, 6, "Monto", "1", 0, "C", false, 0, "")

			} else {

				pdf.SetFont("Arial", "", 12)
				pdf.CellFormat(30, 6, c.Date, "1", 0, "C", false, 0, "")
				pdf.CellFormat(35, 6, c.Invoice, "1", 0, "C", false, 0, "")
				pdf.CellFormat(105, 6, c.Provider, "1", 0, "L", false, 0, "")
				pdf.CellFormat(30, 6, c.Amount, "1", 0, "R", false, 0, "")

			}
			pdf.Ln(-1)
		}
	}

	pdf.SetHeaderFuncMode(func() {
		pdf.Image("src/images/Avatar.jpg", 10, 5, 50, 20, false, "", 0, "")
		pdf.SetY(5)
		pdf.SetFont("Arial", "B", 16)
		pdf.Cell(60, 30, "")
		pdf.CellFormat(0, 5, "Muebleria La Gaby", "1", 0, "C", false, 0, "")
		pdf.Ln(100)
	}, true)

	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 12)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()), "", 0, "R", false, 0, "")
	})

	pdf.Ln(100)
	pdf.AliasNbPages("")
	pdf.AddPage()

	loadData()
	basicTable()
	pdf.AliasNbPages("")
	return pdf.OutputFileAndClose(fileName)
}
