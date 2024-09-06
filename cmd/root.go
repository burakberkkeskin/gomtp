package cmd

import (
	"crypto/tls"
	"net/smtp"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/gomail.v2"
	"gopkg.in/yaml.v2"
)

// CLI flags
var gomtpYamlPath string
var emailTo string
var emailSubject string
var emailBody string

var version string
var commitId string

type EmailConfig struct {
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	From              string `yaml:"from"`
	To                string `yaml:"to"`
	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	SSL               bool   `yaml:"ssl"`
	TLS               bool   `yaml:"tls"`
	Auth              string `yaml:"auth"`
	VerifyCertificate bool   `default:"true" yaml:"verifyCertificate"`
	Subject           string `yaml:"subject"`
	Body              string `yaml:"body"`
}

const usageMessage = `Example Commands: 
  gomtp # Read the gomtp.yaml file and send a test email.
  gomtp -f custom.yaml # Read the custom.yaml file and send a test email.`

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gomtp",
	Short:   "Gomtp is a CLI tool for go that tests SMTP settings easily.",
	Long:    usageMessage,
	Version: version + " " + commitId,
	RunE:    rootRun,
}

func rootRun(cmd *cobra.Command, args []string) error {
	// Read the YAML configuration file
	configFile, err := os.ReadFile(gomtpYamlPath)
	if err != nil {
		return err
	}
	var emailConfig EmailConfig
	err = yaml.Unmarshal(configFile, &emailConfig)
	if err != nil {
		return err
	}
	// Set default values for subject and body if they are empty
	if emailConfig.Subject == "" {
		emailConfig.Subject = "GOMTP Test Subject"
	}
	if emailConfig.Body == "" {
		emailConfig.Body = "This is the test email sent by gomtp."
	}
	if emailConfig.To == "" {
		emailConfig.To = "to@example.com"
	}
	if emailTo != "" {
		emailConfig.To = emailTo
	}
	if emailSubject != "" {
		emailConfig.Subject = emailSubject
	}
	if emailBody != "" {
		emailConfig.Body = emailBody
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
		InsecureSkipVerify: emailConfig.VerifyCertificate,
	}
	d.TLSConfig = tlsConfig

	// Authenticate with the server using the specified method
	switch emailConfig.Auth {
	case "LOGIN":
		d.Auth = smtp.PlainAuth("", emailConfig.Username, emailConfig.Password, emailConfig.Host)
	}

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	cmd.Printf("Email sent successfully!")
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SilenceUsage = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("help", "h", false, "Help menu.")
	rootCmd.Flags().StringVarP(&gomtpYamlPath, "file", "f", "gomtp.yaml", "Configuration file path.")
	rootCmd.Flags().StringVar(&emailTo, "to", "", "Target email address.")
	rootCmd.Flags().StringVarP(&emailSubject, "subject", "s", "", "Subject of the email.")
	rootCmd.Flags().StringVarP(&emailBody, "body", "b", "", "Body of the email.")
}
