package matcher

import (
	"slices"

	"github.com/sjdaws/media-check/internal/config"
	"github.com/sjdaws/media-check/internal/matcher/sources"
	"github.com/sjdaws/media-check/internal/matcher/sources/arr"
	"github.com/sjdaws/media-check/internal/matcher/sources/jellyfin"
	"github.com/sjdaws/media-check/internal/matcher/sources/plex"
)

type Result struct {
	GuidMismatch    [][]*sources.Media
	MultipleFolders []*sources.Media
	NoFolders       []*sources.Media
	Unmatched       []*sources.Media
}

func Check(cfg config.Match) (Result, error) {
	var result Result

	arrMedia, err := arr.Fetch(cfg)
	if err != nil {
		return result, err
	}

	plexMedia, err := plex.Fetch(cfg.Sources.Plex)
	if err != nil {
		return result, err
	}

	findMatches(arrMedia, plexMedia, cfg.Sources.Plex.Exceptions, &result)
	addUnmatched(arrMedia, &result, "plex")

	arrMedia, err = arr.Fetch(cfg)
	if err != nil {
		return result, err
	}

	jellyfinMedia, err := jellyfin.Fetch(cfg.Sources.Jellyfin)
	if err != nil {
		return result, err
	}

	findMatches(arrMedia, jellyfinMedia, cfg.Sources.Jellyfin.Exceptions, &result)
	addUnmatched(arrMedia, &result, "jellyfin")

	return result, nil
}

func addUnmatched(arr *sources.Source, result *Result, player string) {
	for _, item := range arr.Media {
		if !item.Matched {
			item.Unmatched = player

			result.Unmatched = append(result.Unmatched, item)
		}
	}
}

func findMatches(arr *sources.Source, player *sources.Source, exceptions []string, result *Result) {
	if len(player.Multiples) > 1 {
		result.MultipleFolders = append(result.MultipleFolders, player.Multiples...)
	}

	for key, media := range player.Media {
		outcome, match := media.Match(key, arr.Media)

		if match != nil {
			match.Matched = true
		}

		if slices.Contains(exceptions, media.Path) {
			continue
		}

		switch outcome {
		case sources.GuidMismatch:
			result.GuidMismatch = append(result.GuidMismatch, []*sources.Media{match, media})
		case sources.NoMatch:
			result.Unmatched = append(result.Unmatched, media)
		case sources.Match:
		}
	}
}
