package netstack

import (
	"context"
	"net"

	"github.com/amnezia-vpn/amneziawg-go/tun/netstack"
	log "github.com/sirupsen/logrus"
)

type Dialer interface {
	Dial(ctx context.Context, network, addr string) (net.Conn, error)
}

type NSDialer struct {
	net *netstack.Net
}

func NewNSDialer(net *netstack.Net) *NSDialer {
	return &NSDialer{
		net: net,
	}
}

func (d *NSDialer) Dial(ctx context.Context, network, addr string) (net.Conn, error) {
	log.Debugf("dialing %s %s", network, addr)
	conn, err := d.net.Dial(network, addr)
	if err != nil {
		log.Debugf("failed to deal connection: %s", err)
	}
	return conn, err
}
