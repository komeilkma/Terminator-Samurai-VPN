package ipfconfig

import (
	"errors"
	"net"
	"runtime"
)

var (
	errNoGateway      = errors.New("no gateway found")
	errCantParse      = errors.New("can't parse string output")
	errNotImplemented = errors.New("not implemented for OS: " + runtime.GOOS)
)
func DiscoverGatewayIPv4() (ip net.IP, err error) {
	return discoverGatewayOSSpecificIPv4()
}
func DiscoverGatewayIPv6() (ip net.IP, err error) {
	return discoverGatewayOSSpecificIPv6()
}
