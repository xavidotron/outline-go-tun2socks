package tun2socks

import (
	"errors"
	"io"
	"runtime/debug"
	"time"

	"github.com/Jigsaw-Code/outline-go-tun2socks/tun2socks"
)

// AppleTunnel embeds the tun2socks.Tunnel interface so it gets exported by gobind.
type AppleTunnel interface {
	tun2socks.Tunnel
}

// TunWriter is an interface that allows for outputting packets to the TUN (VPN).
type TunWriter interface {
	io.WriteCloser
}

var tunnel AppleTunnel

func init() {
	// Apple VPN extensions have a memory limit of 15MB. Conserve memory by increasing garbage
	// collection frequency and returning memory to the OS every minute.
	debug.SetGCPercent(10)
	ticker := time.NewTicker(time.Minute * 1)
	go func() {
		for _ = range ticker.C {
			debug.FreeOSMemory()
		}
	}()
}

// ConnectSocksTunnel reads packets from a TUN device and routes it to a SOCKS server. Returns an
// AppleTunnel instance that should be used to input packets to the tunnel.
//
// `tunWriter` is used to output packets to the TUN (VPN).
// `host` is  IP address of the SOCKS proxy server.
// `port` is the port of the SOCKS proxy server.
// `isUDPEnabled` indicates whether the tunnel and/or network enable UDP proxying.
//
// Sets an error if the tunnel fails to connect.
func ConnectSocksTunnel(tunWriter TunWriter, host string, port int, isUDPEnabled bool) (AppleTunnel, error) {
	if tunWriter == nil || host == "" || port <= 0 {
		return nil, errors.New("Must provide a TunWriter, a valid SOCKS proxy host and port")
	}
	var err error
	tunnel, err = tun2socks.NewTunnel(host, uint16(port), isUDPEnabled, tunWriter)
	return tunnel, err
}
