package mapping

type FieldType string

const (
	FieldTypeKeyword FieldType = "keyword"
	FieldTypeText    FieldType = "text"
	FieldTypeNumber  FieldType = "number"
	FieldTypeDate    FieldType = "date"
)

var Mapping = map[string]FieldType{
	"publication": FieldTypeKeyword,
	"series":      FieldTypeNumber,
	"episode":     FieldTypeNumber,
	"date":        FieldTypeDate,
	"actor":       FieldTypeKeyword,
	"pos":         FieldTypeNumber,
	"content":     FieldTypeText,
	"type":        FieldTypeKeyword,
}
