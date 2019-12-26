package main

import (
	"time"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	err := generateCertificate("certificate.pdf", "Preslav Mihaylov", "Jonathan Calhoun", time.Now())
	if err != nil {
		panic(err)
	}
}

func generateCertificate(certificate, student, instructor string, date time.Time) error {
	pdf := initPDF()
	generateImages(pdf)
	generateCertificateText(pdf, student)
	generateCertificateDateAndInstructor(pdf, instructor, date)

	return savePDF(pdf, certificate)
}

func initPDF() *gofpdf.Fpdf {
	pdf := gofpdf.New("L", "mm", "A5", "")
	pdf.AddPage()

	return pdf
}

func generateImages(pdf *gofpdf.Fpdf) {
	pageWidth, pageHeight := pdf.GetPageSize()
	pdf.ImageOptions("images/header.png",
		0, 0, pageWidth, 0, false,
		gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	pdf.ImageOptions("images/footer.png",
		0, pageHeight-20, pageWidth, 0, false,
		gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	pdf.ImageOptions("images/gopher.png",
		pageWidth/2-15, pageHeight-50, 30, 0, false,
		gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
}

func generateCertificateText(pdf *gofpdf.Fpdf, student string) {
	pageWidth, _ := pdf.GetPageSize()
	pdf.SetFont("Times", "B", 34)
	pdf.Ln(10)
	pdf.CellFormat(0, 20, "Certificate of Completion", "", 1, "C", false, 0, "")

	pdf.SetFont("Times", "", 18)
	pdf.CellFormat(0, 20, "This certificate is awarded to", "", 1, "C", false, 0, "")

	pdf.SetFont("Times", "", 28)
	pdf.CellFormat(0, 10, student, "", 1, "C", false, 0, "")

	pdf.SetFont("Times", "", 18)
	text := "For successfully completing all twenty programming exercises " +
		"in the Gophercises Go programming course."
	lines := pdf.SplitText(text, pageWidth)
	pdf.CellFormat(0, 20, lines[0], "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 0, lines[1], "", 1, "C", false, 0, "")
}

func generateCertificateDateAndInstructor(pdf *gofpdf.Fpdf, instructor string, date time.Time) {
	pageWidth, pageHeight := pdf.GetPageSize()
	pdf.SetTextColor(105, 105, 105)

	pdf.MoveTo(20, pageHeight-40)
	pdf.SetFont("Times", "", 14)
	pdf.CellFormat(60, 10, date.Format("02/01/2006"), "B", 1, "C", false, 0, "")

	pdf.MoveTo(pageWidth-20-60, pageHeight-40)
	pdf.CellFormat(60, 10, instructor, "B", 1, "C", false, 0, "")

	pdf.MoveTo(20, pageHeight-30)
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(60, 7, "Date", "", 1, "C", false, 0, "")

	pdf.MoveTo(pageWidth-20-60, pageHeight-30)
	pdf.CellFormat(60, 7, "Instructor", "", 1, "C", false, 0, "")
}

func savePDF(pdf *gofpdf.Fpdf, certificate string) error {
	return pdf.OutputFileAndClose(certificate)
}
