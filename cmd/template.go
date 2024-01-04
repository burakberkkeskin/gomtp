package cmd

import (
	_ "embed"
	"os"

	"github.com/spf13/cobra"
)

//go:embed gomtp.yaml
var gomtpYamlTmpl []byte
var gomtpTemplatePath string

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Create a gomtp yaml template file.",
	RunE:  templateCmdFunction,
}

func templateCmdFunction(cmd *cobra.Command, args []string) error {
	err := os.WriteFile(gomtpTemplatePath, gomtpYamlTmpl, 0644)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.Flags().StringVarP(&gomtpTemplatePath, "output", "o", "gomtp.yaml", "Define the output path of gomtp template yaml file.")
}
