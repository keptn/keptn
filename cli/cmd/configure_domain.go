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
	"github.com/keptn/keptn/cli/utils"
	"github.com/spf13/cobra"
)

var configVersion *string

const apiVirtualServiceURL = "/installer/manifests/keptn/keptn-api-virtualservice.yaml"
const domainConfigMapURL = "/installer/manifests/keptn/keptn-domain-configmap.yaml"
const uniformServicesURL = "/installer/manifests/keptn/uniform-services.yaml"

// domainCmd represents the domain command
var domainCmd = &cobra.Command{
	Use:          "domain domain_url",
	Short:        "Configures the domain",
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {

		if len(args) != 1 {
			cmd.SilenceUsage = false
			return errors.New("Requires a domain as argument")
		}

		_, err := url.Parse(args[0])
		if err != nil {
			cmd.SilenceUsage = false
			return errors.New("Cannot parse provided domain")
		}

		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {

		if insecureSkipTLSVerify {
			kubectlOptions = "--insecure-skip-tls-verify=true"
		}

		resourcesAvailable, err := checkResourceAvailablity()
		if err != nil || !resourcesAvailable {
			return errors.New("Resources not found under:\n" +
				getAPIVirtualServiceURL() + "\n" +
				getDomainConfigMapURL() + "\n" +
				getUniformServicesURL())
		}

		kubernetesPlatform := newKubernetesPlatform()
		err = kubernetesPlatform.checkRequirements()
		if err != nil {
			return err
		}

		ctx, _ := getKubeContext()
		fmt.Println("Your .kube-config is configured to: " + strings.TrimSpace(ctx))
		fmt.Println("Would you like to update the keptn domain for this cluster? (y/n)")

		reader := bufio.NewReader(os.Stdin)
		in, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		in = strings.TrimSpace(in)
		if in != "y" && in != "yes" {
			return errors.New(`Please first configure your .kube-config so that it
points to the cluster where you would like to update the keptn domain`)
		}

		fmt.Println("Please note that the domain of already onboarded services is not updated!")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		utils.PrintLog("Starting to configure domain", utils.InfoLevel)

		path, err := keptnutils.GetKeptnDirectory()
		if err != nil {
			return err
		}

		if !mocking {

			err = updateKeptnAPIVirtualService(path, args[0])
			if err != nil {
				return err
			}

			// Generate new certificate
			err = updateCertificate(path, args[0])
			if err != nil {
				return err
			}

			err = reDeployGithubService()
			if err != nil {
				return err
			}

			err = authUsingKube()
			if err != nil {
				return err
			}

			configured, err := checkIfConfiguredUsingKube()
			if err != nil {
				return err
			}
			if !configured {
				fmt.Println("No GitHub configuration found on your keptn installation. Please exectue 'keptn configure'")
			}
		}

		return nil
	},
}

func reDeployGithubService() error {

	o := options{"delete", "deployment", "github-service", "-n", "keptn"}
	o.appendIfNotEmpty(kubectlOptions)
	_, err := keptnutils.ExecuteCommand("kubectl", o)
	if err != nil {
		return err
	}

	o = options{"apply", "-f", getUniformServicesURL()}
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

	// Add config map in keptn namespace that contains the domain - this will be used by other services as well
	options := options{"apply", "-f", keptnDomainConfigMap}
	options.appendIfNotEmpty(kubectlOptions)
	_, err := keptnutils.ExecuteCommand("kubectl", options)
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

func updateCertificate(path, domain string) error {

	template := &x509.Certificate{
		IsCA:                  true,
		BasicConstraintsValid: true,
		SubjectKeyId:          []byte(domain),
		SerialNumber:          big.NewInt(1234),
		Subject: pkix.Name{
			Country:      []string{"Austria"},
			Organization: []string{"keptn"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(1, 0, 0),
		// see http://golang.org/pkg/crypto/x509/#KeyUsage
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	// generate private key
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	publickey := &privatekey.PublicKey

	// create a self-signed certificate. template = parent
	var parent = template
	cert, err := x509.CreateCertificate(rand.Reader, template, parent, publickey, privatekey)
	if err != nil {
		return err
	}

	privateKeyPath := path + "private.key"
	keyfile, _ := os.Create(privateKeyPath)
	var pemkey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privatekey)}
	pem.Encode(keyfile, pemkey)
	keyfile.Close()
	defer os.Remove(privateKeyPath)

	certPath := path + "cert.pem"
	pemfile, _ := os.Create(certPath)
	var pemCert = &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert}
	pem.Encode(pemfile, pemCert)
	pemfile.Close()
	defer os.Remove(certPath)

	// First delete secret
	o := options{"delete", "--namespace", "istio-system", "secret", "istio-ingressgateway-certs"}
	o.appendIfNotEmpty(kubectlOptions)
	_, err = keptnutils.ExecuteCommand("kubectl", o)

	o = options{"create", "--namespace", "istio-system", "secret", "tls", "istio-ingressgateway-certs",
		"--key", privateKeyPath, "--cert", certPath}
	o.appendIfNotEmpty(kubectlOptions)
	_, err = keptnutils.ExecuteCommand("kubectl", o)
	return err
}

func getAPIVirtualServiceURL() string {
	return installerPrefixURL + *configVersion + apiVirtualServiceURL
}

func getDomainConfigMapURL() string {
	return installerPrefixURL + *configVersion + domainConfigMapURL
}

func getUniformServicesURL() string {
	return installerPrefixURL + *configVersion + uniformServicesURL
}

func checkResourceAvailablity() (bool, error) {

	resp, err := http.Get(getAPIVirtualServiceURL())
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
	return true, nil
}

func init() {
	configureCmd.AddCommand(domainCmd)
	configVersion = domainCmd.Flags().StringP("keptn-version", "k", "master",
		"The branch or tag of the version which is used for updating the domain")
	domainCmd.Flags().MarkHidden("keptn-version")
	domainCmd.PersistentFlags().BoolVarP(&insecureSkipTLSVerify, "insecure-skip-tls-verify", "s", false, "Skip tls verification for kubectl commands")
}
