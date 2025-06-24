package jellyfin

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/sjdaws/media-check/internal/config"
	"github.com/sjdaws/media-check/internal/matcher/sources"
)

type metadata struct {
	Guid           string
	Name           string
	Path           string
	ProductionYear int
	ProviderIds    string
	Type           string
}

func Fetch(cfg config.MatchSourcesJellyfin) (*sources.Source, error) {
	items := &sources.Source{
		Media:     make(map[string]*sources.Media),
		Multiples: make([]*sources.Media, 0),
	}

	if cfg.Database == "" {
		return items, nil
	}

	orm, err := gorm.Open(sqlite.Open(cfg.Database), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	types := []string{"MediaBrowser.Controller.Entities.Movies.Movie", "MediaBrowser.Controller.Entities.TV.Series"}

	var result []metadata
	orm.Raw("SELECT guid AS Guid, Name, Path, ProductionYear, ProviderIds, type AS Type FROM TypedBaseItems WHERE `type` IN (?)", types).Scan(&result)

	for _, row := range result {
		media := &sources.Media{
			ID:     row.Guid,
			Source: "jellyfin",
			Title:  row.Name,
			Year:   row.ProductionYear,
		}

		// Get guids
		tags := strings.Split(row.ProviderIds, "|")

		for _, tag := range tags {
			parts := strings.SplitN(tag, "=", 2)
			switch strings.ToLower(parts[0]) {
			case "imdb":
				media.Guids.IMDB = parts[1]
			case "tmdb":
				i, _ := strconv.Atoi(parts[1])
				media.Guids.TMDB = i
			case "tvdb":
				i, _ := strconv.Atoi(parts[1])
				media.Guids.TVDB = i
			}
		}

		// Get folder name
		media.Path = getFolder(cfg, row.Path, row.Type == "MediaBrowser.Controller.Entities.Movies.Movie")
		lowerFolder := strings.ToLower(media.Path)

		if _, ok := items.Media[lowerFolder]; ok {
			items.Multiples = append(items.Multiples, media)

			continue
		}

		items.Media[lowerFolder] = media
	}

	return items, nil
}

func getFolder(cfg config.MatchSourcesJellyfin, path string, movie bool) string {
	if movie {
		// Get directory for movies since path will be full path to file
		return "/movies/" + sources.RewriteFolder(filepath.Dir(path), cfg.Paths.Movies)
	}

	return "/tvshows/" + sources.RewriteFolder(path, cfg.Paths.TVShows)
}
