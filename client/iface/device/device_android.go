//go:build android

package device

import (
	"strings"

	"github.com/amnezia-vpn/amneziawg-go/device"
	"github.com/amnezia-vpn/amneziawg-go/tun"
	"github.com/amnezia-vpn/amneziawg-go/tun/netstack"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"github.com/netbirdio/netbird/client/iface/bind"
	"github.com/netbirdio/netbird/client/iface/configurer"
	"github.com/netbirdio/netbird/client/iface/wgaddr"
)

// WGTunDevice ignore the WGTunDevice interface on Android because the creation of the tun device is different on this platform
type WGTunDevice struct {
	address    wgaddr.Address
	port       int
	key        string
	mtu        int
	iceBind    *bind.ICEBind
	tunAdapter TunAdapter
	disableDNS bool

	name           string
	device         *device.Device
	filteredDevice *FilteredDevice
	udpMux         *bind.UniversalUDPMuxDefault
	configurer     WGConfigurer
}

func NewTunDevice(address wgaddr.Address, port int, key string, mtu int, iceBind *bind.ICEBind, tunAdapter TunAdapter, disableDNS bool) *WGTunDevice {
	return &WGTunDevice{
		address:    address,
		port:       port,
		key:        key,
		mtu:        mtu,
		iceBind:    iceBind,
		tunAdapter: tunAdapter,
		disableDNS: disableDNS,
	}
}

func (t *WGTunDevice) Create(routes []string, dns string, searchDomains []string) (WGConfigurer, error) {
	log.Info("create tun interface")

	routesString := routesToString(routes)
	searchDomainsToString := searchDomainsToString(searchDomains)

	// Skip DNS configuration when DisableDNS is enabled
	if t.disableDNS {
		log.Info("DNS is disabled, skipping DNS and search domain configuration")
		dns = ""
		searchDomainsToString = ""
	}

	fd, err := t.tunAdapter.ConfigureInterface(t.address.String(), t.mtu, dns, searchDomainsToString, routesString)
	if err != nil {
		log.Errorf("failed to create Android interface: %s", err)
		return nil, err
	}

	tunDevice, name, err := tun.CreateUnmonitoredTUNFromFD(fd)
	if err != nil {
		_ = unix.Close(fd)
		log.Errorf("failed to create Android interface: %s", err)
		return nil, err
	}
	t.name = name
	t.filteredDevice = newDeviceFilter(tunDevice)

	log.Debugf("attaching to interface %v", name)
	t.device = device.NewDevice(t.filteredDevice, t.iceBind, device.NewLogger(wgLogLevel(), "[netbird] "))
	// without this property mobile devices can discover remote endpoints if the configured one was wrong.
	// this helps with support for the older NetBird clients that had a hardcoded direct mode
	// t.device.DisableSomeRoamingForBrokenMobileSemantics()

	t.configurer = configurer.NewUSPConfigurer(t.device, t.name, t.iceBind.ActivityRecorder())
	err = t.configurer.ConfigureInterface(t.key, t.port)
	if err != nil {
		t.device.Close()
		t.configurer.Close()
		return nil, err
	}
	return t.configurer, nil
}
func (t *WGTunDevice) Up() (*bind.UniversalUDPMuxDefault, error) {
	err := t.device.Up()
	if err != nil {
		return nil, err
	}

	udpMux, err := t.iceBind.GetICEMux()
	if err != nil {
		return nil, err
	}
	t.udpMux = udpMux
	log.Debugf("device is ready to use: %s", t.name)
	return udpMux, nil
}

func (t *WGTunDevice) UpdateAddr(addr wgaddr.Address) error {
	// todo implement
	return nil
}

func (t *WGTunDevice) Close() error {
	if t.configurer != nil {
		t.configurer.Close()
	}

	if t.device != nil {
		t.device.Close()
		t.device = nil
	}

	if t.udpMux != nil {
		return t.udpMux.Close()

	}
	return nil
}

func (t *WGTunDevice) Device() *device.Device {
	return t.device
}

func (t *WGTunDevice) DeviceName() string {
	return t.name
}

func (t *WGTunDevice) WgAddress() wgaddr.Address {
	return t.address
}

func (t *WGTunDevice) FilteredDevice() *FilteredDevice {
	return t.filteredDevice
}

func (t *WGTunDevice) GetNet() *netstack.Net {
	return nil
}

func routesToString(routes []string) string {
	return strings.Join(routes, ";")
}

func searchDomainsToString(searchDomains []string) string {
	return strings.Join(searchDomains, ";")
}
