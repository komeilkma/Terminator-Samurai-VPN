package UserDatagramProtocol

import (
	"log"
	"net"
	"github.com/golang/snappy"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/cipher"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/config"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/counter"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/netutil"
	"github.com/komeilkma/Terminator-Samurai-VPN/nativewater"
)
func (c *Client) tunToUdp() {
	packet := make([]byte, c.config.BufferSize)
	for {
		n, err := c.iface.Read(packet)
		if err != nil {
			netutil.PrintErr(err, c.config.Verbose)
			break
		}
		b := packet[:n]
		if c.config.Obfs {
			b = cipher.XOR(b)
		}
		if c.config.Compress {
			b = snappy.Encode(nil, b)
		}
		c.localConn.WriteToUDP(b, c.serverAddr)
		counter.IncrWrittenBytes(n)
	}
}
func (c *Client) udpToTun() {
	packet := make([]byte, c.config.BufferSize)
	for {
		n, _, err := c.localConn.ReadFromUDP(packet)
		if err != nil || n == 0 {
			netutil.PrintErr(err, c.config.Verbose)
			continue
		}
		b := packet[:n]
		if c.config.Compress {
			b, err = snappy.Decode(nil, b)
			if err != nil {
				netutil.PrintErr(err, c.config.Verbose)
				continue
			}
		}
		if c.config.Obfs {
			b = cipher.XOR(b)
		}
		c.iface.Write(b)
		counter.IncrReadBytes(n)
	}
}

func StartClient(iface *water.Interface, config config.Config) {
	serverAddr, err := net.ResolveUDPAddr("udp", config.ServerAddr)
	if err != nil {
		log.Fatalln("failed resolve server address:", err)
	}
	localAddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		log.Fatalln("failed get udp socket:", err)
	}
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatalln("failed listen on udp socket:", err)
	}
	defer conn.Close()
	log.Printf("TSVPN udp client started on %v", conn.LocalAddr().String())
	c := &Client{config: config, iface: iface, localConn: conn, serverAddr: serverAddr}
	go c.udpToTun()
	c.tunToUdp()
}
type Client struct {
	config     config.Config
	iface      *water.Interface
	localConn  *net.UDPConn
	serverAddr *net.UDPAddr
}

