package x32

import (
	"github.com/mmmorris1975/simple-logger/logger"
	"net"
	"sync"
	"time"
	"x32-osc-bridge/osc"
)

const (
	OscPort = 10023
)

var Log = logger.StdLogger

type Mixer struct {
	Addr *net.UDPAddr
	Info *osc.XInfo
}

type Client struct {
	Conn *net.UDPConn
	ch   chan bool
	mu   sync.Mutex
}

func NewClient(addr *net.UDPAddr) (*Client, error) {
	var err error
	c := new(Client)

	c.Conn, err = net.DialUDP(addr.Network(), nil, addr)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) StartXremote(t time.Duration) {
	_, _ = c.Conn.Write(osc.WriteString("/xremote"))

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ch != nil {
		Log.Debugln("xremote already active")
		return
	}

	c.ch = make(chan bool)
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
							Log.Errorf("xremote write: %v", err)
							c.StopXremote()
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

func (c *Client) StopXremote() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ch != nil {
		close(c.ch)
		c.ch = nil
	}
}
