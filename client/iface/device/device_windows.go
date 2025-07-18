package device

import (
	"fmt"
	"net"
	"net/netip"

	"github.com/amnezia-vpn/amneziawg-go/device"
	"github.com/amnezia-vpn/amneziawg-go/tun"
	"github.com/amnezia-vpn/amneziawg-go/tun/netstack"
	"github.com/amnezia-vpn/amneziawg-windows/tunnel/winipcfg"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"

	"github.com/netbirdio/netbird/client/iface/bind"
	"github.com/netbirdio/netbird/client/iface/configurer"
	"github.com/netbirdio/netbird/client/iface/wgaddr"
)

const defaultWindowsGUIDSTring = "{f2f29e61-d91f-4d76-8151-119b20c4bdeb}"

type TunDevice struct {
	name    string
	address wgaddr.Address
	port    int
	key     string
	mtu     int
	iceBind *bind.ICEBind

	device          *device.Device
	nativeTunDevice *tun.NativeTun
	filteredDevice  *FilteredDevice
	udpMux          *bind.UniversalUDPMuxDefault
	configurer      WGConfigurer
}

func NewTunDevice(name string, address wgaddr.Address, port int, key string, mtu int, iceBind *bind.ICEBind) *TunDevice {
	return &TunDevice{
		name:    name,
		address: address,
		port:    port,
		key:     key,
		mtu:     mtu,
		iceBind: iceBind,
	}
}

func getGUID() (windows.GUID, error) {
	guidString := defaultWindowsGUIDSTring
	if CustomWindowsGUIDString != "" {
		guidString = CustomWindowsGUIDString
	}
	return windows.GUIDFromString(guidString)
}

func (t *TunDevice) Create() (WGConfigurer, error) {
	guid, err := getGUID()
	if err != nil {
		log.Errorf("failed to get GUID: %s", err)
		return nil, err
	}
	log.Info("create tun interface")
	tunDevice, err := tun.CreateTUNWithRequestedGUID(t.name, &guid, t.mtu)
	if err != nil {
		return nil, fmt.Errorf("error creating tun device: %s", err)
	}
	t.nativeTunDevice = tunDevice.(*tun.NativeTun)
	t.filteredDevice = newDeviceFilter(tunDevice)

	// We need to create a wireguard-go device and listen to configuration requests
	t.device = device.NewDevice(
		t.filteredDevice,
		t.iceBind,
		device.NewLogger(wgLogLevel(), "[netbird] "),
	)

	luid := winipcfg.LUID(t.nativeTunDevice.LUID())

	nbiface, err := luid.IPInterface(windows.AF_INET)
	if err != nil {
		t.device.Close()
		return nil, fmt.Errorf("got error when getting ip interface %s", err)
	}

	nbiface.NLMTU = uint32(t.mtu)

	err = nbiface.Set()
	if err != nil {
		t.device.Close()
		return nil, fmt.Errorf("got error when getting setting the interface mtu: %s", err)
	}
	err = t.assignAddr()
	if err != nil {
		t.device.Close()
		return nil, fmt.Errorf("error assigning ip: %s", err)
	}

	t.configurer = configurer.NewUSPConfigurer(t.device, t.name, t.iceBind.ActivityRecorder())
	err = t.configurer.ConfigureInterface(t.key, t.port)
	if err != nil {
		t.device.Close()
		t.configurer.Close()
		return nil, fmt.Errorf("error configuring interface: %s", err)
	}
	return t.configurer, nil
}

func (t *TunDevice) Up() (*bind.UniversalUDPMuxDefault, error) {
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

func (t *TunDevice) UpdateAddr(address wgaddr.Address) error {
	t.address = address
	return t.assignAddr()
}

func (t *TunDevice) Close() error {
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
func (t *TunDevice) WgAddress() wgaddr.Address {
	return t.address
}

func (t *TunDevice) DeviceName() string {
	return t.name
}

func (t *TunDevice) FilteredDevice() *FilteredDevice {
	return t.filteredDevice
}

// Device returns the wireguard device
func (t *TunDevice) Device() *device.Device {
	return t.device
}

func (t *TunDevice) GetInterfaceGUIDString() (string, error) {
	if t.nativeTunDevice == nil {
		return "", fmt.Errorf("interface has not been initialized yet")
	}

	luid := winipcfg.LUID(t.nativeTunDevice.LUID())
	guid, err := luid.GUID()
	if err != nil {
		return "", err
	}
	return guid.String(), nil
}

// assignAddr Adds IP address to the tunnel interface and network route based on the range provided
func (t *TunDevice) assignAddr() error {
	luid := winipcfg.LUID(t.nativeTunDevice.LUID())
	log.Debugf("adding address %s to interface: %s", t.address.IP, t.name)

	prefix := netip.MustParsePrefix(t.address.String())
	ip := net.IP(prefix.Addr().AsSlice())
	mask := net.CIDRMask(prefix.Bits(), prefix.Addr().BitLen())
	return luid.SetIPAddresses([]net.IPNet{
		{IP: ip, Mask: mask},
	})
}

func (t *TunDevice) GetNet() *netstack.Net {
	return nil
}
