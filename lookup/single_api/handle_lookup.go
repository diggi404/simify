package single_api

import (
	"os"
	"simify/fileutil"
	"sync"

	"golang.org/x/net/proxy"
)

func ProcessLookup(numberChunks <-chan []string, wg *sync.WaitGroup, mutex *sync.Mutex, proxyChannel <-chan proxy.Dialer, apiKey string, carriers *map[string]*os.File, uncheckedNumFile *map[string]*os.File) {
	defer wg.Done()
	numbers := <-numberChunks

	for _, number := range numbers {
		for proxy := range proxyChannel {
			result, statusCode, err := LookupAPI(proxy, number, apiKey)
			if err == nil {
				mutex.Lock()
				carrier := result.Data.Portability.SPIDCarrierName
				if _, exists := (*carriers)[carrier]; !exists {
					file, _ := fileutil.WriteToFile("HRL_Lookup/verified_numbers", carrier+".txt")
					(*carriers)[carrier] = file
					file.WriteString(number + "\n")
				} else {
					file := (*carriers)[carrier]
					file.WriteString(number + "\n")
				}
				mutex.Unlock()
				break
			} else if statusCode == 403 {
				mutex.Lock()
				if _, exists := (*uncheckedNumFile)["unchecked"]; !exists {
					file, _ := fileutil.WriteToFile("HRL_Lookup/unchecked_numbers", "unchecked_numbers.txt")
					(*uncheckedNumFile)["unchecked"] = file
					file.WriteString(number + "\n")
				} else {
					file := (*uncheckedNumFile)["unchecked"]
					file.WriteString(number + "\n")
				}
				mutex.Unlock()
				break
			}
		}
	}
}
