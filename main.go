package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/sjdaws/media-check/internal/config"
	"github.com/sjdaws/media-check/internal/finder"
	"github.com/sjdaws/media-check/internal/matcher"
	"github.com/sjdaws/media-check/pkg/notifier"
)

func main() {
	cfg, err := config.Load("/config/config.yaml")
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	notify := notifier.New(cfg.Notify.Urls)

	findResult := finder.Result{}

	matchResult, err := matcher.Check(cfg.Match)
	if err != nil {
		notify.Message(err.Error())
		log.Fatal(err)
	}

	logResult(findResult, matchResult, notify)
}

func logResult(findResult finder.Result, matchResult matcher.Result, notify *notifier.Notifier) {
	total := len(findResult.Found) + len(matchResult.GuidMismatch) + len(matchResult.MultipleFolders) + len(matchResult.NoFolders) + len(matchResult.Unmatched)
	if total < 1 {
		return
	}

	message := "Found " + strconv.Itoa(total) + " media check violations.\n"

	if len(findResult.Found) > 0 {
		message += "\n"
		message += "-------------------------------------\n"
		message += fmt.Sprintf("Found %d paths that contain blocks:\n", len(findResult.Found))
		message += "-------------------------------------\n"

		for _, violation := range findResult.Found {
			message += "\n"
			message += violation.Path + " contains " + violation.Term + "\n"
		}
		message += "\n"
	}

	if len(matchResult.GuidMismatch) > 0 {
		message += "\n"
		message += "-------------------------------------\n"
		message += fmt.Sprintf("Found %d items with mismatched GUIDs:\n", len(matchResult.GuidMismatch))
		message += "-------------------------------------\n"

		for _, unmatched := range matchResult.GuidMismatch {
			message += "\n"
			message += unmatched[0].Source + " | " + unmatched[1].Source + ":\n"
			message += fmt.Sprintf(" - Title:  %s (%d) | %s (%d)\n", unmatched[0].Title, unmatched[0].Year, unmatched[1].Title, unmatched[1].Year)
			message += " - Folder: " + unmatched[0].Path + " | " + unmatched[1].Path + "\n"
			message += " - IMDB:   " + unmatched[0].Guids.IMDB + " | " + unmatched[1].Guids.IMDB + "\n"
			message += " - TMDB:   " + strconv.Itoa(unmatched[0].Guids.TMDB) + " | " + strconv.Itoa(unmatched[1].Guids.TMDB) + "\n"
			message += " - TVDB:   " + strconv.Itoa(unmatched[0].Guids.TVDB) + " | " + strconv.Itoa(unmatched[1].Guids.TVDB) + "\n"
		}
		message += "\n"
	}

	if len(matchResult.MultipleFolders) > 0 {
		message += "\n"
		message += "-------------------------------------\n"
		message += fmt.Sprintf("Found %d items with multiple folders:\n", len(matchResult.MultipleFolders))
		message += "-------------------------------------\n"

		for _, mismatch := range matchResult.MultipleFolders {
			message += "\n"
			message += mismatch.Title + "\n"
			message += " - " + mismatch.Path + "\n"
		}
		message += "\n"
	}

	if len(matchResult.NoFolders) > 0 {
		message += "\n"
		message += "-------------------------------------\n"
		message += fmt.Sprintf("Found %d items with no folders:\n", len(matchResult.NoFolders))
		message += "-------------------------------------\n"

		for _, mismatch := range matchResult.NoFolders {
			message += "\n"
			message += mismatch.Title + "\n"
		}
		message += "\n"
	}

	if len(matchResult.Unmatched) > 0 {
		message += "\n"
		message += "-------------------------------------\n"
		message += fmt.Sprintf("Found %d unmatched items:\n", len(matchResult.Unmatched))
		message += "-------------------------------------\n"
		message += "\n"
		for _, unmatched := range matchResult.Unmatched {
			title := unmatched.Title
			if unmatched.Guids.IMDB != "" {
				title += " (" + unmatched.Guids.IMDB + ")"
			}

			message += fmt.Sprintf("%s via %s in %s", title, unmatched.Source, unmatched.Path)

			if unmatched.Unmatched != "" {
				message += fmt.Sprintf(" not found on %s", unmatched.Unmatched)
			}

			message += "\n"
		}
		message += "\n"
	}

	notify.Message(message)
	log.Print(message)
}
