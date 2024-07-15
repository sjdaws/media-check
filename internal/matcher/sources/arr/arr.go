package arr

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/sjdaws/media-check/internal/config"
	"github.com/sjdaws/media-check/internal/matcher/sources"
)

func Fetch(cfg config.Match) ([]*sources.Media, error) {
	movies, err := fetchMovies(cfg.Sources.Radarr)
	if err != nil {
		return nil, err
	}

	tvshows, err := fetchTVShows(cfg.Sources.Sonarr)
	if err != nil {
		return nil, err
	}

	return append(movies, tvshows...), nil
}

func callApi(apiKey string, endpoint string, result any) error {
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("X-Api-Key", apiKey)

	client := &http.Client{
		Timeout: 120 * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout: 5 * time.Second,
		},
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return err
	}

	return nil
}

func rewriteFolder(folder string, pathmap []string) string {
	for _, mapping := range pathmap {
		parts := strings.SplitN(mapping, ":", 2)

		if len(parts) == 2 {
			folder = strings.Replace(folder, parts[0], parts[1], 1)
		}
	}

	return folder
}
