package pdf

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
)

// Various constants
const (
	pfSize        float64 = 10
	topBr                 = 1.5
	midBr         float64 = 9
	hSep          float64 = 2
	newLnMaxLen   int     = 30
	HSTMultiplier float32 = 1.15
	moneyFormat           = "#,###.##"
)

// productHdr method
func (p *PDF) productHdr() {
	pdf := p.pdf
	pdf.MoveTo(9, 98)
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(9, 6, "Qty", "", 0, "", false, 0, "")
	pdf.CellFormat(23, 6, "Rms", "", 0, "", false, 0, "")
	pdf.CellFormat(46, 6, "Product", "", 0, "", false, 0, "")
	pdf.CellFormat(82, 6, "Options", "", 0, "", false, 0, "")
	pdf.CellFormat(23, 6, "Each", "", 0, "", false, 0, "")
	pdf.CellFormat(20, 6, "Total", "", 0, "", false, 0, "")
	pdf.SetDrawColor(0, 0, 0)
	pdf.Line(9, 105, 205, 105)
	pdf.Ln(7)
}

// groupList method
func (p *PDF) groupList() {

	pdf := p.pdf
	q := p.q

	var opts string
	if len(q.Items.Group) <= 0 {
		return
	}

	pdf.SetFillColor(100, 100, 100)

	for _, g := range q.Items.Group {

		if nil != g.Specs["options"] {
			opts = g.Specs["options"].(string)
		}

		pdf.Ln(topBr)
		pdf.SetFont("Arial", "", pfSize)
		pdf.CellFormat(5, 6, fmt.Sprintf("%d", g.Qty), "", 0, "", false, 0, "")
		pdf.CellFormat(26, 6, fmt.Sprintf("%s", strings.Join(g.Rooms, ", ")), "", 0, "", false, 0, "")
		xPos := pdf.GetX()
		yPos := pdf.GetY()
		pdf.MultiCell(40, 4.5, setNewLines(g.Specs["groupTypeDescription"]), "", "L", false)
		pdf.MoveTo(xPos+40, yPos)
		pdf.CellFormat(hSep, 6, "", "", 0, "", false, 0, "")
		xPos = pdf.GetX() + 78
		yPos = pdf.GetY()
		pdf.MultiCell(78, 4.5, fmt.Sprintf("%s", opts), "", "", false)
		pdf.MoveTo(xPos, yPos)
		pdf.CellFormat(20, 6, fmt.Sprintf("%.2f", g.Costs.Unit), "", 0, "R", false, 0, "")
		pdf.CellFormat(hSep, 6, "", "", 0, "", false, 0, "")
		pdf.CellFormat(20, 6, fmt.Sprintf("%.2f", g.Costs.Total), "", 0, "R", false, 0, "")
		pdf.Ln(midBr)
		pdf.SetDrawColor(200, 200, 200)
		pdf.CellFormat(0, 6, "", "B", 2, "", false, 0, "")
	}
}

// windowList method
func (p *PDF) windowList() {

	pdf := p.pdf
	q := p.q

	pdf.SetFillColor(100, 100, 100)
	if len(q.Items.Window) <= 0 {
		return
	}

	for _, g := range q.Items.Window {

		pdf.Ln(topBr)
		pdf.SetFont("Arial", "", pfSize)
		pdf.CellFormat(5, 6, fmt.Sprintf("%d", g.Qty), "", 0, "", false, 0, "")
		pdf.CellFormat(26, 6, fmt.Sprintf("%s", strings.Join(g.Rooms, ", ")), "", 0, "", false, 0, "")
		pdf.CellFormat(40, 6, fmt.Sprintf("%s", g.ProductName), "", 0, "", false, 0, "")
		pdf.CellFormat(hSep, 6, "", "", 0, "", false, 0, "")
		xPos := pdf.GetX() + 76
		yPos := pdf.GetY()
		pdf.MultiCell(78, 4.5, fmt.Sprintf("%s", g.Specs["options"]), "", "", false)
		pdf.MoveTo(xPos, yPos)
		pdf.CellFormat(hSep, 6, "", "", 0, "", false, 0, "")
		pdf.CellFormat(20, 6, fmt.Sprintf("%.2f", g.Costs.Unit), "", 0, "R", false, 0, "")
		pdf.CellFormat(hSep, 6, "", "", 0, "", false, 0, "")
		pdf.CellFormat(20, 6, fmt.Sprintf("%.2f", g.Costs.Total), "", 0, "R", false, 0, "")
		pdf.Ln(midBr)
		pdf.SetDrawColor(200, 200, 200)
		pdf.CellFormat(0, 6, "", "B", 2, "", false, 0, "")
	}
}

// otherList method
func (p *PDF) otherList() {

	pdf := p.pdf
	q := p.q

	if len(q.Items.Other) <= 0 {
		return
	}
	var maxLen int
	const descriptionCellW float64 = 65
	const optsCellW float64 = 52

	for _, g := range q.Items.Other {
		maxLen = 20
		locStr := g.Specs.Location
		if len(locStr) > maxLen {
			locStr = locStr[0:maxLen-3] + "..."
		}

		pdf.Ln(topBr)
		pdf.SetFont("Arial", "", pfSize)
		pdf.CellFormat(5, 6, fmt.Sprintf("%d", g.Qty), "", 0, "", false, 0, "")
		pdf.CellFormat(26, 6, fmt.Sprintf("%s", strings.Join(g.Rooms, ", ")), "", 0, "", false, 0, "")
		xPos := pdf.GetX() + descriptionCellW
		yPos := pdf.GetY()
		pdf.MultiCell(descriptionCellW, 4.5, fmt.Sprintf("%s", g.Description), "", "L", false)
		pdf.MoveTo(xPos, yPos)
		pdf.CellFormat(hSep, 6, "", "", 0, "", false, 0, "")
		xPos = pdf.GetX() + optsCellW
		yPos = pdf.GetY()
		pdf.MultiCell(optsCellW, 4.5, fmt.Sprintf("%s", g.Specs.Options), "", "L", false)
		pdf.MoveTo(xPos, yPos)
		pdf.CellFormat(hSep, 6, "", "", 0, "", false, 0, "")
		pdf.CellFormat(20, 6, fmt.Sprintf("%.2f", g.Costs.Unit), "", 0, "R", false, 0, "")
		pdf.CellFormat(hSep, 6, "", "", 0, "", false, 0, "")
		pdf.CellFormat(20, 6, fmt.Sprintf("%.2f", g.Costs.Total), "", 0, "R", false, 0, "")

		pdf.Ln(midBr)
		pdf.SetDrawColor(200, 200, 200)
		pdf.CellFormat(0, 6, "", "B", 2, "", false, 0, "")
	}
}

// summary method
/* func (p *PDF) summary() {

	pdf := p.pdf
	q := p.q

	discount := false
	havePayments := false
	if q.Discount.Total > 0 && q.Invoiced == false {
		discount = true
	}
	if q.Invoiced && q.Price.Payments > 0 {
		havePayments = true
	}
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 15)
	pdf.CellFormat(14, 8, "Summary", "", 2, "", false, 0, "")
	pdf.SetDrawColor(0, 0, 0)
	pdf.CellFormat(0, 0, "", "B", 2, "", false, 0, "")
	pdf.Ln(6)
	pdf.SetFont("Arial", "I", 11)
	pdf.CellFormat(0, 6, "Features:", "", 2, "", false, 0, "")
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 13)
	pdf.MultiCell(0, 6, string(q.Features), "", "", false)
	pdf.Ln(8)
	pdf.CellFormat(0, 5, "", "B", 2, "", false, 0, "")
	pdf.Ln(2)

	if discount {
		pdf.CellFormat(115, 6, "", "", 0, "", false, 0, "")
		xPos := pdf.GetX() + 33
		yPos := pdf.GetY() + 2.75
		pdf.CellFormat(45, 6, "Total", "", 0, "R", false, 0, "")
		pdf.Line(xPos, yPos, xPos+47, yPos)
		pdf.CellFormat(35, 6, formatMoney(q.ItemCosts.Subtotal, ""), "", 2, "R", false, 0, "")
		pdf.Ln(1)
	}
	pdf.CellFormat(115, 6, "", "", 0, "", false, 0, "")
	pdf.CellFormat(45, 6, "Subtotal", "", 0, "R", false, 0, "")
	pdf.CellFormat(35, 6, formatMoney(q.Price.Subtotal, ""), "", 2, "R", false, 0, "")
	pdf.Ln(1)
	pdf.CellFormat(115, 6, "", "", 0, "", false, 0, "")
	pdf.CellFormat(45, 6, "Tax", "", 0, "R", false, 0, "")
	pdf.CellFormat(35, 6, formatMoney(q.Price.Tax, ""), "", 2, "R", false, 0, "")
	pdf.Ln(1)
	pdf.CellFormat(0, 0, "", "B", 2, "", false, 0, "")
	pdf.Ln(2)
	pdf.CellFormat(115, 7, "", "", 0, "", false, 0, "")
	pdf.SetFont("Arial", "B", 13)
	pdf.CellFormat(45, 7, "Total Invoice", "", 0, "R", false, 0, "")
	pdf.CellFormat(35, 7, formatMoney(q.Price.Total, "$"), "", 2, "R", false, 0, "")
	pdf.Ln(1)
	pdf.CellFormat(0, 0, "", "B", 2, "", false, 0, "")

	if !havePayments {
		pdf.CellFormat(0, .75, "", "B", 2, "", false, 0, "")
	}
	if havePayments {
		pdf.Ln(4)
		pdf.SetFont("Arial", "", 13)
		pdf.CellFormat(14, 8, "Payments", "", 2, "", false, 0, "")
		pdf.SetFont("Arial", "", 13)
		pdf.Ln(1)

		for _, p := range q.Payments {
			pdf.CellFormat(115, 6, "", "", 0, "", false, 0, "")
			pdf.CellFormat(45, 6, p.Date.Format(dateFormat), "", 0, "R", false, 0, "")
			pdf.CellFormat(35, 6, formatMoney(p.Amount, ""), "", 2, "R", false, 0, "")
			pdf.Ln(1)
		}

		pdf.CellFormat(0, 0, "", "B", 2, "", false, 0, "")
		pdf.Ln(2)
		pdf.CellFormat(115, 6, "", "", 0, "", false, 0, "")
		pdf.CellFormat(45, 6, "Total Paid", "", 0, "R", false, 0, "")
		pdf.CellFormat(35, 6, formatMoney(q.Price.Payments, ""), "", 2, "R", false, 0, "")
		pdf.Ln(1)
		pdf.SetFont("Arial", "B", 13)
		pdf.CellFormat(115, 6, "", "", 0, "", false, 0, "")
		pdf.CellFormat(45, 6, "Balance Due", "", 0, "R", false, 0, "")
		pdf.CellFormat(35, 6, formatMoney(q.Price.Outstanding, ""), "", 2, "R", false, 0, "")
		pdf.Ln(1)
		pdf.CellFormat(0, 0, "", "B", 2, "", false, 0, "")
		pdf.CellFormat(0, .75, "", "B", 2, "", false, 0, "")
	}
} */

// ================================ Helper Methods

func setNewLines(name interface{}) string {

	str := name.(string)
	if len(str) <= newLnMaxLen {
		return str
	}

	pcs := strings.Split(str, ",")
	strSl := make([]string, len(pcs))
	for i := 0; i < len(pcs); i++ {
		strSl[i] = strings.TrimSpace(pcs[i])
	}
	retVal := strings.Join(strSl, "\n")

	return retVal
}

func formatMoney(num float64, prefix string) string {
	moneyFormat := "#,###.##"
	return fmt.Sprintf("%s%s", prefix, humanize.FormatFloat(moneyFormat, num))
}
