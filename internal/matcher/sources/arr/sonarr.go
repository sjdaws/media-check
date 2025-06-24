package arr

import (
	"strconv"
	"strings"

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

func fetchTVShows(cfg config.MatchSourcesArr, items *sources.Source) error {
	var result []TVShow

	if cfg.ApiKey == "" || cfg.Server == "" {
		return nil
	}

	err := callApi(cfg.ApiKey, cfg.Server+"/api/v3/series", &result)
	if err != nil {
		return err
	}

	for _, tvshow := range result {
		// Ignore shows with no episodes
		if tvshow.Statistics.Episodes < 1 {
			continue
		}

		// Get folder name
		folder := "/tvshows/" + sources.RewriteFolder(tvshow.Path, cfg.Paths)
		lowerFolder := strings.ToLower(folder)

		media := &sources.Media{
			Guids: sources.Guids{
				IMDB: tvshow.IMDBID,
				TMDB: tvshow.TMDBID,
				TVDB: tvshow.TVDBID,
			},
			ID:     strconv.Itoa(tvshow.ID),
			Path:   folder,
			Source: "sonarr",
			Title:  tvshow.Title,
			Year:   tvshow.Year,
		}

		if _, ok := items.Media[lowerFolder]; ok {
			items.Multiples = append(items.Multiples, media)

			continue
		}

		items.Media[lowerFolder] = media
	}

	return nil
}
