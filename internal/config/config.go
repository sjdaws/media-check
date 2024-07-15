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
	Exceptions []string     `yaml:"exceptions"`
	Sources    MatchSources `yaml:"sources"`
}

type MatchSources struct {
	Plex   MatchSourcesPlex `yaml:"plex"`
	Radarr MatchSourcesArr  `yaml:"radarr"`
	Sonarr MatchSourcesArr  `yaml:"sonarr"`
}

type MatchSourcesPlex struct {
	Database string                   `yaml:"database"`
	Sections MatchSourcesPlexSections `yaml:"sections"`
}

type MatchSourcesPlexSections struct {
	Movies  []int `yaml:"movies"`
	TVShows []int `yaml:"tvshows"`
}

type MatchSourcesArr struct {
	ApiKey  string   `yaml:"apikey"`
	PathMap []string `yaml:"pathmap"`
	Server  string   `yaml:"server"`
}

type Notify struct {
	Urls []string `yaml:"urls"`
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
