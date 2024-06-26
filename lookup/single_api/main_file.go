package single_api

import (
	"bufio"
	"fmt"
	"os"
	"simify/fileutil"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/ncruces/zenity"
	"golang.org/x/net/proxy"
)

func SingleAPILookup() {
	red := color.New(color.FgHiRed).PrintfFunc()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nEnter your Telnyx API Key :> ")
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	apiKey = strings.TrimRight(apiKey, "\r\n")
	color.New(color.FgBlue).Println("\nChecking account balance...")

	info, err := CheckBalance(apiKey)
	if err != nil {
		red("\nerr: %v\n", err)
		return
	}
	balanceInfoStr := fmt.Sprintf("\nAccount Balance: %s %s", info.Data.Currency, info.Data.Balance)
	color.New(color.FgHiMagenta).Println(balanceInfoStr)
	balance64, err := strconv.ParseFloat(info.Data.Balance, 64)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	balance := int(balance64)
	if balance <= 0 {
		color.New(color.FgHiRed).Println("\nYour account has a low balance. Kinldy Topup to continue.")
		return
	}
	fmt.Print("\nPress Enter to select your numbers: ")
	_, err = reader.ReadString('\n')
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
	if len(numberList) == 0 {
		color.New(color.FgHiRed).Println("\nEmpty File. Exiting...")
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
	if len(proxyList) == 0 {
		color.New(color.FgHiRed).Println("\nEmtpy File. Exiting...")
		return
	}
	color.New(color.FgHiMagenta).Printf("\nTotal Proxies: %d\n", len(proxyList))

	maxWorkers := 100
	chunkSize := len(numberList) / maxWorkers

	if len(numberList)%maxWorkers != 0 {
		chunkSize++
	}
	numberChunks := make(chan []string, chunkSize)
	proxyChannel := make(chan proxy.Dialer, 1)
	uncheckOutPutChan := make(chan bool)
	var wg sync.WaitGroup
	var mutex sync.Mutex
	carriers := make(map[string]*os.File)
	uncheckedNumFile := make(map[string]*os.File)
	var totalChecks int
	var chanClosed bool

	proxies := ProxyParser(proxyList, proxyType)

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go ProcessLookup(numberChunks, &wg, &mutex, proxyChannel, apiKey, &carriers, &uncheckedNumFile, &totalChecks, uncheckOutPutChan, &chanClosed)
	}

	for i := 0; i < len(numberList); i += chunkSize {
		end := i + chunkSize
		if end > len(numberList) {
			end = len(numberList)
		}
		numberChunks <- numberList[i:end]
	}
	close(numberChunks)

	go func() {
		for {
			for _, proxy := range proxies {
				proxyChannel <- proxy
			}
		}
	}()

	go func() {
		run := true
		for run {
			_, ok := <-uncheckOutPutChan
			if !ok {
				color.New(color.FgHiRed).Printf("\nYou have exhausted your account balance!\nUpdating unchecked_numbers file. Please wait...\n")
				run = false
			}
		}
	}()
	wg.Wait()
	fmt.Println("\nall done.")
}
