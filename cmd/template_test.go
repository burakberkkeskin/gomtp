package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// run with --output custom-path.yaml and check if exit code is 0
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

// run with --output /non/existing/path and check if exit code is non zero
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

// run with --output /etc/gomtp.yaml and check if exit code is non zero
func TestPermissionErrorPathTemplateCommand(t *testing.T) {

}
