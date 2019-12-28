package main

import (
	"flag"
	"strconv"
)

type AnimeIDs []int32

func (a *AnimeIDs) String() string {
	return "Anime ID to scrape. Multiple IDs are allowed by specifying the flag multiple times."
}

func (a *AnimeIDs) Set(v string) error {
	i, err := strconv.Atoi(v)
	if err != nil {
		return err
	}

	*a = append(*a, int32(i))
	return nil
}

func AnimeIDsFlag(name string, value AnimeIDs, usage string) *AnimeIDs {
	flag.Var(&value, name, usage)
	return &value
}
