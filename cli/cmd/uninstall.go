package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/keptn/keptn/cli/pkg/logging"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/spf13/cobra"
)

var uninstallVersion *string

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:          "uninstall",
	Short:        "Uninstalls Keptn from a Kubernetes cluster",
	Long: `Uninstalls Keptn from a Kubernetes cluster.

This command does *not* delete: 

* Istio
* Tiller 
* Dynatrace monitoring
* Prometheus monitoring

Besides, deployed services and the configuration on the Git upstream (i.e., GitHub, GitLab, or Bitbucket) are not deleted. To clean-up created projects and services, instructions are provided [here](../../manage/project#delete-a-project).

**Note:** This command requires a *kubernetes current context* pointing to the cluster where Keptn should get uninstalled.
`,
	Example: `keptn uninstall`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		if insecureSkipTLSVerify {
			kubectlOptions = "--insecure-skip-tls-verify=true"
		}

		ctx, _ := getKubeContext()
		fmt.Println("Your kubernetes current context is configured to cluster: " + strings.TrimSpace(ctx))
		fmt.Println("Would you like to uninstall Keptn from this cluster? (y/n)")

		reader := bufio.NewReader(os.Stdin)
		in, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		in = strings.TrimSpace(in)
		if in != "y" && in != "yes" {
			return nil
		}

		logging.PrintLog("Starting to uninstall Keptn", logging.InfoLevel)

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
		logging.PrintLog("Successfully uninstalled Keptn", logging.InfoLevel)
		logging.PrintLog("Note: Please review the following namespaces and perform manual deletion if necessary:",
			logging.InfoLevel)

		namespaces, err := listAllNamespaces()
		if err != nil {
			return fmt.Errorf("Error when listing all namespaces: %v", err)
		}

		for _, namespace := range namespaces {
			logging.PrintLog(" - "+namespace, logging.InfoLevel)
			if namespace == "default" || namespace == "kube-public" {
				// skip
				logging.PrintLog("      Recommended action: None (default namespace)", logging.InfoLevel)
			} else if namespace == "kube-system" {
				// we need to remove helm / tiller stuff
				logging.PrintLog("      Recommended action: Remove Tiller/Helm using", logging.InfoLevel)
				logging.PrintLog("                          kubectl delete all -l app=helm -n kube-system", logging.InfoLevel)
			} else if namespace == "istio-system" {
				// istio is special, we will refer to the official uninstall docs
				logging.PrintLog("      Please consult the istio Docs at https://istio.io/docs/setup/install/helm/#uninstall on how to remove istio.", logging.InfoLevel)
				logging.PrintLog("      Recommended action: kubectl delete namespace istio-system", logging.InfoLevel)
			} else {
				// just delete the namespace
				logging.PrintLog(fmt.Sprintf("      Please review this namespace in detail using 'kubectl get pods -n %s' before deleting it", namespace), logging.InfoLevel)
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
