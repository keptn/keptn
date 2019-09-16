package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/cli/utils"
	"github.com/spf13/cobra"
)

var uninstallVersion *string

const istioFolder = "/installer/manifests/istio/"

var istioFiles = [...]string{"crd-10.yaml", "crd-11.yaml", "crd-12.yaml",
	"crd-certmanager-10.yaml", "crd-certmanager-11.yaml"}

const tillerPath = "/installer/manifests/tiller/tiller.yaml"

// domainCmd represents the domain command
var uninstallCmd = &cobra.Command{
	Use:          "uninstall",
	Short:        "Uninstalls keptn",
	SilenceUsage: true,

	PreRunE: func(cmd *cobra.Command, args []string) error {

		if insecureSkipTLSVerify {
			kubectlOptions = "--insecure-skip-tls-verify=true"
		}

		resourcesAvailable, err := checkUninstallResourceAvailability()
		if err != nil || !resourcesAvailable {
			return errors.New("Resources not found")
		}

		return nil

	},
	RunE: func(cmd *cobra.Command, args []string) error {

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

		utils.PrintLog("Starting to uninstall keptn", utils.InfoLevel)

		if !mocking {
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

			// Clean up istio CRDs
			for _, val := range istioFiles {
				o := options{"delete", "-f", getIstioCRD(val), "--ignore-not-found"}
				o.appendIfNotEmpty(kubectlOptions)
				out, err := keptnutils.ExecuteCommand("kubectl", o)
				out = strings.TrimSpace(out)
				if out != "" {
					utils.PrintLog(out, utils.VerboseLevel)
				}
				if err != nil {
					return err
				}
			}

			// Clean up tiller
			o := options{"delete", "-f", getTillerResource(), "--ignore-not-found"}
			o.appendIfNotEmpty(kubectlOptions)
			out, err := keptnutils.ExecuteCommand("kubectl", o)
			out = strings.TrimSpace(out)
			if out != "" {
				utils.PrintLog(out, utils.VerboseLevel)
			}
			if err != nil {
				return err
			}
		}
		utils.PrintLog("Successfully uninstalled keptn", utils.InfoLevel)

		return nil
	},
}

func deleteResources(namespace string) error {
	o := options{"delete", "services,deployments,pods,secrets,configmaps", "--all", "-n", namespace, "--ignore-not-found"}
	o.appendIfNotEmpty(kubectlOptions)
	out, err := keptnutils.ExecuteCommand("kubectl", o)
	out = strings.TrimSpace(out)
	if out != "" {
		utils.PrintLog(out, utils.VerboseLevel)
	}
	return err
}

func deleteNamespace(namespace string) error {
	o := options{"delete", "namespace", namespace, "--ignore-not-found"}
	o.appendIfNotEmpty(kubectlOptions)
	out, err := keptnutils.ExecuteCommand("kubectl", o)
	out = strings.TrimSpace(out)
	if out != "" {
		utils.PrintLog(out, utils.VerboseLevel)
	}
	return err
}

func getIstioCRD(fileName string) string {
	return installerPrefixURL + *uninstallVersion + istioFolder + fileName
}

func getTillerResource() string {
	return installerPrefixURL + *uninstallVersion + tillerPath
}

func checkUninstallResourceAvailability() (bool, error) {

	for _, val := range istioFiles {
		resp, err := http.Get(getIstioCRD(val))
		if err != nil {
			return false, err
		}
		if resp.StatusCode != http.StatusOK {
			return false, nil
		}
	}

	resp, err := http.Get(getTillerResource())
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallVersion = uninstallCmd.Flags().StringP("keptn-version", "k", "master",
		"The branch or tag of the version which is used for updating the domain")
	uninstallCmd.Flags().MarkHidden("keptn-version")
	uninstallCmd.PersistentFlags().BoolVarP(&insecureSkipTLSVerify, "insecure-skip-tls-verify", "s", false, "Skip tls verification for kubectl commands")
}
