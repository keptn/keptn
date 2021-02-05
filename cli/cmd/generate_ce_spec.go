package cmd

import (
	"fmt"
	"github.com/keptn/keptn/cli/internal/cespec"
	"github.com/spf13/cobra"
	"os"
)

var generateCeSpecParams *generateCmdParams

// generateCESpecCmd implements the generate cloud-events-spec command
var generateCESpecCmd = &cobra.Command{
	Use:   "cloud-events-spec",
	Short: "Generates the markdown documentation for the Keptn CloudEvents",
	Long: `Generates markdown documentation for the Keptn CloudEvents.

This command can be used to create an up-to-date documentation of the Keptn CloudEvents documentation at
https://github.com/keptn/spec/blob/master/cloudevents.md

It creates one markdown file containing a description of the CloudEvents as well as examples and json schema definitions.
`,
	Example: `keptn generate cloud-events-spec

keptn generate cloud-events-spec --dir=/some/directory`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		outputDir := "./ce-docs"
		if *generateCeSpecParams.Directory != "" {
			outputDir = *generateCeSpecParams.Directory
		}

		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			return fmt.Errorf("error trying to access directory %s. Please make sure the directory exists", outputDir)
		}

		fmt.Println("Generating cloud-events-spec now...")
		cespec.Generate(outputDir)
		fmt.Printf("Docs have been written to: %s\n", outputDir)

		return nil
	},
}

func init() {
	generateCmd.AddCommand(generateCESpecCmd)

	generateCeSpecParams = &generateCmdParams{}
	generateCeSpecParams.Directory = generateCESpecCmd.Flags().StringP("dir", "", "./ce-docs", "directory where the cloud events spec should be written to")
}
