package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MailhogResponse struct {
	Total int              `json:"total"`
	Items []MailhogMessage `json:"messages"`
}

type MailhogMessage struct {
	From    MailhogAddress   `json:"From"`
	To      []MailhogAddress `json:"To"`
	Subject string           `json:"Subject"`
	Snippet string           `json:"Snippet"`
	// Content struct {
	// 	Headers map[string][]string `json:"Headers"`
	// 	Body    string              `json:"Body"`
	// } `json:"Content"`
}

type MailhogAddress struct {
	Name    string `json:"Name"`
	Address string `json:"Address"`
}

func getLatestMessageForRecipient(t *testing.T, recipient string) MailhogMessage {
	resp, err := http.Get("http://localhost:8025/api/v1/messages")
	assert.NoError(t, err, "failed to get messages from MailHog")
	defer resp.Body.Close()

	var mailhogResp MailhogResponse
	err = json.NewDecoder(resp.Body).Decode(&mailhogResp)
	assert.NoError(t, err, "failed to decode MailHog response")

	assert.NotEmpty(t, mailhogResp.Items, "no messages found in MailHog")

	for _, msg := range mailhogResp.Items {
		if msg.To[0].Address == recipient {
			return msg
		}
	}

	t.Fatalf("No message found for recipient: %s", recipient)
	return MailhogMessage{}
}

func clearMailHog(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://localhost:8025/api/v1/messages", nil)
	assert.NoError(t, err, "failed to create DELETE request")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err, "failed to send DELETE request to MailHog")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "unexpected status code when clearing MailHog")
}

// func TestHappyPath(t *testing.T) {
// 	command := rootCmd
// 	command.SetArgs([]string{
// 		"--file", "./tests/gomtpYamls/successConfiguration.yaml",
// 	})
// 	b := bytes.NewBufferString("")
// 	command.SetOut(b)
// 	command.SetErr(b)
// 	command.Execute()
// 	var expected string = "Email sent successfully!"
// 	assert.Equal(t, expected, b.String(), "actual is not expected")
// }

// func TestGomtpYamlNotFound(t *testing.T) {
// 	command := rootCmd
// 	b := bytes.NewBufferString("")
// 	command.SetArgs([]string{
// 		"--file", "./unknown/path.yaml",
// 	})
// 	command.SetOut(b)
// 	command.SetErr(b)
// 	command.Execute()
// 	var expected string = "no such file or directory"
// 	assert.Contains(t, b.String(), expected, "File not found error expected.")
// }

// func TestInvalidCredentialsGoogle(t *testing.T) {
// 	command := rootCmd
// 	b := bytes.NewBufferString("")
// 	command.SetArgs([]string{
// 		"--file", "./tests/gomtpYamls/invalidCredentialsGoogle.yaml",
// 	})
// 	command.SetOut(b)
// 	command.SetErr(b)
// 	command.Execute()
// 	var expected string = "Username and Password not accepted"
// 	assert.Contains(t, b.String(), expected, "Invalid credentials error expected.")
// }

// func TestInvalidCredentialsYandex(t *testing.T) {
// 	command := rootCmd
// 	b := bytes.NewBufferString("")
// 	command.SetArgs([]string{
// 		"--file", "./tests/gomtpYamls/invalidCredentialsYandex.yaml",
// 	})
// 	command.SetOut(b)
// 	command.SetErr(b)
// 	command.Execute()
// 	var expected string = "authentication failed: Invalid user or password"
// 	assert.Contains(t, b.String(), expected, "Invalid credentials error expected.")

// }

// func TestInvalidCredentialsBrevo(t *testing.T) {
// 	command := rootCmd
// 	b := bytes.NewBufferString("")
// 	command.SetArgs([]string{
// 		"--file", "./tests/gomtpYamls/invalidCredentialsBrevo.yaml",
// 	})
// 	command.SetOut(b)
// 	command.SetErr(b)
// 	command.Execute()
// 	var expected string = "Authentication failed"
// 	assert.Contains(t, b.String(), expected, "Invalid credentials error expected.")

// }

// func TestInvalidSSLConfiguration(t *testing.T) {
// 	command := rootCmd
// 	b := bytes.NewBufferString("")
// 	command.SetArgs([]string{
// 		"--file", "./tests/gomtpYamls/nonSslServerWithSslConfiguration.yaml",
// 	})
// 	command.SetOut(b)
// 	command.SetErr(b)
// 	command.Execute()
// 	var expected string = "tls: first record does not look like a TLS handshake"
// 	assert.Contains(t, b.String(), expected, "SSL error expected.")
// }

// func TestNonTLSServerWithTLSConfiguration(t *testing.T) {

// }

// func TestOptionalParameters(t *testing.T) {
// 	command := rootCmd
// 	command.SetArgs([]string{
// 		"--file", "./tests/gomtpYamls/optionalParameters.yaml",
// 	})
// 	b := bytes.NewBufferString("")
// 	command.SetOut(b)
// 	command.SetErr(b)
// 	command.Execute()
// 	var expected string = "Email sent successfully!"
// 	assert.Equal(t, expected, b.String(), "actual is not expected")
// }

// func TestInvalidYaml(t *testing.T) {
// 	command := rootCmd
// 	command.SetArgs([]string{
// 		"--file", "./tests/gomtpYamls/invalid.yaml",
// 	})
// 	b := bytes.NewBufferString("")
// 	command.SetOut(b)
// 	command.SetErr(b)
// 	err := command.Execute()
// 	assert.NotEmpty(t, err)
// }

func TestEmptySubjectYaml(t *testing.T) {
	clearMailHog(t)
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
	//time.Sleep(1 * time.Second)

	// Check the latest message (first in the list)
	latestMessage := getLatestMessageForRecipient(t, "emptySubject@example.com")

	assert.Equal(t, "GOMTP Test Subject", latestMessage.Subject, "unexpected email subject")
	assert.Equal(t, "from@example.com", latestMessage.From.Address, "unexpected sender mailbox")
	assert.Equal(t, "emptySubject@example.com", latestMessage.To[0].Address, "unexpected recipient mailbox")
	assert.Equal(t, "this is line 1 This is line 2", latestMessage.Snippet, "unexpected email body")
}

func TestEmptyToYaml(t *testing.T) {
	clearMailHog(t)

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

	// Wait a moment for MailHog to process the email
	//time.Sleep(1 * time.Second)

	latestMessage := getLatestMessageForRecipient(t, "to@example.com")
	assert.Equal(t, "Testing Email", latestMessage.Subject, "unexpected email subject")
	assert.Equal(t, "emptyToFrom@example.com", latestMessage.From.Address, "unexpected sender domain")
	assert.Equal(t, "to@example.com", latestMessage.To[0].Address, "unexpected recipient domain")
	assert.Equal(t, "Empty To Test Body", latestMessage.Snippet, "unexpected email body")
}

func TestEmptyBodyYaml(t *testing.T) {
	clearMailHog(t)

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

	// Wait a moment for MailHog to process the email
	//time.Sleep(1 * time.Second)

	latestMessage := getLatestMessageForRecipient(t, "emptyBodyTo@example.com")
	assert.Equal(t, "Testing Email For Empty Body", latestMessage.Subject, "unexpected email subject")
	assert.Equal(t, "emptyBodyFrom@example.com", latestMessage.From.Address, "unexpected sender domain")
	assert.Equal(t, "emptyBodyTo@example.com", latestMessage.To[0].Address, "unexpected recipient domain")
	assert.Equal(t, "This is the test email sent by gomtp.", latestMessage.Snippet, "unexpected email body")
}

func TestToFlag(t *testing.T) {
	clearMailHog(t)

	command := rootCmd
	command.SetArgs([]string{
		"--file", "./tests/gomtpYamls/successConfigurationWithoutTo.yaml",
		"--to", "to-flag-test@example.com",
	})
	b := bytes.NewBufferString("")
	command.SetOut(b)
	command.SetErr(b)
	err := command.Execute()
	assert.NoError(t, err, "command execution failed")

	assert.Equal(t, "Email sent successfully!", b.String(), "unexpected command output")

	// Wait a moment for MailHog to process the email
	//time.Sleep(1 * time.Second)

	latestMessage := getLatestMessageForRecipient(t, "to-flag-test@example.com")
	assert.Equal(t, "To Flag Test Subject", latestMessage.Subject, "unexpected email subject")
	assert.Equal(t, "successConfigurationWithoutToFrom@example.com", latestMessage.From.Address, "unexpected sender domain")
	assert.Equal(t, "to-flag-test@example.com", latestMessage.To[0].Address, "unexpected recipient domain")
	assert.Equal(t, "To Flag Test Body", latestMessage.Snippet, "unexpected email body")
}

func TestSubjectFlag(t *testing.T) {
	clearMailHog(t)

	command := rootCmd
	command.SetArgs([]string{
		"--file", "./tests/gomtpYamls/successConfigurationWithoutSubject.yaml",
		"--to", "",
		"--subject", "Subject To Flag Test Subject",
	})
	b := bytes.NewBufferString("")
	command.SetOut(b)
	command.SetErr(b)
	err := command.Execute()
	assert.NoError(t, err, "command execution failed")

	assert.Equal(t, "Email sent successfully!", b.String(), "unexpected command output")

	// Wait a moment for MailHog to process the email
	//time.Sleep(1 * time.Second)

	latestMessage := getLatestMessageForRecipient(t, "successConfigurationWithoutSubjectTo@example.com")
	assert.Equal(t, "Subject To Flag Test Subject", latestMessage.Subject, "unexpected email subject")
	assert.Equal(t, "successConfigurationWithoutSubjectFrom@example.com", latestMessage.From.Address, "unexpected sender domain")
	assert.Equal(t, "successConfigurationWithoutSubjectTo@example.com", latestMessage.To[0].Address, "unexpected recipient domain")
	assert.Equal(t, "Subject To Flag Test Body", latestMessage.Snippet, "unexpected email body")
}

func TestSubjectToBodyFlag(t *testing.T) {
	clearMailHog(t)

	// Setup the command with arguments
	command := rootCmd
	command.SetArgs([]string{
		"--file", "./tests/gomtpYamls/successConfigurationWithoutSubjectBody.yaml",
		"--to", "subjecttobodyflag@example.net",
		"--subject", "Subject To Body Flag Test Subject",
		"--body", "Subject To Body Flag Test Body",
	})

	// Capture the output
	b := bytes.NewBufferString("")
	command.SetOut(b)
	command.SetErr(b)

	// Execute the command
	err := command.Execute()
	assert.NoError(t, err, "command execution failed")

	// Verify the output
	assert.Equal(t, "Email sent successfully!", b.String(), "unexpected command output")

	// Wait a moment for MailHog to process the email
	//time.Sleep(1 * time.Second)

	latestMessage := getLatestMessageForRecipient(t, "subjecttobodyflag@example.net")
	assert.Equal(t, "Subject To Body Flag Test Subject", latestMessage.Subject, "unexpected email subject")
	assert.Equal(t, "from@example.com", latestMessage.From.Address, "unexpected sender domain")
	assert.Equal(t, "subjecttobodyflag@example.net", latestMessage.To[0].Address, "unexpected recipient domain")
	assert.Equal(t, "Subject To Body Flag Test Body", latestMessage.Snippet, "unexpected email body")
}

func TestStdinInput(t *testing.T) {
	clearMailHog(t)
	// Save the original stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Ensure stdin is restored after the test

	// Create a pipe to simulate stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	defer r.Close() // Ensure the read end is closed properly

	os.Stdin = r

	// Write to the pipe asynchronously
	go func() {
		defer w.Close() // Ensure the write end is closed properly
		_, err := w.Write([]byte("Body from stdin"))
		if err != nil {
			t.Errorf("failed to write to pipe: %v", err)
		}
	}()

	// Setup the command with arguments
	command := rootCmd
	command.SetArgs([]string{
		"--file", "./tests/gomtpYamls/successConfigurationWithNoBody.yaml",
		"--subject", "Body From STDIN Test Subject",
		"--body", "",
		"--body-file", "",
		"--to", "bodyFromStdin@example.io",
	})

	// Capture the output
	var outputBuffer bytes.Buffer
	command.SetOut(&outputBuffer)
	command.SetErr(&outputBuffer)

	// Execute the command
	err = command.Execute()
	assert.NoError(t, err, "command execution failed")

	// Check the output
	assert.Contains(t, outputBuffer.String(), "Email sent successfully!", "unexpected command output")

	// Wait a moment for MailHog to process the email
	//time.Sleep(1 * time.Second)

	// Verify the email content
	latestMessage := getLatestMessageForRecipient(t, "bodyFromStdin@example.io")
	assert.Equal(t, "Body From STDIN Test Subject", latestMessage.Subject, "unexpected email subject")
	assert.Equal(t, "from@example.com", latestMessage.From.Address, "unexpected sender domain")
	assert.Equal(t, "bodyFromStdin@example.io", latestMessage.To[0].Address, "unexpected recipient domain")
	assert.Equal(t, "Body from stdin", latestMessage.Snippet, "unexpected email body")
}
