package single_api

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
		if proxyType == "http" {
			dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
			if err != nil {
				continue
			}
			proxyDialers = append(proxyDialers, dialer)
		} else {
			dialer, err := proxy.SOCKS5("tcp", proxyURL.Host, nil, proxy.Direct)
			if err != nil {
				continue
			}
			proxyDialers = append(proxyDialers, dialer)
		}

	}

	return proxyDialers
}
