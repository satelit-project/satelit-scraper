package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

// Application configuration.
type Config struct {
	Serving  *Serving  `yaml:"serving"`
	Scraping *Scraping `yaml:"scraping"`
	AniDB    *AniDB    `yaml:"anidb"`
	Logging  Logging   `yaml:"logging"`
}

// Server configuration.
type Serving struct {
	// Port to listen for incoming connections.
	Port uint `yaml:"port"`

	// Timeout for graceful shutdown.
	HaltTimeout uint64 `yaml:"halt-timeout"`
}

// Scraping configuration.
type Scraping struct {
	// Address of scraping tasks service.
	TaskAddress string `yaml:"task-address"`

	// Timeout for reporting scraping progress to external service.
	ReportTimeout time.Duration `yaml:"report-timeout"`
}

// AniDB specific configuration.
type AniDB struct {
	// Template to create AniDB anime URL from anime ID.
	URLTemplate string `yaml:"url-template"`

	// Timeout for AniDB requests.
	Timeout uint64 `yaml:"timeout"`

	// Delay between AniDB requests.
	Delay uint64 `yaml:"delay"`
}

// Logging configuration.
type Logging struct {
	// Logging profile.
	Profile string `yaml:"profile"`
}

// Returns default app configuration or error if failed to read it.
func Default() (Config, error) {
	data := makeData(os.Environ())
	return AtPath("config/default.yml", data)
}

// Returns app configuration parsed from template with provided data.
func AtPath(path string, data map[string]string) (Config, error) {
	var cfg Config
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	content, err = render(content, data)
	if err != nil {
		return cfg, err
	}

	if err = yaml.Unmarshal(content, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// Renders template with provided data.
func render(cfg []byte, data map[string]string) ([]byte, error) {
	t, err := template.New("config").Parse(string(cfg))
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	err = t.Execute(&b, data)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// Maps provided environment varialbes into template data.
func makeData(env []string) map[string]string {
	data := make(map[string]string, 8)
	for _, env := range env {
		sp := strings.Split(env, "=")
		if len(sp) != 2 {
			continue
		}

		data[sp[0]] = sp[1]
	}

	return data
}
