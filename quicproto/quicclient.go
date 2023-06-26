package quicproto

import (
	"context"
	"crypto/tls"
	"log"
	"time"
	"github.com/golang/snappy"
	"github.com/lucas-clemente/quic-go"
	"github.com/komeilkma/Terminator-Samurai-VPN/nativewater"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/cache"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/cipher"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/config"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/counter"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/netutil"
)

func StartClient(iFace *water.Interface, config config.Config) {
	log.Println("TSVPN quic client started")
	tlsConfig := &tls.Config{
		InsecureSkipVerify: config.TLSInsecureSkipVerify,
		NextProtos:         []string{"TSVPN"},
	}
	if config.TLSSni != "" {
		tlsConfig.ServerName = config.TLSSni
	}
	go tunToQuic(config, iFace)
	for {
		conn, err := quic.DialAddr(config.ServerAddr, tlsConfig, nil)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			time.Sleep(3 * time.Second)
			continue
		}
		stream, err := conn.OpenStreamSync(context.Background())
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			conn.CloseWithError(quic.ApplicationErrorCode(0x01), "closed")
			continue
		}
		cache.GetCache().Set("quicstream", stream, 24*time.Hour)
		quicToTun(config, stream, iFace)
		cache.GetCache().Delete("quicstream")
	}
}

func tunToQuic(config config.Config, iFace *water.Interface) {
	packet := make([]byte, config.BufferSize)
	shb := make([]byte, 2)
	for {
		shn, err := iFace.Read(packet)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			continue
		}
		b := packet[:shn]
		if config.Obfs {
			b = cipher.XOR(b)
		}
		if config.Compress {
			b = snappy.Encode(nil, b)
		}
		shb[0] = byte(shn >> 8 & 0xff)
		shb[1] = byte(shn & 0xff)
		copy(packet[len(shb):len(shb)+len(b)], b)
		copy(packet[:len(shb)], shb)
		if v, ok := cache.GetCache().Get("quicstream"); ok {
			stream := v.(quic.Stream)
			stream.SetWriteDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
			n, err := stream.Write(packet[:len(shb)+len(b)])
			if err != nil {
				netutil.PrintErr(err, config.Verbose)
				continue
			}
			counter.IncrWrittenBytes(n)
		}
	}
}

func quicToTun(config config.Config, stream quic.Stream, iFace *water.Interface) {
	defer stream.Close()
	packet := make([]byte, config.BufferSize)
	shb := make([]byte, 2)
	for {
		stream.SetReadDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
		n, err := stream.Read(shb)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			break
		}
		if n < 2 {
			break
		}
		shn := 0
		shn = ((shn & 0x00) | int(shb[0])) << 8
		shn = shn | int(shb[1])
		splitSize := 99
		var count = 0
		if shn < splitSize {
			stream.SetReadDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
			n, err = stream.Read(packet[:shn])
			if err != nil {
				netutil.PrintErr(err, config.Verbose)
				break
			}
			count = n
		} else {
			for count < shn {
				receiveSize := splitSize
				if shn-count < splitSize {
					receiveSize = shn - count
				}
				stream.SetReadDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
				n, err = stream.Read(packet[count : count+receiveSize])
				if err != nil {
					netutil.PrintErr(err, config.Verbose)
					break
				}
				count += n
			}
		}
		b := packet[:shn]
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
		n, err = iFace.Write(b)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			break
		}
		counter.IncrReadBytes(n)
	}
}
