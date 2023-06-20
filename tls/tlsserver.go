package tls

import (
	"crypto/tls"
	"log"
	"net"
	"time"
	"github.com/golang/snappy"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/cache"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/cipher"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/config"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/counter"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/netutil"
	"github.com/komeilkma/Terminator-Samurai-VPN/nativewater"
)
func StartServer(iFace *water.Interface, config config.Config) {
	log.Printf("TSVPN tls server started on %v", config.LocalAddr)
	cert, err := tls.LoadX509KeyPair(config.TLSCertificateFilePath, config.TLSCertificateKeyFilePath)
	if err != nil {
		log.Panic(err)
	}
	tlsConfig := &tls.Config{
		Certificates:     []tls.Certificate{cert},
		MinVersion:       tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		},
	}
	ln, err := tls.Listen("tcp", config.LocalAddr, tlsConfig)
	if err != nil {
		log.Panic(err)
	}
	go toClient(config, iFace)
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		sniffConn := NewPeekPreDataConn(conn)
		switch sniffConn.Type {
		case TypeHttp:
			if sniffConn.Handle() {
				continue
			}
		case TypeHttp2:
			if sniffConn.Handle() {
				continue
			}
		}
		go toServer(config, sniffConn, iFace)
	}
}
func toClient(config config.Config, iface *water.Interface) {
	packet := make([]byte, config.BufferSize)
	for {
		n, err := iface.Read(packet)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			continue
		}
		b := packet[:n]
		if key := netutil.GetDstKey(b); key != "" {
			if v, ok := cache.GetCache().Get(key); ok {
				if config.Obfs {
					b = cipher.XOR(b)
				}
				if config.Compress {
					b = snappy.Encode(nil, b)
				}
				v.(net.Conn).SetWriteDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
				_, err := v.(net.Conn).Write(b)
				if err != nil {
					cache.GetCache().Delete(key)
					continue
				}
				counter.IncrWrittenBytes(n)
			}
		}
	}
}
func toServer(config config.Config, tlsconn net.Conn, iface *water.Interface) {
	defer tlsconn.Close()
	packet := make([]byte, config.BufferSize)
	for {
		tlsconn.SetReadDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
		n, err := tlsconn.Read(packet)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			break
		}
		b := packet[:n]
		if config.Compress {
			b, err = snappy.Decode(nil, b)
			if err != nil {
				netutil.PrintErr(err, config.Verbose)
				break
			}
		}
		if config.Obfs {
			b = cipher.XOR(b)
		}
		if key := netutil.GetSrcKey(b); key != "" {
			cache.GetCache().Set(key, tlsconn, 24*time.Hour)
			iface.Write(b)
			counter.IncrReadBytes(n)
		}
	}
}
