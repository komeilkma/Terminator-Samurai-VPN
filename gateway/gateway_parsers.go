package gateway

import (
	"log"
	"net"
	"os/exec"
	"strings"
)

type windowsRouteStructIPv4 struct {
	Destination string
	Netmask     string
	Gateway     string
	Interface   string
	Metric      string
}

type windowsRouteStructIPv6 struct {
	If          string
	Metric      string
	Destination string
	Gateway     string
}

func parseToWindowsRouteStructIPv4(output []byte) (windowsRouteStructIPv4, error) {

	lines := strings.Split(string(output), "\n")
	sep := 0
	for idx, line := range lines {
		if sep == 3 {
			if len(lines) <= idx+2 {
				return windowsRouteStructIPv4{}, errNoGateway
			}

			fields := strings.Fields(lines[idx+2])
			if len(fields) < 5 {
				return windowsRouteStructIPv4{}, errCantParse
			}

			return windowsRouteStructIPv4{
				Destination: fields[0],
				Netmask:     fields[1],
				Gateway:     fields[2],
				Interface:   fields[3],
				Metric:      fields[4],
			}, nil
		}
		if strings.HasPrefix(line, "=======") {
			sep++
			continue
		}
	}
	return windowsRouteStructIPv4{}, errNoGateway
}

func parseToWindowsRouteStructIPv6(output []byte) (windowsRouteStructIPv6, error) {

	lines := strings.Split(string(output), "\n")
	sep := 0
	for idx, line := range lines {
		if sep == 3 {
			if len(lines) <= idx+2 {
				return windowsRouteStructIPv6{}, errNoGateway
			}

			fields := strings.Fields(lines[idx+2])
			if len(fields) < 4 {
				return windowsRouteStructIPv6{}, errCantParse
			}

			return windowsRouteStructIPv6{
				If:          fields[0],
				Metric:      fields[1],
				Destination: fields[2],
				Gateway:     fields[3],
			}, nil
		}
		if strings.HasPrefix(line, "=======") {
			sep++
			continue
		}
	}
	return windowsRouteStructIPv6{}, errNoGateway
}

func parseWindowsGatewayIPv4(output []byte) (net.IP, error) {
	parsedOutput, err := parseToWindowsRouteStructIPv4(output)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(parsedOutput.Gateway)
	if ip == nil {
		return nil, errCantParse
	}
	return ip, nil
}

func parseWindowsGatewayIPv6(output []byte) (net.IP, error) {
	parsedOutput, err := parseToWindowsRouteStructIPv6(output)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(parsedOutput.Gateway)
	if ip == nil {
		return nil, errCantParse
	}
	return ip, nil
}

func execCmd(c string, args ...string) string {
	cmd := exec.Command(c, args...)
	out, err := cmd.Output()
	if err != nil {
		log.Println("failed to exec cmd:", err)
	}
	if len(out) == 0 {
		return ""
	}
	s := string(out)
	return strings.ReplaceAll(s, "\n", "")
}
