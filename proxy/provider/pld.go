package provider

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"

	"shitty.moe/satelit-project/satelit-scraper/proxy"
)

// proxy-list.download
type PLD func(proxy.Protocol) ([]proxy.Proxy, error)

func NewPLD() PLD {
	return fetchPLD
}

func (p PLD) Fetch(proto proxy.Protocol) ([]proxy.Proxy, error) {
	return p(proto)
}

func fetchPLD(proto proxy.Protocol) ([]proxy.Proxy, error) {
	url := fmt.Sprintf("https://www.proxy-list.download/api/v1/get?type=%s", proto.String())
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	val, err := parsePLD(resp.Body, proto)
	_ = resp.Body.Close()

	return val, err
}

func parsePLD(buf io.Reader, proto proxy.Protocol) ([]proxy.Proxy, error) {
	proxies := make([]proxy.Proxy, 0)
	reader := bufio.NewReader(buf)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		splits := strings.Split(line, ":")
		if len(splits) != 2 {
			continue
		}

		prx := proxy.Proxy{
			Host:  strings.TrimSpace(splits[0]),
			Port:  strings.TrimSpace(splits[1]),
			Proto: proto,
		}

		proxies = append(proxies, prx)
	}

	return proxies, nil
}
