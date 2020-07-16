package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"helm.sh/helm/v3/pkg/action"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/keptn/keptn/cli/pkg/logging"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstalls Keptn from a Kubernetes cluster",
	Long: `Uninstalls Keptn from a Kubernetes cluster.

This command does *not* delete: 

* Istio
* Dynatrace monitoring
* Prometheus monitoring
* Any (third-party) service installed in addition to Keptn (e.g., notification-service, slackbot-service, ...)

Besides, deployed services and the configuration on the Git upstream (i.e., GitHub, GitLab, or Bitbucket) are not deleted. To clean-up created projects and services, please see [Delete a project](https://keptn.sh/docs/0.7.x/manage/project/#delete-a-project).

**Note:** This command requires a *Kubernetes current context* pointing to the cluster where Keptn should get uninstalled.
`,
	Example:      `keptn uninstall`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		if insecureSkipTLSVerify {
			kubectlOptions = "--insecure-skip-tls-verify=true"
		}

		ctx, _ := getKubeContext()
		fmt.Println("Your Kubernetes current context is configured to cluster: " + strings.TrimSpace(ctx))
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
			if err := uninstallKeptnChart("keptn", "keptn"); err != nil {
				return err
			}
			// Clean up keptn namespace
			if err := deleteNamespace("keptn"); err != nil {
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

func uninstallKeptnChart(releaseName, namespace string) error {
	logging.PrintLog(fmt.Sprintf("Start deleting Helm Chart %s in namespace %s", releaseName, namespace), logging.InfoLevel)
	var kubeconfig string
	if os.Getenv("KUBECONFIG") != "" {
		kubeconfig = keptnutils.ExpandTilde(os.Getenv("KUBECONFIG"))
	} else {
		kubeconfig = filepath.Join(
			keptnutils.UserHomeDir(), ".kube", "config",
		)
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	cfg, err := newActionConfig(config, namespace)
	if err != nil {
		return err
	}

	iCli := action.NewUninstall(cfg)
	_, err = iCli.Run(releaseName)

	if err != nil {
		return fmt.Errorf("Error when deleting Helm Chart %s in namespace %s: %s",
			releaseName, namespace, err.Error())
	}
	return nil
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
	uninstallCmd.PersistentFlags().BoolVarP(&insecureSkipTLSVerify, "insecure-skip-tls-verify", "s", false, "Skip tls verification for kubectl commands")
}
