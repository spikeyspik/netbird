//go:build !android

package iface

import (
	"github.com/amnezia-vpn/amneziawg-go/tun/netstack"

	wgdevice "github.com/amnezia-vpn/amneziawg-go/device"

	"github.com/netbirdio/netbird/client/iface/bind"
	"github.com/netbirdio/netbird/client/iface/device"
	"github.com/netbirdio/netbird/client/iface/wgaddr"
)

type WGTunDevice interface {
	Create() (device.WGConfigurer, error)
	Up() (*bind.UniversalUDPMuxDefault, error)
	UpdateAddr(address wgaddr.Address) error
	WgAddress() wgaddr.Address
	DeviceName() string
	Close() error
	FilteredDevice() *device.FilteredDevice
	Device() *wgdevice.Device
	GetNet() *netstack.Net
}
