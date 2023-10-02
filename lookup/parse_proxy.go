package lookup

import (
	"fmt"
	"net/url"

	"golang.org/x/net/proxy"
)

func ProxyParser(proxies []string, proxyType string) []proxy.Dialer {
	var proxyDialers []proxy.Dialer
	for _, p := range proxies {
		urlStr := fmt.Sprintf("%s://%s", proxyType, p)
		proxyURL, _ := url.Parse(urlStr)
		dialer, err := proxy.SOCKS5("tcp", proxyURL.Host, nil, proxy.Direct)
		if err != nil {
			continue
		}
		proxyDialers = append(proxyDialers, dialer)
	}

	return proxyDialers
}
