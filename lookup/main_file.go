package lookup

import (
	"bufio"
	"fmt"
	"os"
	"simify/fileutil"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/ncruces/zenity"
)

func HRLLOOKUP() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nPress Enter to select your numbers: ")
	_, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	filePath, err := zenity.SelectFile(
		zenity.FileFilters{
			{Patterns: []string{"*.txt"}, CaseFold: false},
		})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	numberList, err := fileutil.ReadFromFile(filePath)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	color.New(color.FgHiMagenta).Printf("\nTotal Numbers: %d\n", len(numberList))

	var proxyTypeStr string
	fmt.Print("\nWhat type of Proxies? [HTTP,SOCKS4,SOCKS5] :> ")
	fmt.Scanln(&proxyTypeStr)
	if len(proxyTypeStr) == 0 {
		fmt.Printf("invalid input.Exiting Program...\n")
		return
	}
	proxyTypeStr = strings.ToLower(strings.TrimSpace(proxyTypeStr))
	var proxyType string
	if strings.Contains(proxyTypeStr, "http") {
		proxyType = "http"
	} else if strings.Contains(proxyTypeStr, "socks4") {
		proxyType = "socks4"
	} else if strings.Contains(proxyTypeStr, "socks5") {
		proxyType = "socks5"
	} else {
		fmt.Printf("invalid choice. Exiting Program...\n")
		return
	}

	fmt.Print("\nPress Enter to load your Proxies: ")
	_, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	proxyFilePath, err := zenity.SelectFile(
		zenity.FileFilters{
			{Patterns: []string{"*.txt"}, CaseFold: false},
		})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	proxyList, err := fileutil.ReadFromFile(proxyFilePath)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	color.New(color.FgHiMagenta).Printf("\nTotal Proxies: %d\n", len(proxyList))

	maxWorkers := 100
	chunkSize := len(numberList) / maxWorkers

	if len(numberList)%maxWorkers != 0 {
		chunkSize++
	}
	numberChunks := make(chan []string, chunkSize)
	var wg sync.WaitGroup
	var mutex sync.Mutex

	proxies := ProxyParser(proxyList, proxyType)

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go ProcessLookup(numberChunks, &wg, &mutex, &proxies)
	}

	for i := 0; i < len(numberList); i += chunkSize {
		end := i + chunkSize
		if end > len(numberList) {
			end = len(numberList)
		}
		numberChunks <- numberList[i:end]
	}
	close(numberChunks)
	wg.Wait()
	fmt.Println("\nall done.")
}