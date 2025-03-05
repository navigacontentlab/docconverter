package docconverter_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"
	"unicode"

	"github.com/navigacontentlab/docconverter/internal"
	"github.com/navigacontentlab/docconverter/newsml"

	"github.com/navigacontentlab/navigadoc/doc"

	"github.com/Infomaker/etree"
	"github.com/navigacontentlab/docconverter"
)

func TestValidateNewsMLDates(t *testing.T) {
	var custom newsml.Options

	file, err := os.Open("./testdata/custom-config.json")
	if err != nil {
		t.Fatal("unable to open ./testdata/dates-config.json")
	}

	dec := json.NewDecoder(file)
	err = dec.Decode(&custom)
	if err != nil {
		t.Fatal("unable to decode ./testdata/dates-config.json")
	}

	tests := []TestData{
		{xml: "./testdata/good-dates.xml", expectError: false},
		{xml: "./testdata/bad-dates.xml", expectError: true},
		{xml: "./examples/full-article.xml", expectError: false},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.xml, func(t *testing.T) {
			xmlDoc := internal.NewEtreeDoc()
			err := xmlDoc.ReadFromFile(test.xml)
			if err != nil {
				t.Fatal(err)
			}
			err = docconverter.ValidateNewsMLDates(xmlDoc, custom.DateElements)
			if !test.expectError && err != nil {
				t.Error(err)
			} else if test.expectError && err == nil {
				t.Error("error was expected")
			}
			err = docconverter.ValidateNewsMLDates(xmlDoc, newsml.DefaultOptions().DateElements)
			if !test.expectError && err != nil {
				t.Error(err)
			} else if test.expectError && err == nil {
				t.Error("error was expected")
			}
		})
	}
}

// Deprecated: UUID handling in OpenContent to be case-insensitive
func TestValidateAndLowercaseNewsMLUUIDs(t *testing.T) {
	tests := []TestData{
		{xml: "./examples/full-article.xml", expectError: false},
		{xml: "./testdata/uuids-invalid.xml", expectError: true},
		{xml: "./testdata/uuids-uppercase.xml", expectError: false},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.xml, func(t *testing.T) {
			xmlDoc := internal.NewEtreeDoc()
			err := xmlDoc.ReadFromFile(test.xml)
			if err != nil {
				t.Fatal(err)
			}
			err = docconverter.FixNewsMLUUIDsAndNamespaces(xmlDoc)
			if !test.expectError && err != nil {
				t.Error(err)
			} else if test.expectError && err == nil {
				t.Error("error was expected")
			}

			if test.expectError && err != nil {
				if errors.Is(err, docconverter.InvalidArgumentError{}) {
					t.Logf("InvalidArgumentError %s", err)
				} else if errors.Is(err, docconverter.RequiredArgumentError{}) {
					t.Logf("RequiredArgumentError %s", err)
				}
			}

			// Test that all UUIDs are lowercase
			err = docconverter.WalkXMLDocumentElements(xmlDoc, nil, func(element *etree.Element, args ...interface{}) error {
				uuidAttr := element.SelectAttr("uuid")
				if element.Tag != "" && uuidAttr == nil {
					return nil
				}

				var uuidValue string
				if uuidAttr != nil {
					uuidValue = uuidAttr.Value
				} else {
					uuidValue = element.Text()
				}

				for _, r := range uuidValue {
					if unicode.IsUpper(r) {
						return fmt.Errorf("uuid contains uppercase: %s", uuidValue)
					}
				}

				return nil
			})
			if err != nil {
				t.Errorf("%s", err)
			}
		})
	}
}

func TestValidateNewsMLUUIDs(t *testing.T) {
	tests := []TestData{
		{xml: "./examples/full-article.xml", expectError: false},
		{xml: "./testdata/uuids-invalid.xml", expectError: true},
		{xml: "./testdata/uuids-uppercase.xml", expectError: false, customConfig: "uc"},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.xml, func(t *testing.T) {
			xmlDoc := internal.NewEtreeDoc()
			err := xmlDoc.ReadFromFile(test.xml)
			if err != nil {
				t.Fatal(err)
			}
			err = docconverter.ValidateNewsMLUUIDsAndFixNamespaces(xmlDoc)
			if !test.expectError && err != nil {
				t.Error(err)
			} else if test.expectError && err == nil {
				t.Error("error was expected")
			}

			if test.expectError && err != nil {
				if errors.Is(err, docconverter.InvalidArgumentError{}) {
					t.Logf("InvalidArgumentError %s", err)
				} else if errors.Is(err, docconverter.RequiredArgumentError{}) {
					t.Logf("RequiredArgumentError %s", err)
				}
			}

			if test.customConfig == "uc" {
				err = docconverter.WalkXMLDocumentElements(xmlDoc, nil, func(element *etree.Element, args ...interface{}) error {
					uuidAttr := element.SelectAttr("uuid")
					if element.Tag != "" && uuidAttr == nil {
						return nil
					}

					var uuidValue string
					if uuidAttr != nil {
						uuidValue = uuidAttr.Value
					} else {
						uuidValue = element.Text()
					}

					for _, r := range uuidValue {
						if unicode.IsLower(r) {
							return fmt.Errorf("uppercase uuid expected: %s", uuidValue)
						}
					}

					return nil
				})
				if err != nil {
					t.Errorf("%s", err)
				}
			}
		})
	}
}

func TestValidateDocumentDates(t *testing.T) {
	var dc newsml.DateConfig
	file, err := os.Open("./testdata/dates-config.json")
	if err != nil {
		t.Fatal("unable to open ./testdata/dates-config.json")
	}

	dec := json.NewDecoder(file)
	err = dec.Decode(&dc)
	if err != nil {
		t.Fatal("unable to decode ./testdata/dates-config.json")
	}

	tests := []TestData{
		{json: "./testdata/text.json", expectError: false},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.xml, func(t *testing.T) {
			testData, err := os.ReadFile(test.json)
			must(t, err, "could not open testfile")

			var document doc.Document
			err = json.Unmarshal(testData, &document)
			must(t, err, "could not unmarshal doc")

			err = docconverter.ValidateDocumentDates(&document, dc)
			if !test.expectError && err != nil {
				t.Error(err)
			} else if test.expectError && err == nil {
				t.Error("error was expected")
			}
		})
	}
}

func TestValidateDate(t *testing.T) {
	var custom newsml.Options

	file, err := os.Open("./testdata/custom-config.json")
	if err != nil {
		t.Fatal("unable to open ./testdata/dates-config.json")
	}

	dec := json.NewDecoder(file)
	err = dec.Decode(&custom)
	if err != nil {
		t.Fatal("unable to decode ./testdata/dates-config.json")
	}

	date := "2021-06-01T04:42:00-05:00"
	err = docconverter.ValidateDate(custom.DateElements, newsml.Block, "firstCreated", date)
	if err != nil {
		t.Fatalf("date %s should have errored", date)
	}

	date = "2021-06-01T04:42:00Z"
	err = docconverter.ValidateDate(custom.DateElements, newsml.Block, "firstCreated", date)
	if err != nil {
		t.Fatalf("date %s should have errored", date)
	}

	date = ""
	err = docconverter.ValidateDate(custom.DateElements, newsml.Block, "firstCreated", date)
	if err == nil {
		t.Fatalf("date %s missing UTC offset should have errored", date)
	}
	// t.Logf("Expected error: %s\n", err)

	date = "2021-99-99T04:42:00Z"
	err = docconverter.ValidateDate(custom.DateElements, newsml.Block, "firstCreated", date)
	if err == nil {
		t.Fatalf("date %s missing UTC offset should have errored", date)
	}
	// t.Logf("Expected error: %s\n", err)

	date = "2021-06-01T04:42:00"
	err = docconverter.ValidateDate(custom.DateElements, newsml.Block, "firstCreated", date)
	if err == nil {
		t.Fatalf("date %s missing UTC offset should have errored", date)
	}
	// t.Logf("Expected error: %s\n", err)

	date = "2021-06-01T04:42:00-16:00"
	err = docconverter.ValidateDate(custom.DateElements, newsml.Block, "firstCreated", date)
	if err == nil {
		t.Fatalf("date %s invalid UTC offset should have errored", date)
	}
	// t.Logf("Expected error: %s\n", err)
}
