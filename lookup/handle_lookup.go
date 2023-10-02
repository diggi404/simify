package lookup

import (
	"sync"

	"golang.org/x/net/proxy"
)

func ProcessLookup(numberChunks <-chan []string, wg *sync.WaitGroup, mutex *sync.Mutex, proxies *[]proxy.Dialer) {
	defer wg.Done()
	numbers := <-numberChunks
	for _, number := range numbers {
		for _, proxy := range *proxies {
			err := LookupAPI(proxy, number)
			if err != nil {
				continue
			}

		}
	}
}
