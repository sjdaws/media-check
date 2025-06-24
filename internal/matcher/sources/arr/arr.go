package arr

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sjdaws/media-check/internal/config"
	"github.com/sjdaws/media-check/internal/matcher/sources"
)

func Fetch(cfg config.Match) (*sources.Source, error) {
	items := sources.Source{
		Media:     make(map[string]*sources.Media),
		Multiples: make([]*sources.Media, 0),
	}

	err := fetchMovies(cfg.Sources.Radarr, &items)
	if err != nil {
		return nil, err
	}

	err = fetchTVShows(cfg.Sources.Sonarr, &items)
	if err != nil {
		return nil, err
	}

	return &items, nil
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
