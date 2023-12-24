package cmd

import (
	"bytes"
	"gomtp/cmd"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHappyPath(t *testing.T) {
	command := cmd.RootCmd
	b := bytes.NewBufferString("")
	command.SetOut(b)
	command.SetErr(b)
	command.Execute()
	var expected string = "Email sent successfully!"
	assert.Equal(t, expected, b.String(), "actual is not expected")
}

func TestGomtpYamlNotFound(t *testing.T) {

}

func TestInvalidCredentials(t *testing.T) {

}

func TestInvalidSSLConfiguration(t *testing.T) {
	command := cmd.RootCmd
	b := bytes.NewBufferString("")
	command.SetArgs([]string{
		"--file", "./tests/gomtpYamls/nonSslServerWithSslConfiguration.yaml",
	})
	command.SetOut(b)
	command.SetErr(b)
	command.Execute()
	var expected string = "tls: first record does not look like a TLS handshake"
	t.Logf("%s", b.String())
	assert.Contains(t, b.String(), expected, "SSL error expected.")
}

func TestNonTLSServerWithTLSConfiguration(t *testing.T) {

}
