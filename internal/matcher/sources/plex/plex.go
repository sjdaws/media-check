package plex

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/sjdaws/media-check/internal/config"
	"github.com/sjdaws/media-check/internal/matcher/sources"
)

type metadata struct {
	ID               int
	LibrarySectionID int
	MetadataType     int
	Title            string
	Year             int
}

type path struct {
	File string
}

type tag struct {
	Tag string
}

func Fetch(cfg config.MatchSourcesPlex) (*sources.Source, error) {
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

	sections := append(cfg.Sections.Movies, cfg.Sections.TVShows...)

	var result []metadata
	orm.Raw("SELECT id, library_section_id, metadata_type, title, year FROM metadata_items WHERE library_section_id IN (?) AND metadata_type IN (1,2)", sections).Scan(&result)

	for _, row := range result {
		media := &sources.Media{
			ID:     strconv.Itoa(row.ID),
			Source: "plex",
			Title:  row.Title,
			Year:   row.Year,
		}

		// Get guids
		var tags []tag
		orm.Raw("SELECT tags.tag FROM metadata_items INNER JOIN taggings, tags ON metadata_items.id = taggings.metadata_item_id AND tags.id = taggings.tag_id WHERE metadata_items.id = ? AND tags.tag_type = 314", row.ID).Scan(&tags)

		for _, tagRow := range tags {
			parts := strings.SplitN(tagRow.Tag, "://", 2)
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

		// Get folders
		var paths []path
		var id any

		equals := "="
		id = row.ID

		if row.MetadataType == 2 {
			equals = "IN"
			id = getMetadataAncestors(orm, []int{row.ID})
		}

		orm.Raw(fmt.Sprintf("SELECT media_parts.file FROM metadata_items INNER JOIN media_parts, media_items ON media_items.metadata_item_id = metadata_items.id AND media_parts.media_item_id = media_items.id WHERE metadata_items.id %s ?", equals), id).Scan(&paths)

		known := make(map[string]string)

		for _, folderRow := range paths {
			media.Path = getFolder(cfg, filepath.Dir(folderRow.File), row.MetadataType == 1)
			lowerFolder := strings.ToLower(media.Path)

			// Only capture paths once
			_, ok := known[lowerFolder]
			if !ok {
				known[lowerFolder] = media.Path

				if _, ok = items.Media[lowerFolder]; ok {
					items.Multiples = append(items.Multiples, media)

					continue
				}

				items.Media[lowerFolder] = media
			}
		}
	}

	return items, nil
}

func getMetadataAncestors(orm *gorm.DB, parents []int) []int {
	var children []int
	orm.Raw("SELECT id FROM metadata_items WHERE parent_id IN (?)", parents).Scan(&children)

	if len(children) == 0 {
		return parents
	}

	return append(parents, getMetadataAncestors(orm, children)...)
}

func getFolder(cfg config.MatchSourcesPlex, path string, movie bool) string {
	if movie {
		return "/movies/" + sources.RewriteFolder(path, cfg.Paths.Movies)
	}

	// Double de-dir tv libraries to remove season folders
	return "/tvshows/" + sources.RewriteFolder(filepath.Dir(path), cfg.Paths.TVShows)
}
