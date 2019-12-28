package proxy

import (
	"math/rand"
	"net/http"
	"time"

	"shitty.moe/satelit-project/satelit-scraper/logging"
)

// Represents a type that can return list of proxy servers.
type Provider interface {
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
	f.log.Infof("fetching proxies")
	list, err := f.provider.Fetch(f.proto)
	if err != nil {
		f.log.Errorf("failed to fetch proxies: %v", err)
		return make([]Proxy, 0)
	}

	shuffle(list)

	f.log.Infof("checking proxies")
	live := make([]Proxy, 0, min(f.limit, len(list)))
	for i := 0; i < len(list) && len(live) < f.limit; i++ {
		proxy := list[i]
		if !proxy.isAvailable(&f.client) {
			f.log.Infof("proxy %v is not reachable", proxy)
			continue
		}

		f.log.Infof("found live proxy: %v", proxy)
		live = append(live, proxy)
	}

	return live
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
