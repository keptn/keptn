package probes

import (
	"net"
	"net/url"
	"time"
)

const (
	DefaultPort      = "80"
	DefaultTransport = "tcp"
)

var DefaultTimeout = time.Duration(3) * time.Second

type Dialer func(network, address string, timeout time.Duration) (net.Conn, error)

type Config struct {
	URL     *url.URL
	Dialer  Dialer
	Timeout time.Duration
	Network string
}

type ReachabilityChecker struct {
	dialer  Dialer
	timeout time.Duration
	network string
	url     *url.URL
	tags    []string
}

func NewReachabilityChecker(cfg *Config) (*ReachabilityChecker, error) {
	t := DefaultTimeout
	if cfg.Timeout != 0 {
		t = cfg.Timeout
	}
	d := net.DialTimeout
	if cfg.Dialer != nil {
		d = cfg.Dialer
	}
	n := DefaultTransport
	if cfg.Network != "" {
		n = cfg.Network
	}
	r := &ReachabilityChecker{
		dialer:  d,
		timeout: t,
		network: n,
		url:     cfg.URL,
	}
	return r, nil
}

func (r *ReachabilityChecker) Status() error {
	port := r.url.Port()
	if len(port) == 0 {
		port = DefaultPort
	}

	conn, err := r.dialer(r.network, r.url.Hostname()+":"+port, r.timeout)
	if err != nil {
		return err
	}
	if conn != nil {
		if errClose := conn.Close(); errClose != nil {
			return errClose
		}
	}
	return nil
}
