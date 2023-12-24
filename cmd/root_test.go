package cmd

import (
	"bytes"
	"testing"

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
	var expected string = "Invalid user"
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
