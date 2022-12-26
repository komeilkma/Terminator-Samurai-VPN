package UserDatagramProtocol

import (
	"log"
	"net"
	"time"
	"github.com/golang/snappy"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/cipher"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/config"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/counter"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/netutil"
	"github.com/patrickmn/go-cache"
)
func (s *Server) tunToUdp() {
	packet := make([]byte, s.config.BufferSize)
	for {
		n, err := s.iface.Read(packet)
		if err != nil {
			netutil.PrintErr(err, s.config.Verbose)
			break
		}
		b := packet[:n]
		if key := netutil.GetDstKey(b); key != "" {
			if v, ok := s.connCache.Get(key); ok {
				if s.config.Obfs {
					b = cipher.XOR(b)
				}
				if s.config.Compress {
					b = snappy.Encode(nil, b)
				}
				_, err := s.localConn.WriteToUDP(b, v.(*net.UDPAddr))
				if err != nil {
					s.connCache.Delete(key)
					continue
				}
				counter.IncrWrittenBytes(n)
			}
		}
	}
}

func StartServer(iface *water.Interface, config config.Config) {
	log.Printf("TSVPN udp server started on %v", config.LocalAddr)
	localAddr, err := net.ResolveUDPAddr("udp", config.LocalAddr)
	if err != nil {
		log.Fatalln("failed to get udp socket:", err)
	}
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatalln("failed to listen on udp socket:", err)
	}
	defer conn.Close()
	s := &Server{config: config, iface: iface, localConn: conn, connCache: cache.New(30*time.Minute, 10*time.Minute)}
	go s.tunToUdp()
	s.udpToTun()
}

type Server struct {
	config    config.Config
	iface     *water.Interface
	localConn *net.UDPConn
	connCache *cache.Cache
}
func (s *Server) udpToTun() {
	packet := make([]byte, s.config.BufferSize)
	for {
		n, cliAddr, err := s.localConn.ReadFromUDP(packet)
		if err != nil || n == 0 {
			netutil.PrintErr(err, s.config.Verbose)
			continue
		}
		b := packet[:n]
		if s.config.Compress {
			b, err = snappy.Decode(nil, b)
			if err != nil {
				netutil.PrintErr(err, s.config.Verbose)
				continue
			}
		}
		if s.config.Obfs {
			b = cipher.XOR(b)
		}
		if key := netutil.GetSrcKey(b); key != "" {
			s.iface.Write(b)
			s.connCache.Set(key, cliAddr, 24*time.Hour)
			counter.IncrReadBytes(n)
		}
	}
}
