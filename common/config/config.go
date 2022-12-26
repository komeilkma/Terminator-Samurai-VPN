package config

type Config struct {
	DeviceName                string
	LocalAddr                 string
	ServerAddr                string
	ServerIP                  string
	ServerIPv6                string
	DNSIP                     string
	CIDR                      string
	CIDRv6                    string
	Key                       string
	Protocol                  string
	WebSocketPath             string
	ServerMode                bool
	GlobalMode                bool
	Obfs                      bool
	Compress                  bool
	MTU                       int
	Timeout                   int
	LocalGateway              string
	LocalGatewayv6            string
	TLSCertificateFilePath    string
	TLSCertificateKeyFilePath string
	TLSSni                    string
	TLSInsecureSkipVerify     bool
	BufferSize                int
	Verbose                   bool
}
