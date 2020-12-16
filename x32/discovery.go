package x32

import (
	"errors"
	"net"
	"sync"
	"time"
	"x32-osc-bridge/osc"
)

// Discover performs auto-discovery of X32 mixers using OSC messages.  If the addr arg is nil, then broadcast
// discovery is used.  If addr is not nil, then discovery is attempted only for that address.  At most, 1 mixer
// will be returned; if no mixers, or more than 1 mixer is found, an error is returned.
func Discover(addr net.IP) (*Mixer, error) {
	var err error
	var wg sync.WaitGroup
	var addrs []net.IP

	if addr != nil {
		addrs = append(addrs, addr)
	} else {
		addrs, err = discoverAddrs()
		if err != nil {
			return nil, err
		}
	}

	var ch = make(chan *Mixer)
	var wg2 sync.WaitGroup
	var mixers []*Mixer
	wg2.Add(1)
	go func() {
		for m := range ch {
			mixers = append(mixers, m)
			Log.Debugf("found mixer %v", *m.Info)
		}
		wg2.Done()
	}()

	for _, a := range addrs {
		wg.Add(1)
		go func(addr net.IP) {
			defer wg.Done()

			var mixer *Mixer
			mixer, err = discover(addr, 3)
			if err != nil {
				Log.Errorf("discover: %v", err)
				return
			}

			ch <- mixer
		}(a)
	}

	wg.Wait()
	close(ch)
	wg2.Wait()

	switch l := len(mixers); {
	case l > 1:
		return nil, errors.New("more than 1 mixer found")
	case l == 0:
		return nil, errors.New("no mixer found")
	default:
		return mixers[0], nil
	}
}

//nolint:gocognit
func discoverAddrs() ([]net.IP, error) {
	var err error
	var discoveryAddrs []net.IP

	var ifaces []net.Interface
	ifaces, err = net.Interfaces()
	if err != nil {
		return discoveryAddrs, err
	}

	for _, iface := range ifaces {
		// candidate network interfaces must be up and support broadcasting
		if iface.Flags&net.FlagUp != net.FlagUp || iface.Flags&net.FlagBroadcast != net.FlagBroadcast {
			Log.Debugf("not considering interface: %s", iface.Name)
			continue
		}

		var addrs []net.Addr
		addrs, err = iface.Addrs()
		if err != nil {
			Log.Errorf("Addrs: %v", err)
			continue
		}

		for _, addr := range addrs {
			if ip, ok := addr.(*net.IPNet); ok {
				// X32 only uses IPv4
				if ipv4 := ip.IP.To4(); ipv4 != nil {
					discoveryAddrs = append(discoveryAddrs, calculateBroadcast(ip))
				}
			}
		}
	}

	return discoveryAddrs, nil
}

func calculateBroadcast(ipNet *net.IPNet) net.IP {
	inverseMask := make(net.IPMask, len(ipNet.Mask))
	for i, v := range ipNet.Mask {
		inverseMask[i] = v ^ 0xff
	}

	// ipNet.IP will likely be a 16-byte/IPv6 representation (even for IPv4), so use the Mask length to check if we
	// should switch to use IPv4 instead (which should always be the case for communication with an X32 console)
	ip := ipNet.IP
	if len(ipNet.Mask) == net.IPv4len {
		ip = ip.To4()
	}

	bcast := make(net.IP, len(ip))
	for i, v := range ip {
		bcast[i] = v | inverseMask[i]
	}

	Log.Debugf("adding discovery interface: %s", bcast)
	return bcast
}

//nolint:gocognit
func discover(ip net.IP, retries int) (*Mixer, error) {
	var err error

	var conn *net.UDPConn
	conn, err = net.ListenUDP("udp4", new(net.UDPAddr))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if _, err = conn.WriteToUDP(osc.WriteString("/xinfo"), &net.UDPAddr{IP: ip, Port: OscPort}); err != nil {
		if e, ok := err.(*net.OpError); ok {
			if e.Temporary() && retries > 0 {
				Log.Debugf("write error, retrying: %v", err)
				return discover(ip, retries-1)
			}
		}
		return nil, err
	}
	_ = conn.SetReadDeadline(time.Now().Add(300 * time.Millisecond))

	var n int
	var raddr *net.UDPAddr
	var buf = make([]byte, 65535)
	n, raddr, err = conn.ReadFromUDP(buf)
	if err != nil {
		if e, ok := err.(*net.OpError); ok {
			if e.Temporary() && retries > 0 {
				Log.Debugf("read error, retrying: %v", err)
				return discover(ip, retries-1)
			}
		}
		return nil, err
	}

	info := new(osc.XInfo)
	if err = info.UnmarshalBinary(buf[:n]); err != nil {
		return nil, err
	}

	return &Mixer{Addr: raddr, Info: info}, nil
}
