package proxy

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"go.uber.org/zap"

	"shitty.moe/satelit-project/satelit-scraper/logging"
)

type Provider interface {
	fmt.Stringer
	Fetch(proto Protocol) ([]Proxy, error)
}

type Fetcher struct {
	provider Provider
	limit    int
	proto    Protocol
	client   *http.Client
	log      *zap.SugaredLogger
}

func NewFetcher(provider Provider, limit int, proto Protocol) *Fetcher {
	log := logging.DefaultLogger().With("provider", provider.String())
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &Fetcher{
		provider: provider,
		limit:    limit,
		proto:    proto,
		client:   client,
		log:      log,
	}
}

func (f *Fetcher) SetProxyTimeout(timeout time.Duration) {
	f.client.Timeout = timeout
}

func (f *Fetcher) Fetch() []Proxy {
	f.log.Info("fetching proxies")
	list, err := f.provider.Fetch(f.proto)
	if err != nil {
		f.log.Errorf("failed to fetch proxies: %v", err)
		return make([]Proxy, 0)
	}

	f.log.Infof("checking proxies")
	live := make([]Proxy, 0, min(f.limit, len(list)))
	for i := 0; i < len(list) && len(live) < f.limit; i++ {
		proxy := list[i]
		if !proxy.isAvailable(f.client) {
			f.log.Infof("proxy %v is not reachable", proxy)
			continue
		}

		live = append(live, proxy)
	}

	shuffle(live)
	return live
}

func (p *Proxy) isAvailable(client *http.Client) bool {
	_, err := client.Head(p.String())
	return err == nil
}

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
