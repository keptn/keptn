package lib

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"

	keptnkubeutils "github.com/keptn/kubernetes-utils/pkg"
	logger "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ICurlValidator interface {
	Validate(request Request, denyList []string, ipAddresses []string) error
	ResolveIPAdresses(curlURL string) []string
	GetConfigDenyList() []string
}

type CurlValidator struct {
	deniedURLs    []string
	curlValidator ICurlValidator
}

func NewCurlValidator(curlValidator ICurlValidator, deniedURLs []string) *CurlValidator {
	validator := &CurlValidator{
		deniedURLs:    deniedURLs,
		curlValidator: curlValidator,
	}
	return validator
}

func (c *CurlValidator) Validate(request Request, denyList []string, ipAddresses []string) error {
	if request.URL == "" {
		return fmt.Errorf("Invalid curl URL: '%s'", request.URL)
	}

	for _, url := range denyList {
		if strings.Contains(request.URL, url) {
			return fmt.Errorf("curl command contains denied URL '%s'", url)
		}
		for _, ip := range ipAddresses {
			if strings.Contains(ip, url) {
				return fmt.Errorf("curl command contains denied IP address '%s'", url)
			}
		}
	}
	return nil
}

func (c *CurlValidator) ResolveIPAdresses(curlURL string) []string {
	ipAddresses := make([]string, 0)
	parsedURL, err := url.Parse(curlURL)
	if err != nil {
		logger.Errorf("Unable to parse URL: %s", curlURL)
		return ipAddresses
	}
	ips, err := net.LookupIP(parsedURL.Hostname())
	if err != nil {
		logger.Errorf("Unable to look up IP for URL: %s", curlURL)
		return ipAddresses
	}
	for _, ip := range ips {
		ipAddresses = append(ipAddresses, ip.String())
	}
	return ipAddresses
}

func (c *CurlValidator) GetConfigDenyList() []string {
	denyList := make([]string, 0)
	kubeAPI, err := keptnkubeutils.GetKubeAPI(false)
	if err != nil {
		logger.Errorf("Unable to read ConfigMap %s: cannot get kubeAPI: %s", WebhookConfigMap, err.Error())
		return c.deniedURLs
	}

	configMap, err := kubeAPI.ConfigMaps(GetNamespaceFromEnvVar()).Get(context.TODO(), WebhookConfigMap, v1.GetOptions{})
	if err != nil {
		logger.Errorf("Unable to get ConfigMap %s: %s", WebhookConfigMap, err.Error())
		return c.deniedURLs
	}

	denyListString := configMap.Data["denyList"]
	denyList = strings.Fields(denyListString)
	return denyList
}
