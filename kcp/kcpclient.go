package kcp

import (
	"crypto/sha1"
	"log"
	"time"
	"github.com/golang/snappy"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/cache"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/cipher"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/config"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/counter"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/netutil"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
)

func StartClient(iFace *water.Interface, config config.Config) {
	log.Println("TS kcp client started")
	key := pbkdf2.Key([]byte(config.Key), []byte("default_salt"), 1024, 32, sha1.New)
	block, err := kcp.NewAESBlockCrypt(key)
	if err != nil {
		netutil.PrintErr(err, config.Verbose)
		return
	}
	go tunToKcp(config, iFace)
	for {
		if session, err := kcp.DialWithOptions(config.ServerAddr, block, 10, 3); err == nil {
			cache.GetCache().Set("kcpconn", session, 24*time.Hour)
			kcpToTun(config, session, iFace)
			cache.GetCache().Delete("kcpconn")
		} else {
			netutil.PrintErr(err, config.Verbose)
			time.Sleep(3 * time.Second)
			continue
		}
	}
}

func tunToKcp(config config.Config, iFace *water.Interface) {
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
		if v, ok := cache.GetCache().Get("kcpconn"); ok {
			session := v.(*kcp.UDPSession)
			session.SetWriteDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
			n, err := session.Write(packet[:len(shb)+len(b)])
			if err != nil {
				netutil.PrintErr(err, config.Verbose)
				continue
			}
			counter.IncrWrittenBytes(n)
		}
	}
}

func kcpToTun(config config.Config, session *kcp.UDPSession, iFace *water.Interface) {
	packet := make([]byte, config.BufferSize)
	shb := make([]byte, 2)
	defer session.Close()
	for {
		session.SetReadDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
		n, err := session.Read(shb)
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
			session.SetReadDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
			n, err = session.Read(packet[:shn])
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
				session.SetReadDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
				n, err = session.Read(packet[count : count+receiveSize])
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
