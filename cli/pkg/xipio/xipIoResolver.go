package xipio

import (
	"context"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/keptn/keptn/cli/pkg/logging"
)

// ResolveXipIo resolves a xip io address
func ResolveXipIo(network, addr string) (net.Conn, error) {
	return ResolveXipIoWithContext(context.Background(), network, addr)
}

// ResolveXipIo resolves a xip io address
func ResolveXipIoWithContext(ctx context.Context, network, addr string) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	if strings.Contains(addr, ".xip.io") {

		regex := `\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?).xip.io\b`
		re := regexp.MustCompile(regex)
		ipWithXipIo := re.FindString(addr)
		ip := ipWithXipIo[:len(ipWithXipIo)-len(".xip.io")]

		regex = `:\d+$`
		re = regexp.MustCompile(regex)
		port := re.FindString(addr)

		var newAddr string
		if port != "" {
			newAddr = ip + port
		}
		logging.PrintLog("Directly resolve "+addr+" to "+newAddr, logging.VerboseLevel)
		addr = newAddr
	}
	return dialer.DialContext(ctx, network, addr)
}
