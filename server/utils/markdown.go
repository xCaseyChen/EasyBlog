package utils

import (
	"bytes"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
)

func MarkdownToHTML(mdString string) (string, error) {
	var htmlContent bytes.Buffer
	highlightGoldmark := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
		),
	)
	if err := highlightGoldmark.Convert([]byte(mdString), &htmlContent); err != nil {
		return "", err
	}
	return htmlContent.String(), nil
}
