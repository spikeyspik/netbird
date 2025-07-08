package bind

import (
	wireguard "github.com/amnezia-vpn/amneziawg-go/conn"

	nbnet "github.com/netbirdio/netbird/util/net"
)

func init() {
	// ControlFns is not thread safe and should only be modified during init.
	*wireguard.ControlFns = append(*wireguard.ControlFns, nbnet.ControlProtectSocket)
}
