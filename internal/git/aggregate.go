package git

import (
	"sort"
)

type AuthorStats struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

type tempAuthorStats struct {
	Name    string
	Lines   int
	Commits map[string]bool
	Files   map[string]bool
}

func AggregateAndSort(result *BlameResult, orderBy string) []AuthorStats {
	statsByAuthor := make(map[string]*tempAuthorStats)

	for commitHash, ci := range result.Commits {
		author := ci.Author

		s, ok := statsByAuthor[author]
		if !ok {
			s = &tempAuthorStats{
				Name:    author,
				Lines:   0,
				Commits: make(map[string]bool),
				Files:   make(map[string]bool),
			}
			statsByAuthor[author] = s
		}
		s.Lines += ci.LineCount
		s.Commits[commitHash] = true

		for f := range ci.Files {
			s.Files[f] = true
		}
	}

	var finalStats []AuthorStats
	for author, temp := range statsByAuthor {
		finalStats = append(finalStats, AuthorStats{
			Name:    author,
			Lines:   temp.Lines,
			Commits: len(temp.Commits),
			Files:   len(temp.Files),
		})
	}

	sort.Slice(finalStats, func(i, j int) bool {
		return compareStats(finalStats[i], finalStats[j], orderBy)
	})

	return finalStats
}
