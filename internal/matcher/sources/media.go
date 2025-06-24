package sources

import (
	"strings"
)

type Guids struct {
	IMDB string
	TMDB int
	TVDB int
}

type Media struct {
	Guids     Guids
	ID        string
	Matched   bool
	Path      string
	Source    string
	Title     string
	Unmatched string
	Year      int
}

type Result int

type Source struct {
	Media     map[string]*Media
	Multiples []*Media
}

const (
	Match Result = iota
	GuidMismatch
	NoMatch
)

func (m *Media) Match(key string, in map[string]*Media) (Result, *Media) {
	found, ok := in[key]

	if !ok {
		return NoMatch, nil
	}

	if (m.Guids == (Guids{}) && found.Guids == (Guids{}) ||
		m.Guids.IMDB != "" && strings.EqualFold(m.Guids.IMDB, found.Guids.IMDB)) ||
		(m.Guids.TMDB != 0 && m.Guids.TMDB == found.Guids.TMDB) ||
		(m.Guids.TVDB != 0 && m.Guids.TVDB == found.Guids.TVDB) {
		return Match, found
	}

	return GuidMismatch, found
}

func RewriteFolder(folder string, paths []string) string {
	for _, path := range paths {
		lowerFolder := strings.ToLower(folder)
		lowerPath := strings.ToLower(path)

		if strings.HasPrefix(lowerFolder, lowerPath) {
			return strings.Trim(folder[len(path):], "/")
		}
	}

	return folder
}
