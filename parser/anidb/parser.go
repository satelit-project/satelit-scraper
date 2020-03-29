package anidb

import (
	"errors"
	"io"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"shitty.moe/satelit-project/satelit-scraper/logging"
	"shitty.moe/satelit-project/satelit-scraper/proto/data"
)

// Parses AniDB anime html page.
type Parser struct {
	url *url.URL
	doc *goquery.Document
	log *logging.Logger
}

// Creates new AniDB anime page parser.
func NewParser(url *url.URL, html io.Reader, log *logging.Logger) (Parser, error) {
	var parser Parser
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return parser, err
	}

	log = log.With("parser", "anidb")
	if id, err := parseSource(url.String(), "/"); err == nil {
		log = log.With("id", id)
	}

	p := Parser{url, doc, log}
	p.url = url
	p.doc = doc
	p.log = log
	return p, nil
}

// Parses and returns Anime instance from AniDB anime page.
func (p *Parser) Anime() (*data.Anime, error) {
	p.log.Infof("parsing strarted")

	anime := data.Anime{
		Source:        p.source(),
		Type:          p.animeType(),
		Title:         p.title(),
		PosterUrl:     p.posterURL(),
		EpisodesCount: p.episodesCount(),
		Episodes:      p.episodes(),
		Tags:          p.tags(),
		Rating:        p.rating(),
		Description:   p.description(),
	}

	if sd := p.startDate(); !sd.IsZero() {
		anime.StartDate = sd.Unix()
	}

	if ed := p.endDate(); !ed.IsZero() {
		anime.EndDate = ed.Unix()
	}

	if anime.Source == nil {
		return nil, errors.New("will not create Anime because Source is not valid")
	}

	return &anime, nil
}
