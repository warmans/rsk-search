package mapping

type FieldType string

const (
	FieldTypeKeyword  FieldType = "keyword"
	FieldTypeText     FieldType = "text"
	FieldTypeNumber   FieldType = "number"
	FieldTypeDate     FieldType = "date"
	FieldTypeShingles FieldType = "shingles"
)

var Mapping = map[string]FieldType{
	"transcript_id": FieldTypeKeyword,
	"publication":   FieldTypeKeyword,
	"series":        FieldTypeNumber,
	"episode":       FieldTypeNumber,
	"date":          FieldTypeDate,
	"actor":         FieldTypeKeyword,
	"pos":           FieldTypeNumber,
	"content":       FieldTypeText,
	"type":          FieldTypeKeyword,
	"special":       FieldTypeKeyword, // true or false
}
