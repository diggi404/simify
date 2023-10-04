package gmailsmtp

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
	"gopkg.in/gomail.v2"
)

func GmailSMTPToSMS() {
	color.New(color.FgHiRed).Print("\n(!) SMTPs should be of the format -> example@gmail.com:password\n(!) Invalid formats will be skipped automatically.\n")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nPress Enter to select your SMTPs: ")
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
	smtpList, err := fileutil.ReadFromFile(filePath)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	if len(smtpList) == 0 {
		color.New(color.FgHiRed).Println("\nEmpty file. Exiting...")
		return
	}

	color.New(color.FgHiMagenta).Printf("\nTotal SMTPs: %d\n", len(smtpList))

	color.New(color.FgHiRed).Print("\n(!) Choose the type of carrier numbers you are sending.\n(!) Wrong selection will lead to undelivered messages.\n")

	carrierList := `
	[1] Verizon Wireless			[2] AT&T			[3] T-Mobile

	[4] Sprint PCS				[5] US Cellular			[6] Metro PCS

	[7] Cricket Wireless			[8] Google Fi
	`
	color.New(color.FgHiGreen).Print(carrierList)

	var choiceStr string
	fmt.Print("\nEnter your Carrier: ")
	fmt.Scanln(&choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil {
		color.New(color.FgHiRed).Println("\nInvalid Choice. Exiting Program...")
		return
	}
	smsDomans := []struct {
		choice int
		domain string
		name   string
	}{
		{1, "@vtext.com", "Verizon Wireless"},
		{2, "@txt.att.net", "AT&T"},
		{3, "@tmomail.net", "T-Mobile"},
		{4, "@messaging.sprintpcs.com", "Sprint PCS"},
		{5, "@email.uscc.net", "US Cellular"},
		{6, "@mymetropcs.com", "Metro PCS"},
		{7, "@sms.cricketwireless.net", "Cricket Wireless"},
		{8, "@msg.fi.google.com", "Google Fi"},
	}
	if choice > len(smsDomans) || choice <= 0 {
		color.New(color.FgHiRed).Println("\nInvalid choice. Exiting Program...")
		return
	}
	var domain string
	var name string
	for _, value := range smsDomans {
		if value.choice == choice {
			domain = value.domain
			name = value.name
		}
	}
	if len(domain) == 0 {
		color.New(color.FgHiRed).Println("\nInvalid choice. Exiting Program...")
	}
	fmt.Print("\nYou selected: ")
	color.New(color.FgHiMagenta).Print(name + "\n")

	color.New(color.FgHiRed).Print("\n(!) Numbers should be of the format -> +13059837812\n(!) Invalid formats will be skipped automatically.\n")

	fmt.Print("\nPress Enter to select your numbers: ")
	_, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	numbersFilePath, err := zenity.SelectFile(
		zenity.FileFilters{
			{Patterns: []string{"*.txt"}, CaseFold: false},
		})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	numberList, err := fileutil.ReadFromFile(numbersFilePath)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	if len(numberList) == 0 {
		color.New(color.FgHiRed).Println("\nEmpty file. Exiting...")
		return
	}
	color.New(color.FgHiMagenta).Printf("\nTotal Numbers: %d\n", len(numberList))

	fmt.Print("\nEnter your Sender Name: ")
	senderName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Print("\nEnter your Message: ")
	messageBody, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	senderName = strings.TrimRight(senderName, "\r\n")
	messageBody = strings.TrimRight(messageBody, "\r\n")

	maxWorkers := 100
	chunkSize := len(numberList) / maxWorkers

	if len(numberList)%maxWorkers != 0 {
		chunkSize++
	}
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var totalSent int
	numbersChan := make(chan []string, chunkSize)
	smtpChan := make(chan string, 1)
	limitExceeded := make(map[string]bool)
	smtpConn := make(map[string]gomail.SendCloser)
	totalSMTPs := len(smtpList)

	fmt.Println()
	fmt.Println()
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go SendMail(numbersChan, &wg, &mutex, smtpChan, domain, &totalSent, senderName, messageBody, &limitExceeded, &smtpConn, totalSMTPs)
	}

	for i := 0; i < len(numberList); i += chunkSize {
		end := i + chunkSize
		if end > len(numberList) {
			end = len(numberList)
		}
		numbersChan <- numberList[i:end]
	}
	close(numbersChan)

	go func() {
		for {
			for _, smtp := range smtpList {
				smtpChan <- smtp
			}
		}
	}()

	wg.Wait()
	fmt.Println("all done.")

}
