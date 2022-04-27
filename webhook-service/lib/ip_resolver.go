package lib

import (
	"net"
	"net/url"

	logger "github.com/sirupsen/logrus"
)

type LookupFunc func(host string) ([]net.IP, error)

type IPResolver struct {
	LookupIP LookupFunc
}

func NewIPResolver(lookUpIPFunc ...LookupFunc) IPResolver {
	resolver := IPResolver{
		LookupIP: net.LookupIP,
	}
	if len(lookUpIPFunc) > 0 {
		resolver.LookupIP = lookUpIPFunc[0]
	}
	return resolver
}

func (i IPResolver) ResolveIPAdresses(curlURL string) []string {
	ipAddresses := GetDeniedURLs(GetEnv())
	parsedURL, err := url.Parse(curlURL)
	if err != nil {
		logger.Errorf("Unable to parse URL: %s", curlURL)
		return ipAddresses
	}
	ips, err := i.LookupIP(parsedURL.Hostname())
	if err != nil {
		logger.Errorf("Unable to look up IP for URL: %s", curlURL)
		return ipAddresses
	}
	for _, ip := range ips {
		ipAddresses = append(ipAddresses, ip.String())
	}
	return ipAddresses
}

func GetDeniedURLs(env map[string]string) []string {
	kubeAPIHostIP := env[KubernetesSvcHostEnvVar]
	kubeAPIPort := env[KubernetesAPIPortEnvVar]

	urls := make([]string, 0)
	if kubeAPIHostIP != "" {
		urls = append(urls, kubeAPIHostIP)
	}
	if kubeAPIPort != "" {
		urls = append(urls, "kubernetes"+":"+kubeAPIPort)
		urls = append(urls, "kubernetes.default"+":"+kubeAPIPort)
		urls = append(urls, "kubernetes.default.svc"+":"+kubeAPIPort)
		urls = append(urls, "kubernetes.default.svc.cluster.local"+":"+kubeAPIPort)
	}
	if kubeAPIHostIP != "" && kubeAPIPort != "" {
		urls = append(urls, kubeAPIHostIP+":"+kubeAPIPort)
	}
	return urls
}
