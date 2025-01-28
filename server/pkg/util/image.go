package util

import (
	"fmt"
	"path"
	"strings"
)

func ThumbPath(dirPath string, fileName string) string {
	return path.Join(dirPath, ThumbName(fileName))
}

func ThumbName(fileName string) string {
	ext := path.Ext(fileName)
	return fmt.Sprintf("%s.thumb%s", strings.TrimSuffix(fileName, ext), ext)
}
