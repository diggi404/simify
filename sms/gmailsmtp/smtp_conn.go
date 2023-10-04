package gmailsmtp

import (
	"gopkg.in/gomail.v2"
)

func CreateSMTPConn(username, password string) (gomail.SendCloser, error) {
	dialer := gomail.NewDialer("smtp.gmail.com", 587, username, password)
	conn, err := dialer.Dial()
	if err != nil {
		return nil, err
	}
	return conn, nil
}
