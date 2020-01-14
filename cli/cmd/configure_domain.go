package cmd

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils"
	"github.com/spf13/cobra"
)

type configureDomainCmdParams struct {
	ConfigVersion *string
	PlatformID    *string
	Branch        string
}

var configureDomainParams *configureDomainCmdParams

const installerPrefixURL = "https://raw.githubusercontent.com/keptn/keptn/"

const apiVirtualServiceSuffix = "/installer/manifests/keptn/keptn-api-virtualservice.yaml"
const apiIngressSuffix = "/installer/manifests/keptn/api-ingress.yaml"
const domainConfigMapSuffix = "/installer/manifests/keptn/keptn-domain-configmap.yaml"
const uniformServicesSuffix = "/installer/manifests/keptn/uniform-services.yaml"
const gatewaySuffix = "/installer/manifests/keptn/keptn-gateway.yaml"

// domainCmd represents the domain command
var domainCmd = &cobra.Command{
	Use:   "domain MY.DOMAIN.COM",
	Short: "Configures the domain",
	Long: `
	
Example:
	keptn configure domain my.domain.com`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {

		if len(args) != 1 {
			cmd.SilenceUsage = false
			return errors.New("Requires a domain as argument")
		}

		if _, err := url.Parse(args[0]); err != nil {
			cmd.SilenceUsage = false
			return errors.New("Cannot parse provided domain")
		}

		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {

		if insecureSkipTLSVerify {
			kubectlOptions = "--insecure-skip-tls-verify=true"
		}

		if (configureDomainParams.ConfigVersion == nil || *configureDomainParams.ConfigVersion == "") &&
			utils.IsOfficialKeptnVersion(Version) {
			configureDomainParams.ConfigVersion = &Version
		} else if configureDomainParams.ConfigVersion == nil || *configureDomainParams.ConfigVersion == "" {
			dev := "develop"
			configureDomainParams.ConfigVersion = &dev
		}

		resourcesAvailable, err := checkConfigureDomainResourceAvailability()
		if err != nil || !resourcesAvailable {
			return errors.New("Resources not found under:\n" +
				getAPIVirtualServiceURL() + "\n" +
				getDomainConfigMapURL() + "\n" +
				getUniformServicesURL())
		}
		logging.PrintLog(fmt.Sprintf("Used version for manifests: %s",
			*configureDomainParams.ConfigVersion), logging.InfoLevel)

		kubernetesPlatform := newKubernetesPlatform()
		return kubernetesPlatform.checkRequirements()

	},
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx, _ := getKubeContext()
		fmt.Println("Your kubernetes current context is configured to cluster: " + strings.TrimSpace(ctx))
		fmt.Println("Would you like to update the keptn domain for this cluster? (y/n)")

		reader := bufio.NewReader(os.Stdin)
		in, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		in = strings.TrimSpace(in)
		if in != "y" && in != "yes" {
			fmt.Println("Please first configure your kubernetes current context so that it" +
				"points to the cluster where you would like to update the keptn domain.")
			return nil
		}

		fmt.Println("Please note that the domain of already onboarded services is not updated!")

		logging.PrintLog("Starting to configure domain", logging.InfoLevel)

		path, err := keptnutils.GetKeptnDirectory()
		if err != nil {
			return err
		}

		ingress, err := getIngressType()
		if err != nil {
			return err
		}

		if !mocking {

			// Generate new certificate
			if err := updateCertificate(path, args[0], ingress); err != nil {
				return err
			}

			if ingress == Istio {
				if err := updateKeptnAPIVirtualService(path, args[0]); err != nil {
					return err
				}
				// Re-deploy gateway, ignore if not found
				if err := reDeployGateway(); err != nil {
					return err
				}
			} else if ingress == Nginx {
				if err := updateKeptnAPIIngress(path, args[0]); err != nil {
					return err
				}
			}

			if err := updateKeptnDomainConfigMap(path, args[0]); err != nil {
				return err
			}

			if err := keptnutils.RestartPodsWithSelector(false, "keptn", "run=api"); err != nil {
				return err
			}

			if err := keptnutils.WaitForPodsWithSelector(false, "keptn", "run=api", 5, 5*time.Second); err != nil {
				return err
			}

			if strings.ToLower(*configureDomainParams.PlatformID) == openshift {
				logging.PrintLog("Successfully configured domain", logging.InfoLevel)
				fmt.Println("Please manually execute the following commands for deleting an old route and creating a new route:")
				fmt.Println("oc delete route istio-wildcard-ingress-secure-keptn -n istio-system")
				fmt.Println("oc create route passthrough istio-wildcard-ingress-secure-keptn --service=istio-ingressgateway --hostname=\"www.keptn.ingress-gateway. " +
					args[0] + "\" --port=https --wildcard-policy=Subdomain --insecure-policy='None' -n istio-system")
				fmt.Println()
				token, err := getAPITokenUsingKube()
				if err != nil {
					return err
				}
				fmt.Println("Afterwards, you can login with 'keptn auth --endpoint=https://api.keptn." + args[0] + " --token=" + token + "'")

			} else {
				var err error
				for retries := 0; retries < 3; retries++ {
					if err = authUsingKube(); err == nil {
						break
					}
					logging.PrintLog("Retry authentication...", logging.InfoLevel)
					if err := keptnutils.RestartPodsWithSelector(false, "keptn", "run=api"); err != nil {
						return err
					}
					if err := keptnutils.WaitForPodsWithSelector(false, "keptn", "run=api", 5, 5*time.Second); err != nil {
						return err
					}
				}
				if err != nil {
					logging.PrintLog("Cannot authenticate to api", logging.QuietLevel)
					return err
				}
				logging.PrintLog("Successfully configured domain", logging.InfoLevel)
			}
		}

		return nil
	},
}

func getIngressType() (Ingress, error) {

	o := options{"get", "ns"}
	o.appendIfNotEmpty(kubectlOptions)
	namespaces, err := keptnutils.ExecuteCommand("kubectl", o)
	if err != nil {
		return Istio, err
	}
	if strings.Contains(namespaces, "istio-system") {
		return Istio, nil
	} else if strings.Contains(namespaces, "ingress-nginx") {
		return Nginx, nil
	}
	return Istio, errors.New("Cannot obtain type of ingress.")
}

func reDeployGateway() error {
	o := options{"delete", "-f", getGatewayURL(), "--ignore-not-found"}
	o.appendIfNotEmpty(kubectlOptions)
	_, err := keptnutils.ExecuteCommand("kubectl", o)
	if err != nil {
		return err
	}

	o = options{"apply", "-f", getGatewayURL()}
	o.appendIfNotEmpty(kubectlOptions)
	_, err = keptnutils.ExecuteCommand("kubectl", o)
	return err
}

func updateKeptnDomainConfigMap(path, domain string) error {

	keptnDomainConfigMap := path + "keptn-domain-configmap.yaml"

	if err := utils.DownloadFile(keptnDomainConfigMap, getDomainConfigMapURL()); err != nil {
		return err
	}

	if err := utils.Replace(keptnDomainConfigMap,
		utils.PlaceholderReplacement{PlaceholderValue: "DOMAIN_PLACEHOLDER", DesiredValue: domain}); err != nil {
		return err
	}

	o := options{"delete", "-f", keptnDomainConfigMap}
	o.appendIfNotEmpty(kubectlOptions)
	_, err := keptnutils.ExecuteCommand("kubectl", o)
	if err != nil {
		return err
	}

	// Add config map in keptn namespace that contains the domain - this will be used by other services as well
	o = options{"apply", "-f", keptnDomainConfigMap}
	o.appendIfNotEmpty(kubectlOptions)
	_, err = keptnutils.ExecuteCommand("kubectl", o)
	return err
}

func updateKeptnAPIVirtualService(path, domain string) error {

	keptnAPIVSPath := path + "keptn-api-virtualservice.yaml"

	if err := utils.DownloadFile(keptnAPIVSPath, getAPIVirtualServiceURL()); err != nil {
		return err
	}

	if err := utils.Replace(keptnAPIVSPath,
		utils.PlaceholderReplacement{PlaceholderValue: "DOMAIN_PLACEHOLDER", DesiredValue: domain}); err != nil {
		return err
	}

	// Delete old api virtual service
	o := options{"delete", "-f", keptnAPIVSPath}
	o.appendIfNotEmpty(kubectlOptions)
	_, err := keptnutils.ExecuteCommand("kubectl", o)
	if err != nil {
		return err
	}

	// Apply new api virtual service
	o = options{"apply", "-f", keptnAPIVSPath}
	o.appendIfNotEmpty(kubectlOptions)
	_, err = keptnutils.ExecuteCommand("kubectl", o)
	return err
}

func updateKeptnAPIIngress(path, domain string) error {

	keptnAPIIngress := path + "api-ingress.yaml"
	if err := utils.DownloadFile(keptnAPIIngress, getAPIIngressURL()); err != nil {
		return err
	}

	if err := utils.Replace(keptnAPIIngress,
		utils.PlaceholderReplacement{PlaceholderValue: "domain.placeholder", DesiredValue: domain}); err != nil {
		return err
	}

	// Delete old api ingress
	o := options{"delete", "-f", keptnAPIIngress}
	o.appendIfNotEmpty(kubectlOptions)
	_, err := keptnutils.ExecuteCommand("kubectl", o)
	if err != nil {
		return err
	}

	// Apply new api virtual service
	o = options{"apply", "-f", keptnAPIIngress}
	o.appendIfNotEmpty(kubectlOptions)
	_, err = keptnutils.ExecuteCommand("kubectl", o)
	return err
}

func updateCertificate(path, domain string, ingress Ingress) error {

	// Source: https://golang.org/src/crypto/tls/generate_cert.go
	// We can verify the generated key with 'openssl rsa -in key.pem -check'
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Keptn"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(1, 0, 0),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{domain},
	}

	// generate private key
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	publickey := &privatekey.PublicKey

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publickey, privatekey)
	if err != nil {
		return err
	}

	privateKeyPath := path + "private.key"
	certPath := path + "cert.pem"

	certOut, err := os.Create(certPath)
	if err != nil {
		return err
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return err
	}
	if err := certOut.Close(); err != nil {
		return err
	}
	defer os.Remove(certPath)

	keyOut, err := os.OpenFile(privateKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privatekey)}); err != nil {
		return err
	}
	if err := keyOut.Close(); err != nil {
		return err
	}
	defer os.Remove(privateKeyPath)

	// First delete secret and afterwards apply new secret with new certificate
	if ingress == Istio {
		o := options{"delete", "--namespace", "istio-system", "secret", "istio-ingressgateway-certs"}
		o.appendIfNotEmpty(kubectlOptions)
		keptnutils.ExecuteCommand("kubectl", o)

		o = options{"create", "--namespace", "istio-system", "secret", "tls", "istio-ingressgateway-certs",
			"--key", privateKeyPath, "--cert", certPath}
		o.appendIfNotEmpty(kubectlOptions)
		_, err = keptnutils.ExecuteCommand("kubectl", o)
		return err
	}
	// Reset secret for nginx
	o := options{"delete", "secret", "sslcerts", "--namespace", "keptn"}
	o.appendIfNotEmpty(kubectlOptions)
	keptnutils.ExecuteCommand("kubectl", o)

	o = options{"create", "secret", "tls", "sslcerts",
		"--key", privateKeyPath, "--cert", certPath, "--namespace", "keptn"}
	o.appendIfNotEmpty(kubectlOptions)
	_, err = keptnutils.ExecuteCommand("kubectl", o)
	return err
}

func getAPIVirtualServiceURL() string {
	return installerPrefixURL + *configureDomainParams.ConfigVersion + apiVirtualServiceSuffix
}

func getAPIIngressURL() string {
	return installerPrefixURL + *configureDomainParams.ConfigVersion + apiIngressSuffix
}

func getDomainConfigMapURL() string {
	return installerPrefixURL + *configureDomainParams.ConfigVersion + domainConfigMapSuffix
}

func getUniformServicesURL() string {
	return installerPrefixURL + *configureDomainParams.ConfigVersion + uniformServicesSuffix
}

func getGatewayURL() string {
	return installerPrefixURL + *configureDomainParams.ConfigVersion + gatewaySuffix
}

func checkConfigureDomainResourceAvailability() (bool, error) {

	resp, err := http.Get(getAPIVirtualServiceURL())
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}
	resp, err = http.Get(getAPIIngressURL())
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	resp, err = http.Get(getDomainConfigMapURL())
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	resp, err = http.Get(getUniformServicesURL())
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	resp, err = http.Get(getGatewayURL())
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}
	return true, nil
}

func init() {
	configureCmd.AddCommand(domainCmd)
	configureDomainParams = &configureDomainCmdParams{}

	configureDomainParams.ConfigVersion = domainCmd.Flags().StringP("keptn-version", "k", "",
		"The branch or tag containing the manifests which are used for updating the domain")
	domainCmd.Flags().MarkHidden("keptn-version")
	domainCmd.PersistentFlags().BoolVarP(&insecureSkipTLSVerify, "insecure-skip-tls-verify", "s", false, "Skip tls verification for kubectl commands")

	configureDomainParams.PlatformID = domainCmd.Flags().StringP("platform", "p", "gke", "The platform on which keptn is running [gke,openshift,aks,kubernetes]")
}
