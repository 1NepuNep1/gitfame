package cmd

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

func Validate(opts Options) error { // Тут я попросил ллмку красивые тесты ошибки мне написать :)
	if opts.Repository == "" {
		return fmt.Errorf("--repository cannot be empty")
	}

	if !strings.HasPrefix(opts.Revision, "HEAD") {
		validHash := regexp.MustCompile(`^[0-9a-fA-F]{7,40}$`)
		validTag := regexp.MustCompile(`^[^\-][\w.\-\/]+$`)

		if !validHash.MatchString(opts.Revision) && !validTag.MatchString(opts.Revision) {
			return fmt.Errorf("invalid revision format %q (must be 'HEAD', 'HEAD~n', a valid commit hash, or a tag)", opts.Revision)
		}
	}

	switch opts.OrderBy {
	case "lines", "commits", "files":
	default:
		return fmt.Errorf("invalid --order-by value %q (must be lines, 'commits', or 'files')", opts.OrderBy)
	}

	switch opts.Format {
	case "tabular", "csv", "json", "json-lines":
	default:
		return fmt.Errorf("invalid --format value %q (must be 'tabular', 'csv', 'json', or 'json-lines')", opts.Format)
	}

	for _, ext := range opts.Extensions {
		if ext == "" {
			return fmt.Errorf("empty extension provided in --extensions")
		}
		if !strings.HasPrefix(ext, ".") {
			return fmt.Errorf("invalid extension %q in --extensions: must start with a dot", ext)
		}
	}

	for _, pattern := range opts.Exclude {
		if pattern == "" {
			return fmt.Errorf("empty pattern provided in --exclude")
		}
		if _, err := filepath.Match(pattern, "bebebe"); err != nil {
			return fmt.Errorf("invalid glob pattern %q in --exclude: %v", pattern, err)
		}
	}

	for _, pattern := range opts.RestrictTo {
		if pattern == "" {
			return fmt.Errorf("empty pattern provided in --restrict-to")
		}
		if _, err := filepath.Match(pattern, "bebebe"); err != nil {
			return fmt.Errorf("invalid glob pattern %q in --restrict-to: %v", pattern, err)
		}
	}

	return nil
}
