package lib

import (
	"net"
	"net/url"

	logger "github.com/sirupsen/logrus"
)

type IIPResolver interface {
	ResolveIPAdresses(curlURL string) []string
}

type LookupFunc func(host string) ([]net.IP, error)

type IPResolver struct {
	LookupIP LookupFunc
}

func NewIPResolver(lookUpIPFunc ...LookupFunc) IPResolver {
	return IPResolver{
		LookupIP: net.LookupIP,
	}
}

func (i IPResolver) ResolveIPAdresses(curlURL string) []string {
	ipAddresses := make([]string, 0)
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
