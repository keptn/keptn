package lib

import (
	"net"
	neturl "net/url"

	logger "github.com/sirupsen/logrus"
)

type IPResolver interface {
	Resolve(url string) []string
}

type LookupFunc func(host string) ([]net.IP, error)
type ParseFunc func(rawURL string) (*neturl.URL, error)

type ipResolver struct {
	LookupIP LookupFunc
	Parse    ParseFunc
}

func NewIPResolver() IPResolver {
	return ipResolver{
		LookupIP: net.LookupIP,
		Parse:    neturl.Parse,
	}
}

func (i ipResolver) Resolve(url string) []string {
	ipAddresses := make([]string, 0)
	parsedURL, err := i.Parse(url)
	if err != nil {
		logger.Errorf("Unable to parse URL: %s", url)
		return ipAddresses
	}
	ips, err := i.LookupIP(parsedURL.Hostname())
	if err != nil {
		logger.Errorf("Unable to look up IP for URL: %s", url)
		return ipAddresses
	}
	for _, ip := range ips {
		ipAddresses = append(ipAddresses, ip.String())
	}
	return ipAddresses
}
