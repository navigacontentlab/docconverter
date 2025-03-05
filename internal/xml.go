package internal

import (
	"bytes"
	"encoding/xml"
	"io"

	"github.com/Infomaker/etree"
)

func NewEtreeDoc() *etree.Document {
	doc := etree.NewDocument()

	doc.ReadSettings.Entity = xml.HTMLEntity

	return doc
}

func NewXMLDecoder(r io.Reader) *xml.Decoder {
	dec := xml.NewDecoder(r)

	dec.Entity = xml.HTMLEntity

	return dec
}

func UnmarshalXML(d []byte, o interface{}) error {
	dec := NewXMLDecoder(bytes.NewReader(d))

	return dec.Decode(o)
}
