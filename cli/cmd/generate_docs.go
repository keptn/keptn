// Inspired by `hugo gen doc`  - see https://github.com/gohugoio/hugo/blob/release-0.69.0/commands/gendoc.go
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)


type generateDocsCmdParams struct {
	Directory  *string
}

var generateDocsParams *generateDocsCmdParams

const gendocFrontmatterTemplate = `---
date: "%s"
title: "%s"
slug: %s
---
`

// crprojectCmd represents the project command
var generateDocsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generates Markdown documentation for the keptn CLI.",
	Long: `Generates Markdown documentation for the keptn CLI.

This command can be used to create an up-to-date documentation of Keptn's command-line interface for https://keptn.sh.

It creates one Markdown file per command, suitable for rendering in Hugo.
`,
	Example: `keptn generate docs
keptn generate docs --dir=/some/directory`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// generate current datetime for prepender
		now := time.Now().Format(time.RFC3339)

		outputDir := "./docs"
		if *generateDocsParams.Directory != "" {
			outputDir = *generateDocsParams.Directory
		}

		// check if output directory exists
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			// outputDir does not exist
			return err
		}

		// define a prepender function for compatibility with Hugo templates
		prepender := func(filename string) string {
			name := filepath.Base(filename)
			base := strings.TrimSuffix(name, path.Ext(name))
			return fmt.Sprintf(gendocFrontmatterTemplate, now, strings.Replace(base, "_", " ", -1), base)
		}

		// links need to be converted to be compatible with hugo
		linkHandler := func(name string) string {
			base := strings.TrimSuffix(name, path.Ext(name))
			return "../" + strings.ToLower(base) + "/"
		}


		// generate docs
		fmt.Println("Generating docs now...")
		doc.GenMarkdownTreeCustom(cmd.Root(), outputDir, prepender, linkHandler)
		fmt.Printf("Docs have been written to %s!\n", outputDir)

		return nil
	},
}

func init() {
	generateCmd.AddCommand(generateDocsCmd)

	generateDocsParams = &generateDocsCmdParams{}
	generateDocsParams.Directory = generateDocsCmd.Flags().StringP("dir", "", "", "directory where the docs should be written to")
}
