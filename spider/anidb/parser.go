package anidb

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"satelit-project/satelit-scraper/proto/scraper"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

type Error struct {
	reason string
}

func (e *Error) Error() string {
	return e.reason
}

type Parser struct {
	url *url.URL
	doc *goquery.Document
}

func NewParser(url *url.URL, html io.Reader) (*Parser, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, err
	}

	p := Parser{url, doc}
	return &p, nil
}

func (p *Parser) Source() *scraper.Anime_Source {
	var source scraper.Anime_Source

	id, err := parseSource(p.url.String(), "aid=")
	if err != nil {
		log.Errorf("anidb id is malformed: %v", err)
		return nil
	}

	source.AnidbId = append(source.AnidbId, id)

	p.doc.Find("div.g_definitionlist tr.resources a[href*=\"myanimelist\"]").Each(func(_ int, s *goquery.Selection) {
		id, err := parseSource(s.AttrOr("href", ""), "/")
		if err != nil {
			log.Warnf("mal id is malformed: %v", err)
		}

		source.MalId = append(source.MalId, id)
	})

	p.doc.Find("div.g_definitionlist tr.resources a[href*=\"animenewsnetwork\"]").Each(func(_ int, s *goquery.Selection) {
		id, err := parseSource(s.AttrOr("href", ""), "id=")
		if err != nil {
			log.Warnf("ann id is malformed: %v", err)
		}

		source.AnnId = append(source.AnnId, id)
 	})

	return &source
}

func (p *Parser) Type() scraper.Anime_Type {
	raw := p.doc.Find("div.g_definitionlist tr.type td.value").First().Text()
	raw = strings.ToLower(raw)

	switch {
	case regexp.MustCompile("tv\\s+series").MatchString(raw):
		return scraper.Anime_TV_SERIES

	case regexp.MustCompile("ova").MatchString(raw):
		return scraper.Anime_OVA

	case regexp.MustCompile("web").MatchString(raw):
		return scraper.Anime_ONA

	case regexp.MustCompile("movie").MatchString(raw):
		return scraper.Anime_MOVIE

	case regexp.MustCompile("tv\\s+special").MatchString(raw):
		return scraper.Anime_SPECIAL

	default:
		return scraper.Anime_UNKNOWN
	}
}

func (p *Parser) Title() string {
	raw := p.doc.Find("div.g_definitionlist tr.romaji td span").First().Text()
	return strings.TrimSpace(raw)
}

func (p *Parser) PosterURL() string {
	raw := p.doc.Find("div.image picture img").First().AttrOr("src", "")
	return strings.TrimSpace(raw)
}

func (p *Parser) EpisodesCount() int32 {
	row := p.doc.Find("div.g_definitionlist tr.type td.value").First()
	prop := strings.TrimSpace(row.Find("span[itemprop=\"numberOfEpisodes\"]").Text())
	if ep, err := strconv.Atoi(prop); err == nil && len(prop) > 0 {
		return int32(ep)
	} else if err != nil {
		log.Warnf("failed to parse itemprop=\"numberOfEpisodes\"]: %v", err)
	}

	// try to parse from row's text
	raw := row.Text()
	raw = strings.ToLower(raw)

	// number after comma or space, usually for TV type
	match := regexp.MustCompile(",?\\s*(\\d+)").FindStringSubmatch(raw)
	if len(match) > 1 {
		log.Warnf("EpisodesCount() regexp found multiple ep numbers: %v", match)
	} else if len(match) == 1 {
		if ep, err := strconv.Atoi(match[0]); err == nil {
			return int32(ep)
		} else {
			log.Warnf("failed to parse episodes count: %v", err)
		}
	}

	// no comma and no numbers, but some text, usually for MOVIE
	match = regexp.MustCompile("^((?![,\\d]).)+$").FindStringSubmatch(raw)
	if len(match) > 0 {
		return 1
	}

	// has comma or space and some text but no numbers, usually number of ep is unknown
	match = regexp.MustCompile("^\\D+,?((?![\\d]).)+$").FindStringSubmatch(raw)
	if len(match) > 0 {
		return 0
	}

	return 0
}

func parseSource(str string, sep string) (int32, error) {
	raw := strings.Split(str, sep)
	if len(raw) == 0 {
		return 0, &Error{fmt.Sprintf("'%v'empty after splitting at %v", str, sep)}
	}

	s, err := strconv.Atoi(raw[len(raw) - 1])
	if err != nil {
		return 0, &Error{fmt.Sprintf("source is not an int: %v", s)}
	}

	return int32(s), nil
}
