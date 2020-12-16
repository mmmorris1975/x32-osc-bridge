package x32

import (
	"github.com/mmmorris1975/simple-logger/logger"
	"net"
	"sync"
	"time"
	"x32-osc-bridge/osc"
)

const (
	// OscPort is the port used by the X32 mixer to receive OSC messages
	OscPort = 10023
)

// Log is the logger used for this package
var Log = logger.StdLogger

// Mixer is the IP address and port, and other information, found during the discovery process
type Mixer struct {
	Addr *net.UDPAddr
	Info *osc.XInfo
}

// Client represents a connection through the bridge/proxy to the mixer on the back end
type Client struct {
	Conn *net.UDPConn
	ch   chan bool
	mu   sync.Mutex
}

// NewClient creates a new connection to the mixer specified in addr
func NewClient(addr *net.UDPAddr) (*Client, error) {
	var err error
	c := new(Client)

	c.Conn, err = net.DialUDP(addr.Network(), nil, addr)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// StartXremote sends the /xremote OSC command to the mixer so that it will send messages when mixer parameters change
//nolint:gocognit
func (c *Client) StartXremote(t time.Duration) {
	_, _ = c.Conn.Write(osc.WriteString("/xremote"))

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ch != nil {
		Log.Debugln("xremote already active")
		return
	}

	c.ch = make(chan bool, 1)
	go func() {
		ticker := time.NewTicker(t)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if _, err := c.Conn.Write(osc.WriteString("/xremote")); err != nil {
					if e, ok := err.(*net.OpError); ok {
						if !e.Temporary() {
							// conn is likely closed, so we should terminate the loop
							Log.Debugf("xremote write: %v", err)
							c.StopXremote()
							return
						}
					}
				}
			case <-c.ch:
				Log.Debugln("stopping xremote")
				return
			}
		}
	}()
}

// StopXremote shuts down the loop which manages the /xremote data subscription
func (c *Client) StopXremote() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ch != nil {
		close(c.ch)
		c.ch = nil
	}
}
