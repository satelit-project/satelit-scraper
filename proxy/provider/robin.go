package provider

import (
	"sync/atomic"

	"shitty.moe/satelit-project/satelit-scraper/proxy"
)

type roundRobinProvider struct {
	providers []proxy.Provider
	index     uint32
}

func NewRoundRobin(providers []proxy.Provider) proxy.Provider {
	return &roundRobinProvider{
		providers: providers,
		index:     0,
	}
}

func (r *roundRobinProvider) Fetch(proto proxy.Protocol) ([]proxy.Proxy, error) {
	index := atomic.AddUint32(&r.index, 1) % uint32(len(r.providers))
	provider := r.providers[index]
	return provider.Fetch(proto)
}

func (r *roundRobinProvider) String() string {
	index := (atomic.LoadUint32(&r.index) + 1) % uint32(len(r.providers))
	return r.providers[index].String()
}
