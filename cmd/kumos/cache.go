package main

import (
	"bytes"

	"shitty.moe/satelit-project/satelit-scraper/proto/data"
)

type noCache struct {}

func (c noCache) AddHTML(data *bytes.Buffer, source data.Source, id int32) error {
	return nil
}
