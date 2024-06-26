package gmailsmtp

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"gopkg.in/gomail.v2"
)

func SendMail(numbersChan <-chan []string, wg *sync.WaitGroup, mutex *sync.Mutex, smtpChan <-chan string, domain string, totalSent *int, senderName string, messageBody string, limitExceeded *map[string]bool, invalidSMTPs *map[string]bool, invalidSMTPFormat *map[string]bool, smtpConn *map[string]gomail.SendCloser, totalSMTPs int, subject string, files []*os.File) {
	defer wg.Done()
	numbers := <-numbersChan
	for main_index, number := range numbers {
		newNumber := strings.TrimPrefix(number, "+1")
		target := fmt.Sprintf("%s%s", newNumber, domain)
		for smtp := range smtpChan {
			if !(*invalidSMTPFormat)[smtp] && !(*limitExceeded)[smtp] && !(*invalidSMTPs)[smtp] {
				splittedSmtpCreds, err := FilterCreds(smtp)
				if err == nil {
					username, password := splittedSmtpCreds[0], splittedSmtpCreds[1]
					var conn gomail.SendCloser
					if _, exists := (*smtpConn)[smtp]; !exists {
						dialer := gomail.NewDialer("smtp.gmail.com", 587, username, password)
						conn1, err := dialer.Dial()
						if err == nil {
							conn = conn1
						} else {
							mutex.Lock()
							if _, exists := (*invalidSMTPs)[smtp]; !exists {
								(*invalidSMTPs)[smtp] = true
								files[3].WriteString(smtp + "\n")
							}
							mutex.Unlock()
							continue
						}
					}
					msg := gomail.NewMessage()
					msg.SetAddressHeader("From", username, senderName)
					msg.SetHeader("To", target)
					msg.SetHeader("Subject", subject)
					msg.SetBody("text/plain", messageBody)
					mutex.Lock()
					if conn != nil {
						(*smtpConn)[smtp] = conn
					}
					err := gomail.Send((*smtpConn)[smtp], msg)
					mutex.Unlock()
					if err == nil {
						mutex.Lock()
						*totalSent++
						color.New(color.FgBlue).Printf("%d -> ", *totalSent)
						color.New(color.FgHiGreen).Printf("%s | SMTP -> %s | Status -> Sent\n", number, username)
						files[0].WriteString(number + "\n")
						mutex.Unlock()
						break
					} else if strings.Contains(err.Error(), "SMTP Daily user sending quota exceeded.") {
						mutex.Lock()
						if _, exists := (*limitExceeded)[smtp]; !exists {
							(*limitExceeded)[smtp] = true
							files[2].WriteString(smtp + "\n")
						}
						mutex.Unlock()
						continue
					} else {
						mutex.Lock()
						duration := 3 * time.Minute
						startTime := time.Now()
						for {
							currentTime := time.Now()
							remainingTime := duration - currentTime.Sub(startTime)
							if remainingTime <= 0 {
								break
							}
							// color.New(color.FgHiMagenta).Printf("\rTime Remaining: %02d:%02d",
							// 	int(remainingTime.Minutes()),
							// 	int(remainingTime.Seconds())%60)
							time.Sleep(time.Second)
						}
						newConn, err := CreateSMTPConn(username, password)
						if err == nil {
							(*smtpConn)[smtp] = newConn
						}
						mutex.Unlock()
						continue
					}
				} else {
					mutex.Lock()
					if _, exists := (*invalidSMTPFormat)[smtp]; !exists {
						(*invalidSMTPFormat)[smtp] = true
						files[4].WriteString(smtp + "\n")
					}
					mutex.Unlock()
					continue
				}
			} else if len((*invalidSMTPFormat)) == totalSMTPs || len((*limitExceeded)) == totalSMTPs || len((*invalidSMTPs)) == totalSMTPs {
				mutex.Lock()
				for i := main_index; i < len(numbers); i++ {
					num := numbers[i]
					files[1].WriteString(num + "\n")
				}
				mutex.Unlock()
				return
			}

		}
	}
}
