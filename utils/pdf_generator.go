package utils

import (
	"bytes"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func HtmlToPDF(htmlContent string) ([]byte, error) {
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	pdfg.AddPage(wkhtml.NewPageReader(bytes.NewReader([]byte(htmlContent))))

	pdfg.PageSize.Set("A4")
	pdfg.MarginTop.Set(10)
	pdfg.MarginBottom.Set(10)

	err = pdfg.Create()
	if err != nil {
		return nil, err
	}

	return pdfg.Bytes(), nil
}
