package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHappyPathTemplateCommand(t *testing.T) {
	command := rootCmd
	customYamlPath := "./custom-path.yaml"
	command.SetArgs([]string{
		"template",
		"--output", customYamlPath,
	})
	err := command.Execute()
	assert.Nil(t, err)
	os.Remove(customYamlPath)
}

func TestInvalidPathTemplateCommand(t *testing.T) {
	command := rootCmd
	customYamlPath := "/non/existing/path"
	command.SetArgs([]string{
		"template",
		"--output", customYamlPath,
	})
	b := bytes.NewBufferString("")
	command.SetOut(b)
	command.SetErr(b)
	var expected string = "no such file or directory"
	err := command.Execute()
	assert.NotNil(t, err)
	assert.Contains(t, b.String(), expected)
	os.Remove(customYamlPath)
}

func TestDefaultTemplateProvider(t *testing.T) {
	command := rootCmd
	customYamlPath := "./test.yaml"
	command.SetArgs([]string{
		"template",
		"--output", customYamlPath,
	})
	b := bytes.NewBufferString("")
	command.SetOut(b)
	command.SetErr(b)
	err := command.Execute()
	assert.Nil(t, err)
	data, err := os.ReadFile(customYamlPath)
	fileContent := string(data)
	assert.Nil(t, err)
	assert.Contains(t, fileContent, "username: ''")
	assert.Contains(t, fileContent, "host: '127.0.0.1'")
	assert.Contains(t, fileContent, "port: 1025")
	assert.Contains(t, fileContent, "ssl: false")
	assert.Contains(t, fileContent, "tls: false")
	assert.Contains(t, fileContent, "auth: 'NO'")
	os.Remove(customYamlPath)
}

func TestGmailTemplateProvider(t *testing.T) {
	command := rootCmd
	customYamlPath := "./test.yaml"
	command.SetArgs([]string{
		"template",
		"--output", customYamlPath,
		"--provider", "gmail",
	})
	b := bytes.NewBufferString("")
	command.SetOut(b)
	command.SetErr(b)
	err := command.Execute()
	assert.Nil(t, err)
	data, err := os.ReadFile(customYamlPath)
	fileContent := string(data)
	assert.Nil(t, err)
	assert.Contains(t, fileContent, "username: 'from@gmail.com'")
	assert.Contains(t, fileContent, "host: 'smtp.gmail.com'")
	assert.Contains(t, fileContent, "port: 587")
	assert.Contains(t, fileContent, "ssl: false")
	assert.Contains(t, fileContent, "tls: true")
	assert.Contains(t, fileContent, "auth: 'LOGIN'")
	os.Remove(customYamlPath)
}

func TestYandexTemplateProvider(t *testing.T) {
	command := rootCmd
	customYamlPath := "./test.yaml"
	command.SetArgs([]string{
		"template",
		"--output", customYamlPath,
		"--provider", "yandex",
	})
	b := bytes.NewBufferString("")
	command.SetOut(b)
	command.SetErr(b)
	err := command.Execute()
	assert.Nil(t, err)
	data, err := os.ReadFile(customYamlPath)
	fileContent := string(data)
	assert.Nil(t, err)
	assert.Contains(t, fileContent, "username: 'example@yandex.com'")
	assert.Contains(t, fileContent, "host: 'smtp.yandex.com'")
	assert.Contains(t, fileContent, "port: 465")
	assert.Contains(t, fileContent, "ssl: true")
	assert.Contains(t, fileContent, "tls: false")
	assert.Contains(t, fileContent, "auth: 'LOGIN'")
	os.Remove(customYamlPath)
}

func TestBrevoTemplateProvider(t *testing.T) {
	command := rootCmd
	customYamlPath := "./test.yaml"
	command.SetArgs([]string{
		"template",
		"--output", customYamlPath,
		"--provider", "brevo",
	})
	b := bytes.NewBufferString("")
	command.SetOut(b)
	command.SetErr(b)
	err := command.Execute()
	assert.Nil(t, err)
	data, err := os.ReadFile(customYamlPath)
	fileContent := string(data)
	assert.Nil(t, err)
	assert.Contains(t, fileContent, "username: 'user@example.com'")
	assert.Contains(t, fileContent, "host: 'smtp-relay.brevo.com'")
	assert.Contains(t, fileContent, "port: 587")
	assert.Contains(t, fileContent, "ssl: false")
	assert.Contains(t, fileContent, "tls: true")
	assert.Contains(t, fileContent, "auth: 'LOGIN'")
	os.Remove(customYamlPath)
}

// run with --output /etc/gomtp.yaml and check if exit code is non zero
func TestPermissionErrorPathTemplateCommand(t *testing.T) {

}
