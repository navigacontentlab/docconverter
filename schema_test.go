package docconverter_test

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/navigacontentlab/docconverter"
	"github.com/navigacontentlab/docconverter/internal"
	"github.com/navigacontentlab/docconverter/newsml"
	"github.com/navigacontentlab/navigadoc"
	"github.com/navigacontentlab/navigadoc/doc"
)

func TestValidateNavigadoc(t *testing.T) {
	testData, err := os.ReadFile(configFile)
	must(t, err, "could not open file")
	customConfig := string(testData)

	var tests = []TestData{
		{xml: "examples/full-article.xml", customConfig: customConfig, expectError: false},
		{xml: "examples/full-picture.xml", customConfig: customConfig, expectError: false},
		{xml: "examples/article-template.xml", customConfig: customConfig, expectError: false},
		{xml: "examples/example-pdf.xml", customConfig: customConfig, expectError: false},
		{xml: "examples/example-text.xml", customConfig: customConfig, expectError: false},
		// FIXME? example-wire.xml winds up with properties missing value, so removed from required list
		{xml: "examples/example-wire.xml", customConfig: customConfig, expectError: false},
		{json: "examples/cca-example.json", expectError: false},
		{json: "testdata/empty-blocks-example.json", expectError: false},
		{json: "testdata/empty-document-example.json", expectError: true},
		{json: "testdata/custom-config.json", expectError: true},
		{json: "testdata/uuids-uppercase.json", expectError: false},
	}

	opts := newsml.DefaultOptions()
	for i := range tests {
		test := tests[i]
		if test.xml != "" {
			t.Run(test.xml, func(t *testing.T) {
				testXML, err := os.ReadFile(test.xml)
				must(t, err, "could not read file")

				var original newsml.NewsItem
				err = internal.UnmarshalXML(testXML, &original)
				must(t, err, "could not unmarshal file")

				if test.customConfig != "" {
					customConfig := newsml.Options{}

					err = json.Unmarshal([]byte(test.customConfig), &customConfig)
					must(t, err, "could not unmarshal file")

					docconverter.MergeOptions(&opts, &customConfig)
				}

				navigaDoc, err := docconverter.NewsItemToDoc(&original, &opts)
				if test.expectError && err != nil {
					return
				}
				must(t, err, "failed NewsItemToDoc")

				docBytes, err := json.MarshalIndent(navigaDoc, "", " ")
				must(t, err, "failed marshal")

				docText := string(docBytes)
				errs, err := navigadoc.ValidateNavigadocJSON(docText)
				if err != nil || len(errs) > 0 || testing.Verbose() {
					if err != nil || len(errs) > 0 {
						t.Errorf("unexpected errors %s", test.xml)
					}
					dumpDocs(t, navigaDoc, nil, testXML, nil, errs...)
				}

				// Test the UUID pattern
				docUUID := navigaDoc.UUID
				navigaDoc.UUID = "JKjhsnajuSH"
				docBytes, err = json.MarshalIndent(navigaDoc, "", " ")
				must(t, err, "failed marshal")

				test.expectError = true
				errs, err = navigadoc.ValidateNavigadocJSON(string(docBytes))
				if err == nil && len(errs) == 0 {
					t.Errorf("expected error on bad uuid %s", test.xml)
				}

				navigaDoc.UUID = docUUID

				// Test the date-time
				rgxp := regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}")
				dates := rgxp.FindAllString(docText, -1)
				badJSON := docText
				for _, d := range dates {
					badJSON = strings.ReplaceAll(badJSON, d, "xxxx-xx:xx:xxT")
					break
				}
				errs, err = navigadoc.ValidateNavigadocJSON(badJSON)
				if err == nil && len(errs) == 0 {
					t.Errorf("expected error on bad date %s", test.xml)
				}

				// Test required ID
				rgxp = regexp.MustCompile("\"type\".*?:.*?\".*?\".*,")
				types := rgxp.FindAllString(docText, -1)
				badJSON = docText
				for _, t := range types {
					badJSON = strings.ReplaceAll(badJSON, t, "")
					break
				}
				errs, err = navigadoc.ValidateNavigadocJSON(badJSON)
				if err == nil && len(errs) == 0 {
					t.Errorf("expected error on required id %s", test.xml)
				}
			})
		}

		if test.json != "" {
			t.Run(test.json, func(t *testing.T) {
				testJSON, err := os.ReadFile(test.json)
				must(t, err, "could not read file")

				errs, err := navigadoc.ValidateNavigadocJSON(string(testJSON))
				if !test.expectError && err != nil || len(errs) > 0 {
					navigaDoc := doc.Document{}
					err := json.Unmarshal(testJSON, &navigaDoc)
					if err != nil {
						errs = append(errs, fmt.Errorf("error unmarshaling %s", test.json))
						navigaDoc = doc.Document{}
					}
					dumpDocs(t, &navigaDoc, nil, nil, nil, errs...)
				}
			})
		}
	}
}
