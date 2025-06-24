package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Debug  string `yaml:"debug"`
	Find   Find   `yaml:"find"`
	Match  Match  `yaml:"match"`
	Notify Notify `yaml:"notify"`
}

type Find struct {
	Block      []string `yaml:"block"`
	Exceptions []string `yaml:"exceptions"`
	Paths      []string `yaml:"paths"`
	Type       string   `yaml:"type"`
}

type Match struct {
	Sources MatchSources `yaml:"sources"`
}

type MatchSources struct {
	Jellyfin MatchSourcesJellyfin `yaml:"jellyfin"`
	Plex     MatchSourcesPlex     `yaml:"plex"`
	Radarr   MatchSourcesArr      `yaml:"radarr"`
	Sonarr   MatchSourcesArr      `yaml:"sonarr"`
}

type MatchSourcesJellyfin struct {
	Database   string   `yaml:"database"`
	Exceptions []string `yaml:"exceptions"`
	Paths      Paths    `yaml:"paths"`
}

type MatchSourcesPlex struct {
	Database   string                   `yaml:"database"`
	Exceptions []string                 `yaml:"exceptions"`
	Paths      Paths                    `yaml:"paths"`
	Sections   MatchSourcesPlexSections `yaml:"sections"`
}

type MatchSourcesPlexSections struct {
	Movies  []int `yaml:"movies"`
	TVShows []int `yaml:"tvshows"`
}

type MatchSourcesArr struct {
	ApiKey string   `yaml:"apikey"`
	Paths  []string `yaml:"paths"`
	Server string   `yaml:"server"`
}

type Notify struct {
	Urls []string `yaml:"urls"`
}

type Paths struct {
	Movies  []string `yaml:"movies"`
	TVShows []string `yaml:"tvshows"`
}

func Load(filename string) (Config, error) {
	config := Config{}

	content, err := os.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("unable to read configuration file: %v", err)
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return config, fmt.Errorf("unable to parse configuration file: %v", err)
	}

	return config, nil
}
