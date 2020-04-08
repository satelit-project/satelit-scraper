package proxy

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"shitty.moe/satelit-project/satelit-scraper/logging"
)

// Represents a type that can return list of proxy servers.
type Provider interface {
	fmt.Stringer

	Fetch(proto Protocol) ([]Proxy, error)
}

// Type that can fetch available proxy servers.
type Fetcher struct {
	provider Provider
	limit    int
	proto    Protocol
	client   http.Client
	log      *logging.Logger
}

// Creates new Fetcher instance for fetching proxies with support for given protocol.
func NewFetcher(provider Provider, limit int, proto Protocol, log *logging.Logger) Fetcher {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	return Fetcher{
		provider: provider,
		limit:    limit,
		proto:    proto,
		client:   client,
		log:      log.With("proxy", "fetch"),
	}
}

// Sets timeout for proxy to determine wherever it's available for usage.
func (f *Fetcher) SetProxyTimeout(timeout time.Duration) {
	f.client.Timeout = timeout
}

// Fetches and returns list of available proxy servers or empty slice if no proxies available.
func (f *Fetcher) Fetch() []Proxy {
	f.log.Infof("fetching proxies from %v", f.provider)
	list, err := f.provider.Fetch(f.proto)
	if err != nil {
		f.log.Errorf("failed to fetch proxies: %v", err)
		return make([]Proxy, 0)
	}

	shuffle(list)

	f.log.Infof("checking proxies")
	return f.filter(list, f.limit)
}

// Returns at least first `limit` live proxies.
func (f *Fetcher) filter(arr []Proxy, limit int) []Proxy {
	var filtered = make([]Proxy, 0, limit)

	for round := 0; len(filtered) < limit && round * limit < len(arr); round += 1 {
		start, end := round * limit, (round + 1) * limit
		if end > len(arr) {
			end = len(arr)
		}

		ch := make(chan Proxy)
		for _, pr := range arr[start:end] {
			go func(proxy Proxy) {
				if !proxy.isAvailable(&f.client) {
					f.log.Debugf("proxy %v is not reachable", proxy)
					ch <- Proxy{}
				}

				ch <- proxy
			}(pr)
		}

		for _, _ = range arr[start:end] {
			pr := <-ch
			if pr.IsValid() {
				f.log.Infof("found live proxy: %v", pr)
				filtered = append(filtered, pr)
			}
		}
	}

	return filtered
}

// Shuffles list of proxies.
func shuffle(arr []Proxy) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
}

// srsly?
func min(x, y int) int {
	if x < y {
		return x
	}

	return y
}
