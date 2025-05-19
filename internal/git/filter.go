package git

import (
	_ "embed"
	"path/filepath"
	"strings"
)

func FilterFiles(files []string, excludes, restricts, languages, exts []string) ([]string, error) {
	var err error
	allowedExts, err := buildAllowedExtensions(languages, exts)

	var newFiles []string

	for _, file := range files {

		if len(excludes) > 0 && matchPattern(file, excludes) {
			continue
		}

		if len(restricts) > 0 && !matchPattern(file, restricts) {
			continue
		}

		if len(allowedExts) > 0 {
			ext := strings.ToLower(filepath.Ext(file))
			if !allowedExts[ext] {
				continue
			}
		}

		newFiles = append(newFiles, file)
	}
	return newFiles, err
}
