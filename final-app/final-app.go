package final-app

import (
	"github.com/komeilkma/Terminator-Samurai-VPN/kcp"
	"github.com/komeilkma/Terminator-Samurai-VPN/quic-proto"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/cipher"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/config"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/netutil"
	"github.com/komeilkma/Terminator-Samurai-VPN/grpc"
	"github.com/komeilkma/Terminator-Samurai-VPN/tls"
	"github.com/komeilkma/Terminator-Samurai-VPN/tun"
	"github.com/komeilkma/Terminator-Samurai-VPN/UserDatagramProtocol"
	"github.com/komeilkma/Terminator-Samurai-VPN/websocket"
	"github.com/komeilkma/Terminator-Samurai-VPN/native-water"
	"log"
)

var _banner = `
  _____              _           _               ___                          _ 
 |_   _|__ _ _ _ __ (_)_ _  __ _| |_ ___ _ _ ___/ __| __ _ _ __ _  _ _ _ __ _(_)
   | |/ -_) '_| '  \| | ' \/ _' |  _/ _ \ '_|___\__ \/ _' | '  \ || | '_/ _' | |
   |_|\___|_| |_|_|_|_|_||_\__,_|\__\___/_|     |___/\__,_|_|_|_\_,_|_| \__,_|_|
                   
%s
`
var _srcUrl = "https://github.com/komeilkma/Terminator-Samurai-VPN"


type App struct {
	Config  *config.Config
	Version string
	Iface   *water.Interface
}

func NewApp(config *config.Config, version string) *App {

	return &App{
		Config:  config,
		Version: version,
	}
}

func (app *App) StartApp() {

	switch app.Config.Protocol {
	case "udp":
		if app.Config.ServerMode {
			udp.StartServer(app.Iface, *app.Config)
		} else {
			udp.StartClient(app.Iface, *app.Config)
		}
	case "ws", "wss":
		if app.Config.ServerMode {
			ws.StartServer(app.Iface, *app.Config)
		} else {
			ws.StartClient(app.Iface, *app.Config)
		}
	case "tls":
		if app.Config.ServerMode {
			tls.StartServer(app.Iface, *app.Config)
		} else {
			tls.StartClient(app.Iface, *app.Config)
		}
	case "grpc":
		if app.Config.ServerMode {
			grpc.StartServer(app.Iface, *app.Config)
		} else {
			grpc.StartClient(app.Iface, *app.Config)
		}
	case "quic":
		if app.Config.ServerMode {
			quic.StartServer(app.Iface, *app.Config)
		} else {
			quic.StartClient(app.Iface, *app.Config)
		}
	case "kcp":
		if app.Config.ServerMode {
			kcp.StartServer(app.Iface, *app.Config)
		}else {
			kcp.StartClient(app.Iface, *app.Config)
		}
	default:
		if app.Config.ServerMode {
			udp.StartServer(app.Iface, *app.Config)
		} else {
			udp.StartClient(app.Iface, *app.Config)
		}
	}
}

func (app *App) InitConfig() {
	log.Printf(_banner, _srcUrl)
	log.Printf("TSVPN version %s", app.Version)
	if !app.Config.ServerMode {
		app.Config.LocalGateway = netutil.DiscoverGateway(true)
		app.Config.LocalGatewayv6 = netutil.DiscoverGateway(false)
	}
	app.Config.BufferSize = 64 * 1024
	cipher.SetKey(app.Config.Key)
	app.Iface = tun.CreateTun(*app.Config)
	log.Printf("initialized config: %+v", app.Config)
	netutil.PrintStats(app.Config.Verbose, app.Config.ServerMode)
}

func (app *App) StopApp() {
	tun.ResetRoute(*app.Config)
	app.Iface.Close()
	log.Println("TSVPN stopped")
}
