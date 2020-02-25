package pdf

import "fmt"

// quoteSummary method
func (p *PDF) quoteSummary() {

	pdf := p.pdf
	q := p.q

	discount := false
	if q.Discount.Total > 0 {
		discount = true
	}
	discountDesc := false
	if q.Discount.Description != "" {
		discountDesc = true
	}
	pdf.AddPage()
	p.features()

	if discount {
		pdf.CellFormat(115, 6, "", "", 0, "", false, 0, "")
		xPos := pdf.GetX() + 33
		yPos := pdf.GetY() + 2.75
		pdf.CellFormat(45, 6, "Total", "", 0, "R", false, 0, "")
		pdf.Line(xPos, yPos, xPos+47, yPos)
		pdf.CellFormat(35, 6, formatMoney(q.ItemCosts.Subtotal, ""), "", 2, "R", false, 0, "")
		pdf.Ln(1)
	}
	p.totals()

	if discountDesc {
		pdf.Ln(2)
		pdf.SetFont("Arial", "I", 11)
		pdf.CellFormat(40, 6, "Discount Description:", "", 0, "", false, 0, "")
		pdf.SetFont("Arial", "", 12)
		pdf.CellFormat(0, 6, fmt.Sprintf("%s", q.Discount.Description), "", 1, "", false, 0, "")
	}

	// Payment methods
	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(40, 6, "Payment Methods: Cash, Cheque, or eTransfer", "", 0, "", false, 0, "")
}

// invoiceSummary method
func (p *PDF) invoiceSummary() {

	pdf := p.pdf
	q := p.q

	havePayments := false
	if q.Price.Payments > 0 {
		havePayments = true
	}

	pdf.AddPage()
	p.features()

	p.totals()

	if !havePayments {
		pdf.CellFormat(0, .75, "", "B", 2, "", false, 0, "")
	}
	if havePayments {
		pdf.Ln(4)
		pdf.SetFont("Arial", "I", 11)
		pdf.CellFormat(14, 8, "Payments", "", 2, "", false, 0, "")
		pdf.SetFont("Arial", "", 13)
		// pdf.Ln(1)

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
}

func (p *PDF) features() {

	pdf := p.pdf
	q := p.q

	pdf.SetFont("Arial", "B", 15)
	pdf.CellFormat(14, 8, "Summary", "", 2, "", false, 0, "")
	pdf.SetDrawColor(0, 0, 0)
	pdf.CellFormat(0, 0, "", "B", 2, "", false, 0, "")
	pdf.Ln(6)
	pdf.SetFont("Arial", "I", 11)
	pdf.CellFormat(0, 6, "Features", "", 2, "", false, 0, "")
	pdf.Ln(2)
	pdf.SetFont("Arial", "", 13)
	pdf.MultiCell(0, 6, string(q.Features), "", "", false)
	pdf.Ln(8)
	pdf.CellFormat(0, 5, "", "B", 2, "", false, 0, "")
	pdf.Ln(2)
}

func (p *PDF) totals() {

	pdf := p.pdf
	q := p.q
	invoiced := q.Invoiced
	totalLabel := "Total Amount"
	if invoiced {
		totalLabel = "Total Invoice"
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
	pdf.CellFormat(45, 7, totalLabel, "", 0, "R", false, 0, "")
	pdf.CellFormat(35, 7, formatMoney(q.Price.Total, ""), "", 2, "R", false, 0, "")
	pdf.Ln(1)
	pdf.CellFormat(0, 0, "", "B", 2, "", false, 0, "")
}
