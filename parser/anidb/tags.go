package anidb

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"shitty.moe/satelit-project/satelit-scraper/proto/data"
)

// Parses and returns a list of tags or empty slice if tags not found.
func (p *Parser) tags() []*data.Anime_Tag {
	tags := make([]*data.Anime_Tag, 0)
	p.doc.Find("div.g_definitionlist tr.tags span.g_tag").Each(func(_ int, s *goquery.Selection) {
		var tag data.Anime_Tag
		name := parseTagName(s)
		if len(name) == 0 {
			p.log.Errorf("failed to parse tag name: %v", tag)
			return
		}

		id, err := parseTagID(s)
		if err != nil {
			p.log.Errorf("failed to parse tag id: %v", err)
			return
		}

		tag.Name = name
		tag.Source = &data.Anime_Tag_AnidbId{AnidbId: id}
		tag.Description = parseTagInfo(s)

		tags = append(tags, &tag)
	})

	return tags
}

// Parses and returns tag name or empty string if name not found.
func parseTagName(s *goquery.Selection) string {
	raw := s.Find("span.tagname").First().Text()
	return strings.TrimSpace(raw)
}

// Parses and returns tag description or empty string if info not found.
func parseTagInfo(s *goquery.Selection) string {
	raw := s.Find("span.text").First().Text()

	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "\"")
	raw = strings.TrimSuffix(raw, "\"")
	raw = strings.TrimSpace(raw)

	return raw
}

// Parses and returns tag id or error if tag id not found.
func parseTagID(s *goquery.Selection) (int32, error) {
	raw := s.Find("a.tooltip").First().AttrOr("href", "")
	match := regexp.MustCompile(`/(\d+)/`).FindStringSubmatch(raw)
	if len(match) != 2 {
		return 0, fmt.Errorf("not tag id: %v", raw)
	}

	id, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, fmt.Errorf("parsing failed: %v", err)
	}

	return int32(id), nil
}
