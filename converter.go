// Package docconverter provides helper functions for converting
// between CCA Doc and Infomaker NewsML (IMNML) format
package docconverter

import (
	"bufio"
	"bytes"

	// embed navigadoc schema
	_ "embed"
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/navigacontentlab/docconverter/internal"
	"github.com/navigacontentlab/docconverter/listpackage"
	"github.com/navigacontentlab/docconverter/newsml"
	"github.com/navigacontentlab/navigadoc"
	"github.com/navigacontentlab/navigadoc/doc"
)

// NavigaDocSchema embedded schema
//
//go:embed schema/navigadoc-schema.json
var NavigaDocSchema string

func DocToConceptItem(document *doc.Document, opts *newsml.Options) (*newsml.ConceptItem, error) {
	if document == nil {
		return nil, ErrEmptyDoc
	}

	err := navigadoc.CheckForEmptyBlocks(document)
	if err != nil {
		return nil, err
	}

	if opts == nil {
		defaults := newsml.DefaultOptions()
		opts = &defaults
	}

	conceptItem, err := newsml.AsConceptItem(document, opts)
	if err != nil {
		return nil, err
	}
	return conceptItem, nil
}

// DocToNewsItem converts a NavigaDoc to a NewsItem
func DocToNewsItem(document *doc.Document, opts *newsml.Options) (*newsml.NewsItem, error) {
	if document == nil {
		return nil, ErrEmptyDoc
	}

	err := navigadoc.CheckForEmptyBlocks(document)
	if err != nil {
		return nil, err
	}

	if opts == nil {
		defaults := newsml.DefaultOptions()
		opts = &defaults
	}

	newsitem, err := newsml.AsNewsItem(document, opts)
	if err != nil {
		return nil, err
	}
	return newsitem, nil
}

func ConceptItemToDoc(conceptItem *newsml.ConceptItem, opts *newsml.Options) (*doc.Document, error) {
	if conceptItem == nil {
		return nil, ErrEmptyConcept
	}

	if opts == nil {
		defaults := newsml.DefaultOptions()
		opts = &defaults
	}

	document, err := newsml.AsDocument(conceptItem, opts)
	if err != nil {
		return nil, err
	}

	return document, nil
}

// NewsItemToDoc converts a NewsItem (Infomaker NewsML) document to a NavigaDoc
func NewsItemToDoc(newsItem *newsml.NewsItem, opts *newsml.Options) (*doc.Document, error) {
	if newsItem == nil {
		return nil, ErrEmptyNewsItem
	}

	if opts == nil {
		defaults := newsml.DefaultOptions()
		opts = &defaults
	}

	document, err := newsml.AsDocument(newsItem, opts)
	if err != nil {
		return nil, err
	}

	if len(document.Properties) == 0 {
		document.Properties = nil
	}

	sort.SliceStable(document.Properties, func(i, j int) bool {
		return document.Properties[i].Name < document.Properties[j].Name
	})

	return document, nil
}

func ListToDoc(list *listpackage.List, opts *newsml.Options) (*doc.Document, error) {
	if list == nil {
		return nil, ErrEmptyList
	}

	if opts == nil {
		defaults := newsml.DefaultOptions()
		opts = &defaults
	}

	document, err := listpackage.AsDocument(list, opts)
	if err != nil {
		return nil, err
	}

	if document.Properties != nil && len(document.Properties) == 0 {
		document.Properties = nil
	}

	sort.SliceStable(document.Properties, func(i, j int) bool {
		return document.Properties[i].Name < document.Properties[j].Name
	})

	return document, nil
}

func DocToList(document *doc.Document, opts *newsml.Options) (*listpackage.List, error) {
	if document == nil {
		return nil, ErrEmptyDoc
	}

	err := navigadoc.CheckForEmptyBlocks(document)
	if err != nil {
		return nil, err
	}

	if opts == nil {
		defaults := newsml.DefaultOptions()
		opts = &defaults
	}

	list, err := listpackage.AsList(document, opts)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func AssignmentToDoc(assignment *newsml.Assignment, opts *newsml.Options) (*doc.Document, error) {
	if assignment == nil {
		return nil, ErrEmptyAssignment
	}

	if opts == nil {
		defaults := newsml.DefaultOptions()
		opts = &defaults
	}

	document, err := assignment.AssignmentToDoc(assignment, opts)
	if err != nil {
		return nil, err
	}

	if document.Properties != nil && len(document.Properties) == 0 {
		document.Properties = nil
	}

	sort.SliceStable(document.Properties, func(i, j int) bool {
		return document.Properties[i].Name < document.Properties[j].Name
	})

	return document, nil
}

func DocToAssignment(document *doc.Document, opts *newsml.Options) (*newsml.Assignment, error) {
	if document == nil {
		return nil, ErrEmptyDoc
	}

	err := navigadoc.CheckForEmptyBlocks(document)
	if err != nil {
		return nil, err
	}

	assignment := &newsml.Assignment{}

	if opts == nil {
		defaults := newsml.DefaultOptions()
		opts = &defaults
	}

	err = assignment.DocToAssignment(document, opts)
	if err != nil {
		return nil, err
	}

	return assignment, nil
}

func PlanningItemToDoc(planningItem *newsml.PlanningItem, opts *newsml.Options) (*doc.Document, error) {
	if planningItem == nil {
		return nil, ErrEmptyPlanningItem
	}

	if opts == nil {
		defaults := newsml.DefaultOptions()
		opts = &defaults
	}

	document, err := planningItem.PlanningItemToDoc(planningItem, opts)
	if err != nil {
		return nil, err
	}

	if len(document.Properties) == 0 {
		document.Properties = nil
	}

	sort.SliceStable(document.Properties, func(i, j int) bool {
		return document.Properties[i].Name < document.Properties[j].Name
	})

	return document, nil
}

func DocToPlanningItem(document *doc.Document, opts *newsml.Options) (*newsml.PlanningItem, error) {
	if document == nil {
		return nil, ErrEmptyDoc
	}

	err := navigadoc.CheckForEmptyBlocks(document)
	if err != nil {
		return nil, err
	}

	planningItem := &newsml.PlanningItem{}

	if opts == nil {
		defaults := newsml.DefaultOptions()
		opts = &defaults
	}

	err = planningItem.DocToPlanningItem(document, opts)
	if err != nil {
		return nil, err
	}

	return planningItem, nil
}

func PackageToDoc(pkg *listpackage.Package, opts *newsml.Options) (*doc.Document, error) {
	if pkg == nil {
		return nil, ErrEmptyPackage
	}

	if opts == nil {
		defaults := newsml.DefaultOptions()
		opts = &defaults
	}

	document, err := listpackage.AsDocument(pkg, opts)
	if err != nil {
		return nil, err
	}

	if len(document.Properties) == 0 {
		document.Properties = nil
	}

	sort.SliceStable(document.Properties, func(i, j int) bool {
		return document.Properties[i].Name < document.Properties[j].Name
	})

	return document, nil
}

func DocToPackage(document *doc.Document, opts *newsml.Options) (*listpackage.Package, error) {
	if document == nil {
		return nil, ErrEmptyDoc
	}

	err := navigadoc.CheckForEmptyBlocks(document)
	if err != nil {
		return nil, err
	}

	if opts == nil {
		defaults := newsml.DefaultOptions()
		opts = &defaults
	}

	pkg, err := listpackage.AsPackage(document, opts)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}

func DocToItemWithDefault(document *doc.Document, opts *newsml.Options, defaultCase func(document *doc.Document, opts *newsml.Options) (interface{}, error)) (interface{}, error) {
	switch document.Type {
	case "x-im/article":
		fallthrough
	case "x-im/image":
		fallthrough
	case "x-im/graphic":
		fallthrough
	case "x-im/pdf":
		return DocToNewsItem(document, opts)
	case "x-im/author":
		fallthrough
	case "x-im/category":
		fallthrough
	case "x-im/channel":
		fallthrough
	case "x-im/content-profile":
		fallthrough
	case "x-im/event":
		fallthrough
	case "x-im/organisation":
		fallthrough
	case "x-im/person":
		fallthrough
	case "x-im/place":
		fallthrough
	case "x-im/section":
		fallthrough
	case "x-im/story":
		fallthrough
	case "x-im/topic":
		return DocToConceptItem(document, opts)
	case "x-im/newscoverage":
		return DocToPlanningItem(document, opts)
	case "x-im/assignment":
		return DocToAssignment(document, opts)
	case "x-im/package":
		return DocToPackage(document, opts)
	case "x-im/list":
		return DocToList(document, opts)
	default:
		return defaultCase(document, opts)
	}
}

func DocToItem(document *doc.Document, opts *newsml.Options) (interface{}, error) {
	return DocToItemWithDefault(document, opts, func(document *doc.Document, opts *newsml.Options) (interface{}, error) {
		return docToCustomItem(document, opts)
	})
}

func docToCustomItem(document *doc.Document, opts *newsml.Options) (item interface{}, err error) {
	// TODO Handle custom concept types here dependent on external config
	_ = opts
	return nil, fmt.Errorf("%w %s", ErrUnsupportedType, document.Type)
}

func XMLToItem(r io.Reader) (interface{}, error) {
	br := bufio.NewReader(r)

	tag, err := peekRootTag(br)
	if err != nil {
		return nil, err
	}

	var item interface{}

	switch tag {
	case "newsItem":
		item = &newsml.NewsItem{}
	case "conceptItem":
		item = &newsml.ConceptItem{}
	case "planningItem":
		item = &newsml.PlanningItem{}
	case "planning":
		item = &newsml.Assignment{}
	case "list":
		item = &listpackage.List{}
	case "package":
		item = &listpackage.Package{}
	default:
		return nil, fmt.Errorf("%w %q", ErrUnsupportedType, tag)
	}

	err = internal.NewXMLDecoder(br).Decode(item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

const peekSize = 1024

func peekRootTag(r *bufio.Reader) (string, error) {
	// Peek at the first bytes of the document so that we can
	// sniff the document type from the root tag.
	head, err := r.Peek(peekSize)
	if err != nil && err != bufio.ErrBufferFull && err != io.EOF {
		return "", fmt.Errorf("failed to read the beginning of the XML document: %w", err)
	}

	dec := internal.NewXMLDecoder(bytes.NewReader(head))

	for {
		// Read tokens from the XML document in a stream.
		t, err := dec.Token()
		if err == io.EOF {
			return "", fmt.Errorf("root tag was not found within the first %d bytes of the document", peekSize)
		}
		if err != nil {
			return "", fmt.Errorf("failed to parse XML: %w", err)
		}

		root, ok := t.(xml.StartElement)
		if !ok {
			continue
		}

		return root.Name.Local, nil
	}
}

func XMLToDoc(xmlDoc string, opts *newsml.Options) (*doc.Document, error) {
	item, err := XMLToItem(strings.NewReader(xmlDoc))
	if err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	return ItemToDoc(item, opts)
}

func ItemToDoc(item interface{}, opts *newsml.Options) (*doc.Document, error) {
	switch ptr := item.(type) {
	case *newsml.NewsItem:
		return NewsItemToDoc(ptr, opts)
	case *newsml.ConceptItem:
		return ConceptItemToDoc(ptr, opts)
	case *newsml.PlanningItem:
		return PlanningItemToDoc(ptr, opts)
	case *newsml.Assignment:
		return AssignmentToDoc(ptr, opts)
	case *listpackage.List:
		return ListToDoc(ptr, opts)
	case *listpackage.Package:
		return PackageToDoc(ptr, opts)
	default:
		return nil, &MalformedDocumentError{"unsupported item type"}
	}
}

func MergeOptions(opts *newsml.Options, customConfig *newsml.Options) {
	for k, v := range customConfig.Option {
		opts.Option[k] = v
	}
	opts.ConceptQcodeType = append(opts.ConceptQcodeType, customConfig.ConceptQcodeType...)
	opts.NewsmlQcodeType = append(opts.NewsmlQcodeType, customConfig.NewsmlQcodeType...)
	opts.AssignmentQcodeType = append(opts.AssignmentQcodeType, customConfig.AssignmentQcodeType...)
	opts.Status = append(opts.Status, customConfig.Status...)
	opts.PropertyType = append(opts.PropertyType, customConfig.PropertyType...)
	opts.ElementType = append(opts.ElementType, customConfig.ElementType...)

	var ce []newsml.LinkObjectExceptionType
	for c := range customConfig.LinkException {
		foundIt := false
		for d := range opts.LinkException {
			if customConfig.LinkException[c].Type == opts.LinkException[d].Type {
				opts.LinkException[d].Rel = customConfig.LinkException[c].Rel
				opts.LinkException[d].Section = mergeSections(opts.LinkException[d].Section, customConfig.LinkException[c].Section)
				foundIt = true
			}
		}
		if !foundIt {
			ce = append(ce, customConfig.LinkException[c])
		}
	}
	opts.LinkException = append(opts.LinkException, ce...)

	var oe []newsml.LinkObjectExceptionType
	for c := range customConfig.ObjectException {
		foundIt := false
		for d := range opts.ObjectException {
			if customConfig.ObjectException[c].Type == opts.ObjectException[d].Type {
				opts.ObjectException[d].Rel = customConfig.ObjectException[c].Rel
				opts.ObjectException[d].Section = mergeSections(opts.ObjectException[d].Section, customConfig.ObjectException[c].Section)
				foundIt = true
			}
		}
		if !foundIt {
			oe = append(oe, customConfig.ObjectException[c])
		}
	}
	opts.ObjectException = append(opts.ObjectException, oe...)

	var pe []newsml.PropertyExceptionType
	for c := range customConfig.PropertyException {
		foundIt := false
		for d := range opts.PropertyException {
			if customConfig.PropertyException[c].Type == opts.PropertyException[d].Type {
				opts.PropertyException[d].Section = mergeSections(opts.PropertyException[d].Section, customConfig.PropertyException[c].Section)
				foundIt = true
			}
		}
		if !foundIt {
			pe = append(pe, customConfig.PropertyException[c])
		}
	}
	opts.PropertyException = append(opts.PropertyException, pe...)

	opts.DataElementConversion = append(opts.DataElementConversion, customConfig.DataElementConversion...)

	for optKey, optValue := range opts.DateElements {
		for customKey, customValue := range customConfig.DateElements {
			if optKey == customKey {
				for customDateKey, customDateValue := range customValue {
					optValue[customDateKey] = customDateValue
				}
			}
		}
	}

	HTMLOpts := customConfig.HTMLSanitizeOptions

	// Make sure we do not overwrite default values with empty/default-struct values
	if HTMLOpts.AllowStandardAttributes != nil {
		opts.HTMLSanitizeOptions.AllowStandardAttributes = HTMLOpts.AllowStandardAttributes
	}

	if HTMLOpts.AllowTables != nil {
		opts.HTMLSanitizeOptions.AllowTables = HTMLOpts.AllowTables
	}

	if HTMLOpts.AllowImages != nil {
		opts.HTMLSanitizeOptions.AllowImages = HTMLOpts.AllowImages
	}

	if HTMLOpts.AllowLists != nil {
		opts.HTMLSanitizeOptions.AllowLists = HTMLOpts.AllowLists
	}

	if HTMLOpts.AllowRelativeURLs != nil {
		opts.HTMLSanitizeOptions.AllowRelativeURLs = HTMLOpts.AllowRelativeURLs
	}

	if HTMLOpts.AllowableURLSchemes != "" {
		opts.HTMLSanitizeOptions.AllowableURLSchemes = HTMLOpts.AllowableURLSchemes
	}

	opts.HTMLSanitizeOptions.ElementsAttributes = append(customConfig.HTMLSanitizeOptions.ElementsAttributes, opts.HTMLSanitizeOptions.ElementsAttributes...)

	// Remove duplicates
	elattrs := make([]newsml.HTMLElementAttributes, 0)
	elms := make(map[string]bool)
	for i := 0; i < len(opts.HTMLSanitizeOptions.ElementsAttributes); i++ {
		if _, ok := elms[opts.HTMLSanitizeOptions.ElementsAttributes[i].Elements]; !ok {
			elms[opts.HTMLSanitizeOptions.ElementsAttributes[i].Elements] = true
			elattrs = append(elattrs, opts.HTMLSanitizeOptions.ElementsAttributes[i])
		}
	}
	opts.HTMLSanitizeOptions.ElementsAttributes = elattrs
}

func mergeSections(dflt []string, cust []string) []string {
	allSections := append(dflt, cust...)

	noDups := map[string]bool{}
	for _, s := range allSections {
		noDups[s] = true
	}

	var keys []string

	for key := range noDups {
		switch key {
		case "itemmeta-newsitem":
			delete(noDups, "contentmeta")
		case "itemmeta-concept":
			delete(noDups, "concept")
		case "itemmeta-conceptitem":
			delete(noDups, "conceptitem")
		case "itemmeta-planning":
			delete(noDups, "planning")
		}
	}
	for key := range noDups {
		keys = append(keys, key)
	}

	return keys
}
