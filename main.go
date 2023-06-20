package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"github.com/komeilkma/Terminator-Samurai-VPN/finalapp"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/config"
)

var _version = "1.0.0"

func main() {
	config := config.Config{}
	flag.StringVar(&config.DeviceName, "dn", "", "device name")
	flag.StringVar(&config.CIDR, "c", "172.18.0.10/24", "tun interface cidr")
	flag.StringVar(&config.CIDRv6, "c6", "fced:9999::9999/64", "tun interface ipv6 cidr")
	flag.IntVar(&config.MTU, "mtu", 1500, "tun mtu")
	flag.StringVar(&config.LocalAddr, "l", ":3000", "local address")
	flag.StringVar(&config.ServerAddr, "s", ":3001", "server address")
	flag.StringVar(&config.ServerIP, "sip", "172.18.0.1", "server ip")
	flag.StringVar(&config.ServerIPv6, "sip6", "fced:9999::1", "server ipv6")
	flag.StringVar(&config.DNSIP, "dip", "8.8.8.8", "dns server ip")
	flag.StringVar(&config.Key, "k", "12312300kma", "key")
	flag.StringVar(&config.Protocol, "p", "udp", "protocol udp/tls/grpc/quic/ws/wss")
	flag.StringVar(&config.WebSocketPath, "path", "/freedom", "websocket path")
	flag.BoolVar(&config.ServerMode, "S", false, "server mode")
	flag.BoolVar(&config.GlobalMode, "g", false, "client global mode")
	flag.BoolVar(&config.Obfs, "obfs", false, "enable data obfuscation")
	flag.BoolVar(&config.Compress, "compress", false, "enable data compression")
	flag.IntVar(&config.Timeout, "t", 30, "dial timeout in seconds")
	flag.StringVar(&config.TLSSni, "sni", "", "tls handshake sni")
	flag.BoolVar(&config.TLSInsecureSkipVerify, "isv", false, "tls insecure skip verify")
	flag.BoolVar(&config.Verbose, "v", false, "enable verbose output")
	flag.Parse()
	app := app.NewApp(&config, _version)
	app.InitConfig()
	go app.StartApp()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.StopApp()
}
