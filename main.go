package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"gopkg.in/gomail.v2"
	"gopkg.in/yaml.v2"
)

type EmailConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
	To       string `yaml:"to"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	SSL      bool   `yaml:"ssl"`
	TLS      bool   `yaml:"tls"`
	Auth     string `yaml:"auth"`
	Subject  string `yaml:"subject"`
	Body     string `yaml:"body"`
}

var version string
var commitId string

func checkVersion() {
	fmt.Println("gomtp version ", version, " ", commitId)
}

func main() {

	if len(os.Args) > 1 {
		if os.Args[1] == "--version" {
			checkVersion()
			os.Exit(0)
		}
	}

	// Read the YAML configuration file
	configFile, err := os.ReadFile("gomtp.yml")
	if err != nil {
		log.Fatal(err)
	}

	var emailConfig EmailConfig
	err = yaml.Unmarshal(configFile, &emailConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Set default values for subject and body if they are empty
	if emailConfig.Subject == "" {
		emailConfig.Subject = "GOMTP Test Subject"
	}

	if emailConfig.Body == "" {
		emailConfig.Body = "This is the test email sent by gomtp."
	}

	// Create the email message
	m := gomail.NewMessage()
	m.SetHeader("From", emailConfig.From)
	m.SetHeader("To", emailConfig.To)
	m.SetHeader("Subject", emailConfig.Subject)
	m.SetBody("text/plain", emailConfig.Body)

	// Dial to the SMTP server
	d := gomail.NewDialer(emailConfig.Host, emailConfig.Port, emailConfig.Username, emailConfig.Password)

	// Enable SSL and TLS
	d.SSL = emailConfig.SSL

	tlsConfig := &tls.Config{
		ServerName:         emailConfig.Host,
		InsecureSkipVerify: !emailConfig.TLS,
	}
	d.TLSConfig = tlsConfig

	// Authenticate with the server using the specified method
	switch emailConfig.Auth {
	case "LOGIN":
		d.Auth = smtp.PlainAuth("", emailConfig.Username, emailConfig.Password, emailConfig.Host)
	}

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Email sent successfully!")
}
