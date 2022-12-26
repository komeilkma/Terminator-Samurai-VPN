//go:build windows
// +build windows

package ipfconfig

import (
	"net"
	"os/exec"
	"syscall"
)

func discoverGatewayOSSpecificIPv4() (ip net.IP, err error) {
	routeCmd := exec.Command("route", "print", "0.0.0.0")
	routeCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := routeCmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return parseWindowsGatewayIPv4(output)
}

func discoverGatewayOSSpecificIPv6() (ip net.IP, err error) {
	routeCmd := exec.Command("route", "print", "-6", "::/0")
	routeCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := routeCmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return parseWindowsGatewayIPv6(output)
}
