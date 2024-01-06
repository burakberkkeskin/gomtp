package cmd

import (
	_ "embed"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

//go:embed embeddedFiles/mailhog.yaml
var mailhogGomtpYamlTmpl []byte

//go:embed embeddedFiles/gmail.yaml
var gmailGomtpYamlTmpl []byte

//go:embed embeddedFiles/yandex.yaml
var yandexGomtpYamlTmpl []byte

//go:embed embeddedFiles/brevo.yaml
var brevoGomtpYamlTmpl []byte

var (
	gomtpTemplatePath string
	providerName      string
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Create a gomtp yaml template file.",
	RunE:  templateCmdFunction,
}

func templateCmdFunction(cmd *cobra.Command, args []string) error {
	template, err := checkProvider(providerName)
	if err != nil {
		return err
	}

	err = os.WriteFile(gomtpTemplatePath, template, 0644)
	if err != nil {
		return err
	}
	return nil
}

func checkProvider(providerName string) ([]byte, error) {
	switch providerName {
	case "mailhog":
		return mailhogGomtpYamlTmpl, nil
	case "gmail":
		return gmailGomtpYamlTmpl, nil
	case "yandex":
		return yandexGomtpYamlTmpl, nil
	case "brevo":
		return brevoGomtpYamlTmpl, nil
	default:
		return nil, errors.New("provider can be one of these: gomtp | gmail | yandex | brevo")
	}
}

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.Flags().StringVarP(&gomtpTemplatePath, "output", "o", "gomtp.yaml", "Output path of gomtp template yaml file.")
	templateCmd.Flags().StringVarP(&providerName, "provider", "p", "mailhog", "Provider for the template.")
}
