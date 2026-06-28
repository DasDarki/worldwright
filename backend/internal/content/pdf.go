package content

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/yuin/goldmark"

	pdfrenderer "github.com/stephenafamo/goldmark-pdf"
)

// PDFEngine renders Worldwright markdown to a PDF document. It's deliberately
// process-wide and lazy: building the underlying goldmark+PDF pipeline takes
// some milliseconds, and the result is concurrency-safe.
type PDFEngine struct {
	once sync.Once
	gm   goldmark.Markdown
	err  error
}

var defaultPDF = &PDFEngine{}

func init() {
	defaultPDF.init()
}

func (e *PDFEngine) init() {
	e.once.Do(func() {
		// Default config: built-in fonts, A4, sensible heading sizes. Avoids
		// pulling in our custom Fraunces/Garamond fonts to keep things
		// dependency-free; the PDF still looks neutral and readable.
		renderer := pdfrenderer.New(
			pdfrenderer.WithContext(context.Background()),
		)
		e.gm = goldmark.New(goldmark.WithRenderer(renderer))
	})
}

// MarkdownToPDF converts the given markdown string into PDF bytes. Returns
// (nil, err) when goldmark or the PDF renderer fail.
func MarkdownToPDF(md string) ([]byte, error) {
	defaultPDF.init()
	if defaultPDF.err != nil {
		return nil, defaultPDF.err
	}
	var buf bytes.Buffer
	if err := defaultPDF.gm.Convert([]byte(md), &buf); err != nil {
		return nil, fmt.Errorf("render pdf: %w", err)
	}
	return buf.Bytes(), nil
}
