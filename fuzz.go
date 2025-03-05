//go:build gofuzz
// +build gofuzz

package docconverter

import (
	"bytes"
	"encoding/xml"

	"github.com/navigacontentlab/docconverter/newsml"
)

func Fuzz(data []byte) int {
	reader := bytes.NewReader(data)
	dec := internal.NewXMLDecoder(reader)

	var article newsml.NewsItem
	if err := dec.Decode(&article); err != nil {
		return 0
	}

	doc, err := NewsItemToDoc(&article, GetDefaultConversionOptions())
	if err != nil {
		if doc != nil {
			panic("document should be nil on error")
		}
		return 0
	}

	item, err := DocToNewsItem(doc, GetDefaultConversionOptions())
	if err != nil {
		if item != nil {
			panic("newsitem should be nil on error")
		}
		return 0
	}

	return 1
}
