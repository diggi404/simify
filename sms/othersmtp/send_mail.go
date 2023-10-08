package othersmtp

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"gopkg.in/gomail.v2"
)

func SendMail(numbersChan <-chan []string, wg *sync.WaitGroup, mutex *sync.Mutex, smtpChan <-chan string, domain string, totalSent *int, senderName string, messageBody string, limitExceeded *map[string]bool, invalidSMTPs *map[string]bool, invalidSMTPformat *map[string]bool, smtpConn *map[string]gomail.SendCloser, totalSMTPs int, subject string, files []*os.File) {
	defer wg.Done()
	numbers := <-numbersChan
	for mainIndex, number := range numbers {
		newNumber := strings.TrimPrefix(number, "+1")
		target := fmt.Sprintf("%s%s", newNumber, domain)
		for smtp := range smtpChan {
			if !(*invalidSMTPformat)[smtp] && !(*limitExceeded)[smtp] && !(*invalidSMTPs)[smtp] {
				splittedCreds, port, err := FilterCreds(smtp)
				if err == nil {
					host, username, password := splittedCreds[0], splittedCreds[2], splittedCreds[3]
					var conn gomail.SendCloser
					if _, exists := (*smtpConn)[smtp]; !exists {
						dialer := gomail.NewDialer(host, port, username, password)
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
						color.New(color.FgHiGreen).Printf("%s | SMTP -> %s:%d | Status -> Sent\n", number, username, port)
						files[0].WriteString(number + "\n")
						mutex.Unlock()
						break
					} else {
						mutex.Lock()
						if _, exists := (*limitExceeded)[smtp]; !exists {
							(*limitExceeded)[smtp] = true
							files[2].WriteString(smtp + "\n")
						}
						mutex.Unlock()
						continue
					}
				} else {
					mutex.Lock()
					if _, exists := (*invalidSMTPformat)[smtp]; !exists {
						(*invalidSMTPformat)[smtp] = true
						files[4].WriteString(smtp + "\n")
					}
					mutex.Unlock()
					continue
				}
			} else if len(*invalidSMTPformat) == totalSMTPs || len(*limitExceeded) == totalSMTPs || len(*invalidSMTPs) == totalSMTPs {
				mutex.Lock()
				for i := mainIndex; i < len(numbers); i++ {
					num := numbers[i]
					files[1].WriteString(num + "\n")
				}
				mutex.Unlock()
				return
			}
		}
	}
}
