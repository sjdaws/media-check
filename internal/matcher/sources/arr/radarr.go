package arr

import (
	"strconv"

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

func fetchMovies(cfg config.MatchSourcesArr) ([]*sources.Media, error) {
	var result []Movie

	err := callApi(cfg.ApiKey, cfg.Server+"/api/v3/movie", &result)
	if err != nil {
		return nil, err
	}

	movies := make([]*sources.Media, 0)

	for _, movie := range result {
		// Ignore missing files
		if !movie.HasFile {
			continue
		}

		movies = append(movies, &sources.Media{
			Folders: []string{rewriteFolder(movie.FolderName, cfg.PathMap)},
			Guids: sources.Guids{
				IMDB: movie.IMDBID,
				TMDB: movie.TMDBID,
			},
			ID:     strconv.Itoa(movie.ID),
			Source: "radarr",
			Title:  movie.Title,
			Year:   movie.Year,
		})
	}

	return movies, nil
}
