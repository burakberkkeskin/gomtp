package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
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
	//clearMailHog(t)

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
	//clearMailHog(t)

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
	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://localhost:8025/api/v2/messages")
	assert.NoError(t, err, "failed to get messages from MailHog")
	defer resp.Body.Close()

	var mailhogResp MailhogResponse
	err = json.NewDecoder(resp.Body).Decode(&mailhogResp)
	assert.NoError(t, err, "failed to decode MailHog response")

	assert.NotEmpty(t, mailhogResp.Items, "no messages found in MailHog")

	latestMessage := mailhogResp.Items[0]
	assert.Equal(t, "Testing Email", latestMessage.Content.Headers["Subject"][0], "unexpected email subject")
	assert.Equal(t, "from", latestMessage.From.Mailbox, "unexpected sender mailbox")
	assert.Equal(t, "example.com", latestMessage.From.Domain, "unexpected sender domain")
	assert.Equal(t, "to", latestMessage.To[0].Mailbox, "unexpected recipient mailbox")
	assert.Equal(t, "example.com", latestMessage.To[0].Domain, "unexpected recipient domain")
	assert.Equal(t, "Empty To Test Body", latestMessage.Content.Body, "unexpected email body")
}

func TestEmptyBodyYaml(t *testing.T) {
	//clearMailHog(t)

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
	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://localhost:8025/api/v2/messages")
	assert.NoError(t, err, "failed to get messages from MailHog")
	defer resp.Body.Close()

	var mailhogResp MailhogResponse
	err = json.NewDecoder(resp.Body).Decode(&mailhogResp)
	assert.NoError(t, err, "failed to decode MailHog response")

	assert.NotEmpty(t, mailhogResp.Items, "no messages found in MailHog")

	latestMessage := mailhogResp.Items[0]
	assert.Equal(t, "Testing Email For Empty Body", latestMessage.Content.Headers["Subject"][0], "unexpected email subject")
	assert.Equal(t, "from", latestMessage.From.Mailbox, "unexpected sender mailbox")
	assert.Equal(t, "example.com", latestMessage.From.Domain, "unexpected sender domain")
	assert.Equal(t, "to", latestMessage.To[0].Mailbox, "unexpected recipient mailbox")
	assert.Equal(t, "example.com", latestMessage.To[0].Domain, "unexpected recipient domain")
	assert.Equal(t, "This is the test email sent by gomtp.", latestMessage.Content.Body, "unexpected email body")
}

func TestToFlag(t *testing.T) {
	//clearMailHog(t)

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
	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://localhost:8025/api/v2/messages")
	assert.NoError(t, err, "failed to get messages from MailHog")
	defer resp.Body.Close()

	var mailhogResp MailhogResponse
	err = json.NewDecoder(resp.Body).Decode(&mailhogResp)
	assert.NoError(t, err, "failed to decode MailHog response")

	assert.NotEmpty(t, mailhogResp.Items, "no messages found in MailHog")

	latestMessage := mailhogResp.Items[0]
	assert.Equal(t, "To Flag Test Subject", latestMessage.Content.Headers["Subject"][0], "unexpected email subject")
	assert.Equal(t, "from", latestMessage.From.Mailbox, "unexpected sender mailbox")
	assert.Equal(t, "example.com", latestMessage.From.Domain, "unexpected sender domain")
	assert.Equal(t, "to-flag-test", latestMessage.To[0].Mailbox, "unexpected recipient mailbox")
	assert.Equal(t, "example.com", latestMessage.To[0].Domain, "unexpected recipient domain")
	assert.Equal(t, "To Flag Test Body", latestMessage.Content.Body, "unexpected email body")
}

func TestSubjectToFlag(t *testing.T) {
	//clearMailHog(t)

	command := rootCmd
	command.SetArgs([]string{
		"--file", "./tests/gomtpYamls/successConfigurationWithoutSubject.yaml",
		"--to", "to2@example.org",
		"--subject", "Subject To Flag Test Subject",
	})
	b := bytes.NewBufferString("")
	command.SetOut(b)
	command.SetErr(b)
	err := command.Execute()
	assert.NoError(t, err, "command execution failed")

	assert.Equal(t, "Email sent successfully!", b.String(), "unexpected command output")

	// Wait a moment for MailHog to process the email
	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://localhost:8025/api/v2/messages")
	assert.NoError(t, err, "failed to get messages from MailHog")
	defer resp.Body.Close()

	var mailhogResp MailhogResponse
	err = json.NewDecoder(resp.Body).Decode(&mailhogResp)
	assert.NoError(t, err, "failed to decode MailHog response")

	assert.NotEmpty(t, mailhogResp.Items, "no messages found in MailHog")

	latestMessage := mailhogResp.Items[0]
	assert.Equal(t, "Subject To Flag Test Subject", latestMessage.Content.Headers["Subject"][0], "unexpected email subject")
	assert.Equal(t, "from", latestMessage.From.Mailbox, "unexpected sender mailbox")
	assert.Equal(t, "example.com", latestMessage.From.Domain, "unexpected sender domain")
	assert.Equal(t, "to2", latestMessage.To[0].Mailbox, "unexpected recipient mailbox")
	assert.Equal(t, "example.org", latestMessage.To[0].Domain, "unexpected recipient domain")
	assert.Equal(t, "Subject To Flag Test Body", latestMessage.Content.Body, "unexpected email body")
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

func getLatestMessageForRecipient(t *testing.T, recipient string) MailhogMessage {
	resp, err := http.Get("http://localhost:8025/api/v2/messages")
	assert.NoError(t, err, "failed to get messages from MailHog")
	defer resp.Body.Close()

	var mailhogResp MailhogResponse
	err = json.NewDecoder(resp.Body).Decode(&mailhogResp)
	assert.NoError(t, err, "failed to decode MailHog response")

	assert.NotEmpty(t, mailhogResp.Items, "no messages found in MailHog")

	for _, msg := range mailhogResp.Items {
		if msg.To[0].Mailbox+"@"+msg.To[0].Domain == recipient {
			return msg
		}
	}

	t.Fatalf("No message found for recipient: %s", recipient)
	return MailhogMessage{}
}

func TestSubjectToBodyFlag(t *testing.T) {
	// Save the original stdin
	oldStdin := os.Stdin
	defer func() {
		// Restore the original stdin after the test
		os.Stdin = oldStdin
	}()

	// Optional: reset stdin if it was modified by other tests
	if os.Stdin != oldStdin {
		t.Log("stdin was modified, resetting")
		os.Stdin = oldStdin
	}

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
	time.Sleep(1 * time.Second)

	latestMessage := getLatestMessageForRecipient(t, "subjecttobodyflag@example.net")
	assert.Equal(t, "Subject To Body Flag Test Subject", latestMessage.Content.Headers["Subject"][0], "unexpected email subject")
	assert.Equal(t, "from", latestMessage.From.Mailbox, "unexpected sender mailbox")
	assert.Equal(t, "example.com", latestMessage.From.Domain, "unexpected sender domain")
	assert.Equal(t, "subjecttobodyflag", latestMessage.To[0].Mailbox, "unexpected recipient mailbox")
	assert.Equal(t, "example.net", latestMessage.To[0].Domain, "unexpected recipient domain")
	assert.Equal(t, "Subject To Body Flag Test Body", latestMessage.Content.Body, "unexpected email body")
}

func TestStdinInput(t *testing.T) {
	// Reset the state before the test

	// Clear MailHog state if necessary
	// clearMailHog(t)

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
		"--to", "stdin@example.io",
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
	time.Sleep(1 * time.Second)

	// Verify the email content
	latestMessage := getLatestMessageForRecipient(t, "stdin@example.io")
	assert.Equal(t, "Body From STDIN Test Subject", latestMessage.Content.Headers["Subject"][0], "unexpected email subject")
	assert.Equal(t, "from", latestMessage.From.Mailbox, "unexpected sender mailbox")
	assert.Equal(t, "example.com", latestMessage.From.Domain, "unexpected sender domain")
	assert.Equal(t, "stdin", latestMessage.To[0].Mailbox, "unexpected recipient mailbox")
	assert.Equal(t, "example.io", latestMessage.To[0].Domain, "unexpected recipient domain")
	assert.Equal(t, "Body from stdin", latestMessage.Content.Body, "unexpected email body")
}
