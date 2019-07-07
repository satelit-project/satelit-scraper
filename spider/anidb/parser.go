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
	return 0
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

func (p *Parser) StartDate() int64 {
	raw := p.doc.Find(`"div.g_definitionlist tr.year span[itemprop="startDate"]"`).First().AttrOr("content", "")
	d, err := parseDate(raw)
	if err != nil {
		log.Errorf("failed to parse start date: %v", err)
		return 0
	}

	return d
}

func (p *Parser) EndDate() int64 {
	raw := p.doc.Find(`"div.g_definitionlist tr.year span[itemprop="endDate"]"`).First().AttrOr("content", "")
	d, err := parseDate(raw)
	if err != nil {
		log.Errorf("failed to parse end date: %v", err)
		return 0
	}

	return d
}

func (p *Parser) Tags() []*scraper.Anime_Tag {
	tags := make([]*scraper.Anime_Tag, 0)
	p.doc.Find("div.g_definitionlist tr.tags span.g_tag").Each(func(_ int, s *goquery.Selection) {
		tag := new(scraper.Anime_Tag)
		name, id := parseTagName(s), parseTagID(s)
		if len(name) == 0 || id < 0 {
			return
		}

		tag.Name = name
		tag.Source = &scraper.Anime_Tag_AnidbId{AnidbId: id}
		tag.Description = parseTagInfo(s)

		tags = append(tags, tag)
	})

	return tags
}

func (p *Parser) Rating() float64 {
	raw := p.doc.Find("div.g_definitionlist tr.tmprating span.value").Text()
	raw = strings.TrimSpace(raw)

	r, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		log.Errorf("failed to parse rating: %v", err)
		return 0
	}

	return r
}

func (p *Parser) Description() string {
	raw := p.doc.Find(`div.desc[itemprop="description"]`).First().Text()
	raw = strings.TrimSpace(raw)

	// remove '* Based on ...' line
	rg := regexp.MustCompile(`(?i)^\*\s?based\son.+$`)
	raw = rg.ReplaceAllString(raw, "")

	// remove 'Source: ...' line
	rg = regexp.MustCompile(`(?i)^source:\s?.+$`)
	raw = rg.ReplaceAllString(raw, "")

	// remove 'Note: ...' line
	rg = regexp.MustCompile(`(?i)^note:\s?.+$`)
	raw = rg.ReplaceAllString(raw, "")

	// reformat
	rg = regexp.MustCompile(`\n\n+`)
	raw = strings.Join(rg.Split(raw, -1), "\n\n")

	return strings.TrimSpace(raw)
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
	d, err := parseDate(raw)
	if err != nil {
		log.Errorf("%v: %v", err, s)
		return 0
	}

	return d
}

func parseDate(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return 0, &Error{"parse date is empty"}
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}

func parseTagName(s *goquery.Selection) string {
	raw := s.Find("span.tagname").First().Text()
	return strings.TrimSpace(raw)
}

func parseTagInfo(s *goquery.Selection) string {
	raw := s.Find("span.text").First().Text()

	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "\"")
	raw = strings.TrimSuffix(raw, "\"")
	raw = strings.TrimSpace(raw)

	return raw
}

func parseTagID(s *goquery.Selection) int32 {
	raw := s.Find("a.tooltip").First().AttrOr("href", "")
	match := regexp.MustCompile(`tagid=(\d+)`).FindStringSubmatch(raw)
	if len(match) != 1 {
		log.Errorf("multiple tag ids found: %v", match)
		return -1
	}

	id, err := strconv.Atoi(match[0])
	if err != nil {
		log.Errorf("failed to parse tag id: %v", err)
		return -1
	}

	return int32(id)
}
