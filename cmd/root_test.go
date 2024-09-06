package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHappyPath(t *testing.T) {
	command := rootCmd
	command.SetArgs([]string{
		"--file", "./tests/gomtpYamls/successConfiguration.yaml",
	})
	b := bytes.NewBufferString("")
	command.SetOut(b)
	command.SetErr(b)
	command.Execute()
	var expected string = "Email sent successfully!"
	assert.Equal(t, expected, b.String(), "actual is not expected")
}

func TestGomtpYamlNotFound(t *testing.T) {
	command := rootCmd
	b := bytes.NewBufferString("")
	command.SetArgs([]string{
		"--file", "./unknown/path.yaml",
	})
	command.SetOut(b)
	command.SetErr(b)
	command.Execute()
	var expected string = "no such file or directory"
	assert.Contains(t, b.String(), expected, "File not found error expected.")
}

func TestInvalidCredentialsGoogle(t *testing.T) {
	command := rootCmd
	b := bytes.NewBufferString("")
	command.SetArgs([]string{
		"--file", "./tests/gomtpYamls/invalidCredentialsGoogle.yaml",
	})
	command.SetOut(b)
	command.SetErr(b)
	command.Execute()
	var expected string = "Username and Password not accepted"
	assert.Contains(t, b.String(), expected, "Invalid credentials error expected.")
}

func TestInvalidCredentialsYandex(t *testing.T) {
	command := rootCmd
	b := bytes.NewBufferString("")
	command.SetArgs([]string{
		"--file", "./tests/gomtpYamls/invalidCredentialsYandex.yaml",
	})
	command.SetOut(b)
	command.SetErr(b)
	command.Execute()
	var expected string = "authentication failed: Invalid user or password"
	assert.Contains(t, b.String(), expected, "Invalid credentials error expected.")

}

func TestInvalidCredentialsBrevo(t *testing.T) {
	command := rootCmd
	b := bytes.NewBufferString("")
	command.SetArgs([]string{
		"--file", "./tests/gomtpYamls/invalidCredentialsBrevo.yaml",
	})
	command.SetOut(b)
	command.SetErr(b)
	command.Execute()
	var expected string = "Authentication failed"
	assert.Contains(t, b.String(), expected, "Invalid credentials error expected.")

}

func TestInvalidSSLConfiguration(t *testing.T) {
	command := rootCmd
	b := bytes.NewBufferString("")
	command.SetArgs([]string{
		"--file", "./tests/gomtpYamls/nonSslServerWithSslConfiguration.yaml",
	})
	command.SetOut(b)
	command.SetErr(b)
	command.Execute()
	var expected string = "tls: first record does not look like a TLS handshake"
	assert.Contains(t, b.String(), expected, "SSL error expected.")
}

func TestNonTLSServerWithTLSConfiguration(t *testing.T) {

}

func TestOptionalParameters(t *testing.T) {
	command := rootCmd
	command.SetArgs([]string{
		"--file", "./tests/gomtpYamls/optionalParameters.yaml",
	})
	b := bytes.NewBufferString("")
	command.SetOut(b)
	command.SetErr(b)
	command.Execute()
	var expected string = "Email sent successfully!"
	assert.Equal(t, expected, b.String(), "actual is not expected")
}

func TestInvalidYaml(t *testing.T) {
	command := rootCmd
	command.SetArgs([]string{
		"--file", "./tests/gomtpYamls/invalid.yaml",
	})
	b := bytes.NewBufferString("")
	command.SetOut(b)
	command.SetErr(b)
	err := command.Execute()
	assert.NotEmpty(t, err)
}

type MailhogResponse struct {
	Total int              `json:"total"`
	Items []MailhogMessage `json:"items"`
}

type MailhogMessage struct {
	From    MailhogAddress   `json:"From"`
	To      []MailhogAddress `json:"To"`
	Content struct {
		Headers map[string][]string `json:"Headers"`
		Body    string              `json:"Body"`
	} `json:"Content"`
}

type MailhogAddress struct {
	Mailbox string `json:"Mailbox"`
	Domain  string `json:"Domain"`
}

func TestEmptySubjectYaml(t *testing.T) {
	command := rootCmd
	command.SetArgs([]string{
		"--file", "./tests/gomtpYamls/emptySubject.yaml",
	})
	b := bytes.NewBufferString("")
	command.SetOut(b)
	command.SetErr(b)
	err := command.Execute()
	assert.NoError(t, err, "command execution failed")

	var expected string = "Email sent successfully!"
	assert.Equal(t, expected, b.String(), "actual is not expected")

	// Wait a moment for MailHog to process the email
	time.Sleep(1 * time.Second)

	// Check MailHog for the sent email
	resp, err := http.Get("http://localhost:8025/api/v2/messages")
	assert.NoError(t, err, "failed to get messages from MailHog")
	defer resp.Body.Close()

	var mailhogResp MailhogResponse
	err = json.NewDecoder(resp.Body).Decode(&mailhogResp)
	assert.NoError(t, err, "failed to decode MailHog response")

	assert.NotEmpty(t, mailhogResp.Items, "no messages found in MailHog")

	// Check the latest message (first in the list)
	latestMessage := mailhogResp.Items[0]
	assert.Equal(t, "GOMTP Test Subject", latestMessage.Content.Headers["Subject"][0], "unexpected email subject")
	assert.Equal(t, "from", latestMessage.From.Mailbox, "unexpected sender mailbox")
	assert.Equal(t, "example.com", latestMessage.From.Domain, "unexpected sender domain")
	assert.Equal(t, "to", latestMessage.To[0].Mailbox, "unexpected recipient mailbox")
	assert.Equal(t, "example.com", latestMessage.To[0].Domain, "unexpected recipient domain")
	assert.Equal(t, "this is line 1\r\nThis is line 2", latestMessage.Content.Body, "unexpected email body")
}

func TestEmptyToYaml(t *testing.T) {
	command := rootCmd
	command.SetArgs([]string{
		"--file", "./tests/gomtpYamls/emptyTo.yaml",
	})
	b := bytes.NewBufferString("")
	command.SetOut(b)
	command.SetErr(b)
	err := command.Execute()
	assert.NoError(t, err, "command execution failed")

	assert.Equal(t, "Email sent successfully!", b.String(), "unexpected command output")

	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://localhost:8025/api/v2/messages")
	assert.NoError(t, err, "failed to get messages from MailHog")
	defer resp.Body.Close()

	var mailhogResp MailhogResponse
	err = json.NewDecoder(resp.Body).Decode(&mailhogResp)
	assert.NoError(t, err, "failed to decode MailHog response")

	assert.NotEmpty(t, mailhogResp.Items, "no messages found in MailHog")

	latestMessage := mailhogResp.Items[0]
	assert.Equal(t, "GOMTP Test Subject", latestMessage.Content.Headers["Subject"][0], "unexpected email subject")
	assert.Equal(t, "from", latestMessage.From.Mailbox, "unexpected sender mailbox")
	assert.Equal(t, "example.com", latestMessage.From.Domain, "unexpected sender domain")
	assert.Empty(t, latestMessage.To, "To field should be empty")
	assert.Equal(t, "This is the test email sent by gomtp.", latestMessage.Content.Body, "unexpected email body")
}

func TestEmptyBodyYaml(t *testing.T) {
	command := rootCmd
	command.SetArgs([]string{
		"--file", "./tests/gomtpYamls/emptyBody.yaml",
	})
	b := bytes.NewBufferString("")
	command.SetOut(b)
	command.SetErr(b)
	err := command.Execute()
	assert.NoError(t, err, "command execution failed")

	assert.Equal(t, "Email sent successfully!", b.String(), "unexpected command output")

	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://localhost:8025/api/v2/messages")
	assert.NoError(t, err, "failed to get messages from MailHog")
	defer resp.Body.Close()

	var mailhogResp MailhogResponse
	err = json.NewDecoder(resp.Body).Decode(&mailhogResp)
	assert.NoError(t, err, "failed to decode MailHog response")

	assert.NotEmpty(t, mailhogResp.Items, "no messages found in MailHog")

	latestMessage := mailhogResp.Items[0]
	assert.Equal(t, "GOMTP Test Subject", latestMessage.Content.Headers["Subject"][0], "unexpected email subject")
	assert.Equal(t, "from", latestMessage.From.Mailbox, "unexpected sender mailbox")
	assert.Equal(t, "example.com", latestMessage.From.Domain, "unexpected sender domain")
	assert.Equal(t, "to", latestMessage.To[0].Mailbox, "unexpected recipient mailbox")
	assert.Equal(t, "example.com", latestMessage.To[0].Domain, "unexpected recipient domain")
	assert.Empty(t, latestMessage.Content.Body, "email body should be empty")
}
