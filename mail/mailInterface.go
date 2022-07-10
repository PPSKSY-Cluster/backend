package mail

import (
	"net/smtp"
	"os"
)

func SendMail(to string, message string) error {
	authentication := smtp.PlainAuth("", os.Getenv("RELAY_USERNAME"),
		os.Getenv("RELAY_PASSWORD"),
		os.Getenv("RELAY_HOST"))

	return smtp.SendMail(os.Getenv("RELAY_HOST")+":"+os.Getenv("RELAY_PORT"),
		authentication,
		os.Getenv("RELAY_USERNAME"),
		[]string{to},
		[]byte(message))
}
