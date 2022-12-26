package clientside

import (
	"log"
	"net"
	"strings"
	"time"
	"github.com/patrickmn/go-cache"
)

var _clientside *cache.Cache

func init() {
	_clientside = cache.New(30*time.Minute, 3*time.Minute)
}

func AddClientIP(ip string) {
	_clientside.Add(ip, 0, cache.DefaultExpiration)
}

func DeleteClientIP(ip string) {
	_clientside.Delete(ip)
}

func ExistClientIP(ip string) bool {
	_, ok := _clientside.Get(ip)
	return ok
}

func KeepAliveClientIP(ip string) {
	if ExistClientIP(ip) {
		_clientside.Increment(ip, 1)
	} else {
		AddClientIP(ip)
	}
}

func PickClientIP(cidr string) (clientIP string, prefixLength string) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Panicf("error cidr %v", cidr)
	}
	total := addressCount(ipNet) - 3
	index := uint64(0)
	ip = incr(ipNet.IP.To4())
	for {
		ip = incr(ip)
		index++
		if index >= total {
			break
		}
		if !ExistClientIP(ip.String()) {
			AddClientIP(ip.String())
			return ip.String(), strings.Split(cidr, "/")[1]
		}
	}
	return "", ""
}

func ListClientIPs() []string {
	result := []string{}
	for k := range _clientside.Items() {
		result = append(result, k)
	}
	return result
}

func addressCount(network *net.IPNet) uint64 {
	prefixLen, bits := network.Mask.Size()
	return 1 << (uint64(bits) - uint64(prefixLen))
}

func incr(IP net.IP) net.IP {
	IP = checkIPv4(IP)
	incIP := make([]byte, len(IP))
	copy(incIP, IP)
	for j := len(incIP) - 1; j >= 0; j-- {
		incIP[j]++
		if incIP[j] > 0 {
			break
		}
	}
	return incIP
}

func checkIPv4(ip net.IP) net.IP {
	if v4 := ip.To4(); v4 != nil {
		return v4
	}
	return ip
}
