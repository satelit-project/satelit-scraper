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

// Parses and returns list of episodes or empty slice if episodes not found.
func (p *Parser) episodes() []*data.Episode {
	eps := make([]*data.Episode, 0)
	p.doc.Find(`table#eplist tr[id*="eid"]`).Each(func(i int, s *goquery.Selection) {
		log := p.log.With("ep_idx", i)
		var ep data.Episode

		ep.Type = parseEpisodeType(s)
		if ep.Type == data.Episode_UNKNOWN {
			log.Infof("skipping episode with unknown type")
			return
		}

		number, err := parseEpisodeNumber(s)
		if err != nil {
			log.Errorf("failed to parse ep number: %v", err)
		}

		name, err := parseEpisodeName(s)
		if err != nil {
			log.Errorf("failed to parse ep name: %v", err)
		}

		duration, err := parseEpisodeDuration(s)
		if err != nil {
			log.Errorf("failed to parse ep duration: %v", err)
		}

		date, err := parseEpisodeDate(s)
		if err != nil {
			log.Errorf("failed to parse ep air date: %v", err)
		}

		ep.Number = number
		ep.Name = name
		ep.Duration = duration.Seconds()
		ep.AirDate = date.Unix()
		eps = append(eps, &ep)
	})

	return eps
}

// Parses episode type of an episode. Episode_UNKNOWN returned on error.
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

// Parses episode number. Error is returned if edisode isn't numbered.
func parseEpisodeNumber(s *goquery.Selection) (int32, error) {
	raw := s.Find("td.eid abbr").First().Text()
	match := regexp.MustCompile(`\d+`).FindStringSubmatch(raw)
	if len(match) == 0 {
		return 0, errors.New(fmt.Sprintf("not found: %v", raw))
	}

	num, err := strconv.Atoi(match[0])
	if err != nil {
		return 0, err
	}

	return int32(num), nil
}

// Parses name of an episode. Empty string is returned edisode doesn't have a name.
func parseEpisodeName(s *goquery.Selection) (string, error) {
	raw := s.Find(`td.name label`).First().Text()

	// generic name like "Episode 1" should be skipped
	if regexp.MustCompile(`episode\s+[\d.]+`).MatchString(strings.ToLower(raw)) {
		return "", errors.New(fmt.Sprintf("generic episode: %v", raw))
	}

	return strings.TrimSpace(raw), nil
}

// Parses episode duration. Zero is returned if episode doesn't have duration.
func parseEpisodeDuration(s *goquery.Selection) (time.Duration, error) {
	raw := s.Find("td.duration").First().Text()
	raw = strings.TrimSpace(raw)
	if len(raw) == 0 {
		return 0, errors.New(fmt.Sprintf("not found: %v", s))
	}

	match := regexp.MustCompile(`(\d+)\s*m`).FindStringSubmatch(raw)
	if len(match) == 0 {
		return 0, errors.New(fmt.Sprintf("not found: %v", raw))
	}

	mins, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("parsing failed: %v", err))
	}

	return time.Duration(mins) * time.Minute, nil
}

// Parses episode air date. Returns zero if episode doesn't have air date.
func parseEpisodeDate(s *goquery.Selection) (time.Time, error) {
	raw := s.Find("td.airdate").First()
	d, err := parseDate(raw.AttrOr("content", raw.Text()))
	if err != nil {
		return time.Time{}, errors.New(fmt.Sprintf("parsing failed: %v", err))
	}

	return d, nil
}
