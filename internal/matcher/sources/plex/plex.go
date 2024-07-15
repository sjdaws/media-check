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

type folder struct {
	File string
}

type metadata struct {
	ID               int
	LibrarySectionID int
	MetadataType     int
	Title            string
	Year             int
}

type tag struct {
	Tag string
}

func Fetch(cfg config.Match) ([]*sources.Media, error) {
	items := make([]*sources.Media, 0)

	orm, err := gorm.Open(sqlite.Open(cfg.Sources.Plex.Database), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sections := append(cfg.Sources.Plex.Sections.Movies, cfg.Sources.Plex.Sections.TVShows...)

	var result []metadata
	orm.Raw("SELECT id, library_section_id, metadata_type, title, year FROM metadata_items WHERE parent_id IS NULL AND library_section_id IN (?) AND metadata_type IN (1,2)", sections).Scan(&result)

	for _, row := range result {
		item := &sources.Media{
			ID:     strconv.Itoa(row.ID),
			Source: "plex",
			Title:  row.Title,
			Year:   row.Year,
		}

		// Get guids
		var tags []tag
		orm.Raw("SELECT tags.tag FROM metadata_items INNER JOIN taggings, tags ON metadata_items.id = taggings.metadata_item_id AND tags.id = taggings.tag_id WHERE metadata_items.id = ? AND metadata_items.library_section_id = ? AND tags.tag_type = 314", row.ID, row.LibrarySectionID).Scan(&tags)

		for _, tagRow := range tags {
			parts := strings.SplitN(tagRow.Tag, "://", 2)
			switch strings.ToLower(parts[0]) {
			case "imdb":
				item.Guids.IMDB = parts[1]
			case "tmdb":
				i, _ := strconv.Atoi(parts[1])
				item.Guids.TMDB = i
			case "tvdb":
				i, _ := strconv.Atoi(parts[1])
				item.Guids.TVDB = i
			}
		}

		// Get folders
		var folders []folder
		var id any

		equals := "="
		id = row.ID

		if row.MetadataType == 2 {
			equals = "IN"
			id = getMetadataAncestors(orm, []int{row.ID})
		}

		orm.Raw(fmt.Sprintf("SELECT media_parts.file FROM metadata_items INNER JOIN media_parts, media_items ON media_items.metadata_item_id = metadata_items.id AND media_parts.media_item_id = media_items.id WHERE metadata_items.id %s ? AND metadata_items.library_section_id = ?", equals), id, row.LibrarySectionID).Scan(&folders)

		paths := make(map[string]string)

		for _, folderRow := range folders {
			path := filepath.Dir(folderRow.File)

			// Double de-dir tv libraries to remove season folders
			if row.MetadataType == 2 {
				path = filepath.Dir(path)
			}

			// Only capture paths once
			_, known := paths[path]
			if !known {
				paths[path] = path
				item.Folders = append(item.Folders, path)
			}
		}

		items = append(items, item)
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
