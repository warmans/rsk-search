package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type SwaggerFile struct {
	Swagger      string              `json:"swagger"`
	Info         map[string]string   `json:"info"`
	Tags         []map[string]string `json:"tags"`
	Scehems      []string            `json:"schemes"`
	Consumes     []string            `json:"consumes"`
	Produces     []string            `json:"produces"`
	ExternalDocs map[string]string   `json:"externalDocs"`

	Paths       map[string]json.RawMessage `json:"paths"`
	Definitions map[string]json.RawMessage `json:"definitions"`
}

func main() {

	base := flag.String("base", "", "file to merge into")
	merge := flag.String("merge", "", "comma separated paths to merge into the root file")
	out := flag.String("out", "./swagger.json", "output file path")

	flag.Parse()

	if base == nil {
		panic("no base file given")
	}
	if merge == nil {
		panic("no merge file/s given")
	}

	baseFile, err := parseFile(*base)
	if err != nil {
		panic(fmt.Sprintf("failed to parse base file %s reason: %s", *base, err.Error()))
	}

	for _, f := range strings.Split(*merge, ",") {
		parsed, err := parseFile(strings.TrimSpace(f))
		if err != nil {
			panic(fmt.Sprintf("failed to parse path %s reason: %s", f, err.Error()))
		}
		for k, v := range parsed.Paths {
			if _, ok := baseFile.Paths[k]; ok {
				panic(fmt.Sprintf("duplicate path: %s", k))
			}
			baseFile.Paths[k] = v
		}
		for k, v := range parsed.Definitions {
			if _, ok := baseFile.Definitions[k]; ok {
				fmt.Println("ignored duplicate definition: ", k)
				continue
			}
			baseFile.Definitions[k] = v
		}
	}
	outFile, err := os.Create(*out)
	if err != nil {
		panic(fmt.Sprintf("failed to create outfile at path %s reason: %s", *out, err.Error()))
	}
	enc := json.NewEncoder(outFile)
	enc.SetIndent("  ", "  ")
	if err := enc.Encode(baseFile); err != nil {
		panic("failed to encode output: " + err.Error())
	}
}

func parseFile(path string) (*SwaggerFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	swaggerData := &SwaggerFile{}
	return swaggerData, json.NewDecoder(f).Decode(swaggerData)
}
