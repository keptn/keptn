package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

var uninstallVersion *string

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:          "uninstall",
	Short:        "Uninstalls keptn",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		if insecureSkipTLSVerify {
			kubectlOptions = "--insecure-skip-tls-verify=true"
		}

		ctx, _ := getKubeContext()
		fmt.Println("Your kubernetes current context is configured to cluster: " + strings.TrimSpace(ctx))
		fmt.Println("Would you like to uninstall keptn from this cluster? (y/n)")

		reader := bufio.NewReader(os.Stdin)
		in, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		in = strings.TrimSpace(in)
		if in != "y" && in != "yes" {
			return nil
		}

		logging.PrintLog("Starting to uninstall keptn", logging.InfoLevel)

		if !mocking {
			// Delete installer pod, ignore if not found
			if err := deleteKeptnInstallerPod("default"); err != nil {
				return err
			}
			// Clean up keptn namespace
			if err := deleteResources("keptn"); err != nil {
				return err
			}
			if err := deleteNamespace("keptn"); err != nil {
				return err
			}
			// Clean up keptn-datastore namespace
			if err := deleteResources("keptn-datastore"); err != nil {
				return err
			}
			if err := deleteNamespace("keptn-datastore"); err != nil {
				return err
			}
		}
		logging.PrintLog("Successfully uninstalled keptn", logging.InfoLevel)

		return nil
	},
}

func deleteKeptnInstallerPod(namespace string) error {
	o := options{"delete", "job", "installer", "-n", namespace, "--ignore-not-found"}
	o.appendIfNotEmpty(kubectlOptions)
	out, err := keptnutils.ExecuteCommand("kubectl", o)
	out = strings.TrimSpace(out)
	if out != "" {
		logging.PrintLog(out, logging.VerboseLevel)
	}
	return err
}

func deleteResources(namespace string) error {
	o := options{"delete", "services,deployments,pods,secrets,configmaps", "--all", "-n", namespace, "--ignore-not-found"}
	o.appendIfNotEmpty(kubectlOptions)
	out, err := keptnutils.ExecuteCommand("kubectl", o)
	out = strings.TrimSpace(out)
	if out != "" {
		logging.PrintLog(out, logging.VerboseLevel)
	}
	return err
}

func deleteNamespace(namespace string) error {
	o := options{"delete", "namespace", namespace, "--ignore-not-found"}
	o.appendIfNotEmpty(kubectlOptions)
	out, err := keptnutils.ExecuteCommand("kubectl", o)
	out = strings.TrimSpace(out)
	if out != "" {
		logging.PrintLog(out, logging.VerboseLevel)
	}
	return err
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallVersion = uninstallCmd.Flags().StringP("keptn-version", "k", "master",
		"The branch or tag of the version which is used for updating the domain")
	uninstallCmd.Flags().MarkHidden("keptn-version")
	uninstallCmd.PersistentFlags().BoolVarP(&insecureSkipTLSVerify, "insecure-skip-tls-verify", "s", false, "Skip tls verification for kubectl commands")
}
