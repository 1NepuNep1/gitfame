package git

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

//go:embed configs/language_extensions.json
var languagesJSON []byte

type languagesDescriptor struct {
	Name       string   // json: "name"
	Type       string   // json: "type"
	Extensions []string // json: "extensions"
}

var allLanguages []languagesDescriptor
var langToExt map[string][]string

func init() {
	if err := json.Unmarshal(languagesJSON, &allLanguages); err != nil {
		panic(fmt.Sprintf("failed to parse languages.json: %v", err))
	}

	langToExt = make(map[string][]string)
	for _, lang := range allLanguages {
		nameLower := strings.ToLower(lang.Name)
		langToExt[nameLower] = lang.Extensions
	}
}

func getExtensionsForLanguage(langName string) []string {
	return langToExt[strings.ToLower(langName)]
}

func buildAllowedExtensions(languages, exts []string) (map[string]bool, error) {
	allowed := make(map[string]bool)

	var err error

	for _, lang := range languages {
		lang = strings.TrimSpace(lang)
		if lang == "" {
			continue
		}

		foundExts := getExtensionsForLanguage(lang)
		if len(foundExts) == 0 {
			err = fmt.Errorf("unknown or extension-less language: %q", lang)
		}

		for _, e := range foundExts {
			eLower := strings.ToLower(e)
			allowed[eLower] = true
		}
	}

	for _, ext := range exts {
		ext = strings.TrimSpace(ext)
		if ext == "" {
			continue
		}
		extLower := strings.ToLower(ext)
		allowed[extLower] = true
	}

	return allowed, err
}

func matchPattern(file string, pattern []string) bool {
	for _, blob := range pattern {
		match, err := filepath.Match(blob, file)
		if err != nil { // Errors are checked in Validate func
			continue
		}
		if match {
			return true
		}
	}
	return false
}

func compareStats(a, b AuthorStats, orderBy string) bool {
	switch orderBy {
	case "lines":
		if a.Lines != b.Lines {
			return a.Lines > b.Lines
		}
		if a.Commits != b.Commits {
			return a.Commits > b.Commits
		}
		if a.Files != b.Files {
			return a.Files > b.Files
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)

	case "commits":
		if a.Commits != b.Commits {
			return a.Commits > b.Commits
		}
		if a.Lines != b.Lines {
			return a.Lines > b.Lines
		}
		if a.Files != b.Files {
			return a.Files > b.Files
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)

	case "files":
		if a.Files != b.Files {
			return a.Files > b.Files
		}
		if a.Lines != b.Lines {
			return a.Lines > b.Lines
		}
		if a.Commits != b.Commits {
			return a.Commits > b.Commits
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)

	default:
		if a.Lines != b.Lines {
			return a.Lines > b.Lines
		}
		if a.Commits != b.Commits {
			return a.Commits > b.Commits
		}
		if a.Files != b.Files {
			return a.Files > b.Files
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	}
}
