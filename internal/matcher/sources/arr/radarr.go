package arr

import (
	"strconv"
	"strings"

	"github.com/sjdaws/media-check/internal/config"
	"github.com/sjdaws/media-check/internal/matcher/sources"
)

type Movie struct {
	FolderName string `json:"folderName"`
	HasFile    bool   `json:"hasFile"`
	ID         int    `json:"id"`
	IMDBID     string `json:"imdbId"`
	TMDBID     int    `json:"tmdbId"`
	TVDBID     int    `json:"tvdbId"`
	Title      string `json:"title"`
	Year       int    `json:"year"`
}

func fetchMovies(cfg config.MatchSourcesArr, items *sources.Source) error {
	var result []Movie

	if cfg.ApiKey == "" || cfg.Server == "" {
		return nil
	}

	err := callApi(cfg.ApiKey, cfg.Server+"/api/v3/movie", &result)
	if err != nil {
		return err
	}

	for _, movie := range result {
		// Ignore unfetched movies
		if !movie.HasFile {
			continue
		}

		// Get folder name
		folder := "/movies/" + sources.RewriteFolder(movie.FolderName, cfg.Paths)
		lowerFolder := strings.ToLower(folder)

		media := &sources.Media{
			Guids: sources.Guids{
				IMDB: movie.IMDBID,
				TMDB: movie.TMDBID,
			},
			ID:     strconv.Itoa(movie.ID),
			Path:   folder,
			Source: "radarr",
			Title:  movie.Title,
			Year:   movie.Year,
		}

		if _, ok := items.Media[lowerFolder]; ok {
			items.Multiples = append(items.Multiples, media)

			continue
		}

		items.Media[lowerFolder] = media
	}

	return nil
}
