package matcher

import (
	"slices"
	"strings"

	"github.com/sjdaws/media-check/internal/config"
	"github.com/sjdaws/media-check/internal/matcher/sources"
	"github.com/sjdaws/media-check/internal/matcher/sources/arr"
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

	plexMedia, err := plex.Fetch(cfg)
	if err != nil {
		return result, err
	}

	findMatches(arrMedia, plexMedia, cfg.Exceptions, &result)
	addUnmatched(arrMedia, &result)

	return result, nil
}

func addUnmatched(arr []*sources.Media, result *Result) {
	for _, item := range arr {
		if !item.Matched {
			result.Unmatched = append(result.Unmatched, item)
		}
	}
}

func findMatches(arr []*sources.Media, plex []*sources.Media, exceptions []string, result *Result) {
	for _, item := range plex {
		if len(item.Folders) < 1 {
			result.NoFolders = append(result.NoFolders, item)

			continue
		}

		if len(item.Folders) > 1 {
			result.MultipleFolders = append(result.MultipleFolders, item)

			continue
		}

		outcome, match := item.Match(arr)

		if match != nil {
			match.Matched = true
		}

		if slices.Contains(exceptions, strings.TrimSuffix(item.Folders[0], "/")) {
			continue
		}

		switch outcome {
		case sources.GuidMismatch:
			result.GuidMismatch = append(result.GuidMismatch, []*sources.Media{match, item})
		case sources.NoMatch:
			result.Unmatched = append(result.Unmatched, item)
		case sources.Match:
		}
	}
}
