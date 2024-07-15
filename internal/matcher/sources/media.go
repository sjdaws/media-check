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
	Guids   Guids
	ID      string
	Folders []string
	Matched bool
	Source  string
	Title   string
	Year    int
}

type Result int

const (
	Match Result = iota
	GuidMismatch
	NoMatch
)

func (m *Media) Match(in []*Media) (Result, *Media) {
	for _, media := range in {
		if strings.TrimSuffix(m.Folders[0], "/") == strings.TrimSuffix(media.Folders[0], "/") {
			if (m.Guids == (Guids{}) && media.Guids == (Guids{}) ||
				m.Guids.IMDB != "" && m.Guids.IMDB == media.Guids.IMDB) ||
				(m.Guids.TMDB != 0 && m.Guids.TMDB == media.Guids.TMDB) ||
				(m.Guids.TVDB != 0 && m.Guids.TVDB == media.Guids.TVDB) {
				return Match, media
			}

			return GuidMismatch, media
		}
	}

	return NoMatch, nil
}
