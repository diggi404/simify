package single_api

import (
	"fmt"
	"os"
	"simify/fileutil"
	"sync"

	"github.com/fatih/color"
	"golang.org/x/net/proxy"
)

func ProcessLookup(numberChunks <-chan []string, wg *sync.WaitGroup, mutex *sync.Mutex, proxyChannel <-chan proxy.Dialer, apiKey string, carriers *map[string]*os.File, uncheckedNumFile *map[string]*os.File, totalChecks *int, uncheckOutPutChan chan bool, chanClosed *bool) {
	defer wg.Done()
	numbers := <-numberChunks

	for main_index, number := range numbers {
		for proxy := range proxyChannel {
			result, statusCode, err := LookupAPI(proxy, number, apiKey)
			if err == nil {
				mutex.Lock()
				lineType := result.Data.Portability.LineType
				carrier := result.Data.Portability.SPIDCarrierName
				state := result.Data.Portability.State
				city := result.Data.Portability.City
				location := fmt.Sprintf("%s, %s", city, state)
				*totalChecks++
				color.New(color.FgBlue).Printf("%d -> ", *totalChecks)
				outputStr := fmt.Sprintf("%s | %s | %s | %s", number, lineType, carrier, location)
				color.New(color.FgHiGreen).Printf("%s\n", outputStr)
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
				if !*chanClosed {
					close(uncheckOutPutChan)
					*chanClosed = true
				}
				if _, exists := (*uncheckedNumFile)["unchecked"]; !exists {
					file, _ := fileutil.WriteToFile("HRL_Lookup/unchecked_numbers", "unchecked_numbers.txt")
					(*uncheckedNumFile)["unchecked"] = file
					for i := main_index; i < len(numbers); i++ {
						num := numbers[i]
						file.WriteString(num + "\n")
					}
				} else {
					file := (*uncheckedNumFile)["unchecked"]
					for i := main_index; i < len(numbers); i++ {
						num := numbers[i]
						file.WriteString(num + "\n")
					}
				}
				mutex.Unlock()
				return
			}
		}
	}
}
