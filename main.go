package main

import (
	"flag"
	"os"
	"os/signal"
)

var _version = "1.0.0"

func main() {
	config := config.Config{}
	flag.StringVar(&config.DeviceName, "dn", "", "device name")
	flag.StringVar(&config.CIDR, "c", "172.18.0.15/24", "tun vpn interface")
	flag.StringVar(&config.CIDRv6, "c6", "fced:9999::9999/64", "tun vpn interface ipv6")
	flag.IntVar(&config.MTU, "mtu", 1500, "tun vpn interface mtu")
	flag.StringVar(&config.LocalAddr, "l", ":3050", "local address")
	flag.StringVar(&config.ServerAddr, "s", ":3051", "server address")
	flag.StringVar(&config.ServerIP, "sip", "172.18.0.1", "server ip")
	flag.StringVar(&config.ServerIPv6, "sip6", "fced:9999::1", "server ipv6")
	flag.StringVar(&config.DNSIP, "dip", "8.8.8.8", "dns server ip")
	flag.Parse()
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.StopApp()
}
