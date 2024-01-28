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

const templateUsageMessage = `Example commands:
  gomtp template # Create a file named gomtp.yaml filled with configuration for mailhog.
  gomtp template -o custom.yaml # Create a file named custom.yaml filled with configuration for mailhog.
  gomtp template -p gmail # Create a file named gomtp.yaml filled with configuration for gmail.
  gomtp template -p yandex # Create a file named gomtp.yaml filled with configuration for gmail.
  gomtp template -p brevo # Create a file named gomtp.yaml filled with configuration for brevo.
  gomtp template -p brevo -o custom.yaml # Create a file named custom.yaml filled with configuration for gmail.
`

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Create a gomtp yaml template file.",
	Long:  templateUsageMessage,
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
		return nil, errors.New("provider can be one of these: mailhog | gmail | yandex | brevo")
	}
}

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.Flags().StringVarP(&gomtpTemplatePath, "output", "o", "gomtp.yaml", "Output path of gomtp template yaml file.")
	templateCmd.Flags().StringVarP(&providerName, "provider", "p", "mailhog", "Provider for the template.")
}
