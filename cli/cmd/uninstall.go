// +build !nokubectl

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/keptn/keptn/cli/pkg/helm"

	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/pkg/platform"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/spf13/cobra"
)

type uninstallCmdParams struct {
	Namespace *string
}

var uninstallParams uninstallCmdParams

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstalls Keptn from a Kubernetes cluster",
	Long: `Uninstalls Keptn from a Kubernetes cluster.

This command does *not* delete: 

* Dynatrace monitoring
* Prometheus monitoring
* Any (third-party) service installed in addition to Keptn (e.g., notification-service, slackbot-service, ...)

Besides, deployed services and the configuration on the Git upstream (i.e., GitHub, GitLab, or Bitbucket) are not deleted. To clean-up created projects and services, please see [Delete a project](https://keptn.sh/docs/` + keptnReleaseDocsURL + `/manage/project/#delete-a-project).

**Note:** This command requires a *Kubernetes current context* pointing to the cluster where Keptn should get uninstalled from.
`,
	Example:      `keptn uninstall`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		if insecureSkipTLSVerify {
			kubectlOptions = "--insecure-skip-tls-verify=true"
		}

		keptnNamespace := *uninstallParams.Namespace

		ctx, _ := platform.GetKubeContext()
		fmt.Println("Your Kubernetes current context is configured to cluster: " + strings.TrimSpace(ctx))
		fmt.Println("Would you like to uninstall Keptn from this cluster? (y/n)")

		reader := bufio.NewReader(os.Stdin)
		in, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		in = strings.ToLower(strings.TrimSpace(in))
		if in != "y" && in != "yes" {
			return nil
		}

		logging.PrintLog("Starting to uninstall Keptn", logging.InfoLevel)

		if !mocking {
			if err := helm.NewHelper().UninstallRelease("keptn", keptnNamespace); err != nil {
				return err
			}
			// Clean up keptn namespace
			if err := deleteNamespace(keptnNamespace); err != nil {
				return err
			}
		}
		logging.PrintLog("Successfully uninstalled Keptn", logging.InfoLevel)
		logging.PrintLog("\nPlease review the following namespaces and perform manual deletion if necessary:",
			logging.InfoLevel)

		namespaces, err := listAllNamespaces()
		if err != nil {
			return fmt.Errorf("Error when listing all namespaces: %v", err)
		}

		for _, namespace := range namespaces {
			logging.PrintLog(" - "+namespace, logging.InfoLevel)
			if namespace == "default" || strings.HasPrefix(namespace, "kube") || strings.HasPrefix(namespace, "openshift") {
				logging.PrintLog("      Recommended action: None (default namespace)", logging.InfoLevel)
			} else {
				// just delete the namespace
				logging.PrintLog(fmt.Sprintf("      Please review this namespace using 'kubectl get pods -n %s' before deleting it", namespace), logging.InfoLevel)
				logging.PrintLog(fmt.Sprintf("      Recommended action: kubectl delete namespace %s", namespace), logging.InfoLevel)
			}
		}

		return nil
	},
}

// returns a list of all namespaces
func listAllNamespaces() ([]string, error) {
	o := options{"get", "namespaces", "-o=jsonpath={.items[*].metadata.name}"}
	o.appendIfNotEmpty(kubectlOptions)
	out, err := keptnutils.ExecuteCommand("kubectl", o)

	if err != nil {
		return nil, err
	}

	out = strings.TrimSpace(out)
	// split by spaces
	arr := strings.Split(out, " ")
	if out != "" {
		logging.PrintLog(out, logging.VerboseLevel)
	}
	return arr, nil
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

	uninstallParams = uninstallCmdParams{}

	uninstallParams.Namespace = uninstallCmd.Flags().StringP("namespace", "n", "keptn",
		"Specify the namespace Keptn should be installed in (default keptn).")

	uninstallCmd.PersistentFlags().BoolVarP(&insecureSkipTLSVerify, "insecure-skip-tls-verify", "s", false, "Skip tls verification for kubectl commands")
}
