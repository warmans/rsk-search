package roadmap

import (
	"embed"
	"fmt"
	"io"
)

//go:embed data/roadmap.md
var roadmap embed.FS

var markdown string

func init() {
	f, err := roadmap.Open("data/roadmap.md")
	if err != nil {
		panic("failed to open embedded roadmap data: " + err.Error())
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("Failed to load roadmap")
		return
	}
	markdown = string(data)
}

func GetMarkdown() string {
	return markdown
}
