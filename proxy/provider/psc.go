package provider

import (
	"fmt"
	"net/http"

	"shitty.moe/satelit-project/satelit-scraper/proxy"
)

// proxyscrape.com
type PSC func(proto proxy.Protocol) ([]proxy.Proxy, error)

func NewPSC() PSC {
	return fetchPSC
}

func (p PSC) Fetch(proto proxy.Protocol) ([]proxy.Proxy, error) {
	return p(proto)
}

func fetchPSC(proto proxy.Protocol) ([]proxy.Proxy, error) {
	url := fmt.Sprintf("https://api.proxyscrape.com/?request=getproxies&proxytype=%s"+
		"&timeout=5000&country=all&ssl=all&anonymity=all", proto.String())

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	val, err := parsePLD(resp.Body, proto)
	_ = resp.Body.Close()

	return val, err
}
