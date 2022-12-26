//go:build darwin
// +build darwin

package water

// MacOSDriverProvider enumerates possible MacOS TUN/TAP implementations
type MacOSDriverProvider int

const (
	// MacOSDriverSystem refers to the default P2P driver
	MacOSDriverSystem MacOSDriverProvider = 0
	// MacOSDriverTunTapOSX refers to the third-party tuntaposx driver
	// see https://sourceforge.net/p/tuntaposx
	MacOSDriverTunTapOSX MacOSDriverProvider = 1
)

// PlatformSpecificParams defines parameters in Config that are specific to
// macOS. A zero-value of such type is valid, yielding an interface
// with OS defined name.
// Currently it is not possible to set the interface name in macOS.
type PlatformSpecificParams struct {
	// Name is the name for the interface to be used.
	// For TunTapOSXDriver, it should be something like "tap0".
	// For SystemDriver, the name should match `utun[0-9]+`, e.g. utun233
	Name string
	// Driver should be set if an alternative driver is desired
	// e.g. TunTapOSXDriver
	Driver MacOSDriverProvider
	// Network is required when creating a TUN interface. The library will call
	// net.ParseCIDR() to parse this string into LocalIP, RemoteNetaddr,
	// RemoteNetmask. The underlying driver will need those to generate ARP
	// response to Windows kernel, to emulate an TUN interface.
	// Please note that it cannot perceive the IP changes caused by DHCP, user
	// configuration to the adapter and etc,. If IP changed, please reconfigure
	// the adapter using syscall, just like openDev().
	// For detail, please refer
	// https://github.com/OpenVPN/tap-windows6/blob/master/src/device.c#L431
	// and https://github.com/songgao/water/pull/13#issuecomment-270341777
	Network []string
}

func defaultPlatformSpecificParams() PlatformSpecificParams {
	return PlatformSpecificParams{}
}
