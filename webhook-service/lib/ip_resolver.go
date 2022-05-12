package lib

import (
	logger "github.com/sirupsen/logrus"
	"net"
	neturl "net/url"
)

type IPResolver interface {
	Resolve(url string) AdrDomainNameMapping
}

type LookupFunc func(host string) ([]net.IP, error)
type LookupAddrFunc func(addr string) (names []string, err error)
type ParseFunc func(rawURL string) (*neturl.URL, error)

type ipResolver struct {
	lookupIP   LookupFunc
	lookupAddr LookupAddrFunc
	parse      ParseFunc
}

func NewIPResolver() IPResolver {
	return ipResolver{
		lookupIP:   net.LookupIP,
		lookupAddr: net.LookupAddr,
		parse:      neturl.Parse,
	}
}

func (i ipResolver) Resolve(url string) AdrDomainNameMapping {
	ipAddresses := make(AdrDomainNameMapping, 0)
	parsedURL, err := i.parse(url)
	if err != nil {
		logger.Errorf("Unable to parse URL: %s", url)
		return ipAddresses
	}

	ips, err := i.lookupIP(parsedURL.Hostname())
	if err != nil {
		logger.Errorf("Unable to look up IP for URL: %s", url)
		return ipAddresses
	}
	for _, ip := range ips {

		// for each ip get all its domains to check if they are among the denied
		hosts, err := i.lookupAddr(ip.String())
		if err != nil {
			logger.Errorf("Unable to look up domains for URL: %s", url)
		}
		ipAddresses[ip.String()] = hosts
	}

	return ipAddresses
}
