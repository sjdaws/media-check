package arr

import (
	"strconv"

	"github.com/sjdaws/media-check/internal/config"
	"github.com/sjdaws/media-check/internal/matcher/sources"
)

type TVShow struct {
	ID         int        `json:"id"`
	IMDBID     string     `json:"imdbId"`
	Path       string     `json:"path"`
	Statistics Statistics `json:"statistics"`
	TMDBID     int        `json:"tmdbId"`
	TVDBID     int        `json:"tvdbId"`
	Title      string     `json:"title"`
	Year       int        `json:"year"`
}

type Statistics struct {
	Episodes int `json:"episodeFileCount"`
}

func fetchTVShows(cfg config.MatchSourcesArr) ([]*sources.Media, error) {
	var result []TVShow

	err := callApi(cfg.ApiKey, cfg.Server+"/api/v3/series", &result)
	if err != nil {
		return nil, err
	}

	shows := make([]*sources.Media, 0)

	for _, show := range result {
		// Ignore shows with no episodes
		if show.Statistics.Episodes < 1 {
			continue
		}

		shows = append(shows, &sources.Media{
			Folders: []string{rewriteFolder(show.Path, cfg.PathMap)},
			Guids: sources.Guids{
				IMDB: show.IMDBID,
				TMDB: show.TMDBID,
				TVDB: show.TVDBID,
			},
			ID:     strconv.Itoa(show.ID),
			Source: "sonarr",
			Title:  show.Title,
			Year:   show.Year,
		})
	}

	return shows, nil
}
