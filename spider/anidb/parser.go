package anidb

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

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

	p.doc.Find(`div.g_definitionlist tr.resources a[href*="myanimelist"]`).Each(func(_ int, s *goquery.Selection) {
		id, err := parseSource(s.AttrOr("href", ""), "/")
		if err != nil {
			log.Errorf("mal id is malformed: %v", err)
		}

		source.MalId = append(source.MalId, id)
	})

	p.doc.Find(`div.g_definitionlist tr.resources a[href*="animenewsnetwork"]`).Each(func(_ int, s *goquery.Selection) {
		id, err := parseSource(s.AttrOr("href", ""), "id=")
		if err != nil {
			log.Errorf("ann id is malformed: %v", err)
		}

		source.AnnId = append(source.AnnId, id)
 	})

	return &source
}

func (p *Parser) Type() scraper.Anime_Type {
	raw := p.doc.Find("div.g_definitionlist tr.type td.value").First().Text()
	raw = strings.ToLower(raw)

	switch {
	case regexp.MustCompile(`tv\s+series`).MatchString(raw):
		return scraper.Anime_TV_SERIES

	case regexp.MustCompile(`ova`).MatchString(raw):
		return scraper.Anime_OVA

	case regexp.MustCompile(`web`).MatchString(raw):
		return scraper.Anime_ONA

	case regexp.MustCompile(`movie`).MatchString(raw):
		return scraper.Anime_MOVIE

	case regexp.MustCompile(`tv\s+special`).MatchString(raw):
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
		log.Errorf("failed to parse itemprop=\"numberOfEpisodes\"]: %v", err)
	}

	// try to parse from row's text
	raw := row.Text()
	raw = strings.ToLower(raw)

	// number after comma, usually for TV type
	match := regexp.MustCompile(`,\s*(\d+)`).FindStringSubmatch(raw)
	if len(match) > 1 {
		log.Warnf("EpisodesCount() regexp found multiple ep numbers: %v", match)
	} else if len(match) == 1 {
		if ep, err := strconv.Atoi(match[0]); err == nil {
			return int32(ep)
		} else {
			log.Errorf("failed to parse episodes count: %v", err)
		}
	}

	// no comma, numbers and questionmark, but some text, usually for 1ep titles
	match = regexp.MustCompile(`^([^,\d?])+$`).FindStringSubmatch(raw)
	if len(match) > 0 {
		return 1
	}

	// probably has comma and some text but no numbers, usually number of ep is unknown
	match = regexp.MustCompile(`^\D+,?([^\d])+$`).FindStringSubmatch(raw)
	if len(match) > 0 {
		log.Infof("unknown episode count for %v", p.url)
		return 0
	}

	log.Errorf("failed to find episode count for %v", p.url)
	return -1
}

func (p *Parser) Episodes() []*scraper.Episode {
	eps := make([]*scraper.Episode, 0)
	p.doc.Find(`table#eplist tr[itemprop="episode"]`).Each(func(_ int, s *goquery.Selection) {
		ep := new(scraper.Episode)
		ep.Type = parseEpisodeType(s)
		if ep.Type == scraper.Episode_UNKNOWN {
			log.Warnf("unknown episode type: %v", s.Text())
			return
		}

		ep.Number = parseEpisodeNumber(s)
		ep.Name = parseEpisodeName(s)
		ep.Duration = parseEpisodeDuration(s)
		ep.AirDate = parseEpisodeDate(s)

		eps = append(eps, ep)
	})

	return eps
}

func parseSource(str string, sep string) (int32, error) {
	raw := strings.Split(str, sep)
	if len(raw) == 0 {
		return -1, &Error{fmt.Sprintf("'%v'empty after splitting at %v", str, sep)}
	}

	s, err := strconv.Atoi(raw[len(raw) - 1])
	if err != nil {
		return -1, &Error{fmt.Sprintf("source is not an int: %v", s)}
	}

	return int32(s), nil
}

func parseEpisodeType(s *goquery.Selection) scraper.Episode_Type {
	raw := s.Find("td abbr").First().AttrOr("title", "")
	raw = strings.TrimSpace(strings.ToLower(raw))

	switch {
	case regexp.MustCompile(`regular\s+episode`).MatchString(raw):
		return scraper.Episode_REGULAR

	case regexp.MustCompile(`special`).MatchString(raw):
		return scraper.Episode_SPECIAL

	default:
		return scraper.Episode_UNKNOWN
	}
}

func parseEpisodeNumber(s *goquery.Selection) int32 {
	raw := s.Find("td abbr").First().Text()
	match := regexp.MustCompile(`\d+`).FindStringSubmatch(raw)
	if len(match) != 1 {
		log.Errorf("multiple episode numbers found: %v", raw)
		return -1
	}

	num, err := strconv.Atoi(match[0])
	if err != nil {
		log.Errorf("episode number is not an int: %v", err)
		return -1
	}

	return int32(num)
}

func parseEpisodeName(s *goquery.Selection) string {
	raw := s.Find("td.name").First().Text()

	// generic name like "Episode 1" should be skipped
	if regexp.MustCompile(`episode\s+[\d.]+`).MatchString(strings.ToLower(raw)) {
		log.Infof("unnamed episode: %v", s)
		return ""
	}

	return strings.TrimSpace(raw)
}

func parseEpisodeDuration(s *goquery.Selection) float64 {
	raw := s.Find("td.duration").First().Text()
	raw = strings.TrimSpace(raw)
	if len(raw) == 0 {
		log.Warnf("episode duration not found: %v", s)
		return 0
	}

	match := regexp.MustCompile(`(\d+)\s*m`).FindStringSubmatch(raw)
	if len(match) == 0 {
		log.Warnf("episode duration not found: %v", raw)
		return 0
	}

	mins, err := strconv.ParseFloat(match[0], 64)
	if err != nil {
		log.Errorf("failed to parse episode duration: %v", err)
		return 0
	}

	return mins * 60
}

func parseEpisodeDate(s *goquery.Selection) int64 {
	raw := s.Find("td.airdate").First().AttrOr("content", "")
	if len(raw) == 0 {
		log.Warnf("episode air date not found: %v", raw)
		return -1
	}

	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		log.Errorf("failed to parse episode air date: %v", raw)
		return -1
	}

	return t.Unix()
}
