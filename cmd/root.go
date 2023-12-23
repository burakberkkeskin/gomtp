package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/gomail.v2"
	"gopkg.in/yaml.v2"
)

// CLI flags
var gomtpYamlPath string

var version string
var commitId string

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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gomtp",
	Short:   "Test SMTP Settings",
	Long:    `gomtp is a CLI tool for go that tests SMTP settings easily.`,
	Version: version + " " + commitId,
	Run:     rootRun,
}

func rootRun(cmd *cobra.Command, args []string) {
	// Read the YAML configuration file
	configFile, err := os.ReadFile(gomtpYamlPath)
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gomtp.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("help", "h", false, "Help menu.")
	rootCmd.Flags().StringVarP(&gomtpYamlPath, "file", "f", "gomtp.yaml", "Configuration file path.")
}
