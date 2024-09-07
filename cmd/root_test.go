package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
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
}

type MailhogAddress struct {
	Name    string `json:"Name"`
	Address string `json:"Address"`
}

func getLatestMessageForRecipient(recipient string) (MailhogMessage, error) {
	resp, err := http.Get("http://localhost:8025/api/v1/messages")

	if err != nil {
		return MailhogMessage{}, err
	}

	defer resp.Body.Close()

	var mailhogResp MailhogResponse
	err = json.NewDecoder(resp.Body).Decode(&mailhogResp)

	if err != nil {
		return MailhogMessage{}, err
	}

	for _, msg := range mailhogResp.Items {
		if msg.To[0].Address == recipient {
			return msg, nil
		}
	}

	return MailhogMessage{}, nil
}

type TestGOMTPSuite struct {
	suite.Suite
	cmd cobra.Command
}

func TestGOMTPSuiteTest(t *testing.T) {
	suite.Run(t, new(TestGOMTPSuite))
}

func (suite *TestGOMTPSuite) SetupTest() {
	suite.cmd = *rootCmd
}

func (suite *TestGOMTPSuite) TestHappyPath() {
	expected := "Email sent successfully!"

	suite.cmd.SetArgs([]string{
		"--file", "../tests/gomtpYamls/successConfiguration.yaml",
		"--body", "",
		"--body-file", "",
		"--subject", "",
		"--to", "",
	})

	b := bytes.NewBufferString("")
	suite.cmd.SetOut(b)
	suite.cmd.SetErr(b)
	suite.cmd.Execute()

	suite.Equal(expected, b.String())
}

func (suite *TestGOMTPSuite) TestBodyFileFlag() {

	// Setup the command with arguments
	suite.cmd.SetArgs([]string{
		"--file", "../tests/gomtpYamls/successConfigurationWithNoBodyFile.yaml",
		"--to", "bodyfileflag@example.net",
		"--subject", "Body File Flag Test Subject",
		"--body-file", "../tests/gomtpYamls/emailBodyFile.log",
		"--body", "",
	})

	// Capture the output
	b := bytes.NewBufferString("")
	suite.cmd.SetOut(b)
	suite.cmd.SetErr(b)

	// Execute the command
	err := suite.cmd.Execute()
	suite.NoError(err)

	// Verify the output
	suite.Equal("Email sent successfully!", b.String())

	// Wait a moment for MailHog to process the email
	//time.Sleep(1 * time.Second)

	latestMessage, err := getLatestMessageForRecipient("bodyfileflag@example.net")
	suite.NoError(err)
	suite.Equal("Body File Flag Test Subject", latestMessage.Subject)
	suite.Equal("from@example.com", latestMessage.From.Address)
	suite.Equal("bodyfileflag@example.net", latestMessage.To[0].Address)
	suite.Equal("Test body file content", latestMessage.Snippet)
}

func (suite *TestGOMTPSuite) TestGomtpYamlNotFound() {
	suite.cmd.SetArgs([]string{
		"--file", "../unknown/path.yaml",
		"--to", "",
		"--body", "",
		"--body-file", "",
		"--subject", "",
	})

	b := bytes.NewBufferString("")
	suite.cmd.SetOut(b)
	suite.cmd.SetErr(b)
	suite.cmd.Execute()

	expected := "no such file or directory"
	suite.Contains(b.String(), expected, "File not found error expected.")
}

func (suite *TestGOMTPSuite) TestInvalidCredentialsGoogle() {
	suite.cmd.SetArgs([]string{
		"--file", "../tests/gomtpYamls/invalidCredentialsGoogle.yaml",
		"--to", "",
		"--body", "",
		"--body-file", "",
		"--subject", "",
	})

	b := bytes.NewBufferString("")
	suite.cmd.SetOut(b)
	suite.cmd.SetErr(b)
	suite.cmd.Execute()

	expected := "Username and Password not accepted"
	suite.Contains(b.String(), expected, "Invalid credentials error expected.")
}

func (suite *TestGOMTPSuite) TestInvalidCredentialsYandex() {
	suite.cmd.SetArgs([]string{
		"--file", "../tests/gomtpYamls/invalidCredentialsYandex.yaml",
		"--to", "",
		"--body", "",
		"--body-file", "",
		"--subject", "",
	})

	b := bytes.NewBufferString("")
	suite.cmd.SetOut(b)
	suite.cmd.SetErr(b)
	suite.cmd.Execute()

	expected := "authentication failed: Invalid user or password"
	suite.Contains(b.String(), expected, "Invalid credentials error expected.")
}

func (suite *TestGOMTPSuite) TestInvalidCredentialsBrevo() {
	suite.cmd.SetArgs([]string{
		"--file", "../tests/gomtpYamls/invalidCredentialsBrevo.yaml",
		"--to", "",
		"--body", "",
		"--body-file", "",
		"--subject", "",
	})

	b := bytes.NewBufferString("")
	suite.cmd.SetOut(b)
	suite.cmd.SetErr(b)
	suite.cmd.Execute()

	expected := "Authentication failed"
	suite.Contains(b.String(), expected, "Invalid credentials error expected.")
}

func (suite *TestGOMTPSuite) TestInvalidSSLConfiguration() {
	suite.cmd.SetArgs([]string{
		"--file", "../tests/gomtpYamls/nonSslServerWithSslConfiguration.yaml",
		"--to", "",
		"--body", "",
		"--body-file", "",
		"--subject", "",
	})

	b := bytes.NewBufferString("")
	suite.cmd.SetOut(b)
	suite.cmd.SetErr(b)
	suite.cmd.Execute()

	expected := "tls: first record does not look like a TLS handshake"
	suite.Contains(b.String(), expected, "SSL error expected.")
}

func (suite *TestGOMTPSuite) TestOptionalParameters() {
	suite.cmd.SetArgs([]string{
		"--file", "../tests/gomtpYamls/optionalParameters.yaml",
		"--to", "",
		"--body", "",
		"--body-file", "",
		"--subject", "",
	})

	b := bytes.NewBufferString("")
	suite.cmd.SetOut(b)
	suite.cmd.SetErr(b)
	suite.cmd.Execute()

	expected := "Email sent successfully!"
	suite.Equal(expected, b.String(), "actual is not expected")
}

func (suite *TestGOMTPSuite) TestInvalidYaml() {
	suite.cmd.SetArgs([]string{
		"--file", "../tests/gomtpYamls/invalid.yaml",
		"--to", "",
		"--body", "",
		"--body-file", "",
		"--subject", "",
	})

	b := bytes.NewBufferString("")
	suite.cmd.SetOut(b)
	suite.cmd.SetErr(b)
	err := suite.cmd.Execute()

	suite.NotEmpty(err)
}

func (suite *TestGOMTPSuite) TestEmptySubjectYaml() {
	suite.cmd.SetArgs([]string{
		"--file", "../tests/gomtpYamls/emptySubject.yaml",
		"--to", "",
		"--body", "",
		"--body-file", "",
		"--subject", "",
	})

	b := bytes.NewBufferString("")
	suite.cmd.SetOut(b)
	suite.cmd.SetErr(b)
	err := suite.cmd.Execute()
	suite.NoError(err)

	expected := "Email sent successfully!"
	suite.Equal(expected, b.String())

	latestMessage, err := getLatestMessageForRecipient("emptySubject@example.com")
	suite.NoError(err)
	suite.Equal("GOMTP Test Subject", latestMessage.Subject)
	suite.Equal("from@example.com", latestMessage.From.Address)
	suite.Equal("emptySubject@example.com", latestMessage.To[0].Address)
}

func (suite *TestGOMTPSuite) TestEmptyToYaml() {
	suite.cmd.SetArgs([]string{
		"--file", "../tests/gomtpYamls/emptyTo.yaml",
		"--to", "",
		"--body", "",
		"--body-file", "",
		"--subject", "",
	})

	b := bytes.NewBufferString("")
	suite.cmd.SetOut(b)
	suite.cmd.SetErr(b)
	err := suite.cmd.Execute()
	suite.NoError(err)

	suite.Equal("Email sent successfully!", b.String())

	latestMessage, err := getLatestMessageForRecipient("to@example.com")
	suite.NoError(err)
	suite.Equal("Testing Email", latestMessage.Subject)
	suite.Equal("emptyToFrom@example.com", latestMessage.From.Address)
	suite.Equal("to@example.com", latestMessage.To[0].Address)
}

func (suite *TestGOMTPSuite) TestEmptyBodyYaml() {
	suite.cmd.SetArgs([]string{
		"--file", "../tests/gomtpYamls/emptyBody.yaml",
		"--to", "",
		"--body", "",
		"--body-file", "",
		"--subject", "",
	})

	b := bytes.NewBufferString("")
	suite.cmd.SetOut(b)
	suite.cmd.SetErr(b)
	err := suite.cmd.Execute()
	suite.NoError(err)

	suite.Equal("Email sent successfully!", b.String())

	latestMessage, err := getLatestMessageForRecipient("emptyBodyTo@example.com")
	suite.NoError(err)
	suite.Equal("Testing Email For Empty Subject", latestMessage.Subject)
	suite.Equal("emptyBodyFrom@example.com", latestMessage.From.Address)
	suite.Equal("emptyBodyTo@example.com", latestMessage.To[0].Address)
}

func (suite *TestGOMTPSuite) TestNonTLSServerWithTLSConfiguration() {

}

func (suite *TestGOMTPSuite) TestToFlag() {

	suite.cmd.SetArgs([]string{
		"--file", "../tests/gomtpYamls/successConfigurationWithoutTo.yaml",
		"--to", "to-flag-test@example.com",
		"--body", "",
		"--body-file", "",
		"--subject", "",
	})

	b := bytes.NewBufferString("")
	suite.cmd.SetOut(b)
	suite.cmd.SetErr(b)
	err := suite.cmd.Execute()

	suite.NoError(err)

	expected := "Email sent successfully!"
	suite.Equal(expected, b.String())

	latestMessage, err := getLatestMessageForRecipient("to-flag-test@example.com")
	suite.NoError(err)
	suite.Equal("To Flag Test Subject", latestMessage.Subject)
	suite.Equal("successConfigurationWithoutToFrom@example.com", latestMessage.From.Address)
	suite.Equal("to-flag-test@example.com", latestMessage.To[0].Address)
}

func (suite *TestGOMTPSuite) TestSubjectFlag() {

	suite.cmd.SetArgs([]string{
		"--file", "../tests/gomtpYamls/successConfigurationWithoutSubject.yaml",
		"--to", "successConfigurationWithoutSubjectTo@example.com",
		"--subject", "Subject To Flag Test Subject",
		"--body", "",
		"--body-file", "",
	})

	b := bytes.NewBufferString("")
	suite.cmd.SetOut(b)
	suite.cmd.SetErr(b)
	err := suite.cmd.Execute()

	suite.NoError(err)

	expected := "Email sent successfully!"
	suite.Equal(expected, b.String())

	latestMessage, err := getLatestMessageForRecipient("successConfigurationWithoutSubjectTo@example.com")
	suite.NoError(err)
	suite.Equal("Subject To Flag Test Subject", latestMessage.Subject)
	suite.Equal("successConfigurationWithoutSubjectFrom@example.com", latestMessage.From.Address)
	suite.Equal("successConfigurationWithoutSubjectTo@example.com", latestMessage.To[0].Address)
}

func (suite *TestGOMTPSuite) TestSubjectToBodyFlag() {

	suite.cmd.SetArgs([]string{
		"--file", "../tests/gomtpYamls/successConfigurationWithoutSubjectBody.yaml",
		"--to", "subjecttobodyflag@example.net",
		"--subject", "Subject To Body Flag Test Subject",
		"--body", "Subject To Body Flag Test Body",
		"--body-file", "",
	})

	b := bytes.NewBufferString("")
	suite.cmd.SetOut(b)
	suite.cmd.SetErr(b)
	err := suite.cmd.Execute()

	suite.NoError(err)

	expected := "Email sent successfully!"
	suite.Equal(expected, b.String())

	latestMessage, err := getLatestMessageForRecipient("subjecttobodyflag@example.net")
	suite.NoError(err)
	suite.Equal("Subject To Body Flag Test Subject", latestMessage.Subject)
	suite.Equal("from@example.com", latestMessage.From.Address)
	suite.Equal("subjecttobodyflag@example.net", latestMessage.To[0].Address)
}

func (suite *TestGOMTPSuite) TestStdinInput() {

	// Save the original stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Create a pipe to simulate stdin
	r, w, err := os.Pipe()
	suite.NoError(err)
	defer r.Close()

	os.Stdin = r

	// Write to the pipe asynchronously
	go func() {
		defer w.Close()
		_, err := w.Write([]byte("Body from stdin"))
		suite.NoError(err)
	}()

	// Setup the command with arguments
	suite.cmd.SetArgs([]string{
		"--file", "../tests/gomtpYamls/successConfigurationWithNoBody.yaml",
		"--subject", "Body From STDIN Test Subject",
		"--body", "",
		"--body-file", "",
		"--to", "bodyFromStdin@example.io",
	})

	// Capture the output
	var outputBuffer bytes.Buffer
	suite.cmd.SetOut(&outputBuffer)
	suite.cmd.SetErr(&outputBuffer)

	err = suite.cmd.Execute()
	suite.NoError(err)

	// Verify the output
	suite.Contains(outputBuffer.String(), "Email sent successfully!", "unexpected command output")

	latestMessage, err := getLatestMessageForRecipient("bodyFromStdin@example.io")
	suite.NoError(err)
	suite.Equal("Body From STDIN Test Subject", latestMessage.Subject)
	suite.Equal("from@example.com", latestMessage.From.Address)
	suite.Equal("bodyFromStdin@example.io", latestMessage.To[0].Address)
}

func (suite *TestGOMTPSuite) TestBodyFileAndBodyFlag() {

	suite.cmd.SetArgs([]string{
		"--file", "../tests/gomtpYamls/successConfigurationWithNoBodyFile.yaml",
		"--body", "should fail body",
		"--to", "bodyfileandbodyflag@example.net",
		"--subject", "Body File and Body Flag",
		"--body-file", "../tests/gomtpYamls/emailBodyFile.log",
	})

	b := bytes.NewBufferString("")
	suite.cmd.SetOut(b)
	suite.cmd.SetErr(b)

	err := suite.cmd.Execute()
	suite.Error(err)

	expected := "cannot specify body via multiple sources simultaneously"
	suite.Contains(b.String(), expected, "unexpected command output")
}
