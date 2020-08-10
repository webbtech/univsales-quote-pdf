package pdf

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jung-kurt/gofpdf"
)

// Quote method
func (p *PDF) Quote() (err error) {

	p.setOutputFileName()
	titleStr := "Quote " + strconv.Itoa(p.q.Number) + " PDF"

	p.pdf = gofpdf.New("P", "mm", "Letter", "")
	pdf := p.pdf
	pdf.SetTitle(titleStr, false)
	pdf.SetAuthor(p.cfg.DocAuthor, false)

	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 9)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d of {nb}", pdf.PageNo()), "", 0, "C", false, 0, "")
	})
	pdf.AliasNbPages("")

	pdf.AddPage()
	p.quoteTitle()
	p.productHdr()
	p.groupList()
	p.windowList()
	p.otherList()
	p.quoteSummary()

	return err
}

func (p *PDF) quoteTitle() {

	pdf := p.pdf
	q := p.q

	var (
		rsp     *http.Response
		tp      string
		imgInfo gofpdf.ImageOptions
	)
	rsp, err := http.Get(p.cfg.LogoURI)
	if err == nil {
		tp = pdf.ImageTypeFromMime(rsp.Header["Content-Type"][0])
		imgInfo = gofpdf.ImageOptions{ImageType: tp}
		pdf.RegisterImageReader(p.cfg.LogoURI, tp, rsp.Body)
	} else {
		pdf.SetError(err)
	}
	customerName := fmt.Sprintf("%s %s", q.Customer.Name.First, q.Customer.Name.Last)
	addressLine1 := fmt.Sprintf("%s", q.JobSheetAddress.Street1)
	addressLine2 := fmt.Sprintf("%s, %s. %s", q.JobSheetAddress.City, q.JobSheetAddress.Province, q.JobSheetAddress.PostalCode)
	quoteNo := fmt.Sprintf("%d R%d", q.Number, q.Revision)

	pdf.SetFont("Arial", "", 12)
	pdf.SetDrawColor(200, 200, 200)
	pdf.ImageOptions(p.cfg.LogoURI, 55, 10, 0, 0, false, imgInfo, 0, fmt.Sprintf("https://%s", coDomain))
	pdf.SetFont("Arial", "", 10)
	pdf.MoveTo(150, 12)
	pdf.CellFormat(0, 5, coAddressStreet, "", 2, "", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("%s, %s %s", coAddressCity, coAddressProvince, coAddressPostal), "", 2, "", false, 0, "")
	pdf.SetTextColor(0, 0, 200)
	pdf.SetFont("Arial", "U", 10)
	pdf.CellFormat(0, 5, coDomain, "", 2, "", false, 0, fmt.Sprintf("https://%s", coDomain))
	pdf.CellFormat(0, 5, coEmail, "", 2, "", false, 0, fmt.Sprintf("mailto:%s", coEmail))

	pdf.MoveTo(75, 42)
	pdf.SetFont("Arial", "B", 18)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 5, "(905) 892-3030", "", 2, "", false, 0, "")

	pdf.MoveTo(10, 50)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "Quote", "", 2, "", false, 0, "")

	pdf.Ln(1)
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 5.5, customerName, "", 2, "", false, 0, "")
	pdf.Ln(1)
	pdf.CellFormat(0, 5.5, addressLine1, "", 2, "", false, 0, "")
	pdf.CellFormat(0, 4, addressLine2, "", 2, "", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 5.5, "(job location)", "", 2, "", false, 0, "")
	pdf.Ln(2)

	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(0, 0, 0)
	if v, ok := q.Customer.PhoneMap["mobile"]; ok {
		pdf.CellFormat(0, 5.5, fmt.Sprintf("Mobile %s", v), "", 2, "", false, 0, "")
	}
	if v, ok := q.Customer.PhoneMap["home"]; ok {
		pdf.CellFormat(0, 5.5, fmt.Sprintf("Home %s", v), "", 2, "", false, 0, "")
	}
	pdf.Ln(2)
	pdf.SetTextColor(0, 0, 200)
	pdf.SetFont("Arial", "U", 12)
	if q.Customer.Email != "" {
		pdf.CellFormat(0, 5.5, q.Customer.Email, "", 2, "", false, 0, fmt.Sprintf("mailto:%s", q.Customer.Email))
	}
	pdf.SetTextColor(0, 0, 0)

	pdf.MoveTo(150, 50)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 6, "Quotation Date", "", 2, "", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 7, q.UpdatedAt.Format(dateFormat), "", 2, "", false, 0, "")
	pdf.CellFormat(0, 5, "", "", 2, "", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 6, "Quotation No/Rev", "", 2, "", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 7, quoteNo, "", 2, "", false, 0, "")
}
