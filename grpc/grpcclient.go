package grpc

import (
	"context"
	"crypto/tls"
	"log"
	"time"
	"github.com/golang/snappy"
	"github.com/komeilkma/Terminator-Samurai-VPN/grpc/proto"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/cache"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/cipher"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/config"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/counter"
	"github.com/komeilkma/Terminator-Samurai-VPN/common/netutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"github.com/komeilkma/Terminator-Samurai-VPN/nativewater"
)



func tunToGrpc(config config.Config, iface *water.Interface) {
	packet := make([]byte, config.BufferSize)
	for {
		n, err := iface.Read(packet)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			break
		}
		if v, ok := cache.GetCache().Get("grpcconn"); ok {
			b := packet[:n]
			if config.Obfs {
				b = cipher.XOR(b)
			}
			if config.Compress {
				b = snappy.Encode(nil, b)
			}
			grpcconn := v.(proto.GrpcServe_TunnelClient)
			err = grpcconn.Send(&proto.PacketData{Data: b})
			if err != nil {
				netutil.PrintErr(err, config.Verbose)
				continue
			}
			counter.IncrWrittenBytes(n)
		}
	}
}

func StartClient(iface *water.Interface, config config.Config) {
	log.Println("TSVPN grpc client started")
	go tunToGrpc(config, iface)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: config.TLSInsecureSkipVerify,
	}
	if config.TLSSni != "" {
		tlsconfig.ServerName = config.TLSSni
	}
	creds := credentials.NewTLS(tlsconfig)
	for {
		conn, err := grpc.Dial(config.ServerAddr, grpc.WithBlock(), grpc.WithTransportCredentials(creds))
		if err != nil {
			time.Sleep(3 * time.Second)
			netutil.PrintErr(err, config.Verbose)
			continue
		}
		streamClient := proto.NewGrpcServeClient(conn)
		stream, err := streamClient.Tunnel(context.Background())
		if err != nil {
			conn.Close()
			netutil.PrintErr(err, config.Verbose)
			continue
		}
		cache.GetCache().Set("grpcconn", stream, 24*time.Hour)
		grpcToTun(config, stream, iface)
		cache.GetCache().Delete("grpcconn")
		conn.Close()
	}
}

func grpcToTun(config config.Config, stream proto.GrpcServe_TunnelClient, iface *water.Interface) {
	for {
		packet, err := stream.Recv()
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			break
		}
		b := packet.Data[:]
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
		_, err = iface.Write(b)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			break
		}
		counter.IncrReadBytes(len(b))
	}
}
