package notifier

import (
	"log"
	"strings"

	"github.com/containrrr/shoutrrr"
)

type Notifier struct {
	urls []string
}

func New(urls []string) *Notifier {
	return &Notifier{
		urls: urls,
	}
}

func (n *Notifier) Message(text string) {
	for _, url := range n.urls {
		url = strings.TrimSpace(url)

		if url == "" {
			continue
		}

		err := shoutrrr.Send(url, text)
		if err != nil {
			log.Printf("notify: %v", err)
		}
	}
}
