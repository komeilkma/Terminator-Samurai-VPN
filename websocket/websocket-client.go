package websocket

import (
	"log"
	"net"
	"time"
	"github.com/gobwas/ws"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/cache"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/cipher"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/config"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/counter"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/netutil"
	"github.com/gobwas/ws/wsutil"
	"github.com/golang/snappy"

)

func StartClient(iface *water.Interface, config config.Config) {
	log.Println("TS websocket client started")
	go tunToWs(config, iface)
	for {
		conn := netutil.ConnectServer(config)
		if conn == nil {
			time.Sleep(3 * time.Second)
			continue
		}
		cache.GetCache().Set("wsconn", conn, 24*time.Hour)
		go wsToTun(config, conn, iface)
		ping(conn, config)
		cache.GetCache().Delete("wsconn")
	}
}

func ping(wsconn net.Conn, config config.Config) {
	defer wsconn.Close()
	for {
		wsconn.SetWriteDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
		err := wsutil.WriteClientMessage(wsconn, ws.OpText, []byte("ping"))
		if err != nil {
			break
		}
		time.Sleep(3 * time.Second)
	}
}

func wsToTun(config config.Config, wsconn net.Conn, iface *water.Interface) {
	defer wsconn.Close()
	for {
		wsconn.SetReadDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
		packet, err := wsutil.ReadServerBinary(wsconn)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			break
		}
		if config.Compress {
			packet, _ = snappy.Decode(nil, packet)
		}
		if config.Obfs {
			packet = cipher.XOR(packet)
		}
		_, err = iface.Write(packet)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			break
		}
		counter.IncrReadBytes(len(packet))
	}
}

func tunToWs(config config.Config, iface *water.Interface) {
	packet := make([]byte, config.BufferSize)
	for {
		n, err := iface.Read(packet)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			break
		}
		if v, ok := cache.GetCache().Get("wsconn"); ok {
			b := packet[:n]
			if config.Obfs {
				b = cipher.XOR(b)
			}
			if config.Compress {
				b = snappy.Encode(nil, b)
			}
			wsconn := v.(net.Conn)
			wsconn.SetWriteDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
			if err = wsutil.WriteClientBinary(wsconn, b); err != nil {
				netutil.PrintErr(err, config.Verbose)
				continue
			}
			counter.IncrWrittenBytes(n)
		}
	}
}
