package anidb

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"shitty.moe/satelit-project/satelit-scraper/proto/data"
)

// Parses and returns info about external DBs locations or nil on error.
func (p *Parser) source() *data.Anime_Source {
	var source data.Anime_Source

	id, err := parseSource(p.url.String(), "/")
	if err != nil {
		p.log.Errorf("anidb id is malformed: %v", err)
		return nil
	}

	source.AnidbIds = append(source.AnidbIds, id)

	p.doc.Find(`div.g_definitionlist tr.resources a[href*="myanimelist"]`).Each(func(_ int, s *goquery.Selection) {
		id, idErr := parseSource(s.AttrOr("href", ""), "/")
		if err != nil {
			err = fmt.Errorf("mal id is malformed: %v", idErr)
			return
		}

		source.MalIds = append(source.MalIds, id)
	})

	p.doc.Find(`div.g_definitionlist tr.resources a[href*="animenewsnetwork"]`).Each(func(_ int, s *goquery.Selection) {
		id, idErr := parseSource(s.AttrOr("href", ""), "id=")
		if err != nil {
			err = fmt.Errorf("ann id is malformed: %v", idErr)
			return
		}

		source.AnnIds = append(source.AnnIds, id)
	})

	if err != nil {
		// since mal and ann IDs are not required we can just log the error
		// and continue parsing
		p.log.Errorf("%v", err)
	}

	return &source
}

// Parses and returns anime entry type. In case of errors Anime_UNKNOWN will be returned.
func (p *Parser) animeType() data.Anime_Type {
	raw := p.doc.Find("div.g_definitionlist tr.type td.value").First().Text()
	raw = strings.ToLower(raw)

	// TODO: regexp as constants
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

// Parses and returns anime title or empty string if not found.
func (p *Parser) title() string {
	raw := p.doc.Find("div.g_definitionlist tr.romaji td span").First().Text()
	return strings.TrimSpace(raw)
}

// Parses and returns poster URL or empty string if not found.
func (p *Parser) posterURL() string {
	raw := p.doc.Find("div.image picture img").First().AttrOr("src", "")
	return strings.TrimSpace(raw)
}

// Parses end returns expected number of show's episode or zero if not found.
func (p *Parser) episodesCount() int32 {
	row := p.doc.Find("div.g_definitionlist tr.type td.value").First()
	prop := strings.TrimSpace(row.Find("span[itemprop=\"numberOfEpisodes\"]").Text())
	if ep, err := strconv.Atoi(prop); err == nil && len(prop) > 0 {
		return int32(ep)
	}

	// try to parse raw node's text
	raw := row.Text()
	count, err := parseRawEpisodesCount(raw)
	if err != nil {
		p.log.Errorf("failed to parse ep count: %v", err)
		return 0
	}

	return count
}

// Parses and returns show's start air date or zero if not found.
func (p *Parser) startDate() time.Time {
	raw := p.doc.Find(`div.g_definitionlist tr.year`).First()
	rawDate := raw.Find("td.value").Text()

	if !strings.Contains(rawDate, "?") {
		prop := raw.Find(` span[itemprop="startDate"]`).AttrOr("content", "")
		d, err := parseDate(prop)
		if err == nil {
			return d
		}

		prop = raw.Find(`span[itemprop="datePublished"]`).AttrOr("content", "")
		d, err = parseDate(prop)
		if err == nil {
			return d
		}
	}

	// try to parse raw node's text
	d, _, err := parseRawAirDate(rawDate)
	if err == nil {
		return d
	}

	p.log.Infof("start air date not found: %s", raw.Text())
	return time.Time{}
}

// Parses and returns show's end air date or zero if not found.
func (p *Parser) endDate() time.Time {
	raw := p.doc.Find(`div.g_definitionlist tr.year`).First()
	rawDate := raw.Find("td.value").Text()

	if !strings.Contains(rawDate, "?") {
		prop := raw.Find(` span[itemprop="endDate"]`).AttrOr("content", "")
		d, err := parseDate(prop)
		if err == nil {
			return d
		}
	}

	_, d, err := parseRawAirDate(rawDate)
	if err == nil {
		return d
	}

	// assuming that if there's "datePublished" at this point then it's a show
	// that aired for one day
	prop := raw.Find(`span[itemprop="datePublished"]`).AttrOr("content", "")
	d, err = parseDate(prop)
	if err == nil {
		return d
	}

	p.log.Infof("failed to parse end date: %v", raw.Text())
	return time.Time{}
}

// Parses and returns show's rating or zero if not found.
func (p *Parser) rating() float64 {
	raw := p.doc.Find("div.g_definitionlist tr.tmprating span.value").Text()
	raw = strings.TrimSpace(raw)
	if len(raw) == 0 {
		return 0
	}

	r, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		p.log.Errorf("failed to parse rating: %v", err)
		return 0
	}

	return r
}

// Parses and returns show's description or empty string if not found.
func (p *Parser) description() string {
	raw := p.doc.Find(`div.desc[itemprop="description"]`).First().Text()
	raw = strings.TrimSpace(raw)

	// remove '* ...' line
	rg := regexp.MustCompile(`(?m)^\*\s+.+$`)
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

// Parses and returns external source ID from external entry URL.
func parseSource(str string, sep string) (int32, error) {
	raw := strings.Split(str, sep)
	if len(raw) == 0 {
		return 0, fmt.Errorf("'%v' is not a source", str)
	}

	s, err := strconv.Atoi(raw[len(raw)-1])
	if err != nil {
		return 0, fmt.Errorf("not an int: %v", s)
	}

	return int32(s), nil
}

// Parses and returns expected show's episode count from AniDB page text or zero if not found.
// Error is returned in case if episode number exists but can't be parsed.
func parseRawEpisodesCount(raw string) (int32, error) {
	raw = strings.ToLower(raw)

	// TODO: regexp as constants
	// number after comma, usually for TV type
	match := regexp.MustCompile(`,\s*(\d+)`).FindStringSubmatch(raw)
	if len(match) > 0 {
		ep, err := strconv.Atoi(match[1])
		return int32(ep), err
	}

	// no comma, numbers or questionmark, but some text, usually for 1ep titles
	match = regexp.MustCompile(`^([^,\d?])+$`).FindStringSubmatch(raw)
	if len(match) > 0 {
		return 1, nil
	}

	// probably has comma and some text but no numbers, usually number of ep is unknown
	match = regexp.MustCompile(`^\D+,?([^\d])+$`).FindStringSubmatch(raw)
	if len(match) > 0 {
		return 0, nil
	}

	return 0, fmt.Errorf("failed to parse episode count from %v", raw)
}

// Parses and returns start and end air date from AniDB page text or zero if air date not found.
// Error is returned in case if air date found but can't be parsed.
func parseRawAirDate(s string) (time.Time, time.Time, error) {
	var zero time.Time

	s = strings.TrimSpace(s)
	match := regexp.MustCompile(`^([\d.]+)(\s*(till|,)(\s+[\d.]+)?)?`).FindStringSubmatch(s)
	if len(match) < 2 {
		return zero, zero, errors.New("raw air date unknown format")
	}

	match = match[1:]
	start, err := parseDate(match[0])
	if len(match) == 1 || len(match[1]) == 0 || err != nil {
		// only one date means it aired in one day
		return start, start, err
	}

	if len(match) < 4 || len(match[3]) == 0 {
		// means that end air date is unknown
		return start, zero, nil
	}

	end, err := parseDate(match[3])
	return start, end, err
}
