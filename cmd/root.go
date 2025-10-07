package cmd

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/gomail.v2"
	"gopkg.in/yaml.v2"
)

// CLI flags
var gomtpYamlPath string
var emailTo string
var emailSubject string
var emailBody string
var emailBodyFile string
var debug bool
var ccList []string

var version string
var commitId string

type EmailConfig struct {
	Username          string   `yaml:"username"`
	Password          string   `yaml:"password"`
	From              string   `yaml:"from"`
	To                string   `yaml:"to"`
	Host              string   `yaml:"host"`
	Port              int      `yaml:"port"`
	SSL               bool     `yaml:"ssl"`
	TLS               bool     `yaml:"tls"`
	Auth              string   `yaml:"auth"`
	VerifyCertificate bool     `default:"true" yaml:"verifyCertificate"`
	Subject           string   `yaml:"subject"`
	Body              string   `yaml:"body"`
	CcList            []string `yaml:"cc"`
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

	// Read the email body from stdin if provided
	stdioBody, err := readBodyFromStdin()
	if err != nil {
		return err
	}

	// Set body input
	err = setBody(&emailConfig, stdioBody)
	if err != nil {
		return err
	}

	setupDefaultEmailConfig(&emailConfig)

	setFlags(&emailConfig)

	// Create the email message
	emailMessage := createEmailMessage(&emailConfig)

	err = sendEmail(&emailConfig, emailMessage)
	if err != nil {
		return err
	}

	cmd.Printf("Email sent successfully!")
	return nil
}

// Setup default values for flags, get the email config pointer.
func setupDefaultEmailConfig(emailConfig *EmailConfig) {
	// Set default values for subject and body if they are empty
	if emailConfig.Subject == "" {
		emailConfig.Subject = "GOMTP Test Subject"
	}
	if emailConfig.To == "" {
		emailConfig.To = "to@example.com"
	}
	if emailConfig.Body == "" {
		emailConfig.Body = "This is the test email sent by gomtp."
	}
	if len(ccList) > 0 {
		emailConfig.CcList = ccList
	}
}

func readBodyFromStdin() (string, error) {
	stdinInfo, _ := os.Stdin.Stat()
	if (stdinInfo.Mode() & os.ModeCharDevice) == 0 {
		stdioBody, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return string(stdioBody), nil
	}
	return "", nil
}

// check --body, --body-file and stdio for email body
func setBody(emailConfig *EmailConfig, stdioBody string) error {
	bodySourceCount := 0

	if stdioBody != "" {
		bodySourceCount++
		emailConfig.Body = stdioBody
	}
	if emailBody != "" {
		bodySourceCount++
		emailConfig.Body = emailBody
	}
	if emailBodyFile != "" {
		bodySourceCount++
		body, err := os.ReadFile(emailBodyFile)
		if err != nil {
			return err
		}
		emailConfig.Body = string(body)
	}
	if bodySourceCount > 1 {
		return fmt.Errorf("cannot specify body via multiple sources simultaneously")
	}
	return nil
}

// Set values from global flags
func setFlags(emailConfig *EmailConfig) {
	if emailTo != "" {
		emailConfig.To = emailTo
	}
	if emailSubject != "" {
		emailConfig.Subject = emailSubject
	}
}

// Create email message from config
func createEmailMessage(emailConfig *EmailConfig) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", emailConfig.From)
	m.SetHeader("To", emailConfig.To)
	m.SetHeader("Subject", emailConfig.Subject)
	m.SetBody("text/plain", emailConfig.Body)
	if len(emailConfig.CcList) > 0 {
		m.SetHeader("Cc", emailConfig.CcList...)
	}
	return m
}

func sendEmail(emailConfig *EmailConfig, m *gomail.Message) error {
	// Validate mode selection
	if emailConfig.SSL && emailConfig.TLS {
		return fmt.Errorf("invalid configuration: both SSL and TLS (STARTTLS) are enabled; choose only one")
	}

	// Render message to bytes once
	var msgBuf bytes.Buffer
	if _, err := m.WriteTo(&msgBuf); err != nil {
		return err
	}

	// Common
	addr := net.JoinHostPort(emailConfig.Host, strconv.Itoa(emailConfig.Port))
	tlsConfig := &tls.Config{
		ServerName:         emailConfig.Host,
		InsecureSkipVerify: !emailConfig.VerifyCertificate,
	}

	if debug {
		fmt.Fprintf(os.Stderr, "[gomtp][debug] host=%s port=%d ssl=%t starttls=%t auth=%s verifyCert=%t\n", emailConfig.Host, emailConfig.Port, emailConfig.SSL, emailConfig.TLS, emailConfig.Auth, emailConfig.VerifyCertificate)
	}

	// Connect
	var (
		conn net.Conn
		c    *smtp.Client
		err  error
	)

	if emailConfig.SSL {
		if debug {
			fmt.Fprintf(os.Stderr, "[gomtp][debug] selected_mode=ssl_implicit\n")
		}
		conn, err = tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return err
		}
		defer conn.Close()
		// Log TLS parameters for implicit TLS
		if debug {
			if st := conn.(*tls.Conn).ConnectionState(); true {
				fmt.Fprintf(os.Stderr, "[gomtp][debug] negotiated_tls=version:%x cipher_suite:%x server_name=%s\n", st.Version, st.CipherSuite, st.ServerName)
				if len(st.PeerCertificates) > 0 {
					cert := st.PeerCertificates[0]
					fmt.Fprintf(os.Stderr, "[gomtp][debug] cert_subject=%s issuer=%s not_before=%s not_after=%s dns_names=%v\n",
						cert.Subject.String(), cert.Issuer.String(), cert.NotBefore.Format(time.RFC3339), cert.NotAfter.Format(time.RFC3339), cert.DNSNames)
				}
			}
		}
		c, err = smtp.NewClient(conn, emailConfig.Host)
		if err != nil {
			return err
		}
	} else {
		if emailConfig.TLS {
			if debug {
				fmt.Fprintf(os.Stderr, "[gomtp][debug] selected_mode=starttls\n")
			}
		} else {
			if debug {
				fmt.Fprintf(os.Stderr, "[gomtp][debug] selected_mode=plain_no_tls\n")
			}
		}
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			return err
		}
		defer conn.Close()
		c, err = smtp.NewClient(conn, emailConfig.Host)
		if err != nil {
			return err
		}
	}
	defer c.Quit()

	// EHLO/HELO
	if err := c.Hello(emailConfig.Host); err != nil {
		return err
	}

	// STARTTLS if requested
	if emailConfig.TLS {
		if ok, _ := c.Extension("STARTTLS"); !ok {
			return fmt.Errorf("server does not support STARTTLS")
		}
		if err := c.StartTLS(tlsConfig); err != nil {
			return err
		}
		if debug {
			if st, ok := c.TLSConnectionState(); ok {
				fmt.Fprintf(os.Stderr, "[gomtp][debug] negotiated_tls=version:%x cipher_suite:%x server_name=%s\n", st.Version, st.CipherSuite, st.ServerName)
				if len(st.PeerCertificates) > 0 {
					cert := st.PeerCertificates[0]
					fmt.Fprintf(os.Stderr, "[gomtp][debug] cert_subject=%s issuer=%s not_before=%s not_after=%s dns_names=%v\n",
						cert.Subject.String(), cert.Issuer.String(), cert.NotBefore.Format(time.RFC3339), cert.NotAfter.Format(time.RFC3339), cert.DNSNames)
				}
			}
		}
	}

	// AUTH if configured and supported
	if emailConfig.Auth == "LOGIN" {
		if ok, _ := c.Extension("AUTH"); ok {
			a := smtp.PlainAuth("", emailConfig.Username, emailConfig.Password, emailConfig.Host)
			if err := c.Auth(a); err != nil {
				return err
			}
		} else if debug {
			fmt.Fprintf(os.Stderr, "[gomtp][debug] server does not advertise AUTH; skipping auth\n")
		}
	}

	// MAIL FROM
	if err := c.Mail(emailConfig.From); err != nil {
		return err
	}

	// RCPT TO (To + Cc)
	recipients := make([]string, 0, 1+len(emailConfig.CcList))
	if emailConfig.To != "" {
		recipients = append(recipients, emailConfig.To)
	}
	recipients = append(recipients, emailConfig.CcList...)
	for _, rcpt := range recipients {
		if rcpt == "" {
			continue
		}
		if err := c.Rcpt(rcpt); err != nil {
			return err
		}
	}

	// DATA
	wc, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := wc.Write(msgBuf.Bytes()); err != nil {
		_ = wc.Close()
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

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
	rootCmd.Flags().StringVar(&emailBodyFile, "body-file", "", "File that contains body of the email.")
	rootCmd.Flags().StringSliceVar(&ccList, "cc", []string{}, "CC email address")
	rootCmd.Flags().BoolVar(&debug, "debug", false, "Enable verbose SMTP/TLS debugging output.")

}
