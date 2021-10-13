package meta

import (
	"embed"
	"encoding/json"
)

//go:embed data/tags.json
var tagIndex embed.FS

var parsedTagIndex = Tags{}

func init() {
	f, err := tagIndex.Open("data/tags.json")
	if err != nil {
		panic("failed to open embedded metadata: " + err.Error())
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	if err := dec.Decode(&parsedTagIndex); err != nil {
		panic("failed to decode metadata: " + err.Error())
	}
}

func GetTag(tagName string) *Tag {
	return parsedTagIndex[tagName]
}

func Unalias(tag string) string {
	t, ok := parsedTagIndex[tag]
	if !ok {
		return ""
	}
	if alias := t.AliasOf; alias != "" {
		return Unalias(alias)
	}
	return tag
}

func GetTagKinds(tag string) []string {
	canonicalTag := Unalias(tag)
	if canonicalTag == "" {
		return nil
	}
	t, ok := parsedTagIndex[canonicalTag]
	if !ok {
		return nil
	}
	return t.Kind
}

type Tags map[string]*Tag

type Tag struct {
	Kind    []string `json:"kind"`
	AliasOf string   `json:"alias"`
}

