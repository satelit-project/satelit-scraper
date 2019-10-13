package anidb

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"

	"satelit-project/satelit-scraper/logging"
	"satelit-project/satelit-scraper/proto/data"
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
	log *zap.SugaredLogger
}

func NewParser(url *url.URL, html io.Reader) (*Parser, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, err
	}

	log := logging.DefaultLogger().With("db", "anidb")
	if id, err := parseSource(url.String(), "aid="); err == nil {
		log = log.With("id", id)
	}

	p := Parser{url, doc, log}
	return &p, nil
}

func (p *Parser) Anime() (*data.Anime, error) {
	anime := data.Anime{
		Source:        p.Source(),
		Type:          p.Type(),
		Title:         p.Title(),
		PosterUrl:     p.PosterURL(),
		EpisodesCount: p.EpisodesCount(),
		Episodes:      p.Episodes(),
		StartDate:     p.StartDate(),
		EndDate:       p.EndDate(),
		Tags:          p.Tags(),
		Rating:        p.Rating(),
		Description:   p.Description(),
	}

	if anime.Source == nil {
		return nil, &Error{"will not create Anime because Source is not valid"}
	}

	return &anime, nil
}

func (p *Parser) Source() *data.Anime_Source {
	var source data.Anime_Source

	id, err := parseSource(p.url.String(), "aid=")
	if err != nil {
		p.log.Warnf("anidb id is malformed: %v", err)
		return nil
	}

	source.AnidbIds = append(source.AnidbIds, id)

	p.doc.Find(`div.g_definitionlist tr.resources a[href*="myanimelist"]`).Each(func(_ int, s *goquery.Selection) {
		id, err := parseSource(s.AttrOr("href", ""), "/")
		if err != nil {
			p.log.Warnf("mal id is malformed: %v", err)
		}

		source.MalIds = append(source.MalIds, id)
	})

	p.doc.Find(`div.g_definitionlist tr.resources a[href*="animenewsnetwork"]`).Each(func(_ int, s *goquery.Selection) {
		id, err := parseSource(s.AttrOr("href", ""), "id=")
		if err != nil {
			p.log.Warnf("ann id is malformed: %v", err)
		}

		source.AnnIds = append(source.AnnIds, id)
	})

	return &source
}

func (p *Parser) Type() data.Anime_Type {
	raw := p.doc.Find("div.g_definitionlist tr.type td.value").First().Text()
	raw = strings.ToLower(raw)

	switch {
	case regexp.MustCompile(`tv\s+series`).MatchString(raw):
		return data.Anime_TV_SERIES

	case regexp.MustCompile(`ova`).MatchString(raw):
		return data.Anime_OVA

	case regexp.MustCompile(`web`).MatchString(raw):
		return data.Anime_ONA

	case regexp.MustCompile(`movie`).MatchString(raw):
		return data.Anime_MOVIE

	case regexp.MustCompile(`tv\s+special`).MatchString(raw):
		return data.Anime_SPECIAL

	default:
		return data.Anime_UNKNOWN
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
		p.log.Infof("ep number is not directly specified: %v", err)
	}

	// try to parse from row's text
	raw := row.Text()
	raw = strings.ToLower(raw)

	count, err := parseEpisodesCount(raw)
	if err != nil {
		p.log.Errorf("failed to parse ep count: %v", err)
		return 0
	}

	return count
}

func (p *Parser) Episodes() []*data.Episode {
	eps := make([]*data.Episode, 0)
	p.doc.Find(`table#eplist tr[id*="eid"]`).Each(func(i int, s *goquery.Selection) {
		log := p.log.With("ep_idx", i)
		ep := new(data.Episode)
		ep.Type = parseEpisodeType(s)
		if ep.Type == data.Episode_UNKNOWN {
			log.Info("unknown episode type")
			return
		}

		number, err := parseEpisodeNumber(s)
		if err != nil {
			log.Warnf("failed to parse ep number: %v", err)
		}

		name, err := parseEpisodeName(s)
		if err != nil {
			log.Warnf("failed to parse ep name: %v", err)
		}

		duration, err := parseEpisodeDuration(s)
		if err != nil {
			log.Warnf("failed to parse ep duration: %v", err)
		}

		date, err := parseEpisodeDate(s)
		if err != nil {
			log.Warnf("failed to parse ep air date: %v", err)
		}

		ep.Number = number
		ep.Name = name
		ep.Duration = duration
		ep.AirDate = date

		eps = append(eps, ep)
	})

	return eps
}

func (p *Parser) StartDate() int64 {
	raw := p.doc.Find(`div.g_definitionlist tr.year`).First()

	prop := raw.Find(` span[itemprop="startDate"]`).AttrOr("content", "")
	d, err := parseDate(prop)
	if err == nil {
		return d
	}

	p.log.Infof("startDate prop not found: %v", err)

	prop = raw.Find(`span[itemprop="datePublished"]`).AttrOr("content", "")
	d, err = parseDate(prop)
	if err == nil {
		return d
	}

	p.log.Infof("datePublished prop not found: %v", err)

	prop = raw.Find("td.value").Text()
	d, _, err = parseRawAirDate(prop)
	if err == nil {
		return d
	}

	p.log.Infof("raw air date not found: %v", err)
	p.log.Warnf("failed to parse start date: %v", raw.Text())

	return 0
}

func (p *Parser) EndDate() int64 {
	raw := p.doc.Find(`div.g_definitionlist tr.year`).First()

	prop := raw.Find(` span[itemprop="endDate"]`).AttrOr("content", "")
	d, err := parseDate(prop)
	if err == nil {
		return d
	}

	p.log.Infof("endDate prop not found: %v", err)

	prop = raw.Find("td.value").Text()
	_, d, err = parseRawAirDate(prop)
	if err == nil {
		return d
	}

	p.log.Infof("raw air date not found: %v", err)

	prop = raw.Find(`span[itemprop="datePublished"]`).AttrOr("content", "")
	d, err = parseDate(prop)
	if err == nil {
		return d
	}

	p.log.Infof("datePublished prop not found: %v", err)
	p.log.Warnf("failed to parse end date: %v", raw.Text())

	return 0
}

func (p *Parser) Tags() []*data.Anime_Tag {
	tags := make([]*data.Anime_Tag, 0)
	p.doc.Find("div.g_definitionlist tr.tags span.g_tag").Each(func(_ int, s *goquery.Selection) {
		tag := new(data.Anime_Tag)
		name := parseTagName(s)
		if len(name) == 0 {
			return
		}

		id, err := parseTagID(s)
		if err != nil {
			p.log.Warnf(" failed to parse tag id: %v", err)
			return
		}

		tag.Name = name
		tag.Source = &data.Anime_Tag_AnidbId{AnidbId: id}
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
		p.log.Warnf("failed to parse rating: %v", err)
		return 0
	}

	return r
}

func (p *Parser) Description() string {
	raw := p.doc.Find(`div.desc[itemprop="description"]`).First().Text()
	raw = strings.TrimSpace(raw)

	// remove '* Based on ...' line
	rg := regexp.MustCompile(`(?mi)^\*\s?based\son.+$`)
	raw = rg.ReplaceAllString(raw, "")

	// remove 'Source: ...' line
	rg = regexp.MustCompile(`(?mi)^source:\s?.+$`)
	raw = rg.ReplaceAllString(raw, "")

	// remove 'Note: ...' line
	rg = regexp.MustCompile(`(?mi)^note:\s?.+$`)
	raw = rg.ReplaceAllString(raw, "")

	// reformat
	rg = regexp.MustCompile(`\n\n+`)
	raw = strings.Join(rg.Split(raw, -1), "\n\n")

	return strings.TrimSpace(raw)
}

func parseSource(str string, sep string) (int32, error) {
	raw := strings.Split(str, sep)
	if len(raw) == 0 {
		return 0, &Error{fmt.Sprintf("'%v' is not a source", str)}
	}

	s, err := strconv.Atoi(raw[len(raw)-1])
	if err != nil {
		return 0, &Error{fmt.Sprintf("not an int: %v", s)}
	}

	return int32(s), nil
}

func parseEpisodesCount(raw string) (int32, error) {
	// number after comma, usually for TV type
	match := regexp.MustCompile(`,\s*(\d+)`).FindStringSubmatch(raw)
	if len(match) > 0 {
		ep, err := strconv.Atoi(match[1])
		return int32(ep), err
	}

	// no comma, numbers and questionmark, but some text, usually for 1ep titles
	match = regexp.MustCompile(`^([^,\d?])+$`).FindStringSubmatch(raw)
	if len(match) > 0 {
		return 1, nil
	}

	// probably has comma and some text but no numbers, usually number of ep is unknown
	match = regexp.MustCompile(`^\D+,?([^\d])+$`).FindStringSubmatch(raw)
	if len(match) > 0 {
		return 0, nil
	}

	return 0, &Error{fmt.Sprintf("failed to parse episode count from %v", raw)}
}

func parseEpisodeType(s *goquery.Selection) data.Episode_Type {
	raw := s.Find("td abbr").First().AttrOr("title", "")
	raw = strings.TrimSpace(strings.ToLower(raw))

	switch {
	case regexp.MustCompile(`regular`).MatchString(raw):
		return data.Episode_REGULAR

	case regexp.MustCompile(`special`).MatchString(raw):
		return data.Episode_SPECIAL

	default:
		return data.Episode_UNKNOWN
	}
}

func parseEpisodeNumber(s *goquery.Selection) (int32, error) {
	raw := s.Find("td.eid abbr").First().Text()
	match := regexp.MustCompile(`\d+`).FindStringSubmatch(raw)
	if len(match) == 0 {
		return 0, &Error{fmt.Sprintf("not found for: %v", raw)}
	}

	num, err := strconv.Atoi(match[0])
	if err != nil {
		return 0, &Error{fmt.Sprintf("not an int: %v", err)}
	}

	return int32(num), nil
}

func parseEpisodeName(s *goquery.Selection) (string, error) {
	raw := s.Find(`td.name label`).First().Text()

	// generic name like "Episode 1" should be skipped
	if regexp.MustCompile(`episode\s+[\d.]+`).MatchString(strings.ToLower(raw)) {
		return "", &Error{fmt.Sprintf("unnamed episode: %v", raw)}
	}

	return strings.TrimSpace(raw), nil
}

func parseEpisodeDuration(s *goquery.Selection) (float64, error) {
	raw := s.Find("td.duration").First().Text()
	raw = strings.TrimSpace(raw)
	if len(raw) == 0 {
		return 0, &Error{fmt.Sprintf("not found: %v", s)}
	}

	match := regexp.MustCompile(`(\d+)\s*m`).FindStringSubmatch(raw)
	if len(match) == 0 {
		return 0, &Error{fmt.Sprintf("not found: %v", raw)}
	}

	mins, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0, &Error{fmt.Sprintf("parsing failed: %v", err)}
	}

	return mins * 60, nil
}

func parseEpisodeDate(s *goquery.Selection) (int64, error) {
	raw := s.Find("td.airdate").First()
	d, err := parseDate(raw.AttrOr("content", ""))
	if err == nil {
		return d, nil
	}

	d, err = parseAltDate(raw.Text())
	if err == nil {
		return d, nil
	}

	return 0, &Error{fmt.Sprintf("parsing failed: %v", err)}
}

func parseDate(s string) (int64, error) {
	s = strings.TrimSpace(s)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}

func parseAltDate(s string) (int64, error) {
	s = strings.TrimSpace(s)
	t, err := time.Parse("02.01.2006", s)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}

func parseRawAirDate(s string) (int64, int64, error) {
	s = strings.TrimSpace(s)
	match := regexp.MustCompile(`(\d{2}.\d{2}.\d{4}).+(till|,).+(\d{2}.\d{2}.\d{4})`).FindStringSubmatch(s)
	if len(match) != 4 {
		return 0, 0, &Error{"raw air date unknown format"}
	}

	start, err := time.Parse("02.01.2006", match[1])
	if err != nil {
		return 0, 0, err
	}

	end, err := time.Parse("02.01.2006", match[3])
	if err != nil {
		return 0, 0, err
	}

	return start.Unix(), end.Unix(), nil
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

func parseTagID(s *goquery.Selection) (int32, error) {
	raw := s.Find("a.tooltip").First().AttrOr("href", "")
	match := regexp.MustCompile(`tagid=(\d+)`).FindStringSubmatch(raw)
	if len(match) == 0 {
		return 0, &Error{fmt.Sprintf("not tag id: %v", raw)}
	}

	id, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, &Error{fmt.Sprintf("parsing failed: %v", err)}
	}

	return int32(id), nil
}
