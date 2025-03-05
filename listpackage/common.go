package listpackage

import (
	"errors"

	"github.com/navigacontentlab/docconverter/newsml"

	"github.com/navigacontentlab/navigadoc/doc"
)

type Product struct {
	UUID string `xml:"uuid,attr,omitempty"`
	Text string `xml:",chardata"`
}

func AsDocument(v interface{}, opts *newsml.Options) (*doc.Document, error) {
	if v == nil {
		return nil, errors.New("nil input")
	}

	document := &doc.Document{}

	switch i := v.(type) {
	case *List:
		err := i.toDoc(document, opts)
		if err != nil {
			return nil, err
		}
	case *Package:
		err := i.toDoc(document, opts)
		if err != nil {
			return nil, err
		}
	}

	return document, nil
}
