package internal

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/mapping"
	index "github.com/blevesearch/bleve_index_api"
	"time"
)

type DialogDocument struct {
	Publication string `json:"publication"`
	Series      int32  `json:"series"`
	Date        string `json:"date"`
	Actor       string `json:"actor"`
	Content     string `json:"content"`
	ContentType string `json:"type"`
}

func RskIndexMapping() (mapping.IndexMapping, error) {

	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	dialogMapping := bleve.NewDocumentMapping()

	dialogMapping.AddFieldMappingsAt("content", englishTextFieldMapping)
	dialogMapping.AddFieldMappingsAt("publication", keywordFieldMapping)
	dialogMapping.AddFieldMappingsAt("series", bleve.NewNumericFieldMapping())
	dialogMapping.AddFieldMappingsAt("date", bleve.NewDateTimeFieldMapping())
	dialogMapping.AddFieldMappingsAt("type", keywordFieldMapping)
	dialogMapping.AddFieldMappingsAt("actor", keywordFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("dialog", dialogMapping)

	indexMapping.DefaultMapping = dialogMapping
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}

type RawDocument struct {
	ID     string
	Fields map[string]interface{}
}

func DecodeDocument(doc index.Document) *RawDocument {
	rv := &RawDocument{
		ID:     doc.ID(),
		Fields: map[string]interface{}{},
	}

	doc.VisitFields(func(field index.Field) {
		var newval interface{}
		switch field := field.(type) {
		case index.TextField:
			newval = field.Text()
		case index.NumericField:
			n, err := field.Number()
			if err == nil {
				newval = n
			}
		case index.DateTimeField:
			d, err := field.DateTime()
			if err == nil {
				newval = d.Format(time.RFC3339Nano)
			}
		}
		existing, existed := rv.Fields[field.Name()]
		if existed {
			switch existing := existing.(type) {
			case []interface{}:
				rv.Fields[field.Name()] = append(existing, newval)
			case interface{}:
				arr := make([]interface{}, 2)
				arr[0] = existing
				arr[1] = newval
				rv.Fields[field.Name()] = arr
			}
		} else {
			rv.Fields[field.Name()] = newval
		}
	})

	return rv
}
