package meta

import "github.com/warmans/rsk-search/gen/api"

type FieldKind string

func (k FieldKind) Proto() api.FieldMeta_Kind {
	switch k {
	case FieldIdentifier:
		return api.FieldMeta_IDENTIFIER
	case FieldKeyword:
		return api.FieldMeta_KEYWORD
	case FieldKeywordList:
		return api.FieldMeta_KEYWORD_LIST
	case FieldText:
		return api.FieldMeta_TEXT
	case FieldInt:
		return api.FieldMeta_INT
	case FieldFloat:
		return api.FieldMeta_FLOAT
	case FieldDate:
		return api.FieldMeta_DATE
	}
	return api.FieldMeta_UNKNOWN
}

const (
	// non-analyzed unique keywords (e.g. dsf342f32f3)
	FieldIdentifier = FieldKind("identifier")

	// non-analyzed keywords (e.g. foo)
	FieldKeyword = FieldKind("keyword")

	// list of keywords e.g. [foo, bar]
	FieldKeywordList = FieldKind("keyword_list")

	// analyzed words (e.g. "foo bar, baz"
	FieldText = FieldKind("text")

	// whole numbers e.g. 1
	FieldInt = FieldKind("int")

	// real numbers e.g. 1.5
	FieldFloat = FieldKind("float")

	// datestamp in RFC339 e.g. 2020-01-25T00:00:00Z
	FieldDate = FieldKind("date")
)

type SearchMeta struct {
	Fields []FieldMeta
}

func (m SearchMeta) Proto() *api.SearchMetadata {
	sm := &api.SearchMetadata{
		Fields: make([]*api.FieldMeta, len(m.Fields)),
	}
	for k, v := range m.Fields {
		sm.Fields[k] = v.Proto()
	}
	return sm
}

type FieldMeta struct {
	Name string
	Kind FieldKind
}

func (m FieldMeta) Proto() *api.FieldMeta {
	return &api.FieldMeta{
		Name: m.Name,
		Kind: m.Kind.Proto(),
	}
}

func GetSearchMeta() SearchMeta {
	return SearchMeta{
		Fields: []FieldMeta{
			{Name: "id", Kind: FieldIdentifier},
			{Name: "publication", Kind: FieldKeyword},
			{Name: "series", Kind: FieldInt},
			{Name: "episode", Kind: FieldInt},
			//{Name: "date", Kind: FieldDate}, // not sure how to do the blieve date matching
			{Name: "actor", Kind: FieldKeyword},
			{Name: "content", Kind: FieldText},
			{Name: "type", Kind: FieldKeyword},
			{Name: "tags", Kind: FieldKeywordList},
		},
	}
}
